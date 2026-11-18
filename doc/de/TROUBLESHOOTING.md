# Problemlösung

| Symptom | Was ist zu überprüfen |
| --- | --- |
| `agentboard` nicht gefunden | Starten Sie PowerShell nach Scoop neu. Verwenden Sie bei der Installation von einer ZIP-Datei den vollständigen Pfad zu `agentboard.exe`. |
| Webseite mit alter Benutzeroberfläche nach dem Erstellen | Stoppen Sie den laufenden Server Strg+C. Führen Sie `./scripts/build.ps1` aus und starten Sie eine neue Binärdatei. Vite muss vor Go erstellt werden, da die Schnittstelle in die EXE-Datei integriert ist. Aktualisieren Sie die Seite Strg+F5. |
| Port 7337 belegt | AgentBoard läuft möglicherweise bereits. Öffnen Sie `http://127.0.0.1:7337` oder beenden Sie den alten Prozess. Für andere Ports verwenden Sie `--addr`. |
| Projekt nicht gefunden | Führen Sie im gewünschten Ordner `agentboard init` aus. Dann `agentboard open` daraus oder `agentboard open C:\путь\к\проекту`. |
| Der Worker sieht die Aufgabe nicht | Die Aufgabe muss diesem bestimmten Arbeiter zugewiesen sein und sich in `Features`, `In progress`, `Testing` oder `Verification` befinden. AI erkennt `Backlog`, `Complete` und benutzerdefinierte Spalten nicht. |
| `NO_PROJECTS` | In der MCP-Datenbank sind keine Projekte vorhanden. Führen Sie `agentboard init` im Projektordner aus und erstellen Sie einen Worker. Der Agent kann den Anruf ohne Neustart wiederholen. Überprüfen Sie, ob MCP und Webschnittstelle denselben `--db` verwenden. |
| `PROJECT_REQUIRED` / `PROJECT_AMBIGUOUS` | Ohne `--project` befinden sich weder der Arbeitsordner des Clients noch der MCP-Startordner in genau demselben registrierten Projekt. Führen Sie den Agenten im Projektordner aus, registrieren Sie ihn über `agentboard init` oder fügen Sie `--project` hinzu. Es ist nicht erforderlich, die Weboberfläche zu starten. |
| `PROJECT_NOT_FOUND` / `WORKER_NOT_FOUND` / `WORKER_REQUIRED` / `NO_WORKERS` | Schauen Sie sich `--project` an, die Arbeiterschnecke in diesem Projekt, und `--db`. Die Nachricht listet registrierte Projekte oder aktivierte Worker auf. Eine benutzerdefinierte MCP-Konfiguration ohne `--project` erfordert einen Worker mit demselben Slug in jedem Projekt. |
| `PROJECT_CONTEXT_MISMATCH` | Der MCP-Startordner oder Client-Arbeitsordner stimmt nicht mit dem angegebenen `--project` überein. Bitte bestätigen Sie die Frage, wenn dies beabsichtigt ist. Andernfalls entfernen Sie `--project` (automatische Erkennung) oder geben Sie das richtige an. |
| MCP-Serverprüfung vorhanden, aber der Client ist nicht verbunden | Bei der Serverprüfung werden die Einstellungen des externen Clients nicht überprüft. Starten Sie den Client neu, überprüfen Sie seine Konfigurationsdatei, den Pfad zu `agentboard.exe`, `--project`, `--worker` und allgemein `--db`. Bitten Sie darum, `get_my_board` anzurufen. |
| Arbeiter offline | Der Client hat MCP möglicherweise beendet oder noch nicht gestartet. Nach 90 Sekunden ohne Heartbeat gilt die Sitzung als getrennt. |
| JSON-Datei wird nicht importiert | Überprüfen Sie `version: 1`, erforderliche `title`, vorhandene `Slug`-Worker und Eigenschaftsnamen. JSON erlaubt keine Kommentare oder nachgestellten Kommas. |
| `agentboard.exe` kann nicht neu erstellt werden | Windows kann eine laufende EXE-Datei nicht ersetzen. Stoppen Sie den Server mit Strg+C und versuchen Sie den Build erneut. |
| Aufgaben sind nach dem Update verschwunden | Stellen Sie sicher, dass `--db` nicht auf eine andere Datei verweist und dass Sie als derselbe Windows-Benutzer angemeldet sind. Die Standardbasis befindet sich in `%AppData%\AgentBoard`. |

Wenn der Fehler nicht beschrieben wird, erfassen Sie den genauen Text der Meldung, die `agentboard version`- und Windows-Versionen sowie die Schritte zum erneuten Versuch. Veröffentlichen Sie keine privaten Projektdaten oder Datenbankinhalte in einer offenen Ausgabe.
