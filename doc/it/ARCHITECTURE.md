# Architettura e API

AgentBoard è un processo Go locale che esegue SQLite. L'API HTTP per gli esseri umani, l'adattatore MCP per l'intelligenza artificiale e l'interfaccia web condividono una logica di servizio comune. SQLite: origine dello stato; il frontend non ignora il server. Gli strumenti AI hanno una superficie di diritti separata e lo stato `Backlog`/`Complete` in MCP non è disponibile. Il record della cronologia da `System` viene creato dall'applicazione stessa.

## Componenti

-`cmd/agentboard`-CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - modelli di dominio ed errori.
- `internal/service`: transizioni di attività, autorizzazioni, importazione e impostazioni.
- `internal/persistence` - SQLite e migrazioni.
- `internal/httpapi` - API HTTP umana.
- `internal/mcpserver` - Strumenti MCP per un lavoratore separato.
- `web` - Interfaccia utente React/TypeScript; l'assembly finisce in `internal/webui/dist` ed è incluso nell'EXE.

## Percorsi HTTP di base

| Metodo e percorso | Destinazione |
| --- | --- |
| `GET /api/health` | Controllo del server. |
| `GET /api/projects` | Progetti registrati. |
| `GET /api/projects/{id}/board` | Scheda del progetto. |
| `GET/PUT /api/projects/{id}/settings` | Impostazioni, ordine e nomi delle colonne. |
| `POST /api/projects/{id}/tasks` | Crea un'attività. |
| `POST /api/projects/{id}/tasks/import` | Importa atomicamente JSON versione 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Carta, cambia, cancella. |
| `POST /api/tasks/{id}/move` | Muoversi da parte di una persona. |
| `GET/POST /api/projects/{id}/workers` | Elenco e creazione di un lavoratore. |
| `GET /api/workers/{id}/mcp/check` | Test interni dell'handshake e degli strumenti MCP. |
| `GET /api/workers/{id}/sessions` | Diagnostica della sessione. |
| `GET/POST /api/projects/{id}/properties` | Proprietà personalizzate. |

HTTP è per l'utente fidato locale. Non pubblicare una porta Web su Internet senza la propria autenticazione, restrizioni di rete e HTTPS. Il server MCP viene avviato utilizzando `stdio` per un progetto e un lavoratore specifici; iniziare con `get_my_board`. Le sue autorizzazioni limitano le azioni all'interno di AgentBoard, ma non sostituiscono la sandbox del file system del client AI.

Le colonne utente vengono memorizzate separatamente da `tasks.state`: il nucleo mantiene l'attività in tale colonna in `backlog` e `tasks.board_column` determina la posizione sulla plancia umana. Ciò mantiene lo stesso modello di transizione dell'IA. La rimozione di una colonna cancella le attività `board_column`.
