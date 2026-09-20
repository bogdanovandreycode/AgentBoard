package persistence

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

const markerDirectory = ".agentboard"

type ProjectMarker struct {
	ProjectID string `json:"project_id"`
}

func (s *Store) InitProject(ctx context.Context, path string) (core.Project, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return core.Project{}, err
	}
	name := filepath.Base(filepath.Clean(abs))
	p, err := s.EnsureProject(ctx, name, abs)
	if err != nil {
		return p, err
	}
	dir := filepath.Join(abs, markerDirectory)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return p, err
	}
	b, _ := json.MarshalIndent(ProjectMarker{ProjectID: p.ID}, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, "project.json"), append(b, '\n'), 0o644); err != nil {
		return p, err
	}
	return p, nil
}

func ReadProjectMarker(path string) (ProjectMarker, error) {
	var marker ProjectMarker
	b, err := os.ReadFile(filepath.Join(path, markerDirectory, "project.json"))
	if err != nil {
		return marker, err
	}
	err = json.Unmarshal(b, &marker)
	return marker, err
}
