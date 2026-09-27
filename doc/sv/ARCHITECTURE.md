# Arkitektur och API

AgentBoard är en lokal Go-process som kör SQLite. HTTP API för människor, MCP-adaptern för AI och webbgränssnittet delar gemensam tjänstelogik. SQLite - tillståndskälla; frontend förbigår inte servern. AI-verktyg har en separat rättighetsyta och statusen `Backlog`/`Complete` i MCP är inte tillgänglig. Historikposten från `System` skapas av applikationen själv.

## Komponenter

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - domänmodeller och fel.
- `internal/service` - uppgiftsövergångar, behörigheter, import och inställningar.
- `internal/persistence` - SQLite och migrering.
- `internal/httpapi` - Human HTTP API.
- `internal/mcpserver` - MCP-verktyg för en separat arbetare.
- `web` - React/TypeScript UI; sammansättningen hamnar i `internal/webui/dist` och ingår i EXE.

## Grundläggande HTTP-rutter

| Metod och väg | Destination |
| --- | --- |
| `GET /api/health` | Serverkontroll. |
| `GET /api/projects` | Registrerade projekt. |
| `GET /api/projects/{id}/board` | Projekt styrelse. |
| `GET/PUT /api/projects/{id}/settings` | Inställningar, ordning och namn på kolumner. |
| `POST /api/projects/{id}/tasks` | Skapa en uppgift. |
| `POST /api/projects/{id}/tasks/import` | Atomiskt importera JSON version 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Kort, ändra, ta bort. |
| `POST /api/tasks/{id}/move` | Att flytta av en person. |
| `GET/POST /api/projects/{id}/workers` | Lista och skapande av en arbetare. |
| `GET /api/workers/{id}/mcp/check` | Intern testning av MCP-handskakning och verktyg. |
| `GET /api/workers/{id}/sessions` | Sessionsdiagnostik. |
| `GET/POST /api/projects/{id}/properties` | Anpassade egenskaper. |

HTTP är för den lokala betrodda användaren. Publicera inte en webbport till Internet utan egen autentisering, nätverksbegränsningar och HTTPS. MCP-servern startas med `stdio` för ett specifikt projekt och arbetare; börja med `get_my_board`. Dess behörigheter begränsar åtgärder inom AgentBoard, men ersätter inte AI-klientens filsystemssandlåda.

Användarkolumner lagras separat från `tasks.state`: kärnan håller uppgiften i en sådan kolumn i `backlog`, och `tasks.board_column` bestämmer platsen på Human-brädet. Detta bibehåller samma AI-övergångsmodell. Om du tar bort en kolumn rensas `board_column`-uppgifter.
