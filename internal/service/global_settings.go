package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/bogdanovandreycode/agentboard/internal/core"
)

type GlobalSettings struct {
	DraftAutosaveEnabled bool `json:"draftAutosaveEnabled"`
	DraftAutosaveSeconds int  `json:"draftAutosaveSeconds"`
	DraftSaveOnBlur      bool `json:"draftSaveOnBlur"`
	CheckProjectContext  bool `json:"checkProjectContext"`
}

func (s *Service) GetGlobalSettings(ctx context.Context) (GlobalSettings, error) {
	value := GlobalSettings{DraftAutosaveEnabled: true, DraftAutosaveSeconds: 10, DraftSaveOnBlur: true, CheckProjectContext: true}
	var raw string
	err := s.Store.DB.QueryRowContext(ctx, `SELECT settings_json FROM global_settings WHERE id=1`).Scan(&raw)
	if err == sql.ErrNoRows {
		return value, nil
	}
	if err != nil {
		return value, err
	}
	err = json.Unmarshal([]byte(raw), &value)
	return value, err
}

func (s *Service) UpdateGlobalSettings(ctx context.Context, value GlobalSettings) (GlobalSettings, error) {
	if value.DraftAutosaveSeconds < 1 || value.DraftAutosaveSeconds > 3600 {
		return GlobalSettings{}, core.ErrInvalidInput
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return value, err
	}
	_, err = s.Store.DB.ExecContext(ctx, `INSERT INTO global_settings(id,settings_json,updated_at) VALUES(1,?,?) ON CONFLICT(id) DO UPDATE SET settings_json=excluded.settings_json,updated_at=excluded.updated_at`, string(raw), core.Now())
	return value, err
}
