# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Förhandsgranska AgentBoard](../../assets/social-preview.png)

**[Ladda ner för Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Dokumentationswebbplats](https://bogdanovandreycode.github.io/AgentBoard/) · [Licens MIT](../../LICENSE)**

AgentBoard är en lokal arbetsgrupp där mänskliga och AI-arbetare arbetar med gemensamma uppgifter, men har olika rättigheter. Applikationen startas med en fil `agentboard.exe`, öppnar webbgränssnittet i webbläsaren och förser arbetare med en separat MCP-server via `stdio`. Datan finns kvar på din dator.

**[Börja från början](START_HERE.md)· [Arbeta med uppgifter](TASKS.md)· [AI-anslutning via MCP](WORKERS_MCP.md)· [Importera JSON](IMPORT.md)· [Inställningar](SETTINGS.md)· [Problemlösning](TROUBLESHOOTING.md)**

## MVP-funktioner

- Lokalt projektkort med uppgiftssökning, JSON-import, anpassade kolumner och egenskaper.
- MCP-åtkomst för en specifik arbetare med begränsade AI-övergångar och slutlig mänsklig acceptans av uppgifter.
- Allmän uppgiftshistorik, kontrollinstruktioner, artefakter, AI-kostnader och diagnostik för arbetaranslutningar.
- Windows installationsprogram, ZIP och Scoop-manifest; webbgränssnittet är inbyggt i den körbara filen.

AgentBoard är designat för en betrodd lokal användare. Applikationen är värd för projekt i molnet och startar inte AI-klienter själv; anslut om nödvändigt en MCP-kompatibel klient till arbetaren.

## Om fem minuter

1. Ladda ner `agentboard-VERSION-windows-amd64-setup.exe` installationsprogrammet från [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Det kommer att föreslå en mapp (som standard `C:\AI\AgentBoard`) och lägga till den i `PATH`. Scoop och ZIP finns också.
2. Öppna PowerShell i din projektmapp, till exempel `C:\Projects\MyApp`.
3. Kör `agentboard init` (för ZIP: fullständig sökväg till `agentboard.exe` och `init`).
4. Kör `agentboard open`. `http://127.0.0.1:7337` öppnas.
5. Lägg till en uppgift med knappen **Ny uppgift**. För en AI-arbetare, öppna **Arbetare → Lägg till arbetare**, välj klientprofilen och kopiera MCP-konfigurationen.

Om du inte redan har en projektmapp, skapa en i Utforskaren i Windows. Ett projekt kan vara vilken mapp som helst, även utan Git och kod.

## Installation via Scoop

I PowerShell med [Scoop](https://scoop.sh/)] redan installerat efter releasen:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

För utvecklare finns [bygga från källan ](INSTALL.md). Utgivningsarbetsflödet skapar ett installationsprogram, ZIP och Scoop-manifest med SHA-256 från samma artefakt. Uppdatering av versionen installerad via Scoop: `scoop update agentboard` efter att ha lagt till manifest i hinken; detaljer - [release förberedelse](SCOOP_RELEASE.md).

## Hur styrelsen är uppbyggd

`Backlog → Features → In progress → Testing → Verification → Complete`

En person kan flytta uppgifter runt styrelsen. AI kan bara flytta `Features → In progress → Testing → Verification`; slutlig acceptans till `Complete` utförs av en människa. Användarkolumner är avsedda för människor: uppgiften i dem förblir i tillståndet `Backlog` för MCP. Testlägen: AI, Human och Hybrid. Historik, tester, artefakter och AI-kostnader är kopplade till uppgiften.

## Lag

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registrerar mappen och skriver endast `.agentboard/project.json` där. SQLite-produktionsdata finns i Windows användarkonfigurationskatalog (`%AppData%\AgentBoard\agentboard.db`), utanför projektet och utanför Scoop-installationen. Att ta bort eller uppdatera paketet bör inte ta bort dessa data. Innan du överför till en annan dator, gör en kopia av databasen medan AgentBoard är stoppad.

Arbetare - logiska konton; AgentBoard själv kör inte Codex, Claude eller någon annan AI-klient. Klienten startar en lokal MCP-process för en specifik arbetare. MCP är applikationens behörighetsgräns, och för filisolering, använd AI-klientsandlådan.

## För utvecklare

Stack: Go, SQLite, officiell MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Bygg frontend först, sedan Gå: `./scripts/build.ps1`. Webbfiler ingår i binären via `go:embed`. Arkitekturen och API beskrivs i [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Licens

AgentBoard är ett gratis och öppen källkodsprojekt under [licens MIT](../../LICENSE). Kommersiell användning, modifiering, splittring och omdistribution är tillåtet förutsatt att upphovsrättsmeddelandet och licensen bibehålls.
