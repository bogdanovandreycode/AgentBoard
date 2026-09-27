# Installatie en start op Windows

## Installatieprogramma (aanbevolen)

Download `agentboard-VERSION-windows-amd64-setup.exe` van de pagina [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). De installatiewizard zal een map voorstellen; de standaardwaarde is `C:\AI\AgentBoard`. Het kopieert `agentboard.exe` en documentatie, maakt een snelkoppeling en voegt de geselecteerde map toe aan het systeem `PATH`. Open na de installatie een nieuwe terminal zodat het commando `agentboard` beschikbaar komt.

Voer in uw projectmap het volgende uit:

```powershell
agentboard init
agentboard open
```

De interface is ingebouwd in de `agentboard.exe`; Er is geen aparte installatie van Go of Node.js nodig. De gegevens worden opgeslagen in `%AppData%\AgentBoard` en blijven behouden wanneer het programma wordt bijgewerkt of verwijderd. Als u de installatie ongedaan maakt via “Geïnstalleerde applicaties” worden snelkoppelingen en de vermelding uit `PATH` verwijderd.

## Scoop installeren

Als Scoop nog niet is geïnstalleerd, opent u PowerShell als uw gewone gebruiker en volgt u de [officiële Scoop](https://scoop.sh/)-instructies. Als er beperkingen gelden op uw bedrijfscomputer, neem dan contact op met uw beheerder; AgentBoard kan ook zonder Scoop vanuit een ZIP worden gestart.

U kunt AgentBoard ook via Scoop installeren:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

De beschikbaarheid van het manifest in GitHub Release kan worden gecontroleerd op [pagina Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Nadat het manifest aan Scoop is toegevoegd, kan de bucket op bucketnaam worden geïnstalleerd en worden bijgewerkt met de opdracht `scoop update agentboard`.

## ZIP zonder schep

Download `agentboard-VERSION-windows-amd64.zip` uit Releases en pak het bijvoorbeeld uit in `C:\Tools\AgentBoard`. In PowerShell in de projectmap:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Voor een AI-client geeft u het volledige pad naar `agentboard.exe` op in de MCP-configuratie als het programma zich niet in `PATH` bevindt.

## Bouw vanuit de bron

Installeer de Go-versie van `go.mod` en Node.js 22 of hoger. In PowerShell in de hoofdmap van de repository:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` assembleert de frontend tot `internal/webui/dist`, voert Go-tests uit en assembleert één `agentboard.exe`. De volgorde is belangrijk: de interface wordt in het binaire bestand ingebouwd wanneer Go wordt gebouwd. Als de oude `agentboard.exe open` actief is, stop deze dan voordat u hem opnieuw opbouwt (Ctrl+C), anders staat Windows u niet toe het bestand te vervangen.

## Waar zijn de gegevens

- `%AppData%\AgentBoard\agentboard.db` - taken, projecten, werknemers en instellingen. U kunt een ander bestand opgeven met de vlag `--db`, maar `init`, `open`/`serve` en `mcp` moeten **hetzelfde pad** hebben.
- `<ваш проект>\.agentboard\project.json` — project-ID. Dit bestand bevat geen taken.
- De server luistert standaard alleen naar `127.0.0.1:7337`. Geef vóór het projectpad nog een adres `--addr` op: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` opent de projectpagina en gebruikt een reeds actieve server op dat adres opnieuw. Als er meerdere projecten in één browser geopend zijn, selecteert u deze via de lijst aan de linkerkant.
