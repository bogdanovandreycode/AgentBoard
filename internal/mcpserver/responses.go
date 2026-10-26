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
	return taskSummary{t.ID, t.Title, t.State, t.Priority, t.TestingMode, t.AssigneeType, t.AssigneeWorkerID, s.Service.AgentWorkflow(ctx, s.Agent, t)}
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
