# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Anteprima AgentBoard](../../assets/social-preview.png)

**[Download per Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Sito della documentazione](https://bogdanovandreycode.github.io/AgentBoard/) · [Licenza MIT](../../LICENSE)**

AgentBoard è un task board locale in cui i lavoratori umani e quelli dell'intelligenza artificiale lavorano su compiti comuni, ma hanno diritti diversi. L'applicazione viene avviata con un file `agentboard.exe`, apre l'interfaccia web nel browser e fornisce ai lavoratori un server MCP separato tramite `stdio`. I dati rimangono sul tuo computer.

**[Inizia da zero](START_HERE.md) · [Lavorare con le attività](TASKS.md) · [Connessione AI tramite MCP](WORKERS_MCP.md) · [Importa JSON](IMPORT.md) · [Impostazioni](SETTINGS.md) · [Risoluzione dei problemi](TROUBLESHOOTING.md)**

## Caratteristiche MVP

- Scheda di progetto locale con ricerca di attività, importazione JSON, colonne e proprietà personalizzate.
- Accesso MCP per un lavoratore specifico con transizioni AI limitate e accettazione umana finale delle attività.
- Cronologia generale delle attività, istruzioni di controllo, artefatti, costi dell'IA e diagnostica della connessione del lavoratore.
- Programma di installazione di Windows, manifest ZIP e Scoop; l'interfaccia web è integrata nel file eseguibile.

AgentBoard è progettato per un utente locale fidato. L'applicazione non ospita progetti nel cloud e non avvia autonomamente client AI; se necessario, collegare al lavoratore un client compatibile MCP.

## Tra cinque minuti

1. Scaricare il programma di installazione `agentboard-VERSION-windows-amd64-setup.exe` da [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Suggerirà una cartella (per impostazione predefinita `C:\AI\AgentBoard`) e la aggiungerà a `PATH`. Sono disponibili anche Scoop e ZIP.
2. Apri PowerShell nella cartella del progetto, ad esempio `C:\Projects\MyApp`.
3. Eseguire `agentboard init` (per ZIP: percorso completo a `agentboard.exe` e `init`).
4. Eseguire `agentboard open`. Si aprirà `http://127.0.0.1:7337`.
5. Aggiungi un'attività utilizzando il pulsante **Nuova attività**. Per un lavoratore AI, apri **Lavoratori → Aggiungi lavoratore**, seleziona il profilo cliente e copia la configurazione MCP.

Se non disponi già di una cartella di progetto, creane una in Esplora risorse. Un progetto può essere qualsiasi cartella, anche senza Git e codice.

## Installazione tramite Scoop

In PowerShell con [Scoop](https://scoop.sh/)] già installato dopo il rilascio:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Per gli sviluppatori c'è [build from source ](INSTALL.md). Il flusso di lavoro di rilascio crea un programma di installazione, un manifest ZIP e Scoop con SHA-256 dallo stesso artefatto. Aggiornamento della versione installata tramite Scoop: `scoop update agentboard` dopo aver aggiunto manifest al bucket; dettagli - [preparazione rilascio](SCOOP_RELEASE.md).

## Come è strutturato il tabellone

`Backlog → Features → In progress → Testing → Verification → Complete`

Una persona può spostare le attività sul tabellone. L'IA può spostare solo `Features → In progress → Testing → Verification`; l'accettazione finale in `Complete` viene eseguita da un essere umano. Le colonne utente sono destinate agli utenti umani: l'attività al loro interno rimane nello stato `Backlog` per MCP. Modalità di test: AI, Umano e Ibrido. La storia, i test, gli artefatti e i costi dell'IA sono allegati all'attività.

## Squadre

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registra la cartella e vi scrive solo `.agentboard/project.json`. I dati di produzione SQLite si trovano nella directory di configurazione utente di Windows (`%AppData%\AgentBoard\agentboard.db`), all'esterno del progetto e all'esterno dell'installazione di Scoop. La rimozione o l'aggiornamento del pacchetto non dovrebbe rimuovere questi dati. Prima del trasferimento su un altro computer, eseguire una copia del database mentre AgentBoard è arrestato.

Lavoratori - conti logici; AgentBoard stesso non esegue Codex, Claude o qualsiasi altro client AI. Il client avvia un processo MCP locale per un lavoratore specifico. MCP è il limite delle autorizzazioni dell'applicazione e, per l'isolamento dei file, utilizza la sandbox del client AI.

## Per gli sviluppatori

Stack: Go, SQLite, SDK ufficiale MCP Go, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Costruisci prima il frontend, poi Vai: `./scripts/build.ps1`. I file Web sono inclusi nel file binario tramite `go:embed`. L'architettura e l'API sono descritte in [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Licenza

AgentBoard è un progetto gratuito e open source con [licenza MIT](../../LICENSE). L'uso commerciale, la modifica, il fork e la ridistribuzione sono consentiti a condizione che vengano mantenuti l'avviso di copyright e la licenza.
