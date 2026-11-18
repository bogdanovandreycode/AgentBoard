package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/mcpserver"
	"github.com/bogdanovandreycode/agentboard/internal/service"
)

// mcpWorkspace records the MCP launch directory and the project containing it.
// A directory outside registered projects is not an error.
func mcpWorkspace(ctx context.Context, svc *service.Service, cwd string, allowMismatch bool) (*mcpserver.WorkspaceContext, error) {
	workspace := &mcpserver.WorkspaceContext{LaunchDirectory: cwd, AllowMismatch: allowMismatch}
	launch, err := svc.ProjectForDirectory(ctx, cwd)
	switch {
	case err == nil:
		workspace.LaunchProject = &launch
	case errors.Is(err, core.ErrNotFound), errors.Is(err, fs.ErrNotExist):
	default:
		return nil, fmt.Errorf("detect MCP launch project: %w", err)
	}
	return workspace, nil
}
