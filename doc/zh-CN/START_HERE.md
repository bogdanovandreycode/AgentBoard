# 从这里开始

## 什么是 AgentBoard

想象一下带有任务卡的普通板。您创建任务、指派专人负责并监督工作。 AI 工作人员仅通过 MCP 接收分配给他们的任务并报告进度。您可以决定任务何时最终准备好。

**项目** - 计算机上的文件夹和单独的板。 **任务** - 一张带有描述、负责人和阶段的卡片。 **Worker**是AI客户端的逻辑名称。 **MCP**是AI客户端连接AgentBoard的方式。 **Scoop** 是 Windows 的软件安装管理器。

无需代码。您需要 Windows、浏览器、PowerShell，并且为了让 AI 工作，还需要安装支持本地 MCP 服务器的 AI 客户端。

## 首次启动

1. 按照[说明](INSTALL.md).
2. 在资源管理器中创建一个项目文件夹，例如`C:\Projects\MyFirstProject`。
3. 在资源管理器中打开此文件夹。单击地址栏中的 ，输入 `powershell` 并按 Enter。
4. 在打开的窗口中，执行以下操作：

```powershell
agentboard init
agentboard open
```

5. 浏览器将打开并显示地址`http://127.0.0.1:7337`。使用白板时，让 PowerShell 窗口保持打开状态。关闭窗口将停止本地服务器，但任务将保留。

如果团队`agentboard`未找到，请关闭 PowerShell，然后在安装 Scoop 后重新打开。从 ZIP 安装时，请使用完整路径`agentboard.exe`。

## 第一个任务

单击“新建任务”，填写“标题”，如有必要，填写“描述”，然后“保存”。新任务在`Backlog`人类可以接触到的。为了让AI开始工作，分配一个工人并将任务转移给`Features`。字段的逐步描述 -[任务.md](TASKS.md)。

## 第一个工人

打开**工人→添加工人**。选择您正在使用的AI客户端（例如Codex或Claude Code），检查名称和短ID`Slug`，单击**保存**。打开创建的worker：有一个MCP配置，一个复制按钮和一个服务器检查。根据复制配置到AI客户端[WORKERS_MCP.md](WORKERS_MCP.md)。连接后，请客户致电`get_my_board`。

## 接下来要读什么

-[任务、测试、故事和专栏](TASKS.md)
-[连接工作人员并检查 MCP](WORKERS_MCP.md)
-[从 JSON 导入任务](IMPORT.md)
-[语言、主题、时区和更新](SETTINGS.md)
-[典型问题](TROUBLESHOOTING.md)
