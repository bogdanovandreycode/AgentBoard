# Importa attività da JSON

La scheda **Importa attività** accetta un file JSON con codifica UTF-8. Usalo se stai trasferendo attività da un altro servizio. Fai clic su **Scarica modello JSON**: il modello include i nomi delle proprietà personalizzate correnti e `Slug` del lavoratore disponibile per questo progetto.

## Passo dopo passo

1. Creare i lavoratori (**Workers**) e le proprietà (**Properties**) necessari prima dell'importazione.
2. Scarica il modello, aprilo in un editor di testo, sostituisci gli esempi con le tue attività e salva il file con l'estensione `.json`.
3. Nella scheda **Importa attività**, seleziona un file. L'interfaccia mostrerà il numero di attività.
4. Fare clic su **Importa attività**. Il server controllerà l'intero file: se viene rilevato un errore, non verrà scritta una singola attività da esso. Correggere il messaggio di errore e riprovare.

File minimo:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Esempio con lavoratore, proprietà e dipendenza:

```json
{
  "version": 1,
  "tasks": [
    {
      "key": "design",
      "title": "Prepare the design",
      "description": "## Goal\nPrepare the home page mockup.",
      "state": "features",
      "priority": "high",
      "testing_mode": "hybrid",
      "assignee": { "type": "worker", "worker": "codex" },
      "ai_test_instructions": "Check the build.",
      "human_test_instructions": "Review the page in a browser.",
      "properties": { "Department": "Design" }
    },
    {
      "title": "Approve the design",
      "depends_on": ["design"],
      "assignee": { "type": "human" }
    }
  ]
}
```

`version` dovrebbe essere `1`, array `tasks` - da 1 a 1000 attività. `title` è obbligatorio. Stadi validi: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Priorità: `critical`, `high`, `medium`, `low`; modalità di prova: `ai`, `human`, `hybrid`. Responsabile: `unassigned`, `human` o `worker` con `worker` esistente (Slug o ID). `properties` utilizza i nomi o gli ID delle proprietà già create. `key` è univoco all'interno del file; `depends_on` si riferisce a tali chiavi. Il server crea autonomamente gli ID attività.

Importare nuovamente lo stesso file creerà nuovi problemi, quindi controlla la scheda prima di fare nuovamente clic. Una volta importate, le attività vengono registrate come create dall'uomo; Nella cronologia viene visualizzata una voce di sistema.
