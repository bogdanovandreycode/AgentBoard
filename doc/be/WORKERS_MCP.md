# Воркеры і падлучэнне MCP

## Крок 1. Стварыце воркера

У AgentBoard адкрыйце **Workers → Add worker**. Абярыце профіль кліента: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot або іншы MCP-кліент. Профіль запаўняе тыповыя магчымасці (код, тэсты, Git); вы можаце іх памяняць. Імя бачна на дошцы, а `Slug` – кароткі ідэнтыфікатар без прабелаў для каманды MCP. Націсніце **Save**.

Адзін воркер адпавядае адной AI-асобы. Для розных кліентаў або каманд стварыце розных воркераў. Профіль магчымасцяў апісвае спецыялізацыю, але не пашырае правы AI на этапы задач.

## Крок 2. Скапіюйце канфігурацыю

Адкрыйце створанага воркера. У блоку **MCP diagnostics** паказаны канфігурацыйны фрагмент і файл, куды яго дадаць. Націсніце **Copy MCP config**. Калі файл ужо існуе, дадайце прапанаваны сервер у наяўны аб'ект `mcpServers`/`servers`/`mcp`, не сціраючы іншыя серверы.

Галоўная каманда выглядае так:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` павінен паказваць на тэчку, якую вы зарэгістравалі праз `agentboard init`. `--worker` - `Slug` створанага воркера. MCP-кліент запускае гэтую каманду сам, калі яму патрэбны інструменты. У браўзэры AgentBoard вэб-сервер можа працаваць асобна.

У блоку дыягностыкі пасля праверкі паказаны абсалютны шлях да запушчанага `agentboard.exe`. Ён асабліва карысны пры ўсталёўцы з ZIP. Пры ўсталёўцы праз Scoop можна выкарыстоўваць каманду `agentboard`, калі кліент бачыць той жа `PATH`.

## Крок 3. Дадайце сервер у свой кліент

У інтэрфейсе воркера ўжо ёсць гатовы фрагмент. Ніжэй тлумачэнне, дзе ён выкарыстоўваецца:

| Кліент | Куды ўставіць | Як праверыць на баку кліента
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, секцыя `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` у тэчцы праекта | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, аб'ект `mcpServers` | `/mcp list` у Gemini CLI |
| Cursor | `.cursor\mcp.json` праекту | спіс MCP-сервераў у наладах Cursor |
| OpenCode | `opencode.json` праекту, аб'ект `mcp` | спіс MCP-інструментаў у OpenCode |
| VS Code Copilot | `.vscode\mcp.json` праекту, аб'ект `servers` | каманда **MCP: List Servers** |

Для Codex і Claude Code прыклад з праектам `C:\Projects\MyFirstProject` і воркерам `codex`:

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

У JSON зваротная касая рыса Windows падвойваецца; гатовы фрагмент з інтэрфейсу робіць гэта аўтаматычна. Калі карыстаецеся `--db` з нестандартнай базай, дадайце яго ў `args` MCP-канфігурацыі і пакажыце той жа шлях, што пры запуску `open`.

**Ollama** дае лакальную мадэль, але не замяняе MCP-кліент. Профіль **Ollama via OpenCode** генеруе MCP-настройку для OpenCode; асобна наладзьце OpenCode на мадэль Ollama. Іншы MCP-сумяшчальны кліент з Ollama таксама падыходзіць.

## Крок 4. Праверце падлучэнне

1. Адкрыйце картку воркера і націсніце **Check again**. **Server check** павінен паказаць лік MCP-інструментаў. Гэта ўнутраная праверка пратаколу і выяўленні прылад сервера.
2. Запусціце або перазапусціце AI-кліент пасля дадання файла канфігурацыі. Папытаеце яго выклікаць `get_my_board`.
3. У AgentBoard з'явіцца **Client connected** і новая сесія ў спісе. Толькі гэта пацвярджае падлучэнне менавіта вашага кліента. Калі ёсць прылады, але кліент не падлучаны, праверце шлях да праграмы, імя файла канфігурацыі і яго JSON/TOML-сінтаксіс.

Пачынайце працоўную сесію з `get_my_board`. AI не атрымлівае задачы з `Backlog` і `Complete`, нават калі ведае іх ID. AI не можа прыкінуцца чалавекам і не мае агульнай каманды "перанесці куды заўгодна". Праца з файлавай сістэмай па-за AgentBoard залежыць ад магчымасцяў і пясочніцы абранага кліента.

Афіцыйныя інструкцыі кліентаў: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
