# Workers and MCP connection

## Step 1. Create a worker

In AgentBoard, open **Workers → Add worker**. Select a client profile: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot or another MCP client. The profile fills in typical features (code, tests, Git); you can change them. The name is visible on the board, and `Slug` is a short identifier without spaces for the MCP command. Click **Save**.

One worker corresponds to one AI personality. Create different workers for different clients or teams. The capability profile describes the specialization, but does not extend the AI's rights to task stages.

## Step 2. Copy the configuration

Open the created worker. The **MCP diagnostics** block shows the configuration fragment and the file where to add it. Click **Copy MCP config**. If the file already exists, add the suggested server to the existing `mcpServers`/`servers`/`mcp` object without erasing other servers.

The main command looks like this:

```text
agentboard mcp --worker codex
```

`--worker` is the `Slug` of the created worker. The MCP client runs this command
itself whenever it needs tools.

### MCP works without the web UI

`agentboard mcp` is a separate stdio process that opens the AgentBoard database
directly. It does not need `agentboard open`/`serve`, does not depend on the folder
where the web UI was started, and is not affected by how many web UI instances are
running. Several agents can work at the same time, each in its own project: every
client starts its own MCP process, and all of them safely share one SQLite database.

### How the project is selected

The MCP server starts in any folder. On the first tool call it binds to a project:

1. `--project PATH`, when given (any folder inside a registered project works);
2. otherwise the client's workspace roots, when the client reports them;
3. otherwise the folder where the client launched the MCP server.

`get_my_board` returns the bound project and worker in `session`, together with
`binding_source`. If the workspace roots change, the next tool call rebinds.
Without `--worker`, the project's only enabled worker is used.

If the server cannot bind, every tool returns an error the agent must relay to the
user, and nothing is accessed:

| Code | Meaning |
| --- | --- |
| `NO_PROJECTS` | The database has no projects. Run `agentboard init` in the project folder and create a worker. No restart is needed. |
| `PROJECT_REQUIRED` | The workspace is outside registered projects. The message lists them. |
| `PROJECT_NOT_FOUND` | `--project` is not inside a registered project. |
| `PROJECT_AMBIGUOUS` | Workspace roots belong to several projects. Add `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Create, enable or choose the worker with `--worker`. |

The configurations generated for **Codex** and **Gemini CLI** are user-level and
shared by every folder, so they omit `--project`: create a worker with the same slug
in each project where that client works. Project-level files (Claude Code, Cursor,
OpenCode, VS Code) include `--project` for their own project.

### Explicit project and context checks

When `--project` is set and the launch folder or the client workspace belongs to a
different project, tools ask for confirmation through MCP elicitation. Clients
without form elicitation receive `PROJECT_CONTEXT_MISMATCH` before any task access.
Correct the configuration, remove `--project` so the project is detected
automatically, or explicitly confirm intentional cross-project work with
`--allow-project-mismatch`. Workspace changes invalidate previous confirmations.
Automatically detected projects need no confirmation.

**Global settings → MCP project context** controls these checks for every project
and worker using the same database. It is enabled by default. When disabled, MCP
uses its configured project without confirmation or mismatch errors. Existing
sessions apply the setting on their next tool call; task visibility, workflow
transitions and human-testing rules remain enforced.

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

Examples for Codex (user-level, project detected automatically) and Claude Code (project file) with the project `C:\Projects\MyFirstProject` and the worker `codex`:

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

In JSON, the Windows backslash is doubled; a ready-made fragment from the interface does this automatically. Generated configurations include the actual `--db` used by the web server. This keeps HTTP and MCP on the same database even when the client has a different environment.

**Ollama** provides a local model, but does not replace the MCP client. The **Ollama via OpenCode** profile generates an MCP configuration for OpenCode; separately configure OpenCode on the Ollama model. Another MCP-compatible client with Ollama is also suitable.

## Step 4: Check your connection

1. Open the worker card and click **Check again**. **Server check** launches the actual stdio executable with this project, worker and database, verifies the handshake and tool discovery, and checks that the project and worker resolve. It does not verify the external AI client's configuration. **Enabled** is permission to connect, not evidence of a successful connection.
2. Start or restart the AI ​​client after adding the configuration file. Have it call `get_my_board`.
3. **Client connected** will appear in AgentBoard and a new session will appear in the list. Only this confirms the connection of your client. If there are tools, but the client is not connected, check the path to the program, the name of the configuration file and its JSON/TOML syntax.

Start your work session with `get_my_board`. AI does not receive tasks from `Backlog` and `Complete`, even if it knows their ID. AI cannot pretend to be human and does not have a general “move anywhere” command. Working with the file system outside of AgentBoard depends on the capabilities and sandbox of the selected client.

`task_id` accepts full UUIDs or unique case-insensitive prefixes of at least eight
characters, scoped to the worker's visible assigned tasks. Ambiguous prefixes return
`AMBIGUOUS_TASK_ID`; hidden tasks do not participate in matching. Source and dependency
IDs in `create_feature` require full UUIDs. Boards and transitions return compact task
summaries; `get_task` returns the full context and native screenshot blocks. JSON is
sent once as text, without a duplicate `structuredContent` payload.

Every tool declares `title` and all four MCP behaviour hints (`readOnlyHint`,
`destructiveHint`, `idempotentHint`, `openWorldHint`) as explicit booleans. Read
tools are read-only and idempotent; write tools only add records or move tasks
forward, so none is destructive, and none reaches outside the local AgentBoard
database.

For human or hybrid testing, finish implementation, move the card to `testing`, add
handoff instructions, and end the session or work on another task. Never wait or poll
for the human. The verification gate remains enforced in the Go service.

`record_test_run.status` accepts `passed`, `failed`, `skipped`, or `manual_required`.
`create_feature.testing_mode` accepts `ai`, `human`, or `hybrid`; priority accepts
`critical`, `high`, `medium`, or `low`. Artifact `kind`, test `type`, and `runner` are
free-form labels. `report_usage.MCPCalls` is a cumulative worker-session snapshot
including that call, not a task total; do not sum snapshots. `WorkerSessionID` and
`mcp_calls_scope` identify its scope.

Official customer instructions: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
