# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![ папярэдні прагляд AgentBoard](../../assets/social-preview.png)

**[Сьцягнуць для Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Сайт дакументацыі](https://bogdanovandreycode.github.io/AgentBoard/) · [Ліцэнзія MIT](../../LICENSE)**]

AgentBoard – лакальная дошка задач, на якой чалавек і AI-воркеры працуюць з агульнымі задачамі, але маюць розныя правы. Прыкладанне запускаецца адным файлам `agentboard.exe`, адчыняе вэб-інтэрфейс у браўзэры і падае воркерам асобны MCP-сервер па `stdio`. Дадзеныя застаюцца на вашым кампутары.

**[Пачаць з нуля](START_HERE.md)· [Праца з задачамі](TASKS.md)· [Падлучэнне AI праз MCP](WORKERS_MCP.md)· [Імпарт JSON](IMPORT.md)· [Налады](SETTINGS.md)· [Рашэнне праблем](TROUBLESHOOTING.md)**

## Магчымасці MVP

- Лакальная дошка праектаў з пошукам задач, імпартам JSON, карыстацкімі калонкамі і ўласцівасцямі.
- MCP-доступ для канкрэтнага воркера з абмежаванымі пераходамі AI і фінальным прыняццем задач чалавекам.
- Агульная гісторыя задач, інструкцыі праверкі, артэфакты, выдаткі AI і дыягностыка падключэння воркераў.
- Усталёўшчык Windows, ZIP і Scoop manifest; вэб-інтэрфейс убудаваны ў выкананы файл.

AgentBoard разлічаны на даверанага лакальнага карыстальніка. Прыкладанне не размяшчае праекты ў воблаку і не запускае AI-кліенты само; пры неабходнасці падключыце сумяшчальны з MCP кліент да воркера.

## За пяць хвілін

1. Запампуйце ўсталёўшчык `agentboard-VERSION-windows-amd64-setup.exe` з [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Ён прапануе тэчку (па змаўчанні `C:\AI\AgentBoard`) і дадасць яе ў `PATH`. Даступныя таксама Scoop і ZIP.
2. Адкрыйце PowerShell у тэчцы вашага праекта, напрыклад `C:\Projects\MyApp`.
3. Выканайце `agentboard init` (для ZIP: поўны шлях да `agentboard.exe` і `init`).
4. Выканайце `agentboard open`. Адкрыецца `http://127.0.0.1:7337`.
5. Дадайце задачу кнопкай **New task**. Для AI-супрацоўніка адкрыйце **Workers → Add worker**, абярыце профіль кліента і скапіруйце MCP-канфігурацыю.

Калі ў вас яшчэ няма тэчкі праекту, стварыце яе ў Правадыру Windows. Праектам можа быць любая тэчка, нават без Git і кода.

## Устаноўка праз Scoop

У PowerShell з ужо ўсталяваным [Scoop](https://scoop.sh/) пасля выхаду рэлізу:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Для распрацоўшчыкаў ёсць [зборка з зыходнікаў](INSTALL.md). Release workflow стварае ўсталёўшчык, ZIP і Scoop manifest з SHA-256 з аднаго і таго ж артэфакта. Абнаўленне ўсталяванай праз Scoop версіі: `scoop update agentboard` пасля дадання manifest у bucket; падрабязнасці - [падрыхтоўка рэлізу ](SCOOP_RELEASE.md).

## Як уладкована дошка

`Backlog → Features → In progress → Testing → Verification → Complete`

Чалавек можа перамяшчаць задачы па дошцы. AI можа рухацца толькі `Features → In progress → Testing → Verification`; фінальнае прыняцце ў `Complete` выконвае чалавек. Карыстальніцкія калонкі прызначаны для чалавека: задача ў іх застаецца ў стане `Backlog` для MCP. Рэжымы тэсціравання: AI, Human і Hybrid. Гісторыя, тэсты, артэфакты і выдаткі AI прымацаваныя да задачы.

## Каманды

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` рэгіструе тэчку і запісвае туды толькі `.agentboard/project.json`. Працоўныя дадзеныя SQLite знаходзяцца ў каталогу канфігурацыі карыстача Windows (`%AppData%\AgentBoard\agentboard.db`), па-за праектам і па-за ўсталёўкай Scoop. Выдаленне ці абнаўленне пакета не павінна выдаляць гэтыя дадзеныя. Перад пераносам на іншы кампутар зрабіце копію базы пры спыненым AgentBoard.

Воркеры - лагічныя ўліковыя запісы; AgentBoard сам не запускае Codex, Claude ці іншы AI-кліент. Кліент запускае лакальны MCP-працэс для канкрэтнага воркера. MCP з'яўляецца мяжой правоў прыкладання, а для ізаляцыі файлаў выкарыстоўвайце пясочніцу AI-кліента.

## Для распрацоўшчыкаў

Стэк: Go, SQLite, афіцыйны MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Спачатку зьбірайце frontend, затым Go: `./scripts/build.ps1`. Вэб-файлы ўключаюцца ў бінарнік праз `go:embed`. Архітэктура і API апісаны ў [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Ліцэнзія

AgentBoard - бясплатны праект з адкрытым зыходным кодам пад [ліцэнзіяй MIT](../../LICENSE). Дазволена камерцыйнае выкарыстанне, змена, стварэнне форкаў і распаўсюджванне пры захаванні апавяшчэння аб аўтарскіх правах і ліцэнзіі.
