package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

//go:embed migrations/001_initial.sql
var initialMigration string

type Store struct{ DB *sql.DB }

func DefaultDBPath() string {
	base, err := os.UserConfigDir()
	if err != nil {
		base = os.TempDir()
	}
	return filepath.Join(base, "AgentBoard", "agentboard.db")
}

func Open(path string) (*Store, error) {
	if path == "" {
		path = DefaultDBPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec("PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000;"); err != nil {
		db.Close()
		return nil, err
	}
	if _, err = db.Exec(initialMigration); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	_, _ = db.Exec("INSERT OR IGNORE INTO schema_migrations(version,applied_at) VALUES(1,?)", core.Now())
	return &Store{DB: db}, nil
}

func (s *Store) Close() error { return s.DB.Close() }
func newID() string           { return uuid.NewString() }

func (s *Store) EnsureProject(ctx context.Context, name, path string) (core.Project, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return core.Project{}, err
	}
	abs = filepath.Clean(abs)
	var p core.Project
	err = s.DB.QueryRowContext(ctx, `SELECT id,name,path,created_at,updated_at FROM projects WHERE path=?`, abs).Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt, &p.UpdatedAt)
	if err == nil {
		return p, nil
	}
	if err != sql.ErrNoRows {
		return p, err
	}
	now := core.Now()
	p = core.Project{ID: newID(), Name: name, Path: abs, CreatedAt: now, UpdatedAt: now}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO projects(id,name,path,created_at,updated_at) VALUES(?,?,?,?,?)`, p.ID, p.Name, p.Path, now, now)
	return p, err
}

func (s *Store) ListProjects(ctx context.Context) ([]core.Project, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id,name,path,created_at,updated_at FROM projects ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.Project
	for rows.Next() {
		var p core.Project
		if err = rows.Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetProject(ctx context.Context, id string) (core.Project, error) {
	var p core.Project
	err := s.DB.QueryRowContext(ctx, `SELECT id,name,path,created_at,updated_at FROM projects WHERE id=?`, id).Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		err = core.ErrNotFound
	}
	return p, err
}
func (s *Store) GetProjectByPath(ctx context.Context, path string) (core.Project, error) {
	abs, _ := filepath.Abs(path)
	var p core.Project
	err := s.DB.QueryRowContext(ctx, `SELECT id,name,path,created_at,updated_at FROM projects WHERE path=?`, filepath.Clean(abs)).Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		err = core.ErrNotFound
	}
	return p, err
}

func (s *Store) CreateWorker(ctx context.Context, projectID string, in core.WorkerInput) (core.Worker, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return core.Worker{}, core.ErrInvalidInput
	}
	if in.Kind == "" {
		in.Kind = "external"
	}
	if in.Capabilities == "" {
		in.Capabilities = "{}"
	}
	now := core.Now()
	w := core.Worker{ID: newID(), ProjectID: projectID, Name: in.Name, Slug: in.Slug, Description: in.Description, Kind: in.Kind, Capabilities: in.Capabilities, Enabled: in.Enabled, CreatedAt: now, UpdatedAt: now}
	_, err := s.DB.ExecContext(ctx, `INSERT INTO workers(id,project_id,name,slug,description,kind,enabled,capabilities,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, w.ID, projectID, w.Name, w.Slug, w.Description, w.Kind, w.Enabled, w.Capabilities, now, now)
	return w, err
}

func scanWorker(row interface{ Scan(...any) error }) (core.Worker, error) {
	var w core.Worker
	err := row.Scan(&w.ID, &w.ProjectID, &w.Name, &w.Slug, &w.Description, &w.Kind, &w.Enabled, &w.Archived, &w.Capabilities, &w.CreatedAt, &w.UpdatedAt, &w.AssignedTaskCount)
	return w, err
}

const workerSelect = `SELECT w.id,w.project_id,w.name,w.slug,w.description,w.kind,w.enabled,w.archived,w.capabilities,w.created_at,w.updated_at,(SELECT count(*) FROM tasks t WHERE t.assignee_worker_id=w.id AND t.state!='complete') FROM workers w`

func (s *Store) GetWorker(ctx context.Context, id string) (core.Worker, error) {
	w, err := scanWorker(s.DB.QueryRowContext(ctx, workerSelect+` WHERE w.id=?`, id))
	if err == sql.ErrNoRows {
		err = core.ErrNotFound
	}
	return w, err
}
func (s *Store) FindWorker(ctx context.Context, projectID, idOrSlug string) (core.Worker, error) {
	w, err := scanWorker(s.DB.QueryRowContext(ctx, workerSelect+` WHERE w.project_id=? AND (w.id=? OR w.slug=?)`, projectID, idOrSlug, idOrSlug))
	if err == sql.ErrNoRows {
		err = core.ErrNotFound
	}
	return w, err
}
func (s *Store) ListWorkers(ctx context.Context, projectID string, enabledOnly bool) ([]core.Worker, error) {
	q := workerSelect + ` WHERE w.project_id=? AND w.archived=0`
	if enabledOnly {
		q += ` AND w.enabled=1`
	}
	q += ` ORDER BY w.name`
	rows, err := s.DB.QueryContext(ctx, q, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.Worker
	for rows.Next() {
		w, e := scanWorker(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, w)
	}
	return out, rows.Err()
}
func (s *Store) UpdateWorker(ctx context.Context, id string, in core.WorkerInput) (core.Worker, error) {
	if in.Kind == "" {
		in.Kind = "external"
	}
	if in.Capabilities == "" {
		in.Capabilities = "{}"
	}
	_, err := s.DB.ExecContext(ctx, `UPDATE workers SET name=?,slug=?,description=?,kind=?,enabled=?,capabilities=?,updated_at=? WHERE id=?`, in.Name, in.Slug, in.Description, in.Kind, in.Enabled, in.Capabilities, core.Now(), id)
	if err != nil {
		return core.Worker{}, err
	}
	return s.GetWorker(ctx, id)
}
func (s *Store) ArchiveWorker(ctx context.Context, id string) error {
	r, err := s.DB.ExecContext(ctx, `UPDATE workers SET archived=1,enabled=0,updated_at=? WHERE id=?`, core.Now(), id)
	if err != nil {
		return err
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return core.ErrNotFound
	}
	return nil
}

func (s *Store) StartSession(ctx context.Context, p core.Project, w core.Worker, client string) (core.WorkerSession, error) {
	if !w.Enabled || w.Archived {
		return core.WorkerSession{}, core.ErrWorkerDisabled
	}
	now := core.Now()
	v := core.WorkerSession{ID: newID(), ProjectID: p.ID, WorkerID: w.ID, StartedAt: now, ClientInfo: client, CreatedAt: now}
	_, err := s.DB.ExecContext(ctx, `INSERT INTO worker_sessions(id,project_id,worker_id,started_at,client_info,created_at) VALUES(?,?,?,?,?,?)`, v.ID, v.ProjectID, v.WorkerID, now, client, now)
	return v, err
}
func (s *Store) EndSession(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE worker_sessions SET ended_at=? WHERE id=?`, core.Now(), id)
	return err
}

func (s *Store) CreateTask(ctx context.Context, projectID string, in core.TaskInput) (core.Task, error) {
	if in.State == "" {
		in.State = "backlog"
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
	if in.Position == 0 {
		_ = s.DB.QueryRowContext(ctx, `SELECT COALESCE(MAX(position),0)+1000 FROM tasks WHERE project_id=? AND state=?`, projectID, in.State).Scan(&in.Position)
	}
	now := core.Now()
	t := core.Task{ID: newID(), ProjectID: projectID, Title: in.Title, Description: in.Description, State: in.State, Position: in.Position, Priority: in.Priority, AssigneeType: in.Assignee.Type, AssigneeWorkerID: in.Assignee.WorkerID, CreatedByType: "human", SourceTaskID: in.SourceTaskID, TestingMode: in.TestingMode, AITestInstructions: in.AITestInstructions, HumanTestInstructions: in.HumanTestInstructions, CreatedAt: now, UpdatedAt: now}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return t, err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO tasks(id,project_id,title,description,state,position,priority,assignee_type,assignee_worker_id,created_by_type,source_task_id,testing_mode,ai_test_instructions,human_test_instructions,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, t.ID, projectID, t.Title, t.Description, t.State, t.Position, t.Priority, t.AssigneeType, t.AssigneeWorkerID, t.CreatedByType, t.SourceTaskID, t.TestingMode, t.AITestInstructions, t.HumanTestInstructions, now, now)
	if err != nil {
		return t, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, newID(), t.ID, "system", "task_created", "Task created by Human", now)
	if err != nil {
		return t, err
	}
	for _, d := range in.DependencyIDs {
		if _, err = tx.ExecContext(ctx, `INSERT INTO task_dependencies(task_id,depends_on_task_id,created_at) VALUES(?,?,?)`, t.ID, d, now); err != nil {
			return t, err
		}
	}
	if err = tx.Commit(); err != nil {
		return t, err
	}
	return s.GetTask(ctx, t.ID)
}

const taskSelect = `SELECT t.id,t.project_id,t.title,t.description,t.state,t.position,t.priority,t.assignee_type,t.assignee_worker_id,t.created_by_type,t.created_by_worker_id,t.created_by_session_id,t.source_task_id,t.testing_mode,t.ai_test_instructions,t.human_test_instructions,t.created_at,t.updated_at,COALESCE(aw.name,''),COALESCE(cw.name,'') FROM tasks t LEFT JOIN workers aw ON aw.id=t.assignee_worker_id LEFT JOIN workers cw ON cw.id=t.created_by_worker_id`

func scanTask(row interface{ Scan(...any) error }) (core.Task, error) {
	var t core.Task
	err := row.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.State, &t.Position, &t.Priority, &t.AssigneeType, &t.AssigneeWorkerID, &t.CreatedByType, &t.CreatedByWorkerID, &t.CreatedBySessionID, &t.SourceTaskID, &t.TestingMode, &t.AITestInstructions, &t.HumanTestInstructions, &t.CreatedAt, &t.UpdatedAt, &t.AssigneeName, &t.CreatorName)
	return t, err
}
func (s *Store) GetTask(ctx context.Context, id string) (core.Task, error) {
	t, err := scanTask(s.DB.QueryRowContext(ctx, taskSelect+` WHERE t.id=?`, id))
	if err == sql.ErrNoRows {
		err = core.ErrNotFound
	}
	return t, err
}
func (s *Store) ListTasks(ctx context.Context, projectID string) ([]core.Task, error) {
	rows, err := s.DB.QueryContext(ctx, taskSelect+` WHERE t.project_id=? ORDER BY t.state,t.position`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.Task
	for rows.Next() {
		t, e := scanTask(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
func (s *Store) ListSpawned(ctx context.Context, taskID string) ([]core.Task, error) {
	rows, err := s.DB.QueryContext(ctx, taskSelect+` WHERE t.source_task_id=? ORDER BY t.created_at`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.Task
	for rows.Next() {
		t, e := scanTask(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func Ptr(s string) *string { return &s }
