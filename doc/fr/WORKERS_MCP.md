# Connexion des travailleurs et du MCP

## Étape 1. Créer un travailleur

Dans AgentBoard, ouvrez **Travailleurs → Ajouter un travailleur**. Sélectionnez un profil client : Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot ou un autre client MCP. Le profil remplit des fonctionnalités typiques (code, tests, Git) ; vous pouvez les changer. Le nom est visible sur la carte et `Slug` est un identifiant court sans espaces pour la commande MCP. Cliquez sur **Enregistrer**.

Un travailleur correspond à une personnalité de l’IA. Créez différents travailleurs pour différents clients ou équipes. Le profil de capacités décrit la spécialisation, mais n'étend pas les droits de l'IA aux étapes de la tâche.

## Étape 2. Copiez la configuration

Ouvrez le travailleur créé. Le bloc **MCP diagnostics** affiche le fragment de configuration et le fichier où l'ajouter. Cliquez sur **Copier la configuration MCP**. Si le fichier existe déjà, ajoutez le serveur suggéré à l'objet `mcpServers`/`servers`/`mcp` existant sans effacer les autres serveurs.

La commande principale ressemble à ceci :

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` du travailleur créé. Le client MCP exécute lui-même cette commande lorsqu'il a besoin d'outils.

### MCP fonctionne sans interface Web

`agentboard mcp` est un processus stdio distinct qui ouvre lui-même la base de données AgentBoard. Il n'a pas besoin de `agentboard open`/`serve`, cela ne dépend pas du dossier à partir duquel l'interface Web est lancée ni du nombre d'instances de l'interface Web en cours d'exécution. Plusieurs agents peuvent s'exécuter simultanément, chacun dans son propre projet : chaque client exécute son propre processus MCP et ils utilisent tous en toute sécurité la même base de données SQLite.

### Comment un projet est sélectionné

Le serveur MCP s'exécute dans n'importe quel dossier. La première fois que vous appelez l'outil, il est lié au projet :

1. `--project PATH`, si spécifié (n'importe quel dossier du projet enregistré convient) ;
2. sinon les dossiers de travail (racines) du client, si le client les signale ;
3. sinon le dossier dans lequel le client a lancé le serveur MCP.

`get_my_board` renvoie le projet et le collaborateur sélectionnés à `session` avec `binding_source`. Si les dossiers de travail du client changent, le prochain appel à l'outil effectue à nouveau la liaison. Sans `--worker`, le seul travailleur activé du projet est utilisé.

Si la liaison n'est pas possible, chaque outil renvoie une erreur que l'agent doit transmettre à l'utilisateur ; les données de la tâche ne sont pas lues dans ce cas :

| Codes | Signification |
| --- | --- |
| `NO_PROJECTS` | Il n'y a aucun projet dans la base de données. Exécutez `agentboard init` dans le dossier du projet et créez un travailleur. Aucun redémarrage nécessaire. |
| `PROJECT_REQUIRED` | Le dossier de travail se trouve en dehors des projets enregistrés. Le message répertorie les projets. |
| `PROJECT_NOT_FOUND` | `--project` ne se trouve pas dans un projet enregistré. |
| `PROJECT_AMBIGUOUS` | Les dossiers de travail d'un client s'étendent sur plusieurs projets. Ajoutez `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Créez, activez ou sélectionnez un travailleur via `--worker`. |

Les configurations pour **Codex** et **Gemini CLI** sont personnalisées et partagées dans tous les dossiers, elles n'ont donc pas `--project` : créez un travailleur avec le même slug dans chaque projet où ce client est en cours d'exécution. Les fichiers au niveau du projet (Claude Code, Cursor, OpenCode, VS Code) contiennent le `--project` de leur projet.

### Conception explicite et vérification du contexte

Si `--project` est spécifié et que le dossier de démarrage ou le dossier de travail du client se trouve dans un autre projet, les outils demandent une confirmation via l'élicitation MCP. Les clients sans élicitation de formulaire reçoivent `PROJECT_CONTEXT_MISMATCH` avant d'accéder aux tâches. Corrigez la configuration, supprimez `--project` pour la détection automatique ou confirmez explicitement que vous avez l'intention de travailler avec un autre projet avec l'indicateur `--allow-project-mismatch`. La modification des dossiers de travail réinitialise les confirmations précédentes. Les projets définis automatiquement ne nécessitent pas de confirmation.

**Paramètres globaux → Contexte du projet MCP** contrôle ces vérifications pour tous les projets et travailleurs de la même base de données. Activé par défaut. S'il est désactivé, MCP utilise le projet spécifié sans confirmation ni erreur de non-concordance. Les sessions déjà ouvertes appliqueront le paramètre la prochaine fois que l'outil sera appelé ; la visibilité des tâches, les transitions et les règles de tests humains continuent d'être suivies.

Après vérification, le bloc de diagnostic affiche le chemin absolu vers le fichier en cours d'exécution.`agentboard.exe`. Ceci est particulièrement utile lors de l’installation à partir d’un ZIP. Lors de l'installation via Scoop, vous pouvez utiliser la commande`agentboard`, si le client voit la même chose`PATH`.

## Étape 3 : Ajoutez un serveur à votre client

L'interface de travail contient déjà un fragment prêt à l'emploi. Vous trouverez ci-dessous une explication de l'endroit où il est utilisé :

| Client | Où insérer | Comment vérifier côté client |
| --- | --- | --- |
| Codex |`%USERPROFILE%\.codex\config.toml`, section`[mcp_servers.agentboard_<slug>]`|`codex mcp list`|
| Claude Code |`.mcp.json`dans le dossier du projet |`claude mcp list`|
| CLI Gémeaux |`%USERPROFILE%\.gemini\settings.json`, objet`mcpServers`|`/mcp list`dans Gemini CLI |
| Curseur |`.cursor\mcp.json`projet | liste des serveurs MCP dans les paramètres du curseur |
| Code Ouvert |`opencode.json`projet, objet`mcp`| liste des outils MCP dans OpenCode |
| Copilote VS Code |`.vscode\mcp.json`projet, objet`servers`| commande **MCP : Liste des serveurs** |

Exemples pour Codex (configuration personnalisée, projet détecté automatiquement) et Claude Code (fichier projet) avec projet`C:\Projects\MyFirstProject`et travailleur`codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

En JSON, la barre oblique inverse de Windows est doublée ; un fragment prêt à l'emploi de l'interface le fait automatiquement. Les configurations générées contiennent le `--db` réel du serveur web, donc HTTP et MCP utilisent la même base, même si l'environnement client est différent.

**Ollama** fournit un modèle local, mais ne remplace pas le client MCP. Le profil **Ollama via OpenCode** génère une configuration MCP pour OpenCode ; configurez séparément OpenCode sur le modèle Ollama. Un autre client compatible MCP avec Ollama convient également.

## Étape 4 : Vérifiez votre connexion

1. Ouvrez la carte de travailleur et cliquez sur **Vérifier à nouveau**. **Server check** exécute un véritable exécutable stdio avec ce projet, ce travailleur et cette base, vérifie la poignée de main, la liste des outils et que le projet et le travailleur sont localisés.
2. Démarrez ou redémarrez le client AI après avoir ajouté le fichier de configuration. Demandez-lui d'appeler `get_my_board`.
3. **Client connecté** apparaîtra dans AgentBoard et une nouvelle session apparaîtra dans la liste. Seulement cela confirme la connexion de votre client. S'il existe des outils, mais que le client n'est pas connecté, vérifiez le chemin d'accès au programme, le nom du fichier de configuration et sa syntaxe JSON/TOML.

Démarrez votre session de travail avec `get_my_board`. L'IA ne reçoit pas de tâches de `Backlog` et `Complete`, même si elle connaît leur identifiant. L’IA ne peut pas prétendre être humaine et n’a pas de commande générale « se déplacer n’importe où ». L'utilisation du système de fichiers en dehors d'AgentBoard dépend des capacités et du bac à sable du client sélectionné.

Chaque outil déclare `title` et les quatre conseils de comportement MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) sous forme de valeurs booléennes explicites. Les outils de lecture sont en lecture seule et idempotents ; les outils d'enregistrement ajoutent uniquement des enregistrements ou font avancer les tâches, de sorte qu'aucun n'est perturbateur ou ne dépasse la base AgentBoard locale.

Instructions officielles du client : [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
