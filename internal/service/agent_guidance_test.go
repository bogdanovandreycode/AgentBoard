package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

func TestAgentShortTaskIDsRespectVisibilityAndCanonicalWrites(t *testing.T) {
	f := setup(t)
	task := createTask(t, f, f.a, "features", "ai")
	prefix := strings.ToUpper(task.ID[:8])
	details, err := f.svc.AgentTask(f.ctx, f.ac, prefix)
	if err != nil || details.ID != task.ID {
		t.Fatalf("resolve: %+v %v", details, err)
	}
	history, err := f.svc.AgentAddHistory(f.ctx, f.ac, prefix, "short ID comment")
	if err != nil || history.TaskID != task.ID {
		t.Fatalf("history: %+v %v", history, err)
	}
	artifact, err := f.svc.AgentAddArtifact(f.ctx, f.ac, prefix, core.ArtifactInput{Name: "Report", Path: "report.txt"})
	if err != nil || artifact.TaskID != task.ID {
		t.Fatalf("artifact: %+v %v", artifact, err)
	}
	run, err := f.svc.AgentRecordTest(f.ctx, f.ac, prefix, core.TestRunInput{Status: "passed"})
	if err != nil || run.TaskID != task.ID {
		t.Fatalf("test: %+v %v", run, err)
	}
	usage, err := f.svc.AgentReportUsage(f.ctx, f.ac, prefix, core.UsageInput{})
	if err != nil || usage.TaskID != task.ID {
		t.Fatalf("usage: %+v %v", usage, err)
	}
	if _, err = f.svc.AgentMove(f.ctx, f.ac, prefix, "in_progress"); err != nil {
		t.Fatal(err)
	}
	details, err = f.svc.AgentTask(f.ctx, f.ac, prefix)
	if err != nil || len(details.History) < 2 || len(details.TestRuns) != 1 || len(details.Artifacts) != 1 || len(details.Usage) != 1 {
		t.Fatalf("canonical context missing: %+v %v", details, err)
	}
	for _, state := range []string{"backlog", "complete"} {
		hidden := createTask(t, f, f.a, state, "ai")
		if _, err := f.svc.ResolveAgentTaskID(f.ctx, f.ac, hidden.ID[:8]); !errors.Is(err, core.ErrTaskNotAccessible) {
			t.Fatalf("hidden %s: %v", state, err)
		}
	}
	foreign := createTask(t, f, f.b, "features", "ai")
	if _, err := f.svc.ResolveAgentTaskID(f.ctx, f.ac, foreign.ID[:8]); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("foreign: %v", err)
	}
	for _, ref := range []string{task.ID[:7], "%", "________", "not-a-task"} {
		if _, err := f.svc.ResolveAgentTaskID(f.ctx, f.ac, ref); !errors.Is(err, core.ErrTaskNotAccessible) {
			t.Fatalf("invalid ref %q: %v", ref, err)
		}
	}
}

func TestAgentShortTaskIDAmbiguityIgnoresHiddenMatches(t *testing.T) {
	f := setup(t)
	visible := createTask(t, f, f.a, "features", "ai")
	hidden := createTask(t, f, f.b, "features", "ai")
	// Only temporary fixture data is edited to create deterministic UUID collisions.
	ids := []string{"12345678-0000-4000-8000-000000000001", "12345678-0000-4000-8000-000000000002"}
	for i, task := range []core.Task{visible, hidden} {
		if _, err := f.store.DB.ExecContext(f.ctx, `DELETE FROM history_entries WHERE task_id=?`, task.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := f.store.DB.ExecContext(f.ctx, `UPDATE tasks SET id=? WHERE id=?`, ids[i], task.ID); err != nil {
			t.Fatal(err)
		}
	}
	if id, err := f.svc.ResolveAgentTaskID(f.ctx, f.ac, "12345678"); err != nil || id != ids[0] {
		t.Fatalf("hidden collision leaked: %s %v", id, err)
	}
	if _, err := f.store.DB.ExecContext(f.ctx, `UPDATE tasks SET assignee_worker_id=? WHERE id=?`, f.a.ID, ids[1]); err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.ResolveAgentTaskID(f.ctx, f.ac, "12345678"); !errors.Is(err, core.ErrAmbiguousTaskID) {
		t.Fatalf("ambiguity: %v", err)
	}
	if id, err := f.svc.ResolveAgentTaskID(f.ctx, f.ac, ids[0]); err != nil || id != ids[0] {
		t.Fatalf("full ID: %s %v", id, err)
	}
}

func TestHumanTestingHandoffKeepsWorkflowGate(t *testing.T) {
	for _, mode := range []string{"human", "hybrid"} {
		t.Run(mode, func(t *testing.T) {
			f := setup(t)
			task := createTask(t, f, f.a, "testing", mode)
			workflow := f.svc.AgentWorkflow(f.ctx, f.ac, task)
			if !workflow.HumanTestingPending || workflow.NextTool != "" || !strings.Contains(workflow.Instructions, "finish the session") {
				t.Fatalf("handoff: %+v", workflow)
			}
			if _, err := f.svc.AgentMove(f.ctx, f.ac, task.ID, "verification"); !errors.Is(err, core.ErrHumanTestPending) {
				t.Fatalf("human gate bypassed: %v", err)
			}
			if _, err := f.svc.AgentRecordTest(f.ctx, f.ac, task.ID, core.TestRunInput{Status: "passed"}); err != nil {
				t.Fatal(err)
			}
			if _, err := f.svc.AgentMove(f.ctx, f.ac, task.ID, "verification"); !errors.Is(err, core.ErrHumanTestPending) {
				t.Fatalf("AI impersonated human: %v", err)
			}
		})
	}
}
