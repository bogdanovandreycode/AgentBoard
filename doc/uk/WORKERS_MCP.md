# Воркери та підключення MCP

## Крок 1. Створіть воркера

В AgentBoard відкрийте **Workers → Add worker**. Виберіть профіль клієнта: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot або інший MCP-клієнт. Профіль заповнює типові можливості (код, тести, Git); Ви можете змінити їх. Ім'я видно на дошці, а `Slug` – короткий ідентифікатор без пропусків для команди MCP. Натисніть **Save**.

Один воркер відповідає одній AI-особі. Для різних клієнтів чи команд створіть різних воркерів. Профіль можливостей описує спеціалізацію, але з розширює права AI на етапи завдань.

## Крок 2. Скопіюйте конфігурацію

Відкрийте створений воркер. У блоці **MCP diagnostics** показано конфігураційний фрагмент та файл, куди його додати. Натисніть **Copy MCP config**. Якщо файл вже існує, додайте запропонований сервер до існуючого об'єкта `mcpServers`/`servers`/`mcp`, не стираючи інші сервери.

Головна команда виглядає так:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` повинен вказувати на папку, яку ви зареєстрували через `agentboard init`. `--worker` - `Slug` створеного воркера. MCP-клієнт запускає цю команду сам, коли йому потрібні інструменти. У браузері AgentBoard веб-сервер може працювати окремо.

У блоці діагностики після перевірки вказано абсолютний шлях до запущеного `agentboard.exe`. Він особливо корисний при встановленні із ZIP. При установці через Scoop можна використовувати команду `agentboard`, якщо клієнт бачить той самий `PATH`.

## Крок 3. Додайте сервер до свого клієнта

В інтерфейсі воркера вже готовий фрагмент. Нижче пояснення, де він використовується:

| Клієнт | куди вставити | Як перевірити на стороні клієнта
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, секція `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` у папці проекту | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, об'єкт `mcpServers` | `/mcp list` у Gemini CLI |
| Cursor | Проекти `.cursor\mcp.json` | список MCP-серверів у налаштуваннях Cursor |
| OpenCode | `opencode.json` проекту, об'єкт `mcp` | список MCP-інструментів у OpenCode |
| VS Code Copilot | проекту `.vscode\mcp.json`, об'єкт `servers` | команда **MCP: List Servers** |

Для Codex та Claude Code приклад із проектом `C:\Projects\MyFirstProject` та воркером `codex`:

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

У JSON зворотна коса характеристика Windows подвоюється; готовий фрагмент з інтерфейсу робить це автоматично. Якщо Ви використовуєте `--db` з нестандартною базою, додайте його до `args` MCP-конфігурації та вкажіть той самий шлях, що під час запуску `open`.

**Ollama** надає локальну модель, але не замінює MCP-клієнт. Профіль **Ollama via OpenCode** генерує MCP-налаштування для OpenCode; окремо налаштуйте OpenCode на Ollama. Інший MCP-сумісний клієнт із Ollama теж підходить.

## Крок 4. Перевірте підключення

1. Відкрийте картку воркера та натисніть **Check again**. **Server check** має показати число MCP-інструментів. Це внутрішня перевірка протоколу та виявлення інструментів сервера.
2. Запустіть або перезапустіть AI-клієнт після додавання конфігураційного файлу. Попросіть виклик виклику `get_my_board`.
3. В AgentBoard з'явиться **Client connected** та нова сесія у списку. Тільки це підтверджує підключення вашого клієнта. Якщо є інструменти, але клієнт не підключений, перевірте шлях до програми, ім'я конфігураційного файлу та його JSON/TOML-синтаксис.

Починайте робочу сесію із `get_my_board`. AI не отримує задач з `Backlog` і `Complete`, навіть якщо знає їх ID. AI не може прикинутися людиною і не має спільної команди "перенести куди завгодно". Робота з файловою системою поза AgentBoard залежить від можливостей та пісочниці обраного клієнта.

Офіційні інструкції клієнтів: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
