package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type WorkspaceContext struct {
	LaunchDirectory string        `json:"launch_directory"`
	LaunchProject   *core.Project `json:"launch_project,omitempty"`
	AllowMismatch   bool          `json:"allow_project_mismatch"`
}

type workspaceSession struct {
	rootsLoaded bool
	roots       []*mcp.Root
	approved    string
}

type workspaceGuard struct {
	sync.Mutex
	sessions map[*mcp.ServerSession]*workspaceSession
}

// addTool applies project-context checks to every worker capability, including
// attachment uploads, before the handler can read or mutate task data.
func addTool[In any](s *Server, tool *mcp.Tool, handler mcp.ToolHandlerFor[In, any]) {
	mcp.AddTool(s.MCP, tool, func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		if result := s.checkWorkspace(ctx, req); result != nil {
			return result, nil, nil
		}
		return handler(ctx, req, in)
	})
}

func (s *Server) invalidateWorkspace(session *mcp.ServerSession) {
	s.guard.Lock()
	defer s.guard.Unlock()
	delete(s.guard.sessions, session)
}

func workspaceFailure(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: message}}}
}

func (s *Server) checkWorkspace(ctx context.Context, req *mcp.CallToolRequest) *mcp.CallToolResult {
	if s.Workspace == nil || s.Workspace.AllowMismatch {
		return nil
	}
	s.guard.Lock()
	state := s.guard.sessions[req.Session]
	if state == nil {
		state = &workspaceSession{}
		s.guard.sessions[req.Session] = state
	}
	loaded, roots, approved := state.rootsLoaded, state.roots, state.approved
	s.guard.Unlock()
	caps := req.ClientCapabilities()
	if caps != nil && caps.RootsV2 != nil && !loaded {
		response, pending, err := workspaceInput(ctx, req, "agentboard_workspace_roots", &mcp.ListRootsParams{})
		if err != nil {
			return workspaceFailure("WORKSPACE_CONTEXT_UNAVAILABLE: cannot read client roots: " + err.Error())
		}
		if pending != nil {
			return pending
		}
		list, ok := response.(*mcp.ListRootsResult)
		if !ok {
			return workspaceFailure("WORKSPACE_CONTEXT_UNAVAILABLE: invalid client roots response")
		}
		roots = list.Roots
		s.guard.Lock()
		state.roots, state.rootsLoaded = roots, true
		s.guard.Unlock()
	}
	var conflicts []string
	if p := s.Workspace.LaunchProject; p != nil && p.ID != s.Agent.Project.ID {
		conflicts = append(conflicts, fmt.Sprintf("MCP launch directory %q belongs to project %q (%s)", s.Workspace.LaunchDirectory, p.Name, p.Path))
	}
	if len(roots) > 0 {
		matches := false
		var paths []string
		for _, root := range roots {
			path, err := rootPath(root.URI)
			if err != nil {
				return workspaceFailure("WORKSPACE_CONTEXT_UNAVAILABLE: " + err.Error())
			}
			paths = append(paths, path)
			p, err := s.Service.ProjectForDirectory(ctx, path)
			if err == nil && p.ID == s.Agent.Project.ID {
				matches = true
			}
		}
		if !matches {
			sort.Strings(paths)
			conflicts = append(conflicts, fmt.Sprintf("client workspace roots are %q", paths))
		}
	}
	if len(conflicts) == 0 {
		return nil
	}
	message := fmt.Sprintf("PROJECT_CONTEXT_MISMATCH: MCP worker %q is bound to project %q (%s), but %s. Continue using this explicitly configured project? To work on the workspace project instead, reconnect with its --project and --worker. The HTTP server's launch directory does not select the MCP project.", s.Agent.Worker.Slug, s.Agent.Project.Name, s.Agent.Project.Path, strings.Join(conflicts, "; "))
	if approved == message {
		return nil
	}
	if caps == nil || caps.Elicitation == nil || (caps.Elicitation.Form == nil && caps.Elicitation.URL != nil) {
		return workspaceFailure(message + " This client cannot ask for confirmation. Correct the configuration or explicitly pass --allow-project-mismatch after confirming the intended project. No task data was accessed.")
	}
	response, pending, err := workspaceInput(ctx, req, "agentboard_confirm_project", &mcp.ElicitParams{
		Mode: "form", Message: message,
		RequestedSchema: json.RawMessage(`{"type":"object","properties":{"confirm":{"type":"boolean","title":"Use the configured AgentBoard project"}},"required":["confirm"]}`),
	})
	if err != nil {
		return workspaceFailure(message + " Confirmation failed: " + err.Error())
	}
	if pending != nil {
		return pending
	}
	answer, ok := response.(*mcp.ElicitResult)
	if !ok || answer.Action != "accept" || answer.Content["confirm"] != true {
		return workspaceFailure(message + " Confirmation was declined or cancelled. No task data was accessed.")
	}
	s.guard.Lock()
	state.approved = message
	s.guard.Unlock()
	return nil
}

// Use multi-round-trip input requests on new protocol versions, and the SDK's
// legacy server-to-client APIs when those requests are supported by the client.
func workspaceInput(ctx context.Context, req *mcp.CallToolRequest, key string, params mcp.InputRequest) (mcp.InputResponse, *mcp.CallToolResult, error) {
	if response := req.Params.InputResponses[key]; response != nil {
		return response, nil, nil
	}
	if req.ProtocolVersion() >= "2026-07-28" {
		return nil, &mcp.CallToolResult{InputRequests: mcp.InputRequestMap{key: params}}, nil
	}
	switch p := params.(type) {
	case *mcp.ListRootsParams:
		v, err := req.Session.ListRoots(ctx, p)
		return v, nil, err
	case *mcp.ElicitParams:
		v, err := req.Session.Elicit(ctx, p)
		return v, nil, err
	default:
		return nil, nil, fmt.Errorf("unsupported workspace request")
	}
}

func rootPath(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme != "file" || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("invalid workspace root URI %q; expected a file URI", uri)
	}
	path := u.Path
	if u.Host != "" && u.Host != "localhost" {
		path = "//" + u.Host + path
	} else if runtime.GOOS == "windows" && len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}
	path = filepath.FromSlash(path)
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("workspace root %q is not an absolute path", uri)
	}
	return filepath.Clean(path), nil
}
