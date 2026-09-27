# Архитектура и API

AgentBoard е локален Go процес, изпълняващ SQLite. HTTP API за хора, MCP адаптерът за AI и уеб интерфейсът споделят обща логика на услугата. SQLite - източник на състояние; интерфейсът не заобикаля сървъра. AI инструментите имат отделна повърхност за права и статусът `Backlog`/`Complete` в MCP не е наличен. Записът на историята от `System` се създава от самото приложение.

## Компоненти

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - модели на домейни и грешки.
- `internal/service` - преходи на задачи, разрешения, импортиране и настройки.
- `internal/persistence` - SQLite и миграции.
- `internal/httpapi` - Човешки HTTP API.
- `internal/mcpserver` - MCP инструменти за отделен работник.
- `web` - React/TypeScript UI; сборката завършва в `internal/webui/dist` и е включена в EXE.

## Основни HTTP маршрути

| Метод и път | Дестинация |
| --- | --- |
| `GET /api/health` | Проверка на сървъра. |
| `GET /api/projects` | Регистрирани проекти. |
| `GET /api/projects/{id}/board` | Проектна дъска. |
| `GET/PUT /api/projects/{id}/settings` | Настройки, ред и имена на колони. |
| `POST /api/projects/{id}/tasks` | Създайте задача. |
| `POST /api/projects/{id}/tasks/import` | Импортирайте атомарно JSON версия 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Карта, промяна, изтриване. |
| `POST /api/tasks/{id}/move` | Придвижване от човек. |
| `GET/POST /api/projects/{id}/workers` | Списък и създаване на работник. |
| `GET /api/workers/{id}/mcp/check` | Вътрешно тестване на MCP ръкостискане и инструменти. |
| `GET /api/workers/{id}/sessions` | Диагностика на сесията. |
| `GET/POST /api/projects/{id}/properties` | Персонализирани свойства. |

HTTP е за локалния доверен потребител. Не публикувайте уеб порт в Интернет без собствено удостоверяване, мрежови ограничения и HTTPS. MCP сървърът се стартира с помощта на `stdio` за конкретен проект и работник; започнете с `get_my_board`. Неговите разрешения ограничават действията в рамките на AgentBoard, но не заместват пясъчника на файловата система на AI клиента.

Потребителските колони се съхраняват отделно от `tasks.state`: ядрото запазва задачата в такава колона в `backlog`, а `tasks.board_column` определя мястото на дъската Human. Това поддържа същия модел на преход на AI. Премахването на колона изчиства задачите `board_column`.
