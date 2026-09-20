package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/google/uuid"
)

type Service struct{ Store *persistence.Store }

func New(store *persistence.Store) *Service { return &Service{Store: store} }
func id() string                            { return uuid.NewString() }

func (s *Service) BeginAgentSession(ctx context.Context, projectPath, workerRef, client string) (core.AgentContext, error) {
	p, err := s.Store.GetProjectByPath(ctx, projectPath)
	if err != nil {
		return core.AgentContext{}, err
	}
	w, err := s.Store.FindWorker(ctx, p.ID, workerRef)
	if err != nil {
		return core.AgentContext{}, err
	}
	if !w.Enabled || w.Archived {
		return core.AgentContext{}, core.ErrWorkerDisabled
	}
	sess, err := s.Store.StartSession(ctx, p, w, client)
	if err != nil {
		return core.AgentContext{}, err
	}
	return core.AgentContext{Project: p, Worker: w, Session: sess}, nil
}
func (s *Service) EndAgentSession(ctx context.Context, a core.AgentContext) error {
	return s.Store.EndSession(ctx, a.Session.ID)
}

func (s *Service) AgentBoard(ctx context.Context, a core.AgentContext) (core.Board, error) {
	rows, err := s.Store.DB.QueryContext(ctx, persistenceTaskSelect()+` WHERE t.project_id=? AND t.assignee_type='worker' AND t.assignee_worker_id=? AND t.state IN ('features','in_progress','testing','verification') ORDER BY t.state,t.position`, a.Project.ID, a.Worker.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	b := core.Board{"features": {}, "in_progress": {}, "testing": {}, "verification": {}}
	var tasks []core.Task
	for rows.Next() {
		t, e := scanTask(rows)
		if e != nil {
			return nil, e
		}
		tasks = append(tasks, t)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range tasks {
		if err := s.loadAgentContext(ctx, &tasks[i]); err != nil {
			return nil, err
		}
		b[tasks[i].State] = append(b[tasks[i].State], tasks[i])
	}
	return b, rows.Err()
}

func (s *Service) AgentTask(ctx context.Context, a core.AgentContext, taskID string) (core.TaskDetails, error) {
	t, err := s.accessibleTask(ctx, a, taskID)
	if err != nil {
		return core.TaskDetails{}, err
	}
	if err = s.loadAgentContext(ctx, &t); err != nil {
		return core.TaskDetails{}, err
	}
	d := core.TaskDetails{Task: t}
	d.History, err = s.ListHistory(ctx, taskID)
	if err != nil {
		return d, err
	}
	d.TestRuns, err = s.ListTests(ctx, taskID)
	if err != nil {
		return d, err
	}
	d.Artifacts = t.Artifacts
	d.SpawnedTasks, err = s.Store.ListSpawned(ctx, taskID)
	if err != nil {
		return d, err
	}
	d.Usage, err = s.ListUsageForWorker(ctx, taskID, a.Worker.ID)
	return d, err
}

func (s *Service) accessibleTask(ctx context.Context, a core.AgentContext, taskID string) (core.Task, error) {
	t, err := s.Store.GetTask(ctx, taskID)
	if err != nil {
		return t, core.ErrTaskNotAccessible
	}
	if t.ProjectID != a.Project.ID || t.AssigneeType != "worker" || t.AssigneeWorkerID == nil || *t.AssigneeWorkerID != a.Worker.ID || t.State == "backlog" || t.State == "complete" {
		return core.Task{}, core.ErrTaskNotAccessible
	}
	return t, nil
}

func (s *Service) CreateFeature(ctx context.Context, a core.AgentContext, in core.FeatureInput) (core.Task, error) {
	if strings.TrimSpace(in.Title) == "" {
		return core.Task{}, core.ErrInvalidInput
	}
	if in.Priority == "" {
		in.Priority = "medium"
	}
	if in.TestingMode == "" {
		in.TestingMode = "ai"
	}
	if in.Assignee.Type == "" {
		in.Assignee.Type = "unassigned"
	}
	if err := s.validateAssignee(ctx, a.Project.ID, in.Assignee, true); err != nil {
		return core.Task{}, err
	}
	if in.SourceTaskID != nil {
		if _, err := s.accessibleOrCreated(ctx, a, *in.SourceTaskID); err != nil {
			return core.Task{}, core.ErrTaskNotAccessible
		}
	}
	for _, dep := range in.DependsOn {
		if _, err := s.accessibleOrCreated(ctx, a, dep); err != nil {
			return core.Task{}, core.ErrTaskNotAccessible
		}
	}
	now := core.Now()
	t := core.Task{ID: id(), ProjectID: a.Project.ID, Title: in.Title, Description: in.Description, State: "features", Position: 1000, Priority: in.Priority, AssigneeType: in.Assignee.Type, AssigneeWorkerID: in.Assignee.WorkerID, CreatedByType: "agent", CreatedByWorkerID: &a.Worker.ID, CreatedBySessionID: &a.Session.ID, SourceTaskID: in.SourceTaskID, TestingMode: in.TestingMode, CreatedAt: now, UpdatedAt: now}
	assigned := in.Assignee.Type
	if in.Assignee.WorkerID != nil {
		if w, e := s.Store.GetWorker(ctx, *in.Assignee.WorkerID); e == nil {
			assigned = w.Name
		}
	}
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return t, err
	}
	defer tx.Rollback()
	_ = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(position),0)+1000 FROM tasks WHERE project_id=? AND state='features'`, a.Project.ID).Scan(&t.Position)
	_, err = tx.ExecContext(ctx, `INSERT INTO tasks(id,project_id,title,description,state,position,priority,assignee_type,assignee_worker_id,created_by_type,created_by_worker_id,created_by_session_id,source_task_id,testing_mode,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, t.ID, t.ProjectID, t.Title, t.Description, t.State, t.Position, t.Priority, t.AssigneeType, t.AssigneeWorkerID, t.CreatedByType, t.CreatedByWorkerID, t.CreatedBySessionID, t.SourceTaskID, t.TestingMode, now, now)
	if err != nil {
		return t, err
	}
	msg := fmt.Sprintf("Feature created by %s\nAssigned to: %s", a.Worker.Name, assigned)
	if in.Reason != "" {
		msg += "\nReason: " + in.Reason
	}
	if in.SourceTaskID != nil {
		msg += "\nSource: " + *in.SourceTaskID
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, id(), t.ID, "system", "task_created", msg, now)
	if err != nil {
		return t, err
	}
	if in.SourceTaskID != nil {
		_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, id(), *in.SourceTaskID, "system", "system_event", "Spawned feature "+t.ID, now)
		if err != nil {
			return t, err
		}
	}
	for _, dep := range in.DependsOn {
		if _, err = tx.ExecContext(ctx, `INSERT INTO task_dependencies(task_id,depends_on_task_id,created_at) VALUES(?,?,?)`, t.ID, dep, now); err != nil {
			return t, err
		}
	}
	if err = tx.Commit(); err != nil {
		return t, err
	}
	return s.Store.GetTask(ctx, t.ID)
}

func (s *Service) accessibleOrCreated(ctx context.Context, a core.AgentContext, taskID string) (core.Task, error) {
	t, err := s.Store.GetTask(ctx, taskID)
	if err != nil || t.ProjectID != a.Project.ID {
		return core.Task{}, core.ErrTaskNotAccessible
	}
	if t.AssigneeWorkerID != nil && *t.AssigneeWorkerID == a.Worker.ID {
		return t, nil
	}
	if t.CreatedByWorkerID != nil && *t.CreatedByWorkerID == a.Worker.ID {
		return t, nil
	}
	return core.Task{}, core.ErrTaskNotAccessible
}
func (s *Service) validateAssignee(ctx context.Context, projectID string, a core.Assignee, enabled bool) error {
	switch a.Type {
	case "human", "unassigned":
		if a.WorkerID != nil {
			return core.ErrInvalidInput
		}
	case "worker":
		if a.WorkerID == nil {
			return core.ErrInvalidInput
		}
		w, err := s.Store.GetWorker(ctx, *a.WorkerID)
		if err != nil || w.ProjectID != projectID || w.Archived || enabled && !w.Enabled {
			return core.ErrInvalidInput
		}
	default:
		return core.ErrInvalidInput
	}
	return nil
}

func (s *Service) AgentAddHistory(ctx context.Context, a core.AgentContext, taskID, content string) (core.HistoryEntry, error) {
	if _, err := s.accessibleTask(ctx, a, taskID); err != nil {
		return core.HistoryEntry{}, err
	}
	if strings.TrimSpace(content) == "" {
		return core.HistoryEntry{}, core.ErrInvalidInput
	}
	h := core.HistoryEntry{ID: id(), TaskID: taskID, ActorType: "agent", EntryType: "agent_comment", Content: content, CreatedAt: core.Now(), WorkerID: &a.Worker.ID, WorkerSessionID: &a.Session.ID, ActorName: a.Worker.Name}
	_, err := s.Store.DB.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,worker_id,worker_session_id,entry_type,content,created_at) VALUES(?,?,?,?,?,?,?,?)`, h.ID, h.TaskID, h.ActorType, h.WorkerID, h.WorkerSessionID, h.EntryType, h.Content, h.CreatedAt)
	return h, err
}

func (s *Service) AgentMove(ctx context.Context, a core.AgentContext, taskID, target string) (core.Task, error) {
	t, err := s.accessibleTask(ctx, a, taskID)
	if err != nil {
		return t, err
	}
	allowed := map[string]string{"features": "in_progress", "in_progress": "testing", "testing": "verification"}
	if allowed[t.State] != target {
		return t, core.ErrInvalidTransition
	}
	if target == "verification" {
		if err = s.canAgentVerify(ctx, t, a.Worker.ID); err != nil {
			return t, err
		}
	}
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return t, err
	}
	defer tx.Rollback()
	now := core.Now()
	r, err := tx.ExecContext(ctx, `UPDATE tasks SET state=?,updated_at=? WHERE id=? AND state=? AND assignee_worker_id=?`, target, now, t.ID, t.State, a.Worker.ID)
	if err != nil {
		return t, err
	}
	n, _ := r.RowsAffected()
	if n != 1 {
		return t, core.ErrInvalidTransition
	}
	msg := fmt.Sprintf("State changed: %s → %s", t.State, target)
	_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,worker_id,worker_session_id,entry_type,content,created_at) VALUES(?,?,?,?,?,?,?,?)`, id(), t.ID, "system", a.Worker.ID, a.Session.ID, "state_transition", msg, now)
	if err != nil {
		return t, err
	}
	if err = tx.Commit(); err != nil {
		return t, err
	}
	return s.Store.GetTask(ctx, t.ID)
}

func (s *Service) canAgentVerify(ctx context.Context, t core.Task, workerID string) error {
	if t.TestingMode == "human" {
		return core.ErrHumanTestPending
	}
	var aiPass int
	_ = s.Store.DB.QueryRowContext(ctx, `SELECT count(*) FROM test_runs WHERE task_id=? AND worker_id=? AND status='passed'`, t.ID, workerID).Scan(&aiPass)
	if aiPass == 0 {
		return core.ErrInvalidTransition
	}
	if t.TestingMode == "hybrid" {
		var humanPass int
		_ = s.Store.DB.QueryRowContext(ctx, `SELECT count(*) FROM test_runs WHERE task_id=? AND worker_id IS NULL AND status='passed'`, t.ID).Scan(&humanPass)
		if humanPass == 0 {
			return core.ErrHumanTestPending
		}
	}
	return nil
}

func (s *Service) HumanAssign(ctx context.Context, taskID string, a core.Assignee) (core.Task, error) {
	t, err := s.Store.GetTask(ctx, taskID)
	if err != nil {
		return t, err
	}
	if err = s.validateAssignee(ctx, t.ProjectID, a, false); err != nil {
		return t, err
	}
	old := assigneeLabel(ctx, s.Store, t.AssigneeType, t.AssigneeWorkerID)
	next := assigneeLabel(ctx, s.Store, a.Type, a.WorkerID)
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return t, err
	}
	defer tx.Rollback()
	now := core.Now()
	_, err = tx.ExecContext(ctx, `UPDATE tasks SET assignee_type=?,assignee_worker_id=?,updated_at=? WHERE id=?`, a.Type, a.WorkerID, now, taskID)
	if err != nil {
		return t, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, id(), taskID, "system", "assignment_changed", fmt.Sprintf("Assignment changed: %s → %s", old, next), now)
	if err != nil {
		return t, err
	}
	if err = tx.Commit(); err != nil {
		return t, err
	}
	return s.Store.GetTask(ctx, taskID)
}
func assigneeLabel(ctx context.Context, s *persistence.Store, typ string, wid *string) string {
	if typ != "worker" || wid == nil {
		return strings.Title(typ)
	}
	w, e := s.GetWorker(ctx, *wid)
	if e != nil {
		return "Archived worker"
	}
	return w.Name
}

func (s *Service) HumanMove(ctx context.Context, taskID, target string) (core.Task, error) {
	t, err := s.Store.GetTask(ctx, taskID)
	if err != nil {
		return t, err
	}
	valid := map[string]bool{"backlog": true, "features": true, "in_progress": true, "testing": true, "verification": true, "complete": true}
	if !valid[target] {
		return t, core.ErrInvalidTransition
	}
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return t, err
	}
	defer tx.Rollback()
	now := core.Now()
	_, err = tx.ExecContext(ctx, `UPDATE tasks SET state=?,updated_at=? WHERE id=?`, target, now, taskID)
	if err != nil {
		return t, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, id(), taskID, "system", "state_transition", fmt.Sprintf("State changed by Human: %s → %s", t.State, target), now)
	if err != nil {
		return t, err
	}
	if err = tx.Commit(); err != nil {
		return t, err
	}
	return s.Store.GetTask(ctx, taskID)
}

func (s *Service) HumanComment(ctx context.Context, taskID, content string) (core.HistoryEntry, error) {
	if _, err := s.Store.GetTask(ctx, taskID); err != nil {
		return core.HistoryEntry{}, err
	}
	h := core.HistoryEntry{ID: id(), TaskID: taskID, ActorType: "human", EntryType: "human_comment", Content: content, CreatedAt: core.Now(), ActorName: "Human"}
	_, err := s.Store.DB.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, h.ID, taskID, h.ActorType, h.EntryType, h.Content, h.CreatedAt)
	return h, err
}

func (s *Service) ListHistory(ctx context.Context, taskID string) ([]core.HistoryEntry, error) {
	rows, err := s.Store.DB.QueryContext(ctx, `SELECT h.id,h.task_id,h.actor_type,h.worker_id,h.worker_session_id,h.entry_type,h.content,h.created_at,COALESCE(w.name,'') FROM history_entries h LEFT JOIN workers w ON w.id=h.worker_id WHERE h.task_id=? ORDER BY h.created_at`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.HistoryEntry
	for rows.Next() {
		var h core.HistoryEntry
		if err = rows.Scan(&h.ID, &h.TaskID, &h.ActorType, &h.WorkerID, &h.WorkerSessionID, &h.EntryType, &h.Content, &h.CreatedAt, &h.ActorName); err != nil {
			return nil, err
		}
		if h.ActorName == "" {
			h.ActorName = strings.Title(h.ActorType)
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

func (s *Service) loadAgentContext(ctx context.Context, t *core.Task) error {
	var err error
	t.Dependencies, err = s.ListDependencies(ctx, t.ID)
	if err != nil {
		return err
	}
	t.Artifacts, err = s.ListArtifacts(ctx, t.ID)
	if err != nil {
		return err
	}
	t.Properties, err = s.ListTaskProperties(ctx, t.ID, true)
	return err
}

func (s *Service) ListDependencies(ctx context.Context, taskID string) ([]core.Dependency, error) {
	rows, err := s.Store.DB.QueryContext(ctx, `SELECT d.task_id,d.depends_on_task_id,t.title,t.state FROM task_dependencies d JOIN tasks t ON t.id=d.depends_on_task_id WHERE d.task_id=?`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.Dependency
	for rows.Next() {
		var d core.Dependency
		if err = rows.Scan(&d.TaskID, &d.DependsOnTaskID, &d.Title, &d.State); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Artifacts, _ = s.ListArtifacts(ctx, out[i].DependsOnTaskID)
	}
	return out, rows.Err()
}

// Kept local so workflow queries stay aligned with persistence scans without exposing SQL internals.
func persistenceTaskSelect() string {
	return `SELECT t.id,t.project_id,t.title,t.description,t.state,t.position,t.priority,t.assignee_type,t.assignee_worker_id,t.created_by_type,t.created_by_worker_id,t.created_by_session_id,t.source_task_id,t.testing_mode,t.ai_test_instructions,t.human_test_instructions,t.created_at,t.updated_at,COALESCE(aw.name,''),COALESCE(cw.name,'') FROM tasks t LEFT JOIN workers aw ON aw.id=t.assignee_worker_id LEFT JOIN workers cw ON cw.id=t.created_by_worker_id`
}
func scanTask(row interface{ Scan(...any) error }) (core.Task, error) {
	var t core.Task
	err := row.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.State, &t.Position, &t.Priority, &t.AssigneeType, &t.AssigneeWorkerID, &t.CreatedByType, &t.CreatedByWorkerID, &t.CreatedBySessionID, &t.SourceTaskID, &t.TestingMode, &t.AITestInstructions, &t.HumanTestInstructions, &t.CreatedAt, &t.UpdatedAt, &t.AssigneeName, &t.CreatorName)
	return t, err
}

var _ = errors.Is
var _ = json.Valid
var _ = sql.ErrNoRows
