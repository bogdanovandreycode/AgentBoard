# Arhitectură și API

AgentBoard este un proces Go local care rulează SQLite. API-ul HTTP pentru oameni, adaptorul MCP pentru AI și interfața web au o logică comună de servicii. SQLite - sursă de stare; frontend-ul nu ocolește serverul. Instrumentele AI au o suprafață de drepturi separată, iar starea `Backlog`/`Complete` în MCP nu este disponibilă. Înregistrarea istorică de la `System` este creată de aplicația însăși.

## Componente

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - modele de domenii și erori.
- `internal/service` - tranziții de sarcini, permisiuni, import și setări.
- `internal/persistence` - SQLite și migrații.
- `internal/httpapi` - API-ul HTTP uman.
- `internal/mcpserver` - Instrumente MCP pentru un lucrător separat.
- `web` - Interfața de utilizare React/TypeScript; ansamblul ajunge în `internal/webui/dist` și este inclus în EXE.

## Rute HTTP de bază

| Metodă și cale | Destinație |
| --- | --- |
| `GET /api/health` | Verificare server. |
| `GET /api/projects` | Proiecte înregistrate. |
| `GET /api/projects/{id}/board` | Tabloul de proiect. |
| `GET/PUT /api/projects/{id}/settings` | Setări, ordine și numele coloanelor. |
| `POST /api/projects/{id}/tasks` | Creați o sarcină. |
| `POST /api/projects/{id}/tasks/import` | Importă atomic JSON versiunea 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Card, schimba, șterge. |
| `POST /api/tasks/{id}/move` | Mutarea de către o persoană. |
| `GET/POST /api/projects/{id}/workers` | Lista și crearea unui lucrător. |
| `GET /api/workers/{id}/mcp/check` | Testarea internă a strângerii de mână MCP și a instrumentelor. |
| `GET /api/workers/{id}/sessions` | Diagnosticarea sesiunii. |
| `GET/POST /api/projects/{id}/properties` | Proprietăți personalizate. |

HTTP este pentru utilizatorul local de încredere. Nu publicați un port web pe Internet fără propria sa autentificare, restricții de rețea și HTTPS. Serverul MCP este lansat folosind `stdio` pentru un anumit proiect și lucrător; începe cu `get_my_board`. Permisiunile sale restricționează acțiunile din AgentBoard, dar nu înlocuiesc sandbox-ul sistemului de fișiere al clientului AI.

Coloanele utilizator sunt stocate separat de `tasks.state`: nucleul păstrează sarcina într-o astfel de coloană în `backlog`, iar `tasks.board_column` determină locul pe tabloul uman. Acest lucru menține același model de tranziție AI. Eliminarea unei coloane șterge sarcinile `board_column`.
