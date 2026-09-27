# Inställningar

Öppna **Inställningar** i sidomenyn för det valda projektet. Ändringar sparas med knappen **Spara inställningar** och lagras för projektet i AgentBoard-databasen.

## Språk och utseende

Fältet **Språk** öppnar en sökrullgardinslista. Som standard är alternativet **Följ system** valt, vilket tar webbläsarspråket. Tillgängliga språk från skärmdumpar: arabiska, portugisiska (Brasilien), förenklad kinesiska, tjeckiska, danska, holländska, engelska, finska, franska, tyska, italienska, japanska, koreanska, norska bokmål, polska, ryska, spanska, svenska, turkiska, ukrainska, vietnamesiska; dessutom vitryska, rumänska och bulgariska. **Följ system** tar webbläsarspråket. Huvudsignaturerna översätts manuellt, de återstående UI-raderna har en preliminär automatisk översättning. Innan den offentliggörs är det tillrådligt att korrekturläsa översättningar av modersmålstalare; klientnamn, kommandon, JSON-fält och användardata förblir utan översättning.

**Färgschema**: Mörk (originaltema), Ljus, Svart, Ubuntu och Windows. **Tidszon** styr visningen av datum; uppgifterna fortsätter att lagras i UTC. **Systemets tidszon** använder datorinställningar.

## Kolumner

De fyra stegen `Features`, `In progress`, `Testing`, `Verification` är fixerade i denna ordning: de kan inte bytas om eller raderas. Andra kolumner kan ordnas om med hjälp av de tillgängliga pilarna. Ange ett namn och klicka på **Lägg till kolumn** för att skapa en anpassad kolumn. Den är designad för mänskliga uppgifter: AI ser en uppgift som `Backlog` och tar inte emot den via MCP. Om du tar bort en kolumn medan du sparar överförs dess uppgifter till den vanliga `Backlog`.

## Webb och MCP

**Breduppdatering** ställer in uppdateringsfrekvensen för brädan och korten i sekunder (1–60). **Arbetaruppdatering** uppdaterar arbetarnas status (2–120 sekunder). Detta är webbgränssnittsundersökning, inte AI-utlösarfrekvens. Arbetare startar inte automatiskt.

Webbserveradressen ställs in vid start av CLI, till exempel `agentboard open --addr 127.0.0.1:7444`. För att ändra adressen måste servern startas om. Standard är `127.0.0.1:7337`. MCP fungerar genom ett separat lokalt kommando `agentboard mcp --project ... --worker ...` och är oberoende av webbporten. För en icke-standardbas, ange samma `--db` i alla kommandon. Kopiera konfigurationen för varje klient från arbetarkortet.
