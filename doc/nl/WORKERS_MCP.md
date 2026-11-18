# Werknemers en MCP-verbinding

## Stap 1. Maak een werknemer aan

Open in AgentBoard **Werknemers → Werknemer toevoegen**. Selecteer een klantprofiel: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot of een andere MCP-client. Het profiel vult typische kenmerken in (code, tests, Git); je kunt ze veranderen. De naam is zichtbaar op het bord en `Slug` is een korte identificatie zonder spaties voor de MCP-opdracht. Klik op **Opslaan**.

Eén werknemer komt overeen met één AI-persoonlijkheid. Creëer verschillende werknemers voor verschillende klanten of teams. Het bekwaamheidsprofiel beschrijft de specialisatie, maar breidt de rechten van de AI niet uit naar taakfasen.

## Stap 2. Kopieer de configuratie

Open de gemaakte werknemer. Het **MCP diagnostics**-blok toont het configuratiefragment en het bestand waar het moet worden toegevoegd. Klik op **MCP-configuratie kopiëren**. Als het bestand al bestaat, voegt u de voorgestelde server toe aan het bestaande `mcpServers`/`servers`/`mcp`-object zonder de andere servers te wissen.

Het hoofdcommando ziet er als volgt uit:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` van de gemaakte werknemer. De MCP-client voert deze opdracht zelf uit wanneer deze hulpmiddelen nodig heeft.

### MCP werkt zonder webinterface

`agentboard mcp` is een afzonderlijk stdio-proces dat zelf de AgentBoard-database opent. Het heeft geen `agentboard open`/`serve` nodig, het is niet afhankelijk van de map van waaruit de webinterface wordt gestart en van het aantal exemplaren van de webinterface dat actief is. Er kunnen meerdere agenten tegelijkertijd draaien, elk in hun eigen project: elke klant voert zijn eigen MCP-proces uit en ze gebruiken allemaal veilig dezelfde SQLite-database.

### Hoe een project wordt geselecteerd

De MCP-server draait in elke map. De eerste keer dat u de tool aanroept, is deze gebonden aan het project:

1. `--project PATH`, indien gespecificeerd (elke map binnen het geregistreerde project is geschikt);
2. anders de werkmappen (roots) van de opdrachtgever, indien de opdrachtgever deze rapporteert;
3. anders de map waarin de client de MCP-server heeft gestart.

`get_my_board` retourneert het geselecteerde project en de geselecteerde werknemer samen met `session` naar `binding_source`. Als de werkmappen van de klant veranderen, voert de volgende aanroep van de tool de binding opnieuw uit. Zonder `--worker` wordt de enige ingeschakelde werker van het project gebruikt.

Als binding niet mogelijk is, retourneert elk hulpprogramma een fout die de agent aan de gebruiker moet doorgeven; taakgegevens worden in dit geval niet gelezen:

| Code | Betekenis |
| --- | --- |
| `NO_PROJECTS` | Er zijn geen projecten in de database. Voer `agentboard init` uit in de projectmap en maak een werknemer aan. Geen herstart nodig. |
| `PROJECT_REQUIRED` | De werkmap bevindt zich buiten de geregistreerde projecten. In het bericht worden projecten vermeld. |
| `PROJECT_NOT_FOUND` | `--project` bevindt zich niet in een geregistreerd project. |
| `PROJECT_AMBIGUOUS` | De werkmappen van een klant omvatten meerdere projecten. Voeg `--project` toe. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Maak, activeer of selecteer een werknemer via `--worker`. |

De configuraties voor **Codex** en **Gemini CLI** zijn aangepast en gedeeld in alle mappen, dus ze hebben geen `--project`: Maak een werker met dezelfde slug in elk project waar deze client wordt uitgevoerd. Bestanden op projectniveau (Claude Code, Cursor, OpenCode, VS Code) bevatten `--project` van hun project.

### Expliciet ontwerp en contextcontrole

Als `--project` is opgegeven en de opstartmap of werkmap van de klant zich in een ander project bevindt, vragen de tools om bevestiging via MCP-elicitatie. Klanten zonder formulieruitschakeling ontvangen `PROJECT_CONTEXT_MISMATCH` voordat ze toegang krijgen tot taken. Corrigeer de configuratie, verwijder `--project` voor automatische detectie of bevestig expliciet dat u van plan bent met een ander project te werken met de vlag `--allow-project-mismatch`. Als u werkmappen wijzigt, worden eerdere bevestigingen gereset. Automatisch gedefinieerde projecten vereisen geen bevestiging.

**Algemene instellingen → MCP-projectcontext** bestuurt deze controles voor alle projecten en werknemers van dezelfde database. Standaard ingeschakeld. Indien uitgeschakeld, gebruikt MCP het opgegeven project zonder bevestigingen of mismatch-fouten. Reeds geopende sessies passen de instelling toe de volgende keer dat de tool wordt aangeroepen; taakzichtbaarheid, overgangen en regels voor menselijke tests worden nog steeds gevolgd.

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

Voorbeelden voor Codex (gebruikersconfiguratie, project automatisch gedetecteerd) en Claude Code (projectbestand) met project `C:\Projects\MyFirstProject` en werknemer `codex`:

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

In JSON wordt de Windows-backslash verdubbeld; een kant-en-klaar fragment uit de interface doet dit automatisch. De gegenereerde configuraties bevatten de daadwerkelijke `--db` van de webserver, dus HTTP en MCP gebruiken dezelfde basis, zelfs als de clientomgeving anders is.

**Ollama** biedt een lokaal model, maar vervangt de MCP-client niet. Het profiel **Ollama via OpenCode** genereert een MCP-configuratie voor OpenCode; configureer OpenCode afzonderlijk op het Ollama-model. Een andere MCP-compatibele client met Ollama is ook geschikt.

## Stap 4: Controleer uw verbinding

1. Open de werknemerskaart en klik op **Opnieuw controleren**. **Servercontrole** voert een echt uitvoerbaar stdio-bestand uit met dit project, de werker en de basis, controleert de handshake, de lijst met tools en controleert of het project en de werker zich bevinden.
2. Start of herstart de AI-client na het toevoegen van het configuratiebestand. Vraag hem om `get_my_board` te bellen.
3. **Client verbonden** verschijnt in AgentBoard en een nieuwe sessie verschijnt in de lijst. Alleen hiermee wordt de verbinding van uw klant bevestigd. Als er tools zijn, maar de client is niet verbonden, controleer dan het pad naar het programma, de naam van het configuratiebestand en de JSON/TOML-syntaxis ervan.

Start uw werksessie met `get_my_board`. AI ontvangt geen taken van `Backlog` en `Complete`, ook al kent het hun ID. AI kan zich niet voordoen als een mens en kent geen algemeen ‘beweeg ergens naartoe’-commando. Het werken met het bestandssysteem buiten AgentBoard is afhankelijk van de mogelijkheden en de sandbox van de geselecteerde client.

Elke tool declareert `title` en alle vier MCP-gedraghints (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) als expliciete Booleaanse waarden. Leeshulpmiddelen zijn alleen-lezen en idempotent; opnametools voegen alleen records toe of verplaatsen taken vooruit, dus geen enkele is storend of gaat verder dan de lokale AgentBoard-basis.

Officiële klantinstructies: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
