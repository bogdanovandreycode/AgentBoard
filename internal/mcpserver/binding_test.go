package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func bindingClient(t *testing.T, svc *service.Service, binding Binding, launch string, roots ...*mcp.Root) (*mcp.ClientSession, *mcp.Client) {
	t.Helper()
	srv := NewWithBinding(svc, binding, "test")
	srv.Workspace = &WorkspaceContext{LaunchDirectory: launch}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	t.Cleanup(cancel)
	st, ct := mcp.NewInMemoryTransports()
	ss, err := srv.MCP.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ss.Close() })
	c := mcp.NewClient(&mcp.Implementation{Name: "binding-test", Version: "1"}, nil)
	c.AddRoots(roots...)
	cs, err := c.Connect(ctx, ct, &mcp.ClientSessionOptions{ProtocolVersion: "2025-06-18"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs, c
}

func callText(t *testing.T, cs *mcp.ClientSession, name string, wantError bool) string {
	t.Helper()
	r, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	text := r.Content[0].(*mcp.TextContent).Text
	if r.IsError != wantError {
		t.Fatalf("%s error=%v: %s", name, r.IsError, text)
	}
	return text
}

func TestMCPStartsWithoutProjectsAndAsksAgentToTellUser(t *testing.T) {
	st, err := persistence.Open(filepath.Join(t.TempDir(), "empty.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	cs, _ := bindingClient(t, service.New(st), Binding{WorkerRef: "codex"}, t.TempDir())
	tools, err := cs.ListTools(context.Background(), nil)
	if err != nil || len(tools.Tools) == 0 {
		t.Fatalf("tools must be listed before any project exists: %v", err)
	}
	text := callText(t, cs, "get_my_board", true)
	if !strings.HasPrefix(text, "NO_PROJECTS:") || !strings.Contains(text, "agentboard init") || !strings.Contains(text, "Tell the user") {
		t.Fatalf("unhelpful empty-database error: %s", text)
	}
	resource, err := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{URI: "agentboard://current-worker"})
	if err != nil || !strings.Contains(resource.Contents[0].Text, "not bound yet") {
		t.Fatalf("unbound identity: %+v %v", resource, err)
	}
}

func TestMCPBindsFromAnyFolderUsingRootsLaunchDirectoryOrProject(t *testing.T) {
	ctx := context.Background()
	st, err := persistence.Open(filepath.Join(t.TempDir(), "board.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	svc := service.New(st)
	one, _ := st.EnsureProject(ctx, "one", t.TempDir())
	two, _ := st.EnsureProject(ctx, "two", t.TempDir())
	codexOne, _ := st.CreateWorker(ctx, one.ID, core.WorkerInput{Name: "Codex", Slug: "codex", Enabled: true})
	codexTwo, _ := st.CreateWorker(ctx, two.ID, core.WorkerInput{Name: "Codex", Slug: "codex", Enabled: true})
	if _, err := st.CreateWorker(ctx, two.ID, core.WorkerInput{Name: "Gemini", Slug: "gemini", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	neutral := t.TempDir()
	session := func(text string) map[string]string {
		t.Helper()
		var v struct{ Session map[string]string }
		if err := json.Unmarshal([]byte(text), &v); err != nil {
			t.Fatal(err)
		}
		return v.Session
	}

	// Outside any project, the error lists registered projects for the user.
	cs, _ := bindingClient(t, svc, Binding{WorkerRef: "codex"}, neutral)
	if text := callText(t, cs, "get_my_board", true); !strings.HasPrefix(text, "PROJECT_REQUIRED:") || !strings.Contains(text, one.Path) || !strings.Contains(text, two.Path) {
		t.Fatalf("project required: %s", text)
	}

	// Client roots select the project even when launched from another folder;
	// two agents on different projects run side by side on one database.
	nested := filepath.Join(two.Path, "src")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	csTwo, clientTwo := bindingClient(t, svc, Binding{WorkerRef: "codex", Client: "test"}, one.Path, &mcp.Root{URI: rootURI(nested)})
	csOne, _ := bindingClient(t, svc, Binding{WorkerRef: "codex", Client: "test"}, filepath.Join(one.Path))
	if got := session(callText(t, csTwo, "get_my_board", false)); got["worker_id"] != codexTwo.ID || got["binding_source"] != "client_roots" {
		t.Fatalf("roots binding: %+v", got)
	}
	if got := session(callText(t, csOne, "get_my_board", false)); got["worker_id"] != codexOne.ID || got["binding_source"] != "launch_directory" {
		t.Fatalf("launch directory binding: %+v", got)
	}

	// A workspace change rebinds the session to the new project.
	clientTwo.RemoveRoots(rootURI(nested))
	clientTwo.AddRoots(&mcp.Root{URI: rootURI(one.Path)})
	deadline := time.Now().Add(2 * time.Second)
	for session(callText(t, csTwo, "get_my_board", false))["project_id"] != one.ID {
		if time.Now().After(deadline) {
			t.Fatal("roots change did not rebind")
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Worker selection problems are reported, never guessed.
	cs, _ = bindingClient(t, svc, Binding{}, two.Path)
	if text := callText(t, cs, "get_my_board", true); !strings.HasPrefix(text, "WORKER_REQUIRED:") || !strings.Contains(text, "gemini") {
		t.Fatalf("worker required: %s", text)
	}
	cs, _ = bindingClient(t, svc, Binding{}, one.Path)
	if got := session(callText(t, cs, "get_my_board", false)); got["worker_id"] != codexOne.ID {
		t.Fatalf("single worker default: %+v", got)
	}
	cs, _ = bindingClient(t, svc, Binding{WorkerRef: "missing"}, one.Path)
	if text := callText(t, cs, "get_my_board", true); !strings.HasPrefix(text, "WORKER_NOT_FOUND:") || !strings.Contains(text, "codex") {
		t.Fatalf("worker not found: %s", text)
	}

	// An explicit --project works from any folder; unknown paths are reported.
	cs, _ = bindingClient(t, svc, Binding{ProjectPath: two.Path, WorkerRef: "codex"}, neutral)
	if got := session(callText(t, cs, "get_my_board", false)); got["project_id"] != two.ID || got["binding_source"] != "--project" {
		t.Fatalf("explicit project: %+v", got)
	}
	cs, _ = bindingClient(t, svc, Binding{ProjectPath: neutral, WorkerRef: "codex"}, neutral)
	if text := callText(t, cs, "get_my_board", true); !strings.HasPrefix(text, "PROJECT_NOT_FOUND:") {
		t.Fatalf("explicit unknown project: %s", text)
	}
}

func TestEveryToolDeclaresAllBehaviourHints(t *testing.T) {
	f := newImageFixture(t)
	tools, err := f.cs.ListTools(f.ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) < 15 {
		t.Fatalf("unexpected tool count %d", len(tools.Tools))
	}
	writes := map[string]bool{"create_feature": true, "add_history": true, "add_artifact": true, "record_test_run": true, "report_usage": true, "upload_image": true, "upload_file": true, "move_to_in_progress": true, "move_to_testing": true, "move_to_verification": true}
	for _, tool := range tools.Tools {
		raw, err := json.Marshal(tool)
		if err != nil {
			t.Fatal(err)
		}
		var wire struct {
			Title       string
			Annotations map[string]any
		}
		if err := json.Unmarshal(raw, &wire); err != nil {
			t.Fatal(err)
		}
		for _, hint := range []string{"readOnlyHint", "destructiveHint", "idempotentHint", "openWorldHint"} {
			if _, ok := wire.Annotations[hint].(bool); !ok {
				t.Fatalf("%s: %s is not an explicit boolean: %s", tool.Name, hint, raw)
			}
		}
		if wire.Title == "" || wire.Annotations["readOnlyHint"] == writes[tool.Name] || wire.Annotations["destructiveHint"] != false || wire.Annotations["openWorldHint"] != false {
			t.Fatalf("%s: wrong hints %s", tool.Name, raw)
		}
	}
}

func TestStoresInSeveralProcessesShareTheDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shared.db")
	first, err := persistence.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := persistence.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	var mode string
	if err := second.DB.QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil || !strings.EqualFold(mode, "wal") {
		t.Fatalf("journal mode %q: %v", mode, err)
	}
	ctx := context.Background()
	p, err := first.EnsureProject(ctx, "shared", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 2)
	for i, st := range []*persistence.Store{first, second} {
		go func(i int, st *persistence.Store) {
			for n := 0; n < 20; n++ {
				if _, err := st.CreateTask(ctx, p.ID, core.TaskInput{Title: "task", State: "features", Assignee: core.Assignee{Type: "unassigned"}, TestingMode: "ai"}); err != nil {
					done <- err
					return
				}
			}
			done <- nil
		}(i, st)
	}
	for range 2 {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
	tasks, err := second.ListTasks(ctx, p.ID)
	if err != nil || len(tasks) != 40 {
		t.Fatalf("tasks %d: %v", len(tasks), err)
	}
}
