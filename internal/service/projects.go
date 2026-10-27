package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
)

type ProjectUpdate struct{ Name, Path string }

// UpdateProject changes the registered binding, never moves the user's files.
func (s *Service) UpdateProject(ctx context.Context, projectID string, in ProjectUpdate) (core.Project, error) {
	p, err := s.Store.GetProject(ctx, projectID)
	if err != nil {
		return p, err
	}
	name := strings.TrimSpace(in.Name)
	if name == "" || len([]rune(name)) > 120 || strings.TrimSpace(in.Path) == "" {
		return p, core.ErrInvalidInput
	}
	path, err := filepath.Abs(in.Path)
	if err != nil {
		return p, fmt.Errorf("%w: project folder", core.ErrInvalidInput)
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return p, fmt.Errorf("%w: project folder must be an existing directory", core.ErrInvalidInput)
	}
	path = filepath.Clean(path)
	sameDirectory := p.Path == path
	if oldInfo, err := os.Stat(p.Path); err == nil {
		sameDirectory = os.SameFile(info, oldInfo)
	}
	if existing, err := s.Store.GetProjectByPath(ctx, path); err == nil && existing.ID != projectID {
		return p, fmt.Errorf("%w: this folder belongs to another project", core.ErrInvalidInput)
	} else if err != nil && !errors.Is(err, core.ErrNotFound) {
		return p, err
	}
	marker, markerErr := persistence.ReadProjectMarker(path)
	if markerErr == nil && marker.ProjectID != projectID {
		return p, fmt.Errorf("%w: selected folder contains another project's marker", core.ErrInvalidInput)
	}
	if markerErr != nil && !errors.Is(markerErr, os.ErrNotExist) {
		return p, fmt.Errorf("%w: invalid project marker: %v", core.ErrInvalidInput, markerErr)
	}
	createdMarker := false
	markerPath := filepath.Join(path, ".agentboard", "project.json")
	if errors.Is(markerErr, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(markerPath), 0755); err != nil {
			return p, err
		}
		data, _ := json.MarshalIndent(persistence.ProjectMarker{ProjectID: projectID}, "", "  ")
		f, err := os.OpenFile(markerPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return p, err
		}
		_, writeErr := f.Write(append(data, '\n'))
		closeErr := f.Close()
		if writeErr != nil || closeErr != nil {
			_ = os.Remove(markerPath)
			if writeErr != nil {
				return p, writeErr
			}
			return p, closeErr
		}
		createdMarker = true
	}
	_, err = s.Store.DB.ExecContext(ctx, `UPDATE projects SET name=?,path=?,updated_at=? WHERE id=?`, name, path, core.Now(), projectID)
	if err != nil {
		if createdMarker {
			_ = os.Remove(markerPath)
		}
		return p, err
	}
	if !sameDirectory {
		if old, err := persistence.ReadProjectMarker(p.Path); err == nil && old.ProjectID == projectID {
			_ = os.Remove(filepath.Join(p.Path, ".agentboard", "project.json"))
		}
	}
	return s.Store.GetProject(ctx, projectID)
}
