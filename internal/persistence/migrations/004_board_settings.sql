CREATE TABLE IF NOT EXISTS project_settings (
 project_id TEXT PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
 columns_json TEXT NOT NULL,
 updated_at TEXT NOT NULL
);
ALTER TABLE tasks ADD COLUMN board_column TEXT NOT NULL DEFAULT '';
ALTER TABLE workers ADD COLUMN harness TEXT NOT NULL DEFAULT 'custom';
CREATE INDEX IF NOT EXISTS idx_tasks_board_column ON tasks(project_id,board_column,position);
