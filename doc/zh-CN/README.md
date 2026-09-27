# AgentBoard

[🌐 Languages](../LANGUAGES.md)

AgentBoard 是一个本地任务板，人类和人工智能工作人员可以在其中执行常见任务，但拥有不同的权限。该应用程序通过一个文件 `agentboard.exe` 启动，在浏览器中打开 Web 界面，并通过 `stdio` 为工作人员提供单独的 MCP 服务器。数据保留在您的计算机上。

**[从头开始](START_HERE.md)·[使用任务](TASKS.md)·[通过MCP](WORKERS_MCP.md)连接AI]·[导入JSON](IMPORT.md)·[设置](SETTINGS.md)·[解决问题](TROUBLESHOOTING.md)]]

## 五分钟后

1. 从 [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases).zip] 下载 `agentboard-VERSION-windows-amd64-setup.exe` 安装程序。它将建议一个文件夹（默认为 `C:\AI\AgentBoard`）并将其添加到 `PATH`。还提供 Scoop 和 ZIP。
2. 在项目文件夹中打开 PowerShell，例如 `C:\Projects\MyApp`。
3. 执行 `agentboard init`（对于 ZIP：`agentboard.exe` 和 `init` 的完整路径）。
4. 执行`agentboard open`。 `http://127.0.0.1:7337` 将打开。
5. 使用 **新任务** 按钮添加任务。对于 AI 工作人员，打开 **工作人员 → 添加工作人员**，选择客户端配置文件并复制 MCP 配置。

如果您还没有项目文件夹，请在 Windows 资源管理器中创建一个。项目可以是任何文件夹，即使没有 Git 和代码。

## 通过 Scoop 安装

在发布后已安装 [Scoop](https://scoop.sh/)] 的 PowerShell 中：

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

对于开发人员来说，有[从源代码 ](INSTALL.md) 构建。发布工作流程使用同一工件中的 SHA-256 创建安装程序、ZIP 和 Scoop 清单。将清单添加到存储桶后，更新通过 Scoop 安装的版本：`scoop update agentboard`；详细信息 - [发布准备](SCOOP_RELEASE.md)。

## 董事会的结构如何

`Backlog → Features → In progress → Testing → Verification → Complete`

一个人可以在板上移动任务。 AI只能移动`Features → In progress → Testing → Verification`； `Complete` 的最终接受是由人执行的。用户列适用于人类：其中的任务对于 MCP 仍处于 `Backlog` 状态。测试模式：人工智能、人类和混合。历史记录、测试、工件和人工智能成本都附加在任务中。

## 团队

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` 注册文件夹并仅在其中写入 `.agentboard/project.json`。 SQLite 生产数据位于 Windows 用户配置目录 (`%AppData%\AgentBoard\agentboard.db`) 中，位于项目外部和 Scoop 安装外部。删除或更新包不应删除此数据。在传输到另一台计算机之前，请在 AgentBoard 停止时复制数据库。

工人——逻辑账户； AgentBoard 本身不运行 Codex、Claude 或任何其他 AI 客户端。客户端为特定工作人员启动本地 MCP 进程。 MCP是应用程序权限边界，并且为了文件隔离，使用AI客户端沙箱。

## 对于开发者

Stack：Go、SQLite、官方 MCP Go SDK、React、TypeScript、Vite、PrimeReact、TanStack Query、dnd-kit。首先构建前端，然后转到：`./scripts/build.ps1`。 Web 文件通过 `go:embed` 包含在二进制文件中。架构和 API 在 [doc/ARCHITECTURE.md](ARCHITECTURE.md) 中进行了描述。
