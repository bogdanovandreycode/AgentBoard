# Forbereder en udgivelse til Scoop

## Hvad er allerede automatiseret

`scripts/package-scoop.ps1 -Version 0.2.0` udfører `npm ci`, frontend build, `go test ./...`, Windows build `agentboard.exe` med versionsnummer, opretter en ZIP og beregner SHA-256 for den ZIP. Fra **den samme** fil oprettes `release/agentboard.json` med URL, hash, CLI-shim, genvej, `checkver` og `autoupdate`. Databasen lever uden for installationsmappen, så `persist` er ikke nødvendig i manifestet.

`.github/workflows/release.yml` på tagget `vX.Y.Z` kører den samme pakke i Windows runner, bygger Inno Setup installationsprogrammet og vedhæfter ZIP, installationsprogrammet og manifestet til GitHub Release. Manuel kørsel af en arbejdsgang skaber kun en artefakt til test uden at udgive en udgivelse.

## Bogføringsordre for vedligeholder

1. Tjek at kode, dokumentation og versionsnummer er klar. Definer projektets licens: manifest angiver i øjeblikket `Unknown`, fordi licensfilen ikke er defineret i lageret. Hvis du vælger en licens, skal du tilføje `LICENSE` og opdatere `license` i scriptet før udgivelse.
2. Kør `./scripts/package-scoop.ps1 -Version X.Y.Z` lokalt. Kontroller `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` og `agentboard version` output efter udpakning. Ændre ikke ZIP efter hashen er blevet beregnet.
3. Opret og indsend tag `vX.Y.Z`. GitHub Actions vil udgive en udgivelse med ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` og manifest. Tjek alle tre filer på udgivelsessiden og SHA-256 ZIP fra manifest.
4. På en ren Windows-maskine med Scoop skal du køre `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, derefter `agentboard version`, `agentboard init` og `agentboard open` i testmappen.
5. For en permanent opdateringsfeed skal du placere den genererede `agentboard.json` i din egen Scoop-spand eller tilbyde den i en passende offentlig spand. Kom tilbage til `scoop update agentboard` efter næste udgivelse. Installationskommandoen fra en URL er velegnet til det første bekendtskab, mens bøtten er mere praktisk til opdateringer.

Udskift ikke den tilfældige værdi `hash` manuelt: Scoop kontrollerer indholdet af den downloadede ZIP.

For manifesttjek skal du fokusere på [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [oprette manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) og [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
