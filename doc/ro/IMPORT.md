# Importați sarcini din JSON

Fila **Import tasks** acceptă un fișier JSON în codificare UTF-8. Utilizați-l dacă transferați sarcini de la alt serviciu. Faceți clic pe **Descărcați șablonul JSON**: șablonul include numele proprietăților personalizate curente și `Slug` ale lucrătorului disponibil pentru acest proiect.

## Pas cu pas

1. Creați lucrătorii necesari (**Lucrători**) și proprietățile (**Proprietăți**) înainte de a importa.
2. Descărcați șablonul, deschideți-l într-un editor de text, înlocuiți exemplele cu propriile sarcini și salvați fișierul cu extensia `.json`.
3. În fila **Importați sarcini**, selectați un fișier. Interfața va afișa numărul de sarcini.
4. Faceți clic pe **Importați sarcini**. Serverul va verifica întregul fișier: dacă este detectată o eroare, nu va fi scrisă nicio sarcină din aceasta. Corectați mesajul de eroare și încercați din nou.

Fișier minim:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Exemplu cu lucrător, proprietate și dependență:

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

`version` ar trebui să fie `1`, matrice `tasks` - de la 1 la 1000 de sarcini. `title` este necesar. Etape valide: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Priorități: `critical`, `high`, `medium`, `low`; moduri de testare: `ai`, `human`, `hybrid`. Responsabil: `unassigned`, `human` sau `worker` cu `worker` existent (Slug sau ID). `properties` utilizează numele sau ID-urile proprietăților deja create. `key` este unic în fișier; `depends_on` se referă la astfel de chei. Serverul creează el însuși ID-uri de activitate.

Importarea aceluiași fișier din nou va crea probleme noi, așa că verificați tabloul înainte de a da din nou clic. Când sunt importate, sarcinile sunt înregistrate ca fiind create de om; O intrare de sistem apare în istoric.
