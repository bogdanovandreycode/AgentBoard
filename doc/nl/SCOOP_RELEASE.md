# Een release voor Scoop voorbereiden

## Wat al geautomatiseerd is

`scripts/package-scoop.ps1 -Version 0.2.0` voert `npm ci` uit, frontend build, `go test ./...`, Windows build `agentboard.exe` met versienummer, maakt een ZIP en berekent de SHA-256 van die ZIP. Vanuit **hetzelfde** bestand wordt `release/agentboard.json` gemaakt met URL, hash, CLI-shim, snelkoppeling, `checkver` en `autoupdate`. De database bevindt zich buiten de installatiemap, dus `persist` is niet nodig in het manifest.

`.github/workflows/release.yml` op de tag `vX.Y.Z` voert dezelfde verpakking uit in Windows Runner, bouwt het Inno Setup-installatieprogramma en koppelt de ZIP, het installatieprogramma en het manifest aan de GitHub-release. Het handmatig uitvoeren van een workflow creëert alleen een artefact om te testen, zonder een release te publiceren.

## Order plaatsen voor onderhouder

1. Controleer of de code, documentatie en versienummer gereed zijn. Definieer de licentie van het project: manifest geeft momenteel `Unknown` aan omdat het licentiebestand niet is gedefinieerd in de repository. Als u een licentie kiest, voeg dan `LICENSE` toe en update `license` in het script voordat u het uitbrengt.
2. Voer `./scripts/package-scoop.ps1 -Version X.Y.Z` lokaal uit. Controleer de uitvoer van de `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` en `agentboard version` na het uitpakken. Wijzig de ZIP niet nadat de hash is berekend.
3. Maak en verzend tag `vX.Y.Z`. GitHub Actions zal een release publiceren met ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` en manifest. Controleer alle drie de bestanden op de Release-pagina en de SHA-256 ZIP uit het manifest.
4. Voer op een schone Windows-machine met Scoop `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json` uit, daarna `agentboard version`, `agentboard init` en `agentboard open` in de testmap.
5. Voor een permanente updatefeed plaatst u de gegenereerde `agentboard.json` in uw eigen Scoop-bucket of biedt u deze aan in een geschikte openbare bucket. Kom later terug voor `scoop update agentboard` na de volgende release. Het installatiecommando vanaf een URL is geschikt voor de eerste kennismaking, terwijl de bucket handiger is voor updates.

Vervang de willekeurige waarde `hash` niet handmatig: Scoop controleert de inhoud van de gedownloade ZIP.

Voor manifestcontroles concentreert u zich op [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) maken en [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
