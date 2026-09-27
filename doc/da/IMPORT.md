# Importer opgaver fra JSON

Fanen **Importer opgaver** accepterer en JSON-fil i UTF-8-kodning. Brug den, hvis du overfører opgaver fra en anden tjeneste. Klik på **Download JSON-skabelon**: Skabelonen inkluderer de aktuelle brugerdefinerede egenskabsnavne og `Slug` for den tilgængelige arbejder for dette projekt.

## Trin for skridt

1. Opret de nødvendige arbejdere (**Workers**) og egenskaber (**Properties**) før import.
2. Download skabelonen, åbn den i en teksteditor, udskift eksemplerne med dine egne opgaver og gem filen med `.json`-udvidelsen.
3. Vælg en fil på fanen **Importer opgaver**. Grænsefladen vil vise antallet af opgaver.
4. Klik på **Importer opgaver**. Serveren vil kontrollere hele filen: Hvis der opdages en fejl, vil der ikke blive skrevet en eneste opgave fra den. Ret fejlmeddelelsen og prøv igen.

Minimum fil:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Eksempel med arbejder, ejendom og afhængighed:

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

`version` skal være `1`, array `tasks` - fra 1 til 1000 opgaver. `title` er påkrævet. Gyldige stadier: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioriteter: `critical`, `high`, `medium`, `low`; testtilstande: `ai`, `human`, `hybrid`. Ansvarlig: `unassigned`, `human` eller `worker` med eksisterende `worker` (Slug eller ID). `properties` bruger navnene eller id'erne på allerede oprettede egenskaber. `key` er unik i filen; `depends_on` refererer til sådanne nøgler. Serveren opretter selv opgave-id'er.

Import af den samme fil igen vil skabe nye problemer, så tjek tavlen, før du klikker igen. Når de importeres, registreres opgaver som menneskeskabte; En systempost vises i historikken.
