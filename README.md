# AgentBoard

**Documentation:** [English](README.md) · [Русский](doc/ru/README.md) · [all 24 languages](doc/LANGUAGES.md)

AgentBoard is a local task board where human and AI workers work on common tasks, but have different rights. The application is launched with a single file `agentboard.exe`, opens the web interface in the browser and provides workers with a separate MCP server via `stdio`. The data remains on your computer.

**[Start from scratch](doc/START_HERE.md) · [Working with tasks](doc/TASKS.md) · [Connecting AI via MCP](doc/WORKERS_MCP.md) · [Import JSON](doc/IMPORT.md) · [Settings](doc/SETTINGS.md) · [Solving problems](doc/TROUBLESHOOTING.md)**

## In five minutes

1. Download the `agentboard-VERSION-windows-amd64-setup.exe` installer from [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). It will suggest a folder (by default `C:\AI\AgentBoard`) and add it to `PATH`. Scoop and ZIP are also available.
2. Open PowerShell in your project folder, for example `C:\Projects\MyApp`.
3. Run `agentboard init` (for ZIP: full path to `agentboard.exe` and `init`).
4. Execute `agentboard open`. This will open `http://127.0.0.1:7337`.
5. Add a task using the **New task** button. For an AI worker, open **Workers → Add worker**, select the client profile and copy the MCP configuration.

If you don't already have a project folder, create one in Windows Explorer. A project can be any folder, even without Git and code.

## Installation via Scoop

In PowerShell with [Scoop](https://scoop.sh/) already installed after the release:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

For developers there is [build from source](doc/INSTALL.md). The release workflow creates an installer, ZIP and Scoop manifest with SHA-256 from the same artifact. Updating the version installed via Scoop: `scoop update agentboard` after adding manifest to the bucket; details - [release preparation](doc/SCOOP_RELEASE.md).

## How the board is structured

`Backlog → Features → In progress → Testing → Verification → Complete`

A person can move tasks around the board. AI can only move `Features → In progress → Testing → Verification`; the final acceptance in `Complete` is performed by a human. User columns are intended for humans: the task in them remains in the `Backlog` state for MCP. Test modes: AI, Human and Hybrid. History, tests, artifacts and AI costs are attached to the task.

## Commands

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registers the folder and writes only `.agentboard/project.json` to it. Application data is stored in the Windows user configuration directory (`%AppData%\AgentBoard\agentboard.db`), outside the project and Scoop installation. Removing or updating the package does not remove this data. Before transferring it to another computer, stop AgentBoard and copy the database.

Workers - logical accounts; AgentBoard itself does not run Codex, Claude or any other AI client. The client starts a local MCP process for a specific worker. MCP is the application permission boundary, and for file isolation, use the AI ​​client sandbox.

## For developers

Stack: Go, SQLite, official MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Build frontend first, then Go: `./scripts/build.ps1`. Web files are included in the binary via `go:embed`. The architecture and API are described in [doc/ARCHITECTURE.md](doc/ARCHITECTURE.md).
