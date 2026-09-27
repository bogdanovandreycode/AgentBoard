# Архітектура та API

AgentBoard - локальний процес Go з SQLite. HTTP API для людини, MCP-адаптер для AI та веб-інтерфейс використовують загальну сервісну логіку. SQLite - джерело стану; frontend не оминає сервер. AI-інструменти мають окрему поверхню прав, а стан `Backlog`/`Complete` у MCP недоступний. Запис історії від `System` створює сам додаток.

## Компоненти

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - моделі та помилки домену.
- `internal/service` - переходи завдань, дозволу, імпорт та налаштування.
- `internal/persistence` - SQLite та міграції.
- `internal/httpapi` - Human HTTP API.
- `internal/mcpserver` - MCP-інструменти окремого воркера.
- `web` - React/TypeScript UI; збірка потрапляє до `internal/webui/dist` та включається до EXE.

## Основні HTTP-маршрути

| Метод та шлях | Призначення |
| --- | --- |
| `GET /api/health` | Перевірка сервера. |
| `GET /api/projects` | Зареєстровані проекти. |
| `GET /api/projects/{id}/board` | Проект проекту. |
| `GET/PUT /api/projects/{id}/settings` | Налаштування, порядок та назви колонок. |
| `POST /api/projects/{id}/tasks` | Створити завдання. |
| `POST /api/projects/{id}/tasks/import` | Атомарно імпортувати JSON версії 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Картка, зміна, видалення. |
| `POST /api/tasks/{id}/move` | Переміщення людиною. |
| `GET/POST /api/projects/{id}/workers` | Список та створення воркера. |
| `GET /api/workers/{id}/mcp/check` | Внутрішня перевірка MCP handshake та інструментів. |
| `GET /api/workers/{id}/sessions` | Діагностика сесій. |
| `GET/POST /api/projects/{id}/properties` | Користувацькі властивості. |

HTTP призначено для локального довіреного користувача. Не публікуйте веб-порт в Інтернеті без власної автентифікації, мережевих обмежень та HTTPS. MCP-сервер запускається по `stdio` для конкретного проекту та воркера; починайте з `get_my_board`. Його дозволи обмежують дії всередині AgentBoard, але не замінюють пісочницю AI-клієнта для файлової системи.

Користувальницькі колонки зберігаються окремо від `tasks.state`: завдання в такій колонці core тримає `backlog`, а `tasks.board_column` визначає місце на Human-дошці. Це зберігає попередню модель AI-переходів. Видалення колонки очищує задачі `board_column`.
