# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Náhled AgentBoard](../../assets/social-preview.png)

**[Stáhnout pro Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Web s dokumentací](https://bogdanovandreycode.github.io/AgentBoard/) · [Licence MIT](../../LICENSE)**

AgentBoard je místní panel úkolů, kde lidé a AI pracovníci pracují na společných úkolech, ale mají různá práva. Aplikace se spouští s jedním souborem `agentboard.exe`, otevře webové rozhraní v prohlížeči a poskytuje pracovníkům samostatný MCP server přes `stdio`. Data zůstávají ve vašem počítači.

**[Začít od nuly](START_HERE.md)· [Práce s úkoly](TASKS.md)· [Připojení AI přes MCP](WORKERS_MCP.md)· [Import JSON](IMPORT.md)· [Nastavení](SETTINGS.md)· [Řešení problémů](TROUBLESHOOTING.md)**

## Funkce MVP

- Místní projektová deska s vyhledáváním úkolů, importem JSON, vlastními sloupci a vlastnostmi.
- Přístup MCP pro konkrétního pracovníka s omezenými přechody AI a konečným lidským přijetím úkolů.
- Obecná historie úkolů, kontrolní pokyny, artefakty, náklady na AI a diagnostika připojení pracovníků.
- Instalační program systému Windows, manifest ZIP a Scoop; webové rozhraní je zabudováno do spustitelného souboru.

AgentBoard je navržen pro důvěryhodného místního uživatele. Aplikace nehostí projekty v cloudu a sama nespouští AI klienty; v případě potřeby připojte k pracovníkovi klienta kompatibilního s MCP.

## Za pět minut

1. Stáhněte si instalační program `agentboard-VERSION-windows-amd64-setup.exe` z [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Navrhne složku (ve výchozím nastavení `C:\AI\AgentBoard`) a přidá ji do `PATH`. Scoop a ZIP jsou také k dispozici.
2. Otevřete PowerShell ve složce projektu, například `C:\Projects\MyApp`.
3. Spusťte `agentboard init` (pro ZIP: úplná cesta k `agentboard.exe` a `init`).
4. Spusťte `agentboard open`. Otevře se `http://127.0.0.1:7337`.
5. Přidejte úkol pomocí tlačítka **Nový úkol**. Pro AI pracovníka otevřete **Workers → Add worker**, vyberte profil klienta a zkopírujte konfiguraci MCP.

Pokud ještě nemáte složku projektu, vytvořte si ji v Průzkumníkovi Windows. Projekt může být jakákoli složka, dokonce i bez Gitu a kódu.

## Instalace přes Scoop

V PowerShellu s již nainstalovaným [Scoop](https://scoop.sh/)] po vydání:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Pro vývojáře existuje [sestavení ze zdroje ](INSTALL.md). Pracovní postup vydání vytvoří instalační program, ZIP a Scoop manifest s SHA-256 ze stejného artefaktu. Aktualizace verze nainstalované přes Scoop: `scoop update agentboard` po přidání manifestu do bucketu; podrobnosti - [příprava k vydání](SCOOP_RELEASE.md).

## Jak je strukturována deska

`Backlog → Features → In progress → Testing → Verification → Complete`

Osoba může přesouvat úkoly po tabuli. AI může pohybovat pouze `Features → In progress → Testing → Verification`; konečné přijetí do `Complete` provádí člověk. Uživatelské sloupce jsou určeny pro lidi: úloha v nich zůstává pro MCP ve stavu `Backlog`. Testovací režimy: AI, Human a Hybrid. Historie, testy, artefakty a náklady na AI jsou připojeny k úkolu.

## Týmy

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` zaregistruje složku a zapíše tam pouze `.agentboard/project.json`. Produkční data SQLite jsou umístěna v adresáři uživatelské konfigurace Windows (`%AppData%\AgentBoard\agentboard.db`), mimo projekt a mimo instalaci Scoop. Odstraněním nebo aktualizací balíčku by tato data neměla být odstraněna. Před přenosem do jiného počítače vytvořte kopii databáze, zatímco je AgentBoard zastaven.

Pracovníci - logické účty; Samotný AgentBoard nespouští Codex, Claude ani žádný jiný AI klient. Klient spustí místní proces MCP pro konkrétního pracovníka. MCP je hranice oprávnění aplikace a pro izolaci souborů použijte karanténu klienta AI.

## Pro vývojáře

Stack: Go, SQLite, oficiální MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Nejprve sestavte frontend, pak jděte: `./scripts/build.ps1`. Webové soubory jsou zahrnuty v binárním formátu přes `go:embed`. Architektura a API jsou popsány v [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Licence

AgentBoard je bezplatný a open source projekt pod [licence MIT](../../LICENSE). Komerční použití, úpravy, větvení a redistribuce jsou povoleny za předpokladu, že jsou zachovány autorská práva a licence.
