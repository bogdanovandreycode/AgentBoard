# Gestion de projet et rédaction de tâches

La page **Tous les projets** affiche le dossier enregistré, le nombre de projets publiés
tâches et dernière activité. L'activité prend en compte les changements dans les tâches, l'historique et
dossiers de dépenses. Le bouton crayon vous permet de renommer le projet ou de sélectionner
dossier existant dans la boîte de dialogue système. La modification du chemin met à jour l'enregistrement et le marqueur ;
Les ID de projet, les travailleurs, les tâches et les paramètres sont enregistrés, les fichiers ne sont pas déplacés.
Les dossiers déjà enregistrés dans un autre projet ne peuvent pas être sélectionnés.

## Paramètres globaux

Ouvrez les **Paramètres globaux** à partir de la page ou de l'onglet de sélection de projet.
**Paramètres globaux** dans les paramètres de n'importe quelle carte. Ces paramètres sont communs à tous les projets
une base AgentBoard :

- L'enregistrement automatique des brouillons par minuterie est activé par défaut, toutes les **10 secondes**.
- **La sauvegarde en quittant le terrain** est activée par défaut et peut être configurée séparément.
- Les vérifications **Contexte du projet MCP** sont activées et désactivées globalement. Leur région
  les actions et un indicateur CLI distinct sont décrits dans la section [Workers and MCP](WORKERS_MCP.md).

## Cycle de vie du brouillon

La boîte de dialogue Nouvelle tâche enregistre le travail inachevé sous forme de brouillon. Ligne d'état
fait la distinction entre les modifications non enregistrées, l'enregistrement en cours, le brouillon enregistré et l'erreur
conservation. Le brouillon n'apparaît pas sur le tableau et n'est pas disponible pour les travailleurs de MCP jusqu'à ce que
publié.

Avant de publier, utilisez le bouton **Créer une tâche** pour renseigner **Titre, État,
Responsable et Description**. Avant cela, **Enregistrer le brouillon** enregistre un brouillon sans créer
tâches. Un formulaire vide est supprimé par défaut une fois fermé. Si effacé
Le brouillon existant est laissé vide, il est supprimé lors de son prochain enregistrement.
La fermeture d'une boîte de dialogue non vide l'enregistre également.

Le bouton **Brouillons** sur le tableau vous permet de continuer ou de supprimer les brouillons enregistrés.
Plusieurs brouillons pris en charge. Après un redémarrage, l'ouverture
brouillon. Les brouillons sont stockés dans SQLite ; A chaque changement, le local est mis à jour
copier dans le navigateur, afin que les modifications non enregistrées survivent aux redémarrages ou
panne temporaire du serveur. Si le stockage du navigateur n'est pas disponible, une boîte de dialogue vous avertit
que le brouillon doit être sauvegardé sur le serveur avant la fermeture.

La rédaction de brouillons utilise des révisions pour empêcher une fenêtre d'écraser le travail
un autre. La réponse de sauvegarde perdue est cohérente s'il y en a déjà une sur le serveur
données identiques. Si le contenu est différent, un conflit est affiché et
**Récupérer en tant que nouveau brouillon** enregistre le contenu local séparément. Publication
atomique et idempotent : republiez le même brouillon après sa perte
La réponse renvoie la même tâche plutôt que de créer une copie.
