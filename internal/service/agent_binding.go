package service

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

// AgentBindingRequest describes how an MCP worker process selects its project
// and worker. Empty ProjectPath means the project is detected from the client's
// workspace roots first and the MCP launch directory second. Empty WorkerRef
// selects the only enabled worker of that project.
type AgentBindingRequest struct {
	ProjectPath     string
	WorkerRef       string
	Roots           []string
	LaunchDirectory string
}

type AgentBinding struct {
	Project core.Project
	Worker  core.Worker
	// Source is "--project", "client_roots" or "launch_directory".
	Source string
}

// BindingError is a configuration problem the agent cannot fix by itself. Its
// message is written to be relayed to the human user verbatim.
type BindingError struct {
	Code    string
	Message string
}

func (e *BindingError) Error() string { return e.Code + ": " + e.Message }

const relayHint = " Tell the user about this; do not guess another project or worker, and do not track the work elsewhere instead."

func bindingError(code, format string, args ...any) error {
	return &BindingError{Code: code, Message: fmt.Sprintf(format, args...) + relayHint}
}

// ResolveAgentBinding selects the project and worker for an MCP worker session.
// It never depends on whether the web UI is running or where it was started.
func (s *Service) ResolveAgentBinding(ctx context.Context, req AgentBindingRequest) (AgentBinding, error) {
	projects, err := s.Store.ListProjects(ctx)
	if err != nil {
		return AgentBinding{}, err
	}
	if len(projects) == 0 {
		return AgentBinding{}, bindingError("NO_PROJECTS", "AgentBoard has no registered projects yet (database %s). The user must register the project folder with `agentboard init <folder>` or in the AgentBoard web UI, and create a worker for it; then call get_my_board again (no restart is needed).", s.Store.DatabasePath())
	}
	var binding AgentBinding
	switch {
	case req.ProjectPath != "":
		p, err := s.ProjectForDirectory(ctx, req.ProjectPath)
		if err != nil {
			if !errors.Is(err, core.ErrNotFound) && !isMissingDirectory(err) {
				return AgentBinding{}, err
			}
			return AgentBinding{}, bindingError("PROJECT_NOT_FOUND", "--project \"%s\" is not inside a registered AgentBoard project. Registered projects: %s. The user must fix --project in the MCP configuration or run `agentboard init` in that folder.", req.ProjectPath, projectList(projects))
		}
		binding = AgentBinding{Project: p, Source: "--project"}
	default:
		found := map[string]core.Project{}
		for _, root := range req.Roots {
			p, err := s.ProjectForDirectory(ctx, root)
			if err == nil {
				found[p.ID] = p
			} else if !errors.Is(err, core.ErrNotFound) && !isMissingDirectory(err) {
				return AgentBinding{}, err
			}
		}
		if len(found) > 1 {
			var matched []core.Project
			for _, p := range found {
				matched = append(matched, p)
			}
			return AgentBinding{}, bindingError("PROJECT_AMBIGUOUS", "the client workspace roots \"%s\" belong to several AgentBoard projects: %s. The user must add --project PATH to this MCP server configuration.", strings.Join(req.Roots, "\", \""), projectList(matched))
		}
		for _, p := range found {
			binding = AgentBinding{Project: p, Source: "client_roots"}
		}
		if binding.Project.ID == "" && req.LaunchDirectory != "" {
			p, err := s.ProjectForDirectory(ctx, req.LaunchDirectory)
			if err == nil {
				binding = AgentBinding{Project: p, Source: "launch_directory"}
			} else if !errors.Is(err, core.ErrNotFound) && !isMissingDirectory(err) {
				return AgentBinding{}, err
			}
		}
		if binding.Project.ID == "" {
			where := fmt.Sprintf("the MCP launch directory \"%s\"", req.LaunchDirectory)
			if len(req.Roots) > 0 {
				where = fmt.Sprintf("the client workspace roots \"%s\" or %s", strings.Join(req.Roots, "\", \""), where)
			}
			return AgentBinding{}, bindingError("PROJECT_REQUIRED", "no registered AgentBoard project contains %s. Registered projects: %s. The user must start the agent inside one of these folders, register the current folder with `agentboard init`, or add --project PATH to the MCP configuration.", where, projectList(projects))
		}
	}
	w, err := s.resolveBindingWorker(ctx, binding.Project, req.WorkerRef)
	if err != nil {
		return AgentBinding{}, err
	}
	binding.Worker = w
	return binding, nil
}

func (s *Service) resolveBindingWorker(ctx context.Context, p core.Project, ref string) (core.Worker, error) {
	workers, err := s.Store.ListWorkers(ctx, p.ID, true)
	if err != nil {
		return core.Worker{}, err
	}
	if ref == "" {
		switch len(workers) {
		case 1:
			return workers[0], nil
		case 0:
			return core.Worker{}, bindingError("NO_WORKERS", "project %q (%s) has no enabled workers. The user must create a worker in AgentBoard → Workers and put its slug into --worker in the MCP configuration.", p.Name, p.Path)
		default:
			return core.Worker{}, bindingError("WORKER_REQUIRED", "project %q (%s) has several enabled workers: %s. The user must choose one with --worker SLUG in the MCP configuration.", p.Name, p.Path, workerList(workers))
		}
	}
	w, err := s.Store.FindWorker(ctx, p.ID, ref)
	if errors.Is(err, core.ErrNotFound) {
		return core.Worker{}, bindingError("WORKER_NOT_FOUND", "worker %q does not exist in project %q (%s). Enabled workers: %s. The user must create this worker in AgentBoard → Workers or change --worker.", ref, p.Name, p.Path, workerList(workers))
	}
	if err != nil {
		return core.Worker{}, err
	}
	if !w.Enabled || w.Archived {
		return core.Worker{}, bindingError("WORKER_DISABLED", "worker %q in project %q is disabled or archived. The user must enable it in AgentBoard → Workers or configure another --worker.", w.Slug, p.Name)
	}
	return w, nil
}

func isMissingDirectory(err error) bool {
	return errors.Is(err, core.ErrInvalidInput) || errors.Is(err, fs.ErrNotExist)
}

func projectList(projects []core.Project) string {
	sort.Slice(projects, func(i, j int) bool { return projects[i].Name < projects[j].Name })
	const limit = 20
	var parts []string
	for i, p := range projects {
		if i == limit {
			parts = append(parts, fmt.Sprintf("and %d more", len(projects)-limit))
			break
		}
		parts = append(parts, fmt.Sprintf("%s (%s)", p.Name, p.Path))
	}
	return strings.Join(parts, ", ")
}

func workerList(workers []core.Worker) string {
	if len(workers) == 0 {
		return "none"
	}
	var parts []string
	for _, w := range workers {
		parts = append(parts, w.Slug)
	}
	return strings.Join(parts, ", ")
}
