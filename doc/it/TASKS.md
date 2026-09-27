# Lavorare con le attività

## Aggiungi attività

Nella scheda **Bacheca**, fai clic su **Nuova attività**. Compila il titolo e la descrizione. La descrizione e le istruzioni di test supportano Markdown: evidenzia il testo e utilizza la barra di formattazione per grassetto, corsivo, collegamento ed elenco. Fare clic su **Salva**.

Campi:

| Campo | Cosa significa |
| --- | --- |
| Titolo | Nome breve dell'attività. |
| Descrizione | Cosa bisogna fare e come capire che il lavoro è pronto. |
| Stato | Fase di lavoro; una nuova attività solitamente inizia da `Backlog`. |
| Priorità | `critical`, `high`, `medium` o `low`. |
| Responsabile | Una persona, un lavoratore specifico o senza appuntamento. |
| Modalità di prova | `AI` - controlla l'IA; `Human` - assegni umani; `Hybrid` - entrambi. |
| Dipendenze | Attività che dovrebbero essere completate prima. |
| Istruzioni per test AI/umani | Istruzioni per il revisore appropriato. |
| Proprietà personalizzate | Campi aggiuntivi creati nella scheda **Proprietà**. |

Per un'attività AI, crea prima un lavoratore, selezionalo in Responsabile, quindi trasferisci la carta su `Features`. L'intelligenza artificiale vede solo le attività assegnategli in quattro fasi da `Features` a `Verification`.

## Sposta attività

Trascina la carta tra le colonne. Puoi passare dalla visualizzazione a tutte le colonne alla visualizzazione a colonne larghe con scorrimento orizzontale. `Backlog` e `Complete` sono controllati dall'uomo. L'intelligenza artificiale può portare avanti l'attività `Features → In progress → Testing → Verification` solo tramite speciali strumenti MCP. Un'attività con test manuali non completati non dovrebbe superare la verifica AI.

## Scheda attività

Fai clic su una scheda per visualizzare la descrizione, il proprietario, le istruzioni del test e le schede per la cronologia, i test, gli artefatti e i costi dell'IA. **Modifica attività** modifica il contenuto. Il commento della persona viene aggiunto alla storia generale. Cerca nei primi filtri ricerche per titolo, descrizione, ID e lavoratore; un pulsante separato apre una ricerca ampia.

## Colonne e proprietà

Nella scheda **Impostazioni** puoi modificare l'ordine delle colonne disponibili per una persona e aggiungerne di tue. Le quattro fasi dell'IA sono fisse e procedono nello stesso ordine. La colonna utente è un luogo per le attività rinviate da una persona: per AI, tale attività ha lo stato `Backlog`. Quando una colonna viene eliminata, le sue attività tornano al normale `Backlog`.

Nella scheda **Proprietà** puoi aggiungere campi come testo, numero, contrassegno, data, selezione e URL. `Human only` nasconde il campo all'AI; `Agent read` consente la lettura, `Agent read/write` consente anche la scrittura tramite strumenti supportati. Ciò non modifica i diritti dell'IA sulle fasi dell'attività.
