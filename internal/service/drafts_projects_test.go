package service

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/google/uuid"
)

func TestDraftPersistenceReadinessAndAtomicPublication(t *testing.T) {
	f := setup(t)
	settings, err := f.svc.GetGlobalSettings(f.ctx)
	if err != nil || settings.DraftAutosaveSeconds != 10 || !settings.DraftSaveOnBlur || !settings.DraftAutosaveEnabled || !settings.CheckProjectContext {
		t.Fatalf("defaults: %+v %v", settings, err)
	}
	id := uuid.NewString()
	empty, err := f.svc.SaveTaskDraft(f.ctx, f.project.ID, id, DraftSave{Input: core.TaskInput{Priority: "medium", TestingMode: "ai", Assignee: core.Assignee{Type: "unassigned"}}})
	if err != nil || empty.Saved || empty.Revision != 0 {
		t.Fatalf("empty persisted: %+v %v", empty, err)
	}
	partial := core.TaskInput{Title: "Work in progress"}
	draft, err := f.svc.SaveTaskDraft(f.ctx, f.project.ID, id, DraftSave{Input: partial})
	if err != nil || !draft.Saved || draft.Ready || draft.Revision != 1 || len(draft.MissingFields) != 3 {
		t.Fatalf("partial: %+v %v", draft, err)
	}
	if _, err := f.svc.AgentTask(f.ctx, f.ac, id); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("draft exposed: %v", err)
	}
	if _, err := f.svc.PublishTaskDraft(f.ctx, f.project.ID, id, 1); !errors.Is(err, core.ErrDraftIncomplete) {
		t.Fatalf("incomplete published: %v", err)
	}
	if _, err := f.svc.SaveTaskDraft(f.ctx, f.project.ID, id, DraftSave{Input: partial}); !errors.Is(err, core.ErrDraftConflict) {
		t.Fatalf("stale overwrite: %v", err)
	}
	complete := core.TaskInput{Title: "Ready", Description: "Task details", State: "features", Assignee: workerAssignee(f.a.ID), TestingMode: "ai"}
	draft, err = f.svc.SaveTaskDraft(f.ctx, f.project.ID, id, DraftSave{Input: complete, Revision: 1})
	if err != nil || !draft.Ready || draft.Revision != 2 {
		t.Fatalf("ready: %+v %v", draft, err)
	}
	task, err := f.svc.PublishTaskDraft(f.ctx, f.project.ID, id, 2)
	if err != nil || task.CreatedByType != "human" || task.State != "features" {
		t.Fatalf("publish: %+v %v", task, err)
	}
	retry, err := f.svc.PublishTaskDraft(f.ctx, f.project.ID, id, 2)
	if err != nil || retry.ID != task.ID {
		t.Fatalf("duplicate: %+v %v", retry, err)
	}
	if _, err := f.svc.SaveTaskDraft(f.ctx, f.project.ID, id, DraftSave{Input: complete, Revision: 2}); !errors.Is(err, core.ErrDraftConflict) {
		t.Fatalf("resurrected: %v", err)
	}
	drafts, err := f.svc.ListTaskDrafts(f.ctx, f.project.ID)
	if err != nil || len(drafts) != 0 {
		t.Fatalf("published listed: %+v %v", drafts, err)
	}
	tasks, err := f.store.ListTasks(f.ctx, f.project.ID)
	if err != nil || len(tasks) != 1 {
		t.Fatalf("task count: %+v %v", tasks, err)
	}
	badID := uuid.NewString()
	complete.Properties = map[string]string{"missing-property": "x"}
	bad, err := f.svc.SaveTaskDraft(f.ctx, f.project.ID, badID, DraftSave{Input: complete})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.PublishTaskDraft(f.ctx, f.project.ID, badID, bad.Revision); err == nil {
		t.Fatal("invalid draft published")
	}
	drafts, err = f.svc.ListTaskDrafts(f.ctx, f.project.ID)
	if err != nil || len(drafts) != 1 {
		t.Fatalf("failed draft lost: %+v %v", drafts, err)
	}
	if _, err := f.svc.SaveTaskDraft(f.ctx, f.project.ID, badID, DraftSave{Revision: bad.Revision}); err != nil {
		t.Fatal(err)
	}
	drafts, err = f.svc.ListTaskDrafts(f.ctx, f.project.ID)
	if err != nil || len(drafts) != 0 {
		t.Fatalf("cleared draft retained: %+v %v", drafts, err)
	}
}

func TestGlobalSettingsSharedAndPersisted(t *testing.T) {
	f := setup(t)
	value, _ := f.svc.GetGlobalSettings(f.ctx)
	value.DraftAutosaveSeconds, value.DraftSaveOnBlur, value.CheckProjectContext = 23, false, false
	if _, err := f.svc.UpdateGlobalSettings(f.ctx, value); err != nil {
		t.Fatal(err)
	}
	otherStore, err := persistence.Open(f.store.DatabasePath())
	if err != nil {
		t.Fatal(err)
	}
	defer otherStore.Close()
	got, err := New(otherStore).GetGlobalSettings(f.ctx)
	if err != nil || got != value {
		t.Fatalf("not shared: %+v %v", got, err)
	}
	value.DraftAutosaveSeconds = 0
	if _, err := f.svc.UpdateGlobalSettings(f.ctx, value); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("bad interval accepted: %v", err)
	}
}

func TestProjectStatisticsAndFolderRebinding(t *testing.T) {
	f := setup(t)
	if _, err := f.store.InitProject(f.ctx, f.project.Path); err != nil {
		t.Fatal(err)
	}
	task := createTask(t, f, f.a, "features", "ai")
	if _, err := f.store.DB.ExecContext(f.ctx, `UPDATE tasks SET updated_at='2020-01-01T00:00:00Z' WHERE id=?`, task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := f.store.DB.ExecContext(f.ctx, `UPDATE history_entries SET created_at='2021-01-01T00:00:00Z' WHERE task_id=?`, task.ID); err != nil {
		t.Fatal(err)
	}
	projects, err := f.store.ListProjects(f.ctx)
	if err != nil || len(projects) != 1 || projects[0].TaskCount != 1 || projects[0].LastActivityAt != "2021-01-01T00:00:00Z" {
		t.Fatalf("history not counted: %+v %v", projects, err)
	}
	usage, err := f.svc.AgentReportUsage(f.ctx, f.ac, task.ID, core.UsageInput{})
	if err != nil {
		t.Fatal(err)
	}
	projects, _ = f.store.ListProjects(f.ctx)
	if projects[0].LastActivityAt != usage.CreatedAt {
		t.Fatalf("usage not counted: %+v", projects)
	}
	destination := filepath.Join(t.TempDir(), "New project folder")
	if err := os.MkdirAll(destination, 0700); err != nil {
		t.Fatal(err)
	}
	updated, err := f.svc.UpdateProject(f.ctx, f.project.ID, ProjectUpdate{Name: "Renamed", Path: destination})
	if err != nil || updated.ID != f.project.ID || updated.Name != "Renamed" || updated.Path != destination {
		t.Fatalf("rebind: %+v %v", updated, err)
	}
	marker, err := persistence.ReadProjectMarker(destination)
	if err != nil || marker.ProjectID != f.project.ID {
		t.Fatalf("new marker: %+v %v", marker, err)
	}
	if _, err := persistence.ReadProjectMarker(f.project.Path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("old marker retained: %v", err)
	}
	if _, err := f.svc.BeginAgentSession(f.ctx, destination, f.a.Slug, "test"); err != nil {
		t.Fatalf("worker not rebound: %v", err)
	}
	other, err := f.store.InitProject(f.ctx, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.UpdateProject(f.ctx, f.project.ID, ProjectUpdate{Name: "Wrong", Path: other.Path}); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("collision: %v", err)
	}
	if _, err := f.svc.UpdateProject(f.ctx, f.project.ID, ProjectUpdate{Name: "Wrong", Path: filepath.Join(destination, "missing")}); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("missing folder: %v", err)
	}
	if _, err := f.svc.UpdateProject(f.ctx, f.project.ID, ProjectUpdate{Name: "Again", Path: destination}); err != nil {
		t.Fatal(err)
	}
	if _, err := persistence.ReadProjectMarker(destination); err != nil {
		t.Fatal(err)
	}
	if got, err := f.store.GetTask(f.ctx, task.ID); err != nil || got.ProjectID != f.project.ID {
		t.Fatalf("task lost: %+v %v", got, err)
	}
}
