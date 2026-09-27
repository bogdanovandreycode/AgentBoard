# Vorbereiten einer Veröffentlichung für Scoop

## Was bereits automatisiert ist

`scripts/package-scoop.ps1 -Version X.Y.Z` führt `npm ci`, Frontend-Build, `go test ./...`, Windows-Build `agentboard.exe` mit Versionsnummer aus, erstellt eine ZIP und berechnet den SHA-256 dieser ZIP. Aus **der gleichen** Datei wird `release/agentboard.json` mit URL, Hash, MIT-Lizenz, CLI-Shim, Verknüpfung, `checkver` und `autoupdate` erstellt. Die ZIP-Datei und das Installationsprogramm enthalten die Datei `LICENSE`. Die Datenbank befindet sich außerhalb des Installationsverzeichnisses, daher wird `persist` im Manifest nicht benötigt.

`.github/workflows/release.yml` auf Tag `vX.Y.Z` führt das gleiche Paket in Windows Runner aus, erstellt das Inno Setup-Installationsprogramm und hängt die ZIP-Datei, das Installationsprogramm und das Manifest an die GitHub-Version an. Durch die manuelle Ausführung eines Workflows wird nur ein Artefakt zum Testen erstellt, ohne dass eine Version veröffentlicht wird.

## Buchungsauftrag für Betreuer

1. Stellen Sie sicher, dass der Code, die Dokumentation, die Versionsnummer und die Datei `LICENSE` bereit sind.
2. Führen Sie `./scripts/package-scoop.ps1 -Version X.Y.Z` lokal aus. Überprüfen Sie die Ausgabe von `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` und `agentboard version` nach dem Entpacken. Ändern Sie die Postleitzahl nicht, nachdem der Hash berechnet wurde.
3. Erstellen und übermitteln Sie das Tag `vX.Y.Z`. GitHub Actions wird ein Release mit ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` und Manifest veröffentlichen. Überprüfen Sie alle drei Dateien auf der Release-Seite und die SHA-256-ZIP-Datei aus dem Manifest.
4. Führen Sie auf einem sauberen Windows-Computer mit Scoop `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, dann `agentboard version`, `agentboard init` und `agentboard open` im Testordner aus.
5. Für einen permanenten Update-Feed platzieren Sie den generierten `agentboard.json` in Ihrem eigenen Scoop-Bucket oder bieten Sie ihn in einem geeigneten öffentlichen Bucket an. Schauen Sie nach der nächsten Veröffentlichung erneut nach `scoop update agentboard`. Der Installationsbefehl von einer URL eignet sich für das erste Kennenlernen, während der Bucket für Updates praktischer ist.

Ersetzen Sie den Zufallswert nicht manuell. `hash`: Scoop überprüft den Inhalt der heruntergeladenen ZIP-Datei.

Konzentrieren Sie sich bei Manifestprüfungen auf [formatieren Sie Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [erstellen von Manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) und [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
