# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Перетворити AgentBoard](../../assets/social-preview.png)

**[Скачати для Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Сайт документації](https://bogdanovandreycode.github.io/AgentBoard/)] [Ліцензія MIT](../../LICENSE)**

AgentBoard — локальна дошка завдань, де людина і AI-воркери працюють із спільними завданнями, але мають різні права. Програма запускається одним файлом `agentboard.exe`, відкриває веб-інтерфейс у браузері та надає воркерам окремий MCP-сервер з `stdio`. Дані залишаються на вашому комп'ютері.

**[Почати з нуля](START_HERE.md)· [Робота із завданнями](TASKS.md)· [Підключення AI через MCP](WORKERS_MCP.md)· [Імпорт JSON](IMPORT.md)· [Налаштування](SETTINGS.md)· [Вирішення проблем](TROUBLESHOOTING.md)**

## Можливості MVP

- Локальна дошка проектів з пошуком завдань, імпортом JSON, колонками користувача та властивостями.
- MCP-доступ для конкретного воркера з обмеженими переходами AI та фінальним прийняттям завдань людиною.
- Загальна історія завдань, інструкції перевірки, артефакти, витрати AI та діагностика підключення воркерів.
- інсталятор Windows, ZIP і Scoop manifest; веб-інтерфейс вбудований у виконуваний файл.

AgentBoard розрахований на довіреного локального користувача. Додаток не розміщує проекти у хмарі та не запускає AI-клієнти саме; у разі потреби підключіть сумісний з MCP клієнт до воркера.

## За п'ять хвилин

1. Завантажте інсталятор `agentboard-VERSION-windows-amd64-setup.exe` з [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Він запропонує папку (за замовчуванням `C:\AI\AgentBoard`) та додасть її до `PATH`. Доступні також Scoop та ZIP.
2. Відкрийте PowerShell у папці проекту, наприклад `C:\Projects\MyApp`.
3. Виконайте `agentboard init` (для ZIP: повний шлях до `agentboard.exe` та `init`).
4. Виконайте `agentboard open`. Відкриється `http://127.0.0.1:7337`.
5. Додайте завдання кнопкою **New task**. Для AI-співробітника відкрийте **Workers → Add worker**, виберіть профіль клієнта та скопіюйте MCP-конфігурацію.

Якщо у вас ще немає папки проекту, створіть її у Провіднику Windows. Проектом може бути будь-яка папка, навіть без Git та коду.

## Установка через Scoop

PowerShell з уже встановленим [Scoop](https://scoop.sh/) після виходу релізу:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Для розробників є [складання з вихідних кодів ](INSTALL.md). Release workflow створює установник, ZIP та Scoop manifest із SHA-256 з одного і того ж артефакту. Оновлення встановленої через Scoop версії: `scoop update agentboard` після додавання manifest до bucket; подробиці - [підготовка релізу ](SCOOP_RELEASE.md).

## Як влаштована дошка

`Backlog → Features → In progress → Testing → Verification → Complete`

Людина може переміщати завдання на дошці. AI може рухатися лише `Features → In progress → Testing → Verification`; фінальне прийняття `Complete` виконує людина. Користувальницькі колонки призначені для людини: завдання залишається в стані `Backlog` для MCP. Режими тестування: AI, Human та Hybrid. Історія, тести, артефакти та витрати AI прикріплені до завдання.

## Команди

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` реєструє папку та записує туди лише `.agentboard/project.json`. Робочі дані SQLite знаходяться в каталозі конфігурації користувача Windows (`%AppData%\AgentBoard\agentboard.db`), поза проектом та поза встановлення Scoop. Видалення або оновлення пакета не повинно видаляти ці дані. Перед перенесенням на інший комп'ютер зробіть копію бази під час зупиненого AgentBoard.

Воркери - логічні облікові записи; AgentBoard сам не запускає Codex, Claude чи інший AI-клієнт. Клієнт запускає локальний MCP-процес для конкретного воркера. MCP є межею прав програми, а для ізоляції файлів використовуйте пісочницю AI-клієнта.

## Для розробників

Стек: Go, SQLite, офіційний MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Спочатку збирайте frontend, а потім Go: `./scripts/build.ps1`. Веб-файли включаються до бінарника через `go:embed`. Архітектура та API описані в [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Ліцензія

AgentBoard - безкоштовний проект з відкритим вихідним кодом під [ліцензією MIT](../../LICENSE). Дозволено комерційне використання, зміну, створення форків та розповсюдження при збереженні повідомлення про авторські права та ліцензії.
