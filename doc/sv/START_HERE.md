# Börja här

## Vad är AgentBoard

Föreställ dig en vanlig tavla med uppgiftskort. Du skapar uppgifter, utser någon ansvarig och övervakar arbetet. AI-arbetare får endast sina tilldelade uppgifter via MCP och rapporterar framsteg. Du bestämmer när uppgiften äntligen är klar.

**Projekt** - en mapp på datorn och ett separat kort. **Uppgift** - ett kort med beskrivning, ansvarig person och scen. **Worker** är det logiska namnet på AI-klienten. **MCP** är sättet som AI-klienten ansluter till AgentBoard. **Scoop** är en mjukvaruinstallationshanterare för Windows.

Ingen kod krävs. Du behöver Windows, en webbläsare, PowerShell och, för att AI ska fungera, en installerad AI-klient med stöd för lokala MCP-servrar.

## Första lanseringen

1. Installera applikationen enligt [instruktioner](INSTALL.md).
2. Skapa en projektmapp i Utforskaren, till exempel `C:\Projects\MyFirstProject`.
3. Öppna den här mappen i Utforskaren. Klicka i adressfältet, skriv `powershell` och tryck på Retur.
4. I fönstret som öppnas gör du:

```powershell
agentboard init
agentboard open
```

5. En webbläsare öppnas med adressen `http://127.0.0.1:7337`. Lämna PowerShell-fönstret öppet medan du använder whiteboardtavlan. Om du stänger fönstret stoppas den lokala servern, men uppgifterna kvarstår.

Om kommandot `agentboard` inte hittas, stäng PowerShell och öppna det igen efter installation av Scoop. När du installerar från en ZIP, använd den fullständiga sökvägen till `agentboard.exe`.

## Första uppgiften

Klicka på **Ny uppgift**, fyll i **Titel**, vid behov **Beskrivning**, sedan **Spara**. En ny uppgift i `Backlog` är tillgänglig för människor. För att AI ska börja arbeta, tilldela en arbetare och överför uppgiften till `Features`. Steg-för-steg beskrivning av fälten - [TASKS.md](TASKS.md).

## Första arbetaren

Öppna **Arbetare → Lägg till arbetare**. Välj den AI-klient du använder (till exempel Codex eller Claude Code), kontrollera namnet och kort ID `Slug`, klicka på **Spara**. Öppna den skapade arbetaren: det finns en MCP-konfiguration, en kopieringsknapp och en serverkontroll. Kopiera konfigurationen till AI-klienten enligt [WORKERS_MCP.md](WORKERS_MCP.md). När du är ansluten ber du klienten att ringa `get_my_board`.

## Vad du ska läsa härnäst

- [Uppgifter, tester, berättelser och kolumner](TASKS.md)
- [Ansluter arbetare och kontrollerar MCP](WORKERS_MCP.md)
- [Importera uppgifter från JSON](IMPORT.md)
- [Språk, tema, tidszon och uppdatering](SETTINGS.md)
- [Typiska problem](TROUBLESHOOTING.md)
