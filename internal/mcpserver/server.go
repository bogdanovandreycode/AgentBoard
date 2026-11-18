package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	MCP       *mcp.Server
	Service   *service.Service
	Workspace *WorkspaceContext
	// Binding selects the project and worker on the first tool call. It is nil
	// when the server was created for an already started agent session.
	Binding *Binding
	guard   workspaceGuard

	agentMu       sync.Mutex
	agent         *core.AgentContext
	bindingSource string
	rebind        bool
}

// Binding is the MCP launch configuration. Empty ProjectPath detects the project
// from client workspace roots, then from the launch directory. Empty WorkerRef
// selects the project's only enabled worker.
type Binding struct {
	ProjectPath string
	WorkerRef   string
	Client      string
}

const instructions = `AgentBoard is the task source of truth for this worker.
At the start of a work session call get_my_board. It also reports which AgentBoard
project and worker this session is bound to; confirm it matches your workspace.
The MCP server works without the AgentBoard web UI and selects the project from
--project, your client workspace roots, or its launch folder. If a tool returns
NO_PROJECTS, PROJECT_REQUIRED, PROJECT_NOT_FOUND, PROJECT_AMBIGUOUS, NO_WORKERS,
WORKER_REQUIRED, WORKER_NOT_FOUND or WORKER_DISABLED, stop using AgentBoard tools
and relay that message to the user: only the user can fix the configuration. Continue in_progress tasks first;
otherwise select an assigned feature. Use AgentBoard tools for task state, history,
testing and usage. Only move forward using the explicit workflow tools and respect
human testing requirements. Never read or modify AgentBoard SQLite storage directly.
For human/hybrid tasks, move completed implementation to testing, leave handoff
instructions, then finish the session or continue another task. Never wait, poll,
or repeatedly retry while human testing is pending. Do not bypass this gate.
task_id accepts a full UUID or a unique case-insensitive prefix of at least 8
characters among your visible assigned tasks. Source/dependency IDs require full UUIDs.
get_my_board returns summaries; use get_task for full descriptions and history.
Responses contain one JSON text payload, plus native image blocks when applicable.
report_usage MCPCalls is a cumulative session snapshot including that call,
not a per-task count. Do not sum these snapshots across tasks or usage events.
AgentBoard MCP is the only supported worker interface. These capabilities enforce
workflow permissions; OS-level isolation requires a sandbox.
get_task includes images referenced in the task description, testing instructions
and comments as native image content. Read screenshots with their Markdown links.
If an image is omitted from get_task, call get_image with its task_id and image_id.
To share your screenshot, call upload_image with the assigned task_id and a local
image path or base64 data, then include its returned Markdown link in add_history.
Images are stored as compressed WebP files by AgentBoard; never use local file
paths or base64 data as comment attachments.
For other attachments use upload_file, then include its returned Markdown link in
add_history. get_task lists file metadata in Files and History[].Files. Use get_file
with task_id and file_id to read a file in bounded UTF-8 or base64 chunks; advance
offset to next_offset until eof. Files are stored unchanged, separately from images.`

// New creates a server for an already started agent session.
func New(s *service.Service, a core.AgentContext, version string) *Server {
	x := newServer(s, version)
	x.agent = &a
	return x
}

// NewWithBinding creates a server that resolves its project and worker lazily,
// so it can start in any folder and report configuration problems to the agent.
func NewWithBinding(s *service.Service, binding Binding, version string) *Server {
	x := newServer(s, version)
	x.Binding = &binding
	return x
}

func newServer(s *service.Service, version string) *Server {
	x := &Server{Service: s, guard: workspaceGuard{sessions: make(map[*mcp.ServerSession]*workspaceSession)}}
	x.MCP = mcp.NewServer(&mcp.Implementation{Name: "agentboard", Title: "AgentBoard", Version: version}, &mcp.ServerOptions{Instructions: instructions,
		RootsListChangedHandler: func(_ context.Context, req *mcp.RootsListChangedRequest) {
			x.invalidateWorkspace(req.Session)
			x.agentMu.Lock()
			x.rebind = true
			x.agentMu.Unlock()
		},
	})
	x.MCP.AddResource(&mcp.Resource{URI: "agentboard://current-worker", Name: "Current AgentBoard worker", MIMEType: "application/json"}, func(ctx context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		b, err := json.Marshal(x.identity())
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: "agentboard://current-worker", MIMEType: "application/json", Text: string(b)}}}, nil
	})
	x.register()
	return x
}

func (s *Server) identity() map[string]any {
	s.agentMu.Lock()
	a, source := s.agent, s.bindingSource
	s.agentMu.Unlock()
	v := map[string]any{"workspace": s.Workspace, "instructions": instructions}
	if a == nil {
		v["status"] = "not bound yet: call get_my_board to bind this MCP session to an AgentBoard project and worker"
		return v
	}
	v["status"] = "bound"
	v["project"], v["project_id"], v["project_path"] = a.Project.Name, a.Project.ID, a.Project.Path
	v["worker"], v["worker_id"], v["worker_slug"] = a.Worker.Name, a.Worker.ID, a.Worker.Slug
	if source != "" {
		v["binding_source"] = source
	}
	return v
}

func (s *Server) currentAgent() *core.AgentContext {
	s.agentMu.Lock()
	defer s.agentMu.Unlock()
	return s.agent
}

// Run serves the transport and ends the worker session it started.
func (s *Server) Run(ctx context.Context, transport mcp.Transport) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer func() {
		if a := s.currentAgent(); a != nil {
			_ = s.Service.EndAgentSession(context.Background(), *a)
		}
	}()
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if a := s.currentAgent(); a != nil {
					_ = s.Service.Store.TouchSession(ctx, a.Session.ID)
				}
			}
		}
	}()
	return s.MCP.Run(ctx, transport)
}

type agentKey struct{}

func withAgent(ctx context.Context, a core.AgentContext) context.Context {
	return context.WithValue(ctx, agentKey{}, a)
}

// agentFrom returns the agent bound by addTool before the handler runs.
func agentFrom(ctx context.Context) core.AgentContext {
	a, _ := ctx.Value(agentKey{}).(core.AgentContext)
	return a
}

func (s *Server) count(ctx context.Context) {
	_, _ = s.Service.Store.DB.ExecContext(ctx, `UPDATE worker_sessions SET mcp_calls=mcp_calls+1,last_activity_at=?,last_seen_at=? WHERE id=?`, core.Now(), core.Now(), agentFrom(ctx).Session.ID)
}

// annotations declares all four MCP behaviour hints explicitly. AgentBoard tools
// only touch the local AgentBoard database (closed world) and never delete or
// overwrite data: history is append-only and transitions only move forward.
func annotations(readOnly, idempotent bool) *mcp.ToolAnnotations {
	destructive, openWorld := false, false
	return &mcp.ToolAnnotations{ReadOnlyHint: readOnly, DestructiveHint: &destructive, IdempotentHint: idempotent, OpenWorldHint: &openWorld}
}

func output(v any) (*mcp.CallToolResult, any, error) {
	b, _ := json.Marshal(v)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil, nil
}
func failure(err error) (*mcp.CallToolResult, any, error) {
	message := err.Error()
	switch {
	case errors.Is(err, core.ErrHumanTestPending):
		message = "HUMAN_TEST_PENDING: leave the task in testing, provide handoff instructions, and finish this session or continue another task. Do not wait, poll, or retry verification until human testing is resolved"
	case errors.Is(err, core.ErrInvalidTransition):
		message += ": AI moves tasks one step forward only (features → in_progress → testing → verification). move_to_verification also needs a passed record_test_run from this worker. Call get_task and follow its Workflow.next_tool."
	case errors.Is(err, core.ErrTaskNotAccessible):
		message += ": the task does not exist, is not assigned to this worker, or is in Backlog/Complete. Use task IDs from get_my_board."
	case errors.Is(err, core.ErrWorkerDisabled):
		message += ": this worker was disabled or archived in AgentBoard. Stop and tell the user."
	}
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: message}}}, nil, nil
}

type taskIDInput struct {
	TaskID string `json:"task_id" jsonschema:"Full UUID or a unique case-insensitive prefix of at least 8 characters from get_my_board"`
}
type historyInput struct {
	TaskID  string `json:"task_id" jsonschema:"Full UUID or a unique case-insensitive prefix of at least 8 characters from get_my_board"`
	Content string `json:"content" jsonschema:"Markdown comment, for example progress notes, handoff instructions for a human tester, or links returned by upload_image/upload_file"`
}
type artifactInput struct {
	TaskID      string `json:"task_id" jsonschema:"Full UUID or a unique case-insensitive prefix of at least 8 characters from get_my_board"`
	Name        string `json:"name" jsonschema:"Short human-readable deliverable name"`
	Kind        string `json:"kind" jsonschema:"Free-form artifact category, for example file, screenshot, url, build, or report"`
	Path        string `json:"path,omitempty" jsonschema:"Local path of the deliverable; provide path, url, or both"`
	URL         string `json:"url,omitempty" jsonschema:"URL of the deliverable, for example a pull request or deployed preview"`
	Description string `json:"description,omitempty"`
}
type featureInput struct {
	Title               string            `json:"title"`
	Description         string            `json:"description,omitempty"`
	Reason              string            `json:"reason" jsonschema:"Why this work is needed, for example the finding that triggered it"`
	AssigneeType        string            `json:"assignee_type"`
	AssigneeWorkerID    *string           `json:"assignee_worker_id,omitempty" jsonschema:"Worker ID from list_workers; required when assignee_type is worker"`
	Priority            string            `json:"priority,omitempty"`
	TestingMode         string            `json:"testing_mode,omitempty"`
	SourceTaskID        *string           `json:"source_task_id,omitempty"`
	DependsOn           []string          `json:"depends_on,omitempty"`
	SuggestedProperties map[string]string `json:"suggested_properties,omitempty"`
}
type testInput struct {
	TaskID        string `json:"task_id" jsonschema:"Full UUID or a unique case-insensitive prefix of at least 8 characters from get_my_board"`
	Runner        string `json:"runner" jsonschema:"Free-form tool that ran the tests, for example go test, pytest, or playwright"`
	Type          string `json:"type" jsonschema:"Free-form test category, for example automated, unit, integration, or manual. Agent results never satisfy human testing."`
	Command       string `json:"command,omitempty" jsonschema:"Exact command that was run"`
	Status        string `json:"status"`
	Summary       string `json:"summary" jsonschema:"One or two sentences describing the result"`
	OutputExcerpt string `json:"output_excerpt,omitempty"`
	DurationMS    *int64 `json:"duration_ms,omitempty"`
}
type usageInput struct {
	TaskID          string  `json:"task_id" jsonschema:"Full UUID or a unique case-insensitive prefix of at least 8 characters from get_my_board"`
	Provider        *string `json:"provider,omitempty"`
	Model           *string `json:"model,omitempty"`
	InputTokens     *int64  `json:"input_tokens,omitempty"`
	OutputTokens    *int64  `json:"output_tokens,omitempty"`
	ReasoningTokens *int64  `json:"reasoning_tokens,omitempty"`
	CachedTokens    *int64  `json:"cached_tokens,omitempty"`
	ReasoningEffort *string `json:"reasoning_effort,omitempty"`
	ModelCalls      *int64  `json:"model_calls,omitempty"`
	ToolCalls       *int64  `json:"tool_calls,omitempty"`
	WallTimeMS      *int64  `json:"wall_time_ms,omitempty"`
	Source          *string `json:"source,omitempty"`
}

func (s *Server) register() {
	addTool(s, &mcp.Tool{Name: "get_my_board", Title: "Get my AgentBoard tasks", Annotations: annotations(true, true), Description: "Start every work session here. Returns the bound AgentBoard project and worker (session), compact summaries of your assigned tasks grouped by state (board), and a suggested next_step. Confirm session.project_path matches your workspace. Use get_task for full details."}, func(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentBoard(ctx, agentFrom(ctx))
		if e != nil {
			return failure(e)
		}
		board := make(map[string][]taskSummary, len(v))
		for state, tasks := range v {
			board[state] = make([]taskSummary, 0, len(tasks))
			for _, task := range tasks {
				board[state] = append(board[state], s.summary(ctx, task))
			}
		}
		return output(boardResponse{Session: s.sessionInfo(agentFrom(ctx)), Board: board, NextStep: boardNextStep(board)})
	})
	addTool(s, &mcp.Tool{Name: "get_task", Title: "Get task details", Annotations: annotations(true, true), Description: "Get assigned task details and view screenshots from its description, testing instructions and comments. Images are returned as native image content; use get_image for any omitted images."}, func(ctx context.Context, _ *mcp.CallToolRequest, in taskIDInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentTask(ctx, agentFrom(ctx), in.TaskID)
		if e != nil {
			return failure(e)
		}
		return s.taskOutput(ctx, v)
	})
	s.registerImages()
	s.registerFiles()
	addTool(s, &mcp.Tool{Name: "list_workers", Title: "List workers", Annotations: annotations(true, true), Description: "List enabled Workers available when assigning a newly-created Feature."}, func(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.Store.ListWorkers(ctx, agentFrom(ctx).Project.ID, true)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	addTool(s, &mcp.Tool{Name: "create_feature", Title: "Create feature task", Annotations: annotations(false, false), InputSchema: enumSchema[featureInput](map[string][]any{"assignee_type": {"unassigned", "human", "worker"}, "priority": {"critical", "high", "medium", "low"}, "testing_mode": {"ai", "human", "hybrid"}}), Description: "Create a Feature attributed to the current Worker and optionally assign it. source_task_id and depends_on require full UUIDs."}, func(ctx context.Context, _ *mcp.CallToolRequest, in featureInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.CreateFeature(ctx, agentFrom(ctx), core.FeatureInput{Title: in.Title, Description: in.Description, Reason: in.Reason, Assignee: core.Assignee{Type: in.AssigneeType, WorkerID: in.AssigneeWorkerID}, Priority: in.Priority, TestingMode: in.TestingMode, SourceTaskID: in.SourceTaskID, DependsOn: in.DependsOn, SuggestedProperties: in.SuggestedProperties})
		if e != nil {
			return failure(e)
		}
		return output(s.summary(ctx, v))
	})
	addTool(s, &mcp.Tool{Name: "add_history", Title: "Add task comment", Annotations: annotations(false, false), Description: "Add an agent comment to an assigned task. To attach screenshots, first call upload_image and insert its returned Markdown link in content."}, func(ctx context.Context, _ *mcp.CallToolRequest, in historyInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentAddHistory(ctx, agentFrom(ctx), in.TaskID, in.Content)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	addTool(s, &mcp.Tool{Name: "add_artifact", Title: "Add task artifact", Annotations: annotations(false, false), Description: "Register a path or URL deliverable for an assigned task."}, func(ctx context.Context, _ *mcp.CallToolRequest, in artifactInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentAddArtifact(ctx, agentFrom(ctx), in.TaskID, core.ArtifactInput{Name: in.Name, Kind: in.Kind, Path: in.Path, URL: in.URL, Description: in.Description})
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	s.addMove("move_to_in_progress", "in_progress")
	s.addMove("move_to_testing", "testing")
	s.addMove("move_to_verification", "verification")
	addTool(s, &mcp.Tool{Name: "record_test_run", Title: "Record test run", Annotations: annotations(false, false), InputSchema: enumSchema[testInput](map[string][]any{"status": {"passed", "failed", "skipped", "manual_required"}}), Description: "Record an agent test result. status: passed, failed, skipped, manual_required. runner and type are free-form labels. Agent results never count as human testing."}, func(ctx context.Context, _ *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentRecordTest(ctx, agentFrom(ctx), in.TaskID, core.TestRunInput{Runner: in.Runner, Type: in.Type, Command: in.Command, Status: in.Status, Summary: in.Summary, OutputExcerpt: in.OutputExcerpt, DurationMS: in.DurationMS})
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	addTool(s, &mcp.Tool{Name: "report_usage", Title: "Report usage", Annotations: annotations(false, false), Description: "Report usage for a task. MCPCalls is a server-controlled cumulative SESSION snapshot including this call, not a per-task count. Do not sum snapshots. WorkerSessionID identifies the session."}, func(ctx context.Context, _ *mcp.CallToolRequest, in usageInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		var calls int64
		_ = s.Service.Store.DB.QueryRowContext(ctx, `SELECT mcp_calls FROM worker_sessions WHERE id=?`, agentFrom(ctx).Session.ID).Scan(&calls)
		v, e := s.Service.AgentReportUsage(ctx, agentFrom(ctx), in.TaskID, core.UsageInput{Provider: in.Provider, Model: in.Model, InputTokens: in.InputTokens, OutputTokens: in.OutputTokens, ReasoningTokens: in.ReasoningTokens, CachedTokens: in.CachedTokens, ReasoningEffort: in.ReasoningEffort, ModelCalls: in.ModelCalls, ToolCalls: in.ToolCalls, WallTimeMS: in.WallTimeMS, Source: in.Source})
		if e != nil {
			return failure(e)
		}
		v.MCPCalls = &calls
		_, _ = s.Service.Store.DB.ExecContext(ctx, `UPDATE usage_events SET mcp_calls=? WHERE id=?`, calls, v.ID)
		return output(struct {
			core.UsageEvent
			MCPCallsScope string `json:"mcp_calls_scope"`
		}{v, "session_total_including_this_call"})
	})
}

var moveDescriptions = map[string]string{
	"in_progress":  "Start work on an assigned task: features → in_progress. Call before implementing.",
	"testing":      "Hand off finished implementation: in_progress → testing. Then run checks and record_test_run. For human/hybrid testing, add handoff instructions with add_history, then finish the session or continue another task; never wait or poll.",
	"verification": "Request final human acceptance: testing → verification. Requires a passed record_test_run from this worker and, for human/hybrid tasks, a passed human test. On HUMAN_TEST_PENDING leave the task in testing; never wait, poll or retry.",
}

func (s *Server) addMove(name, target string) {
	addTool(s, &mcp.Tool{Name: name, Title: "Move task to " + target, Annotations: annotations(false, false), Description: moveDescriptions[target]}, func(ctx context.Context, _ *mcp.CallToolRequest, in taskIDInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentMove(ctx, agentFrom(ctx), in.TaskID, target)
		if e != nil {
			return failure(e)
		}
		return output(s.summary(ctx, v))
	})
}
