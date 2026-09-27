# Préparation d'une release pour Scoop

## Ce qui est déjà automatisé

`scripts/package-scoop.ps1 -Version X.Y.Z` exécute `npm ci`, build frontend, `go test ./...`, build Windows `agentboard.exe` avec numéro de version, crée un ZIP et calcule le SHA-256 de ce ZIP. À partir du **même** fichier, il crée `release/agentboard.json` avec URL, hachage, licence MIT, CLI-shim, raccourci, `checkver` et `autoupdate`. Le ZIP et le programme d'installation contiennent le fichier `LICENSE`. La base de données se trouve en dehors du répertoire d'installation, donc `persist` n'est pas nécessaire dans le manifeste.

`.github/workflows/release.yml` sur la balise `vX.Y.Z` exécute le même package dans Windows Runner, crée le programme d'installation d'Inno Setup et attache le ZIP, le programme d'installation et le manifeste à la version GitHub. L'exécution manuelle d'un flux de travail crée uniquement un artefact à tester, sans publier de version.

## Ordre de publication pour le responsable

1. Vérifiez que le code, la documentation, le numéro de version et le fichier `LICENSE` sont prêts.
2. Exécutez `./scripts/package-scoop.ps1 -Version X.Y.Z` localement. Vérifiez la sortie `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` et `agentboard version` après le déballage. Ne modifiez pas le ZIP une fois le hachage calculé.
3. Créez et soumettez la balise `vX.Y.Z`. GitHub Actions publiera une version avec ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` et manifeste. Vérifiez les trois fichiers sur la page Release et le ZIP SHA-256 du manifeste.
4. Sur une machine Windows propre avec Scoop, exécutez `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, puis `agentboard version`, `agentboard init` et `agentboard open` dans le dossier de test.
5. Pour un flux de mise à jour permanent, placez le `agentboard.json` généré dans votre propre bucket Scoop ou proposez-le dans un bucket public approprié. Revenez pour `scoop update agentboard` après la prochaine version. La commande d'installation à partir d'une URL convient à la première connaissance, tandis que le bucket est plus pratique pour les mises à jour.

Ne remplacez pas manuellement la valeur aléatoire `hash` : Scoop vérifie le contenu du ZIP téléchargé.

Pour les vérifications de manifeste, concentrez-vous sur [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [création de manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) et [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
