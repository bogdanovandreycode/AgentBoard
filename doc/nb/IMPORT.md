# Importer oppgaver fra JSON

Kategorien **Importer oppgaver** godtar en JSON-fil i UTF-8-koding. Bruk den hvis du overfører oppgaver fra en annen tjeneste. Klikk på **Last ned JSON-mal**: Malen inkluderer gjeldende egendefinerte egenskapsnavn og `Slug` for den tilgjengelige arbeideren for dette prosjektet.

## Trinn for trinn

1. Opprett de nødvendige arbeiderne (**Arbeidere**) og egenskaper (**Egenskaper**) før import.
2. Last ned malen, åpne den i et tekstredigeringsprogram, bytt ut eksemplene med dine egne oppgaver og lagre filen med filtypen `.json`.
3. På fanen **Importer oppgaver** velger du en fil. Grensesnittet vil vise antall oppgaver.
4. Klikk på **Importer oppgaver**. Serveren vil sjekke hele filen: hvis en feil oppdages, vil ikke en eneste oppgave fra den bli skrevet. Rett feilmeldingen og prøv igjen.

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

Eksempel med arbeider, eiendom og avhengighet:

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

`version` skal være `1`, array `tasks` - fra 1 til 1000 oppgaver. `title` kreves. Gyldige stadier: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioriteter: `critical`, `high`, `medium`, `low`; testmoduser: `ai`, `human`, `hybrid`. Ansvarlig: `unassigned`, `human` eller `worker` med eksisterende `worker` (Slug eller ID). `properties` bruker navnene eller IDene til allerede opprettede egenskaper. `key` er unik i filen; `depends_on` refererer til slike nøkler. Serveren lager oppgave-IDer selv.

Å importere den samme filen på nytt vil skape nye problemer, så sjekk tavlen før du klikker igjen. Når de importeres, registreres oppgaver som menneskeskapte; En systemoppføring vises i loggen.
