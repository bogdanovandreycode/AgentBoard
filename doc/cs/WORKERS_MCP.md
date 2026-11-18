# Pracovníci a připojení MCP

## Krok 1. Vytvořte pracovníka

V AgentBoard otevřete **Workers → Add worker**. Vyberte profil klienta: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama přes OpenCode, VS Code Copilot nebo jiný klient MCP. Profil vyplňuje typické vlastnosti (kód, testy, Git); můžete je změnit. Název je viditelný na desce a `Slug` je krátký identifikátor bez mezer pro příkaz MCP. Klikněte na **Uložit**.

Jeden pracovník odpovídá jedné osobnosti AI. Vytvořte různé pracovníky pro různé klienty nebo týmy. Profil schopností popisuje specializaci, ale nerozšiřuje práva AI na fáze úkolů.

## Krok 2. Zkopírujte konfiguraci

Otevřete vytvořeného pracovníka. Blok **Diagnostika MCP** zobrazuje fragment konfigurace a soubor, kam jej přidat. Klikněte na **Kopírovat konfiguraci MCP**. Pokud soubor již existuje, přidejte navrhovaný server ke stávajícímu objektu `mcpServers`/`servers`/`mcp` bez vymazání ostatních serverů.

Hlavní příkaz vypadá takto:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` vytvořeného pracovníka. Klient MCP spustí tento příkaz sám, když potřebuje nástroje.

### MCP funguje bez webového rozhraní

`agentboard mcp` je samostatný proces stdio, který sám otevírá databázi AgentBoard. Nepotřebuje `agentboard open`/`serve`, nezáleží na složce, ze které se webové rozhraní spouští, a na tom, kolik instancí webového rozhraní běží. Více agentů může běžet současně, každý ve svém vlastním projektu: každý klient provozuje svůj vlastní proces MCP a všichni bezpečně používají stejnou databázi SQLite.

### Jak se vybírá projekt

Server MCP běží v libovolné složce. Když poprvé zavoláte nástroj, je vázán na projekt:

1. `--project PATH`, pokud je zadán (vhodná je jakákoli složka v registrovaném projektu);
2. jinak pracovní složky (kořeny) klienta, pokud je klient nahlásí;
3. jinak složka, ve které klient spustil MCP server.

`get_my_board` vrátí vybraný projekt a pracovníka do `session` spolu s `binding_source`. Pokud se pracovní složky klienta změní, další volání nástroje provede vazbu znovu. Bez `--worker` je použit jediný povolený pracovník projektu.

Pokud vazba není možná, každý nástroj vrátí chybu, kterou musí agent předat uživateli; data úlohy se v tomto případě nepřečtou:

| Kód | Význam |
| --- | --- |
| `NO_PROJECTS` | V databázi nejsou žádné projekty. Spusťte `agentboard init` ve složce projektu a vytvořte pracovníka. Není potřeba restart. |
| `PROJECT_REQUIRED` | Pracovní složka je mimo registrované projekty. Zpráva obsahuje seznam projektů. |
| `PROJECT_NOT_FOUND` | `--project` není v registrovaném projektu. |
| `PROJECT_AMBIGUOUS` | Pracovní složky klienta zahrnují více projektů. Přidejte `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Vytvořte, povolte nebo vyberte pracovníka prostřednictvím `--worker`. |

Konfigurace pro **Codex** a **Gemini CLI** jsou vlastní a sdílené ve všech složkách, takže nemají `--project`: Vytvořte pracovníka se stejným slug v každém projektu, kde tento klient běží. Soubory na úrovni projektu (Claude Code, Cursor, OpenCode, VS Code) obsahují `--project` jejich projektu.

### Explicitní návrh a kontrola kontextu

Pokud je zadáno `--project` a spouštěcí nebo pracovní složka klienta je v jiném projektu, nástroje požádají o potvrzení prostřednictvím vyvolání MCP. Klienti bez vyvolání formuláře obdrží `PROJECT_CONTEXT_MISMATCH` před přístupem k úlohám. Opravte konfiguraci, odeberte `--project` pro automatickou detekci nebo výslovně potvrďte, že máte v úmyslu pracovat s jiným projektem s příznakem `--allow-project-mismatch`. Změna pracovních složek resetuje předchozí potvrzení. Automaticky definované projekty nevyžadují potvrzení.

**Globální nastavení → Kontext projektu MCP** řídí tyto kontroly pro všechny projekty a pracovníky stejné databáze. Ve výchozím nastavení povoleno. Pokud je zakázáno, MCP použije zadaný projekt bez potvrzení nebo chyb nesouladu. Již otevřené relace použijí nastavení při příštím volání nástroje; Viditelnost úkolu, přechody a pravidla testování na lidech se nadále dodržují.

Po kontrole diagnostický blok ukazuje absolutní cestu k běžícímu `agentboard.exe`. Je to užitečné zejména při instalaci ze ZIP. Při instalaci přes Scoop můžete použít příkaz `agentboard`, pokud klient vidí stejný `PATH`.

## Krok 3: Přidejte server ke svému klientovi

Pracovní rozhraní již má připravený fragment. Níže je vysvětlení, kde se používá:

| Klient | Kam vložit | Jak zkontrolovat na straně klienta |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, sekce `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` ve složce projektu | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, objekt `mcpServers` | `/mcp list` v Gemini CLI |
| Kurzor | Projekt `.cursor\mcp.json` | seznam MCP serverů v nastavení kurzoru |
| OpenCode | Projekt `opencode.json`, objekt `mcp` | seznam nástrojů MCP v OpenCode |
| VS Code Copilot | Projekt `.vscode\mcp.json`, objekt `servers` | příkaz **MCP: Seznam serverů** |

Příklady pro Codex (uživatelská konfigurace, projekt detekován automaticky) a Claude Code (soubor projektu) s projektem `C:\Projects\MyFirstProject` a pracovníkem `codex`:

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

V JSON je zpětné lomítko Windows zdvojené; hotový fragment z rozhraní to provede automaticky. Vygenerované konfigurace obsahují skutečné `--db` webového serveru, takže HTTP a MCP používají stejný základ, i když je klientské prostředí odlišné.

**Ollama** poskytuje místní model, ale nenahrazuje klienta MCP. Profil **Ollama via OpenCode** generuje konfiguraci MCP pro OpenCode; samostatně nakonfigurovat OpenCode na modelu Ollama. Vhodný je také další klient kompatibilní s MCP s Ollama.

## Krok 4: Zkontrolujte připojení

1. Otevřete kartu pracovníka a klikněte na **Zkontrolovat znovu**. **Kontrola serveru** spouští skutečný stdio spustitelný soubor s tímto projektem, pracovníkem a základnou, kontroluje handshake, seznam nástrojů a že se projekt a pracovník nacházejí.
2. Po přidání konfiguračního souboru spusťte nebo restartujte klienta AI. Požádejte ho, aby zavolal na `get_my_board`.
3. V AgentBoard se objeví **Klient připojen** a v seznamu se objeví nová relace. Pouze to potvrdí připojení vašeho klienta. Pokud existují nástroje, ale klient není připojen, zkontrolujte cestu k programu, název konfiguračního souboru a jeho syntaxi JSON/TOML.

Začněte svou pracovní relaci s `get_my_board`. AI nepřijímá úkoly od `Backlog` a `Complete`, i když zná jejich ID. Umělá inteligence nemůže předstírat, že je člověk, a nemá obecný příkaz „přesunout se kamkoli“. Práce se systémem souborů mimo AgentBoard závisí na možnostech a karanténě vybraného klienta.

Každý nástroj deklaruje `title` a všechny čtyři rady pro chování MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) jako explicitní booleovské hodnoty. Nástroje pro čtení jsou pouze pro čtení a idempotentní; nahrávací nástroje pouze přidávají záznamy nebo posouvají úkoly dopředu, takže žádné nejsou rušivé ani nepřesahují místní základnu AgentBoard.

Oficiální pokyny pro zákazníky: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
