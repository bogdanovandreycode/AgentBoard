package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

// ResolveAgentTaskID resolves a full ID or a unique prefix within the worker's
// visible tasks. Hidden tasks never participate in matching or ambiguity errors.
func (s *Service) ResolveAgentTaskID(ctx context.Context, a core.AgentContext, ref string) (string, error) {
	ref = strings.ToLower(strings.TrimSpace(ref))
	if len(ref) < 8 || len(ref) > 36 || strings.Trim(ref, "0123456789abcdef-") != "" {
		return "", core.ErrTaskNotAccessible
	}
	rows, err := s.Store.DB.QueryContext(ctx, `SELECT id FROM tasks WHERE project_id=? AND assignee_type='worker' AND assignee_worker_id=? AND state IN ('features','in_progress','testing','verification') AND id LIKE ? LIMIT 2`, a.Project.ID, a.Worker.ID, ref+"%")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return "", err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	switch len(ids) {
	case 0:
		return "", core.ErrTaskNotAccessible
	case 1:
		return ids[0], nil
	default:
		return "", core.ErrAmbiguousTaskID
	}
}

type AgentWorkflow struct {
	NextTool            string `json:"next_tool,omitempty"`
	HumanTestingPending bool   `json:"human_testing_pending"`
	Instructions        string `json:"instructions"`
}

// AgentWorkflow describes the same policy enforced by AgentMove.
func (s *Service) AgentWorkflow(ctx context.Context, a core.AgentContext, t core.Task) AgentWorkflow {
	if t.AssigneeType != "worker" || t.AssigneeWorkerID == nil || *t.AssigneeWorkerID != a.Worker.ID {
		return AgentWorkflow{Instructions: "This feature is assigned to another actor; continue your own assigned tasks."}
	}
	switch t.State {
	case "features":
		return AgentWorkflow{NextTool: "move_to_in_progress", Instructions: "Start implementation."}
	case "in_progress":
		return AgentWorkflow{NextTool: "move_to_testing", Instructions: "When implementation is ready, move to testing. For human or hybrid testing, provide handoff instructions, then finish the session or work on another task; never wait or poll for a human."}
	case "testing":
		err := s.canAgentVerify(ctx, t, a.Worker.ID)
		if errors.Is(err, core.ErrHumanTestPending) {
			return AgentWorkflow{HumanTestingPending: true, Instructions: "Leave this task in testing. Provide handoff instructions and finish the session or work on another task. Do not wait, poll, retry verification, or perform testing on behalf of a human."}
		}
		if err != nil {
			return AgentWorkflow{NextTool: "record_test_run", Instructions: "Run the configured AI checks and record their results before requesting verification."}
		}
		return AgentWorkflow{NextTool: "move_to_verification", Instructions: "Required testing has passed; move to verification."}
	default:
		return AgentWorkflow{Instructions: "Implementation is handed off. Finish the session or work on another task; final acceptance belongs to the human."}
	}
}
