CREATE TABLE IF NOT EXISTS images (
 id TEXT PRIMARY KEY,
 project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
 name TEXT NOT NULL,
 path TEXT NOT NULL UNIQUE,
 width INTEGER NOT NULL CHECK(width > 0),
 height INTEGER NOT NULL CHECK(height > 0),
 uploaded_task_id TEXT REFERENCES tasks(id) ON DELETE SET NULL,
 created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_images_project ON images(project_id);
