# Problem solving

| Symptom | What to check |
| --- | --- |
| `agentboard` not found | Restart PowerShell after Scoop. When installing from a ZIP, use the full path to `agentboard.exe`. |
| Web page showing old interface after build | Stop the running server Ctrl+C. Execute `./scripts/build.ps1`, launch the new binary. Vite must build before Go because the interface is built into the EXE. Refresh the page Ctrl+F5. |
| Port 7337 busy | AgentBoard may already be running. Open `http://127.0.0.1:7337` or terminate the old process. For a different port, use `--addr`. |
| Project not found | In the desired folder, run `agentboard init`. Then run `agentboard open` there or `agentboard open C:\Projects\MyProject`. |
| The worker does not see the task | The task must be assigned to this particular worker and be in `Features`, `In progress`, `Testing` or `Verification`. AI does not see `Backlog`, `Complete` and custom columns. |
| MCP server check exists, but the client is not connected | The server check does not check the external client settings. Restart the client, check its configuration file, path to `agentboard.exe`, `--project`, `--worker` and general `--db`. Ask to call `get_my_board`. |
| Worker Offline | The client may have terminated or may not have started MCP yet. After 90 seconds without a heartbeat, the session is considered disconnected. |
| JSON file not importing | Check `version: 1`, required `title`, existing `Slug` workers and property names. JSON does not allow comments or trailing commas. |
| Cannot rebuild `agentboard.exe` | Windows cannot replace a running EXE. Stop the server Ctrl+C and try the build again. |
| Tasks disappeared after update | Check that `--db` does not point to a different file and that you are logged in as the same Windows user. The default database is located in `%AppData%\AgentBoard`. |

If the error is not described, collect the exact text of the message, the `agentboard version` and Windows versions, and steps to retry. Do not publish private project data or database contents in an open issue.
