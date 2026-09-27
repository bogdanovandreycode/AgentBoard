# Risoluzione dei problemi

| Sintomo | Cosa controllare |
| --- | --- |
| `agentboard` non trovato | Riavvia PowerShell dopo Scoop. Quando si installa da un ZIP, utilizzare il percorso completo di `agentboard.exe`. |
| Pagina Web che mostra la vecchia interfaccia dopo la compilazione | Arresta il server in esecuzione Ctrl+C. Esegui `./scripts/build.ps1`, avvia un nuovo binario. Vite deve essere creato prima di Go perché l'interfaccia è incorporata nell'EXE. Aggiorna la pagina Ctrl+F5. |
| Porta 7337 occupata | AgentBoard potrebbe essere già in esecuzione. Apri `http://127.0.0.1:7337` o termina il vecchio processo. Per altre porte utilizzare `--addr`. |
| Progetto non trovato | Nella cartella desiderata, eseguire `agentboard init`. Quindi `agentboard open` da esso o `agentboard open C:\путь\к\проекту`. |
| Il lavoratore non vede l'attività | L'attività deve essere assegnata a questo particolare lavoratore e trovarsi in `Features`, `In progress`, `Testing` o `Verification`. L'intelligenza artificiale non vede `Backlog`, `Complete` e le colonne personalizzate. |
| Il controllo del server MCP esiste, ma il client non è connesso | Il controllo del server non verifica le impostazioni del client esterno. Riavvia il client, controlla il suo file di configurazione, il percorso di `agentboard.exe`, `--project`, `--worker` e `--db` generale. Chiedi di chiamare `get_my_board`. |
| Lavoratore offline | Il client potrebbe essere terminato o non aver ancora avviato MCP. Dopo 90 secondi senza battito cardiaco, la sessione viene considerata disconnessa. |
| File JSON non importato | Controlla `version: 1`, `title` richiesto, lavoratori `Slug` esistenti e nomi di proprietà. JSON non consente commenti o virgole finali. |
| Impossibile ricostruire `agentboard.exe` | Windows non può sostituire un EXE in esecuzione. Arresta il server Ctrl+C e riprova la compilazione. |
| Le attività sono scomparse dopo l'aggiornamento | Verifica che `--db` non punti a un altro file e di aver effettuato l'accesso come lo stesso utente Windows. La base predefinita è `%AppData%\AgentBoard`. |

Se l'errore non viene descritto, raccogli il testo esatto del messaggio, le versioni `agentboard version` e Windows e i passaggi per riprovare. Non pubblicare dati di progetti privati ​​o contenuti di database in un numero aperto.
