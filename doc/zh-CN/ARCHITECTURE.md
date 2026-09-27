# 架构和 API

AgentBoard 是一个运行 SQLite 的本地 Go 进程。用于人类的 HTTP API、用于 AI 的 MCP 适配器以及 Web 界面共享通用的服务逻辑。 SQLite - 状态源；前端不会绕过服务器。 AI工具具有单独的权限表面，并且MCP中的`Backlog`/`Complete`状态不可用。 `System` 的历史记录是由应用程序本身创建的。

## 组件

- `cmd/agentboard` - CLI `init`、`open`、`serve`、`mcp`、`version`。
- `internal/core` - 域模型和错误。
- `internal/service` - 任务转换、权限、导入和设置。
- `internal/persistence` - SQLite 和迁移。
- `internal/httpapi` - 人类 HTTP API。
- `internal/mcpserver` - 用于单独工作人员的 MCP 工具。
- `web` - React/TypeScript UI；该程序集最终为 `internal/webui/dist` 并包含在 EXE 中。

## 基本 HTTP 路由

|方法与路径|目的地 |
| --- | --- |
|`GET /api/health`|服务器检查。 |
|`GET /api/projects`|已注册项目。 |
|`GET /api/projects/{id}/board`|项目板。 |
|`GET/PUT /api/projects/{id}/settings`|列的设置、顺序和名称。 |
|`POST /api/projects/{id}/tasks`|创建任务。 |
|`POST /api/projects/{id}/tasks/import`|原子导入 JSON 版本 1。
|`GET/PATCH/DELETE /api/tasks/{id}`|卡、更改、删​​除。 |
|`POST /api/tasks/{id}/move`|一个人移动。 |
|`GET/POST /api/projects/{id}/workers`|工人的列表和创建。 |
|`GET /api/workers/{id}/mcp/check`| MCP 握手和工具的内部测试。 |
|`GET /api/workers/{id}/sessions`|会话诊断。 |
|`GET/POST /api/projects/{id}/properties`|自定义属性。 |

HTTP 适用于本地可信用户。不要将没有其自身身份验证、网络限制和 HTTPS 的 Web 端口发布到 Internet。 MCP 服务器由`stdio`对于特定项目和工人；开始于`get_my_board`。其权限限制AgentBoard内的操作，但不会取代AI客户端的文件系统沙箱。

用户列分别存储`tasks.state`：这样一栏的任务是由核心保存的`backlog`, 一个`tasks.board_column`决定在人类板上的位置。这保持了相同的人工智能转换模型。删除列会清除`board_column`任务。
