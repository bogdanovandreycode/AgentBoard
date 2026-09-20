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
agentboard serve [--addr 127.0.0.1:7337]
agentboard version
```

Operational data is stored in the user's configuration directory, outside managed projects. `agentboard init` writes only `.agentboard/project.json` in the project and is safe to run repeatedly. Workers are logical identities; AgentBoard never launches AI processes automatically.

## Status

Local AgentBoard MVP.
