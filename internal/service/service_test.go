package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
)

type fixture struct {
	ctx     context.Context
	store   *persistence.Store
	svc     *Service
	project core.Project
	a, b    core.Worker
	ac, bc  core.AgentContext
}

func setup(t *testing.T) fixture {
	t.Helper()
	ctx := context.Background()
	st, err := persistence.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	p, err := st.EnsureProject(ctx, "test", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	a, err := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Codex", Slug: "codex", Kind: "external", Capabilities: `{"coding":true}`, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	b, err := st.CreateWorker(ctx, p.ID, core.WorkerInput{Name: "Gemma", Slug: "gemma", Kind: "external", Capabilities: `{"image_generation":true}`, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	svc := New(st)
	as, _ := st.StartSession(ctx, p, a, "test")
	bs, _ := st.StartSession(ctx, p, b, "test")
	return fixture{ctx, st, svc, p, a, b, core.AgentContext{Project: p, Worker: a, Session: as}, core.AgentContext{Project: p, Worker: b, Session: bs}}
}
func workerAssignee(id string) core.Assignee { return core.Assignee{Type: "worker", WorkerID: &id} }
func createTask(t *testing.T, f fixture, w core.Worker, state, mode string) core.Task {
	t.Helper()
	task, err := f.store.CreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Task for " + w.Name, State: state, Assignee: workerAssignee(w.ID), TestingMode: mode})
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func TestWorkerVisibilityAndDirectAccess(t *testing.T) {
	f := setup(t)
	ta := createTask(t, f, f.a, "features", "ai")
	tb := createTask(t, f, f.b, "features", "ai")
	board, err := f.svc.AgentBoard(f.ctx, f.ac)
	if err != nil {
		t.Fatal(err)
	}
	if len(board["features"]) != 1 || board["features"][0].ID != ta.ID {
		t.Fatalf("wrong scoped board: %#v", board)
	}
	if _, err = f.svc.AgentTask(f.ctx, f.ac, tb.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("get foreign: %v", err)
	}
	if _, err = f.svc.AgentAddHistory(f.ctx, f.ac, tb.ID, "no"); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("history foreign: %v", err)
	}
	if _, err = f.svc.AgentAddArtifact(f.ctx, f.ac, tb.ID, core.ArtifactInput{Name: "x", Path: "x"}); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("artifact foreign: %v", err)
	}
	if _, err = f.svc.AgentRecordTest(f.ctx, f.ac, tb.ID, core.TestRunInput{Status: "passed"}); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("test foreign: %v", err)
	}
	if _, err = f.svc.AgentReportUsage(f.ctx, f.ac, tb.ID, core.UsageInput{}); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("usage foreign: %v", err)
	}
	if _, err = f.svc.AgentMove(f.ctx, f.ac, tb.ID, "in_progress"); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("move foreign: %v", err)
	}
}

func TestAgentCreatedFeatureAndHumanAssignment(t *testing.T) {
	f := setup(t)
	source := createTask(t, f, f.a, "in_progress", "ai")
	feature, err := f.svc.CreateFeature(f.ctx, f.ac, core.FeatureInput{Title: "Create portraits", Reason: "Needed", Assignee: workerAssignee(f.b.ID), SourceTaskID: &source.ID})
	if err != nil {
		t.Fatal(err)
	}
	if feature.State != "features" || feature.CreatedByType != "agent" || feature.CreatedByWorkerID == nil || *feature.CreatedByWorkerID != f.a.ID || feature.AssigneeWorkerID == nil || *feature.AssigneeWorkerID != f.b.ID {
		t.Fatalf("bad feature: %#v", feature)
	}
	ab, _ := f.svc.AgentBoard(f.ctx, f.ac)
	bb, _ := f.svc.AgentBoard(f.ctx, f.bc)
	if len(ab["features"]) != 0 || len(bb["features"]) != 1 {
		t.Fatalf("visibility A=%d B=%d", len(ab["features"]), len(bb["features"]))
	}
	human, err := f.svc.CreateFeature(f.ctx, f.ac, core.FeatureInput{Title: "Human follow-up", Assignee: core.Assignee{Type: "human"}})
	if err != nil {
		t.Fatal(err)
	}
	if human.AssigneeType != "human" {
		t.Fatal("not human")
	}
	ab, _ = f.svc.AgentBoard(f.ctx, f.ac)
	if len(ab["features"]) != 0 {
		t.Fatal("human task leaked")
	}
}

func TestTransitionsAndTestingPolicies(t *testing.T) {
	for _, tc := range []struct {
		mode      string
		wantErr   error
		humanPass bool
	}{{"ai", nil, false}, {"human", core.ErrHumanTestPending, true}, {"hybrid", core.ErrHumanTestPending, false}, {"hybrid", nil, true}} {
		t.Run(tc.mode, func(t *testing.T) {
			f := setup(t)
			task := createTask(t, f, f.a, "features", tc.mode)
			var err error
			task, err = f.svc.AgentMove(f.ctx, f.ac, task.ID, "in_progress")
			if err != nil {
				t.Fatal(err)
			}
			task, err = f.svc.AgentMove(f.ctx, f.ac, task.ID, "testing")
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.svc.AgentRecordTest(f.ctx, f.ac, task.ID, core.TestRunInput{Runner: "go", Type: "automated", Status: "passed"})
			if err != nil {
				t.Fatal(err)
			}
			if tc.humanPass {
				_, _ = f.svc.HumanRecordTest(f.ctx, task.ID, core.TestRunInput{Runner: "human", Type: "manual", Status: "passed"})
			}
			_, err = f.svc.AgentMove(f.ctx, f.ac, task.ID, "verification")
			if tc.wantErr == nil && err != nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestDisabledWorkerAndSessionAttribution(t *testing.T) {
	f := setup(t)
	_, err := f.store.UpdateWorker(f.ctx, f.a.ID, core.WorkerInput{Name: f.a.Name, Slug: f.a.Slug, Kind: f.a.Kind, Capabilities: f.a.Capabilities, Enabled: false})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.svc.BeginAgentSession(f.ctx, f.project.Path, f.a.Slug, "test"); !errors.Is(err, core.ErrWorkerDisabled) {
		t.Fatalf("expected disabled, got %v", err)
	}
	task := createTask(t, f, f.b, "in_progress", "ai")
	h, err := f.svc.AgentAddHistory(f.ctx, f.bc, task.ID, "working")
	if err != nil {
		t.Fatal(err)
	}
	if h.WorkerID == nil || *h.WorkerID != f.b.ID || h.WorkerSessionID == nil || *h.WorkerSessionID != f.bc.Session.ID {
		t.Fatalf("bad attribution: %#v", h)
	}
}

func TestArtifactsAndDependencyOutputs(t *testing.T) {
	f := setup(t)
	graphic := createTask(t, f, f.b, "verification", "ai")
	art, err := f.svc.AgentAddArtifact(f.ctx, f.bc, graphic.ID, core.ArtifactInput{Name: "Portrait", Kind: "image", Path: "assets/portrait.png"})
	if err != nil {
		t.Fatal(err)
	}
	consumer, err := f.store.CreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Use portrait", State: "features", Assignee: workerAssignee(f.a.ID), TestingMode: "ai", DependencyIDs: []string{graphic.ID}})
	if err != nil {
		t.Fatal(err)
	}
	detail, err := f.svc.AgentTask(f.ctx, f.ac, consumer.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.Dependencies) != 1 || len(detail.Dependencies[0].Artifacts) != 1 || detail.Dependencies[0].Artifacts[0].ID != art.ID {
		t.Fatalf("missing output: %#v", detail.Dependencies)
	}
}

func TestCustomPropertyVisibility(t *testing.T) {
	f := setup(t)
	hidden, err := f.svc.CreatePropertyDefinition(f.ctx, f.project.ID, core.PropertyDefinitionInput{Name: "Budget", Type: "number", Visibility: "human_only"})
	if err != nil {
		t.Fatal(err)
	}
	visible, err := f.svc.CreatePropertyDefinition(f.ctx, f.project.ID, core.PropertyDefinitionInput{Name: "Asset style", Type: "text", Visibility: "agent_read_write"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := f.store.CreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Properties", State: "features", Assignee: workerAssignee(f.a.ID), TestingMode: "ai", Properties: map[string]string{hidden.ID: "100", visible.ID: "low-poly"}})
	if err != nil {
		t.Fatal(err)
	}
	agent, err := f.svc.AgentTask(f.ctx, f.ac, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(agent.Properties) != 1 || agent.Properties[0].Name != "Asset style" {
		t.Fatalf("agent properties: %#v", agent.Properties)
	}
	human, err := f.svc.HumanTask(f.ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(human.Properties) != 2 {
		t.Fatalf("human properties: %#v", human.Properties)
	}
}

func TestPropertyDefaultsAndValidation(t *testing.T) {
	f := setup(t)
	definition, err := f.svc.CreatePropertyDefinition(f.ctx, f.project.ID, core.PropertyDefinitionInput{
		Name: "Ticket", Type: "text", Visibility: "agent_read", Placeholder: "ABC-123", Regex: `^[A-Z]+-[0-9]+$`, DefaultValue: "ABC-123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if definition.Placeholder != "ABC-123" || definition.Regex == "" {
		t.Fatalf("definition fields: %#v", definition)
	}
	task, err := f.svc.HumanCreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Default", State: "backlog", Assignee: core.Assignee{Type: "unassigned"}, TestingMode: "ai"})
	if err != nil {
		t.Fatal(err)
	}
	details, err := f.svc.HumanTask(f.ctx, task.ID)
	if err != nil || len(details.Properties) != 1 || details.Properties[0].Value != "ABC-123" {
		t.Fatalf("default value: %#v, %v", details.Properties, err)
	}
	invalid := map[string]string{definition.ID: "wrong"}
	if _, err = f.svc.HumanUpdateTask(f.ctx, task.ID, core.TaskUpdate{Properties: &invalid}); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("invalid update: %v", err)
	}
	if _, err = f.svc.HumanCreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Invalid", State: "backlog", Assignee: core.Assignee{Type: "unassigned"}, TestingMode: "ai", Properties: invalid}); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("invalid create: %v", err)
	}
	valid := map[string]string{definition.ID: "XYZ-9"}
	if _, err = f.svc.HumanUpdateTask(f.ctx, task.ID, core.TaskUpdate{Properties: &valid}); err != nil {
		t.Fatal(err)
	}
	details, err = f.svc.HumanTask(f.ctx, task.ID)
	if err != nil || details.Properties[0].Value != "XYZ-9" {
		t.Fatalf("updated value: %#v, %v", details.Properties, err)
	}
	if _, err = f.svc.CreatePropertyDefinition(f.ctx, f.project.ID, core.PropertyDefinitionInput{Name: "Broken", Type: "text", Visibility: "human_only", Regex: "["}); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("invalid regex: %v", err)
	}
}
