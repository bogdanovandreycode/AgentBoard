package persistence

import (
	"context"
	"database/sql"
	"github.com/bogdanovandreycode/agentboard/internal/core"
	"path/filepath"
	"testing"
)

func TestSessionActivityAcrossConnections(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "sessions.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	other, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	p, err := st.EnsureProject(ctx, "test", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	w, err := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Codex", Slug: "codex", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	session, err := st.StartSession(ctx, p, w, "test")
	if err != nil {
		t.Fatal(err)
	}
	check := func(active int64) {
		t.Helper()
		got, err := other.GetWorker(ctx, w.ID)
		if err != nil || got.ActiveSessionCount != active || got.SessionCount != 1 {
			t.Fatalf("worker: %+v %v", got, err)
		}
	}
	check(1)
	if _, err := st.DB.Exec(`UPDATE worker_sessions SET last_seen_at='2000-01-01T00:00:00Z' WHERE id=?`, session.ID); err != nil {
		t.Fatal(err)
	}
	check(0)
	if err := st.TouchSession(ctx, session.ID); err != nil {
		t.Fatal(err)
	}
	check(1)
	if err := st.EndSession(ctx, session.ID); err != nil {
		t.Fatal(err)
	}
	check(0)
	sessions, err := other.ListWorkerSessions(ctx, w.ID)
	if err != nil || len(sessions) != 1 || sessions[0].EndedAt == nil {
		t.Fatalf("sessions: %+v %v", sessions, err)
	}
	if err := st.TouchSession(ctx, session.ID); err != nil {
		t.Fatal(err)
	}
	check(0)
}

func TestUpgradePreservesSessionsAndTracksHistory(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "legacy.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(initialMigration); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(sessionCallsMigration); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO schema_migrations VALUES(1,'old'),(2,'old');
 INSERT INTO projects VALUES('p','demo','demo','old','old');
 INSERT INTO workers(id,project_id,name,slug,kind,created_at,updated_at) VALUES('w','p','Codex','codex','external','old','old');
 INSERT INTO worker_sessions(id,project_id,worker_id,started_at,created_at,mcp_calls) VALUES('s','p','w','old','old',7);`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	worker, err := st.GetWorker(ctx, "w")
	if err != nil || worker.MCPCalls != 7 || worker.SessionCount != 1 || worker.ActiveSessionCount != 0 || worker.LastActivityAt != nil {
		t.Fatalf("legacy worker: %+v %v", worker, err)
	}
	task, err := st.CreateTask(ctx, "p", core.TaskInput{Title: "test", Assignee: core.Assignee{Type: "human"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.DB.Exec(`INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES('h',?,'human','human_comment','test','2099-01-01T00:00:00Z')`, task.ID); err != nil {
		t.Fatal(err)
	}
	var updated string
	if err := st.DB.QueryRow(`SELECT updated_at FROM tasks WHERE id=?`, task.ID).Scan(&updated); err != nil {
		t.Fatal(err)
	}
	if updated != "2099-01-01T00:00:00Z" {
		t.Fatalf("task activity: %s", updated)
	}
}
