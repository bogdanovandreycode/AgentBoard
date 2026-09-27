# Workers and MCP connection

## Step 1. Create a worker

In AgentBoard, open **Workers → Add worker**. Select a client profile: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot or another MCP client. The profile fills in typical features (code, tests, Git); you can change them. The name is visible on the board, and `Slug` is a short identifier without spaces for the MCP command. Click **Save**.

One worker corresponds to one AI personality. Create different workers for different clients or teams. The capability profile describes the specialization, but does not extend the AI's rights to task stages.

## Step 2. Copy the configuration

Open the created worker. The **MCP diagnostics** block shows the configuration fragment and the file where to add it. Click **Copy MCP config**. If the file already exists, add the suggested server to the existing `mcpServers`/`servers`/`mcp` object without erasing other servers.

The main command looks like this:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` should point to the folder you registered via `agentboard init`. `--worker` — `Slug` of the created worker. The MCP client runs this command itself when it needs tools. In the AgentBoard browser, the web server can run separately.

After checking, the diagnostic block contains the absolute path to the running `agentboard.exe`. It is especially useful when installing from a ZIP. When installing via Scoop, you can use the `agentboard` command if the client sees the same `PATH`.

## Step 3: Add a server to your client

The worker interface already has a ready-made fragment. Below is an explanation of where it is used:

| Client | Where to insert | How to check on the client side |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, section `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` in project folder | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, object `mcpServers` | `/mcp list` in Gemini CLI |
| Cursor | `.cursor\mcp.json` project | list of MCP servers in Cursor settings |
| OpenCode | `opencode.json` project, object `mcp` | list of MCP tools in OpenCode |
| VS Code Copilot | `.vscode\mcp.json` project, object `servers` | command **MCP: List Servers** |

For Codex and Claude Code, an example with the project `C:\Projects\MyFirstProject` and the worker `codex`:

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

In JSON, the Windows backslash is doubled; a ready-made fragment from the interface does this automatically. If you are using `--db` with a non-standard database, add it to the `args` of the MCP configuration and specify the same path as when running `open`.

**Ollama** provides a local model, but does not replace the MCP client. The **Ollama via OpenCode** profile generates an MCP configuration for OpenCode; separately configure OpenCode on the Ollama model. Another MCP-compatible client with Ollama is also suitable.

## Step 4: Check your connection

1. Open the worker card and click **Check again**. **Server check** should show the number of MCP tools. This is an internal protocol check and detection of server tools.
2. Start or restart the AI ​​client after adding the configuration file. Have it call `get_my_board`.
3. **Client connected** will appear in AgentBoard and a new session will appear in the list. Only this confirms the connection of your client. If there are tools, but the client is not connected, check the path to the program, the name of the configuration file and its JSON/TOML syntax.

Start your work session with `get_my_board`. AI does not receive tasks from `Backlog` and `Complete`, even if it knows their ID. AI cannot pretend to be human and does not have a general “move anywhere” command. Working with the file system outside of AgentBoard depends on the capabilities and sandbox of the selected client.

Official customer instructions: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
