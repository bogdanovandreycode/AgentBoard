# Устаноўка і запуск на Windows

## Усталёўшчык (рэкамендуецца)

Запампуйце `agentboard-VERSION-windows-amd64-setup.exe` са старонкі [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Майстар усталёўкі прапануе тэчку; па змаўчанні гэта `C:\AI\AgentBoard`. Ён скапіюе `agentboard.exe` і дакументацыю, створыць ярлык і дадасць выбраную тэчку ў сістэмны `PATH`. Пасля ўстаноўкі адкрыйце новы тэрмінал, каб каманда `agentboard` стала даступна.

У тэчцы вашага праекта выканайце:

```powershell
agentboard init
agentboard open
```

Інтэрфейс убудаваны ў `agentboard.exe`; асобная ўстаноўка Go ці Node.js не патрэбна. Дадзеныя захоўваюцца ў `%AppData%\AgentBoard` і захоўваюцца пры абнаўленні ці выдаленні праграмы. Выдаленне праз «Усталяваныя прыкладанні» прыбірае цэтлікі і запіс з `PATH`.

## Устаноўка Scoop

Калі Scoop яшчэ не ўсталяваны, адкрыйце PowerShell ад свайго звычайнага карыстальніка і прытрымлівайцеся [афіцыйнай інструкцыі Scoop](https://scoop.sh/). Пры абмежаваннях карпаратыўнага кампутара звернецеся да адміністратара; AgentBoard таксама можна запусціць з ZIP без Scoop.

Альтэрнатыўна ўсталюеце AgentBoard праз Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Наяўнасць manifest у GitHub Release можна праверыць на [старонцы Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Пасля дадання manifest у Scoop bucket можна ўсталёўваць па імі bucket і абнаўляць камандай `scoop update agentboard`.

## ZIP без Scoop

Запампуйце `agentboard-VERSION-windows-amd64.zip` з Releases, распакуйце, напрыклад, у `C:\Tools\AgentBoard`. У PowerShell у тэчцы праекта:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Для AI-кліента пакажыце ў яго MCP-канфігурацыі поўны шлях да `agentboard.exe`, калі праграма не знаходзіцца ў `PATH`.

## Зборка з зыходнікаў

Усталюйце Go версіі з `go.mod` і Node.js 22 ці навей. У PowerShell у корані рэпазітара:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` збірае frontend у `internal/webui/dist`, запускае Go-тэсты і збірае адзін `agentboard.exe`. Парадак важны: інтэрфейс убудоўваецца ў бінарнік пры зборцы Go. Калі стары `agentboard.exe open` запушчаны, спыніце яго перад перазборкай (Ctrl+C), інакш Windows не дасць замяніць файл.

## Дзе дадзеныя

- `%AppData%\AgentBoard\agentboard.db` - задачы, праекты, воркеры і налады. Можна задаць іншы файл сцягам `--db`, але ў `init`, `open`/`serve` і `mcp` павінен быць **адзін і той жа шлях**.
- `<ваш проект>\.agentboard\project.json` - ідэнтыфікатар праекта. Гэты файл не змяшчае задач.
- Сервер слухае толькі `127.0.0.1:7337` па змаўчанні. Іншы адрас укажыце `--addr` да шляху праекта: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` адчыняе старонку праекту і паўторна выкарыстае ўжо які працуе сервер на гэтым адрасе. Калі ў адным браўзэры адчынена некалькі праектаў, выбірайце іх праз спіс злева.
