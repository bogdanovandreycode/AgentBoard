# Architektura a API

AgentBoard je místní proces Go, na kterém běží SQLite. HTTP API pro lidi, MCP adaptér pro AI a webové rozhraní sdílejí společnou logiku služeb. SQLite - stavový zdroj; frontend neobchází server. Nástroje AI mají samostatný povrch práv a stav `Backlog`/`Complete` v MCP není k dispozici. Záznam historie ze `System` vytváří samotná aplikace.

## Komponenty

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - doménové modely a chyby.
- `internal/service` - přechody úloh, oprávnění, import a nastavení.
- `internal/persistence` - SQLite a migrace.
- `internal/httpapi` - Human HTTP API.
- `internal/mcpserver` - Nástroje MCP pro samostatného pracovníka.
- `web` - uživatelské rozhraní React/TypeScript; sestava končí v `internal/webui/dist` a je součástí EXE.

## Základní HTTP trasy

| Metoda a cesta | Destinace |
| --- | --- |
| `GET /api/health` | Kontrola serveru. |
| `GET /api/projects` | Registrované projekty. |
| `GET /api/projects/{id}/board` | Projektová deska. |
| `GET/PUT /api/projects/{id}/settings` | Nastavení, pořadí a názvy sloupců. |
| `POST /api/projects/{id}/tasks` | Vytvořte úkol. |
| `POST /api/projects/{id}/tasks/import` | Atomicky importujte JSON verze 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Karta, změna, smazání. |
| `POST /api/tasks/{id}/move` | Stěhování osobou. |
| `GET/POST /api/projects/{id}/workers` | Seznam a vytvoření pracovníka. |
| `GET /api/workers/{id}/mcp/check` | Interní testování MCP handshake a nástrojů. |
| `GET /api/workers/{id}/sessions` | Diagnostika relace. |
| `GET/POST /api/projects/{id}/properties` | Vlastní vlastnosti. |

HTTP je pro místního důvěryhodného uživatele. Nezveřejňujte webový port na internet bez vlastní autentizace, síťových omezení a HTTPS. Server MCP se spouští pomocí `stdio` pro konkrétní projekt a pracovníka; začněte s `get_my_board`. Jeho oprávnění omezují akce v rámci AgentBoard, ale nenahrazují karanténu souborového systému klienta AI.

Uživatelské sloupce jsou uloženy odděleně od `tasks.state`: jádro uchovává úlohu v takovém sloupci v `backlog` a `tasks.board_column` určuje místo na Human boardu. To zachovává stejný model přechodu AI. Odstraněním sloupce vymažete úlohy `board_column`.
