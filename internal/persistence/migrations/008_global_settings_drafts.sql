CREATE TABLE IF NOT EXISTS global_settings (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    settings_json TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS task_drafts (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    input_json TEXT NOT NULL,
    revision INTEGER NOT NULL DEFAULT 1,
    published_task_id TEXT REFERENCES tasks(id) ON DELETE SET NULL,
    discarded INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS task_drafts_project ON task_drafts(project_id, discarded, published_task_id);
