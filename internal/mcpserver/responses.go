package mcpserver

import (
	"context"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/service"
	"github.com/google/jsonschema-go/jsonschema"
)

type taskSummary struct {
	ID               string
	Title            string
	State            string
	Priority         string
	TestingMode      string
	AssigneeType     string
	AssigneeWorkerID *string `json:",omitempty"`
	Workflow         service.AgentWorkflow
}

func (s *Server) summary(ctx context.Context, t core.Task) taskSummary {
	return taskSummary{t.ID, t.Title, t.State, t.Priority, t.TestingMode, t.AssigneeType, t.AssigneeWorkerID, s.Service.AgentWorkflow(ctx, agentFrom(ctx), t)}
}

type sessionInfo struct {
	Project       string `json:"project"`
	ProjectID     string `json:"project_id"`
	ProjectPath   string `json:"project_path"`
	Worker        string `json:"worker"`
	WorkerSlug    string `json:"worker_slug"`
	WorkerID      string `json:"worker_id"`
	SessionID     string `json:"worker_session_id"`
	BindingSource string `json:"binding_source,omitempty"`
}

type boardResponse struct {
	Session  sessionInfo              `json:"session"`
	Board    map[string][]taskSummary `json:"board"`
	NextStep string                   `json:"next_step"`
}

func (s *Server) sessionInfo(a core.AgentContext) sessionInfo {
	s.agentMu.Lock()
	source := s.bindingSource
	s.agentMu.Unlock()
	return sessionInfo{a.Project.Name, a.Project.ID, a.Project.Path, a.Worker.Name, a.Worker.Slug, a.Worker.ID, a.Session.ID, source}
}

func boardNextStep(board map[string][]taskSummary) string {
	switch {
	case len(board["in_progress"]) > 0:
		return "Continue the in_progress task: call get_task for its full description and history."
	case len(board["features"]) > 0:
		return "Pick a feature, call get_task to read it, then move_to_in_progress before starting work."
	case len(board["testing"]) > 0 || len(board["verification"]) > 0:
		return "Only testing/verification tasks remain. Finish AI testing where the workflow allows; never wait or poll for human testing. If nothing is actionable, tell the user."
	default:
		return "No tasks are assigned to this worker in this project. Tell the user; a human assigns Features in the AgentBoard web UI. You may propose work with create_feature only if the user asks."
	}
}

func enumSchema[T any](fields map[string][]any) *jsonschema.Schema {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(err)
	}
	for field, values := range fields {
		schema.Properties[field].Enum = values
	}
	return schema
}
