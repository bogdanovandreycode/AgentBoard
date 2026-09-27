# Förbereder en release för Scoop

## Vad är redan automatiserat

`scripts/package-scoop.ps1 -Version 0.2.0` kör `npm ci`, frontend build, `go test ./...`, Windows build `agentboard.exe` med versionsnummer, skapar en ZIP och beräknar SHA-256 för den ZIP. Från **samma** fil skapas `release/agentboard.json` med URL, hash, CLI-shim, genväg, `checkver` och `autoupdate`. Databasen finns utanför installationskatalogen, så `persist` behövs inte i manifestet.

`.github/workflows/release.yml` på taggen `vX.Y.Z` kör samma paket i Windows runner, bygger installationsprogrammet för Inno Setup och bifogar ZIP, installationsprogrammet och manifestet till GitHub Release. Att köra ett arbetsflöde manuellt skapar bara en artefakt för testning, utan att publicera en version.

## Bokföringsorder för underhållare

1. Kontrollera att koden, dokumentationen och versionsnumret är klara. Definiera projektets licens: manifest anger för närvarande `Unknown` eftersom licensfilen inte är definierad i arkivet. Om du väljer en licens, lägg till `LICENSE` och uppdatera `license` i skriptet innan release.
2. Kör `./scripts/package-scoop.ps1 -Version X.Y.Z` lokalt. Kontrollera `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` och `agentboard version` utgången efter uppackning. Ändra inte ZIP efter att hashen har beräknats.
3. Skapa och skicka taggen `vX.Y.Z`. GitHub Actions kommer att publicera en release med ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` och manifest. Kontrollera alla tre filerna på releasesidan och SHA-256 ZIP från manifest.
4. På en ren Windows-maskin med Scoop, kör `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, sedan `agentboard version`, `agentboard init` och `agentboard open` i testmappen.
5. För ett permanent uppdateringsflöde, placera den genererade `agentboard.json` i din egen Scoop-hink eller bjud den i en lämplig offentlig hink. Kom tillbaka för `scoop update agentboard` efter nästa utgåva. Installationskommandot från en URL är lämpligt för den första bekantskapen, medan hinken är bekvämare för uppdateringar.

Byt inte ut det slumpmässiga värdet `hash` manuellt: Scoop kontrollerar innehållet i den nedladdade ZIP-filen.

För manifestkontroller, fokusera på [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [skapa manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) och [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
