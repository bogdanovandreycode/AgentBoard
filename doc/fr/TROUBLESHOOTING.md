# Résolution de problèmes

| Symptôme | Que vérifier |
| --- | --- |
| `agentboard` introuvable | Redémarrez PowerShell après Scoop. Lors de l'installation à partir d'un ZIP, utilisez le chemin complet vers `agentboard.exe`. |
| Page Web montrant l'ancienne interface après la construction | Arrêtez le serveur en cours d'exécution Ctrl+C. Exécutez `./scripts/build.ps1`, lancez un nouveau binaire. Vite doit être construit avant Go car l'interface est intégrée à l'EXE. Actualisez la page Ctrl+F5. |
| Port 7337 occupé | AgentBoard est peut-être déjà en cours d'exécution. Ouvrez `http://127.0.0.1:7337` ou terminez l'ancien processus. Pour un autre port, utilisez `--addr`. |
| Projet introuvable | Dans le dossier souhaité, exécutez `agentboard init`. Puis `agentboard open` ou `agentboard open C:\путь\к\проекту`. |
| Le travailleur ne voit pas la tâche | La tâche doit être affectée à ce travailleur particulier et se trouver dans `Features`, `In progress`, `Testing` ou `Verification`. AI ne voit pas `Backlog`, `Complete` ni les colonnes personnalisées. |
| `NO_PROJECTS` | Il n'y a aucun projet dans la base de données MCP. Exécutez `agentboard init` dans le dossier du projet et créez un travailleur ; l'agent peut répéter l'appel sans redémarrer. Vérifiez que MCP et l'interface Web utilisent le même `--db`. |
| `PROJECT_REQUIRED` / `PROJECT_AMBIGUOUS` | Sans `--project`, ni le dossier de travail du client ni le dossier de démarrage MCP ne se trouvent exactement dans le même projet enregistré. Exécutez l'agent dans le dossier du projet, enregistrez-le via `agentboard init` ou ajoutez `--project`. Il n'est pas nécessaire de lancer l'interface Web. |
| `PROJECT_NOT_FOUND` / `WORKER_NOT_FOUND` / `WORKER_REQUIRED` / `NO_WORKERS` | Découvrez `--project`, le slug de travail dans ce projet et `--db`. Le message répertorie les projets enregistrés ou les travailleurs activés. Une configuration MCP personnalisée sans `--project` nécessite un travailleur avec le même slug dans chaque projet. |
| `PROJECT_CONTEXT_MISMATCH` | Le dossier de démarrage MCP ou le dossier de travail client ne correspond pas au `--project` spécifié. Veuillez confirmer la question si cela est intentionnel, sinon supprimez `--project` (détection automatique) ou indiquez la bonne. |
| La vérification du serveur MCP existe, mais le client n'est pas connecté | La vérification du serveur ne vérifie pas les paramètres du client externe. Redémarrez le client, vérifiez son fichier de configuration, le chemin d'accès à `agentboard.exe`, `--project`, `--worker` et `--db` général. Demandez à appeler `get_my_board`. |
| Travailleur hors ligne | Le client a peut-être terminé ou n'a peut-être pas encore démarré MCP. Après 90 secondes sans battement de cœur, la session est considérée comme déconnectée. |
| Le fichier JSON n'est pas importé | Vérifiez `version: 1`, `title` requis, les travailleurs `Slug` existants et les noms de propriété. JSON n'autorise pas les commentaires ni les virgules finales. |
| Impossible de reconstruire `agentboard.exe` | Windows ne peut pas remplacer un EXE en cours d'exécution. Arrêtez le serveur Ctrl+C et réessayez la compilation. |
| Les tâches ont disparu après la mise à jour | Vérifiez que `--db` ne pointe pas vers un autre fichier et que vous êtes connecté en tant que même utilisateur Windows. La base par défaut se trouve dans `%AppData%\AgentBoard`. |

Si l'erreur n'est pas décrite, collectez le texte exact du message, les versions `agentboard version` et Windows, ainsi que les étapes à suivre pour réessayer. Ne publiez pas de données de projet privé ou de contenu de base de données dans un numéro ouvert.
