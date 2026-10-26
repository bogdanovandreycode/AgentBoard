package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/mcpserver"
	"github.com/bogdanovandreycode/agentboard/internal/service"
)

func resolveMCPProject(ctx context.Context, svc *service.Service, cwd, target string, allowMismatch bool) (core.Project, *mcpserver.WorkspaceContext, error) {
	workspace := &mcpserver.WorkspaceContext{LaunchDirectory: cwd, AllowMismatch: allowMismatch}
	launch, err := svc.ProjectForDirectory(ctx, cwd)
	if err == nil {
		workspace.LaunchProject = &launch
	} else if !errors.Is(err, core.ErrNotFound) {
		return core.Project{}, nil, fmt.Errorf("detect MCP launch project: %w", err)
	}
	if target == "" {
		if workspace.LaunchProject == nil {
			return core.Project{}, nil, fmt.Errorf("PROJECT_REQUIRED: %q is outside registered projects; use --project PATH --worker SLUG (and --db PATH for a custom database)", cwd)
		}
		return launch, workspace, nil
	}
	p, err := svc.ProjectForDirectory(ctx, target)
	if err != nil {
		return p, nil, fmt.Errorf("PROJECT_LOOKUP_FAILED: --project %q: %w; use a registered project directory and check --db", target, err)
	}
	return p, workspace, nil
}
