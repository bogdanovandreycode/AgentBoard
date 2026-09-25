# Архитектура и API

AgentBoard — локальный процесс Go с SQLite. HTTP API для человека, MCP-адаптер для AI и веб-интерфейс используют общую сервисную логику. SQLite — источник состояния; frontend не обходит сервер. AI-инструменты имеют отдельную поверхность прав, а состояние `Backlog`/`Complete` в MCP недоступно. Запись истории от `System` создаёт само приложение.

## Компоненты

- `cmd/agentboard` — CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` — модели и ошибки домена.
- `internal/service` — переходы задач, разрешения, импорт и настройки.
- `internal/persistence` — SQLite и миграции.
- `internal/httpapi` — Human HTTP API.
- `internal/mcpserver` — MCP-инструменты отдельного воркера.
- `web` — React/TypeScript UI; сборка попадает в `internal/webui/dist` и включается в EXE.

## Основные HTTP-маршруты

| Метод и путь | Назначение |
| --- | --- |
| `GET /api/health` | Проверка сервера. |
| `GET /api/projects` | Зарегистрированные проекты. |
| `GET /api/projects/{id}/board` | Доска проекта. |
| `GET/PUT /api/projects/{id}/settings` | Настройки, порядок и названия колонок. |
| `POST /api/projects/{id}/tasks` | Создать задачу. |
| `POST /api/projects/{id}/tasks/import` | Атомарно импортировать JSON версии 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Карточка, изменение, удаление. |
| `POST /api/tasks/{id}/move` | Перемещение человеком. |
| `GET/POST /api/projects/{id}/workers` | Список и создание воркера. |
| `GET /api/workers/{id}/mcp/check` | Внутренняя проверка MCP handshake и инструментов. |
| `GET /api/workers/{id}/sessions` | Диагностика сессий. |
| `GET/POST /api/projects/{id}/properties` | Пользовательские свойства. |

HTTP предназначен для локального доверенного пользователя. Не публикуйте веб-порт в Интернет без собственной аутентификации, сетевых ограничений и HTTPS. MCP-сервер запускается по `stdio` для конкретного проекта и воркера; начинайте с `get_my_board`. Его разрешения ограничивают действия внутри AgentBoard, но не заменяют песочницу AI-клиента для файловой системы.

Пользовательские колонки хранятся отдельно от `tasks.state`: задачу в такой колонке core держит в `backlog`, а `tasks.board_column` определяет место на Human-доске. Это сохраняет прежнюю модель AI-переходов. Удаление колонки очищает `board_column` задач.
