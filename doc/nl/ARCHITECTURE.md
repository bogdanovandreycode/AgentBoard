# Architectuur en API

AgentBoard is een lokaal Go-proces waarop SQLite wordt uitgevoerd. De HTTP API voor mensen, de MCP-adapter voor AI en de webinterface delen gemeenschappelijke servicelogica. SQLite - statusbron; frontend omzeilt de server niet. AI-tools hebben een apart rechtenoppervlak en de `Backlog`/`Complete`-status in MCP is niet beschikbaar. Het historierecord van `System` wordt door de applicatie zelf aangemaakt.

## Componenten

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - domeinmodellen en fouten.
- `internal/service` - taakovergangen, machtigingen, import en instellingen.
- `internal/persistence` - SQLite en migraties.
- `internal/httpapi` - Menselijke HTTP-API.
- `internal/mcpserver` - MCP-tools voor een afzonderlijke werknemer.
- `web` - React/TypeScript-gebruikersinterface; de assembly komt terecht in `internal/webui/dist` en wordt opgenomen in de EXE.

## Basis HTTP-routes

| Methode en pad | Bestemming |
| --- | --- |
| `GET /api/health` | Servercontrole. |
| `GET /api/projects` | Geregistreerde projecten. |
| `GET /api/projects/{id}/board` | Projectbord. |
| `GET/PUT /api/projects/{id}/settings` | Instellingen, volgorde en namen van kolommen. |
| `POST /api/projects/{id}/tasks` | Maak een taak. |
| `POST /api/projects/{id}/tasks/import` | JSON-versie 1 atomair importeren. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Kaart, wijzigen, verwijderen. |
| `POST /api/tasks/{id}/move` | Bewegen door een persoon. |
| `GET/POST /api/projects/{id}/workers` | Lijst en creatie van een werknemer. |
| `GET /api/workers/{id}/mcp/check` | Intern testen van MCP-handshake en tools. |
| `GET /api/workers/{id}/sessions` | Sessiediagnostiek. |
| `GET/POST /api/projects/{id}/properties` | Aangepaste eigenschappen. |

HTTP is voor de lokale vertrouwde gebruiker. Publiceer geen webpoort naar internet zonder eigen authenticatie, netwerkbeperkingen en HTTPS. De MCP-server wordt gestart met `stdio` voor een specifiek project en een specifieke medewerker; begin met `get_my_board`. De machtigingen beperken acties binnen het AgentBoard, maar vervangen niet de sandbox van het bestandssysteem van de AI-client.

Gebruikerskolommen worden apart van `tasks.state` opgeslagen: de kern bewaart de taak in zo'n kolom in `backlog`, en `tasks.board_column` bepaalt de plaats op het Human-bord. Dit handhaaft hetzelfde AI-transitiemodel. Als u een kolom verwijdert, worden de `board_column`-taken gewist.
