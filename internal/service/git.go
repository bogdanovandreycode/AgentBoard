package service

import (
	"context"
	"os/exec"
	"strings"
)

type GitInfo struct {
	Available    bool
	Branch, Head string
	Dirty        bool
	ChangedFiles []string
}

func (s *Service) ProjectGitInfo(ctx context.Context, projectID string) (GitInfo, error) {
	p, err := s.Store.GetProject(ctx, projectID)
	if err != nil {
		return GitInfo{}, err
	}
	run := func(args ...string) (string, error) {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = p.Path
		b, e := cmd.Output()
		return strings.TrimSpace(string(b)), e
	}
	head, err := run("rev-parse", "HEAD")
	if err != nil {
		return GitInfo{Available: false}, nil
	}
	branch, _ := run("branch", "--show-current")
	status, _ := run("status", "--porcelain")
	info := GitInfo{Available: true, Branch: branch, Head: head, Dirty: status != ""}
	for _, line := range strings.Split(status, "\n") {
		if len(line) > 3 {
			info.ChangedFiles = append(info.ChangedFiles, strings.TrimSpace(line[3:]))
		}
	}
	return info, nil
}
