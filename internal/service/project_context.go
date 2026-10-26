package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

// ProjectForDirectory finds the nearest registered project containing a folder.
// It never depends on the HTTP server's working directory.
func (s *Service) ProjectForDirectory(ctx context.Context, path string) (core.Project, error) {
	dir, err := filepath.Abs(path)
	if err != nil {
		return core.Project{}, err
	}
	info, err := os.Stat(dir)
	if err != nil {
		return core.Project{}, err
	}
	if !info.IsDir() {
		return core.Project{}, core.ErrInvalidInput
	}
	for {
		p, err := s.Store.GetProjectByPath(ctx, dir)
		if err == nil || !errors.Is(err, core.ErrNotFound) {
			return p, err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return core.Project{}, core.ErrNotFound
		}
		dir = parent
	}
}
