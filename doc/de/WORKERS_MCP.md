# Worker- und MCP-Verbindung

## Schritt 1. Erstellen Sie einen Worker

Öffnen Sie in AgentBoard **Arbeiter → Arbeiter hinzufügen**. Wählen Sie ein Client-Profil aus: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama über OpenCode, VS Code Copilot oder einen anderen MCP-Client. Das Profil füllt typische Funktionen (Code, Tests, Git) aus; Sie können sie ändern. Der Name ist auf der Platine sichtbar und `Slug` ist eine kurze Kennung ohne Leerzeichen für den MCP-Befehl. Klicken Sie auf **Speichern**.

Ein Arbeiter entspricht einer KI-Persönlichkeit. Erstellen Sie unterschiedliche Mitarbeiter für unterschiedliche Kunden oder Teams. Das Fähigkeitsprofil beschreibt die Spezialisierung, erweitert jedoch nicht die Rechte der KI auf Aufgabenschritte.

## Schritt 2. Kopieren Sie die Konfiguration

Öffnen Sie den erstellten Worker. Der Block **MCP-Diagnose** zeigt das Konfigurationsfragment und die Datei, in der es hinzugefügt werden soll. Klicken Sie auf **MCP-Konfiguration kopieren**. Wenn die Datei bereits vorhanden ist, fügen Sie den vorgeschlagenen Server zum vorhandenen Objekt `mcpServers`/`servers`/`mcp` hinzu, ohne die anderen Server zu löschen.

Der Hauptbefehl sieht so aus:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` des erstellten Workers. Der MCP-Client führt diesen Befehl selbst aus, wenn er Tools benötigt.

### MCP funktioniert ohne Webinterface

`agentboard mcp` ist ein separater stdio-Prozess, der selbst die AgentBoard-Datenbank öffnet. Es ist kein `agentboard open`/`serve` erforderlich, es hängt nicht vom Ordner ab, aus dem die Webschnittstelle gestartet wird, und auch nicht davon, wie viele Instanzen der Webschnittstelle ausgeführt werden. Mehrere Agenten können gleichzeitig ausgeführt werden, jeder in seinem eigenen Projekt: Jeder Client führt seinen eigenen MCP-Prozess aus und alle verwenden sicher dieselbe SQLite-Datenbank.

### So wird ein Projekt ausgewählt

Der MCP-Server läuft in einem beliebigen Ordner. Beim ersten Aufruf des Tools wird es an das Projekt gebunden:

1. `--project PATH`, falls angegeben (jeder Ordner innerhalb des registrierten Projekts ist geeignet);
2. andernfalls die Arbeitsordner (Roots) des Kunden, sofern diese vom Kunden gemeldet werden;
3. andernfalls der Ordner, in dem der Client den MCP-Server gestartet hat.

`get_my_board` gibt das ausgewählte Projekt und den ausgewählten Worker zusammen mit `session` an `binding_source` zurück. Wenn sich die Arbeitsordner des Clients ändern, führt der nächste Aufruf des Tools die Bindung erneut durch. Ohne `--worker` wird der einzige aktivierte Worker des Projekts verwendet.

Wenn die Bindung nicht möglich ist, gibt jedes Tool einen Fehler zurück, den der Agent an den Benutzer weitergeben muss; Auftragsdaten werden in diesem Fall nicht gelesen:

| Code | Bedeutung |
| --- | --- |
| `NO_PROJECTS` | Es sind keine Projekte in der Datenbank vorhanden. Führen Sie `agentboard init` im Projektordner aus und erstellen Sie einen Worker. Kein Neustart erforderlich. |
| `PROJECT_REQUIRED` | Der Arbeitsordner liegt außerhalb der registrierten Projekte. Die Nachricht listet Projekte auf. |
| `PROJECT_NOT_FOUND` | `--project` befindet sich nicht in einem registrierten Projekt. |
| `PROJECT_AMBIGUOUS` | Die Arbeitsordner eines Kunden umfassen mehrere Projekte. Fügen Sie `--project` hinzu. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Erstellen, aktivieren oder wählen Sie einen Worker über `--worker` aus. |

Die Konfigurationen für **Codex** und **Gemini CLI** sind benutzerdefiniert und werden in allen Ordnern gemeinsam genutzt, sodass sie nicht über `--project` verfügen: Erstellen Sie in jedem Projekt, in dem dieser Client ausgeführt wird, einen Worker mit demselben Slug. Dateien auf Projektebene (Claude Code, Cursor, OpenCode, VS Code) enthalten `--project` ihres Projekts.

### Explizites Design und Kontextprüfung

Wenn `--project` angegeben ist und sich der Startordner oder Arbeitsordner des Clients in einem anderen Projekt befindet, fordern die Tools eine Bestätigung per MCP-Ermittlung an. Kunden ohne Formularerhebung erhalten `PROJECT_CONTEXT_MISMATCH`, bevor sie auf Aufgaben zugreifen. Korrigieren Sie die Konfiguration, entfernen Sie `--project` für die automatische Erkennung oder bestätigen Sie explizit, dass Sie beabsichtigen, mit einem anderen Projekt mit dem Flag `--allow-project-mismatch` zu arbeiten. Durch das Ändern von Arbeitsordnern werden vorherige Bestätigungen zurückgesetzt. Automatisch definierte Projekte erfordern keine Bestätigung.

**Globale Einstellungen → MCP-Projektkontext** steuert diese Prüfungen für alle Projekte und Worker derselben Datenbank. Standardmäßig aktiviert. Wenn diese Option deaktiviert ist, verwendet MCP das angegebene Projekt ohne Bestätigungen oder Nichtübereinstimmungsfehler. Bereits geöffnete Sitzungen übernehmen die Einstellung beim nächsten Aufruf des Tools; Aufgabensichtbarkeit, Übergänge und menschliche Testregeln werden weiterhin befolgt.

Nach der Überprüfung zeigt der Diagnoseblock den absoluten Pfad zum laufenden `agentboard.exe` an. Dies ist besonders nützlich, wenn Sie von einer ZIP-Datei installieren. Bei der Installation über Scoop können Sie den Befehl `agentboard` verwenden, wenn der Client denselben `PATH` sieht.

## Schritt 3: Fügen Sie Ihrem Client einen Server hinzu

Die Worker-Schnittstelle verfügt bereits über ein vorgefertigtes Fragment. Nachfolgend finden Sie eine Erklärung, wo es verwendet wird:

| Kunde | Wo soll | eingefügt werden? So überprüfen Sie auf der Clientseite |
| --- | --- | --- |
| Kodex | `%USERPROFILE%\.codex\config.toml`, Abschnitt `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` im Projektordner | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, Objekt `mcpServers` | `/mcp list` in Gemini CLI |
| Cursor | `.cursor\mcp.json`-Projekt | Liste der MCP-Server in den Cursor-Einstellungen |
| OpenCode | `opencode.json` Projekt, Objekt `mcp` | Liste der MCP-Tools in OpenCode |
| VS-Code-Copilot | `.vscode\mcp.json` Projekt, Objekt `servers` | Befehl **MCP: Server auflisten** |

Beispiele für Codex (Benutzerkonfiguration, Projekt automatisch erkannt) und Claude Code (Projektdatei) mit Projekt `C:\Projects\MyFirstProject` und Worker `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

In JSON wird der Windows-Backslash verdoppelt; Ein vorgefertigtes Fragment aus der Schnittstelle erledigt dies automatisch. Die generierten Konfigurationen enthalten den tatsächlichen `--db` des Webservers, sodass HTTP und MCP dieselbe Basis verwenden, auch wenn die Clientumgebung unterschiedlich ist.

**Ollama** stellt ein lokales Modell bereit, ersetzt jedoch nicht den MCP-Client. Das Profil **Ollama via OpenCode** generiert eine MCP-Konfiguration für OpenCode; Konfigurieren Sie OpenCode separat auf dem Ollama-Modell. Ein weiterer MCP-kompatibler Client mit Ollama ist ebenfalls geeignet.

## Schritt 4: Überprüfen Sie Ihre Verbindung

1. Öffnen Sie die Arbeitnehmerkarte und klicken Sie auf **Erneut prüfen**. **Serverprüfung** führt eine echte ausführbare stdio-Datei mit diesem Projekt, Worker und Basis aus, prüft den Handshake, die Liste der Tools und ob das Projekt und der Worker gefunden werden.
2. Starten oder starten Sie den AI-Client neu, nachdem Sie die Konfigurationsdatei hinzugefügt haben. Bitten Sie ihn, `get_my_board` anzurufen.
3. **Client verbunden** wird in AgentBoard angezeigt und eine neue Sitzung wird in der Liste angezeigt. Erst dadurch wird die Verbindung Ihres Clients bestätigt. Wenn Tools vorhanden sind, der Client jedoch nicht verbunden ist, überprüfen Sie den Pfad zum Programm, den Namen der Konfigurationsdatei und deren JSON/TOML-Syntax.

Beginnen Sie Ihre Arbeitssitzung mit `get_my_board`. AI empfängt keine Aufgaben von `Backlog` und `Complete`, auch wenn es deren ID kennt. KI kann nicht vorgeben, ein Mensch zu sein, und verfügt nicht über einen allgemeinen Befehl, sich irgendwohin zu bewegen. Die Arbeit mit dem Dateisystem außerhalb von AgentBoard hängt von den Fähigkeiten und der Sandbox des ausgewählten Clients ab.

Jedes Tool deklariert `title` und alle vier MCP-Verhaltenshinweise (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) als explizite boolesche Werte. Lesewerkzeuge sind schreibgeschützt und idempotent; Aufzeichnungstools fügen nur Datensätze hinzu oder verschieben Aufgaben nach vorne, sodass keine störend ist oder über die lokale AgentBoard-Basis hinausgeht.

Offizielle Kundenanweisungen: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
