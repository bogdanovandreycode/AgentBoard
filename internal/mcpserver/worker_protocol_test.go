package mcpserver

import (
	"context"
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func decodeJSON(t *testing.T, r *mcp.CallToolResult, out any) {
	t.Helper()
	if r.StructuredContent != nil {
		t.Fatal("response duplicated JSON in structuredContent")
	}
	if err := json.Unmarshal([]byte(r.Content[0].(*mcp.TextContent).Text), out); err != nil {
		t.Fatal(err)
	}
}

func TestCompactBoardTransitionsSchemasAndUsageScope(t *testing.T) {
	f := newImageFixture(t)
	secretDescription := "A long description that belongs only in get_task"
	if _, err := f.svc.HumanUpdateTask(f.ctx, f.task.ID, core.TaskUpdate{Description: &secretDescription}); err != nil {
		t.Fatal(err)
	}
	var board map[string][]map[string]any
	decodeJSON(t, f.call(t, "get_my_board", map[string]any{}, false), &board)
	if len(board["features"]) != 1 || board["features"][0]["Description"] != nil {
		t.Fatalf("non-compact board: %+v", board)
	}
	var moved map[string]any
	decodeJSON(t, f.call(t, "move_to_in_progress", map[string]any{"task_id": strings.ToUpper(f.task.ID[:8])}, false), &moved)
	if moved["ID"] != f.task.ID || moved["Description"] != nil || moved["State"] != "in_progress" {
		t.Fatalf("non-compact transition: %+v", moved)
	}
	var usage map[string]any
	decodeJSON(t, f.call(t, "report_usage", map[string]any{"task_id": f.task.ID[:8]}, false), &usage)
	if usage["mcp_calls_scope"] != "session_total_including_this_call" || usage["WorkerSessionID"] != f.agent.Session.ID || usage["MCPCalls"] != float64(4) {
		t.Fatalf("usage scope: %+v", usage)
	}
	tools, err := f.cs.ListTools(f.ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range tools.Tools {
		if tool.Name != "record_test_run" && tool.Name != "create_feature" {
			continue
		}
		data, _ := json.Marshal(tool.InputSchema)
		if !strings.Contains(string(data), `"enum"`) {
			t.Fatalf("missing enums: %s", data)
		}
	}
	f.call(t, "record_test_run", map[string]any{"task_id": f.task.ID, "runner": "go", "type": "automated", "summary": "bad status", "status": "success"}, true)
}

func workspaceClient(t *testing.T, f imageFixture, version string, workspace *WorkspaceContext, options *mcp.ClientOptions, roots ...*mcp.Root) (*mcp.ClientSession, *mcp.Client) {
	t.Helper()
	srv := New(f.svc, f.agent, "test")
	srv.Workspace = workspace
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	t.Cleanup(cancel)
	st, ct := mcp.NewInMemoryTransports()
	ss, err := srv.MCP.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ss.Close() })
	c := mcp.NewClient(&mcp.Implementation{Name: "workspace-test", Version: "1"}, options)
	c.AddRoots(roots...)
	cs, err := c.Connect(ctx, ct, &mcp.ClientSessionOptions{ProtocolVersion: version})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs, c
}

func rootURI(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return (&url.URL{Scheme: "file", Path: path}).String()
}

func TestWorkspaceMismatchConfirmationAndNoQuestionFallback(t *testing.T) {
	for _, version := range []string{"2025-03-26", "2026-07-28"} {
		for _, action := range []string{"unsupported", "accept", "decline", "cancel"} {
			t.Run(version+"/"+action, func(t *testing.T) {
				f := newImageFixture(t)
				other, err := f.svc.Store.EnsureProject(f.ctx, "other", t.TempDir())
				if err != nil {
					t.Fatal(err)
				}
				questions := 0
				var options *mcp.ClientOptions
				if action != "unsupported" {
					options = &mcp.ClientOptions{ElicitationHandler: func(_ context.Context, req *mcp.ElicitRequest) (*mcp.ElicitResult, error) {
						questions++
						if !strings.Contains(req.Params.Message, "PROJECT_CONTEXT_MISMATCH") {
							t.Error("missing mismatch context")
						}
						return &mcp.ElicitResult{Action: action, Content: map[string]any{"confirm": action == "accept"}}, nil
					}}
				}
				cs, _ := workspaceClient(t, f, version, &WorkspaceContext{LaunchDirectory: other.Path, LaunchProject: &other}, options, &mcp.Root{URI: rootURI(other.Path)})
				before, _ := f.svc.Store.GetWorker(f.ctx, f.agent.Worker.ID)
				r, err := cs.CallTool(f.ctx, &mcp.CallToolParams{Name: "get_my_board", Arguments: map[string]any{}})
				if err != nil {
					t.Fatal(err)
				}
				if action == "accept" {
					if r.IsError {
						t.Fatalf("accepted: %+v", r.Content)
					}
					if _, err := cs.CallTool(f.ctx, &mcp.CallToolParams{Name: "get_my_board", Arguments: map[string]any{}}); err != nil {
						t.Fatal(err)
					}
					if questions != 1 {
						t.Fatalf("confirmation repeated %d times", questions)
					}
				} else {
					if !r.IsError || !strings.Contains(r.Content[0].(*mcp.TextContent).Text, "PROJECT_CONTEXT_MISMATCH") {
						t.Fatalf("missing actionable error: %+v", r)
					}
					after, _ := f.svc.Store.GetWorker(f.ctx, f.agent.Worker.ID)
					if before.MCPCalls != after.MCPCalls {
						t.Fatal("task handler ran before confirmation")
					}
				}
			})
		}
	}
}

func TestWorkspaceRootsOverrideAndRootChange(t *testing.T) {
	f := newImageFixture(t)
	other, err := f.svc.Store.EnsureProject(f.ctx, "other", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	// A neutral launch directory still detects a conflicting client workspace.
	cs, c := workspaceClient(t, f, "2025-03-26", &WorkspaceContext{LaunchDirectory: t.TempDir()}, nil, &mcp.Root{URI: rootURI(other.Path)})
	r, err := cs.CallTool(f.ctx, &mcp.CallToolParams{Name: "get_my_board", Arguments: map[string]any{}})
	if err != nil || !r.IsError {
		t.Fatalf("client roots mismatch: %+v %v", r, err)
	}
	c.RemoveRoots(rootURI(other.Path))
	c.AddRoots(&mcp.Root{URI: rootURI(f.agent.Project.Path)})
	// Notifications are asynchronous; use a bounded condition wait in this test.
	deadline := time.Now().Add(2 * time.Second)
	for {
		r, err = cs.CallTool(f.ctx, &mcp.CallToolParams{Name: "get_my_board", Arguments: map[string]any{}})
		if err != nil {
			t.Fatal(err)
		}
		if !r.IsError {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("workspace cache not invalidated")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cs, _ = workspaceClient(t, f, "2025-03-26", &WorkspaceContext{LaunchDirectory: other.Path, LaunchProject: &other, AllowMismatch: true}, nil)
	r, err = cs.CallTool(f.ctx, &mcp.CallToolParams{Name: "get_my_board", Arguments: map[string]any{}})
	if err != nil || r.IsError {
		t.Fatalf("explicit mismatch override: %+v %v", r, err)
	}
}

func TestGlobalProjectContextSwitchAffectsExistingSession(t *testing.T) {
	f := newImageFixture(t)
	other, err := f.svc.Store.EnsureProject(f.ctx, "other", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	cs, _ := workspaceClient(t, f, "2025-03-26", &WorkspaceContext{LaunchDirectory: other.Path, LaunchProject: &other}, nil)
	call := func(name string, args map[string]any) *mcp.CallToolResult {
		t.Helper()
		result, err := cs.CallTool(f.ctx, &mcp.CallToolParams{Name: name, Arguments: args})
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	if result := call("get_my_board", map[string]any{}); !result.IsError {
		t.Fatal("enabled policy did not block")
	}
	settings, err := f.svc.GetGlobalSettings(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	settings.CheckProjectContext = false
	if _, err := f.svc.UpdateGlobalSettings(f.ctx, settings); err != nil {
		t.Fatal(err)
	}
	if result := call("get_my_board", map[string]any{}); result.IsError {
		t.Fatalf("disabled policy blocked: %+v", result.Content)
	}
	if result := call("get_task", map[string]any{"task_id": "00000000"}); !result.IsError {
		t.Fatal("visibility bypassed")
	}
	settings.CheckProjectContext = true
	if _, err := f.svc.UpdateGlobalSettings(f.ctx, settings); err != nil {
		t.Fatal(err)
	}
	if result := call("get_my_board", map[string]any{}); !result.IsError {
		t.Fatal("re-enabled policy did not apply")
	}
}
