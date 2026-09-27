# AgentBoard

[🌐 Languages](../LANGUAGES.md)

AgentBoard est un tableau de tâches local où les travailleurs humains et IA travaillent sur des tâches communes, mais ont des droits différents. L'application est lancée avec un fichier `agentboard.exe`, ouvre l'interface Web dans le navigateur et fournit aux travailleurs un serveur MCP distinct via `stdio`. Les données restent sur votre ordinateur.

**[Partir de zéro](START_HERE.md) · [Travailler avec des tâches](TASKS.md) · [Connexion de l'IA via MCP](WORKERS_MCP.md) · [Importer JSON](IMPORT.md) · [Paramètres](SETTINGS.md) · [Résoudre des problèmes](TROUBLESHOOTING.md)**

## Dans cinq minutes

1. Téléchargez le programme d'installation de `agentboard-VERSION-windows-amd64-setup.exe` à partir de [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Il proposera un dossier (par défaut `C:\AI\AgentBoard`) et l'ajoutera à `PATH`. Scoop et ZIP sont également disponibles.
2. Ouvrez PowerShell dans votre dossier de projet, par exemple `C:\Projects\MyApp`.
3. Exécutez `agentboard init` (pour ZIP : chemin complet vers `agentboard.exe` et `init`).
4. Exécutez `agentboard open`. `http://127.0.0.1:7337` s'ouvrira.
5. Ajoutez une tâche à l'aide du bouton **Nouvelle tâche**. Pour un travailleur IA, ouvrez **Travailleurs → Ajouter un travailleur**, sélectionnez le profil client et copiez la configuration MCP.

Si vous n'avez pas encore de dossier de projet, créez-en un dans l'Explorateur Windows. Un projet peut être n'importe quel dossier, même sans Git ni code.

## Installation via Scoop

Dans PowerShell avec [Scoop](https://scoop.sh/)] déjà installé après la version :

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Pour les développeurs, il existe [build from source ](INSTALL.md). Le workflow de publication crée un programme d'installation, un manifeste ZIP et Scoop avec SHA-256 à partir du même artefact. Mise à jour de la version installée via Scoop : `scoop update agentboard` après ajout du manifeste au bucket ; détails - [préparation de la version](SCOOP_RELEASE.md).

## Comment le conseil d'administration est structuré

`Backlog → Features → In progress → Testing → Verification → Complete`

Une personne peut déplacer des tâches sur le tableau. L'IA ne peut déplacer que `Features → In progress → Testing → Verification` ; l'acceptation finale dans `Complete` est effectuée par un humain. Les colonnes utilisateur sont destinées aux humains : la tâche qu'elles contiennent reste dans l'état `Backlog` pour MCP. Modes de test : IA, Humain et Hybride. L’historique, les tests, les artefacts et les coûts de l’IA sont attachés à la tâche.

## Équipes

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` enregistre le dossier et y écrit uniquement `.agentboard/project.json`. Les données de production SQLite se trouvent dans le répertoire de configuration utilisateur Windows (`%AppData%\AgentBoard\agentboard.db`), en dehors du projet et en dehors de l'installation de Scoop. La suppression ou la mise à jour du package ne devrait pas supprimer ces données. Avant de transférer vers un autre ordinateur, faites une copie de la base de données pendant qu'AgentBoard est arrêté.

Travailleurs - comptes logiques ; AgentBoard lui-même n'exécute pas Codex, Claude ou tout autre client IA. Le client démarre un processus MCP local pour un travailleur spécifique. MCP est la limite d'autorisation de l'application, et pour l'isolation des fichiers, utilisez le bac à sable client AI.

## Pour les développeurs

Pile : Go, SQLite, SDK MCP Go officiel, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Créez d'abord l'interface, puis allez : `./scripts/build.ps1`. Les fichiers Web sont inclus dans le binaire via `go:embed`. L'architecture et l'API sont décrites dans [doc/ARCHITECTURE.md](ARCHITECTURE.md).
