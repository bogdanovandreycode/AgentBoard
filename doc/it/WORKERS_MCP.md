# Lavoratori e connessione MCP

## Passaggio 1. Crea un lavoratore

In AgentBoard, apri **Lavoratori → Aggiungi lavoratore**. Seleziona un profilo client: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama tramite OpenCode, VS Code Copilot o un altro client MCP. Il profilo riempie le funzionalità tipiche (codice, test, Git); puoi cambiarli. Il nome è visibile sulla scheda, e `Slug` è un identificatore breve senza spazi per il comando MCP. Fare clic su **Salva**.

Un lavoratore corrisponde a una personalità AI. Crea lavoratori diversi per clienti o team diversi. Il profilo delle capacità descrive la specializzazione, ma non estende i diritti dell'IA alle fasi delle attività.

## Passaggio 2. Copia la configurazione

Apri il lavoratore creato. Il blocco **Diagnostica MCP** mostra il frammento di configurazione e il file dove aggiungerlo. Fai clic su **Copia configurazione MCP**. Se il file esiste già, aggiungere il server suggerito all'oggetto `mcpServers`/`servers`/`mcp` esistente senza cancellare gli altri server.

Il comando principale è simile al seguente:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` del lavoratore creato. Il client MCP esegue questo comando da solo quando necessita di strumenti.

### MCP funziona senza interfaccia web

`agentboard mcp` è un processo stdio separato che apre esso stesso il database AgentBoard. Non necessita di `agentboard open`/`serve`, non dipende dalla cartella da cui viene avviata l'interfaccia web e da quante istanze dell'interfaccia web sono in esecuzione. È possibile eseguire più agenti contemporaneamente, ciascuno nel proprio progetto: ogni client esegue il proprio processo MCP e tutti utilizzano in modo sicuro lo stesso database SQLite.

### Come viene selezionato un progetto

Il server MCP viene eseguito in qualsiasi cartella. La prima volta che chiami lo strumento, questo è associato al progetto:

1. `--project PATH`, se specificato (è adatta qualsiasi cartella all'interno del progetto registrato);
2. altrimenti le cartelle di lavoro del cliente (root), se il cliente le segnala;
3. altrimenti la cartella in cui il client ha avviato il server MCP.

`get_my_board` restituisce il progetto e il lavoratore selezionati a `session` insieme a `binding_source`. Se le cartelle di lavoro del client cambiano, la chiamata successiva allo strumento esegue nuovamente l'associazione. Senza `--worker` viene utilizzato l'unico lavoratore abilitato del progetto.

Se il collegamento non è possibile, ogni strumento restituisce un errore che l'agente deve passare all'utente; i dati dell'attività non vengono letti in questo caso:

| Codice | Significato |
| --- | --- |
| `NO_PROJECTS` | Non ci sono progetti nel database. Esegui `agentboard init` nella cartella del progetto e crea un lavoratore. Non è necessario riavviare. |
| `PROJECT_REQUIRED` | La cartella di lavoro è esterna ai progetti registrati. Il messaggio elenca i progetti. |
| `PROJECT_NOT_FOUND` | `--project` non è all'interno di un progetto registrato. |
| `PROJECT_AMBIGUOUS` | Le cartelle di lavoro di un cliente abbracciano più progetti. Aggiungi `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Crea, abilita o seleziona un lavoratore tramite `--worker`. |

Le configurazioni per **Codex** e **Gemini CLI** sono personalizzate e condivise in tutte le cartelle, quindi non hanno `--project`: crea un lavoratore con lo stesso slug in ogni progetto in cui è in esecuzione questo client. I file a livello di progetto (Claude Code, Cursor, OpenCode, VS Code) contengono `--project` del loro progetto.

### Progettazione esplicita e controllo del contesto

Se viene specificato `--project` e la cartella di avvio o di lavoro del client si trova in un progetto diverso, gli strumenti richiedono conferma tramite elicitazione MCP. I client senza richiesta di modulo ricevono `PROJECT_CONTEXT_MISMATCH` prima di accedere alle attività. Correggi la configurazione, rimuovi `--project` per il rilevamento automatico o conferma esplicitamente che intendi lavorare con un altro progetto con il flag `--allow-project-mismatch`. La modifica delle cartelle di lavoro ripristina le conferme precedenti. I progetti definiti automaticamente non richiedono conferma.

**Impostazioni globali → Contesto del progetto MCP** controlla questi controlli per tutti i progetti e i lavoratori dello stesso database. Abilitato per impostazione predefinita. Se disabilitato, MCP utilizza il progetto specificato senza conferme o errori di mancata corrispondenza. Le sessioni già aperte applicano l'impostazione al successivo richiamo dello strumento; la visibilità delle attività, le transizioni e le regole dei test umani continuano a essere seguite.

Dopo il controllo, il blocco diagnostico mostra il percorso assoluto dello `agentboard.exe` in esecuzione. È particolarmente utile quando si installa da un ZIP. Durante l'installazione tramite Scoop, è possibile utilizzare il comando `agentboard` se il client vede lo stesso `PATH`.

## Passaggio 3: aggiungi un server al tuo client

L'interfaccia lavoratore ha già un frammento già pronto. Di seguito è riportata una spiegazione di dove viene utilizzato:

| Cliente | Dove inserire | Come verificare sul lato client |
| --- | --- | --- |
| Codice | `%USERPROFILE%\.codex\config.toml`, sezione `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Codice Claude | `.mcp.json` nella cartella del progetto | `claude mcp list` |
| Gemelli CLI | `%USERPROFILE%\.gemini\settings.json`, oggetto `mcpServers` | `/mcp list` nella CLI Gemini |
| Cursore | Progetto `.cursor\mcp.json` | elenco dei server MCP in Impostazioni cursore |
| OpenCode | Progetto `opencode.json`, oggetto `mcp` | elenco degli strumenti MCP in OpenCode |
| Copilota VS Code | Progetto `.vscode\mcp.json`, oggetto `servers` | comando **MCP: Elenco server** |

Esempi di Codex (configurazione utente, progetto rilevato automaticamente) e Claude Code (file di progetto) con progetto `C:\Projects\MyFirstProject` e lavoratore `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

In JSON, la barra rovesciata di Windows è raddoppiata; un frammento già pronto dell'interfaccia lo fa automaticamente. Le configurazioni generate contengono l'effettivo `--db` del server web, quindi HTTP e MCP utilizzano la stessa base, anche se l'ambiente client è diverso.

**Ollama** fornisce un modello locale, ma non sostituisce il client MCP. Il profilo **Ollama tramite OpenCode** genera una configurazione MCP per OpenCode; configurare separatamente OpenCode sul modello Ollama. È adatto anche un altro client compatibile con MCP con Ollama.

## Passaggio 4: controlla la connessione

1. Apri la carta lavoratore e fai clic su **Controlla di nuovo**. **Verifica del server** esegue un vero eseguibile stdio con questo progetto, lavoratore e base, controlla l'handshake, l'elenco degli strumenti e la posizione del progetto e del lavoratore.
2. Avvia o riavvia il client AI dopo aver aggiunto il file di configurazione. Chiedigli di chiamare `get_my_board`.
3. **Client connesso** verrà visualizzato in AgentBoard e nell'elenco verrà visualizzata una nuova sessione. Solo questo conferma la connessione del tuo cliente. Se sono presenti strumenti, ma il client non è connesso, controlla il percorso del programma, il nome del file di configurazione e la sua sintassi JSON/TOML.

Inizia la tua sessione di lavoro con `get_my_board`. L'intelligenza artificiale non riceve attività da `Backlog` e `Complete`, anche se ne conosce l'ID. L’intelligenza artificiale non può fingere di essere umana e non ha un comando generale “muoviti ovunque”. L'utilizzo del file system all'esterno di AgentBoard dipende dalle funzionalità e dalla sandbox del client selezionato.

Ciascuno strumento dichiara `title` e tutti e quattro i suggerimenti sul comportamento MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) come valori booleani espliciti. Gli strumenti di lettura sono di sola lettura e idempotenti; gli strumenti di registrazione si limitano ad aggiungere record o a spostare attività in avanti, quindi nessuno disturba o va oltre la base AgentBoard locale.

Istruzioni ufficiali per il cliente: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
