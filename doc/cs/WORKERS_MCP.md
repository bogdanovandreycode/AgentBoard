# Pracovníci a připojení MCP

## Krok 1. Vytvořte pracovníka

V AgentBoard otevřete **Workers → Add worker**. Vyberte profil klienta: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama přes OpenCode, VS Code Copilot nebo jiný klient MCP. Profil vyplňuje typické vlastnosti (kód, testy, Git); můžete je změnit. Název je viditelný na desce a `Slug` je krátký identifikátor bez mezer pro příkaz MCP. Klikněte na **Uložit**.

Jeden pracovník odpovídá jedné osobnosti AI. Vytvořte různé pracovníky pro různé klienty nebo týmy. Profil schopností popisuje specializaci, ale nerozšiřuje práva AI na fáze úkolů.

## Krok 2. Zkopírujte konfiguraci

Otevřete vytvořeného pracovníka. Blok **Diagnostika MCP** zobrazuje fragment konfigurace a soubor, kam jej přidat. Klikněte na **Kopírovat konfiguraci MCP**. Pokud soubor již existuje, přidejte navrhovaný server ke stávajícímu objektu `mcpServers`/`servers`/`mcp` bez vymazání ostatních serverů.

Hlavní příkaz vypadá takto:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` by měl ukazovat na složku, kterou jste zaregistrovali prostřednictvím `agentboard init`. `--worker` — `Slug` vytvořeného pracovníka. Klient MCP spustí tento příkaz sám, když potřebuje nástroje. V prohlížeči AgentBoard může webový server běžet samostatně.

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

Pro Codex a Claude Code, příklad s projektem `C:\Projects\MyFirstProject` a pracovníkem `codex`:

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

V JSON je zpětné lomítko Windows zdvojené; hotový fragment z rozhraní to provede automaticky. Pokud používáte `--db` s vlastní základnou, přidejte ji do konfigurace `args` MCP a zadejte stejnou cestu jako při spouštění `open`.

**Ollama** poskytuje místní model, ale nenahrazuje klienta MCP. Profil **Ollama via OpenCode** generuje konfiguraci MCP pro OpenCode; samostatně nakonfigurovat OpenCode na modelu Ollama. Vhodný je také další klient kompatibilní s MCP s Ollama.

## Krok 4: Zkontrolujte připojení

1. Otevřete kartu pracovníka a klikněte na **Zkontrolovat znovu**. **Kontrola serveru** by měla ukazovat počet nástrojů MCP. Jedná se o interní kontrolu protokolu a detekci serverových nástrojů.
2. Po přidání konfiguračního souboru spusťte nebo restartujte klienta AI. Požádejte ho, aby zavolal na `get_my_board`.
3. V AgentBoard se objeví **Klient připojen** a v seznamu se objeví nová relace. Pouze to potvrdí připojení vašeho klienta. Pokud existují nástroje, ale klient není připojen, zkontrolujte cestu k programu, název konfiguračního souboru a jeho syntaxi JSON/TOML.

Začněte svou pracovní relaci s `get_my_board`. AI nepřijímá úkoly od `Backlog` a `Complete`, i když zná jejich ID. Umělá inteligence nemůže předstírat, že je člověk, a nemá obecný příkaz „přesunout se kamkoli“. Práce se systémem souborů mimo AgentBoard závisí na možnostech a karanténě vybraného klienta.

Oficiální pokyny pro zákazníky: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
