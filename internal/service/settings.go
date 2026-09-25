package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

type BoardColumn struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type ProjectSettings struct {
	Columns []BoardColumn `json:"columns"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
	Theme string `json:"theme"`
	BoardRefreshSeconds int `json:"boardRefreshSeconds"`
	WorkerRefreshSeconds int `json:"workerRefreshSeconds"`
}

var builtInColumns = []BoardColumn{
	{"backlog", "Backlog"}, {"features", "Features"}, {"in_progress", "In progress"},
	{"testing", "Testing"}, {"verification", "Verification"}, {"complete", "Complete"},
}
var customColumnID = regexp.MustCompile(`^custom-[a-z0-9][a-z0-9-]{0,39}$`)

func defaultSettings() ProjectSettings {
	return ProjectSettings{Columns: append([]BoardColumn(nil), builtInColumns...), Timezone: "local", Language: "system", Theme: "dark", BoardRefreshSeconds: 1, WorkerRefreshSeconds: 4}
}

func (s *Service) GetProjectSettings(ctx context.Context, projectID string) (ProjectSettings, error) {
	if _, err := s.Store.GetProject(ctx, projectID); err != nil { return ProjectSettings{}, err }
	settings := defaultSettings()
	var raw string
	err := s.Store.DB.QueryRowContext(ctx, `SELECT columns_json FROM project_settings WHERE project_id=?`, projectID).Scan(&raw)
	if err == sql.ErrNoRows { return settings, nil }
	if err != nil { return settings, err }
	if err = json.Unmarshal([]byte(raw), &settings); err != nil { return ProjectSettings{}, err }
	return settings, nil
}

func (s *Service) UpdateProjectSettings(ctx context.Context, projectID string, in ProjectSettings) (ProjectSettings, error) {
	if _, err := s.Store.GetProject(ctx, projectID); err != nil { return ProjectSettings{}, err }
	if len(in.Columns) < len(builtInColumns) || len(in.Columns) > 24 { return ProjectSettings{}, core.ErrInvalidInput }
	seen := make(map[string]bool, len(in.Columns))
	var fixedOrder []string
	for _, c := range in.Columns {
		if seen[c.ID] || strings.TrimSpace(c.Name) == "" || len([]rune(c.Name)) > 40 { return ProjectSettings{}, core.ErrInvalidInput }
		seen[c.ID] = true
		fixed := false
		for _, builtin := range builtInColumns {
			if c.ID == builtin.ID {
				if c.Name != builtin.Name { return ProjectSettings{}, core.ErrInvalidInput }
				fixed = true
				if c.ID != "backlog" && c.ID != "complete" { fixedOrder = append(fixedOrder, c.ID) }
				break
			}
		}
		if !fixed && !customColumnID.MatchString(c.ID) { return ProjectSettings{}, core.ErrInvalidInput }
	}
	if len(seen) < len(builtInColumns) { return ProjectSettings{}, core.ErrInvalidInput }
	for _, builtin := range builtInColumns { if !seen[builtin.ID] { return ProjectSettings{}, core.ErrInvalidInput } }
	if strings.Join(fixedOrder, ",") != "features,in_progress,testing,verification" { return ProjectSettings{}, core.ErrInvalidInput }
	if in.Timezone == "" { in.Timezone = "local" }
	if in.Timezone != "local" { if _, err := time.LoadLocation(in.Timezone); err != nil { return ProjectSettings{}, fmt.Errorf("%w: timezone", core.ErrInvalidInput) } }
	if in.Language == "" { in.Language = "system" }
	if len(in.Language) > 20 || len(in.Language) < 2 { return ProjectSettings{}, core.ErrInvalidInput }
	if in.Theme == "" { in.Theme = "dark" }
	validThemes := map[string]bool{"dark": true, "light": true, "black": true, "ubuntu": true, "windows": true}
	if !validThemes[in.Theme] || in.BoardRefreshSeconds < 1 || in.BoardRefreshSeconds > 60 || in.WorkerRefreshSeconds < 2 || in.WorkerRefreshSeconds > 120 { return ProjectSettings{}, core.ErrInvalidInput }
	previous, err := s.GetProjectSettings(ctx, projectID)
	if err != nil { return ProjectSettings{}, err }
	raw, err := json.Marshal(in)
	if err != nil { return ProjectSettings{}, err }
	tx, err := s.Store.DB.BeginTx(ctx, nil)
	if err != nil { return ProjectSettings{}, err }
	defer tx.Rollback()
	for _, c := range previous.Columns {
		if strings.HasPrefix(c.ID, "custom-") && !seen[c.ID] {
			if _, err = tx.ExecContext(ctx, `UPDATE tasks SET board_column='',updated_at=? WHERE project_id=? AND board_column=?`, core.Now(), projectID, c.ID); err != nil { return ProjectSettings{}, err }
		}
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO project_settings(project_id,columns_json,updated_at) VALUES(?,?,?) ON CONFLICT(project_id) DO UPDATE SET columns_json=excluded.columns_json,updated_at=excluded.updated_at`, projectID, string(raw), core.Now())
	if err != nil { return ProjectSettings{}, err }
	if err = tx.Commit(); err != nil { return ProjectSettings{}, err }
	return in, nil
}

func (s *Service) resolveHumanColumn(ctx context.Context, projectID, target string) (string, string, error) {
	for _, c := range builtInColumns { if target == c.ID { return target, "", nil } }
	settings, err := s.GetProjectSettings(ctx, projectID)
	if err != nil { return "", "", err }
	for _, c := range settings.Columns {
		if target == c.ID && strings.HasPrefix(target, "custom-") { return "backlog", target, nil }
	}
	return "", "", core.ErrInvalidTransition
}
