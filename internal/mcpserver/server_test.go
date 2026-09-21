package mcpserver

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestCompleteMCPToolSurface(t *testing.T) {
	ctx := context.Background()
	st, err := persistence.Open(filepath.Join(t.TempDir(), "mcp.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	p, _ := st.EnsureProject(ctx, "demo", t.TempDir())
	a, _ := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Codex", Slug: "codex", Enabled: true})
	b, _ := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Gemma", Slug: "gemma", Enabled: true})
	task, err := st.CreateTask(ctx, p.ID, core.TaskInput{Title: "Build", State: "features", Assignee: core.Assignee{Type: "worker", WorkerID: &a.ID}, TestingMode: "ai"})
	if err != nil {
		t.Fatal(err)
	}
	sess, _ := st.StartSession(ctx, p, a, "integration-test")
	srv := New(service.New(st), core.AgentContext{Project: p, Worker: a, Session: sess}, "test")
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "test"}, nil)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	ss, err := srv.MCP.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ss.Close()
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()
	if !strings.Contains(cs.InitializeResult().Instructions, "get_my_board") || !strings.Contains(cs.InitializeResult().Instructions, "Never read or modify") {
		t.Fatal("missing worker instructions")
	}
	resources, err := cs.ListResources(ctx, nil)
	if err != nil || len(resources.Resources) != 1 || resources.Resources[0].URI != "agentboard://current-worker" {
		t.Fatalf("resources: %v %v", resources, err)
	}
	resource, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "agentboard://current-worker"})
	if err != nil {
		t.Fatal(err)
	}
	var identity map[string]string
	if err := json.Unmarshal([]byte(resource.Contents[0].Text), &identity); err != nil {
		t.Fatal(err)
	}
	if identity["worker_id"] != a.ID || identity["project_id"] != p.ID {
		t.Fatalf("wrong identity: %v", identity)
	}
	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 11 {
		t.Fatalf("got %d tools", len(tools.Tools))
	}
	call := func(name string, args map[string]any) {
		t.Helper()
		r, e := cs.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
		if e != nil {
			t.Fatalf("%s: %v", name, e)
		}
		if r.IsError {
			t.Fatalf("%s returned error: %#v", name, r.Content)
		}
	}
	call("get_my_board", map[string]any{})
	call("get_task", map[string]any{"task_id": task.ID})
	call("list_workers", map[string]any{})
	call("create_feature", map[string]any{"title": "Graphics", "reason": "needed", "assignee_type": "worker", "assignee_worker_id": b.ID, "source_task_id": task.ID})
	call("add_history", map[string]any{"task_id": task.ID, "content": "Started"})
	call("add_artifact", map[string]any{"task_id": task.ID, "name": "Patch", "kind": "file", "path": "patch.diff"})
	call("move_to_in_progress", map[string]any{"task_id": task.ID})
	call("move_to_testing", map[string]any{"task_id": task.ID})
	call("record_test_run", map[string]any{"task_id": task.ID, "runner": "go", "type": "automated", "status": "passed", "summary": "ok"})
	call("move_to_verification", map[string]any{"task_id": task.ID})
	call("report_usage", map[string]any{"task_id": task.ID, "provider": "openai", "model": "gpt", "input_tokens": 10.0})
	activity, err := st.GetWorker(ctx, a.ID)
	if err != nil || activity.MCPCalls != 11 || activity.LastActivityAt == nil || activity.ActiveSessionCount != 1 {
		t.Fatalf("activity: %+v %v", activity, err)
	}
	other, err := st.GetWorker(ctx, b.ID)
	if err != nil || other.MCPCalls != 0 || other.LastActivityAt != nil {
		t.Fatalf("activity leaked: %+v %v", other, err)
	}
}
