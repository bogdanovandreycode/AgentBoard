package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/persistence"
	"github.com/bogdanovandreycode/agentboard/internal/service"
)

func TestMCPProjectDiscoveryAndExplicitCrossProjectContext(t *testing.T) {
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
	other, err := st.EnsureProject(ctx, "two", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	subdir := filepath.Join(p.Path, "source", "nested")
	if err := os.MkdirAll(subdir, 0700); err != nil {
		t.Fatal(err)
	}
	got, workspace, err := resolveMCPProject(ctx, svc, subdir, "", false)
	if err != nil || got.ID != p.ID || workspace.LaunchProject.ID != p.ID {
		t.Fatalf("nested discovery: %+v %+v %v", got, workspace, err)
	}
	neutral := t.TempDir()
	got, workspace, err = resolveMCPProject(ctx, svc, neutral, p.Path, false)
	if err != nil || got.ID != p.ID || workspace.LaunchProject != nil {
		t.Fatalf("neutral launch: %+v %+v %v", got, workspace, err)
	}
	got, workspace, err = resolveMCPProject(ctx, svc, other.Path, p.Path, false)
	if err != nil || got.ID != p.ID || workspace.LaunchProject.ID != other.ID || workspace.AllowMismatch {
		t.Fatalf("cross-project context lost: %+v %+v %v", got, workspace, err)
	}
	if _, _, err := resolveMCPProject(ctx, svc, neutral, "", false); err == nil || !strings.Contains(err.Error(), "PROJECT_REQUIRED") {
		t.Fatalf("unknown directory: %v", err)
	}
	if runtime.GOOS == "windows" {
		if got, _, err := resolveMCPProject(ctx, svc, neutral, strings.ToLower(p.Path), false); err != nil || got.ID != p.ID {
			t.Fatalf("Windows path casing: %+v %v", got, err)
		}
	}
}
