# Paramètres

Ouvrez **Paramètres** dans le menu latéral du projet sélectionné. Les modifications sont enregistrées avec le bouton **Enregistrer les paramètres** et stockées pour le projet dans la base de données AgentBoard.

## Langue et apparence

Le champ **Langue** ouvre une liste déroulante de recherche. Par défaut, l'option **Suivre le système** est sélectionnée, qui prend la langue du navigateur. Langues disponibles à partir des captures d'écran : arabe, portugais (Brésil), chinois simplifié, tchèque, danois, néerlandais, anglais, finnois, français, allemand, italien, japonais, coréen, bokmål norvégien, polonais, russe, espagnol, suédois, turc, ukrainien, vietnamien ; ainsi que le biélorusse, le roumain et le bulgare. **Suivez le système** prend la langue du navigateur. Les signatures principales sont traduites manuellement, les lignes restantes de l'interface utilisateur ont une traduction automatique préliminaire. Avant la diffusion publique, il est conseillé de relire les traductions par des locuteurs natifs ; les noms de clients, les commandes, les champs JSON et les données utilisateur restent sans traduction.

**Schéma de couleurs** : Sombre (thème original), Clair, Noir, Ubuntu et Windows. **Fuseau horaire** contrôle l'affichage des dates ; les données continuent d'être stockées en UTC. **Le fuseau horaire du système** utilise les paramètres de l'ordinateur.

## Colonnes

Les quatre étapes `Features`, `In progress`, `Testing`, `Verification` sont fixées dans cet ordre : elles ne peuvent pas être renommées ou supprimées. D'autres colonnes peuvent être réorganisées à l'aide des flèches disponibles. Saisissez un nom et cliquez sur **Ajouter une colonne** pour créer une colonne personnalisée. Il est conçu pour les tâches humaines : l'IA voit une tâche comme `Backlog` et ne la reçoit pas via MCP. La suppression d'une colonne lors de la sauvegarde transfère ses tâches vers le `Backlog` standard.

## Web et MCP

**Actualisation du tableau** définit le taux de rafraîchissement du tableau et des cartes en secondes (1 à 60). **Actualisation des travailleurs** met à jour le statut des travailleurs (2 à 120 secondes). Il s'agit d'une interrogation de l'interface Web, pas de la fréquence de déclenchement de l'IA. Les travailleurs ne démarrent pas automatiquement.

L'adresse du serveur Web est définie lors du démarrage de la CLI, par exemple `agentboard open --addr 127.0.0.1:7444`. Pour modifier l'adresse, le serveur doit être redémarré. La valeur par défaut est `127.0.0.1:7337`. MCP fonctionne via une commande locale distincte `agentboard mcp --project ... --worker ...` et est indépendant du port Web. Pour une base non standard, spécifiez le même `--db` dans toutes les commandes. Copiez la configuration de chaque client à partir de la carte du travailleur.
