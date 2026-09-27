# Arkitektur og API

AgentBoard er en lokal Go-proces, der kører SQLite. HTTP API for mennesker, MCP-adapteren til AI og webgrænsefladen deler fælles servicelogik. SQLite - tilstandskilde; frontend omgår ikke serveren. AI-værktøjer har en separat rettighedsoverflade, og `Backlog`/`Complete`-statussen i MCP er ikke tilgængelig. Historieposten fra `System` oprettes af selve applikationen.

## Komponenter

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - domænemodeller og fejl.
- `internal/service` - opgaveovergange, tilladelser, import og indstillinger.
- `internal/persistence` - SQLite og migreringer.
- `internal/httpapi` - Human HTTP API.
- `internal/mcpserver` - MCP-værktøjer til en separat arbejder.
- `web` - React/TypeScript UI; samlingen ender i `internal/webui/dist` og er inkluderet i EXE.

## Grundlæggende HTTP-ruter

| Metode og sti | Destination |
| --- | --- |
| `GET /api/health` | Tjek server. |
| `GET /api/projects` | Registrerede projekter. |
| `GET /api/projects/{id}/board` | Projekt bestyrelse. |
| `GET/PUT /api/projects/{id}/settings` | Indstillinger, rækkefølge og navne på kolonner. |
| `POST /api/projects/{id}/tasks` | Opret en opgave. |
| `POST /api/projects/{id}/tasks/import` | Atomisk import af JSON version 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Kort, skift, slet. |
| `POST /api/tasks/{id}/move` | Flytning af en person. |
| `GET/POST /api/projects/{id}/workers` | Liste og oprettelse af en arbejder. |
| `GET /api/workers/{id}/mcp/check` | Intern test af MCP-håndtryk og værktøjer. |
| `GET /api/workers/{id}/sessions` | Sessionsdiagnostik. |
| `GET/POST /api/projects/{id}/properties` | Brugerdefinerede egenskaber. |

HTTP er for den lokale betroede bruger. Udgiv ikke en webport til internettet uden dens egen godkendelse, netværksbegrænsninger og HTTPS. MCP-serveren lanceres ved hjælp af `stdio` til et specifikt projekt og arbejder; start med `get_my_board`. Dens tilladelser begrænser handlinger i AgentBoard, men erstatter ikke AI-klientens filsystemsandbox.

Brugerkolonner gemmes separat fra `tasks.state`: kernen opbevarer opgaven i en sådan kolonne i `backlog`, og `tasks.board_column` bestemmer pladsen på Human board. Dette bevarer den samme AI-overgangsmodel. Fjernelse af en kolonne sletter `board_column`-opgaver.
