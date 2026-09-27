# Installation et lancement sous Windows

## Installateur (recommandé)

Téléchargez `agentboard-VERSION-windows-amd64-setup.exe` à partir de la page [Sorties](https://github.com/bogdanovandreycode/AgentBoard/releases). L'assistant d'installation proposera un dossier ; la valeur par défaut est `C:\AI\AgentBoard`. Il copiera `agentboard.exe` et la documentation, créera un raccourci et ajoutera le dossier sélectionné au système `PATH`. Après l'installation, ouvrez un nouveau terminal pour que la commande `agentboard` devienne disponible.

Dans le dossier de votre projet, exécutez :

```powershell
agentboard init
agentboard open
```

L'interface est intégrée au `agentboard.exe` ; Aucune installation séparée de Go ou Node.js n'est nécessaire. Les données sont stockées dans `%AppData%\AgentBoard` et sont conservées lorsque le programme est mis à jour ou désinstallé. La désinstallation via « Applications installées » supprime les raccourcis et l'entrée de `PATH`.

## Installation de Scoop

Si Scoop n'est pas déjà installé, ouvrez PowerShell en tant qu'utilisateur régulier et suivez les [instructions officielles de Scoop](https://scoop.sh/). S'il existe des restrictions sur votre ordinateur d'entreprise, contactez votre administrateur ; AgentBoard peut également être lancé depuis un ZIP sans Scoop.

Vous pouvez également installer AgentBoard via Scoop :

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

La disponibilité du manifeste dans GitHub Release peut être vérifiée sur [page Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Après avoir ajouté le manifeste à Scoop, le bucket peut être installé par nom de bucket et mis à jour avec la commande `scoop update agentboard`.

## ZIP sans Scoop

Téléchargez `agentboard-VERSION-windows-amd64.zip` à partir des versions, décompressez, par exemple, dans `C:\Tools\AgentBoard`. En PowerShell dans le dossier du projet :

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Pour un client AI, spécifiez le chemin complet vers `agentboard.exe` dans sa configuration MCP si le programme ne se trouve pas dans `PATH`.

## Construire à partir des sources

Installez la version Go à partir de `go.mod` et Node.js 22 ou version ultérieure. En PowerShell à la racine du dépôt :

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` assemble le frontend en `internal/webui/dist`, exécute des tests Go et assemble un `agentboard.exe`. L'ordre est important : l'interface est intégrée au binaire lors de la construction de Go. Si l'ancien `agentboard.exe open` est en cours d'exécution, arrêtez-le avant de le reconstruire (Ctrl+C), sinon Windows ne vous permettra pas de remplacer le fichier.

## Où sont les données

- `%AppData%\AgentBoard\agentboard.db` - tâches, projets, travailleurs et paramètres. Vous pouvez spécifier un fichier différent avec l'indicateur `--db`, mais `init`, `open`/`serve` et `mcp` doivent avoir le **même chemin**.
- `<ваш проект>\.agentboard\project.json` — identifiant du projet. Ce fichier ne contient pas de tâches.
- Le serveur n'écoute que `127.0.0.1:7337` par défaut. Spécifiez une autre adresse `--addr` avant le chemin du projet : `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` ouvre la page du projet et réutilise un serveur déjà en cours d'exécution à cette adresse. Si plusieurs projets sont ouverts dans un navigateur, sélectionnez-les dans la liste de gauche.
