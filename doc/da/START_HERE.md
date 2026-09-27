# Start her

## Hvad er AgentBoard

Forestil dig en almindelig tavle med opgavekort. Du opretter opgaver, tildeler en ansvarlig og fører tilsyn med arbejdet. AI-medarbejdere modtager kun deres tildelte opgaver via MCP og rapporterer fremskridt. Du bestemmer selv, hvornår opgaven endelig er klar.

**Projekt** - en mappe på computeren og et separat bord. **Opgave** - et kort med beskrivelse, ansvarlig person og scene. **Worker** er det logiske navn på AI-klienten. **MCP** er den måde, AI-klienten forbinder til AgentBoard. **Scoop** er en softwareinstallationsmanager til Windows.

Ingen kode påkrævet. Du skal bruge Windows, en browser, PowerShell og, for at AI skal fungere, en installeret AI-klient med understøttelse af lokale MCP-servere.

## Første lancering

1. Installer applikationen i henhold til [instruktioner](INSTALL.md).
2. Opret en projektmappe i Stifinder, for eksempel `C:\Projects\MyFirstProject`.
3. Åbn denne mappe i Stifinder. Klik i adresselinjen, skriv `powershell` og tryk på Enter.
4. I det vindue, der åbnes, gør du:

```powershell
agentboard init
agentboard open
```

5. En browser åbnes med adressen `http://127.0.0.1:7337`. Lad PowerShell-vinduet være åbent, mens du bruger tavlen. Lukning af vinduet vil stoppe den lokale server, men opgaverne forbliver.

Hvis kommandoen `agentboard` ikke findes, skal du lukke PowerShell og åbne den igen efter installation af Scoop. Når du installerer fra en ZIP, skal du bruge den fulde sti til `agentboard.exe`.

## Første opgave

Klik på **Ny opgave**, udfyld **Titel**, om nødvendigt **Beskrivelse** og derefter **Gem**. En ny opgave i `Backlog` er tilgængelig for mennesker. For at AI kan begynde at arbejde, skal du tildele en arbejder og overføre opgaven til `Features`. Trin-for-trin beskrivelse af felterne - [TASKS.md](TASKS.md).

## Første arbejder

Åbn **Medarbejdere → Tilføj arbejder**. Vælg den AI-klient, du bruger (for eksempel Codex eller Claude Code), tjek navnet og det korte ID `Slug`, klik på **Gem**. Åbn den oprettede arbejder: der er en MCP-konfiguration, en kopiknap og en serverkontrol. Kopier konfigurationen til AI-klienten i henhold til [WORKERS_MCP.md](WORKERS_MCP.md). Når du er tilsluttet, skal du bede klienten om at ringe til `get_my_board`.

## Hvad skal du læse næste gang

- [Opgaver, test, historier og spalter](TASKS.md)
- [Forbinder arbejdere og kontrollerer MCP](WORKERS_MCP.md)
- [Importér opgaver fra JSON](IMPORT.md)
- [Sprog, tema, tidszone og opdatering](SETTINGS.md)
- [Typiske problemer](TROUBLESHOOTING.md)
