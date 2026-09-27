# Installation og start på Windows

## Installationsprogram (anbefales)

Download `agentboard-VERSION-windows-amd64-setup.exe` fra siden [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Installationsguiden vil foreslå en mappe; standarden er `C:\AI\AgentBoard`. Den vil kopiere `agentboard.exe` og dokumentation, oprette en genvej og tilføje den valgte mappe til systemet `PATH`. Efter installationen skal du åbne en ny terminal, så kommandoen `agentboard` bliver tilgængelig.

Kør i din projektmappe:

```powershell
agentboard init
agentboard open
```

Interfacet er indbygget i `agentboard.exe`; Ingen separat installation af Go eller Node.js er nødvendig. Dataene gemmes i `%AppData%\AgentBoard` og bevares, når programmet opdateres eller afinstalleres. Afinstallation via "Installerede applikationer" fjerner genveje og indgangen fra `PATH`.

## Installerer Scoop

Hvis Scoop ikke allerede er installeret, skal du åbne PowerShell som din almindelige bruger og følge [officielle Scoop](https://scoop.sh/) instruktioner. Hvis der er begrænsninger på din virksomhedscomputer, skal du kontakte din administrator. AgentBoard kan også startes fra en ZIP uden Scoop.

Alternativt kan du installere AgentBoard via Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Tilgængeligheden af ​​manifest i GitHub Release kan kontrolleres på [side Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Efter tilføjelse af manifest til Scoop, kan bucket installeres efter bucket navn og opdateres med `scoop update agentboard` kommandoen.

## ZIP uden Scoop

Download `agentboard-VERSION-windows-amd64.zip` fra Releases, pak for eksempel ud i `C:\Tools\AgentBoard`. I PowerShell i projektmappen:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

For en AI-klient skal du angive den fulde sti til `agentboard.exe` i dens MCP-konfiguration, hvis programmet ikke er placeret i `PATH`.

## Byg fra kilde

Installer Go-versionen fra `go.mod` og Node.js 22 eller nyere. I PowerShell i roden af ​​depotet:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` samler frontenden til `internal/webui/dist`, kører Go-tests og samler en `agentboard.exe`. Rækkefølgen er vigtig: grænsefladen er indbygget i det binære, når Go er bygget. Hvis den gamle `agentboard.exe open` kører, skal du stoppe den før genopbygning (Ctrl+C), ellers vil Windows ikke tillade dig at erstatte filen.

## Hvor er dataene

- `%AppData%\AgentBoard\agentboard.db` - opgaver, projekter, arbejdere og indstillinger. Du kan angive en anden fil med `--db`-flaget, men `init`, `open`/`serve` og `mcp` skal have **samme sti**.
- `<ваш проект>\.agentboard\project.json` — projektidentifikator. Denne fil indeholder ikke opgaver.
- Serveren lytter kun til `127.0.0.1:7337` som standard. Angiv en anden adresse `--addr` før projektstien: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` åbner projektsiden og genbruger en allerede kørende server på den adresse. Hvis flere projekter er åbne i en browser, skal du vælge dem gennem listen til venstre.
