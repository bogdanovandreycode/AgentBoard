# Воркеры и подключение MCP

## Шаг 1. Создайте воркера

В AgentBoard откройте **Workers → Add worker**. Выберите профиль клиента: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot или другой MCP-клиент. Профиль заполняет типичные возможности (код, тесты, Git); вы можете их поменять. Имя видно на доске, а `Slug` — короткий идентификатор без пробелов для команды MCP. Нажмите **Save**.

Один воркер соответствует одной AI-личности. Для разных клиентов или команд создайте разных воркеров. Профиль возможностей описывает специализацию, но не расширяет права AI на этапы задач.

## Шаг 2. Скопируйте конфигурацию

Откройте созданного воркера. В блоке **MCP diagnostics** показан конфигурационный фрагмент и файл, куда его добавить. Нажмите **Copy MCP config**. Если файл уже существует, добавьте предложенный сервер в существующий объект `mcpServers`/`servers`/`mcp`, не стирая другие серверы.

Главная команда выглядит так:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` должен указывать на папку, которую вы зарегистрировали через `agentboard init`. `--worker` — `Slug` созданного воркера. MCP-клиент запускает эту команду сам, когда ему нужны инструменты. В браузере AgentBoard веб-сервер может работать отдельно.

В блоке диагностики после проверки указан абсолютный путь к запущенному `agentboard.exe`. Он особенно полезен при установке из ZIP. При установке через Scoop можно использовать команду `agentboard`, если клиент видит тот же `PATH`.

## Шаг 3. Добавьте сервер в свой клиент

В интерфейсе воркера уже есть готовый фрагмент. Ниже пояснение, где он используется:

| Клиент | Куда вставить | Как проверить на стороне клиента |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, секция `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` в папке проекта | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, объект `mcpServers` | `/mcp list` в Gemini CLI |
| Cursor | `.cursor\mcp.json` проекта | список MCP-серверов в настройках Cursor |
| OpenCode | `opencode.json` проекта, объект `mcp` | список MCP-инструментов в OpenCode |
| VS Code Copilot | `.vscode\mcp.json` проекта, объект `servers` | команда **MCP: List Servers** |

Для Codex и Claude Code пример с проектом `C:\Projects\MyFirstProject` и воркером `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
```

```json
// Содержимое .mcp.json (уберите эту строку-комментарий при копировании)
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

В JSON обратная косая черта Windows удваивается; готовый фрагмент из интерфейса делает это автоматически. Если используете `--db` с нестандартной базой, добавьте его в `args` MCP-конфигурации и укажите тот же путь, что при запуске `open`.

**Ollama** предоставляет локальную модель, но не заменяет MCP-клиент. Профиль **Ollama via OpenCode** генерирует MCP-настройку для OpenCode; отдельно настройте OpenCode на модель Ollama. Другой MCP-совместимый клиент с Ollama тоже подходит.

## Шаг 4. Проверьте подключение

1. Откройте карточку воркера и нажмите **Check again**. **Server check** должен показать число MCP-инструментов. Это внутренняя проверка протокола и обнаружения инструментов сервера.
2. Запустите или перезапустите AI-клиент после добавления файла конфигурации. Попросите его вызвать `get_my_board`.
3. В AgentBoard появится **Client connected** и новая сессия в списке. Только это подтверждает подключение именно вашего клиента. Если есть инструменты, но клиент не подключён, проверьте путь к программе, имя файла конфигурации и его JSON/TOML-синтаксис.

Начинайте рабочую сессию с `get_my_board`. AI не получает задачи из `Backlog` и `Complete`, даже если знает их ID. AI не может притвориться человеком и не имеет общей команды «перенести куда угодно». Работа с файловой системой вне AgentBoard зависит от возможностей и песочницы выбранного клиента.

Официальные инструкции клиентов: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
