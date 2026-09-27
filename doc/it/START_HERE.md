# Inizia da qui

## Cos'è AgentBoard

Immagina una normale bacheca con le carte attività. Crei compiti, assegni qualcuno responsabile e supervisioni il lavoro. I lavoratori AI ricevono solo i compiti assegnati tramite MCP e segnalano i progressi. Decidi tu quando l'attività è finalmente pronta.

**Progetto** - una cartella sul computer e una scheda separata. **Compito** - una carta con una descrizione, la persona responsabile e la fase. **Lavoratore** è il nome logico del client AI. **MCP** è il modo in cui il client AI si connette ad AgentBoard. **Scoop** è un gestore di installazione software per Windows.

Nessun codice richiesto. Hai bisogno di Windows, un browser, PowerShell e, affinché l'intelligenza artificiale funzioni, un client AI installato con supporto per i server MCP locali.

## Primo lancio

1. Installare l'applicazione secondo le [istruzioni](INSTALL.md).
2. Creare una cartella di progetto in Explorer, ad esempio `C:\Projects\MyFirstProject`.
3. Aprire questa cartella in Esplora risorse. Fare clic sulla barra degli indirizzi, digitare `powershell` e premere Invio.
4. Nella finestra che si apre, esegui:

```powershell
agentboard init
agentboard open
```

5. Si aprirà un browser con l'indirizzo `http://127.0.0.1:7337`. Lascia aperta la finestra di PowerShell mentre usi la lavagna. La chiusura della finestra arresterà il server locale, ma le attività rimarranno.

Se il comando `agentboard` non viene trovato, chiudi PowerShell e riaprilo dopo aver installato Scoop. Quando si installa da un ZIP, utilizzare il percorso completo di `agentboard.exe`.

## Primo compito

Fai clic su **Nuova attività**, inserisci **Titolo**, se necessario **Descrizione**, quindi **Salva**. Un nuovo compito in `Backlog` è disponibile per gli esseri umani. Affinché l'IA inizi a funzionare, assegna un lavoratore e trasferisci l'attività su `Features`. Descrizione passo passo dei campi - [TASKS.md](TASKS.md).

## Primo operaio

Apri **Lavoratori → Aggiungi lavoratore**. Seleziona il client AI che stai utilizzando (ad esempio Codex o Claude Code), controlla il nome e l'ID breve `Slug`, fai clic su **Salva**. Apri il lavoratore creato: c'è una configurazione MCP, un pulsante di copia e un controllo del server. Copia la configurazione sul client AI secondo [WORKERS_MCP.md](WORKERS_MCP.md). Una volta connesso, chiedi al cliente di chiamare `get_my_board`.

## Cosa leggere dopo

- [Compiti, test, storie e rubriche](TASKS.md)
- [Connessione dei lavoratori e verifica MCP](WORKERS_MCP.md)
- [Importa attività da JSON](IMPORT.md)
- [Lingua, tema, fuso orario e aggiornamento](SETTINGS.md)
- [Problemi tipici](TROUBLESHOOTING.md)
