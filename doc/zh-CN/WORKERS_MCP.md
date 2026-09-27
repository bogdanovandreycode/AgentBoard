# Workers 和 MCP 连接

## 步骤1.创建一个worker

在 AgentBoard 中，打开 **工作人员 → 添加工作人员**。选择客户端配置文件：Codex、Claude Code、Gemini CLI、Cursor、OpenCode、Ollama via OpenCode、VS Code Copilot 或其他 MCP 客户端。配置文件填写典型功能（代码、测试、Git）；你可以改变它们。该名称在板上可见，`Slug` 是 MCP 命令的短标识符，不带空格。单击**保存**。

一名工人对应一名人工智能人格。为不同的客户或团队创建不同的工作人员。能力概况描述了专业化，但没有将人工智能的权利扩展到任务阶段。

## 步骤2.复制配置

打开创建的worker。 **MCP 诊断** 块显示配置片段以及添加它的文件。单击 **复制 MCP 配置**。如果该文件已存在，请将建议的服务器添加到现有的 `mcpServers`/`servers`/`mcp` 对象，而不删除其他服务器。

主要命令如下所示：

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` 应指向您通过 `agentboard init` 注册的文件夹。 `--worker` — 创建的工作线程的 `Slug`。 MCP 客户端在需要工具时自行运行此命令。在AgentBoard浏览器中，Web服务器可以单独运行。

检查后，诊断块显示正在运行的`agentboard.exe`的绝对路径。从 ZIP 安装时它特别有用。通过 Scoop 安装时，如果客户端看到相同的 `agentboard`，则可以使用 `PATH` 命令。

## 步骤 3：向您的客户端添加服务器

工作界面已经有一个现成的片段。下面是它的使用位置的解释：

|客户|在哪里插入 |如何在客户端查看 |
| --- | --- | --- |
|法典| `%USERPROFILE%\.codex\config.toml`，`[mcp_servers.agentboard_<slug>]` 节 | `codex mcp list` |
|克劳德·代码 |项目文件夹中的`.mcp.json` | `claude mcp list` |
|双子座 CLI | `%USERPROFILE%\.gemini\settings.json`，对象 `mcpServers` | Gemini CLI 中的 `/mcp list` |
|光标| `.cursor\mcp.json`项目|光标设置中的 MCP 服务器列表 |
|开放代码 | `opencode.json` 项目，对象 `mcp` | OpenCode 中的 MCP 工具列表 |
| VS Code 副驾驶 | `.vscode\mcp.json` 项目，对象 `servers` |命令 **MCP：列出服务器** |

对于 Codex 和 Claude Code，以项目 `C:\Projects\MyFirstProject` 和工作人员 `codex` 为例：

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

在 JSON 中，Windows 反斜杠加倍；界面中的现成片段会自动执行此操作。如果你使用`--db`使用自定义基础，将其添加到`args`MCP配置并指定与启动时相同的路径`open`。

**Ollama** 提供本地模型，但不会取代 MCP 客户端。 **Ollama via OpenCode** 配置文件生成 OpenCode 的 MCP 配置；在 Ollama 模型上单独配置 OpenCode。另一个与 Ollama 兼容的 MCP 客户端也适用。

## 步骤 4：检查您的连接

1、打开工人卡，点击**再次检查**。 **服务器检查**应显示 MCP 工具的数量。这是一个内部协议检查和检测服务器的工具。
2、添加配置文件后启动或重启AI客户端。请他打电话`get_my_board`。
3. **客户端已连接**将出现在 AgentBoard 中，并且新会话将出现在列表中。只有这样才能确认您的客户端的连接。如果有工具，但客户端未连接，请检查程序路径、配置文件名称及其JSON/TOML语法。

开始你的工作会议`get_my_board`。 AI 不接收任务`Backlog`和`Complete`，即使他知道他们的 ID。人工智能无法假装人类，也没有通用的“移动到任何地方”的命令。在 AgentBoard 外部使用文件系统取决于所选客户端的功能和沙箱。

官方客户说明：[Codex](https://developers.openai.com/learn/docs-mcp)、[Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp)、[Gemini CLI](https://geminicli.com/docs/tools/mcp-server/)]、[Cursor](https://prod.cursor.com/docs/cli/mcp)、[OpenCode](https://opencode.ai/docs/mcp-servers/)、[VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers)]。
