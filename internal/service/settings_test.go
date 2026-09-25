package service

import (
	"errors"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

func TestCustomColumnStaysHumanOnlyAndReturnsTasksToBacklog(t *testing.T) {
	f := setup(t)
	settings, err := f.svc.GetProjectSettings(f.ctx, f.project.ID)
	if err != nil { t.Fatal(err) }
	settings.Columns = append(settings.Columns, BoardColumn{ID: "custom-review", Name: "Needs review"})
	if _, err = f.svc.UpdateProjectSettings(f.ctx, f.project.ID, settings); err != nil { t.Fatal(err) }
	task := createTask(t, f, f.a, "features", "ai")
	moved, err := f.svc.HumanMove(f.ctx, task.ID, "custom-review")
	if err != nil { t.Fatal(err) }
	if moved.State != "backlog" || moved.BoardColumn != "custom-review" { t.Fatalf("wrong state: %#v", moved) }
	board, err := f.svc.HumanBoard(f.ctx, f.project.ID)
	if err != nil { t.Fatal(err) }
	if len(board["custom-review"]) != 1 || len(board["backlog"]) != 0 { t.Fatalf("wrong board: %#v", board) }
	if _, err = f.svc.AgentTask(f.ctx, f.ac, task.ID); !errors.Is(err, core.ErrTaskNotAccessible) { t.Fatalf("agent read custom task: %v", err) }
	settings.Columns = settings.Columns[:len(settings.Columns)-1]
	if _, err = f.svc.UpdateProjectSettings(f.ctx, f.project.ID, settings); err != nil { t.Fatal(err) }
	moved, err = f.store.GetTask(f.ctx, task.ID)
	if err != nil { t.Fatal(err) }
	if moved.BoardColumn != "" || moved.State != "backlog" { t.Fatalf("delete did not return task to backlog: %#v", moved) }
}

func TestSettingsCannotReorderOrRenameAIWorkflow(t *testing.T) {
	f := setup(t)
	settings, err := f.svc.GetProjectSettings(f.ctx, f.project.ID)
	if err != nil { t.Fatal(err) }
	settings.Columns[1], settings.Columns[2] = settings.Columns[2], settings.Columns[1]
	if _, err = f.svc.UpdateProjectSettings(f.ctx, f.project.ID, settings); !errors.Is(err, core.ErrInvalidInput) { t.Fatalf("reorder accepted: %v", err) }
	settings, _ = f.svc.GetProjectSettings(f.ctx, f.project.ID)
	settings.Columns[1].Name = "Renamed"
	if _, err = f.svc.UpdateProjectSettings(f.ctx, f.project.ID, settings); !errors.Is(err, core.ErrInvalidInput) { t.Fatalf("rename accepted: %v", err) }
}
