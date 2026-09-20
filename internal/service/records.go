package service

import (
	"context"
	"database/sql"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

func (s *Service) AgentAddArtifact(ctx context.Context, a core.AgentContext, taskID string, in core.ArtifactInput) (core.Artifact, error) {
	if _, err := s.accessibleTask(ctx, a, taskID); err != nil {
		return core.Artifact{}, err
	}
	return s.addArtifact(ctx, taskID, in, "agent", &a.Worker.ID, &a.Session.ID)
}
func (s *Service) HumanAddArtifact(ctx context.Context, taskID string, in core.ArtifactInput) (core.Artifact, error) {
	if _, err := s.Store.GetTask(ctx, taskID); err != nil {
		return core.Artifact{}, err
	}
	return s.addArtifact(ctx, taskID, in, "human", nil, nil)
}
func (s *Service) addArtifact(ctx context.Context, taskID string, in core.ArtifactInput, actor string, workerID, sessionID *string) (core.Artifact, error) {
	if strings.TrimSpace(in.Name) == "" || (strings.TrimSpace(in.Path) == "" && strings.TrimSpace(in.URL) == "") {
		return core.Artifact{}, core.ErrInvalidInput
	}
	a := core.Artifact{ID: id(), TaskID: taskID, Name: in.Name, Kind: in.Kind, Description: in.Description, CreatedByType: actor, CreatedByWorkerID: workerID, CreatedBySessionID: sessionID, CreatedAt: core.Now()}
	if in.Path != "" {
		a.Path = &in.Path
	}
	if in.URL != "" {
		a.URL = &in.URL
	}
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return a, err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO task_artifacts(id,task_id,name,kind,path,url,description,created_by_type,created_by_worker_id,created_by_session_id,created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, a.ID, a.TaskID, a.Name, a.Kind, a.Path, a.URL, a.Description, a.CreatedByType, a.CreatedByWorkerID, a.CreatedBySessionID, a.CreatedAt)
	if err != nil {
		return a, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,worker_id,worker_session_id,entry_type,content,created_at) VALUES(?,?,?,?,?,?,?,?)`, id(), taskID, "system", workerID, sessionID, "artifact_event", "Artifact added: "+a.Name, a.CreatedAt)
	if err != nil {
		return a, err
	}
	err = tx.Commit()
	return a, err
}
func (s *Service) ListArtifacts(ctx context.Context, taskID string) ([]core.Artifact, error) {
	rows, err := s.Store.DB.QueryContext(ctx, `SELECT a.id,a.task_id,a.name,a.kind,a.path,a.url,a.description,a.created_by_type,a.created_by_worker_id,a.created_by_session_id,a.created_at,COALESCE(w.name,'') FROM task_artifacts a LEFT JOIN workers w ON w.id=a.created_by_worker_id WHERE a.task_id=? ORDER BY a.created_at`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.Artifact
	for rows.Next() {
		var a core.Artifact
		if err = rows.Scan(&a.ID, &a.TaskID, &a.Name, &a.Kind, &a.Path, &a.URL, &a.Description, &a.CreatedByType, &a.CreatedByWorkerID, &a.CreatedBySessionID, &a.CreatedAt, &a.CreatorName); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Service) AgentRecordTest(ctx context.Context, a core.AgentContext, taskID string, in core.TestRunInput) (core.TestRun, error) {
	if _, err := s.accessibleTask(ctx, a, taskID); err != nil {
		return core.TestRun{}, err
	}
	return s.recordTest(ctx, taskID, in, &a.Worker.ID, &a.Session.ID)
}
func (s *Service) HumanRecordTest(ctx context.Context, taskID string, in core.TestRunInput) (core.TestRun, error) {
	if _, err := s.Store.GetTask(ctx, taskID); err != nil {
		return core.TestRun{}, err
	}
	return s.recordTest(ctx, taskID, in, nil, nil)
}
func (s *Service) recordTest(ctx context.Context, taskID string, in core.TestRunInput, wid, sid *string) (core.TestRun, error) {
	valid := map[string]bool{"passed": true, "failed": true, "skipped": true, "manual_required": true}
	if !valid[in.Status] {
		return core.TestRun{}, core.ErrInvalidInput
	}
	r := core.TestRun{ID: id(), TaskID: taskID, WorkerID: wid, WorkerSessionID: sid, Runner: in.Runner, Type: in.Type, Status: in.Status, Summary: in.Summary, OutputExcerpt: in.OutputExcerpt, DurationMS: in.DurationMS, CreatedAt: core.Now()}
	if in.Command != "" {
		r.Command = &in.Command
	}
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return r, err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO test_runs(id,task_id,worker_id,worker_session_id,runner,type,command,status,summary,output_excerpt,duration_ms,created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`, r.ID, r.TaskID, r.WorkerID, r.WorkerSessionID, r.Runner, r.Type, r.Command, r.Status, r.Summary, r.OutputExcerpt, r.DurationMS, r.CreatedAt)
	if err != nil {
		return r, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,worker_id,worker_session_id,entry_type,content,created_at) VALUES(?,?,?,?,?,?,?,?)`, id(), taskID, "system", wid, sid, "test_run", "Test run "+r.Status+": "+r.Summary, r.CreatedAt)
	if err != nil {
		return r, err
	}
	err = tx.Commit()
	return r, err
}
func (s *Service) ListTests(ctx context.Context, taskID string) ([]core.TestRun, error) {
	rows, err := s.Store.DB.QueryContext(ctx, `SELECT r.id,r.task_id,r.worker_id,r.worker_session_id,r.runner,r.type,r.command,r.status,r.summary,r.output_excerpt,r.duration_ms,r.created_at,COALESCE(w.name,'') FROM test_runs r LEFT JOIN workers w ON w.id=r.worker_id WHERE r.task_id=? ORDER BY r.created_at DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.TestRun
	for rows.Next() {
		var r core.TestRun
		if err = rows.Scan(&r.ID, &r.TaskID, &r.WorkerID, &r.WorkerSessionID, &r.Runner, &r.Type, &r.Command, &r.Status, &r.Summary, &r.OutputExcerpt, &r.DurationMS, &r.CreatedAt, &r.ActorName); err != nil {
			return nil, err
		}
		if r.ActorName == "" {
			r.ActorName = "Human"
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Service) AgentReportUsage(ctx context.Context, a core.AgentContext, taskID string, in core.UsageInput) (core.UsageEvent, error) {
	if _, err := s.accessibleTask(ctx, a, taskID); err != nil {
		return core.UsageEvent{}, err
	}
	u := core.UsageEvent{ID: id(), ProjectID: a.Project.ID, TaskID: taskID, WorkerID: a.Worker.ID, WorkerSessionID: a.Session.ID, Provider: in.Provider, Model: in.Model, InputTokens: in.InputTokens, OutputTokens: in.OutputTokens, ReasoningTokens: in.ReasoningTokens, CachedTokens: in.CachedTokens, ReasoningEffort: in.ReasoningEffort, ModelCalls: in.ModelCalls, ToolCalls: in.ToolCalls, WallTimeMS: in.WallTimeMS, Source: in.Source, CreatedAt: core.Now(), WorkerName: a.Worker.Name}
	_, err := s.Store.DB.ExecContext(ctx, `INSERT INTO usage_events(id,project_id,task_id,worker_id,worker_session_id,provider,model,input_tokens,output_tokens,reasoning_tokens,cached_tokens,reasoning_effort,model_calls,tool_calls,mcp_calls,wall_time_ms,source,created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, u.ID, u.ProjectID, u.TaskID, u.WorkerID, u.WorkerSessionID, u.Provider, u.Model, u.InputTokens, u.OutputTokens, u.ReasoningTokens, u.CachedTokens, u.ReasoningEffort, u.ModelCalls, u.ToolCalls, u.MCPCalls, u.WallTimeMS, u.Source, u.CreatedAt)
	return u, err
}
func (s *Service) ListUsage(ctx context.Context, taskID string) ([]core.UsageEvent, error) {
	return s.listUsage(ctx, taskID, "")
}
func (s *Service) ListUsageForWorker(ctx context.Context, taskID, workerID string) ([]core.UsageEvent, error) {
	return s.listUsage(ctx, taskID, workerID)
}
func (s *Service) listUsage(ctx context.Context, taskID, workerID string) ([]core.UsageEvent, error) {
	q := `SELECT u.id,u.project_id,u.task_id,u.worker_id,u.worker_session_id,u.provider,u.model,u.input_tokens,u.output_tokens,u.reasoning_tokens,u.cached_tokens,u.reasoning_effort,u.model_calls,u.tool_calls,u.mcp_calls,u.wall_time_ms,u.source,u.created_at,w.name FROM usage_events u JOIN workers w ON w.id=u.worker_id WHERE u.task_id=?`
	args := []any{taskID}
	if workerID != "" {
		q += ` AND u.worker_id=?`
		args = append(args, workerID)
	}
	q += ` ORDER BY u.created_at`
	rows, err := s.Store.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.UsageEvent
	for rows.Next() {
		var u core.UsageEvent
		if err = rows.Scan(&u.ID, &u.ProjectID, &u.TaskID, &u.WorkerID, &u.WorkerSessionID, &u.Provider, &u.Model, &u.InputTokens, &u.OutputTokens, &u.ReasoningTokens, &u.CachedTokens, &u.ReasoningEffort, &u.ModelCalls, &u.ToolCalls, &u.MCPCalls, &u.WallTimeMS, &u.Source, &u.CreatedAt, &u.WorkerName); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (s *Service) ListTaskProperties(ctx context.Context, taskID string, agentVisible bool) ([]core.TaskPropertyValue, error) {
	q := `SELECT v.task_id,v.property_definition_id,d.name,d.type,d.visibility,v.value,v.updated_at FROM task_property_values v JOIN property_definitions d ON d.id=v.property_definition_id WHERE v.task_id=?`
	if agentVisible {
		q += ` AND d.visibility!='human_only'`
	}
	q += ` ORDER BY d.name`
	rows, err := s.Store.DB.QueryContext(ctx, q, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []core.TaskPropertyValue
	for rows.Next() {
		var v core.TaskPropertyValue
		if err = rows.Scan(&v.TaskID, &v.PropertyDefinitionID, &v.Name, &v.Type, &v.Visibility, &v.Value, &v.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

var _ = sql.ErrNoRows
