# Agentstyret

[🌐 Languages](../LANGUAGES.md)

AgentBoard er et lokalt oppgavebord der menneskelige og AI-arbeidere jobber med vanlige oppgaver, men har ulike rettigheter. Applikasjonen starter med én fil`agentboard.exe`, åpner nettgrensesnittet i nettleseren og gir arbeiderne en egen MCP-server via`stdio`. Dataene forblir på datamaskinen din.

**[Start fra bunnen av](START_HERE.md)· [Jobber med oppgaver](TASKS.md)· [AI-tilkobling via MCP](WORKERS_MCP.md)· [Importer JSON](IMPORT.md)· [Innstillinger](SETTINGS.md)· [Problemløsning](TROUBLESHOOTING.md)**

## Om fem minutter

1. Last ned installasjonsprogrammet`agentboard-VERSION-windows-amd64-setup.exe`fra [Utgivelser](https://github.com/bogdanovandreycode/AgentBoard/releases). Det vil foreslå en mappe (standard`C:\AI\AgentBoard`) og vil legge den til`PATH`. Scoop og ZIP er også tilgjengelig.
2. Åpne PowerShell i prosjektmappen som`C:\Projects\MyApp`.
3. Utfør`agentboard init`(for ZIP: full bane til`agentboard.exe`Og`init`).
4. Utfør`agentboard open`. Vil åpne`http://127.0.0.1:7337`.
5. Legg til en oppgave ved å bruke knappen **Ny oppgave**. For en AI-arbeider, åpne **Arbeidere → Legg til arbeider**, velg klientprofilen og kopier MCP-konfigurasjonen.

Hvis du ikke allerede har en prosjektmappe, lag en i Windows Utforsker. Et prosjekt kan være hvilken som helst mappe, selv uten Git og kode.

## Installasjon via Scoop

I PowerShell med [Scoop] allerede installert](https://scoop.sh/)etter utgivelsen:

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
