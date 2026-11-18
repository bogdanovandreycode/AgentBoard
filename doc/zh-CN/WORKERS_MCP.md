# Workers 和 MCP 连接

## 步骤1.创建一个worker

在 AgentBoard 中，打开 **工作人员 → 添加工作人员**。选择客户端配置文件：Codex、Claude Code、Gemini CLI、Cursor、OpenCode、Ollama via OpenCode、VS Code Copilot 或其他 MCP 客户端。配置文件填写典型功能（代码、测试、Git）；你可以改变它们。该名称在板上可见，`Slug` 是 MCP 命令的短标识符，不带空格。单击**保存**。

一名工人对应一名人工智能人格。为不同的客户或团队创建不同的工作人员。能力概况描述了专业化，但没有将人工智能的权利扩展到任务阶段。

## 步骤2.复制配置

打开创建的worker。 **MCP 诊断** 块显示配置片段以及添加它的文件。单击 **复制 MCP 配置**。如果该文件已存在，请将建议的服务器添加到现有的 `mcpServers`/`servers`/`mcp` 对象，而不删除其他服务器。

主要命令如下所示：

```text
agentboard mcp --worker codex
```

`--worker` — 创建的工作线程的 `Slug`。 MCP 客户端在需要工具时自行运行此命令。

### MCP 无需 Web 界面即可工作

`agentboard mcp` 是一个单独的 stdio 进程，它本身打开 AgentBoard 数据库。它不需要 `agentboard open`/`serve`，它不依赖于启动 Web 界面的文件夹，也不依赖于正在运行的 Web 界面实例的数量。多个代理可以同时运行，每个代理都在自己的项目中：每个客户端运行自己的 MCP 进程，并且它们都安全地使用相同的 SQLite 数据库。

### 如何选择项目

MCP 服务器在任何文件夹中运行。第一次调用该工具时，它与项目绑定：

1. `--project PATH`，如果指定（注册项目内的任何文件夹都适用）；
2. 否则，如果客户报告了客户的工作文件夹（根目录）；
3. 否则客户端启动 MCP 服务器的文件夹。

`get_my_board` 将选定的项目和工作人员与 `session` 一起返回到 `binding_source`。如果客户端的工作文件夹发生更改，则下次调用该工具时将再次执行绑定。如果没有 `--worker`，则使用项目中唯一启用的工作线程。

如果无法绑定，则每个工具都会返回一个错误，代理必须将该错误传递给用户；在这种情况下不会读取任务数据：

|代码|意义|
| ---| ---|
| `NO_PROJECTS` |数据库中没有项目。在项目文件夹中执行`agentboard init`并创建worker。无需重新启动。 |
| `PROJECT_REQUIRED` |工作文件夹位于注册项目之外。该消息列出了项目。 |
| `PROJECT_NOT_FOUND` | `--project` 不在注册项目内。 |
| `PROJECT_AMBIGUOUS` |客户的工作文件夹跨越多个项目。添加`--project`。 |
| `NO_WORKERS`、`WORKER_REQUIRED`、`WORKER_NOT_FOUND`、`WORKER_DISABLED` |通过 `--worker` 创建、启用或选择工作人员。 |

**Codex** 和 **Gemini CLI** 的配置是自定义的，并在所有文件夹之间共享，因此它们没有 `--project`：在运行此客户端的每个项目中创建具有相同 slug 的工作程序。项目级文件（Claude Code、Cursor、OpenCode、VS Code）包含其项目的 `--project`。

### 显式设计和上下文检查

如果指定了 `--project` 并且客户端的启动文件夹或工作文件夹位于不同的项目中，则工具会要求通过 MCP 获取进行确认。没有表单诱导的客户端在访问任务之前会收到 `PROJECT_CONTEXT_MISMATCH`。更正配置，删除 `--project` 以进行自动检测，或明确确认您打算使用带有 `--allow-project-mismatch` 标志的另一个项目。更改工作文件夹会重置之前的确认。自动定义的项目不需要确认。

**全局设置 → MCP 项目上下文** 控制对同一数据库的所有项目和工作人员的这些检查。默认启用。如果禁用，MCP 将使用指定的项目，而不会出现确认或不匹配错误。已经打开的会话会在下次调用该工具时应用该设置；继续遵循任务可见性、转换和人工测试规则。

检查后，诊断块显示正在运行的`agentboard.exe`的绝对路径。从 ZIP 安装时它特别有用。通过 Scoop 安装时，如果客户端看到相同的 `agentboard`，则可以使用 `PATH` 命令。

## 步骤 3：向您的客户端添加服务器

工作界面已经有一个现成的片段。下面是它的使用位置的解释：

|客户|在哪里插入|如何在客户端查看 |
| ---| ---| ---|
|法典| `%USERPROFILE%\.codex\config.toml`，`[mcp_servers.agentboard_<slug>]` 节 | `codex mcp list` |
|克劳德·代码 |项目文件夹中的`.mcp.json` | `claude mcp list` |
|双子座 CLI | `%USERPROFILE%\.gemini\settings.json`，对象 `mcpServers` | Gemini CLI 中的 `/mcp list` |
|光标| `.cursor\mcp.json`项目|光标设置中的 MCP 服务器列表 |
|开放代码 | `opencode.json` 项目，对象 `mcp` | OpenCode 中的 MCP 工具列表 |
| VS Code 副驾驶 | `.vscode\mcp.json` 项目，对象 `servers` |命令 **MCP：列出服务器** |

Codex（用户配置、自动检测的项目）和 Claude Code（项目文件）以及项目 `C:\Projects\MyFirstProject` 和工作器 `codex` 的示例：

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

在 JSON 中，Windows 反斜杠加倍；界面中的现成片段会自动执行此操作。生成的配置包含实际的`--db`Web 服务器，因此 HTTP 和 MCP 使用相同的基础，即使客户端环境不同。

**Ollama** 提供本地模型，但不会取代 MCP 客户端。 **Ollama via OpenCode** 配置文件生成 OpenCode 的 MCP 配置；在 Ollama 模型上单独配置 OpenCode。另一个与 Ollama 兼容的 MCP 客户端也适用。

## 步骤 4：检查您的连接

1、打开工人卡，点击**再次检查**。 **服务器检查** 使用该项目、工作人员和基础运行真正的 stdio 可执行文件，检查握手、工具列表以及项目和工作人员的位置。
2、添加配置文件后启动或重启AI客户端。请他打电话`get_my_board`。
3. **客户端已连接**将出现在 AgentBoard 中，并且新会话将出现在列表中。只有这样才能确认您的客户端的连接。如果有工具，但客户端未连接，请检查程序路径、配置文件名称及其JSON/TOML语法。

开始你的工作会议`get_my_board`。 AI 不接收任务`Backlog`和`Complete`，即使他知道他们的 ID。人工智能无法假装人类，也没有通用的“移动到任何地方”的命令。在 AgentBoard 外部使用文件系统取决于所选客户端的功能和沙箱。

每个工具都将 `title` 和所有四个 MCP 行为提示（`readOnlyHint`、`destructiveHint`、`idempotentHint`、`openWorldHint`）声明为显式布尔值。读取工具是只读的、幂等的；记录工具仅添加记录或向前推进任务，因此不会造成破坏或超出本地 AgentBoard 基础。

官方客户说明：[Codex](https://developers.openai.com/learn/docs-mcp)、[Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp)、[Gemini CLI](https://geminicli.com/docs/tools/mcp-server/)]、[Cursor](https://prod.cursor.com/docs/cli/mcp)、[OpenCode](https://opencode.ai/docs/mcp-servers/)、[VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers)]。
