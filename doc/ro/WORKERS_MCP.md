# Lucrători și conexiune MCP

## Pasul 1. Creați un lucrător

În AgentBoard, deschideți **Workers → Add worker**. Selectați un profil de client: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama prin OpenCode, VS Code Copilot sau alt client MCP. Profilul completează caracteristici tipice (cod, teste, Git); le poti schimba. Numele este vizibil pe placă, iar `Slug` este un identificator scurt fără spații pentru comanda MCP. Faceți clic pe **Salvați**.

Un lucrător corespunde unei personalități AI. Creați lucrători diferiți pentru clienți sau echipe diferiți. Profilul de capacitate descrie specializarea, dar nu extinde drepturile AI la etapele sarcinii.

## Pasul 2. Copiați configurația

Deschideți lucrătorul creat. Blocul **Diagnoza MCP** arată fragmentul de configurare și fișierul unde să îl adăugați. Faceți clic pe **Copiați configurația MCP**. Dacă fișierul există deja, adăugați serverul sugerat la obiectul existent `mcpServers`/`servers`/`mcp` fără a șterge celelalte servere.

Comanda principală arată astfel:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` ar trebui să indice folderul pe care l-ați înregistrat prin `agentboard init`. `--worker` — `Slug` al lucrătorului creat. Clientul MCP rulează singur această comandă atunci când are nevoie de instrumente. În browserul AgentBoard, serverul web poate rula separat.

După verificare, blocul de diagnosticare arată calea absolută către `agentboard.exe` care rulează. Este util mai ales când se instalează dintr-un ZIP. Când instalați prin Scoop, puteți utiliza comanda `agentboard` dacă clientul vede același `PATH`.

## Pasul 3: Adăugați un server la clientul dvs

Interfața de lucru are deja un fragment gata făcut. Mai jos este o explicație a locului în care este utilizat:

| Client | Unde se introduce | Cum se verifică pe partea clientului |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, secțiunea `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Cod | `.mcp.json` în folderul de proiect | `claude mcp list` |
| Gemeni CLI | `%USERPROFILE%\.gemini\settings.json`, obiect `mcpServers` | `/mcp list` în Gemini CLI |
| Cursor | Proiect `.cursor\mcp.json` | lista de servere MCP în Setări cursor |
| OpenCode | Proiect `opencode.json`, obiect `mcp` | lista de instrumente MCP în OpenCode |
| Copilotul cod VS | Proiect `.vscode\mcp.json`, obiect `servers` | comanda **MCP: Listă servere** |

Pentru Codex și Claude Code, un exemplu cu proiectul `C:\Projects\MyFirstProject` și lucrătorul `codex`:

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

În JSON, bara oblică inversă Windows este dublată; un fragment gata făcut din interfață face acest lucru automat. Dacă utilizați `--db` cu o bază personalizată, adăugați-o la configurația MCP `args` și specificați aceeași cale ca la pornirea `open`.

**Ollama** oferă un model local, dar nu înlocuiește clientul MCP. Profilul **Ollama prin OpenCode** generează o configurație MCP pentru OpenCode; configurați separat OpenCode pe modelul Ollama. Un alt client compatibil MCP cu Ollama este de asemenea potrivit.

## Pasul 4: Verificați-vă conexiunea

1. Deschideți cardul de lucrător și faceți clic pe **Verificați din nou**. **Verificarea serverului** ar trebui să arate numărul de instrumente MCP. Aceasta este o verificare internă a protocolului și o detectare a instrumentelor de server.
2. Porniți sau reporniți clientul AI după adăugarea fișierului de configurare. Rugați-l să sune la `get_my_board`.
3. **Client conectat** va apărea în AgentBoard și va apărea o nouă sesiune în listă. Doar asta confirmă conexiunea clientului tău. Dacă există instrumente, dar clientul nu este conectat, verificați calea către program, numele fișierului de configurare și sintaxa acestuia JSON/TOML.

Începeți sesiunea de lucru cu `get_my_board`. AI nu primește sarcini de la `Backlog` și `Complete`, chiar dacă le cunoaște ID-ul. AI nu poate pretinde a fi uman și nu are o comandă generală de „mutare oriunde”. Lucrul cu sistemul de fișiere în afara AgentBoard depinde de capabilitățile și sandbox-ul clientului selectat.

Instrucțiuni oficiale pentru clienți: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
