# Impostazioni

Apri **Impostazioni** nel menu laterale del progetto selezionato. Le modifiche vengono salvate con il pulsante **Salva impostazioni** e archiviate per il progetto nel database AgentBoard.

## Linguaggio e aspetto

Il campo **Lingua** apre un elenco a discesa di ricerca. Per impostazione predefinita, è selezionata l'opzione **Segui sistema**, che accetta la lingua del browser. Lingue disponibili dagli screenshot: arabo, portoghese (Brasile), cinese semplificato, ceco, danese, olandese, inglese, finlandese, francese, tedesco, italiano, giapponese, coreano, norvegese Bokmål, polacco, russo, spagnolo, svedese, turco, ucraino, vietnamita; inoltre bielorusso, rumeno e bulgaro. **Segui sistema** accetta la lingua del browser. Le firme principali vengono tradotte manualmente, le restanti righe dell'interfaccia utente hanno una traduzione automatica preliminare. Prima del rilascio al pubblico, è consigliabile correggere le traduzioni da parte di madrelingua; i nomi dei client, i comandi, i campi JSON e i dati utente rimangono senza traduzione.

**Schema colori**: Scuro (tema originale), Chiaro, Nero, Ubuntu e Windows. **Fuso orario** controlla la visualizzazione delle date; i dati continuano ad essere archiviati in UTC. Il **Fuso orario del sistema** utilizza le impostazioni del computer.

## Colonne

I quattro stadi `Features`, `In progress`, `Testing`, `Verification` sono fissati in questo ordine: non possono essere rinominati o cancellati. Altre colonne possono essere riorganizzate utilizzando le frecce disponibili. Inserisci un nome e fai clic su **Aggiungi colonna** per creare una colonna personalizzata. È progettato per le attività umane: l'intelligenza artificiale vede un'attività come `Backlog` e non la riceve tramite MCP. L'eliminazione di una colonna durante il salvataggio trasferisce le sue attività al normale `Backlog`.

## Web e MCP

**Aggiornamento tabellone** imposta la frequenza di aggiornamento del tabellone e delle carte in secondi (1–60). **Aggiornamento lavoratore** aggiorna lo stato dei lavoratori (2-120 secondi). Si tratta di un polling dell'interfaccia web, non della frequenza di attivazione dell'AI. I lavoratori non si avviano automaticamente.

L'indirizzo del server web viene impostato all'avvio della CLI, ad esempio `agentboard open --addr 127.0.0.1:7444`. Per modificare l'indirizzo è necessario riavviare il server. Il valore predefinito è `127.0.0.1:7337`. MCP funziona tramite un comando locale separato `agentboard mcp --project ... --worker ...` ed è indipendente dalla porta web. Per una base non standard, specificare lo stesso `--db` in tutti i comandi. Copia la configurazione di ciascun cliente dalla scheda lavoratore.
