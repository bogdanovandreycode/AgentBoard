# Begin hier

## Wat is AgentBoard

Stel je een gewoon bord voor met taakkaarten. Je maakt taken aan, wijst een verantwoordelijke aan en houdt toezicht op de werkzaamheden. AI-werknemers ontvangen alleen de hun toegewezen taken via MCP en rapporteren de voortgang. Jij bepaalt wanneer de taak eindelijk klaar is.

**Project** - een map op de computer en een apart bord. **Taak** - een kaart met een beschrijving, verantwoordelijke persoon en fase. **Worker** is de logische naam van de AI-client. **MCP** is de manier waarop de AI-client verbinding maakt met het AgentBoard. **Scoop** is een software-installatiemanager voor Windows.

Geen code vereist. Je hebt Windows, een browser, PowerShell en, om AI te laten werken, een geïnstalleerde AI-client nodig met ondersteuning voor lokale MCP-servers.

## Eerste lancering

1. Installeer de applicatie volgens [instructies](INSTALL.md).
2. Maak een projectmap in Verkenner, bijvoorbeeld `C:\Projects\MyFirstProject`.
3. Open deze map in Verkenner. Klik in de adresbalk, typ `powershell` en druk op Enter.
4. In het geopende venster doet u het volgende:

```powershell
agentboard init
agentboard open
```

5. Er wordt een browser geopend met het adres `http://127.0.0.1:7337`. Laat het PowerShell-venster open terwijl u het whiteboard gebruikt. Als u het venster sluit, stopt de lokale server, maar blijven de taken bestaan.

Als de opdracht `agentboard` niet wordt gevonden, sluit u PowerShell en opent u deze opnieuw nadat u Scoop hebt geïnstalleerd. Wanneer u installeert vanuit een ZIP, gebruik dan het volledige pad naar `agentboard.exe`.

## Eerste taak

Klik op **Nieuwe taak**, vul **Titel** in, indien nodig **Beschrijving** en vervolgens **Opslaan**. Een nieuwe taak in `Backlog` is beschikbaar voor mensen. Om de AI te laten werken, wijst u een werknemer toe en draagt ​​u de taak over naar `Features`. Stapsgewijze beschrijving van de velden - [TASKS.md](TASKS.md).

## Eerste werknemer

Open **Werknemers → Werknemer toevoegen**. Selecteer de AI-client die u gebruikt (bijvoorbeeld Codex of Claude Code), controleer de naam en korte ID `Slug`, klik op **Opslaan**. Open de aangemaakte werker: er is een MCP-configuratie, een kopieerknop en een servercontrole. Kopieer de configuratie naar de AI-client volgens [WORKERS_MCP.md](WORKERS_MCP.md). Eenmaal verbonden, vraagt ​​u de klant om `get_my_board` te bellen.

## Wat moet je nu lezen

- [Taken, tests, verhalen en kolommen](TASKS.md)
- [Werknemers verbinden en MCP](WORKERS_MCP.md) controleren
- [Importeer taken uit JSON](IMPORT.md)
- [Taal, thema, tijdzone en update](SETTINGS.md)
- [Typische problemen](TROUBLESHOOTING.md)
