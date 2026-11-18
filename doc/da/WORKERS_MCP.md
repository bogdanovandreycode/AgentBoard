# Arbejdere og MCP-forbindelse

## Trin 1. Opret en arbejder

I AgentBoard skal du åbne **Medarbejdere → Tilføj arbejder**. Vælg en klientprofil: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot eller en anden MCP-klient. Profilen udfylder typiske funktioner (kode, test, Git); du kan ændre dem. Navnet er synligt på tavlen, og `Slug` er en kort identifikator uden mellemrum til MCP-kommandoen. Klik på **Gem**.

Én arbejder svarer til én AI-personlighed. Opret forskellige medarbejdere til forskellige kunder eller teams. Kompetenceprofilen beskriver specialiseringen, men udvider ikke AI'ens rettigheder til opgavestadier.

## Trin 2. Kopier konfigurationen

Åbn den oprettede arbejder. **MCP-diagnostik**-blokken viser konfigurationsfragmentet og filen, hvor det skal tilføjes. Klik på **Kopiér MCP-konfiguration**. Hvis filen allerede findes, skal du tilføje den foreslåede server til det eksisterende `mcpServers`/`servers`/`mcp`-objekt uden at slette de andre servere.

Hovedkommandoen ser sådan ud:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` af den oprettede arbejder. MCP-klienten kører selv denne kommando, når den har brug for værktøjer.

### MCP fungerer uden en webgrænseflade

`agentboard mcp` er en separat stdio-proces, der selv åbner AgentBoard-databasen. Det behøver ikke `agentboard open`/`serve`, det afhænger ikke af mappen, hvorfra webgrænsefladen startes, og hvor mange forekomster af webgrænsefladen, der kører. Flere agenter kan køre samtidigt, hver i deres eget projekt: hver klient kører sin egen MCP-proces, og de bruger alle sikkert den samme SQLite-database.

### Hvordan et projekt udvælges

MCP-serveren kører i enhver mappe. Første gang du ringer til værktøjet, er det bundet til projektet:

1. `--project PATH`, hvis specificeret (enhver mappe inde i det registrerede projekt er egnet);
2. ellers klientens arbejdsmapper (rødder), hvis klienten rapporterer dem;
3. ellers mappen, hvor klienten startede MCP-serveren.

`get_my_board` returnerer det valgte projekt og arbejder til `session` sammen med `binding_source`. Hvis klientens arbejdsmapper ændres, udfører det næste kald til værktøjet bindingen igen. Uden `--worker` bruges den eneste aktiverede arbejder i projektet.

Hvis binding ikke er mulig, returnerer hvert værktøj en fejl, som agenten skal videregive til brugeren; opgavedata læses ikke i dette tilfælde:

| Kode | Betydning |
| --- | --- |
| `NO_PROJECTS` | Der er ingen projekter i databasen. Udfør `agentboard init` i projektmappen og opret en arbejder. Ingen genstart nødvendig. |
| `PROJECT_REQUIRED` | Arbejdsmappen er uden for registrerede projekter. Meddelelsen viser projekter. |
| `PROJECT_NOT_FOUND` | `--project` er ikke inde i et registreret projekt. |
| `PROJECT_AMBIGUOUS` | En kundes arbejdsmapper spænder over flere projekter. Tilføj `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Opret, aktiver eller vælg en arbejder via `--worker`. |

Konfigurationerne for **Codex** og **Gemini CLI** er tilpasset og delt på tværs af alle mapper, så de ikke har `--project`: Opret en arbejder med den samme slug i hvert projekt, hvor denne klient kører. Projektniveaufiler (Claude Code, Cursor, OpenCode, VS Code) indeholder `--project` af deres projekt.

### Eksplicit design og kontekstkontrol

Hvis `--project` er angivet, og klientens startmappe eller arbejdsmappe er i et andet projekt, beder værktøjerne om bekræftelse via MCP-elicitation. Klienter uden formfremkaldelse modtager `PROJECT_CONTEXT_MISMATCH`, før de får adgang til opgaver. Ret konfigurationen, fjern `--project` for automatisk registrering, eller bekræft udtrykkeligt, at du har til hensigt at arbejde med et andet projekt med `--allow-project-mismatch`-flaget. Ændring af arbejdsmapper nulstiller tidligere bekræftelser. Automatisk definerede projekter kræver ikke bekræftelse.

**Globale indstillinger → MCP-projektkontekst** kontrollerer disse kontroller for alle projekter og arbejdere i den samme database. Aktiveret som standard. Hvis deaktiveret, bruger MCP det angivne projekt uden bekræftelser eller uoverensstemmelsesfejl. Allerede åbne sessioner anvender indstillingen, næste gang værktøjet kaldes; opgavesynlighed, overgange og menneskelige testregler følges fortsat.

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

Eksempler på Codex (brugerkonfiguration, projekt registreret automatisk) og Claude Code (projektfil) med projekt `C:\Projects\MyFirstProject` og arbejder `codex`:

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

I JSON er Windows-omvendt skråstreg fordoblet; et færdigt fragment fra grænsefladen gør dette automatisk. De genererede konfigurationer indeholder den faktiske `--db` for webserveren, så HTTP og MCP bruger den samme base, selvom klientmiljøet er anderledes.

**Ollama** leverer en lokal model, men erstatter ikke MCP-klienten. **Ollama via OpenCode**-profilen genererer en MCP-konfiguration til OpenCode; konfigurer OpenCode separat på Ollama-modellen. En anden MCP-kompatibel klient med Ollama er også velegnet.

## Trin 4: Tjek din forbindelse

1. Åbn arbejderkortet, og klik på **Kontroller igen**. **Servercheck** kører en ægte stdio eksekverbar med dette projekt, arbejder og base, kontrollerer håndtrykket, listen over værktøjer og at projektet og arbejderen er placeret.
2. Start eller genstart AI-klienten efter tilføjelse af konfigurationsfilen. Bed ham om at ringe til `get_my_board`.
3. **Kunde tilsluttet** vises i AgentBoard, og en ny session vises på listen. Kun dette bekræfter din klients forbindelse. Hvis der er værktøjer, men klienten ikke er tilsluttet, skal du kontrollere stien til programmet, navnet på konfigurationsfilen og dens JSON/TOML-syntaks.

Start din arbejdssession med `get_my_board`. AI modtager ikke opgaver fra `Backlog` og `Complete`, selvom den kender deres ID. AI kan ikke foregive at være et menneske og har ikke en generel "flytt hvor som helst"-kommando. Arbejdet med filsystemet uden for AgentBoard afhænger af den valgte klients muligheder og sandkasse.

Hvert værktøj erklærer `title` og alle fire MCP-adfærdstip (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) som eksplicitte booleske værdier. Læseværktøjer er skrivebeskyttede og idempotente; optageværktøjer tilføjer kun poster eller flytter opgaver fremad, så ingen er forstyrrende eller går ud over den lokale AgentBoard-base.

Officielle kundevejledninger: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
