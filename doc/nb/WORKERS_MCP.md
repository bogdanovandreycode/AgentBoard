# Arbeidere og MCP-tilkobling

## Trinn 1. Opprett en arbeider

I AgentBoard åpner du **Arbeidere → Legg til arbeider**. Velg en klientprofil: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot eller en annen MCP-klient. Profilen fyller ut typiske funksjoner (kode, tester, Git); du kan endre dem. Navnet er synlig på tavlen, og `Slug` er en kort identifikator uten mellomrom for MCP-kommandoen. Klikk på **Lagre**.

Én arbeider tilsvarer én AI-personlighet. Lag forskjellige arbeidere for forskjellige kunder eller team. Kompetanseprofilen beskriver spesialiseringen, men utvider ikke AIs rettigheter til oppgavestadier.

## Trinn 2. Kopier konfigurasjonen

Åpne den opprettede arbeideren. **MCP-diagnostikkblokken viser konfigurasjonsfragmentet og filen der det skal legges til. Klikk på **Kopier MCP-konfigurasjon**. Hvis filen allerede eksisterer, legg til den foreslåtte serveren til det eksisterende `mcpServers`/`servers`/`mcp`-objektet uten å slette de andre serverne.

Hovedkommandoen ser slik ut:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` skal peke til mappen du registrerte via `agentboard init`. `--worker` — `Slug` til den opprettede arbeideren. MCP-klienten kjører denne kommandoen selv når den trenger verktøy. I AgentBoard-nettleseren kan webserveren kjøres separat.

Etter kontroll viser diagnoseblokken den absolutte banen til den kjørende `agentboard.exe`. Det er spesielt nyttig når du installerer fra en ZIP. Når du installerer via Scoop, kan du bruke `agentboard`-kommandoen hvis klienten ser den samme `PATH`.

## Trinn 3: Legg til en server til klienten din

Arbeidergrensesnittet har allerede et ferdig fragment. Nedenfor er en forklaring på hvor det brukes:

| Kunde | Hvor skal du sette inn | Hvordan sjekke på klientsiden |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, seksjon `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` i prosjektmappen | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, objekt `mcpServers` | `/mcp list` i Gemini CLI |
| Markør | `.cursor\mcp.json` prosjekt | liste over MCP-servere i markørinnstillinger |
| OpenCode | `opencode.json`-prosjekt, objekt `mcp` | liste over MCP-verktøy i OpenCode |
| VS Code Copilot | `.vscode\mcp.json`-prosjekt, objekt `servers` | kommando **MCP: List servere** |

For Codex og Claude Code, et eksempel med prosjektet `C:\Projects\MyFirstProject` og arbeideren `codex`:

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

I JSON er Windows-omvendt skråstrek doblet; et ferdig fragment fra grensesnittet gjør dette automatisk. Hvis du bruker `--db` med en tilpasset base, legg den til `args` MCP-konfigurasjon og spesifiser samme bane som når du starter `open`.

**Ollama** gir en lokal modell, men erstatter ikke MCP-klienten. **Ollama via OpenCode**-profilen genererer en MCP-konfigurasjon for OpenCode; konfigurere OpenCode separat på Ollama-modellen. En annen MCP-kompatibel klient med Ollama er også egnet.

## Trinn 4: Sjekk tilkoblingen

1. Åpne arbeiderkortet og klikk på **Sjekk igjen**. **Serversjekk** skal vise antall MCP-verktøy. Dette er en intern protokollsjekk og gjenkjenning av serververktøy.
2. Start eller restart AI-klienten etter å ha lagt til konfigurasjonsfilen. Be ham ringe `get_my_board`.
3. **Klient tilkoblet** vises i AgentBoard og en ny økt vil vises i listen. Bare dette bekrefter tilkoblingen til klienten din. Hvis det finnes verktøy, men klienten ikke er tilkoblet, sjekk banen til programmet, navnet på konfigurasjonsfilen og dens JSON/TOML-syntaks.

Start arbeidsøkten med `get_my_board`. AI mottar ikke oppgaver fra `Backlog` og `Complete`, selv om den kjenner deres ID. AI kan ikke late som om det er et menneske og har ikke en generell "flytt hvor som helst"-kommando. Arbeid med filsystemet utenfor AgentBoard avhenger av egenskapene og sandkassen til den valgte klienten.

Offisielle kundeinstruksjoner: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
