# Arbejdere og MCP-forbindelse

## Trin 1. Opret en arbejder

I AgentBoard skal du åbne **Medarbejdere → Tilføj arbejder**. Vælg en klientprofil: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot eller en anden MCP-klient. Profilen udfylder typiske funktioner (kode, test, Git); du kan ændre dem. Navnet er synligt på tavlen, og `Slug` er en kort identifikator uden mellemrum til MCP-kommandoen. Klik på **Gem**.

Én arbejder svarer til én AI-personlighed. Opret forskellige medarbejdere til forskellige kunder eller teams. Kompetenceprofilen beskriver specialiseringen, men udvider ikke AI'ens rettigheder til opgavestadier.

## Trin 2. Kopier konfigurationen

Åbn den oprettede arbejder. **MCP-diagnostik**-blokken viser konfigurationsfragmentet og filen, hvor det skal tilføjes. Klik på **Kopiér MCP-konfiguration**. Hvis filen allerede findes, skal du tilføje den foreslåede server til det eksisterende `mcpServers`/`servers`/`mcp`-objekt uden at slette de andre servere.

Hovedkommandoen ser sådan ud:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` skal pege på den mappe, du har registreret via `agentboard init`. `--worker` — `Slug` af den oprettede arbejder. MCP-klienten kører selv denne kommando, når den har brug for værktøjer. I AgentBoard-browseren kan webserveren køre separat.

Efter kontrol viser diagnoseblokken den absolutte sti til den kørende `agentboard.exe`. Det er især nyttigt, når du installerer fra en ZIP. Når du installerer via Scoop, kan du bruge kommandoen `agentboard`, hvis klienten ser den samme `PATH`.

## Trin 3: Tilføj en server til din klient

Arbejdergrænsefladen har allerede et færdigt fragment. Nedenfor er en forklaring på, hvor det bruges:

| Kunde | Hvor skal du indsætte | Sådan tjekker du på klientsiden |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, afsnit `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude kode | `.mcp.json` i projektmappen | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, objekt `mcpServers` | `/mcp list` i Gemini CLI |
| Markør | `.cursor\mcp.json` projekt | liste over MCP-servere i markørindstillinger |
| OpenCode | `opencode.json` projekt, objekt `mcp` | liste over MCP-værktøjer i OpenCode |
| VS Code Copilot | `.vscode\mcp.json` projekt, objekt `servers` | kommando **MCP: Liste over servere** |

For Codex og Claude Code, et eksempel med projektet `C:\Projects\MyFirstProject` og arbejderen `codex`:

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

I JSON er Windows-omvendt skråstreg fordoblet; et færdigt fragment fra grænsefladen gør dette automatisk. Hvis du bruger `--db` med en brugerdefineret base, skal du tilføje den til `args` MCP-konfiguration og angive den samme sti, som når du starter `open`.

**Ollama** leverer en lokal model, men erstatter ikke MCP-klienten. **Ollama via OpenCode**-profilen genererer en MCP-konfiguration til OpenCode; konfigurer OpenCode separat på Ollama-modellen. En anden MCP-kompatibel klient med Ollama er også velegnet.

## Trin 4: Tjek din forbindelse

1. Åbn arbejderkortet, og klik på **Kontroller igen**. **Servercheck** bør vise antallet af MCP-værktøjer. Dette er en intern protokolkontrol og detektering af serverværktøjer.
2. Start eller genstart AI-klienten efter tilføjelse af konfigurationsfilen. Bed ham om at ringe til `get_my_board`.
3. **Kunde tilsluttet** vises i AgentBoard, og en ny session vises på listen. Kun dette bekræfter din klients forbindelse. Hvis der er værktøjer, men klienten ikke er tilsluttet, skal du kontrollere stien til programmet, navnet på konfigurationsfilen og dens JSON/TOML-syntaks.

Start din arbejdssession med `get_my_board`. AI modtager ikke opgaver fra `Backlog` og `Complete`, selvom den kender deres ID. AI kan ikke foregive at være et menneske og har ikke en generel "flytt hvor som helst"-kommando. At arbejde med filsystemet uden for AgentBoard afhænger af den valgte klients muligheder og sandkasse.

Officielle kundevejledninger: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
