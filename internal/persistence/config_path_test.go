package persistence

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestWindowsDefaultDBPathWithoutAPPDATA(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows profile fallback")
	}
	home := t.TempDir()
	t.Setenv("APPDATA", "")
	t.Setenv("USERPROFILE", home)
	path, err := DefaultDBPath()
	if err != nil || path != filepath.Join(home, "AppData", "Roaming", "AgentBoard", "agentboard.db") {
		t.Fatalf("wrong fallback: %s %v", path, err)
	}
	t.Setenv("USERPROFILE", "")
	if _, err := DefaultDBPath(); err == nil {
		t.Fatal("missing profile silently selected another database")
	}
}
