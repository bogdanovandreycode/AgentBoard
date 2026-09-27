# Connexion des travailleurs et du MCP

## Étape 1. Créer un travailleur

Dans AgentBoard, ouvrez **Travailleurs → Ajouter un travailleur**. Sélectionnez un profil client : Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot ou un autre client MCP. Le profil remplit des fonctionnalités typiques (code, tests, Git) ; vous pouvez les changer. Le nom est visible sur la carte et `Slug` est un identifiant court sans espaces pour la commande MCP. Cliquez sur **Enregistrer**.

Un travailleur correspond à une personnalité de l’IA. Créez différents travailleurs pour différents clients ou équipes. Le profil de capacités décrit la spécialisation, mais n'étend pas les droits de l'IA aux étapes de la tâche.

## Étape 2. Copiez la configuration

Ouvrez le travailleur créé. Le bloc **MCP diagnostics** affiche le fragment de configuration et le fichier où l'ajouter. Cliquez sur **Copier la configuration MCP**. Si le fichier existe déjà, ajoutez le serveur suggéré à l'objet `mcpServers`/`servers`/`mcp` existant sans effacer les autres serveurs.

La commande principale ressemble à ceci :

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` doit pointer vers le dossier que vous avez enregistré via `agentboard init`. `--worker` — `Slug` du travailleur créé. Le client MCP exécute lui-même cette commande lorsqu'il a besoin d'outils. Dans le navigateur AgentBoard, le serveur Web peut s'exécuter séparément.

Après vérification, le bloc de diagnostic affiche le chemin absolu vers le `agentboard.exe` en cours d'exécution. Ceci est particulièrement utile lors de l’installation à partir d’un ZIP. Lors de l'installation via Scoop, vous pouvez utiliser la commande `agentboard` si le client voit le même `PATH`.

## Étape 3 : Ajoutez un serveur à votre client

L'interface de travail contient déjà un fragment prêt à l'emploi. Vous trouverez ci-dessous une explication de l'endroit où il est utilisé :

| Client | Où insérer | Comment vérifier côté client |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, coupe `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` dans le dossier du projet | `claude mcp list` |
| CLI Gémeaux | `%USERPROFILE%\.gemini\settings.json`, objet `mcpServers` | `/mcp list` dans la CLI Gemini |
| Curseur | Projet `.cursor\mcp.json` | liste des serveurs MCP dans les paramètres du curseur |
| Code Ouvert | Projet `opencode.json`, objet `mcp` | liste des outils MCP dans OpenCode |
| Copilote VS Code | Projet `.vscode\mcp.json`, objet `servers` | commande **MCP : Liste des serveurs** |

Pour Codex et Claude Code, un exemple avec le projet `C:\Projects\MyFirstProject` et le travailleur `codex` :

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

En JSON, la barre oblique inverse de Windows est doublée ; un fragment prêt à l'emploi de l'interface le fait automatiquement. Si vous utilisez `--db` avec une base personnalisée, ajoutez-la à la configuration `args` MCP et spécifiez le même chemin que lors du démarrage de `open`.

**Ollama** fournit un modèle local, mais ne remplace pas le client MCP. Le profil **Ollama via OpenCode** génère une configuration MCP pour OpenCode ; configurez séparément OpenCode sur le modèle Ollama. Un autre client compatible MCP avec Ollama convient également.

## Étape 4 : Vérifiez votre connexion

1. Ouvrez la carte de travailleur et cliquez sur **Vérifier à nouveau**. **La vérification du serveur** devrait indiquer le nombre d'outils MCP. Il s'agit d'une vérification du protocole interne et d'une détection des outils du serveur.
2. Démarrez ou redémarrez le client AI après avoir ajouté le fichier de configuration. Demandez-lui d'appeler `get_my_board`.
3. **Client connecté** apparaîtra dans AgentBoard et une nouvelle session apparaîtra dans la liste. Seulement cela confirme la connexion de votre client. S'il existe des outils, mais que le client n'est pas connecté, vérifiez le chemin d'accès au programme, le nom du fichier de configuration et sa syntaxe JSON/TOML.

Démarrez votre session de travail avec `get_my_board`. L'IA ne reçoit pas de tâches de `Backlog` et `Complete`, même si elle connaît leur identifiant. L’IA ne peut pas prétendre être humaine et n’a pas de commande générale « se déplacer n’importe où ». L'utilisation du système de fichiers en dehors d'AgentBoard dépend des capacités et du bac à sable du client sélectionné.

Instructions officielles du client : [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
