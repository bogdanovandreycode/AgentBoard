# Importujte úkoly z JSON

Karta **Úlohy importu** přijímá soubor JSON v kódování UTF-8. Použijte jej, pokud přenášíte úkoly z jiné služby. Klikněte na **Stáhnout šablonu JSON**: Šablona obsahuje aktuální názvy vlastních vlastností a `Slug` dostupného pracovníka pro tento projekt.

## Krok za krokem

1. Před importem vytvořte potřebné pracovníky (**Workers**) a vlastnosti (**Properties**).
2. Stáhněte si šablonu, otevřete ji v textovém editoru, nahraďte příklady vlastními úkoly a uložte soubor s příponou `.json`.
3. Na kartě **Úlohy importu** vyberte soubor. Rozhraní zobrazí počet úkolů.
4. Klikněte na **Importovat úlohy**. Server zkontroluje celý soubor: pokud je zjištěna chyba, nebude z něj zapsána ani jedna úloha. Opravte chybovou zprávu a zkuste to znovu.

Minimální soubor:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Příklad s pracovníkem, majetkem a závislostí:

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

`version` by mělo být `1`, pole `tasks` - od 1 do 1000 úloh. Je vyžadován `title`. Platné stupně: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Priority: `critical`, `high`, `medium`, `low`; testovací režimy: `ai`, `human`, `hybrid`. Zodpovídá: `unassigned`, `human` nebo `worker` se stávajícím `worker` (Slug nebo ID). `properties` používá názvy nebo ID již vytvořených vlastností. `key` je v rámci souboru jedinečný; `depends_on` označuje takové klíče. Server sám vytváří ID úloh.

Opětovný import stejného souboru vytvoří nové problémy, takže před dalším kliknutím zkontrolujte desku. Při importu jsou úkoly zaznamenány jako vytvořené člověkem; V historii se objeví systémový záznam.
