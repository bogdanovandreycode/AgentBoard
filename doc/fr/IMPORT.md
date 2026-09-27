# Importer des tâches depuis JSON

L'onglet **Importer des tâches** accepte un fichier JSON au format UTF-8. Utilisez-le si vous transférez des tâches depuis un autre service. Cliquez sur **Télécharger le modèle JSON** : Le modèle inclut les noms de propriétés personnalisées actuelles et `Slug` du travailleur disponible pour ce projet.

## Pas à pas

1. Créez les travailleurs (**Workers**) et les propriétés (**Properties**) nécessaires avant l'importation.
2. Téléchargez le modèle, ouvrez-le dans un éditeur de texte, remplacez les exemples par vos propres tâches et enregistrez le fichier avec l'extension `.json`.
3. Dans l'onglet **Importer des tâches**, sélectionnez un fichier. L'interface affichera le nombre de tâches.
4. Cliquez sur **Importer des tâches**. Le serveur vérifiera l'intégralité du fichier : si une erreur est détectée, aucune tâche ne sera écrite. Corrigez le message d'erreur et réessayez.

Fichier minimum :

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Exemple avec travailleur, propriété et dépendance :

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

`version` doit être `1`, tableau `tasks` - de 1 à 1 000 tâches. `title` est requis. Étapes valides : `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Priorités : `critical`, `high`, `medium`, `low` ; modes de test : `ai`, `human`, `hybrid`. Responsable : `unassigned`, `human` ou `worker` avec `worker` existant (Slug ou ID). `properties` utilise les noms ou ID des propriétés déjà créées. `key` est unique dans le fichier ; `depends_on` fait référence à ces clés. Le serveur crée lui-même les ID de tâches.

Importer à nouveau le même fichier créera de nouveaux problèmes, alors vérifiez le tableau avant de cliquer à nouveau. Lors de l'importation, les tâches sont enregistrées comme créées par l'homme ; Une entrée système apparaît dans l'historique.
