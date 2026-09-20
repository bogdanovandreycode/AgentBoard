package httpapi

import (
	"bytes"
	"context"
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
