# Pregătirea unei versiuni pentru Scoop

## Ce este deja automatizat

`scripts/package-scoop.ps1 -Version X.Y.Z` execută `npm ci`, construirea frontend, `go test ./...`, versiunea Windows `agentboard.exe` cu numărul de versiune, creează un ZIP și calculează SHA-256 al acelui ZIP. Din **același** fișier creează `release/agentboard.json` cu URL, hash, licență MIT, CLI-shim, scurtătură, `checkver` și `autoupdate`. ZIP și programul de instalare conțin fișierul `LICENSE`. Baza de date se află în afara directorului de instalare, deci `persist` nu este necesar în manifest.

`.github/workflows/release.yml` pe eticheta `vX.Y.Z` rulează același pachet în Windows Runner, creează programul de instalare Inno Setup și atașează ZIP, programul de instalare și manifest la versiunea GitHub. Rularea manuală a unui flux de lucru creează doar un artefact pentru testare, fără a publica o versiune.

## Ordin de postare pentru întreținător

1. Verificați dacă codul, documentația, numărul versiunii și fișierul `LICENSE` sunt gata.
2. Rulați `./scripts/package-scoop.ps1 -Version X.Y.Z` local. Verificați ieșirea `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` și `agentboard version` după despachetare. Nu modificați codul ZIP după ce hash-ul a fost calculat.
3. Creați și trimiteți eticheta `vX.Y.Z`. GitHub Actions va publica o versiune cu ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` și manifest. Verificați toate cele trei fișiere pe pagina de lansare și SHA-256 ZIP din manifest.
4. Pe o mașină Windows curată cu Scoop, rulați `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, apoi `agentboard version`, `agentboard init` și `agentboard open` în folderul de testare.
5. Pentru un feed de actualizare permanentă, plasați `agentboard.json` generat în propria dvs. găleată Scoop sau oferiți-l într-o găleată publică adecvată. Reveniți pentru `scoop update agentboard` după următoarea lansare. Comanda de instalare de la o adresă URL este potrivită pentru prima cunoștință, în timp ce găleata este mai convenabilă pentru actualizări.

Nu înlocuiți manual valoarea aleatorie `hash`: Scoop verifică conținutul fișierului ZIP descărcat.

Pentru verificările manifestului, concentrați-vă pe [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [crearea manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) și [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
