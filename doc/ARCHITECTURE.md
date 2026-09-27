# Architecture and API

AgentBoard runs locally as a Go process backed by SQLite. The Human HTTP API, AI MCP adapter and web interface use the same core services. SQLite is the source of truth. AI tools have separate permissions and cannot access `Backlog` or `Complete`. Only AgentBoard creates `System` history entries.

## Components

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - domain models and errors.
- `internal/service` - task transitions, permissions, imports and settings.
- `internal/persistence` - SQLite and migrations.
- `internal/httpapi` - Human HTTP API.
- `internal/mcpserver` - MCP tools for each worker.
- `web` - React/TypeScript UI; its build is copied to `internal/webui/dist` and embedded in the EXE.

## Basic HTTP routes

| Method and path | Purpose |
| --- | --- |
| `GET /api/health` | Server check. |
| `GET /api/projects` | Registered projects. |
| `GET /api/projects/{id}/board` | Project board. |
| `GET/PUT /api/projects/{id}/settings` | Settings, order and names of columns. |
| `POST /api/projects/{id}/tasks` | Create a task. |
| `POST /api/projects/{id}/tasks/import` | Atomically import JSON version 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Card, change, delete. |
| `POST /api/tasks/{id}/move` | Move a task as a human. |
| `GET/POST /api/projects/{id}/workers` | List and creation of a worker. |
| `GET /api/workers/{id}/mcp/check` | Internal testing of MCP handshake and tools. |
| `GET /api/workers/{id}/sessions` | Session diagnostics. |
| `GET/POST /api/projects/{id}/properties` | Custom properties. |

HTTP is intended for a trusted local user. Do not expose the web port to the Internet without authentication, network restrictions and HTTPS. The MCP server runs over `stdio` for a specific project and worker; start with `get_my_board`. Its permissions apply within AgentBoard and do not replace the AI client's filesystem sandbox.

Custom columns are stored separately from `tasks.state`: the core keeps these tasks in `backlog`, while `tasks.board_column` determines their position on the Human board. Removing a column clears `board_column` on its tasks.
