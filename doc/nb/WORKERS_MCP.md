# Arbeidere og MCP-tilkobling

## Trinn 1. Opprett en arbeider

I AgentBoard åpner du **Arbeidere → Legg til arbeider**. Velg en klientprofil: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot eller en annen MCP-klient. Profilen fyller ut typiske funksjoner (kode, tester, Git); du kan endre dem. Navnet er synlig på tavlen, og `Slug` er en kort identifikator uten mellomrom for MCP-kommandoen. Klikk på **Lagre**.

Én arbeider tilsvarer én AI-personlighet. Lag forskjellige arbeidere for forskjellige kunder eller team. Kompetanseprofilen beskriver spesialiseringen, men utvider ikke AIs rettigheter til oppgavestadier.

## Trinn 2. Kopier konfigurasjonen

Åpne den opprettede arbeideren. **MCP-diagnostikkblokken viser konfigurasjonsfragmentet og filen der det skal legges til. Klikk på **Kopier MCP-konfigurasjon**. Hvis filen allerede eksisterer, legg til den foreslåtte serveren til det eksisterende `mcpServers`/`servers`/`mcp`-objektet uten å slette de andre serverne.

Hovedkommandoen ser slik ut:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` til den opprettede arbeideren. MCP-klienten kjører denne kommandoen selv når den trenger verktøy.

### MCP fungerer uten nettgrensesnitt

`agentboard mcp` er en egen stdio-prosess som selv åpner AgentBoard-databasen. Den trenger ikke `agentboard open`/`serve`, den avhenger ikke av mappen som nettgrensesnittet startes fra, og hvor mange forekomster av nettgrensesnittet som kjører. Flere agenter kan kjøres samtidig, hver i sitt eget prosjekt: hver klient kjører sin egen MCP-prosess, og de bruker alle sikkert den samme SQLite-databasen.

### Hvordan et prosjekt velges

MCP-serveren kjører i hvilken som helst mappe. Første gang du ringer verktøyet, er det bundet til prosjektet:

1. `--project PATH`, hvis spesifisert (enhver mappe i det registrerte prosjektet er egnet);
2. ellers klientens arbeidsmapper (røtter), hvis klienten rapporterer dem;
3. ellers mappen der klienten startet MCP-serveren.

`get_my_board` returnerer det valgte prosjektet og arbeideren til `session` sammen med `binding_source`. Hvis klientens arbeidsmapper endres, utfører neste oppkall til verktøyet bindingen på nytt. Uten `--worker` brukes den eneste aktiverte arbeideren i prosjektet.

Hvis binding ikke er mulig, returnerer hvert verktøy en feil som agenten må sende til brukeren; oppgavedata leses ikke i dette tilfellet:

| Kode | Betydning |
| --- | --- |
| `NO_PROJECTS` | Det er ingen prosjekter i databasen. Kjør `agentboard init` i prosjektmappen og opprett en arbeider. Ingen omstart nødvendig. |
| `PROJECT_REQUIRED` | Arbeidsmappen er utenfor registrerte prosjekter. Meldingen viser prosjekter. |
| `PROJECT_NOT_FOUND` | `--project` er ikke inne i et registrert prosjekt. |
| `PROJECT_AMBIGUOUS` | En klients arbeidsmapper spenner over flere prosjekter. Legg til `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Opprett, aktiver eller velg en arbeider via `--worker`. |

Konfigurasjonene for **Codex** og **Gemini CLI** er tilpasset og delt på tvers av alle mapper, slik at de ikke har `--project`: Opprett en arbeider med samme slug i hvert prosjekt der denne klienten kjører. Prosjektnivåfiler (Claude Code, Cursor, OpenCode, VS Code) inneholder `--project` av prosjektet deres.

### Eksplisitt design og kontekstkontroll

Hvis `--project` er spesifisert og klientens oppstartsmappe eller arbeidsmappe er i et annet prosjekt, ber verktøyene om bekreftelse via MCP-elicitation. Klienter uten skjemafremkalling mottar `PROJECT_CONTEXT_MISMATCH` før de får tilgang til oppgaver. Korriger konfigurasjonen, fjern `--project` for automatisk gjenkjenning, eller bekreft eksplisitt at du har tenkt å jobbe med et annet prosjekt med `--allow-project-mismatch`-flagget. Endring av arbeidsmapper tilbakestiller tidligere bekreftelser. Automatisk definerte prosjekter krever ikke bekreftelse.

**Globale innstillinger → MCP-prosjektkontekst** kontrollerer disse sjekkene for alle prosjekter og arbeidere i samme database. Aktivert som standard. Hvis deaktivert, bruker MCP det spesifiserte prosjektet uten bekreftelser eller mismatchfeil. Allerede åpne økter bruker innstillingen neste gang verktøyet kalles; oppgavesynlighet, overganger og regler for menneskelig testing følges fortsatt.

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

Eksempler for Codex (brukerkonfigurasjon, prosjekt oppdaget automatisk) og Claude Code (prosjektfil) med prosjekt `C:\Projects\MyFirstProject` og arbeider `codex`:

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

I JSON er Windows-omvendt skråstrek doblet; et ferdig fragment fra grensesnittet gjør dette automatisk. De genererte konfigurasjonene inneholder den faktiske `--db` til webserveren, så HTTP og MCP bruker samme base, selv om klientmiljøet er forskjellig.

**Ollama** gir en lokal modell, men erstatter ikke MCP-klienten. **Ollama via OpenCode**-profilen genererer en MCP-konfigurasjon for OpenCode; konfigurere OpenCode separat på Ollama-modellen. En annen MCP-kompatibel klient med Ollama er også egnet.

## Trinn 4: Sjekk tilkoblingen

1. Åpne arbeiderkortet og klikk på **Sjekk igjen**. **Server check** kjører en ekte stdio kjørbar med dette prosjektet, arbeider og base, sjekker håndtrykket, listen over verktøy og at prosjektet og arbeideren er lokalisert.
2. Start eller restart AI-klienten etter å ha lagt til konfigurasjonsfilen. Be ham ringe `get_my_board`.
3. **Klient tilkoblet** vises i AgentBoard og en ny økt vil vises i listen. Bare dette bekrefter tilkoblingen til klienten din. Hvis det finnes verktøy, men klienten ikke er tilkoblet, sjekk banen til programmet, navnet på konfigurasjonsfilen og dens JSON/TOML-syntaks.

Start arbeidsøkten med `get_my_board`. AI mottar ikke oppgaver fra `Backlog` og `Complete`, selv om den kjenner deres ID. AI kan ikke late som om det er et menneske og har ikke en generell "flytt hvor som helst"-kommando. Arbeid med filsystemet utenfor AgentBoard avhenger av egenskapene og sandkassen til den valgte klienten.

Hvert verktøy erklærer `title` og alle fire MCP-adferdstips (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) som eksplisitte boolske verdier. Leseverktøy er skrivebeskyttet og idempotente; opptaksverktøy legger bare til poster eller flytter oppgaver fremover, så ingen er forstyrrende eller går utover den lokale AgentBoard-basen.

Offisielle kundeinstruksjoner: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
