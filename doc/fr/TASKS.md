# Travailler avec des tâches

## Ajouter une tâche

Dans l'onglet **Tableau**, cliquez sur **Nouvelle tâche**. Remplissez le titre et la description. La description et les instructions de test prennent en charge Markdown : mettez en surbrillance le texte et utilisez la barre de formatage pour le gras, l'italique, les liens et la liste. Cliquez sur **Enregistrer**.

Champs :

| Champ | Qu'est-ce que |
| --- | --- |
| Titre | Nom court de la tâche. |
| Descriptif | Que faut-il faire et comment comprendre que le travail est prêt. |
| État | Étape de travail ; une nouvelle tâche commence généralement à `Backlog`. |
| Priorité | `critical`, `high`, `medium` ou `low`. |
| Responsable | Une personne, un travailleur spécifique ou sans rendez-vous. |
| Mode test | `AI` - vérifie l'IA ; `Human` - contrôles humains ; `Hybrid` - les deux. |
| Dépendances | Tâches qui devraient être terminées plus tôt. |
| Instructions de test IA/humain | Instructions à l'intention de l'examinateur approprié. |
| Propriétés personnalisées | Champs supplémentaires créés dans l'onglet **Propriétés**. |

Pour une tâche IA, créez d'abord un travailleur, sélectionnez-le dans Responsable, puis transférez la carte vers `Features`. L'IA ne voit que les tâches qui lui sont assignées en quatre étapes, de `Features` à `Verification`.

## Déplacer la tâche

Faites glisser la carte entre les colonnes. Vous pouvez basculer entre une vue de toutes les colonnes et des colonnes larges avec défilement horizontal. `Backlog` et `Complete` sont contrôlés par l'homme. L'IA ne peut faire avancer la tâche `Features → In progress → Testing → Verification` que via des outils MCP spéciaux. Une tâche avec des tests manuels inachevés ne devrait pas réussir la vérification de l'IA.

## Carte de tâche

Cliquez sur une carte pour voir la description, le propriétaire, les instructions de test et les onglets pour l'historique, les tests, les artefacts et les coûts de l'IA. **Modifier la tâche** modifie le contenu. Le commentaire de la personne est ajouté à l'histoire globale. Recherchez dans les filtres supérieurs les recherches par titre, description, identifiant et travailleur ; un bouton séparé ouvre une grande recherche.

## Colonnes et propriétés

Dans l'onglet **Paramètres**, vous pouvez modifier l'ordre des colonnes disponibles pour une personne et ajouter la vôtre. Les quatre étapes de l'IA sont fixes et se déroulent dans le même ordre. La colonne utilisateur est un lieu pour les tâches reportées par une personne : pour l'IA, une telle tâche a le statut `Backlog`. Lorsqu'une colonne est supprimée, ses tâches reviennent à la normale `Backlog`.

Dans l'onglet **Propriétés**, vous pouvez ajouter des champs tels que du texte, un numéro, un indicateur, une date, une sélection et une URL. `Human only` cache le champ à l'IA ; `Agent read` permet la lecture, `Agent read/write` permet également l'écriture via les outils pris en charge. Cela ne change pas les droits de l'IA sur les étapes de la tâche.
