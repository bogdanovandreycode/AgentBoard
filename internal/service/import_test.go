package service

import (
	"errors"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

func TestHumanImportTasksIsAtomicAndResolvesProjectNames(t *testing.T) {
	f := setup(t)
	property, err := f.svc.CreatePropertyDefinition(f.ctx, f.project.ID, core.PropertyDefinitionInput{Name: "Estimate", Type: "text", Visibility: "agent_read"})
	if err != nil {
		t.Fatal(err)
	}
	doc := ImportDocument{Version: 1, Tasks: []ImportTask{
		{Key: "first", Title: "Imported first", State: "features", Assignee: ImportAssignee{Type: "worker", Worker: f.a.Slug}, Properties: map[string]string{"Estimate": "3 days"}},
		{Key: "second", Title: "Imported second", DependsOn: []string{"first"}},
	}}
	count, err := f.svc.HumanImportTasks(f.ctx, f.project.ID, doc)
	if err != nil || count != 2 {
		t.Fatalf("import: count=%d err=%v", count, err)
	}
	board, err := f.svc.HumanBoard(f.ctx, f.project.ID)
	if err != nil {
		t.Fatal(err)
	}
	first := board["features"][0]
	second := board["backlog"][0]
	if first.AssigneeWorkerID == nil || *first.AssigneeWorkerID != f.a.ID || len(second.Dependencies) != 1 || second.Dependencies[0].DependsOnTaskID != first.ID {
		t.Fatalf("worker or dependency unresolved: first=%+v second=%+v", first, second)
	}
	details, err := f.svc.HumanTask(f.ctx, first.ID)
	if err != nil || len(details.Properties) != 1 || details.Properties[0].PropertyDefinitionID != property.ID || details.Properties[0].Value != "3 days" || len(details.History) != 1 || details.History[0].ActorType != "system" {
		t.Fatalf("property or history unresolved: details=%+v err=%v", details, err)
	}
	doc.Tasks[1].Properties = map[string]string{"Unknown": "value"}
	_, err = f.svc.HumanImportTasks(f.ctx, f.project.ID, doc)
	if !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
	board, err = f.svc.HumanBoard(f.ctx, f.project.ID)
	if err != nil || len(board["features"]) != 1 || len(board["backlog"]) != 1 {
		t.Fatalf("invalid file changed board: %+v err=%v", board, err)
	}
}
