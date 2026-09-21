ALTER TABLE worker_sessions ADD COLUMN last_activity_at TEXT;
ALTER TABLE worker_sessions ADD COLUMN last_seen_at TEXT;
CREATE INDEX IF NOT EXISTS worker_sessions_worker_started ON worker_sessions(worker_id,started_at);

-- Comments, test results and artifacts are task activity too.
CREATE TRIGGER IF NOT EXISTS history_task_activity AFTER INSERT ON history_entries
BEGIN
 UPDATE tasks SET updated_at=NEW.created_at WHERE id=NEW.task_id;
END;
