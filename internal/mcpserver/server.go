package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	MCP       *mcp.Server
	Service   *service.Service
	Agent     core.AgentContext
	Workspace *WorkspaceContext
	guard     workspaceGuard
}

const instructions = `AgentBoard is the task source of truth for this worker.
At the start of a work session call get_my_board. Continue in_progress tasks first;
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

func New(s *service.Service, a core.AgentContext, version string) *Server {
	x := &Server{Service: s, Agent: a, guard: workspaceGuard{sessions: make(map[*mcp.ServerSession]*workspaceSession)}}
	x.MCP = mcp.NewServer(&mcp.Implementation{Name: "agentboard", Version: version}, &mcp.ServerOptions{Instructions: instructions,
		RootsListChangedHandler: func(_ context.Context, req *mcp.RootsListChangedRequest) { x.invalidateWorkspace(req.Session) },
	})
	x.MCP.AddResource(&mcp.Resource{URI: "agentboard://current-worker", Name: "Current AgentBoard worker", MIMEType: "application/json"}, func(ctx context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		b, err := json.Marshal(map[string]any{"project": a.Project.Name, "project_id": a.Project.ID, "project_path": a.Project.Path, "worker": a.Worker.Name, "worker_id": a.Worker.ID, "workspace": x.Workspace, "instructions": instructions})
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: "agentboard://current-worker", MIMEType: "application/json", Text: string(b)}}}, nil
	})
	x.register()
	return x
}
func (s *Server) Run(ctx context.Context, transport mcp.Transport) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = s.Service.Store.TouchSession(ctx, s.Agent.Session.ID)
			}
		}
	}()
	return s.MCP.Run(ctx, transport)
}
func (s *Server) count(ctx context.Context) {
	_, _ = s.Service.Store.DB.ExecContext(ctx, `UPDATE worker_sessions SET mcp_calls=mcp_calls+1,last_activity_at=?,last_seen_at=? WHERE id=?`, core.Now(), core.Now(), s.Agent.Session.ID)
}
func output(v any) (*mcp.CallToolResult, any, error) {
	b, _ := json.Marshal(v)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil, nil
}
func failure(err error) (*mcp.CallToolResult, any, error) {
	if errors.Is(err, core.ErrHumanTestPending) {
		err = errors.New("HUMAN_TEST_PENDING: leave the task in testing, provide handoff instructions, and finish this session or continue another task. Do not wait, poll, or retry verification until human testing is resolved")
	}
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}}}, nil, nil
}

type taskIDInput struct {
	TaskID string `json:"task_id" jsonschema:"Full UUID or a unique case-insensitive prefix of at least 8 characters from get_my_board"`
}
type historyInput struct {
	TaskID  string `json:"task_id"`
	Content string `json:"content"`
}
type artifactInput struct {
	TaskID      string `json:"task_id"`
	Name        string `json:"name"`
	Kind        string `json:"kind" jsonschema:"Free-form artifact category, for example file, screenshot, url, build, or report"`
	Path        string `json:"path,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}
type featureInput struct {
	Title               string            `json:"title"`
	Description         string            `json:"description,omitempty"`
	Reason              string            `json:"reason"`
	AssigneeType        string            `json:"assignee_type"`
	AssigneeWorkerID    *string           `json:"assignee_worker_id,omitempty"`
	Priority            string            `json:"priority,omitempty"`
	TestingMode         string            `json:"testing_mode,omitempty"`
	SourceTaskID        *string           `json:"source_task_id,omitempty"`
	DependsOn           []string          `json:"depends_on,omitempty"`
	SuggestedProperties map[string]string `json:"suggested_properties,omitempty"`
}
type testInput struct {
	TaskID        string `json:"task_id"`
	Runner        string `json:"runner"`
	Type          string `json:"type" jsonschema:"Free-form test category, for example automated, unit, integration, or manual. Agent results never satisfy human testing."`
	Command       string `json:"command,omitempty"`
	Status        string `json:"status"`
	Summary       string `json:"summary"`
	OutputExcerpt string `json:"output_excerpt,omitempty"`
	DurationMS    *int64 `json:"duration_ms,omitempty"`
}
type usageInput struct {
	TaskID          string  `json:"task_id"`
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
	addTool(s, &mcp.Tool{Name: "get_my_board", Description: "List compact summaries of assigned AI-visible tasks. Use get_task for full details. Verify agentboard://current-worker matches your intended project."}, func(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentBoard(ctx, s.Agent)
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
		return output(board)
	})
	addTool(s, &mcp.Tool{Name: "get_task", Description: "Get assigned task details and view screenshots from its description, testing instructions and comments. Images are returned as native image content; use get_image for any omitted images."}, func(ctx context.Context, _ *mcp.CallToolRequest, in taskIDInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentTask(ctx, s.Agent, in.TaskID)
		if e != nil {
			return failure(e)
		}
		return s.taskOutput(ctx, v)
	})
	s.registerImages()
	s.registerFiles()
	addTool(s, &mcp.Tool{Name: "list_workers", Description: "List enabled Workers available when assigning a newly-created Feature."}, func(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.Store.ListWorkers(ctx, s.Agent.Project.ID, true)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	addTool(s, &mcp.Tool{Name: "create_feature", InputSchema: enumSchema[featureInput](map[string][]any{"assignee_type": {"unassigned", "human", "worker"}, "priority": {"critical", "high", "medium", "low"}, "testing_mode": {"ai", "human", "hybrid"}}), Description: "Create a Feature attributed to the current Worker and optionally assign it. source_task_id and depends_on require full UUIDs."}, func(ctx context.Context, _ *mcp.CallToolRequest, in featureInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.CreateFeature(ctx, s.Agent, core.FeatureInput{Title: in.Title, Description: in.Description, Reason: in.Reason, Assignee: core.Assignee{Type: in.AssigneeType, WorkerID: in.AssigneeWorkerID}, Priority: in.Priority, TestingMode: in.TestingMode, SourceTaskID: in.SourceTaskID, DependsOn: in.DependsOn, SuggestedProperties: in.SuggestedProperties})
		if e != nil {
			return failure(e)
		}
		return output(s.summary(ctx, v))
	})
	addTool(s, &mcp.Tool{Name: "add_history", Description: "Add an agent comment to an assigned task. To attach screenshots, first call upload_image and insert its returned Markdown link in content."}, func(ctx context.Context, _ *mcp.CallToolRequest, in historyInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentAddHistory(ctx, s.Agent, in.TaskID, in.Content)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	addTool(s, &mcp.Tool{Name: "add_artifact", Description: "Register a path or URL deliverable for an assigned task."}, func(ctx context.Context, _ *mcp.CallToolRequest, in artifactInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentAddArtifact(ctx, s.Agent, in.TaskID, core.ArtifactInput{Name: in.Name, Kind: in.Kind, Path: in.Path, URL: in.URL, Description: in.Description})
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	s.addMove("move_to_in_progress", "in_progress")
	s.addMove("move_to_testing", "testing")
	s.addMove("move_to_verification", "verification")
	addTool(s, &mcp.Tool{Name: "record_test_run", InputSchema: enumSchema[testInput](map[string][]any{"status": {"passed", "failed", "skipped", "manual_required"}}), Description: "Record an agent test result. status: passed, failed, skipped, manual_required. runner and type are free-form labels. Agent results never count as human testing."}, func(ctx context.Context, _ *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentRecordTest(ctx, s.Agent, in.TaskID, core.TestRunInput{Runner: in.Runner, Type: in.Type, Command: in.Command, Status: in.Status, Summary: in.Summary, OutputExcerpt: in.OutputExcerpt, DurationMS: in.DurationMS})
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	addTool(s, &mcp.Tool{Name: "report_usage", Description: "Report usage for a task. MCPCalls is a server-controlled cumulative SESSION snapshot including this call, not a per-task count. Do not sum snapshots. WorkerSessionID identifies the session."}, func(ctx context.Context, _ *mcp.CallToolRequest, in usageInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		var calls int64
		_ = s.Service.Store.DB.QueryRowContext(ctx, `SELECT mcp_calls FROM worker_sessions WHERE id=?`, s.Agent.Session.ID).Scan(&calls)
		v, e := s.Service.AgentReportUsage(ctx, s.Agent, in.TaskID, core.UsageInput{Provider: in.Provider, Model: in.Model, InputTokens: in.InputTokens, OutputTokens: in.OutputTokens, ReasoningTokens: in.ReasoningTokens, CachedTokens: in.CachedTokens, ReasoningEffort: in.ReasoningEffort, ModelCalls: in.ModelCalls, ToolCalls: in.ToolCalls, WallTimeMS: in.WallTimeMS, Source: in.Source})
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
func (s *Server) addMove(name, target string) {
	addTool(s, &mcp.Tool{Name: name, Description: "Move an assigned task forward to " + target + ". Leave pending human/hybrid tasks in testing, provide handoff instructions, then finish the session or continue another task; never wait or poll."}, func(ctx context.Context, _ *mcp.CallToolRequest, in taskIDInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentMove(ctx, s.Agent, in.TaskID, target)
		if e != nil {
			return failure(e)
		}
		return output(s.summary(ctx, v))
	})
}
