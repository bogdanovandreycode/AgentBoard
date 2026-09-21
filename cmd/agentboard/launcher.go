package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/persistence"
)

func launcherURL(base, path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("project path must be a directory: %s", abs)
	}
	for {
		marker, err := persistence.ReadProjectMarker(abs)
		if err == nil {
			if marker.ProjectID == "" {
				return "", errors.New("project marker has no project_id")
			}
			return base + "/?" + url.Values{"project": {marker.ProjectID}}.Encode(), nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("read project marker: %w", err)
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return base, nil
		}
		abs = parent
	}
}

func agentBoardRunning(base string) bool {
	client := &http.Client{Timeout: time.Second}
	response, err := client.Get(base + "/api/health")
	if err != nil {
		return false
	}
	defer response.Body.Close()
	var health struct{ Service string }
	return response.StatusCode == http.StatusOK && json.NewDecoder(response.Body).Decode(&health) == nil && health.Service == "agentboard"
}
