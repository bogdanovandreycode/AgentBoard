# Projektledning och uppdragsutkast

Sidan **Alla projekt** visar den registrerade mappen, antalet publicerade
uppgifter och senaste aktivitet. Aktiviteten tar hänsyn till förändringar i uppgifter, historik och
utgiftsregister. Pennknappen låter dig byta namn på projektet eller välja
befintlig mapp i systemdialogrutan. Genom att ändra sökvägen uppdateras registreringen och markören;
Projekt-ID, arbetare, uppgifter och inställningar sparas, filer flyttas inte.
Mappar som redan är registrerade i ett annat projekt kan inte väljas.

## Globala inställningar

Öppna **Globala inställningar** från projektvalssidan eller -fliken
**Globala inställningar** i inställningarna på alla brädor. Dessa parametrar är gemensamma för alla projekt
en AgentBoard-bas:

- Automatisk lagring av utkast med timer är aktiverat som standard, var **10:e sekund**.
- **Spara när du lämnar fältet** är aktiverat som standard och kan konfigureras separat.
- Kontroller av **MCP-projektkontext** är aktiverade och inaktiverade globalt. Deras område
  åtgärder och en separat CLI-flagga beskrivs i avsnittet [Arbetare och MCP](WORKERS_MCP.md).

## Utkast till livscykel

Dialogrutan Ny uppgift sparar oavslutat arbete som ett utkast. Statusrad
skiljer på osparade ändringar, lagring pågår, sparat utkast och fel
bevarande. Utkastet visas inte på tavlan och är inte tillgängligt för MCP-arbetare förrän
publiceras.

Innan du publicerar, använd knappen **Skapa uppgift** för att fylla i **Titel, Status,
Ansvarig och beskrivning**. Innan detta sparar **Spara utkast** ett utkast utan att skapa
uppgifter. Ett tomt formulär kasseras som standard när det stängs. Om rensas
befintligt utkast till tomt tillstånd, raderas det nästa gång det sparas.
Om du stänger en icke-tom dialogruta sparas den också.

Knappen **Utkast** på tavlan låter dig fortsätta eller ta bort sparade utkast.
Flera utkast stöds. Efter en omstart öppnas
utkast. Utkast lagras i SQLite; Med varje ändring uppdateras den lokala
kopiera i webbläsaren, så osparade redigeringar överlever omstarter eller
tillfälligt serverfel. Om webbläsarlagring inte är tillgängligt, varnar en dialogruta
att utkastet måste sparas på servern innan stängning.

Att skriva utkast använder revisioner för att förhindra ett fönster från att skriva över arbete
en annan. Förlorat sparsvar är konsekvent om det redan finns ett på servern
identiska uppgifter. Om innehållet är annorlunda visas en konflikt och
**Återställ som ett nytt utkast** sparar lokalt innehåll separat. Publicering
atomär och idempotent: återpublicera samma utkast efter att det har gått förlorat
svar returnerar samma uppgift snarare än att skapa en kopia.
