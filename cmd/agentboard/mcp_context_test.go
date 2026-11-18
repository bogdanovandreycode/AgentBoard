package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
)

func TestMCPWorkspaceDetectsLaunchProjectWithoutRequiringOne(t *testing.T) {
	ctx := context.Background()
	st, err := persistence.Open(filepath.Join(t.TempDir(), "board.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	p, err := st.EnsureProject(ctx, "one", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	subdir := filepath.Join(p.Path, "source", "nested")
	if err := os.MkdirAll(subdir, 0700); err != nil {
		t.Fatal(err)
	}
	workspace, err := mcpWorkspace(ctx, svc, subdir, false)
	if err != nil || workspace.LaunchProject == nil || workspace.LaunchProject.ID != p.ID {
		t.Fatalf("nested discovery: %+v %v", workspace, err)
	}
	workspace, err = mcpWorkspace(ctx, svc, t.TempDir(), true)
	if err != nil || workspace.LaunchProject != nil || !workspace.AllowMismatch {
		t.Fatalf("neutral launch directory must not fail: %+v %v", workspace, err)
	}
	workspace, err = mcpWorkspace(ctx, svc, filepath.Join(t.TempDir(), "missing"), false)
	if err != nil || workspace.LaunchProject != nil {
		t.Fatalf("missing launch directory must not fail: %+v %v", workspace, err)
	}
}
