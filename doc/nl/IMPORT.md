# Importeer taken uit JSON

Het tabblad **Importtaken** accepteert een JSON-bestand in UTF-8-codering. Gebruik het als u taken van een andere dienst overdraagt. Klik op **JSON-sjabloon downloaden**: De sjabloon bevat de huidige namen van aangepaste eigenschappen en `Slug` van de beschikbare werker voor dit project.

## Stap voor stap

1. Maak de benodigde werknemers (**Workers**) en eigenschappen (**Eigenschappen**) aan voordat u importeert.
2. Download het sjabloon, open het in een teksteditor, vervang de voorbeelden door uw eigen taken en sla het bestand op met de extensie `.json`.
3. Selecteer op het tabblad **Importtaken** een bestand. De interface toont het aantal taken.
4. Klik op **Taken importeren**. De server controleert het hele bestand: als er een fout wordt gedetecteerd, wordt er geen enkele taak uit geschreven. Corrigeer de foutmelding en probeer het opnieuw.

Minimaal bestand:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Voorbeeld met werknemer, eigendom en afhankelijkheid:

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

`version` moet `1` zijn, array `tasks` - van 1 tot 1000 taken. `title` is vereist. Geldige fasen: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioriteiten: `critical`, `high`, `medium`, `low`; testmodi: `ai`, `human`, `hybrid`. Verantwoordelijk: `unassigned`, `human` of `worker` met bestaande `worker` (Slug of ID). `properties` gebruikt de namen of ID's van reeds gemaakte eigenschappen. `key` is uniek binnen het bestand; `depends_on` verwijst naar dergelijke sleutels. De server maakt zelf taak-ID's aan.

Als u hetzelfde bestand opnieuw importeert, ontstaan ​​er nieuwe problemen, dus controleer het bord voordat u opnieuw klikt. Wanneer ze worden geïmporteerd, worden taken geregistreerd als door mensen gemaakt; Er verschijnt een systeemvermelding in de geschiedenis.
