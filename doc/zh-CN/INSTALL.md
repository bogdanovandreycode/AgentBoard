# 在 Windows 上安装并启动

## 安装程序（推荐）

从[Releases](https://github.com/bogdanovandreycode/AgentBoard/releases).zip]页面下载`agentboard-VERSION-windows-amd64-setup.exe`。安装向导会建议一个文件夹；默认为 `C:\AI\AgentBoard`。它将复制 `agentboard.exe` 和文档，创建快捷方式并将所选文件夹添加到系统 `PATH`。安装后，打开一个新终端，以便命令 `agentboard` 可用。

在您的项目文件夹中运行：

```powershell
agentboard init
agentboard open
```

该接口内置于`agentboard.exe`中；无需单独安装 Go 或 Node.js。数据存储在`%AppData%\AgentBoard`中，并在程序更新或卸载时保留。通过“已安装的应用程序”卸载会删除 `PATH` 中的快捷方式和条目。

## 安装勺子

如果尚未安装 Scoop，请以普通用户身份打开 PowerShell 并按照[官方 Scoop](https://scoop.sh/) 说明进行操作。如果您的公司计算机有限制，请联系您的管理员； AgentBoard 也可以从 ZIP 启动，无需 Scoop。

或者通过 Scoop 安装 AgentBoard：

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

GitHub Release 中清单的可用性可以在 [page Releases](https://github.com/bogdanovandreycode/AgentBoard/releases)] 上检查。将清单添加到 Scoop 后，可以通过存储桶名称安装存储桶，并使用 `scoop update agentboard` 命令更新存储桶。

## ZIP 不带勺子

从发行版下载 `agentboard-VERSION-windows-amd64.zip`，例如解压到 `C:\Tools\AgentBoard`。在项目文件夹中的 PowerShell 中：

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

对于 AI 客户端，如果程序不在 `agentboard.exe` 中，请在其 MCP 配置中指定 `PATH` 的完整路径。

## 从源代码构建

从 `go.mod` 和 Node.js 22 或更高版本安装 Go 版本。在存储库根目录的 PowerShell 中：

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` 将前端组装到 `internal/webui/dist` 中，运行 Go 测试并组装一个 `agentboard.exe`。顺序很重要：当 Go 构建时，接口就内置到二进制文件中。如果旧的 `agentboard.exe open` 正在运行，请在重建 (Ctrl+C) 之前停止它，否则 Windows 将不允许您替换该文件。

## 数据在哪里

- `%AppData%\AgentBoard\agentboard.db` - 任务、项目、工作人员和设置。您可以使用 `--db` 标志指定不同的文件，但 `init`、`open`/`serve` 和 `mcp` 必须具有**相同的路径**。
- `<ваш проект>\.agentboard\project.json` — 项目标识符。该文件不包含任务。
- 服务器默认只监听`127.0.0.1:7337`。在项目路径之前指定另一个地址`--addr`：`agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`。

`open` 打开项目页面并重用该地址处已运行的服务器。如果在一个浏览器中打开了多个项目，请通过左侧列表选择它们。
