# Architektura i API

AgentBoard to lokalny proces Go z uruchomionym SQLite. Interfejs API HTTP dla ludzi, adapter MCP dla sztucznej inteligencji i interfejs sieciowy mają wspólną logikę usług. SQLite - źródło stanu; frontend nie omija serwera. Narzędzia AI mają osobną powierzchnię uprawnień, a status `Backlog`/`Complete` w MCP nie jest dostępny. Zapis historii z `System` jest tworzony przez samą aplikację.

## Komponenty

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - modele domen i błędy.
- `internal/service` - przejścia zadań, uprawnienia, import i ustawienia.
- `internal/persistence` - SQLite i migracje.
- `internal/httpapi` - Ludzkie API HTTP.
- `internal/mcpserver` - Narzędzia MCP dla osobnego pracownika.
- `web` - Interfejs użytkownika React/TypeScript; zespół kończy się w `internal/webui/dist` i jest zawarty w pliku EXE.

## Podstawowe trasy HTTP

| Metoda i ścieżka | Miejsce docelowe |
| --- | --- |
| `GET /api/health` | Kontrola serwera. |
| `GET /api/projects` | Zarejestrowane projekty. |
| `GET /api/projects/{id}/board` | Tablica projektowa. |
| `GET/PUT /api/projects/{id}/settings` | Ustawienia, kolejność i nazwy kolumn. |
| `POST /api/projects/{id}/tasks` | Utwórz zadanie. |
| `POST /api/projects/{id}/tasks/import` | Atomowo importuj wersję JSON 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Karta, zmień, usuń. |
| `POST /api/tasks/{id}/move` | Poruszanie się przez osobę. |
| `GET/POST /api/projects/{id}/workers` | Lista i utworzenie pracownika. |
| `GET /api/workers/{id}/mcp/check` | Wewnętrzne testowanie uzgadniania i narzędzi MCP. |
| `GET /api/workers/{id}/sessions` | Diagnostyka sesji. |
| `GET/POST /api/projects/{id}/properties` | Właściwości niestandardowe. |

HTTP jest przeznaczony dla lokalnego zaufanego użytkownika. Nie publikuj portu internetowego w Internecie bez własnego uwierzytelnienia, ograniczeń sieciowych i protokołu HTTPS. Serwer MCP jest uruchamiany przy użyciu `stdio` dla konkretnego projektu i pracownika; zacznij od `get_my_board`. Jego uprawnienia ograniczają działania w obrębie AgentBoard, ale nie zastępują piaskownicy systemu plików klienta AI.

Kolumny użytkownika są przechowywane oddzielnie od `tasks.state`: rdzeń przechowuje zadanie w takiej kolumnie w `backlog`, a `tasks.board_column` określa miejsce na tablicy ludzi. Pozwala to zachować ten sam model przejścia na sztuczną inteligencję. Usunięcie kolumny usuwa zadania `board_column`.
