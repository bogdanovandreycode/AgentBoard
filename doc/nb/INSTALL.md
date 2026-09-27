# Installasjon og start på Windows

## Installer (anbefalt)

Last ned `agentboard-VERSION-windows-amd64-setup.exe` fra siden [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Installasjonsveiviseren vil foreslå en mappe; standard er `C:\AI\AgentBoard`. Den vil kopiere `agentboard.exe` og dokumentasjon, lage en snarvei og legge til den valgte mappen til systemet `PATH`. Etter installasjonen åpner du en ny terminal slik at kommandoen `agentboard` blir tilgjengelig.

Kjør i prosjektmappen din:

```powershell
agentboard init
agentboard open
```

Grensesnittet er innebygd i `agentboard.exe`; Ingen separat installasjon av Go eller Node.js er nødvendig. Dataene lagres i `%AppData%\AgentBoard` og beholdes når programmet oppdateres eller avinstalleres. Avinstallering via "Installerte applikasjoner" fjerner snarveier og oppføringen fra `PATH`.

## Installerer Scoop

Hvis Scoop ikke allerede er installert, åpne PowerShell som din vanlige bruker og følg [offisielle Scoop](https://scoop.sh/) instruksjoner. Hvis det er begrensninger på bedriftsdatamaskinen din, kontakt administratoren. AgentBoard kan også startes fra en ZIP uten Scoop.

Alternativt kan du installere AgentBoard via Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Tilgjengeligheten av manifest i GitHub-utgivelsen kan sjekkes på [side Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Etter å ha lagt til manifest til Scoop, kan bøtten installeres etter bøttenavn og oppdateres med `scoop update agentboard`-kommandoen.

## ZIP uten Scoop

Last ned `agentboard-VERSION-windows-amd64.zip` fra utgivelser, pakk ut for eksempel i `C:\Tools\AgentBoard`. I PowerShell i prosjektmappen:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

For en AI-klient, spesifiser hele banen til `agentboard.exe` i MCP-konfigurasjonen hvis programmet ikke er plassert i `PATH`.

## Bygg fra kilden

Installer Go-versjonen fra `go.mod` og Node.js 22 eller nyere. I PowerShell ved roten av depotet:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` monterer frontenden til `internal/webui/dist`, kjører Go-tester og monterer én `agentboard.exe`. Rekkefølgen er viktig: grensesnittet er innebygd i binæren når Go bygges. Hvis den gamle `agentboard.exe open` kjører, stopp den før gjenoppbygging (Ctrl+C), ellers vil ikke Windows tillate deg å erstatte filen.

## Hvor er dataene

- `%AppData%\AgentBoard\agentboard.db` - oppgaver, prosjekter, arbeidere og innstillinger. Du kan spesifisere en annen fil med `--db`-flagget, men `init`, `open`/`serve` og `mcp` må ha **samme bane**.
- `<ваш проект>\.agentboard\project.json` — project identifier. This file does not contain tasks.
- Serveren lytter kun til `127.0.0.1:7337` som standard. Spesifiser en annen adresse `--addr` før prosjektbanen: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` åpner prosjektsiden og gjenbruker en server som allerede kjører på den adressen. Hvis flere prosjekter er åpne i en nettleser, velg dem gjennom listen til venstre.
