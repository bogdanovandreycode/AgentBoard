# Installazione e avvio su Windows

## Installatore (consigliato)

Scarica `agentboard-VERSION-windows-amd64-setup.exe` dalla pagina [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). La procedura guidata di installazione suggerirà una cartella; il valore predefinito è `C:\AI\AgentBoard`. Copierà `agentboard.exe` e la documentazione, creerà un collegamento e aggiungerà la cartella selezionata al sistema `PATH`. Dopo l'installazione, aprire un nuovo terminale in modo che diventi disponibile il comando `agentboard`.

Nella cartella del tuo progetto esegui:

```powershell
agentboard init
agentboard open
```

L'interfaccia è integrata nello `agentboard.exe`; Non è necessaria alcuna installazione separata di Go o Node.js. I dati vengono memorizzati in `%AppData%\AgentBoard` e vengono conservati quando il programma viene aggiornato o disinstallato. La disinstallazione tramite "Applicazioni installate" rimuove i collegamenti e la voce da `PATH`.

## Installazione di Scoop

Se Scoop non è già installato, apri PowerShell come utente normale e segui le [istruzioni ufficiali Scoop](https://scoop.sh/). Se sono presenti restrizioni sul tuo computer aziendale, contatta il tuo amministratore; AgentBoard può essere lanciato anche da ZIP senza Scoop.

In alternativa, installa AgentBoard tramite Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

La disponibilità del manifest nella versione GitHub può essere verificata su [pagina Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Dopo aver aggiunto il manifest a Scoop, il bucket può essere installato in base al nome del bucket e aggiornato con il comando `scoop update agentboard`.

## ZIP senza scoop

Scarica `agentboard-VERSION-windows-amd64.zip` da Releases, scompattalo, ad esempio, in `C:\Tools\AgentBoard`. In PowerShell nella cartella del progetto:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Per un client AI, specificare il percorso completo di `agentboard.exe` nella configurazione MCP se il programma non si trova in `PATH`.

## Compila dal sorgente

Installa la versione Go da `go.mod` e Node.js 22 o successiva. In PowerShell nella radice del repository:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` assembla il frontend in `internal/webui/dist`, esegue i test Go e assembla uno `agentboard.exe`. L'ordine è importante: l'interfaccia è integrata nel binario quando viene creato Go. Se il vecchio `agentboard.exe open` è in esecuzione, interromperlo prima della ricostruzione (Ctrl+C), altrimenti Windows non consentirà di sostituire il file.

## Dove sono i dati

- `%AppData%\AgentBoard\agentboard.db` - attività, progetti, lavoratori e impostazioni. È possibile specificare un file diverso con il flag `--db`, ma `init`, `open`/`serve` e `mcp` devono avere lo **stesso percorso**.
- `<ваш проект>\.agentboard\project.json` — identificatore del progetto. Questo file non contiene attività.
- Per impostazione predefinita, il server ascolta solo `127.0.0.1:7337`. Specificare un altro indirizzo `--addr` prima del percorso del progetto: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` apre la pagina del progetto e riutilizza un server già in esecuzione a quell'indirizzo. Se più progetti sono aperti in un browser, selezionali dall'elenco a sinistra.
