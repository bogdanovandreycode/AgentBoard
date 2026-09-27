# Agentenbord

[🌐 Languages](../LANGUAGES.md)

AgentBoard is een lokaal taakbord waar mensen en AI-werkers aan gemeenschappelijke taken werken, maar verschillende rechten hebben. De applicatie wordt gestart met één bestand `agentboard.exe`, opent de webinterface in de browser en biedt werknemers via `stdio` een aparte MCP-server. De gegevens blijven op uw computer staan.

**[Van nul beginnen](START_HERE.md) · [Werken met taken](TASKS.md) · [AI verbinden via MCP](WORKERS_MCP.md) · [JSON](IMPORT.md) importeren · [Instellingen](SETTINGS.md) · [Problemen oplossen](TROUBLESHOOTING.md)**

## Over vijf minuten

1. Download het `agentboard-VERSION-windows-amd64-setup.exe`-installatieprogramma van [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Er wordt een map voorgesteld (standaard `C:\AI\AgentBoard`) en deze wordt toegevoegd aan `PATH`. Scoop en ZIP zijn ook beschikbaar.
2. Open PowerShell in uw projectmap, bijvoorbeeld `C:\Projects\MyApp`.
3. Voer `agentboard init` uit (voor ZIP: volledig pad naar `agentboard.exe` en `init`).
4. Voer `agentboard open` uit. `http://127.0.0.1:7337` wordt geopend.
5. Voeg een taak toe via de knop **Nieuwe taak**. Voor een AI-werknemer opent u **Werknemers → Werknemer toevoegen**, selecteert u het klantprofiel en kopieert u de MCP-configuratie.

Als u nog geen projectmap heeft, maakt u er een in Windows Verkenner. Een project kan elke map zijn, zelfs zonder Git en code.

## Installatie via Scoop

In PowerShell met [Scoop](https://scoop.sh/)] al geïnstalleerd na de release:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Voor ontwikkelaars is er [build from source ](INSTALL.md). De releaseworkflow maakt een installatieprogramma, ZIP en Scoop-manifest met SHA-256 op basis van hetzelfde artefact. Updaten van de versie geïnstalleerd via Scoop: `scoop update agentboard` na toevoegen van manifest aan de bucket; details - [release voorbereiding](SCOOP_RELEASE.md).

## Hoe het bestuur is gestructureerd

`Backlog → Features → In progress → Testing → Verification → Complete`

Een persoon kan taken over het bord verplaatsen. AI kan alleen `Features → In progress → Testing → Verification` verplaatsen; de definitieve acceptatie in `Complete` wordt uitgevoerd door een mens. Gebruikerskolommen zijn bedoeld voor mensen: de taak daarin blijft in de status `Backlog` voor MCP. Testmodi: AI, menselijk en hybride. Geschiedenis, tests, artefacten en AI-kosten zijn aan de taak verbonden.

## Teams

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registreert de map en schrijft daar alleen `.agentboard/project.json`. De SQLite-productiegegevens bevinden zich in de Windows-gebruikersconfiguratiemap (`%AppData%\AgentBoard\agentboard.db`), buiten het project en buiten de Scoop-installatie. Als u het pakket verwijdert of bijwerkt, worden deze gegevens niet verwijderd. Voordat u de database overzet naar een andere computer, maakt u een kopie van de database terwijl AgentBoard is gestopt.

Werknemers - logische accounts; AgentBoard zelf voert geen Codex, Claude of enige andere AI-client uit. De klant start een lokaal MCP-proces voor een specifieke medewerker. MCP is de grens van de toepassingsrechten en gebruik voor bestandsisolatie de AI-clientsandbox.

## Voor ontwikkelaars

Stack: Go, SQLite, officiële MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Bouw eerst de frontend en ga dan naar: `./scripts/build.ps1`. Webbestanden worden via `go:embed` in het binaire bestand opgenomen. De architectuur en API worden beschreven in [doc/ARCHITECTURE.md](ARCHITECTURE.md).
