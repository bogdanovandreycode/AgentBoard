# Installation and launch on Windows

## Installer (recommended)

Download `agentboard-VERSION-windows-amd64-setup.exe` from the [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases) page. The installation wizard will suggest a folder; the default is `C:\AI\AgentBoard`. It will copy `agentboard.exe` and documentation, create a shortcut and add the selected folder to the system `PATH`. After installation, open a new terminal so that the `agentboard` command becomes available.

In your project folder run:

```powershell
agentboard init
agentboard open
```

The interface is built into `agentboard.exe`; No separate installation of Go or Node.js is needed. The data is stored in `%AppData%\AgentBoard` and is retained when the program is updated or uninstalled. Uninstalling via "Installed Applications" removes shortcuts and the entry from `PATH`.

## Installing Scoop

If Scoop is not already installed, open PowerShell as your regular user and follow the [official Scoop instructions](https://scoop.sh/). If there are restrictions on your corporate computer, contact your administrator; AgentBoard can also be launched from a ZIP without Scoop.

Alternatively install AgentBoard via Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

The presence of manifest in GitHub Release can be checked on the [Releases page](https://github.com/bogdanovandreycode/AgentBoard/releases). After adding manifest to Scoop, the bucket can be installed by bucket name and updated with the `scoop update agentboard` command.

## ZIP without Scoop

Download `agentboard-VERSION-windows-amd64.zip` from Releases, extract, for example, to `C:\Tools\AgentBoard`. In PowerShell in the project folder:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

For an AI client, specify the full path to `agentboard.exe` in its MCP configuration if the program is not in `PATH`.

## Build from source

Install the Go version from `go.mod` and Node.js 22 or later. In PowerShell at the root of the repository:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` builds the frontend into `internal/webui/dist`, runs Go tests and builds one `agentboard.exe`. The order is important: the interface is built into the binary when Go is built. If the old `agentboard.exe open` is running, stop it before rebuilding (Ctrl+C), otherwise Windows will not allow you to replace the file.

## Where is the data

- `%AppData%\AgentBoard\agentboard.db` - tasks, projects, workers and settings. You can specify a different file with the `--db` flag, but `init`, `open`/`serve` and `mcp` must have the **same path**.
- `<your project>\.agentboard\project.json` - project identifier. This file does not contain tasks.
- The server only listens to `127.0.0.1:7337` by default. Specify another address `--addr` before the project path: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` opens the project page and reuses an already running server at that address. If several projects are open in one browser, select them through the list on the left.
