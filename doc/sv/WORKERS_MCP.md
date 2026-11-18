# Arbetare och MCP-anslutning

## Steg 1. Skapa en arbetare

Öppna **Arbetare → Lägg till arbetare** i AgentBoard. Välj en klientprofil: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot eller annan MCP-klient. Profilen fyller i typiska funktioner (kod, tester, Git); du kan ändra dem. Namnet är synligt på kortet och `Slug` är en kort identifierare utan mellanslag för MCP-kommandot. Klicka på **Spara**.

En arbetare motsvarar en AI-personlighet. Skapa olika arbetare för olika kunder eller team. Kapacitetsprofilen beskriver specialiseringen, men utökar inte AI:s rättigheter till uppgiftsstadier.

## Steg 2. Kopiera konfigurationen

Öppna den skapade arbetaren. Blocket **MCP-diagnostik** visar konfigurationsfragmentet och filen där det ska läggas till. Klicka på **Kopiera MCP-konfiguration**. Om filen redan finns, lägg till den föreslagna servern till det befintliga `mcpServers`/`servers`/`mcp`-objektet utan att radera de andra servrarna.

Huvudkommandot ser ut så här:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` för den skapade arbetaren. MCP-klienten kör detta kommando själv när den behöver verktyg.

### MCP fungerar utan ett webbgränssnitt

`agentboard mcp` är en separat stdio-process som själv öppnar AgentBoard-databasen. Det behöver inte `agentboard open`/`serve`, det beror inte på mappen från vilken webbgränssnittet startas och på hur många instanser av webbgränssnittet som körs. Flera agenter kan köras samtidigt, var och en i sitt eget projekt: varje klient kör sin egen MCP-process, och de använder alla säkert samma SQLite-databas.

### Hur ett projekt väljs ut

MCP-servern körs i valfri mapp. Första gången du anropar verktyget är det bundet till projektet:

1. `--project PATH`, om angivet (vilken mapp som helst i det registrerade projektet är lämplig);
2. annars klientens arbetsmappar (rötter), om klienten rapporterar dem;
3. annars mappen där klienten startade MCP-servern.

`get_my_board` returnerar det valda projektet och arbetaren till `session` tillsammans med `binding_source`. Om klientens arbetsmappar ändras, utför nästa anrop till verktyget bindningen igen. Utan `--worker` används den enda aktiverade arbetaren i projektet.

Om bindning inte är möjlig returnerar varje verktyg ett fel som agenten måste skicka till användaren; uppgiftsdata läses inte i detta fall:

| Kod | Betydelse |
| --- | --- |
| `NO_PROJECTS` | Det finns inga projekt i databasen. Kör `agentboard init` i projektmappen och skapa en arbetare. Ingen omstart behövs. |
| `PROJECT_REQUIRED` | Arbetsmappen ligger utanför registrerade projekt. Meddelandet listar projekt. |
| `PROJECT_NOT_FOUND` | `--project` finns inte i ett registrerat projekt. |
| `PROJECT_AMBIGUOUS` | En kunds arbetsmappar spänner över flera projekt. Lägg till `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Skapa, aktivera eller välj en arbetare via `--worker`. |

Konfigurationerna för **Codex** och **Gemini CLI** är anpassade och delas över alla mappar, så de har inte `--project`: Skapa en arbetare med samma slug i varje projekt där den här klienten körs. Filer på projektnivå (Claude Code, Cursor, OpenCode, VS Code) innehåller `--project` för deras projekt.

### Explicit design och kontextkontroll

Om `--project` anges och klientens startmapp eller arbetsmapp finns i ett annat projekt, ber verktygen om bekräftelse via MCP-elicitation. Klienter utan formulärframkallande får `PROJECT_CONTEXT_MISMATCH` innan de får åtkomst till uppgifter. Korrigera konfigurationen, ta bort `--project` för automatisk upptäckt, eller bekräfta uttryckligen att du tänker arbeta med ett annat projekt med `--allow-project-mismatch`-flaggan. Att ändra arbetsmappar återställer tidigare bekräftelser. Automatiskt definierade projekt kräver ingen bekräftelse.

**Globala inställningar → MCP-projektkontext** kontrollerar dessa kontroller för alla projekt och arbetare i samma databas. Aktiverad som standard. Om det är inaktiverat använder MCP det angivna projektet utan bekräftelser eller missmatchningsfel. Redan öppna sessioner tillämpar inställningen nästa gång verktyget anropas; uppgiftssynlighet, övergångar och mänskliga testregler fortsätter att följas.

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

Exempel på Codex (användarkonfiguration, projekt identifieras automatiskt) och Claude Code (projektfil) med projekt `C:\Projects\MyFirstProject` och arbetaren `codex`:

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

I JSON fördubblas Windows-omvänt snedstreck; ett färdigt fragment från gränssnittet gör detta automatiskt. De genererade konfigurationerna innehåller den faktiska `--db` för webbservern, så HTTP och MCP använder samma bas, även om klientmiljön är annorlunda.

**Ollama** tillhandahåller en lokal modell, men ersätter inte MCP-klienten. **Ollama via OpenCode**-profilen genererar en MCP-konfiguration för OpenCode; konfigurera OpenCode separat på Ollama-modellen. En annan MCP-kompatibel klient med Ollama är också lämplig.

## Steg 4: Kontrollera din anslutning

1. Öppna arbetarkortet och klicka på **Kontrollera igen**. **Servercheck** kör en riktig körbar stdio med detta projekt, arbetare och bas, kontrollerar handskakningen, listan över verktyg och att projektet och arbetaren finns.
2. Starta eller starta om AI-klienten efter att ha lagt till konfigurationsfilen. Be honom ringa `get_my_board`.
3. **Klientansluten** visas i AgentBoard och en ny session visas i listan. Endast detta bekräftar anslutningen av din klient. Om det finns verktyg, men klienten inte är ansluten, kontrollera sökvägen till programmet, namnet på konfigurationsfilen och dess JSON/TOML-syntax.

Starta din arbetssession med `get_my_board`. AI tar inte emot uppgifter från `Backlog` och `Complete`, även om den känner till deras ID. AI kan inte låtsas vara mänsklig och har inte ett allmänt "flytta någonstans"-kommando. Att arbeta med filsystemet utanför AgentBoard beror på kapaciteten och sandlådan för den valda klienten.

Varje verktyg deklarerar `title` och alla fyra MCP-beteendetips (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) som explicita booleska värden. Läsverktyg är skrivskyddade och idempotenta; inspelningsverktyg lägger bara till poster eller flyttar uppgifter framåt, så ingen är störande eller går utöver den lokala AgentBoard-basen.

Officiella kundinstruktioner: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
