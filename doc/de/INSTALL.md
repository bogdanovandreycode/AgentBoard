# Installation und Start unter Windows

## Installer (empfohlen)

Laden Sie `agentboard-VERSION-windows-amd64-setup.exe` von der Seite [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases)] herunter. Der Installationsassistent schlägt einen Ordner vor; Der Standardwert ist `C:\AI\AgentBoard`. Es kopiert `agentboard.exe` und die Dokumentation, erstellt eine Verknüpfung und fügt den ausgewählten Ordner dem System `PATH` hinzu. Öffnen Sie nach der Installation ein neues Terminal, damit der Befehl `agentboard` verfügbar wird.

Führen Sie in Ihrem Projektordner Folgendes aus:

```powershell
agentboard init
agentboard open
```

Die Schnittstelle ist im `agentboard.exe` integriert; Es ist keine separate Installation von Go oder Node.js erforderlich. Die Daten werden in `%AppData%\AgentBoard` gespeichert und bleiben erhalten, wenn das Programm aktualisiert oder deinstalliert wird. Durch die Deinstallation über „Installierte Anwendungen“ werden Verknüpfungen und der Eintrag von `PATH` entfernt.

## Scoop installieren

Wenn Scoop noch nicht installiert ist, öffnen Sie PowerShell als Ihr regulärer Benutzer und befolgen Sie die offiziellen Anweisungen von Scoop](https://scoop.sh/). Wenn auf Ihrem Firmencomputer Einschränkungen bestehen, wenden Sie sich an Ihren Administrator. AgentBoard kann auch ohne Scoop aus einer ZIP-Datei gestartet werden.

Alternativ installieren Sie AgentBoard über Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Die Verfügbarkeit des Manifests im GitHub-Release kann auf [Seite Releases](https://github.com/bogdanovandreycode/AgentBoard/releases)] überprüft werden. Nach dem Hinzufügen des Manifests zu Scoop kann der Bucket nach Bucket-Namen installiert und mit dem Befehl `scoop update agentboard` aktualisiert werden.

## ZIP ohne Scoop

Laden Sie `agentboard-VERSION-windows-amd64.zip` von Releases herunter und entpacken Sie es beispielsweise in `C:\Tools\AgentBoard`. In PowerShell im Projektordner:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Geben Sie für einen AI-Client den vollständigen Pfad zu `agentboard.exe` in seiner MCP-Konfiguration an, wenn sich das Programm nicht in `PATH` befindet.

## Aus dem Quellcode erstellen

Installieren Sie die Go-Version von `go.mod` und Node.js 22 oder höher. In PowerShell im Stammverzeichnis des Repositorys:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` baut das Frontend in `internal/webui/dist` zusammen, führt Go-Tests durch und baut einen `agentboard.exe` zusammen. Die Reihenfolge ist wichtig: Die Schnittstelle wird beim Erstellen von Go in die Binärdatei integriert. Wenn der alte `agentboard.exe open` ausgeführt wird, stoppen Sie ihn vor dem Neuaufbau (Strg+C), andernfalls lässt Windows das Ersetzen der Datei nicht zu.

## Wo sind die Daten?

- `%AppData%\AgentBoard\agentboard.db` – Aufgaben, Projekte, Arbeiter und Einstellungen. Sie können eine andere Datei mit dem Flag `--db` angeben, aber `init`, `open`/`serve` und `mcp` müssen den **gleichen Pfad** haben.
- `<ваш проект>\.agentboard\project.json` – Projektkennung. Diese Datei enthält keine Aufgaben.
- Der Server hört standardmäßig nur auf `127.0.0.1:7337`. Geben Sie vor dem Projektpfad eine andere Adresse `--addr` an: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` öffnet die Projektseite und verwendet einen bereits laufenden Server unter dieser Adresse erneut. Wenn mehrere Projekte in einem Browser geöffnet sind, wählen Sie diese über die Liste links aus.
