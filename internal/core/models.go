package core

import (
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("NOT_FOUND")
	ErrTaskNotAccessible = errors.New("TASK_NOT_ACCESSIBLE")
	ErrInvalidTransition = errors.New("INVALID_TRANSITION")
	ErrHumanTestPending  = errors.New("HUMAN_TEST_PENDING")
	ErrWorkerDisabled    = errors.New("WORKER_DISABLED")
	ErrInvalidInput      = errors.New("INVALID_INPUT")
)

type Project struct {
	ID, Name, Path, CreatedAt, UpdatedAt string
}

type Worker struct {
	ID, ProjectID, Name, Slug, Description, Kind, Capabilities string
	Enabled, Archived                                          bool
	CreatedAt, UpdatedAt                                       string
	AssignedTaskCount                                          int
	ActiveSessionCount, SessionCount, MCPCalls                 int64
	LastActivityAt                                             *string
}

type WorkerSession struct {
	ID, ProjectID, WorkerID, StartedAt, ClientInfo, CreatedAt string
	EndedAt                                                   *string
	LastActivityAt, LastSeenAt                                *string
	MCPCalls                                                  int64
}

type Task struct {
	ID, ProjectID, Title, Description, State, Priority, AssigneeType string
	Position                                                         float64
	AssigneeWorkerID, CreatedByWorkerID, CreatedBySessionID          *string
	CreatedByType                                                    string
	SourceTaskID                                                     *string
	TestingMode, AITestInstructions, HumanTestInstructions           string
	CreatedAt, UpdatedAt                                             string
	AssigneeName, CreatorName                                        string
	Dependencies                                                     []Dependency
	Artifacts                                                        []Artifact
	Properties                                                       []TaskPropertyValue
}

type Dependency struct {
	TaskID, DependsOnTaskID, Title, State string
	Artifacts                             []Artifact
}

type Artifact struct {
	ID, TaskID, Name, Kind, Description, CreatedByType, CreatedAt string
	Path, URL, CreatedByWorkerID, CreatedBySessionID              *string
	CreatorName                                                   string
}

type HistoryEntry struct {
	ID, TaskID, ActorType, EntryType, Content, CreatedAt string
	WorkerID, WorkerSessionID                            *string
	ActorName                                            string
}

type TestRun struct {
	ID, TaskID, Runner, Type, Status, Summary, OutputExcerpt, CreatedAt string
	WorkerID, WorkerSessionID, Command                                  *string
	DurationMS                                                          *int64
	ActorName                                                           string
}

type UsageEvent struct {
	ID, ProjectID, TaskID, WorkerID, WorkerSessionID, CreatedAt string
	Provider, Model, ReasoningEffort, Source                    *string
	InputTokens, OutputTokens, ReasoningTokens, CachedTokens    *int64
	ModelCalls, ToolCalls, MCPCalls, WallTimeMS                 *int64
	WorkerName                                                  string
}

type PropertyDefinition struct {
	ID, ProjectID, Name, Type, Options, Visibility, CreatedAt, UpdatedAt string
}

type TaskPropertyValue struct {
	TaskID, PropertyDefinitionID, Name, Type, Visibility, Value, UpdatedAt string
}

type Assignee struct {
	Type     string  `json:"type"`
	WorkerID *string `json:"worker_id,omitempty"`
}

type TaskInput struct {
	Title, Description, State, Priority, TestingMode string
	Assignee                                         Assignee
	AITestInstructions, HumanTestInstructions        string
	SourceTaskID                                     *string
	Position                                         float64
	DependencyIDs                                    []string
	Properties                                       map[string]string
}

type TaskUpdate struct {
	Title, Description, Priority, TestingMode *string
	AITestInstructions, HumanTestInstructions *string
	Position                                  *float64
	DependencyIDs                             *[]string
	Properties                                *map[string]string
}

type WorkerInput struct {
	Name, Slug, Description, Kind, Capabilities string
	Enabled                                     bool
}

type PropertyDefinitionInput struct {
	Name, Type, Options, Visibility string
}

type FeatureInput struct {
	Title, Description, Reason, Priority, TestingMode string
	Assignee                                          Assignee
	SourceTaskID                                      *string
	DependsOn                                         []string
	SuggestedProperties                               map[string]string
}

type ArtifactInput struct{ Name, Kind, Path, URL, Description string }

type TestRunInput struct {
	Runner, Type, Command, Status, Summary, OutputExcerpt string
	DurationMS                                            *int64
}

type UsageInput struct {
	Provider, Model, ReasoningEffort, Source                 *string
	InputTokens, OutputTokens, ReasoningTokens, CachedTokens *int64
	ModelCalls, ToolCalls, WallTimeMS                        *int64
}

type AgentContext struct {
	Project Project
	Worker  Worker
	Session WorkerSession
}

type Board map[string][]Task

type TaskDetails struct {
	Task
	History      []HistoryEntry
	TestRuns     []TestRun
	Usage        []UsageEvent
	SpawnedTasks []Task
}

func Now() string { return time.Now().UTC().Format(time.RFC3339Nano) }
