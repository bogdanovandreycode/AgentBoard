# 解决问题

|症状|检查什么 |
| ---| ---|
| `agentboard` 未找到 | Scoop 后重新启动 PowerShell。从 ZIP 安装时，请使用 `agentboard.exe` 的完整路径。 |
|构建后显示旧界面的网页 |停止正在运行的服务器 Ctrl+C。执行`./scripts/build.ps1`，启动一个新的二进制文件。 Vite 必须在 Go 之前构建，因为接口内置于 EXE 中。刷新页面 Ctrl+F5。 |
| 7337端口繁忙| AgentBoard 可能已经在运行。打开`http://127.0.0.1:7337`或结束旧进程。对于其他端口，请使用 `--addr`。 |
|未找到项目 |在所需的文件夹中，执行 `agentboard init`。然后是 `agentboard open` 或 `agentboard open C:\путь\к\проекту`。 |
|工作人员看不到任务 |该任务必须分配给该特定工作人员并位于 `Features`、`In progress`、`Testing` 或 `Verification` 中。 AI 看不到 `Backlog`、`Complete` 和自定义列。 |
| `NO_PROJECTS` | MCP 数据库中没有项目。在项目文件夹中执行`agentboard init`，创建worker；座席可以重复呼叫而无需重新启动。检查 MCP 和 Web 界面是否使用相同的 `--db`。 |
| `PROJECT_REQUIRED` / `PROJECT_AMBIGUOUS` |如果没有 `--project`，客户端的工作文件夹和 MCP 启动文件夹都不会位于完全相同的注册项目中。在项目文件夹中运行代理，通过`agentboard init`注册它或添加`--project`。无需启动 Web 界面。 |
| `PROJECT_NOT_FOUND` / `WORKER_NOT_FOUND` / `WORKER_REQUIRED` / `NO_WORKERS` |查看 `--project`，该项目中的工人 slug 和 `--db`。该消息列出了已注册的项目或已启用的工作人员。没有 `--project` 的自定义 MCP 配置需要每个项目中具有相同 slug 的工作人员。 |
| `PROJECT_CONTEXT_MISMATCH` | MCP 启动文件夹或客户端工作文件夹与指定的 `--project` 不匹配。请确认该问题是否是故意的，否则请删除 `--project`（自动检测）或指出正确的问题。 |
| MCP服务器检查存在，但客户端未连接 |服务器检查不检查外部客户端设置。重启客户端，检查其配置文件，以及`agentboard.exe`、`--project`、`--worker`和通用`--db`的路径。要求致电`get_my_board`。 |
|工人离线 |客户端可能已终止或尚未启动 MCP。 90 秒无心跳后，会话将被视为已断开连接。 |
| JSON 文件未导入 |检查`version: 1`、必需的`title`、现有的`Slug`工人和财产名称。 JSON 不允许注释或尾随逗号。 |
|无法重建`agentboard.exe` | Windows 无法替换正在运行的 EXE。停止服务器 Ctrl+C 并再次尝试构建。 |
|更新后任务消失了 |检查 `--db` 是否未指向另一个文件，并且您是否以同一 Windows 用户身份登录。默认底座位于 `%AppData%\AgentBoard`。 |

如果未描述错误，请收集消息的确切文本、`agentboard version` 和 Windows 版本以及重试步骤。不要在未解决的问题中发布私人项目数据或数据库内容。
