# Architektur und API

AgentBoard ist ein lokaler Go-Prozess, auf dem SQLite ausgeführt wird. Die HTTP-API für Menschen, der MCP-Adapter für KI und die Webschnittstelle verwenden eine gemeinsame Dienstlogik. SQLite – Statusquelle; Das Frontend umgeht den Server nicht. KI-Tools verfügen über eine eigene Rechteoberfläche und der Status `Backlog`/`Complete` im MCP ist nicht verfügbar. Der Verlaufsdatensatz von `System` wird von der Anwendung selbst erstellt.

## Komponenten

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` – Domänenmodelle und Fehler.
- `internal/service` – Aufgabenübergänge, Berechtigungen, Import und Einstellungen.
- `internal/persistence` - SQLite und Migrationen.
- `internal/httpapi` – Menschliche HTTP-API.
- `internal/mcpserver` – MCP-Tools für einen separaten Arbeiter.
– `web` – React/TypeScript-Benutzeroberfläche; Die Assembly landet in `internal/webui/dist` und ist in der EXE-Datei enthalten.

## Grundlegende HTTP-Routen

| Methode und Pfad | Ziel |
| --- | --- |
| `GET /api/health` | Serverüberprüfung. |
| `GET /api/projects` | Registrierte Projekte. |
| `GET /api/projects/{id}/board` | Projektvorstand. |
| `GET/PUT /api/projects/{id}/settings` | Einstellungen, Reihenfolge und Namen der Spalten. |
| `POST /api/projects/{id}/tasks` | Erstellen Sie eine Aufgabe. |
| `POST /api/projects/{id}/tasks/import` | JSON Version 1 atomar importieren. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Karte, ändern, löschen. |
| `POST /api/tasks/{id}/move` | Umzug durch eine Person. |
| `GET/POST /api/projects/{id}/workers` | Liste und Erstellung eines Arbeiters. |
| `GET /api/workers/{id}/mcp/check` | Interne Tests von MCP-Handshake und -Tools. |
| `GET /api/workers/{id}/sessions` | Sitzungsdiagnose. |
| `GET/POST /api/projects/{id}/properties` | Benutzerdefinierte Eigenschaften. |

HTTP ist für den lokalen vertrauenswürdigen Benutzer. Veröffentlichen Sie keinen Webport im Internet ohne eigene Authentifizierung, Netzwerkbeschränkungen und HTTPS. Der MCP-Server wird mit `stdio` für ein bestimmtes Projekt und einen bestimmten Worker gestartet. Beginnen Sie mit `get_my_board`. Seine Berechtigungen schränken Aktionen innerhalb des AgentBoard ein, ersetzen jedoch nicht die Dateisystem-Sandbox des AI-Clients.

Benutzerspalten werden getrennt von `tasks.state` gespeichert: Der Kern speichert die Aufgabe in einer solchen Spalte in `backlog` und `tasks.board_column` bestimmt den Platz auf dem Human-Board. Dadurch wird das gleiche KI-Übergangsmodell beibehalten. Durch das Entfernen einer Spalte werden `board_column`-Aufgaben gelöscht.
