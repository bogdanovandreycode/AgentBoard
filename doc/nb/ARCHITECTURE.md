# Arkitektur og API

AgentBoard er en lokal Go-prosess som kjører SQLite. HTTP API for mennesker, MCP-adapteren for AI og nettgrensesnittet deler felles tjenestelogikk. SQLite - tilstandskilde; frontend omgår ikke serveren. AI-verktøy har en egen rettighetsoverflate og `Backlog`/`Complete`-statusen i MCP er ikke tilgjengelig. Historieposten fra `System` opprettes av applikasjonen selv.

## Komponenter

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - domenemodeller og feil.
- `internal/service` - oppgaveoverganger, tillatelser, import og innstillinger.
- `internal/persistence` - SQLite og migreringer.
- `internal/httpapi` - Human HTTP API.
- `internal/mcpserver` - MCP-verktøy for en separat arbeider.
- `web` - React/TypeScript UI; sammenstillingen ender opp i `internal/webui/dist` og er inkludert i EXE.

## Grunnleggende HTTP-ruter

| Metode og vei | Destinasjon |
| --- | --- |
| `GET /api/health` | Serversjekk. |
| `GET /api/projects` | Registrerte prosjekter. |
| `GET /api/projects/{id}/board` | Prosjektstyre. |
| `GET/PUT /api/projects/{id}/settings` | Innstillinger, rekkefølge og navn på kolonner. |
| `POST /api/projects/{id}/tasks` | Lag en oppgave. |
| `POST /api/projects/{id}/tasks/import` | Atomisk importer JSON versjon 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Kort, endre, slett. |
| `POST /api/tasks/{id}/move` | Flytte av en person. |
| `GET/POST /api/projects/{id}/workers` | Liste og opprettelse av en arbeider. |
| `GET /api/workers/{id}/mcp/check` | Intern testing av MCP-håndtrykk og verktøy. |
| `GET /api/workers/{id}/sessions` | Sesjonsdiagnostikk. |
| `GET/POST /api/projects/{id}/properties` | Egendefinerte egenskaper. |

HTTP er for den lokale pålitelige brukeren. Ikke publiser en nettport til Internett uten egen autentisering, nettverksbegrensninger og HTTPS. MCP-serveren lanseres ved hjelp av `stdio` for et spesifikt prosjekt og arbeider; start med `get_my_board`. Tillatelsene begrenser handlinger i AgentBoard, men erstatter ikke AI-klientens filsystemsandkasse.

Brukerkolonner lagres separat fra `tasks.state`: kjernen holder oppgaven i en slik kolonne i `backlog`, og `tasks.board_column` bestemmer plassen på Human-tavlen. Dette opprettholder den samme AI-overgangsmodellen. Fjerning av en kolonne sletter `board_column`-oppgaver.
