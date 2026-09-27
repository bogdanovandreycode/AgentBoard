# Архітэктура і API

AgentBoard - лакальны працэс Go з SQLite. HTTP API для чалавека, MCP-адаптар для AI і вэб-інтэрфейс выкарыстоўваюць агульную сэрвісную логіку. SQLite - крыніца стану; frontend не абыходзіць сервер. AI-інструменты маюць асобную паверхню правоў, а стан `Backlog`/`Complete` у MCP недаступны. Запіс гісторыі ад `System` стварае само дадатак.

## Кампаненты

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - мадэлі і памылкі дамена.
- `internal/service` - пераходы задач, дазволу, імпарт і налады.
- `internal/persistence` - SQLite і міграцыі.
- `internal/httpapi` - Human HTTP API.
- `internal/mcpserver` - MCP-інструменты асобнага воркера.
- `web` - React / TypeScript UI; зборка пападае ў `internal/webui/dist` і ўключаецца ў EXE.

## Асноўныя HTTP-маршруты

| Метад і шлях | Прызначэнне |
| --- | --- |
| `GET /api/health` | Праверка сервера. |
| `GET /api/projects` | Зарэгістраваныя праекты. |
| `GET /api/projects/{id}/board` | Дошка праекту. |
| `GET/PUT /api/projects/{id}/settings` | Настройкі, парадак і назвы калонак. |
| `POST /api/projects/{id}/tasks` | Стварыць задачу. |
| `POST /api/projects/{id}/tasks/import` | Атамарна імпартаваць JSON версіі 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Картка, змяненне, выдаленне. |
| `POST /api/tasks/{id}/move` | Перамяшчэнне чалавекам. |
| `GET/POST /api/projects/{id}/workers` | Спіс і стварэнне воркера. |
| `GET /api/workers/{id}/mcp/check` | Унутраная праверка MCP handshake і інструментаў. |
| `GET /api/workers/{id}/sessions` | Дыягностыка сесій. |
| `GET/POST /api/projects/{id}/properties` | Карыстальніцкія ўласцівасці. |

HTTP прызначаны для лакальнага даверанага карыстальніка. Не публікуйце вэб-порт у Інтэрнэт без уласнай аўтэнтыфікацыі, сеткавых абмежаванняў і HTTPS. MCP-сервер запускаецца па `stdio` для канкрэтнага праекта і воркера; пачынайце з `get_my_board`. Яго дазволы абмяжоўваюць дзеянні ўсярэдзіне AgentBoard, але не замяняюць пясочніцу AI-кліента для файлавай сістэмы.

Карыстальніцкія калонкі захоўваюцца асобна ад `tasks.state`: задачу ў такой калонцы core трымае ў `backlog`, а `tasks.board_column` вызначае месца на Human-дошцы. Гэта захоўвае ранейшую мадэль AI-пераходаў. Выдаленне калонкі чысціць `board_column` задач.
