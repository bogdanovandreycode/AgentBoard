# Vorbereiten einer Veröffentlichung für Scoop

## Was bereits automatisiert ist

`scripts/package-scoop.ps1 -Version 0.2.0` führt `npm ci`, Frontend-Build, `go test ./...`, Windows-Build `agentboard.exe` mit Versionsnummer aus, erstellt eine ZIP und berechnet den SHA-256 dieser ZIP. Aus **der gleichen** Datei wird `release/agentboard.json` mit URL, Hash, CLI-Shim, Verknüpfung, `checkver` und `autoupdate` erstellt. Die Datenbank befindet sich außerhalb des Installationsverzeichnisses, daher wird `persist` im Manifest nicht benötigt.

`.github/workflows/release.yml` auf dem Tag `vX.Y.Z` führt das gleiche Paket in Windows Runner aus, erstellt das Inno Setup-Installationsprogramm und hängt die ZIP-Datei, das Installationsprogramm und das Manifest an die GitHub-Version an. Durch die manuelle Ausführung eines Workflows wird nur ein Artefakt zum Testen erstellt, ohne dass eine Version veröffentlicht wird.

## Buchungsauftrag für Betreuer

1. Überprüfen Sie, ob Code, Dokumentation und Versionsnummer vorliegen. Definieren Sie die Lizenz des Projekts: Das Manifest zeigt derzeit `Unknown` an, da die Lizenzdatei nicht im Repository definiert ist. Wenn Sie sich für eine Lizenz entscheiden, fügen Sie bitte vor der Veröffentlichung `LICENSE` hinzu und aktualisieren Sie `license` im Skript.
2. Führen Sie `./scripts/package-scoop.ps1 -Version X.Y.Z` lokal aus. Überprüfen Sie die Ausgabe von `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` und `agentboard version` nach dem Auspacken. Ändern Sie die Postleitzahl nicht, nachdem der Hash berechnet wurde.
3. Erstellen und übermitteln Sie das Tag `vX.Y.Z`. GitHub Actions wird ein Release mit ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` und Manifest veröffentlichen. Überprüfen Sie alle drei Dateien auf der Release-Seite und die SHA-256-ZIP-Datei aus dem Manifest.
4. Führen Sie auf einem sauberen Windows-Computer mit Scoop zunächst `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, dann `agentboard version`, `agentboard init` und `agentboard open` im Testordner aus.
5. Für einen dauerhaften Update-Feed platzieren Sie den generierten `agentboard.json` in Ihrem eigenen Scoop-Bucket oder bieten Sie ihn in einem geeigneten öffentlichen Bucket an. Schauen Sie nach der nächsten Veröffentlichung noch einmal nach `scoop update agentboard`. Der Installationsbefehl von einer URL eignet sich für das erste Kennenlernen, während der Bucket für Updates praktischer ist.

Ersetzen Sie den Zufallswert nicht manuell. `hash`: Scoop überprüft den Inhalt der heruntergeladenen ZIP-Datei.

Konzentrieren Sie sich bei Manifestprüfungen auf [Formatieren von Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [Erstellen von Manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) und [Autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
