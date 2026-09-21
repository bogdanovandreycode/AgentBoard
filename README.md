# AgentBoard

AgentBoard is a local-first task control plane designed for collaboration between a human project owner and AI agents.

## Core idea

Human and AI work with the same tasks and shared task history, but have different workflow capabilities.

Task workflow:

backlog -> features -> in_progress -> testing -> verification -> complete

AI can only move forward through:

features -> in_progress -> testing -> verification

Final acceptance belongs to the human project owner.

## Features planned for MVP

- Persistent local projects
- Kanban task board
- Shared Human / AI / System task history
- Strict AI workflow capabilities
- Human-only final acceptance
- AI / Human / Hybrid testing modes
- Structured test results
- MCP integration
- Human HTTP API
- AI usage telemetry
- Custom task properties
- Task dependencies
- Git metadata integration
- Embedded React Web UI

## Stack

Backend:
- Go
- SQLite
- database/sql
- chi
- official Model Context Protocol Go SDK

Frontend:
- React
- TypeScript
- Vite
- TanStack Query
- dnd-kit
- Tailwind CSS

## Architecture

AgentBoard Core owns:
- workflow rules;
- persistence;
- permissions;
- history;
- testing rules;
- telemetry.

Adapters:
- HTTP API for Human interaction;
- MCP for AI interaction;
- CLI for local administration.

SQLite is stored outside managed workspaces.

Managed projects contain only lightweight AgentBoard project metadata.

## Quick start

```powershell
go build -o agentboard.exe ./cmd/agentboard
./agentboard.exe init
./agentboard.exe open
```

The UI opens at `http://127.0.0.1:7337`. Create one or more Workers, assign tasks, and launch an external agent with a Worker-scoped MCP server:

```powershell
./agentboard.exe mcp --project C:\path\to\project --worker codex
```

Other commands:

```text
agentboard open [--addr 127.0.0.1:7337] [project-path]
agentboard serve [--addr 127.0.0.1:7337]
agentboard version
```

Operational data is stored in the user's configuration directory, outside managed projects. `agentboard init` writes only `.agentboard/project.json` in the project and is safe to run repeatedly. Workers are logical identities; AgentBoard never launches AI processes automatically.

## Status

Local AgentBoard MVP.

## Project launcher and live updates

`agentboard open` finds `.agentboard/project.json` in the current directory or its
parents and opens `/?project=<id>`. `agentboard open .` and an explicit directory
work too. Outside a registered project, the browser opens a project chooser.
Invalid markers are reported; an unknown project ID in the URL shows the chooser
with an explanation. Browser Back/Forward and the project selector update the URL.
An already running AgentBoard server is reused. CLI flags go before the directory.

The visible board and open task (including history, tests, artifacts and usage)
refresh every second. Workers refresh every 4 seconds and projects every 10 seconds.
Polling pauses in background tabs and resumes on focus; unsaved form input is kept.

Worker forms provide Generic, Coding, Graphics and Reviewer capability presets,
with editable JSON for custom capabilities. Capabilities describe the worker;
they do not grant workflow permissions.

## MCP discovery and diagnostics

The MCP initialize response includes worker instructions. The read-only resource
`agentboard://current-worker` identifies the configured worker and project and
points clients to `get_my_board`. All 11 worker tools and existing permissions
remain available. Resources are discovery metadata, not a substitute for tools.

Workers show active/offline status, total session and MCP-call counts, and the last
MCP call time. Open a worker to inspect its latest 50 sessions or copy its Codex
TOML configuration (requires `agentboard` on PATH). The human REST API adds
`GET /api/workers/{workerID}/sessions`; worker responses add `ActiveSessionCount`,
`SessionCount`, `MCPCalls`, and nullable `LastActivityAt`. Sessions include nullable
`LastActivityAt` and `LastSeenAt`, plus `MCPCalls`. `GET /api/health` identifies the server.

MCP processes send a heartbeat every 15 seconds. Sessions without a heartbeat for
90 seconds are offline, including processes that crashed without setting `EndedAt`.
Last MCP activity measures tool calls, not heartbeats or resource discovery.
Pre-upgrade sessions have unknown activity timestamps and are not marked active
until a new process sends heartbeats. Migration 3 preserves all existing data.

Never read or modify AgentBoard SQLite storage directly for worker tasks: MCP is
the supported worker interface. This is a workflow/capability boundary, not an OS
security boundary. Filesystem isolation requires running the worker in a sandbox.

For self-dogfooding, initialize this repository, register a worker through the UI,
and use MCP to create/advance its assigned features, record tests and add history.
Final acceptance in Complete remains human-owned.
