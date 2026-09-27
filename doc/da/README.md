# Agentbestyrelse

[🌐 Languages](../LANGUAGES.md)

AgentBoard er et lokalt opgavebord, hvor menneskelige og AI-medarbejdere arbejder med fælles opgaver, men har forskellige rettigheder. Ansøgningen starter med én fil`agentboard.exe`, åbner webgrænsefladen i browseren og giver arbejderne en separat MCP-server via`stdio`. Dataene forbliver på din computer.

**[Start fra bunden](START_HERE.md)· [Arbejde med opgaver](TASKS.md)· [AI-forbindelse via MCP](WORKERS_MCP.md)· [Importer JSON](IMPORT.md)· [Indstillinger](SETTINGS.md)· [Problemløsning](TROUBLESHOOTING.md)**

## Om fem minutter

1. Download installationsprogrammet`agentboard-VERSION-windows-amd64-setup.exe`fra [Udgivelser](https://github.com/bogdanovandreycode/AgentBoard/releases). Det vil foreslå en mappe (standard`C:\AI\AgentBoard`) og vil tilføje det til`PATH`. Scoop og ZIP er også tilgængelige.
2. Åbn PowerShell i din projektmappe som f.eks`C:\Projects\MyApp`.
3. Udfør`agentboard init`(for ZIP: fuld sti til`agentboard.exe`Og`init`).
4. Udfør`agentboard open`. Vil åbne`http://127.0.0.1:7337`.
5. Tilføj en opgave ved hjælp af knappen **Ny opgave**. For en AI-medarbejder skal du åbne **Workers → Add worker**, vælge klientprofilen og kopiere MCP-konfigurationen.

Hvis du ikke allerede har en projektmappe, skal du oprette en i Windows Stifinder. Et projekt kan være en hvilken som helst mappe, selv uden Git og kode.

## Installation via Scoop

I PowerShell med [Scoop] allerede installeret](https://scoop.sh/)efter udgivelsen:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

For udviklere er der [build from source ](INSTALL.md). Udgivelsesworkflowet opretter et installationsprogram, ZIP og Scoop-manifest med SHA-256 fra den samme artefakt. Opdatering af versionen installeret via Scoop: `scoop update agentboard` efter tilføjelse af manifest til bøtten; detaljer - [udgivelsesforberedelse](SCOOP_RELEASE.md).

## Hvordan bestyrelsen er opbygget

`Backlog → Features → In progress → Testing → Verification → Complete`

En person kan flytte opgaver rundt på tavlen. AI kan kun flytte `Features → In progress → Testing → Verification`; endelig accept i `Complete` udføres af et menneske. Brugerkolonner er beregnet til mennesker: opgaven i dem forbliver i tilstanden `Backlog` for MCP. Testtilstande: AI, Human og Hybrid. Historik, test, artefakter og AI-omkostninger er knyttet til opgaven.

## Hold

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registrerer mappen og skriver kun `.agentboard/project.json` der. SQLite-produktionsdataene er placeret i Windows-brugerkonfigurationsbiblioteket (`%AppData%\AgentBoard\agentboard.db`), uden for projektet og uden for Scoop-installationen. Fjernelse eller opdatering af pakken bør ikke fjerne disse data. Før du overfører til en anden computer, skal du lave en kopi af databasen, mens AgentBoard er stoppet.

Arbejdere - logiske konti; AgentBoard selv kører ikke Codex, Claude eller nogen anden AI-klient. Klienten starter en lokal MCP-proces for en specifik arbejder. MCP er grænsen for applikationstilladelser, og til filisolering skal du bruge AI-klientsandkassen.

## For udviklere

Stack: Go, SQLite, officiel MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Byg frontend først, derefter Go: `./scripts/build.ps1`. Webfiler er inkluderet i binæren via `go:embed`. Arkitekturen og API er beskrevet i [doc/ARCHITECTURE.md](ARCHITECTURE.md).
