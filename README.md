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

## Status

MVP development.
