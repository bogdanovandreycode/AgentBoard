package persistence

import (
	"context"
	"github.com/bogdanovandreycode/agentboard/internal/core"
)

func (s *Store) TouchSession(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE worker_sessions SET last_seen_at=? WHERE id=? AND ended_at IS NULL`, core.Now(), id)
	return err
}

// ListWorkerSessions returns the latest 50 sessions for human diagnostics.
func (s *Store) ListWorkerSessions(ctx context.Context, workerID string) ([]core.WorkerSession, error) {
	if _, err := s.GetWorker(ctx, workerID); err != nil {
		return nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT id,project_id,worker_id,started_at,client_info,created_at,ended_at,last_activity_at,last_seen_at,mcp_calls FROM worker_sessions WHERE worker_id=? ORDER BY started_at DESC LIMIT 50`, workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []core.WorkerSession{}
	for rows.Next() {
		var v core.WorkerSession
		if err := rows.Scan(&v.ID, &v.ProjectID, &v.WorkerID, &v.StartedAt, &v.ClientInfo, &v.CreatedAt, &v.EndedAt, &v.LastActivityAt, &v.LastSeenAt, &v.MCPCalls); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
