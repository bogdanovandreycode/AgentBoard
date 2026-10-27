package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func (a *API) pickFolder(w http.ResponseWriter, r *http.Request) {
	var in struct{ Path string }
	if !decode(w, r, &in) {
		return
	}
	path, err := selectFolder(r.Context(), in.Path)
	if err != nil {
		write(w, 500, map[string]string{"error": "Folder picker failed; you can enter the path manually. " + err.Error()})
		return
	}
	write(w, 200, map[string]string{"path": path})
}

func selectFolder(ctx context.Context, initial string) (string, error) {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		return selectWindowsFolder(ctx, initial)
	case "darwin":
		command = exec.CommandContext(ctx, "osascript", "-e", `POSIX path of (choose folder with prompt "Select project folder")`)
	default:
		command = exec.CommandContext(ctx, "zenity", "--file-selection", "--directory", "--title=Select project folder", "--filename="+initial+string(filepath.Separator))
	}
	output, err := command.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok && ((runtime.GOOS != "windows" && exit.ExitCode() == 1) || strings.Contains(string(exit.Stderr), "-128")) {
			return "", nil
		}
		return "", fmt.Errorf("%w", err)
	}
	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", nil
	}
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		return "", fmt.Errorf("selected folder is unavailable")
	}
	return filepath.Clean(path), nil
}
