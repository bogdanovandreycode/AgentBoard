package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bogdanovandreycode/agentboard/internal/core"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
)

func TestHumanAPIProjectWorkerTaskFlow(t *testing.T) {
	st, err := persistence.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	p, err := st.EnsureProject(context.Background(), "demo", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	h := New(service.New(st))
	call := func(method, path, body string) int {
		r := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code >= 400 {
			t.Log(w.Body.String())
		}
		return w.Code
	}
	if got := call("POST", "/api/projects/"+p.ID+"/workers", `{"Name":"Codex","Slug":"codex","Kind":"external","Capabilities":"{}","Enabled":true}`); got != 201 {
		t.Fatalf("worker: %d", got)
	}
	if got := call("POST", "/api/projects/"+p.ID+"/tasks", `{"Title":"Implement","State":"features","Assignee":{"type":"human"},"TestingMode":"hybrid"}`); got != 201 {
		t.Fatalf("task: %d", got)
	}
	if got := call("GET", "/api/projects/"+p.ID+"/board", ""); got != 200 {
		t.Fatalf("board: %d", got)
	}
	if got := call("GET", "/api/projects", ""); got != 200 {
		t.Fatalf("projects: %d", got)
	}
	_ = http.StatusOK
}

func TestWorkerSessionDiagnostics(t *testing.T) {
	ctx := context.Background()
	st, err := persistence.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	p, err := st.EnsureProject(ctx, "test", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	worker, err := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Codex", Slug: "codex", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	session, err := st.StartSession(ctx, p, worker, "test client")
	if err != nil {
		t.Fatal(err)
	}
	handler := New(service.New(st))
	request := httptest.NewRequest("GET", "/api/workers/"+worker.ID+"/sessions", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("sessions: %d %s", response.Code, response.Body.String())
	}
	var sessions []core.WorkerSession
	if err := json.Unmarshal(response.Body.Bytes(), &sessions); err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].ID != session.ID || sessions[0].LastSeenAt == nil {
		t.Fatalf("sessions: %+v", sessions)
	}
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/api/workers/missing/sessions", nil))
	if response.Code != 404 {
		t.Fatalf("missing worker: %d", response.Code)
	}
}
