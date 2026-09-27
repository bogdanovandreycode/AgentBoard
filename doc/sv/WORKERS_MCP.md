# Arbetare och MCP-anslutning

## Steg 1. Skapa en arbetare

Öppna **Arbetare → Lägg till arbetare** i AgentBoard. Välj en klientprofil: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot eller annan MCP-klient. Profilen fyller i typiska funktioner (kod, tester, Git); du kan ändra dem. Namnet är synligt på kortet och `Slug` är en kort identifierare utan mellanslag för MCP-kommandot. Klicka på **Spara**.

En arbetare motsvarar en AI-personlighet. Skapa olika arbetare för olika kunder eller team. Kapacitetsprofilen beskriver specialiseringen, men utökar inte AI:s rättigheter till uppgiftsstadier.

## Steg 2. Kopiera konfigurationen

Öppna den skapade arbetaren. Blocket **MCP-diagnostik** visar konfigurationsfragmentet och filen där det ska läggas till. Klicka på **Kopiera MCP-konfiguration**. Om filen redan finns, lägg till den föreslagna servern till det befintliga `mcpServers`/`servers`/`mcp`-objektet utan att radera de andra servrarna.

Huvudkommandot ser ut så här:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` bör peka på mappen du registrerade via `agentboard init`. `--worker` — `Slug` för den skapade arbetaren. MCP-klienten kör detta kommando själv när den behöver verktyg. I AgentBoard-webbläsaren kan webbservern köras separat.

Efter kontroll visar diagnosblocket den absoluta sökvägen till den körande `agentboard.exe`. Det är särskilt användbart när du installerar från en ZIP. När du installerar via Scoop kan du använda kommandot `agentboard` om klienten ser samma `PATH`.

## Steg 3: Lägg till en server till din klient

Arbetargränssnittet har redan ett färdigt fragment. Nedan finns en förklaring av var den används:

| Kund | Var ska man infoga | Så här kontrollerar du på klientsidan |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, sektion `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` i projektmappen | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, objekt `mcpServers` | `/mcp list` i Gemini CLI |
| Markör | `.cursor\mcp.json`-projekt | lista över MCP-servrar i markörinställningar |
| OpenCode | `opencode.json`-projekt, objekt `mcp` | lista över MCP-verktyg i OpenCode |
| VS Code Copilot | `.vscode\mcp.json`-projekt, objekt `servers` | kommando **MCP: Lista servrar** |

För Codex och Claude Code, ett exempel med projektet `C:\Projects\MyFirstProject` och arbetaren `codex`:

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

I JSON fördubblas Windows-omvänt snedstreck; ett färdigt fragment från gränssnittet gör detta automatiskt. Om du använder `--db` med en anpassad bas, lägg till den i `args` MCP-konfiguration och ange samma sökväg som när du startar `open`.

**Ollama** tillhandahåller en lokal modell, men ersätter inte MCP-klienten. **Ollama via OpenCode**-profilen genererar en MCP-konfiguration för OpenCode; konfigurera OpenCode separat på Ollama-modellen. En annan MCP-kompatibel klient med Ollama är också lämplig.

## Steg 4: Kontrollera din anslutning

1. Öppna arbetarkortet och klicka på **Kontrollera igen**. **Serverkontroll** bör visa antalet MCP-verktyg. Detta är en intern protokollkontroll och upptäckt av serververktyg.
2. Starta eller starta om AI-klienten efter att ha lagt till konfigurationsfilen. Be honom ringa `get_my_board`.
3. **Kunden ansluten** visas i AgentBoard och en ny session visas i listan. Endast detta bekräftar anslutningen av din klient. Om det finns verktyg, men klienten inte är ansluten, kontrollera sökvägen till programmet, namnet på konfigurationsfilen och dess JSON/TOML-syntax.

Starta din arbetssession med `get_my_board`. AI tar inte emot uppgifter från `Backlog` och `Complete`, även om den känner till deras ID. AI kan inte låtsas vara människa och har inte ett allmänt "flytta någonstans"-kommando. Att arbeta med filsystemet utanför AgentBoard beror på kapaciteten och sandlådan för den valda klienten.

Officiella kundinstruktioner: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
