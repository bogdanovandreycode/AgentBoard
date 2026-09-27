# Commencez ici

## Qu'est-ce qu'AgentBoard

Imaginez un tableau ordinaire avec des cartes de tâches. Vous créez des tâches, affectez une personne responsable et supervisez le travail. Les travailleurs de l'IA reçoivent uniquement les tâches qui leur sont assignées via MCP et signalent les progrès. Vous décidez quand la tâche est enfin prête.

**Projet** - un dossier sur l'ordinateur et une carte séparée. **Tâche** - une carte avec une description, une personne responsable et une étape. **Worker** est le nom logique du client IA. **MCP** est la manière dont le client AI se connecte à AgentBoard. **Scoop** est un gestionnaire d'installation de logiciels pour Windows.

Aucun code requis. Vous avez besoin de Windows, d'un navigateur, de PowerShell et, pour que l'IA fonctionne, d'un client IA installé prenant en charge les serveurs MCP locaux.

## Premier lancement

1. Installez l'application conformément aux [instructions](INSTALL.md).
2. Créez un dossier de projet dans l'Explorateur, par exemple `C:\Projects\MyFirstProject`.
3. Ouvrez ce dossier dans l'Explorateur. Cliquez dans la barre d'adresse, tapez `powershell` et appuyez sur Entrée.
4. Dans la fenêtre qui s'ouvre, faites :

```powershell
agentboard init
agentboard open
```

5. Un navigateur s'ouvrira avec l'adresse `http://127.0.0.1:7337`. Laissez la fenêtre PowerShell ouverte pendant que vous utilisez le tableau blanc. La fermeture de la fenêtre arrêtera le serveur local, mais les tâches resteront.

Si la commande `agentboard` est introuvable, fermez PowerShell et rouvrez-le après avoir installé Scoop. Lors de l'installation à partir d'un ZIP, utilisez le chemin complet vers `agentboard.exe`.

## Première tâche

Cliquez sur **Nouvelle tâche**, remplissez **Titre**, si nécessaire **Description**, puis **Enregistrer**. Une nouvelle tâche dans `Backlog` est disponible pour les humains. Pour que l'IA commence à travailler, affectez un travailleur et transférez la tâche vers `Features`. Description étape par étape des champs - [TASKS.md](TASKS.md).

## Premier travailleur

Ouvrez **Travailleurs → Ajouter un travailleur**. Sélectionnez le client AI que vous utilisez (par exemple Codex ou Claude Code), vérifiez le nom et l'ID court `Slug`, cliquez sur **Enregistrer**. Ouvrez le travailleur créé : il y a une configuration MCP, un bouton de copie et une vérification du serveur. Copiez la configuration sur le client AI selon [WORKERS_MCP.md](WORKERS_MCP.md). Une fois connecté, demandez au client d'appeler `get_my_board`.

## Que lire ensuite

- [Tâches, tests, histoires et colonnes](TASKS.md)
- [Connexion des travailleurs et vérification de MCP](WORKERS_MCP.md)
- [Importer des tâches depuis JSON](IMPORT.md)
- [Langue, thème, fuseau horaire et mise à jour](SETTINGS.md)
- [Problèmes typiques](TROUBLESHOOTING.md)
