package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

// ImportTask is the portable JSON format accepted by the Human import endpoint.
// Keys and dependencies are local to one import file; no external IDs are trusted.
type ImportTask struct {
	Key                   string            `json:"key"`
	Title                 string            `json:"title"`
	Description           string            `json:"description"`
	State                 string            `json:"state"`
	Priority              string            `json:"priority"`
	TestingMode           string            `json:"testing_mode"`
	Assignee              ImportAssignee    `json:"assignee"`
	AITestInstructions    string            `json:"ai_test_instructions"`
	HumanTestInstructions string            `json:"human_test_instructions"`
	Properties            map[string]string `json:"properties"`
	DependsOn             []string          `json:"depends_on"`
}

type ImportAssignee struct {
	Type   string `json:"type"`
	Worker string `json:"worker"`
}

type ImportDocument struct {
	Version int          `json:"version"`
	Tasks   []ImportTask `json:"tasks"`
}

func (s *Service) HumanImportTasks(ctx context.Context, projectID string, doc ImportDocument) (int, error) {
	if _, err := s.Store.GetProject(ctx, projectID); err != nil {
		return 0, err
	}
	if doc.Version != 1 || len(doc.Tasks) == 0 || len(doc.Tasks) > 1000 {
		return 0, fmt.Errorf("%w: expected version 1 and 1-1000 tasks", core.ErrInvalidInput)
	}
	workers, err := s.Store.ListWorkers(ctx, projectID, false)
	if err != nil {
		return 0, err
	}
	definitions, err := s.ListPropertyDefinitions(ctx, projectID)
	if err != nil {
		return 0, err
	}
	workerIDs := map[string]string{}
	for _, worker := range workers {
		workerIDs[worker.ID] = worker.ID
		workerIDs[worker.Slug] = worker.ID
	}
	propertyIDs := map[string]string{}
	for _, definition := range definitions {
		propertyIDs[definition.ID] = definition.ID
		propertyIDs[definition.Name] = definition.ID
	}
	keys := map[string]string{}
	for _, task := range doc.Tasks {
		if task.Key == "" {
			continue
		}
		if _, exists := keys[task.Key]; exists {
			return 0, fmt.Errorf("%w: duplicate key %q", core.ErrInvalidInput, task.Key)
		}
		keys[task.Key] = id()
	}
	ids := make([]string, len(doc.Tasks))
	for i, task := range doc.Tasks {
		if task.Key != "" {
			ids[i] = keys[task.Key]
		} else {
			ids[i] = id()
		}
	}
	positions := map[string]float64{}
	for _, state := range []string{"backlog", "features", "in_progress", "testing", "verification", "complete"} {
		var position float64
		if err := s.Store.DB.QueryRowContext(ctx, `SELECT COALESCE(MAX(position),0) FROM tasks WHERE project_id=? AND state=?`, projectID, state).Scan(&position); err != nil {
			return 0, err
		}
		positions[state] = position
	}
	// Validate the entire file before opening the write transaction.
	for i := range doc.Tasks {
		task := &doc.Tasks[i]
		if strings.TrimSpace(task.Title) == "" {
			return 0, fmt.Errorf("%w: task %d has no title", core.ErrInvalidInput, i+1)
		}
		if task.State == "" {
			task.State = "backlog"
		}
		if _, ok := positions[task.State]; !ok {
			return 0, fmt.Errorf("%w: task %d has invalid state", core.ErrInvalidInput, i+1)
		}
		if task.Priority == "" {
			task.Priority = "medium"
		}
		if task.Priority != "critical" && task.Priority != "high" && task.Priority != "medium" && task.Priority != "low" {
			return 0, fmt.Errorf("%w: task %d has invalid priority", core.ErrInvalidInput, i+1)
		}
		if task.TestingMode == "" {
			task.TestingMode = "ai"
		}
		if task.TestingMode != "ai" && task.TestingMode != "human" && task.TestingMode != "hybrid" {
			return 0, fmt.Errorf("%w: task %d has invalid testing_mode", core.ErrInvalidInput, i+1)
		}
		if task.Assignee.Type == "" {
			task.Assignee.Type = "unassigned"
		}
		if task.Assignee.Type == "worker" {
			if _, ok := workerIDs[task.Assignee.Worker]; !ok {
				return 0, fmt.Errorf("%w: task %d has unknown worker %q", core.ErrInvalidInput, i+1, task.Assignee.Worker)
			}
		} else if (task.Assignee.Type != "human" && task.Assignee.Type != "unassigned") || task.Assignee.Worker != "" {
			return 0, fmt.Errorf("%w: task %d has invalid assignee", core.ErrInvalidInput, i+1)
		}
		for name := range task.Properties {
			if _, ok := propertyIDs[name]; !ok {
				return 0, fmt.Errorf("%w: task %d has unknown property %q", core.ErrInvalidInput, i+1, name)
			}
		}
		for _, dependency := range task.DependsOn {
			if _, ok := keys[dependency]; !ok || dependency == task.Key {
				return 0, fmt.Errorf("%w: task %d has invalid dependency %q", core.ErrInvalidInput, i+1, dependency)
			}
		}
	}
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	now := core.Now()
	for i, task := range doc.Tasks {
		positions[task.State] += 1000
		var workerID any
		if task.Assignee.Type == "worker" {
			workerID = workerIDs[task.Assignee.Worker]
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO tasks(id,project_id,title,description,state,position,priority,assignee_type,assignee_worker_id,created_by_type,testing_mode,ai_test_instructions,human_test_instructions,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, ids[i], projectID, task.Title, task.Description, task.State, positions[task.State], task.Priority, task.Assignee.Type, workerID, "human", task.TestingMode, task.AITestInstructions, task.HumanTestInstructions, now, now)
		if err != nil {
			return 0, err
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, id(), ids[i], "system", "task_created", "Task imported by Human", now)
		if err != nil {
			return 0, err
		}
		for name, value := range task.Properties {
			_, err = tx.ExecContext(ctx, `INSERT INTO task_property_values(task_id,property_definition_id,value,updated_at) VALUES(?,?,?,?)`, ids[i], propertyIDs[name], value, now)
			if err != nil {
				return 0, err
			}
		}
	}
	for i, task := range doc.Tasks {
		for _, dependency := range task.DependsOn {
			_, err = tx.ExecContext(ctx, `INSERT INTO task_dependencies(task_id,depends_on_task_id,created_at) VALUES(?,?,?)`, ids[i], keys[dependency], now)
			if err != nil {
				return 0, err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(doc.Tasks), nil
}
