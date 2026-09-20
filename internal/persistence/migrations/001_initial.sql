PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS projects (
 id TEXT PRIMARY KEY, name TEXT NOT NULL, path TEXT NOT NULL UNIQUE, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS workers (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), name TEXT NOT NULL, slug TEXT NOT NULL,
 description TEXT NOT NULL DEFAULT '', kind TEXT NOT NULL DEFAULT 'external', enabled INTEGER NOT NULL DEFAULT 1,
 archived INTEGER NOT NULL DEFAULT 0, capabilities TEXT NOT NULL DEFAULT '{}', created_at TEXT NOT NULL, updated_at TEXT NOT NULL,
 UNIQUE(project_id, slug)
);
CREATE TABLE IF NOT EXISTS worker_sessions (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), worker_id TEXT NOT NULL REFERENCES workers(id),
 started_at TEXT NOT NULL, ended_at TEXT, client_info TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS tasks (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), title TEXT NOT NULL, description TEXT NOT NULL DEFAULT '',
 state TEXT NOT NULL CHECK(state IN ('backlog','features','in_progress','testing','verification','complete')),
 position REAL NOT NULL DEFAULT 1000, priority TEXT NOT NULL DEFAULT 'medium',
 assignee_type TEXT NOT NULL CHECK(assignee_type IN ('unassigned','human','worker')),
 assignee_worker_id TEXT REFERENCES workers(id), created_by_type TEXT NOT NULL CHECK(created_by_type IN ('human','agent')),
 created_by_worker_id TEXT REFERENCES workers(id), created_by_session_id TEXT REFERENCES worker_sessions(id),
 source_task_id TEXT REFERENCES tasks(id), testing_mode TEXT NOT NULL CHECK(testing_mode IN ('ai','human','hybrid')),
 ai_test_instructions TEXT NOT NULL DEFAULT '', human_test_instructions TEXT NOT NULL DEFAULT '',
 created_at TEXT NOT NULL, updated_at TEXT NOT NULL,
 CHECK((assignee_type='worker' AND assignee_worker_id IS NOT NULL) OR (assignee_type!='worker' AND assignee_worker_id IS NULL)),
 CHECK((created_by_type='agent' AND created_by_worker_id IS NOT NULL AND created_by_session_id IS NOT NULL) OR created_by_type='human')
);
CREATE INDEX IF NOT EXISTS idx_tasks_board ON tasks(project_id,state,position);
CREATE INDEX IF NOT EXISTS idx_tasks_worker ON tasks(assignee_worker_id,state);
CREATE TABLE IF NOT EXISTS task_dependencies (
 task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE, depends_on_task_id TEXT NOT NULL REFERENCES tasks(id), created_at TEXT NOT NULL,
 PRIMARY KEY(task_id, depends_on_task_id), CHECK(task_id != depends_on_task_id)
);
CREATE TABLE IF NOT EXISTS task_artifacts (
 id TEXT PRIMARY KEY, task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE, name TEXT NOT NULL, kind TEXT NOT NULL,
 path TEXT, url TEXT, description TEXT NOT NULL DEFAULT '', created_by_type TEXT NOT NULL,
 created_by_worker_id TEXT REFERENCES workers(id), created_by_session_id TEXT REFERENCES worker_sessions(id), created_at TEXT NOT NULL,
 CHECK(path IS NOT NULL OR url IS NOT NULL)
);
CREATE TABLE IF NOT EXISTS history_entries (
 id TEXT PRIMARY KEY, task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE, actor_type TEXT NOT NULL,
 worker_id TEXT REFERENCES workers(id), worker_session_id TEXT REFERENCES worker_sessions(id), entry_type TEXT NOT NULL,
 content TEXT NOT NULL, created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_history_task ON history_entries(task_id,created_at);
CREATE TABLE IF NOT EXISTS test_runs (
 id TEXT PRIMARY KEY, task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE, worker_id TEXT REFERENCES workers(id),
 worker_session_id TEXT REFERENCES worker_sessions(id), runner TEXT NOT NULL, type TEXT NOT NULL, command TEXT,
 status TEXT NOT NULL CHECK(status IN ('passed','failed','skipped','manual_required')), summary TEXT NOT NULL DEFAULT '',
 output_excerpt TEXT NOT NULL DEFAULT '', duration_ms INTEGER, created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS usage_events (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
 worker_id TEXT NOT NULL REFERENCES workers(id), worker_session_id TEXT NOT NULL REFERENCES worker_sessions(id),
 provider TEXT, model TEXT, input_tokens INTEGER, output_tokens INTEGER, reasoning_tokens INTEGER, cached_tokens INTEGER,
 reasoning_effort TEXT, model_calls INTEGER, tool_calls INTEGER, mcp_calls INTEGER, wall_time_ms INTEGER, source TEXT, created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS property_definitions (
 id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id), name TEXT NOT NULL, type TEXT NOT NULL,
 options TEXT NOT NULL DEFAULT '[]', visibility TEXT NOT NULL CHECK(visibility IN ('human_only','agent_read','agent_read_write')),
 created_at TEXT NOT NULL, updated_at TEXT NOT NULL, UNIQUE(project_id,name)
);
CREATE TABLE IF NOT EXISTS task_property_values (
 task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
 property_definition_id TEXT NOT NULL REFERENCES property_definitions(id) ON DELETE CASCADE,
 value TEXT NOT NULL, updated_at TEXT NOT NULL, PRIMARY KEY(task_id,property_definition_id)
);
