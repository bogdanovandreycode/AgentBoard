CREATE TABLE IF NOT EXISTS files (
 id TEXT PRIMARY KEY,
 project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
 name TEXT NOT NULL,
 path TEXT NOT NULL UNIQUE,
 mime_type TEXT NOT NULL,
 size INTEGER NOT NULL CHECK(size >= 0),
 uploaded_task_id TEXT REFERENCES tasks(id) ON DELETE SET NULL,
 created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_files_project ON files(project_id);
