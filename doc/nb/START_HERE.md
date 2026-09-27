# Start her

## Hva er AgentBoard

Se for deg et vanlig brett med oppgavekort. Du lager oppgaver, tildeler en ansvarlig og fører tilsyn med arbeidet. AI-arbeidere mottar kun sine tildelte oppgaver via MCP og rapporterer fremgang. Du bestemmer når oppgaven endelig er klar.

**Prosjekt** - en mappe på datamaskinen og et eget styre. **Oppgave** - et kort med beskrivelse, ansvarlig person og scene. **Worker** er det logiske navnet på AI-klienten. **MCP** er måten AI-klienten kobler til AgentBoard. **Scoop** er en programvareinstallasjonsbehandling for Windows.

Ingen kode kreves. Du trenger Windows, en nettleser, PowerShell og, for at AI skal fungere, en installert AI-klient med støtte for lokale MCP-servere.

## Første lansering

1. Installer applikasjonen i henhold til [instruksjoner ](INSTALL.md).
2. Opprett en prosjektmappe i Utforsker, for eksempel `C:\Projects\MyFirstProject`.
3. Åpne denne mappen i Utforsker. Klikk i adressefeltet, skriv inn `powershell` og trykk Enter.
4. I vinduet som åpnes gjør du:

```powershell
agentboard init
agentboard open
```

5. En nettleser åpnes med adressen `http://127.0.0.1:7337`. La PowerShell-vinduet være åpent mens du bruker tavlen. Å lukke vinduet vil stoppe den lokale serveren, men oppgavene vil forbli.

Hvis kommandoen `agentboard` ikke blir funnet, lukker du PowerShell og åpner den på nytt etter at du har installert Scoop. Når du installerer fra en ZIP, bruk hele banen til `agentboard.exe`.

## Første oppgave

Klikk på **Ny oppgave**, fyll inn **Tittel**, om nødvendig **Beskrivelse**, deretter **Lagre**. En ny oppgave i `Backlog` er tilgjengelig for mennesker. For at AI skal begynne å jobbe, tilordne en arbeider og overføre oppgaven til `Features`. Trinn-for-trinn beskrivelse av feltene - [TASKS.md](TASKS.md).

## Første arbeider

Åpne **Arbeidere → Legg til arbeider**. Velg AI-klienten du bruker (for eksempel Codex eller Claude Code), sjekk navnet og kort ID `Slug`, klikk **Lagre**. Åpne den opprettede arbeideren: det er en MCP-konfigurasjon, en kopiknapp og en serversjekk. Kopier konfigurasjonen til AI-klienten i henhold til [WORKERS_MCP.md](WORKERS_MCP.md). Når du er tilkoblet, ber du klienten ringe `get_my_board`.

## Hva du skal lese videre

- [Oppgaver, tester, historier og spalter](TASKS.md)
- [Koble til arbeidere og sjekke MCP](WORKERS_MCP.md)
- [Importer oppgaver fra JSON](IMPORT.md)
- [Språk, tema, tidssone og oppdatering](SETTINGS.md)
- [Typiske problemer](TROUBLESHOOTING.md)
