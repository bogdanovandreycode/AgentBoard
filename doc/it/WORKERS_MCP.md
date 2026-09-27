# Lavoratori e connessione MCP

## Passaggio 1. Crea un lavoratore

In AgentBoard, apri **Lavoratori → Aggiungi lavoratore**. Seleziona un profilo client: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama tramite OpenCode, VS Code Copilot o un altro client MCP. Il profilo riempie le funzionalità tipiche (codice, test, Git); puoi cambiarli. Il nome è visibile sulla scheda, e `Slug` è un identificatore breve senza spazi per il comando MCP. Fare clic su **Salva**.

Un lavoratore corrisponde a una personalità AI. Crea lavoratori diversi per clienti o team diversi. Il profilo delle capacità descrive la specializzazione, ma non estende i diritti dell'IA alle fasi delle attività.

## Passaggio 2. Copia la configurazione

Apri il lavoratore creato. Il blocco **Diagnostica MCP** mostra il frammento di configurazione e il file dove aggiungerlo. Fai clic su **Copia configurazione MCP**. Se il file esiste già, aggiungere il server suggerito all'oggetto `mcpServers`/`servers`/`mcp` esistente senza cancellare gli altri server.

Il comando principale è simile al seguente:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` dovrebbe puntare alla cartella registrata tramite `agentboard init`. `--worker` — `Slug` del lavoratore creato. Il client MCP esegue questo comando da solo quando necessita di strumenti. Nel browser AgentBoard, il server Web può essere eseguito separatamente.

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

Per Codex e Claude Code, un esempio con il progetto `C:\Projects\MyFirstProject` e il lavoratore `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

In JSON, la barra rovesciata di Windows è raddoppiata; un frammento già pronto dell'interfaccia lo fa automaticamente. Se stai utilizzando `--db` con una base personalizzata, aggiungila alla configurazione MCP `args` e specifica lo stesso percorso di quando avvii `open`.

**Ollama** fornisce un modello locale, ma non sostituisce il client MCP. Il profilo **Ollama tramite OpenCode** genera una configurazione MCP per OpenCode; configurare separatamente OpenCode sul modello Ollama. È adatto anche un altro client compatibile con MCP con Ollama.

## Passaggio 4: controlla la connessione

1. Apri la carta lavoratore e fai clic su **Controlla di nuovo**. Il **controllo del server** dovrebbe mostrare il numero di strumenti MCP. Si tratta di un controllo del protocollo interno e del rilevamento degli strumenti del server.
2. Avvia o riavvia il client AI dopo aver aggiunto il file di configurazione. Chiedigli di chiamare `get_my_board`.
3. **Client connesso** verrà visualizzato in AgentBoard e nell'elenco verrà visualizzata una nuova sessione. Solo questo conferma la connessione del tuo cliente. Se sono presenti strumenti, ma il client non è connesso, controlla il percorso del programma, il nome del file di configurazione e la sua sintassi JSON/TOML.

Inizia la tua sessione di lavoro con `get_my_board`. L'intelligenza artificiale non riceve attività da `Backlog` e `Complete`, anche se ne conosce l'ID. L’intelligenza artificiale non può fingere di essere umana e non ha un comando generale “muoviti ovunque”. L'utilizzo del file system all'esterno di AgentBoard dipende dalle capacità e dalla sandbox del client selezionato.

Istruzioni ufficiali per il cliente: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
