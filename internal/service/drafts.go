package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/google/uuid"
)

type TaskDraft struct {
	ID, ProjectID        string
	Input                core.TaskInput
	Revision             int
	Saved, Ready         bool
	MissingFields        []string
	CreatedAt, UpdatedAt string
}

type DraftSave struct {
	Input    core.TaskInput
	Revision int
}

func describeDraft(d *TaskDraft) {
	d.MissingFields = []string{}
	if strings.TrimSpace(d.Input.Title) == "" {
		d.MissingFields = append(d.MissingFields, "Title")
	}
	if d.Input.State == "" && d.Input.BoardColumn == "" {
		d.MissingFields = append(d.MissingFields, "State")
	}
	if d.Input.Assignee.Type == "" || d.Input.Assignee.Type == "unassigned" || (d.Input.Assignee.Type == "worker" && (d.Input.Assignee.WorkerID == nil || *d.Input.Assignee.WorkerID == "")) {
		d.MissingFields = append(d.MissingFields, "Responsible")
	}
	if strings.TrimSpace(d.Input.Description) == "" {
		d.MissingFields = append(d.MissingFields, "Description")
	}
	d.Ready = len(d.MissingFields) == 0
}

func meaningfulDraft(in core.TaskInput) bool {
	if strings.TrimSpace(in.Title+in.Description+in.AITestInstructions+in.HumanTestInstructions) != "" || in.State != "" || in.BoardColumn != "" || (in.Priority != "" && in.Priority != "medium") || (in.TestingMode != "" && in.TestingMode != "ai") || (in.Assignee.Type != "" && in.Assignee.Type != "unassigned") || len(in.DependencyIDs) > 0 {
		return true
	}
	for _, value := range in.Properties {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}

func (s *Service) ListTaskDrafts(ctx context.Context, projectID string) ([]TaskDraft, error) {
	if _, err := s.Store.GetProject(ctx, projectID); err != nil {
		return nil, err
	}
	rows, err := s.Store.DB.QueryContext(ctx, `SELECT id,input_json,revision,created_at,updated_at FROM task_drafts WHERE project_id=? AND discarded=0 AND published_task_id IS NULL ORDER BY updated_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []TaskDraft{}
	for rows.Next() {
		d := TaskDraft{ProjectID: projectID, Saved: true}
		var raw string
		if err := rows.Scan(&d.ID, &raw, &d.Revision, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(raw), &d.Input); err != nil {
			return nil, err
		}
		describeDraft(&d)
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *Service) SaveTaskDraft(ctx context.Context, projectID, draftID string, in DraftSave) (TaskDraft, error) {
	d := TaskDraft{ID: draftID, ProjectID: projectID, Input: in.Input, Revision: in.Revision}
	if _, err := uuid.Parse(draftID); err != nil || in.Revision < 0 {
		return d, core.ErrInvalidInput
	}
	if _, err := s.Store.GetProject(ctx, projectID); err != nil {
		return d, err
	}
	describeDraft(&d)
	raw, err := json.Marshal(in.Input)
	if err != nil {
		return d, err
	}
	if len(raw) > 1<<20 {
		return d, core.ErrInvalidInput
	}
	now := core.Now()
	meaningful := meaningfulDraft(in.Input)
	if in.Revision == 0 {
		if !meaningful {
			return d, nil
		}
		result, err := s.Store.DB.ExecContext(ctx, `INSERT INTO task_drafts(id,project_id,input_json,revision,created_at,updated_at) VALUES(?,?,?,1,?,?) ON CONFLICT(id) DO NOTHING`, draftID, projectID, string(raw), now, now)
		if err != nil {
			return d, err
		}
		count, _ := result.RowsAffected()
		if count != 1 {
			return d, core.ErrDraftConflict
		}
		d.Revision, d.CreatedAt = 1, now
	} else {
		result, err := s.Store.DB.ExecContext(ctx, `UPDATE task_drafts SET input_json=?,revision=revision+1,discarded=?,updated_at=? WHERE id=? AND project_id=? AND revision=? AND discarded=0 AND published_task_id IS NULL`, string(raw), !meaningful, now, draftID, projectID, in.Revision)
		if err != nil {
			return d, err
		}
		count, _ := result.RowsAffected()
		if count != 1 {
			return d, core.ErrDraftConflict
		}
		d.Revision++
	}
	d.Saved, d.UpdatedAt = meaningful, now
	return d, nil
}

func (s *Service) DiscardTaskDraft(ctx context.Context, projectID, draftID string) error {
	_, err := s.Store.DB.ExecContext(ctx, `UPDATE task_drafts SET discarded=1,revision=revision+1,updated_at=? WHERE id=? AND project_id=? AND published_task_id IS NULL`, core.Now(), draftID, projectID)
	return err
}

// PublishTaskDraft atomically creates one task and retires its draft. Repeated
// publish requests return the same task, including after a lost HTTP response.
func (s *Service) PublishTaskDraft(ctx context.Context, projectID, draftID string, revision int) (core.Task, error) {
	var raw string
	var published sql.NullString
	var storedRevision, discarded int
	err := s.Store.DB.QueryRowContext(ctx, `SELECT input_json,revision,published_task_id,discarded FROM task_drafts WHERE id=? AND project_id=?`, draftID, projectID).Scan(&raw, &storedRevision, &published, &discarded)
	if err == sql.ErrNoRows {
		return core.Task{}, core.ErrNotFound
	}
	if err != nil {
		return core.Task{}, err
	}
	if published.Valid {
		return s.Store.GetTask(ctx, published.String)
	}
	if discarded != 0 || storedRevision != revision {
		return core.Task{}, core.ErrDraftConflict
	}
	var input core.TaskInput
	if err := json.Unmarshal([]byte(raw), &input); err != nil {
		return core.Task{}, err
	}
	d := TaskDraft{Input: input}
	describeDraft(&d)
	if !d.Ready {
		return core.Task{}, core.ErrDraftIncomplete
	}
	input, err = s.prepareHumanTask(ctx, projectID, input)
	if err != nil {
		return core.Task{}, err
	}
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return core.Task{}, err
	}
	defer tx.Rollback()
	// Take the write lock before rechecking revision/publication.
	if _, err := tx.ExecContext(ctx, `UPDATE task_drafts SET revision=revision WHERE id=?`, draftID); err != nil {
		return core.Task{}, err
	}
	err = tx.QueryRowContext(ctx, `SELECT revision,published_task_id,discarded FROM task_drafts WHERE id=? AND project_id=?`, draftID, projectID).Scan(&storedRevision, &published, &discarded)
	if err != nil {
		return core.Task{}, err
	}
	if published.Valid {
		_ = tx.Rollback()
		return s.Store.GetTask(ctx, published.String)
	}
	if discarded != 0 || storedRevision != revision {
		return core.Task{}, core.ErrDraftConflict
	}
	task, err := s.Store.CreateTaskTx(ctx, tx, projectID, input)
	if err != nil {
		return task, err
	}
	_, err = tx.ExecContext(ctx, `UPDATE task_drafts SET published_task_id=?,discarded=1,revision=revision+1,updated_at=? WHERE id=?`, task.ID, core.Now(), draftID)
	if err != nil {
		return task, err
	}
	if err := tx.Commit(); err != nil {
		return task, err
	}
	return s.Store.GetTask(ctx, task.ID)
}
