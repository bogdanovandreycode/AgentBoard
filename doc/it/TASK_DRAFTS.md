# Gestione del progetto e bozze di attività

La pagina **Tutti i progetti** mostra la cartella registrata, il numero di quelli pubblicati
compiti e le ultime attività. L'attività tiene conto dei cambiamenti nelle attività, nella cronologia e
registri delle spese. Il pulsante matita consente di rinominare il progetto o selezionarlo
cartella esistente nella finestra di dialogo del sistema. La modifica del percorso aggiorna la registrazione e il contrassegno;
Gli ID progetto, i lavoratori, le attività e le impostazioni vengono salvati, i file non vengono spostati.
Non è possibile selezionare cartelle già registrate su un altro progetto.

## Impostazioni globali

Apri **Impostazioni globali** dalla pagina o dalla scheda di selezione del progetto
**Impostazioni globali** nelle impostazioni di qualsiasi scheda. Questi parametri sono comuni a tutti i progetti
una base AgentBoard:

- Il salvataggio automatico delle bozze tramite timer è abilitato per impostazione predefinita, ogni **10 secondi**.
- **Il salvataggio all'uscita dal campo** è abilitato per impostazione predefinita e può essere configurato separatamente.
- I controlli del **Contesto del progetto MCP** sono abilitati e disabilitati a livello globale. La loro zona
  azioni e un flag CLI separato sono descritti nella sezione [Workers e MCP](WORKERS_MCP.md).

## Bozza del ciclo di vita

La finestra di dialogo Nuova attività salva il lavoro non completato come bozza. Linea di stato
distingue tra modifiche non salvate, salvataggio in corso, bozza salvata ed errore
conservazione. La bozza non appare sul forum e non sarà disponibile per i lavoratori MCP fino al
pubblicato.

Prima della pubblicazione, utilizza il pulsante **Crea attività** per inserire **Titolo, Stato,
Responsabile e Descrizione**. Prima di ciò, **Salva bozza** salva una bozza senza crearla
compiti. Un modulo vuoto viene scartato per impostazione predefinita quando viene chiuso. Se cancellato
bozza esistente in uno stato vuoto, verrà eliminata la prossima volta che verrà salvata.
Anche la chiusura di una finestra di dialogo non vuota la salva.

Il pulsante **Bozze** sulla bacheca ti consente di continuare o eliminare le bozze salvate.
Sono supportate più bozze. Dopo un riavvio, il file open
bozza. Le bozze sono archiviate in SQLite; Ad ogni modifica, quella locale viene aggiornata
copia nel browser, in modo che le modifiche non salvate sopravvivano ai riavvii o
guasto temporaneo del server. Se la memoria del browser non è disponibile, viene visualizzata una finestra di dialogo
che la bozza deve essere salvata sul server prima della chiusura.

La scrittura di bozze utilizza le revisioni per impedire che una finestra sovrascriva il lavoro
un altro. La risposta di salvataggio persa è coerente se ce n'è già una sul server
dati identici. Se il contenuto è diverso, viene mostrato un conflitto e
**Ripristina come nuova bozza** salva i contenuti locali separatamente. Pubblicazione
atomico e idempotente: ripubblicare la stessa bozza dopo che è andata perduta
la risposta restituisce la stessa attività anziché crearne una copia.
