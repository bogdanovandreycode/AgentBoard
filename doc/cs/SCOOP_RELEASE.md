# Příprava vydání pro Scoop

## Co je již automatizované

`scripts/package-scoop.ps1 -Version 0.2.0` provede `npm ci`, frontend build, `go test ./...`, Windows build `agentboard.exe` s číslem verze, vytvoří ZIP a vypočítá SHA-256 tohoto ZIP. Z **stejného** souboru vytvoří `release/agentboard.json` s URL, hash, CLI-shim, zástupce, `checkver` a `autoupdate`. Databáze žije mimo instalační adresář, takže `persist` není v manifestu potřeba.

`.github/workflows/release.yml` na štítku `vX.Y.Z` spouští stejné balení v programu Windows Runner, sestavuje instalační program Inno Setup a připojuje ZIP, instalační program a manifest k vydání GitHub. Ruční spuštění pracovního postupu pouze vytvoří artefakt pro testování, bez publikování vydání.

## Příkaz k zaúčtování pro správce

1. Zkontrolujte, zda jsou připraveny kód, dokumentace a číslo verze. Definujte licenci projektu: manifest aktuálně označuje `Unknown`, protože licenční soubor není definován v úložišti. Pokud zvolíte licenci, přidejte prosím `LICENSE` a aktualizujte `license` ve skriptu před vydáním.
2. Spusťte `./scripts/package-scoop.ps1 -Version X.Y.Z` lokálně. Po vybalení zkontrolujte výstup `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` a `agentboard version`. Po výpočtu hashe neměňte ZIP.
3. Vytvořte a odešlete značku `vX.Y.Z`. Akce GitHub zveřejní vydání se ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` a manifestem. Zkontrolujte všechny tři soubory na stránce Release a SHA-256 ZIP z manifestu.
4. Na čistém počítači se systémem Windows pomocí nástroje Scoop spusťte `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, poté `agentboard version`, `agentboard init` a `agentboard open` v testovací složce.
5. Pro trvalý aktualizační zdroj umístěte vygenerovaný `agentboard.json` do svého vlastního kbelíku Scoop nebo jej nabídněte ve vhodném veřejném kbelíku. Po dalším vydání se vraťte na `scoop update agentboard`. Instalační příkaz z URL je vhodný pro první seznámení, zatímco bucket je pohodlnější pro aktualizace.

Nenahrazujte ručně náhodnou hodnotu `hash`: Scoop kontroluje obsah staženého ZIP.

Pro kontroly manifestů se zaměřte na [formát Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [vytvoření manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) a [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
