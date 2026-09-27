# Instellingen

Open **Instellingen** in het zijmenu van het geselecteerde project. Wijzigingen worden opgeslagen met de knop **Instellingen opslaan** en voor het project opgeslagen in de AgentBoard-database.

## Taal en uiterlijk

Het veld **Taal** opent een vervolgkeuzelijst voor zoeken. Standaard is de optie **Systeem volgen** geselecteerd, waarbij de browsertaal wordt gebruikt. Beschikbare talen uit screenshots: Arabisch, Portugees (Brazilië), Vereenvoudigd Chinees, Tsjechisch, Deens, Nederlands, Engels, Fins, Frans, Duits, Italiaans, Japans, Koreaans, Noors Bokmål, Pools, Russisch, Spaans, Zweeds, Turks, Oekraïens, Vietnamees; daarnaast Wit-Russisch, Roemeens en Bulgaars. **Volgsysteem** gebruikt de browsertaal. De belangrijkste handtekeningen worden handmatig vertaald, de overige UI-regels hebben een voorlopige automatische vertaling. Voordat ze openbaar worden gemaakt, is het raadzaam vertalingen door moedertaalsprekers te proeflezen; clientnamen, opdrachten, JSON-velden en gebruikersgegevens blijven zonder vertaling.

**Kleurenschema**: Donker (origineel thema), Licht, Zwart, Ubuntu en Windows. **Tijdzone** regelt de weergave van datums; de gegevens blijven opgeslagen in UTC. **Systeemtijdzone** gebruikt computerinstellingen.

## Kolommen

De vier fasen `Features`, `In progress`, `Testing`, `Verification` zijn in deze volgorde vastgelegd: ze kunnen niet worden hernoemd of verwijderd. Andere kolommen kunnen opnieuw worden gerangschikt met behulp van de beschikbare pijlen. Voer een naam in en klik op **Kolom toevoegen** om een ​​aangepaste kolom te maken. Het is ontworpen voor menselijke taken: AI ziet een taak als `Backlog` en ontvangt deze niet via MCP. Als u een kolom verwijdert tijdens het opslaan, worden de taken ervan overgebracht naar de reguliere `Backlog`.

## Web en MCP

**Bordvernieuwing** stelt de vernieuwingssnelheid van het bord en de kaarten in seconden in (1–60). **Werknemer vernieuwen** werkt de status van werknemers bij (2-120 seconden). Dit is webinterface-polling, geen AI-triggerfrequentie. Werknemers beginnen niet automatisch.

Het webserveradres wordt ingesteld bij het starten van de CLI, bijvoorbeeld `agentboard open --addr 127.0.0.1:7444`. Om het adres te wijzigen, moet de server opnieuw worden opgestart. De standaardwaarde is `127.0.0.1:7337`. MCP werkt via een afzonderlijk lokaal commando `agentboard mcp --project ... --worker ...` en is onafhankelijk van de webpoort. Voor een niet-standaard basis specificeert u in alle opdrachten dezelfde `--db`. Kopieer de configuratie van elke client van de werknemerskaart.
