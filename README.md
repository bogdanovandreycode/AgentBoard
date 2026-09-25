# AgentBoard

AgentBoard — локальная доска задач, на которой человек и AI-воркеры работают с общими задачами, но имеют разные права. Приложение запускается одним файлом `agentboard.exe`, открывает веб-интерфейс в браузере и предоставляет воркерам отдельный MCP-сервер по `stdio`. Данные остаются на вашем компьютере.

**[Начать с нуля](doc/START_HERE.md) · [Работа с задачами](doc/TASKS.md) · [Подключение AI через MCP](doc/WORKERS_MCP.md) · [Импорт JSON](doc/IMPORT.md) · [Настройки](doc/SETTINGS.md) · [Решение проблем](doc/TROUBLESHOOTING.md)**

## За пять минут

1. Установите AgentBoard через Scoop после публикации первого релиза (инструкция ниже) или скачайте ZIP из [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases) и распакуйте его.
2. Откройте PowerShell в папке вашего проекта, например `C:\Projects\MyApp`.
3. Выполните `agentboard init` (для ZIP: полный путь к `agentboard.exe` и `init`).
4. Выполните `agentboard open`. Откроется `http://127.0.0.1:7337`.
5. Добавьте задачу кнопкой **New task**. Для AI-сотрудника откройте **Workers → Add worker**, выберите профиль клиента и скопируйте MCP-конфигурацию.

Если у вас ещё нет папки проекта, создайте её в Проводнике Windows. Проектом может быть любая папка, даже без Git и кода.

## Установка через Scoop

В PowerShell с уже установленным [Scoop](https://scoop.sh/) после выхода релиза:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Первого GitHub Release пока нет: до его публикации эта команда не сработает. Для разработчиков есть [сборка из исходников](doc/INSTALL.md). Release workflow создаёт ZIP и Scoop manifest с SHA-256 из одного и того же артефакта. Обновление установленной версии: `scoop update agentboard` после добавления manifest в bucket; подробности — [подготовка релиза](doc/SCOOP_RELEASE.md).

## Как устроена доска

`Backlog → Features → In progress → Testing → Verification → Complete`

Человек может перемещать задачи по доске. AI может двигаться только `Features → In progress → Testing → Verification`; финальное принятие в `Complete` выполняет человек. Пользовательские колонки предназначены для человека: задача в них остаётся в состоянии `Backlog` для MCP. Режимы тестирования: AI, Human и Hybrid. История, тесты, артефакты и затраты AI прикреплены к задаче.

## Команды

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` регистрирует папку и записывает туда только `.agentboard/project.json`. Рабочие данные SQLite находятся в каталоге конфигурации пользователя Windows (`%AppData%\AgentBoard\agentboard.db`), вне проекта и вне установки Scoop. Удаление или обновление пакета не должно удалять эти данные. Перед переносом на другой компьютер сделайте копию базы при остановленном AgentBoard.

Воркеры — логические учётные записи; AgentBoard сам не запускает Codex, Claude или другой AI-клиент. Клиент запускает локальный MCP-процесс для конкретного воркера. MCP является границей прав приложения, а для изоляции файлов используйте песочницу AI-клиента.

## Для разработчиков

Стек: Go, SQLite, официальный MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Сначала собирайте frontend, затем Go: `./scripts/build.ps1`. Веб-файлы включаются в бинарник через `go:embed`. Архитектура и API описаны в [doc/ARCHITECTURE.md](doc/ARCHITECTURE.md).

Documentation in Russian is under [`doc/`](doc/START_HERE.md). The app is local-first, ships as one Windows executable, and exposes a worker-scoped stdio MCP server. English UI is available in Settings.
