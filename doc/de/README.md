# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Vorschau AgentBoard](../../assets/social-preview.png)

**[Download für Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Dokumentationsseite](https://bogdanovandreycode.github.io/AgentBoard/) · [Lizenz MIT](../../LICENSE)**

AgentBoard ist ein lokales Taskboard, in dem menschliche und KI-Mitarbeiter an gemeinsamen Aufgaben arbeiten, jedoch unterschiedliche Rechte haben. Die Anwendung wird mit einer Datei `agentboard.exe` gestartet, öffnet die Weboberfläche im Browser und stellt den Mitarbeitern über `stdio` einen separaten MCP-Server zur Verfügung. Die Daten verbleiben auf Ihrem Computer.

**[Von vorne beginnen](START_HERE.md) · [Mit Aufgaben arbeiten](TASKS.md) · [KI über MCP verbinden](WORKERS_MCP.md) · [JSON importieren](IMPORT.md) · [Einstellungen](SETTINGS.md) · [Probleme lösen](TROUBLESHOOTING.md)**

## MVP-Funktionen

- Lokales Projektboard mit Aufgabensuche, JSON-Import, benutzerdefinierten Spalten und Eigenschaften.
- MCP-Zugriff für einen bestimmten Mitarbeiter mit begrenzten KI-Übergängen und endgültiger menschlicher Übernahme von Aufgaben.
- Allgemeiner Aufgabenverlauf, Prüfanweisungen, Artefakte, KI-Kosten und Arbeiterverbindungsdiagnose.
- Windows-Installer, ZIP- und Scoop-Manifest; Die Weboberfläche ist in die ausführbare Datei integriert.

AgentBoard ist für einen vertrauenswürdigen lokalen Benutzer konzipiert. Die Anwendung hostet keine Projekte in der Cloud und startet selbst keine KI-Clients; Verbinden Sie ggf. einen MCP-kompatiblen Client mit dem Worker.

## In fünf Minuten

1. Laden Sie das `agentboard-VERSION-windows-amd64-setup.exe`-Installationsprogramm von [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases)] herunter. Es wird ein Ordner vorgeschlagen (standardmäßig `C:\AI\AgentBoard`) und zu `PATH` hinzugefügt. Scoop und ZIP sind ebenfalls verfügbar.
2. Öffnen Sie PowerShell in Ihrem Projektordner, zum Beispiel `C:\Projects\MyApp`.
3. Führen Sie `agentboard init` aus (für ZIP: vollständiger Pfad zu `agentboard.exe` und `init`).
4. Führen Sie `agentboard open` aus. `http://127.0.0.1:7337` wird geöffnet.
5. Fügen Sie über die Schaltfläche **Neue Aufgabe** eine Aufgabe hinzu. Öffnen Sie für einen AI-Worker **Worker → Worker hinzufügen**, wählen Sie das Client-Profil aus und kopieren Sie die MCP-Konfiguration.

Wenn Sie noch keinen Projektordner haben, erstellen Sie einen im Windows Explorer. Ein Projekt kann ein beliebiger Ordner sein, auch ohne Git und Code.

## Installation über Scoop

In PowerShell mit [Scoop](https://scoop.sh/)] bereits nach der Veröffentlichung installiert:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Für Entwickler gibt es [Build from Source ](INSTALL.md). Der Release-Workflow erstellt ein Installationsprogramm, ein ZIP- und Scoop-Manifest mit SHA-256 aus demselben Artefakt. Aktualisieren der über Scoop installierten Version: `scoop update agentboard` nach dem Hinzufügen des Manifests zum Bucket; Details - [Release-Vorbereitung](SCOOP_RELEASE.md).

## Wie das Board aufgebaut ist

`Backlog → Features → In progress → Testing → Verification → Complete`

Eine Person kann Aufgaben auf der Tafel verschieben. KI kann nur `Features → In progress → Testing → Verification` bewegen; Die endgültige Aufnahme in `Complete` erfolgt durch einen Menschen. Benutzerspalten sind für Menschen gedacht: Die darin enthaltene Aufgabe bleibt für MCP im Status `Backlog`. Testmodi: KI, Mensch und Hybrid. Der Aufgabe sind Verlauf, Tests, Artefakte und KI-Kosten beigefügt.

## Teams

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registriert den Ordner und schreibt dort nur `.agentboard/project.json`. Die SQLite-Produktionsdaten befinden sich im Windows-Benutzerkonfigurationsverzeichnis (`%AppData%\AgentBoard\agentboard.db`), außerhalb des Projekts und außerhalb der Scoop-Installation. Durch das Entfernen oder Aktualisieren des Pakets sollten diese Daten nicht entfernt werden. Erstellen Sie vor der Übertragung auf einen anderen Computer eine Kopie der Datenbank, während AgentBoard gestoppt ist.

Arbeiter – logische Konten; AgentBoard selbst führt weder Codex, Claude noch einen anderen KI-Client aus. Der Client startet einen lokalen MCP-Prozess für einen bestimmten Worker. MCP ist die Anwendungsberechtigungsgrenze. Verwenden Sie für die Dateiisolierung die AI-Client-Sandbox.

## Für Entwickler

Stack: Go, SQLite, offizielles MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Zuerst das Frontend erstellen, dann los: `./scripts/build.ps1`. Webdateien werden über `go:embed` in die Binärdatei eingebunden. Die Architektur und API sind in [doc/ARCHITECTURE.md](ARCHITECTURE.md)] beschrieben.

## Lizenz

AgentBoard ist ein kostenloses Open-Source-Projekt unter der Lizenz MIT](../../LICENSE). Kommerzielle Nutzung, Änderung, Verzweigung und Weiterverbreitung sind gestattet, sofern der Urheberrechtsvermerk und die Lizenz eingehalten werden.
