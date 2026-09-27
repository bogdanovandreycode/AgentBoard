# Aufgaben aus JSON importieren

Die Registerkarte **Importaufgaben** akzeptiert eine JSON-Datei in UTF-8-Kodierung. Verwenden Sie es, wenn Sie Aufgaben von einem anderen Dienst übertragen. Klicken Sie auf **JSON-Vorlage herunterladen**: Die Vorlage enthält die aktuellen benutzerdefinierten Eigenschaftsnamen und `Slug` des verfügbaren Workers für dieses Projekt.

## Schritt für Schritt

1. Erstellen Sie vor dem Import die erforderlichen Worker (**Workers**) und Eigenschaften (**Properties**).
2. Laden Sie die Vorlage herunter, öffnen Sie sie in einem Texteditor, ersetzen Sie die Beispiele durch Ihre eigenen Aufgaben und speichern Sie die Datei mit der Erweiterung `.json`.
3. Wählen Sie auf der Registerkarte **Importaufgaben** eine Datei aus. Die Benutzeroberfläche zeigt die Anzahl der Aufgaben an.
4. Klicken Sie auf **Aufgaben importieren**. Der Server prüft die gesamte Datei: Wenn ein Fehler festgestellt wird, wird keine einzige Aufgabe daraus geschrieben. Korrigieren Sie die Fehlermeldung und versuchen Sie es erneut.

Mindestdatei:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Beispiel mit Arbeiter, Eigentum und Abhängigkeit:

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

`version` sollte `1` sein, Array `tasks` – von 1 bis 1000 Aufgaben. `title` ist erforderlich. Gültige Stufen: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioritäten: `critical`, `high`, `medium`, `low`; Testmodi: `ai`, `human`, `hybrid`. Verantwortlich: `unassigned`, `human` oder `worker` mit vorhandenem `worker` (Slug oder ID). `properties` verwendet die Namen oder IDs bereits erstellter Eigenschaften. `key` ist innerhalb der Datei eindeutig; `depends_on` bezieht sich auf solche Schlüssel. Der Server erstellt selbst Task-IDs.

Das erneute Importieren derselben Datei führt zu neuen Problemen. Überprüfen Sie daher das Forum, bevor Sie erneut klicken. Beim Import werden Aufgaben als von Menschen erstellt erfasst; Im Verlauf erscheint ein Systemeintrag.
