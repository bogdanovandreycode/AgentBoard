# Preparare una versione per Scoop

## Cosa è già automatizzato

`scripts/package-scoop.ps1 -Version 0.2.0` esegue `npm ci`, build frontend, `go test ./...`, build Windows `agentboard.exe` con numero di versione, crea uno ZIP e calcola lo SHA-256 di quello ZIP. Dallo **stesso** file crea `release/agentboard.json` con URL, hash, CLI-shim, scorciatoia, `checkver` e `autoupdate`. Il database si trova all'esterno della directory di installazione, quindi `persist` non è necessario nel manifest.

`.github/workflows/release.yml` sul tag `vX.Y.Z` esegue lo stesso pacchetto nel runner di Windows, crea il programma di installazione Inno Setup e allega ZIP, programma di installazione e manifest alla versione GitHub. L'esecuzione manuale di un flusso di lavoro crea solo un artefatto per il test, senza pubblicare una versione.

## Ordine di pubblicazione per il manutentore

1. Verificare che il codice, la documentazione e il numero di versione siano pronti. Definisci la licenza del progetto: manifest attualmente indica `Unknown` perché il file di licenza non è definito nel repository. Se scegli una licenza, aggiungi `LICENSE` e aggiorna `license` nello script prima del rilascio.
2. Eseguire `./scripts/package-scoop.ps1 -Version X.Y.Z` localmente. Controllare l'output di `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` e `agentboard version` dopo il disimballaggio. Non modificare il codice postale dopo aver calcolato l'hash.
3. Crea e invia il tag `vX.Y.Z`. GitHub Actions pubblicherà una release con ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` e manifest. Controlla tutti e tre i file nella pagina Release e lo ZIP SHA-256 dal manifest.
4. Su un computer Windows pulito con Scoop, eseguire `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, quindi `agentboard version`, `agentboard init` e `agentboard open` nella cartella test.
5. Per un feed di aggiornamento permanente, posiziona lo `agentboard.json` generato nel tuo bucket Scoop o offrilo in un bucket pubblico adatto. Controlla nuovamente per `scoop update agentboard` dopo la prossima versione. Il comando di installazione da un URL è adatto per la prima conoscenza, mentre il bucket è più conveniente per gli aggiornamenti.

Non sostituire manualmente il valore casuale `hash`: Scoop controlla il contenuto dello ZIP scaricato.

Per i controlli del manifest, concentrati su [formato Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [creazione di manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) e [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
