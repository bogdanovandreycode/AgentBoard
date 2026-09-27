# AgentBoard

[🌐 Languages](../LANGUAGES.md)

AgentBoard е локален съвет за задачи, където хората и работниците с AI работят по общи задачи, но имат различни права. Приложението се стартира с един файл `agentboard.exe`, отваря уеб интерфейса в браузъра и предоставя на работниците отделен MCP сървър чрез `stdio`. Данните остават на вашия компютър.

**[Започнете от нулата](START_HERE.md) · [Работа със задачи](TASKS.md) · [Свързване на AI чрез MCP](WORKERS_MCP.md) · [Импортиране на JSON](IMPORT.md) · [Настройки](SETTINGS.md) · [Разрешаване на проблеми](TROUBLESHOOTING.md)**

## След пет минути

1. Изтеглете инсталатора на `agentboard-VERSION-windows-amd64-setup.exe` от [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Той ще предложи папка (по подразбиране `C:\AI\AgentBoard`) и ще я добави към `PATH`. Предлагат се също Scoop и ZIP.
2. Отворете PowerShell в папката на вашия проект, например `C:\Projects\MyApp`.
3. Изпълнете `agentboard init` (за ZIP: пълен път до `agentboard.exe` и `init`).
4. Изпълнете `agentboard open`. `http://127.0.0.1:7337` ще се отвори.
5. Добавете задача чрез бутона **Нова задача**. За AI работник отворете **Работници → Добавяне на работник**, изберете клиентския профил и копирайте MCP конфигурацията.

Ако все още нямате папка на проекта, създайте такава в Windows Explorer. Проектът може да бъде всяка папка, дори без Git и код.

## Инсталиране чрез Scoop

В PowerShell с [Scoop](https://scoop.sh/)] вече инсталиран след изданието:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

За разработчиците има [изграждане от източник ](INSTALL.md). Работният процес на освобождаване създава инсталатор, ZIP и Scoop манифест с SHA-256 от един и същи артефакт. Актуализиране на версията, инсталирана чрез Scoop: `scoop update agentboard` след добавяне на манифест към кофата; подробности - [подготовка за освобождаване](SCOOP_RELEASE.md).

## Как е структурирана дъската

`Backlog → Features → In progress → Testing → Verification → Complete`

Човек може да мести задачи по дъската. AI може да движи само `Features → In progress → Testing → Verification`; окончателното приемане в `Complete` се извършва от човек. Потребителските колони са предназначени за хора: задачата в тях остава в състояние `Backlog` за MCP. Тестови режими: AI, Human и Hybrid. Към задачата са приложени история, тестове, артефакти и разходи за AI.

## Отбори

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` регистрира папката и записва там само `.agentboard/project.json`. Данните за производство на SQLite се намират в директорията за потребителска конфигурация на Windows (`%AppData%\AgentBoard\agentboard.db`), извън проекта и извън инсталацията на Scoop. Премахването или актуализирането на пакета не трябва да премахва тези данни. Преди да прехвърлите на друг компютър, направете копие на базата данни, докато AgentBoard е спрян.

Работници - логически сметки; Самият AgentBoard не изпълнява Codex, Claude или който и да е друг AI клиент. Клиентът стартира локален MCP процес за конкретен работник. MCP е границата на разрешението за приложение и за изолиране на файлове използвайте пясъчната среда на клиента на AI.

## За разработчици

Стек: Go, SQLite, официален MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Първо изградете интерфейс, след това Go: `./scripts/build.ps1`. Уеб файловете са включени в двоичния файл чрез `go:embed`. Архитектурата и API са описани в [doc/ARCHITECTURE.md](ARCHITECTURE.md).
