package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	MCP     *mcp.Server
	Service *service.Service
	Agent   core.AgentContext
}

func New(s *service.Service, a core.AgentContext, version string) *Server {
	x := &Server{Service: s, Agent: a}
	x.MCP = mcp.NewServer(&mcp.Implementation{Name: "agentboard", Version: version}, nil)
	x.register()
	return x
}
func (s *Server) Run(ctx context.Context, transport mcp.Transport) error {
	return s.MCP.Run(ctx, transport)
}
func (s *Server) count(ctx context.Context) {
	_, _ = s.Service.Store.DB.ExecContext(ctx, `UPDATE worker_sessions SET mcp_calls=mcp_calls+1 WHERE id=?`, s.Agent.Session.ID)
}
func output(v any) (*mcp.CallToolResult, any, error) {
	b, _ := json.Marshal(v)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}, StructuredContent: v}, nil, nil
}
func failure(err error) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}}}, nil, nil
}

type taskIDInput struct {
	TaskID string `json:"task_id" jsonschema:"Task ID"`
}
type historyInput struct {
	TaskID  string `json:"task_id"`
	Content string `json:"content"`
}
type artifactInput struct {
	TaskID      string `json:"task_id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
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
	Type          string `json:"type"`
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
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "get_my_board", Description: "List only tasks assigned to the current Worker in AI-visible workflow states."}, func(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentBoard(ctx, s.Agent)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "get_task", Description: "Get details for a task assigned to the current Worker."}, func(ctx context.Context, _ *mcp.CallToolRequest, in taskIDInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentTask(ctx, s.Agent, in.TaskID)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "list_workers", Description: "List enabled Workers available when assigning a newly-created Feature."}, func(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.Store.ListWorkers(ctx, s.Agent.Project.ID, true)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "create_feature", Description: "Create a Feature attributed to the current Worker and optionally assign it."}, func(ctx context.Context, _ *mcp.CallToolRequest, in featureInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.CreateFeature(ctx, s.Agent, core.FeatureInput{Title: in.Title, Description: in.Description, Reason: in.Reason, Assignee: core.Assignee{Type: in.AssigneeType, WorkerID: in.AssigneeWorkerID}, Priority: in.Priority, TestingMode: in.TestingMode, SourceTaskID: in.SourceTaskID, DependsOn: in.DependsOn, SuggestedProperties: in.SuggestedProperties})
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "add_history", Description: "Add an agent comment to an assigned task."}, func(ctx context.Context, _ *mcp.CallToolRequest, in historyInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentAddHistory(ctx, s.Agent, in.TaskID, in.Content)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "add_artifact", Description: "Register a path or URL deliverable for an assigned task."}, func(ctx context.Context, _ *mcp.CallToolRequest, in artifactInput) (*mcp.CallToolResult, any, error) {
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
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "record_test_run", Description: "Record a structured automated test result for an assigned task."}, func(ctx context.Context, _ *mcp.CallToolRequest, in testInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentRecordTest(ctx, s.Agent, in.TaskID, core.TestRunInput{Runner: in.Runner, Type: in.Type, Command: in.Command, Status: in.Status, Summary: in.Summary, OutputExcerpt: in.OutputExcerpt, DurationMS: in.DurationMS})
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
	mcp.AddTool(s.MCP, &mcp.Tool{Name: "report_usage", Description: "Report model and tool usage for an assigned task. Worker identity and MCP call count are server-controlled."}, func(ctx context.Context, _ *mcp.CallToolRequest, in usageInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		var calls int64
		_ = s.Service.Store.DB.QueryRowContext(ctx, `SELECT mcp_calls FROM worker_sessions WHERE id=?`, s.Agent.Session.ID).Scan(&calls)
		v, e := s.Service.AgentReportUsage(ctx, s.Agent, in.TaskID, core.UsageInput{Provider: in.Provider, Model: in.Model, InputTokens: in.InputTokens, OutputTokens: in.OutputTokens, ReasoningTokens: in.ReasoningTokens, CachedTokens: in.CachedTokens, ReasoningEffort: in.ReasoningEffort, ModelCalls: in.ModelCalls, ToolCalls: in.ToolCalls, WallTimeMS: in.WallTimeMS, Source: in.Source})
		if e != nil {
			return failure(e)
		}
		v.MCPCalls = &calls
		_, _ = s.Service.Store.DB.ExecContext(ctx, `UPDATE usage_events SET mcp_calls=? WHERE id=?`, calls, v.ID)
		return output(v)
	})
}
func (s *Server) addMove(name, target string) {
	mcp.AddTool(s.MCP, &mcp.Tool{Name: name, Description: "Move an assigned task forward to " + target + "."}, func(ctx context.Context, _ *mcp.CallToolRequest, in taskIDInput) (*mcp.CallToolResult, any, error) {
		s.count(ctx)
		v, e := s.Service.AgentMove(ctx, s.Agent, in.TaskID, target)
		if e != nil {
			return failure(e)
		}
		return output(v)
	})
}
