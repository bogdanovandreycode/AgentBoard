# Problem solving

| Symptom | What to check |
| --- | --- |
| `agentboard` not found | Restart PowerShell after Scoop. When installing from a ZIP, use the full path to `agentboard.exe`. |
| Web page showing old interface after build | Stop the running server Ctrl+C. Execute `./scripts/build.ps1`, launch the new binary. Vite must build before Go because the interface is built into the EXE. Refresh the page Ctrl+F5. |
| Port 7337 busy | AgentBoard may already be running. Open `http://127.0.0.1:7337` or terminate the old process. For a different port, use `--addr`. |
| Project not found | In the desired folder, run `agentboard init`. Then run `agentboard open` there or `agentboard open C:\Projects\MyProject`. |
| `PROJECT_LOOKUP_FAILED` / `WORKER_LOOKUP_FAILED` | Check the project path, this project's worker slug, and `--db`. A global MCP configuration pinned to another project will not select the current project. Copy the configuration from that project's worker card. |
| `PROJECT_CONTEXT_MISMATCH` | MCP launch directory or client workspace conflicts with the configured project. Confirm the question if intentional, otherwise reconnect with the correct `--project` and `--worker`. Clients without form elicitation must correct their configuration or explicitly use `--allow-project-mismatch`. |
| `PROJECT_REQUIRED` | MCP was launched outside registered projects without `--project`. Supply the registered project path explicitly. |
| Missing `APPDATA` in an MCP host | On Windows, AgentBoard resolves `%USERPROFILE%\AppData\Roaming` instead of opening a temporary database. For redirected profiles or custom storage, use the exact `--db` from the worker configuration. If no profile can be resolved, startup returns an actionable error. |
| The worker does not see the task | The task must be assigned to this particular worker and be in `Features`, `In progress`, `Testing` or `Verification`. AI does not see `Backlog`, `Complete` and custom columns. |
| MCP server check exists, but the client is not connected | The server check does not check the external client settings. Restart the client, check its configuration file, path to `agentboard.exe`, `--project`, `--worker` and general `--db`. Ask to call `get_my_board`. |
| Worker Offline | The client may have terminated or may not have started MCP yet. After 90 seconds without a heartbeat, the session is considered disconnected. |
| JSON file not importing | Check `version: 1`, required `title`, existing `Slug` workers and property names. JSON does not allow comments or trailing commas. |
| Cannot rebuild `agentboard.exe` | Windows cannot replace a running EXE. Stop the server Ctrl+C and try the build again. |
| Tasks disappeared after update | Check that `--db` does not point to a different file and that you are logged in as the same Windows user. The default database is located in `%AppData%\AgentBoard`. |

If the error is not described, collect the exact text of the message, the `agentboard version` and Windows versions, and steps to retry. Do not publish private project data or database contents in an open issue.
