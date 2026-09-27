# Instalace a spuštění v systému Windows

## Instalační program (doporučeno)

Stáhněte si `agentboard-VERSION-windows-amd64-setup.exe` ze stránky [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Průvodce instalací navrhne složku; výchozí je `C:\AI\AgentBoard`. Zkopíruje `agentboard.exe` a dokumentaci, vytvoří zástupce a přidá vybranou složku do systému `PATH`. Po instalaci otevřete nový terminál, aby byl k dispozici příkaz `agentboard`.

Ve složce projektu spusťte:

```powershell
agentboard init
agentboard open
```

Rozhraní je zabudováno do `agentboard.exe`; Není potřeba žádná samostatná instalace Go nebo Node.js. Data jsou uložena v `%AppData%\AgentBoard` a jsou zachována, když je program aktualizován nebo odinstalován. Odinstalace pomocí „Installed Applications“ odstraní zástupce a položku ze `PATH`.

## Instalace Scoop

Pokud Scoop ještě není nainstalován, otevřete PowerShell jako váš běžný uživatel a postupujte podle [oficiálních pokynů Scoop](https://scoop.sh/). Pokud jsou na vašem podnikovém počítači nějaká omezení, obraťte se na správce; AgentBoard lze také spustit ze ZIP bez Scoop.

Alternativně nainstalujte AgentBoard přes Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Dostupnost manifestu ve verzi GitHub lze zkontrolovat na [stránce Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Po přidání manifestu do Scoop lze bucket nainstalovat podle názvu bucketu a aktualizovat pomocí příkazu `scoop update agentboard`.

## ZIP bez lopatky

Stáhněte si `agentboard-VERSION-windows-amd64.zip` z Releases, rozbalte například do `C:\Tools\AgentBoard`. V PowerShellu ve složce projektu:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

U klienta AI zadejte úplnou cestu k `agentboard.exe` v jeho konfiguraci MCP, pokud se program nenachází v `PATH`.

## Sestavte ze zdroje

Nainstalujte verzi Go ze `go.mod` a Node.js 22 nebo novější. V PowerShellu v kořenovém adresáři úložiště:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` sestaví frontend do `internal/webui/dist`, spustí testy Go a sestaví jeden `agentboard.exe`. Pořadí je důležité: rozhraní je zabudováno do binárního kódu při sestavení Go. Pokud je spuštěn starý `agentboard.exe open`, zastavte jej před přestavbou (Ctrl+C), jinak vám systém Windows nedovolí soubor nahradit.

## Kde jsou data

- `%AppData%\AgentBoard\agentboard.db` - úkoly, projekty, pracovníci a nastavení. Pomocí příznaku `--db` můžete zadat jiný soubor, ale `init`, `open`/`serve` a `mcp` musí mít **stejnou cestu**.
- `<ваш проект>\.agentboard\project.json` — identifikátor projektu. Tento soubor neobsahuje úkoly.
- Server ve výchozím nastavení naslouchá pouze `127.0.0.1:7337`. Před cestou projektu zadejte jinou adresu `--addr`: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` otevře stránku projektu a znovu použije již běžící server na dané adrese. Pokud je v jednom prohlížeči otevřeno několik projektů, vyberte je ze seznamu vlevo.
