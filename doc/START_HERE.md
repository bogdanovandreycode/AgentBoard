# Start here

## What is AgentBoard

Imagine a regular board with task cards. You create tasks, assign someone responsible and oversee the work. AI workers receive only their assigned tasks via MCP and report progress. You decide when the task is finally ready.

**Project** - a folder on the computer and a separate board. **Task** - a card with a description, responsible person and stage. **Worker** is the logical name of the AI ​​client. **MCP** is the way the AI ​​client connects to the AgentBoard. **Scoop** is a software installation manager for Windows.

No code required. You need Windows, a browser, PowerShell and, for AI to work, an installed AI client with support for local MCP servers.

## First launch

1. Install the application according to [instructions](INSTALL.md).
2. Create a project folder in Explorer, for example `C:\Projects\MyFirstProject`.
3. Open this folder in Explorer. Click in the address bar, type `powershell` and press Enter.
4. In the window that opens, do:

```powershell
agentboard init
agentboard open
```

5. A browser will open at `http://127.0.0.1:7337`. Leave the PowerShell window open while you use the board. Closing the window will stop the local server, but the tasks will remain.

If the `agentboard` command is not found, close PowerShell and reopen after installing Scoop. When installing from a ZIP, use the full path to `agentboard.exe`.

## First task

Click **New task**, fill in **Title**, if necessary **Description**, then **Save**. A new task in `Backlog` is available to a person. For the AI ​​to start working, assign a worker and move the task to `Features`. Step-by-step description of the fields - [TASKS.md](TASKS.md).

## First worker

Open **Workers → Add worker**. Select the AI ​​client you are using (for example Codex or Claude Code), check the name and short identifier `Slug`, click **Save**. Open the created worker: there is an MCP configuration, a copy button and a server check. Copy the configuration to the AI ​​client according to [WORKERS_MCP.md](WORKERS_MCP.md). Once connected, have the client call `get_my_board`.

## What to read next

- [Tasks, tests, history and columns](TASKS.md)
- [Connecting workers and checking MCP](WORKERS_MCP.md)
- [Import tasks from JSON](IMPORT.md)
- [Language, theme, time zone and update](SETTINGS.md)
- [Typical problems](TROUBLESHOOTING.md)
