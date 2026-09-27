# Importera uppgifter från JSON

Fliken **Importera uppgifter** accepterar en JSON-fil i UTF-8-kodning. Använd den om du överför uppgifter från en annan tjänst. Klicka på **Ladda ner JSON-mall**: Mallen innehåller de aktuella anpassade egenskapsnamnen och `Slug` för den tillgängliga arbetaren för detta projekt.

## Steg för steg

1. Skapa nödvändiga arbetare (**Arbetare**) och egenskaper (**Egenskaper**) innan du importerar.
2. Ladda ner mallen, öppna den i en textredigerare, ersätt exemplen med dina egna uppgifter och spara filen med tillägget `.json`.
3. Välj en fil på fliken **Importera uppgifter**. Gränssnittet kommer att visa antalet uppgifter.
4. Klicka på **Importera uppgifter**. Servern kommer att kontrollera hela filen: om ett fel upptäcks kommer inte en enda uppgift från den att skrivas. Korrigera felmeddelandet och försök igen.

Minsta fil:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Exempel med arbetare, egendom och beroende:

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

`version` ska vara `1`, array `tasks` - från 1 till 1000 uppgifter. `title` krävs. Giltiga stadier: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioriteter: `critical`, `high`, `medium`, `low`; testlägen: `ai`, `human`, `hybrid`. Ansvarig: `unassigned`, `human` eller `worker` med befintlig `worker` (Slug eller ID). `properties` använder namn eller ID för redan skapade egenskaper. `key` är unik i filen; `depends_on` hänvisar till sådana nycklar. Servern skapar själv uppgifts-ID:n.

Att importera samma fil igen kommer att skapa nya problem, så kontrollera tavlan innan du klickar igen. När de importeras registreras uppgifter som skapade av människor; En systempost visas i historiken.
