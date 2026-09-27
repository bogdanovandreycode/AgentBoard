# Forbereder en utgivelse for Scoop

## Hva er allerede automatisert

`scripts/package-scoop.ps1 -Version X.Y.Z` kjører `npm ci`, frontend build, `go test ./...`, Windows build `agentboard.exe` med versjonsnummer, oppretter en ZIP og beregner SHA-256 for den ZIP. Fra **den samme** filen opprettes `release/agentboard.json` med URL, hash, MIT-lisens, CLI-shim, snarvei, `checkver` og `autoupdate`. ZIP og installasjonsprogrammet inneholder filen `LICENSE`. Databasen lever utenfor installasjonskatalogen, så `persist` er ikke nødvendig i manifestet.

`.github/workflows/release.yml` på taggen `vX.Y.Z` kjører den samme pakken i Windows runner, bygger Inno Setup-installasjonsprogrammet og legger ved ZIP, installasjonsprogrammet og manifestet til GitHub-utgivelsen. Manuell kjøring av en arbeidsflyt skaper bare en artefakt for testing, uten å publisere en utgivelse.

## Posteringsordre for vedlikeholder

1. Bekreft at koden, dokumentasjonen, versjonsnummeret og filen `LICENSE` er klare.
2. Kjør `./scripts/package-scoop.ps1 -Version X.Y.Z` lokalt. Sjekk `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` og `agentboard version` utgang etter utpakking. Ikke endre ZIP etter at hashen er beregnet.
3. Opprett og send inn taggen `vX.Y.Z`. GitHub Actions vil publisere en utgivelse med ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` og manifest. Sjekk alle tre filene på utgivelsessiden og SHA-256 ZIP fra manifest.
4. På en ren Windows-maskin med Scoop, kjør `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, deretter `agentboard version`, `agentboard init` og `agentboard open` i testmappen.
5. For en permanent oppdateringsfeed, plasser den genererte `agentboard.json` i din egen Scoop-bøtte eller tilby den i en passende offentlig bøtte. Sjekk tilbake for `scoop update agentboard` etter neste utgivelse. Installasjonskommandoen fra en URL passer for det første bekjentskapet, mens bøtta er mer praktisk for oppdateringer.

Ikke bytt ut den tilfeldige verdien `hash` manuelt: Scoop sjekker innholdet i den nedlastede ZIP-filen.

For manifestsjekker, fokuser på [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [oppretter manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) og [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
