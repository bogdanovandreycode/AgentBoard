# Lucrători și conexiune MCP

## Pasul 1. Creați un lucrător

În AgentBoard, deschideți **Workers → Add worker**. Selectați un profil de client: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama prin OpenCode, VS Code Copilot sau alt client MCP. Profilul completează caracteristici tipice (cod, teste, Git); le poti schimba. Numele este vizibil pe placă, iar `Slug` este un identificator scurt fără spații pentru comanda MCP. Faceți clic pe **Salvați**.

Un lucrător corespunde unei personalități AI. Creați lucrători diferiți pentru clienți sau echipe diferiți. Profilul de capacitate descrie specializarea, dar nu extinde drepturile AI la etapele sarcinii.

## Pasul 2. Copiați configurația

Deschideți lucrătorul creat. Blocul **Diagnoza MCP** arată fragmentul de configurare și fișierul unde să îl adăugați. Faceți clic pe **Copiați configurația MCP**. Dacă fișierul există deja, adăugați serverul sugerat la obiectul existent `mcpServers`/`servers`/`mcp` fără a șterge celelalte servere.

Comanda principală arată astfel:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` al lucrătorului creat. Clientul MCP rulează singur această comandă atunci când are nevoie de instrumente.

### MCP funcționează fără o interfață web

`agentboard mcp` este un proces stdio separat care deschide el însuși baza de date AgentBoard. Nu are nevoie de `agentboard open`/`serve`, nu depinde de folderul din care este lansată interfața web și de câte instanțe ale interfeței web rulează. Mai mulți agenți pot rula simultan, fiecare în propriul proiect: fiecare client rulează propriul său proces MCP și toți folosesc în siguranță aceeași bază de date SQLite.

### Cum este selectat un proiect

Serverul MCP rulează în orice folder. Prima dată când apelați instrumentul, acesta este legat de proiect:

1. `--project PATH`, dacă este specificat (orice folder din interiorul proiectului înregistrat este potrivit);
2. în caz contrar folderele de lucru (roots) ale clientului, dacă clientul le raportează;
3. altfel folderul în care clientul a lansat serverul MCP.

`get_my_board` returnează proiectul selectat și lucrătorul la `session` împreună cu `binding_source`. Dacă folderele de lucru ale clientului se modifică, următorul apel către instrument efectuează din nou legarea. Fără `--worker`, este utilizat singurul lucrător activat al proiectului.

Dacă legarea nu este posibilă, fiecare instrument returnează o eroare pe care agentul trebuie să o transmită utilizatorului; datele sarcinii nu sunt citite în acest caz:

| Cod | Înțeles |
| --- | --- |
| `NO_PROJECTS` | Nu există proiecte în baza de date. Executați `agentboard init` în folderul de proiect și creați un lucrător. Nu este nevoie de repornire. |
| `PROJECT_REQUIRED` | Dosarul de lucru se află în afara proiectelor înregistrate. Mesajul listează proiecte. |
| `PROJECT_NOT_FOUND` | `--project` nu se află într-un proiect înregistrat. |
| `PROJECT_AMBIGUOUS` | Dosarele de lucru ale unui client acoperă mai multe proiecte. Adăugați `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Creați, activați sau selectați un lucrător prin `--worker`. |

Configurațiile pentru **Codex** și **Gemini CLI** sunt personalizate și partajate în toate folderele, deci nu au `--project`: creați un lucrător cu același slug în fiecare proiect în care rulează acest client. Fișierele la nivel de proiect (Claude Code, Cursor, OpenCode, VS Code) conțin `--project` al proiectului lor.

### Design explicit și verificarea contextului

Dacă este specificat `--project` și folderul de pornire al clientului sau folderul de lucru se află într-un proiect diferit, instrumentele solicită confirmarea prin elicitarea MCP. Clienții fără elicitare de formular primesc `PROJECT_CONTEXT_MISMATCH` înainte de a accesa sarcini. Corectați configurația, eliminați `--project` pentru detectarea automată sau confirmați în mod explicit că intenționați să lucrați cu un alt proiect cu indicatorul `--allow-project-mismatch`. Schimbarea folderelor de lucru resetează confirmările anterioare. Proiectele definite automat nu necesită confirmare.

**Setări globale → contextul proiectului MCP** controlează aceste verificări pentru toate proiectele și lucrătorii aceleiași baze de date. Activat implicit. Dacă este dezactivat, MCP utilizează proiectul specificat fără confirmări sau erori de nepotrivire. Sesiunile deja deschise aplică setarea data viitoare când instrumentul este apelat; Vizibilitatea sarcinilor, tranzițiile și regulile de testare umană continuă să fie respectate.

După verificare, blocul de diagnosticare arată calea absolută către `agentboard.exe` care rulează. Este util mai ales când se instalează dintr-un ZIP. Când instalați prin Scoop, puteți utiliza comanda `agentboard` dacă clientul vede același `PATH`.

## Pasul 3: Adăugați un server la clientul dvs

Interfața de lucru are deja un fragment gata făcut. Mai jos este o explicație a locului în care este utilizat:

| Client | Unde se introduce | Cum se verifică pe partea clientului |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, secțiunea `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Cod | `.mcp.json` în folderul de proiect | `claude mcp list` |
| Gemeni CLI | `%USERPROFILE%\.gemini\settings.json`, obiect `mcpServers` | `/mcp list` în Gemeni CLI |
| Cursor | Proiect `.cursor\mcp.json` | lista de servere MCP în Setări cursor |
| OpenCode | Proiect `opencode.json`, obiect `mcp` | lista de instrumente MCP în OpenCode |
| Copilotul cod VS | Proiect `.vscode\mcp.json`, obiect `servers` | comanda **MCP: Listă servere** |

Exemple pentru Codex (configurație utilizator, proiect detectat automat) și Claude Code (fișier proiect) cu proiectul `C:\Projects\MyFirstProject` și lucrător `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
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

În JSON, bara oblică inversă Windows este dublată; un fragment gata făcut din interfață face acest lucru automat. Configurațiile generate conțin `--db` real al serverului web, astfel încât HTTP și MCP folosesc aceeași bază, chiar dacă mediul client este diferit.

**Ollama** oferă un model local, dar nu înlocuiește clientul MCP. Profilul **Ollama prin OpenCode** generează o configurație MCP pentru OpenCode; configurați separat OpenCode pe modelul Ollama. Un alt client compatibil MCP cu Ollama este de asemenea potrivit.

## Pasul 4: Verificați-vă conexiunea

1. Deschideți cardul de lucrător și faceți clic pe **Verificați din nou**. **Verificare server** rulează un executabil stdio real cu acest proiect, lucrător și bază, verifică strângerea de mână, lista de instrumente și dacă proiectul și lucrătorul sunt localizați.
2. Porniți sau reporniți clientul AI după adăugarea fișierului de configurare. Rugați-l să sune la `get_my_board`.
3. **Client conectat** va apărea în AgentBoard și va apărea o nouă sesiune în listă. Doar asta confirmă conexiunea clientului tău. Dacă există instrumente, dar clientul nu este conectat, verificați calea către program, numele fișierului de configurare și sintaxa acestuia JSON/TOML.

Începeți sesiunea de lucru cu `get_my_board`. AI nu primește sarcini de la `Backlog` și `Complete`, chiar dacă le cunoaște ID-ul. AI nu poate pretinde a fi uman și nu are o comandă generală de „mutare oriunde”. Lucrul cu sistemul de fișiere în afara AgentBoard depinde de capabilitățile și sandbox-ul clientului selectat.

Fiecare instrument declară `title` și toate cele patru indicii de comportament MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) ca valori booleene explicite. Instrumentele de citire sunt numai pentru citire și idempotente; instrumentele de înregistrare doar adaugă înregistrări sau mută sarcini înainte, astfel încât niciuna nu perturbă sau nu depășește baza locală AgentBoard.

Instrucțiuni oficiale pentru clienți: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
