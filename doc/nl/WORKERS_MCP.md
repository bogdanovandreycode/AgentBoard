# Werknemers en MCP-verbinding

## Stap 1. Maak een werknemer aan

Open in AgentBoard **Werknemers → Werknemer toevoegen**. Selecteer een klantprofiel: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot of een andere MCP-client. Het profiel vult typische kenmerken in (code, tests, Git); je kunt ze veranderen. De naam is zichtbaar op het bord en `Slug` is een korte identificatie zonder spaties voor de MCP-opdracht. Klik op **Opslaan**.

Eén werknemer komt overeen met één AI-persoonlijkheid. Creëer verschillende werknemers voor verschillende klanten of teams. Het bekwaamheidsprofiel beschrijft de specialisatie, maar breidt de rechten van de AI niet uit naar taakfasen.

## Stap 2. Kopieer de configuratie

Open de gemaakte werknemer. Het **MCP diagnostics**-blok toont het configuratiefragment en het bestand waar het moet worden toegevoegd. Klik op **MCP-configuratie kopiëren**. Als het bestand al bestaat, voegt u de voorgestelde server toe aan het bestaande `mcpServers`/`servers`/`mcp`-object zonder de andere servers te wissen.

Het hoofdcommando ziet er als volgt uit:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` moet verwijzen naar de map die u hebt geregistreerd via `agentboard init`. `--worker` — `Slug` van de gemaakte werknemer. De MCP-client voert deze opdracht zelf uit wanneer deze hulpmiddelen nodig heeft. In de AgentBoard-browser kan de webserver afzonderlijk worden uitgevoerd.

Na controle toont het diagnoseblok het absolute pad naar de actieve `agentboard.exe`. Het is vooral handig bij het installeren vanuit een ZIP. Bij installatie via Scoop kunt u de opdracht `agentboard` gebruiken als de client dezelfde `PATH` ziet.

## Stap 3: Voeg een server toe aan uw client

De werkinterface heeft al een kant-en-klaar fragment. Hieronder vindt u een uitleg van waar het wordt gebruikt:

| Klant | Waar invoegen | Hoe u dit aan de clientzijde kunt controleren |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, profiel `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` in de projectmap | `claude mcp list` |
| Tweeling CLI | `%USERPROFILE%\.gemini\settings.json`, object `mcpServers` | `/mcp list` in Gemini CLI |
| Cursor | `.cursor\mcp.json`-project | lijst met MCP-servers in Cursorinstellingen |
| OpenCode | `opencode.json`-project, object `mcp` | lijst met MCP-tools in OpenCode |
| VS Code Copiloot | `.vscode\mcp.json`-project, object `servers` | commando **MCP: Lijstservers** |

Voor Codex en Claude Code, een voorbeeld met het project `C:\Projects\MyFirstProject` en de werknemer `codex`:

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

In JSON wordt de Windows-backslash verdubbeld; een kant-en-klaar fragment uit de interface doet dit automatisch. Als u de `--db` met een aangepaste basis gebruikt, voegt u deze toe aan de `args` MCP-configuratie en geeft u hetzelfde pad op als bij het starten van `open`.

**Ollama** biedt een lokaal model, maar vervangt de MCP-client niet. Het profiel **Ollama via OpenCode** genereert een MCP-configuratie voor OpenCode; configureer OpenCode afzonderlijk op het Ollama-model. Een andere MCP-compatibele client met Ollama is ook geschikt.

## Stap 4: Controleer uw verbinding

1. Open de werknemerskaart en klik op **Opnieuw controleren**. **Servercontrole** moet het aantal MCP-tools tonen. Dit is een interne protocolcontrole en detectie van servertools.
2. Start of herstart de AI-client na het toevoegen van het configuratiebestand. Vraag hem om `get_my_board` te bellen.
3. **Client verbonden** verschijnt in AgentBoard en een nieuwe sessie verschijnt in de lijst. Alleen hiermee wordt de verbinding van uw klant bevestigd. Als er tools zijn, maar de client is niet verbonden, controleer dan het pad naar het programma, de naam van het configuratiebestand en de JSON/TOML-syntaxis ervan.

Start uw werksessie met `get_my_board`. AI ontvangt geen taken van `Backlog` en `Complete`, ook al kent het hun ID. AI kan zich niet voordoen als een mens en kent geen algemeen ‘beweeg ergens naartoe’-commando. Het werken met het bestandssysteem buiten AgentBoard is afhankelijk van de mogelijkheden en de sandbox van de geselecteerde client.

Officiële klantinstructies: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
