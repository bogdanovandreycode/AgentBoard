# Architecture et API

AgentBoard est un processus Go local exécutant SQLite. L'API HTTP pour les humains, l'adaptateur MCP pour l'IA et l'interface Web partagent une logique de service commune. SQLite - source d'état ; le frontend ne contourne pas le serveur. Les outils d'IA ont une surface de droits distincte et le statut `Backlog`/`Complete` dans MCP n'est pas disponible. L'enregistrement historique de `System` est créé par l'application elle-même.

## Composants

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - modèles de domaine et erreurs.
- `internal/service` - transitions de tâches, autorisations, importation et paramètres.
- `internal/persistence` - SQLite et migrations.
- `internal/httpapi` - API HTTP humaine.
- `internal/mcpserver` - Outils MCP pour un travailleur distinct.
- `web` - Interface utilisateur React/TypeScript ; l'assembly se retrouve dans `internal/webui/dist` et est inclus dans l'EXE.

## Routes HTTP de base

| Méthode et chemin | Destination |
| --- | --- |
| `GET /api/health` | Vérification du serveur. |
| `GET /api/projects` | Projets enregistrés. |
| `GET /api/projects/{id}/board` | Tableau de projet. |
| `GET/PUT /api/projects/{id}/settings` | Paramètres, ordre et noms des colonnes. |
| `POST /api/projects/{id}/tasks` | Créez une tâche. |
| `POST /api/projects/{id}/tasks/import` | Importer atomiquement JSON version 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Carte, modifier, supprimer. |
| `POST /api/tasks/{id}/move` | Déménagement par une personne. |
| `GET/POST /api/projects/{id}/workers` | Liste et création d'un travailleur. |
| `GET /api/workers/{id}/mcp/check` | Tests internes de la poignée de main et des outils MCP. |
| `GET /api/workers/{id}/sessions` | Diagnostic de session. |
| `GET/POST /api/projects/{id}/properties` | Propriétés personnalisées. |

HTTP est destiné à l'utilisateur de confiance local. Ne publiez pas un port Web sur Internet sans sa propre authentification, ses restrictions réseau et HTTPS. Le serveur MCP est lancé à l'aide de `stdio` pour un projet et un travailleur spécifiques ; commencez par `get_my_board`. Ses autorisations restreignent les actions au sein d'AgentBoard, mais ne remplacent pas le bac à sable du système de fichiers du client AI.

Les colonnes utilisateur sont stockées séparément de `tasks.state` : le noyau conserve la tâche dans une telle colonne dans `backlog`, et `tasks.board_column` détermine la place sur le tableau humain. Cela maintient le même modèle de transition d’IA. La suppression d'une colonne efface les tâches `board_column`.
