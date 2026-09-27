# Worker- und MCP-Verbindung

## Schritt 1. Erstellen Sie einen Worker

Öffnen Sie in AgentBoard **Arbeiter → Arbeiter hinzufügen**. Wählen Sie ein Client-Profil aus: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama über OpenCode, VS Code Copilot oder einen anderen MCP-Client. Das Profil füllt typische Funktionen (Code, Tests, Git) aus; Sie können sie ändern. Der Name ist auf der Platine sichtbar und `Slug` ist eine kurze Kennung ohne Leerzeichen für den MCP-Befehl. Klicken Sie auf **Speichern**.

Ein Arbeiter entspricht einer KI-Persönlichkeit. Erstellen Sie unterschiedliche Mitarbeiter für unterschiedliche Kunden oder Teams. Das Fähigkeitsprofil beschreibt die Spezialisierung, erweitert jedoch nicht die Rechte der KI auf Aufgabenschritte.

## Schritt 2. Kopieren Sie die Konfiguration

Öffnen Sie den erstellten Worker. Der Block **MCP-Diagnose** zeigt das Konfigurationsfragment und die Datei, in der es hinzugefügt werden soll. Klicken Sie auf **MCP-Konfiguration kopieren**. Wenn die Datei bereits vorhanden ist, fügen Sie den vorgeschlagenen Server zum vorhandenen Objekt `mcpServers`/`servers`/`mcp` hinzu, ohne die anderen Server zu löschen.

Der Hauptbefehl sieht so aus:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` sollte auf den Ordner verweisen, den Sie über `agentboard init` registriert haben. `--worker` – `Slug` des erstellten Workers. Der MCP-Client führt diesen Befehl selbst aus, wenn er Tools benötigt. Im AgentBoard-Browser kann der Webserver separat ausgeführt werden.

Nach der Prüfung zeigt der Diagnoseblock den absoluten Pfad zum laufenden `agentboard.exe` an. Dies ist besonders nützlich, wenn Sie von einer ZIP-Datei installieren. Bei der Installation über Scoop können Sie den Befehl `agentboard` verwenden, wenn der Client denselben `PATH` sieht.

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

Für Codex und Claude Code ein Beispiel mit dem Projekt `C:\Projects\MyFirstProject` und dem Worker `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

In JSON wird der Windows-Backslash verdoppelt; Ein vorgefertigtes Fragment aus der Schnittstelle erledigt dies automatisch. Wenn Sie `--db` mit einer benutzerdefinierten Basis verwenden, fügen Sie diese zur MCP-Konfiguration `args` hinzu und geben Sie denselben Pfad wie beim Starten von `open` an.

**Ollama** stellt ein lokales Modell bereit, ersetzt jedoch nicht den MCP-Client. Das Profil **Ollama via OpenCode** generiert eine MCP-Konfiguration für OpenCode; Konfigurieren Sie OpenCode separat auf dem Ollama-Modell. Ein weiterer MCP-kompatibler Client mit Ollama ist ebenfalls geeignet.

## Schritt 4: Überprüfen Sie Ihre Verbindung

1. Öffnen Sie die Arbeitnehmerkarte und klicken Sie auf **Erneut prüfen**. **Serverprüfung** sollte die Anzahl der MCP-Tools anzeigen. Hierbei handelt es sich um eine interne Protokollprüfung und Erkennung von Servertools.
2. Starten oder starten Sie den AI-Client neu, nachdem Sie die Konfigurationsdatei hinzugefügt haben. Bitten Sie ihn, `get_my_board` anzurufen.
3. **Client verbunden** wird in AgentBoard angezeigt und eine neue Sitzung wird in der Liste angezeigt. Erst dadurch wird die Verbindung Ihres Clients bestätigt. Wenn Tools vorhanden sind, der Client jedoch nicht verbunden ist, überprüfen Sie den Pfad zum Programm, den Namen der Konfigurationsdatei und deren JSON/TOML-Syntax.

Beginnen Sie Ihre Arbeitssitzung mit `get_my_board`. AI empfängt keine Aufgaben von `Backlog` und `Complete`, auch wenn es deren ID kennt. KI kann nicht vorgeben, ein Mensch zu sein, und verfügt nicht über einen allgemeinen Befehl, sich irgendwohin zu bewegen. Die Arbeit mit dem Dateisystem außerhalb von AgentBoard hängt von den Fähigkeiten und der Sandbox des ausgewählten Clients ab.

Offizielle Kundenanweisungen: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
