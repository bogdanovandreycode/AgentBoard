# Agentstyret

![Forhåndsvisning AgentBoard](../../assets/social-preview.png)

**[Last ned for Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Dokumentasjonsside](https://bogdanovandreycode.github.io/AgentBoard/) · [Lisens MIT](../../LICENSE)**

AgentBoard er et lokalt oppgavebord der menneskelige og AI-arbeidere jobber med vanlige oppgaver, men har ulike rettigheter. Applikasjonen lanseres med én fil `agentboard.exe`, åpner nettgrensesnittet i nettleseren og gir arbeidere en separat MCP-server via `stdio`. Dataene forblir på datamaskinen din.

**[Start fra bunnen av](START_HERE.md)· [Jobber med oppgaver](TASKS.md)· [AI-tilkobling via MCP](WORKERS_MCP.md)· [Importer JSON](IMPORT.md)· [Innstillinger](SETTINGS.md)· [Problemløsning](TROUBLESHOOTING.md)**

## MVP-funksjoner

- Lokalt prosjekttavle med oppgavesøk, JSON-import, egendefinerte kolonner og egenskaper.
- MCP-tilgang for en spesifikk arbeider med begrensede AI-overganger og endelig menneskelig aksept av oppgaver.
- Generell oppgavehistorikk, sjekkinstruksjoner, artefakter, AI-kostnader og diagnostikk for arbeidertilkobling.
- Windows installasjonsprogram, ZIP og Scoop-manifest; webgrensesnittet er innebygd i den kjørbare filen.

AgentBoard er designet for en pålitelig lokal bruker. Applikasjonen er ikke vert for prosjekter i skyen og lanserer ikke AI-klienter selv; om nødvendig, koble en MCP-kompatibel klient til arbeideren.

## Om fem minutter

1. Last ned `agentboard-VERSION-windows-amd64-setup.exe`-installasjonsprogrammet fra [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Den vil foreslå en mappe (som standard `C:\AI\AgentBoard`) og legge den til `PATH`. Scoop og ZIP er også tilgjengelig.
2. Åpne PowerShell i prosjektmappen din, for eksempel `C:\Projects\MyApp`.
3. Kjør `agentboard init` (for ZIP: full bane til `agentboard.exe` og `init`).
4. Kjør `agentboard open`. `http://127.0.0.1:7337` åpnes.
5. Legg til en oppgave ved å bruke knappen **Ny oppgave**. For en AI-arbeider, åpne **Arbeidere → Legg til arbeider**, velg klientprofilen og kopier MCP-konfigurasjonen.

Hvis du ikke allerede har en prosjektmappe, lag en i Windows Utforsker. Et prosjekt kan være hvilken som helst mappe, selv uten Git og kode.

## Installasjon via Scoop

I PowerShell med [Scoop](https://scoop.sh/)] allerede installert etter utgivelsen:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

For utviklere er det [bygg fra kilden ](INSTALL.md). Utgivelsesarbeidsflyten oppretter et installasjonsprogram, ZIP og Scoop-manifest med SHA-256 fra samme artefakt. Oppdatering av versjonen installert via Scoop: `scoop update agentboard` etter å ha lagt til manifest i bøtta; detaljer - [utgivelsesforberedelse](SCOOP_RELEASE.md).

## Hvordan styret er bygget opp

`Backlog → Features → In progress → Testing → Verification → Complete`

En person kan flytte oppgaver rundt på brettet. AI kan bare flytte `Features → In progress → Testing → Verification`; endelig aksept i `Complete` utføres av et menneske. Brukerkolonner er ment for mennesker: oppgaven i dem forblir i `Backlog`-tilstanden for MCP. Testmoduser: AI, Human og Hybrid. Historikk, tester, artefakter og AI-kostnader er knyttet til oppgaven.

## Lag

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registrerer mappen og skriver kun `.agentboard/project.json` der. SQLite-produksjonsdataene er plassert i Windows-brukerkonfigurasjonskatalogen (`%AppData%\AgentBoard\agentboard.db`), utenfor prosjektet og utenfor Scoop-installasjonen. Fjerning eller oppdatering av pakken bør ikke fjerne disse dataene. Før du overfører til en annen datamaskin, lag en kopi av databasen mens AgentBoard er stoppet.

Arbeidere - logiske kontoer; AgentBoard selv kjører ikke Codex, Claude eller noen annen AI-klient. Klienten starter en lokal MCP-prosess for en spesifikk arbeider. MCP er applikasjonstillatelsesgrensen, og for filisolering, bruk AI-klientsandboksen.

## For utviklere

Stack: Go, SQLite, offisiell MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Bygg frontend først, deretter gå: `./scripts/build.ps1`. Nettfiler er inkludert i binæren via `go:embed`. Arkitekturen og API er beskrevet i [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Lisens

AgentBoard er et gratis og åpen kildekode-prosjekt under [lisens MIT](../../LICENSE). Kommersiell bruk, modifikasjon, fordeling og redistribuering er tillatt forutsatt at opphavsrettserklæringen og lisensen opprettholdes.
