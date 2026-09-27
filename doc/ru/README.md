# AgentBoard

![Превью AgentBoard](../../assets/social-preview.png)

[🌐 Languages](../LANGUAGES.md)

**[Скачать для Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Сайт документации](https://bogdanovandreycode.github.io/AgentBoard/) · [Лицензия MIT](../../LICENSE)**

AgentBoard — локальная доска задач, на которой человек и AI-воркеры работают с общими задачами, но имеют разные права. Приложение запускается одним файлом `agentboard.exe`, открывает веб-интерфейс в браузере и предоставляет воркерам отдельный MCP-сервер по `stdio`. Данные остаются на вашем компьютере.

**[Начать с нуля](START_HERE.md) · [Работа с задачами](TASKS.md) · [Подключение AI через MCP](WORKERS_MCP.md) · [Импорт JSON](IMPORT.md) · [Настройки](SETTINGS.md) · [Решение проблем](TROUBLESHOOTING.md)**

## Возможности MVP

- Локальная доска проектов с поиском задач, импортом JSON, пользовательскими колонками и свойствами.
- MCP-доступ для конкретного воркера с ограниченными переходами AI и финальным принятием задач человеком.
- Общая история задач, инструкции проверки, артефакты, затраты AI и диагностика подключения воркеров.
- Установщик Windows, ZIP и Scoop manifest; веб-интерфейс встроен в исполняемый файл.

AgentBoard рассчитан на доверенного локального пользователя. Приложение не размещает проекты в облаке и не запускает AI-клиенты само; при необходимости подключите совместимый с MCP клиент к воркеру.

## За пять минут

1. Скачайте установщик `agentboard-VERSION-windows-amd64-setup.exe` из [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Он предложит папку (по умолчанию `C:\AI\AgentBoard`) и добавит её в `PATH`. Доступны также Scoop и ZIP.
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

Для разработчиков есть [сборка из исходников](INSTALL.md). Release workflow создаёт установщик, ZIP и Scoop manifest с SHA-256 из одного и того же артефакта. Обновление установленной через Scoop версии: `scoop update agentboard` после добавления manifest в bucket; подробности — [подготовка релиза](SCOOP_RELEASE.md).

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

Стек: Go, SQLite, официальный MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Сначала собирайте frontend, затем Go: `./scripts/build.ps1`. Веб-файлы включаются в бинарник через `go:embed`. Архитектура и API описаны в [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Лицензия

AgentBoard — бесплатный проект с открытым исходным кодом под [лицензией MIT](../../LICENSE). Разрешено коммерческое использование, изменение, создание форков и распространение при сохранении уведомления об авторских правах и лицензии.
