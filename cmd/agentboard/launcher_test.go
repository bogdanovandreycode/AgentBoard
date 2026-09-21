package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLauncherProjectResolution(t *testing.T) {
	root := t.TempDir()
	base := "http://127.0.0.1:7337"
	if got, err := launcherURL(base, root); err != nil || got != base {
		t.Fatalf("chooser: %s %v", got, err)
	}
	markerDir := filepath.Join(root, ".agentboard")
	if err := os.MkdirAll(markerDir, 0755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(markerDir, "project.json")
	if err := os.WriteFile(marker, []byte(`{"project_id":"project-123"}`), 0644); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(root, "nested", "folder")
	if err := os.MkdirAll(child, 0755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{root, child} {
		got, err := launcherURL(base, path)
		if err != nil || got != base+"/?project=project-123" {
			t.Fatalf("project: %s %v", got, err)
		}
	}
	if err := os.WriteFile(marker, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := launcherURL(base, root); err == nil || !strings.Contains(err.Error(), "project_id") {
		t.Fatalf("empty marker: %v", err)
	}
	if err := os.WriteFile(marker, []byte(`broken`), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := launcherURL(base, root); err == nil {
		t.Fatal("accepted malformed marker")
	}
	if _, err := launcherURL(base, filepath.Join(root, "missing")); err == nil {
		t.Fatal("accepted missing path")
	}
}
