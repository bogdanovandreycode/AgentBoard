# Projectbeheer en taakconcepten

De pagina **Alle projecten** toont de geregistreerde map en het aantal gepubliceerde mappen
taken en laatste activiteit. Activiteit houdt rekening met veranderingen in taken, geschiedenis en
onkostenregistratie. Met de potloodknop kunt u het project hernoemen of selecteren
bestaande map in het systeemdialoogvenster. Als u het pad wijzigt, worden de registratie en markering bijgewerkt;
Project-ID's, werknemers, taken en instellingen worden opgeslagen, bestanden worden niet verplaatst.
Mappen die al bij een ander project zijn geregistreerd, kunnen niet worden geselecteerd.

## Algemene instellingen

Open **Algemene instellingen** vanaf de projectselectiepagina of het tabblad
**Algemene instellingen** in de instellingen van elk bord. Deze parameters zijn gemeenschappelijk voor alle projecten
één AgentBoard-basis:

- Het automatisch opslaan van concepten via de timer is standaard ingeschakeld, elke **10 seconden**.
- **Opslaan bij het verlaten van het veld** is standaard ingeschakeld en kan afzonderlijk worden geconfigureerd.
- **MCP-projectcontext**-controles worden wereldwijd in- en uitgeschakeld. Hun gebied
  acties en een aparte CLI-vlag worden beschreven in de sectie [Werknemers en MCP](WORKERS_MCP.md).

## Ontwerp levenscyclus

In het dialoogvenster Nieuwe taak wordt onvoltooid werk opgeslagen als concept. Statusregel
maakt onderscheid tussen niet-opgeslagen wijzigingen, bezig met opslaan, opgeslagen concept en fouten
behoud. Het concept verschijnt niet op het bord en is pas beschikbaar voor MCP-werknemers
gepubliceerd.

Gebruik vóór publicatie de knop **Taak maken** om **Titel, Staat,
Verantwoordelijk en beschrijving**. Voordien slaat **Concept opslaan** een concept op zonder dat u het hoeft te maken
taken. Een leeg formulier wordt bij het sluiten standaard weggegooid. Indien gewist
bestaande concept naar een lege staat, wordt het verwijderd wanneer het de volgende keer wordt opgeslagen.
Als u een niet-lege dialoog sluit, wordt deze ook opgeslagen.

Met de knop **Concepten** op het bord kunt u doorgaan met opgeslagen concepten of deze verwijderen.
Meerdere concepten ondersteund. Na een herstart verschijnt het bestand open
ontwerp. Concepten worden opgeslagen in SQLite; Bij elke wijziging wordt de lokale bijgewerkt
kopiëren in de browser, zodat niet-opgeslagen bewerkingen het opnieuw opstarten overleven of
tijdelijke serverstoring. Als er geen browseropslag beschikbaar is, wordt er een dialoogvenster gewaarschuwd
dat het concept vóór het sluiten op de server moet worden opgeslagen.

Bij het schrijven van concepten worden revisies gebruikt om te voorkomen dat één venster het werk overschrijft
een ander. Een verloren opslagreactie is consistent als er al een op de server staat
identieke gegevens. Als de inhoud anders is, wordt er een conflict getoond en
**Herstellen als een nieuw concept** slaat lokale inhoud afzonderlijk op. Publicatie
atomair en idempotent: hetzelfde concept opnieuw publiceren nadat het verloren is gegaan
response retourneert dezelfde taak in plaats van een kopie te maken.
