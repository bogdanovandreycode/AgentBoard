# Trabalhadores e conexão MCP

## Etapa 1. Crie um trabalhador

No AgentBoard, abra **Trabalhadores → Adicionar trabalhador**. Selecione um perfil de cliente: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot ou outro cliente MCP. O perfil preenche recursos típicos (código, testes, Git); você pode alterá-los. O nome fica visível na placa e `Slug` é um identificador curto sem espaços para o comando MCP. Clique em **Salvar**.

Um trabalhador corresponde a uma personalidade de IA. Crie diferentes trabalhadores para diferentes clientes ou equipes. O perfil de capacidade descreve a especialização, mas não estende os direitos da IA ​​aos estágios da tarefa.

## Etapa 2. Copie a configuração

Abra o trabalhador criado. O bloco **MCP diagnostics** mostra o fragmento de configuração e o arquivo onde adicioná-lo. Clique em **Copiar configuração do MCP**. Se o arquivo já existir, adicione o servidor sugerido ao objeto `mcpServers`/`servers`/`mcp` existente sem apagar os outros servidores.

O comando principal é assim:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` deve apontar para a pasta que você registrou via `agentboard init`. `--worker` — `Slug` do trabalhador criado. O próprio cliente MCP executa esse comando quando precisa de ferramentas. No navegador AgentBoard, o servidor web pode ser executado separadamente.

Após a verificação, o bloco de diagnóstico mostra o caminho absoluto para o `agentboard.exe` em execução. É especialmente útil ao instalar a partir de um ZIP. Ao instalar via Scoop, você pode usar o comando `agentboard` se o cliente vir o mesmo `PATH`.

## Etapa 3: Adicione um servidor ao seu cliente

A interface do trabalhador já possui um fragmento pronto. Abaixo está uma explicação de onde ele é usado:

| Cliente | Onde inserir | Como verificar do lado do cliente |
| --- | --- | --- |
| Códice | `%USERPROFILE%\.codex\config.toml`, seção `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Código Cláudio | `.mcp.json` na pasta do projeto | `claude mcp list` |
| Gêmeos CLI | `%USERPROFILE%\.gemini\settings.json`, objeto `mcpServers` | `/mcp list` na CLI Gemini |
| Cursor | Projeto `.cursor\mcp.json` | lista de servidores MCP nas configurações do Cursor |
| Código aberto | Projeto `opencode.json`, objeto `mcp` | lista de ferramentas MCP em OpenCode |
| Copiloto do Código VS | Projeto `.vscode\mcp.json`, objeto `servers` | comando **MCP: listar servidores** |

Para Codex e Claude Code, um exemplo com o projeto `C:\Projects\MyFirstProject` e o trabalhador `codex`:

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

Em JSON, a barra invertida do Windows é duplicada; um fragmento pronto da interface faz isso automaticamente. Se você estiver usando `--db` com uma base personalizada, adicione-o à configuração do MCP `args` e especifique o mesmo caminho ao iniciar `open`.

**Ollama** fornece um modelo local, mas não substitui o cliente MCP. O perfil **Ollama via OpenCode** gera uma configuração MCP para OpenCode; configure separadamente o OpenCode no modelo Ollama. Outro cliente compatível com MCP com Ollama também é adequado.

## Etapa 4: verifique sua conexão

1. Abra o cartão de trabalhador e clique em **Verificar novamente**. **Verificação do servidor** deve mostrar o número de ferramentas MCP. Esta é uma verificação de protocolo interno e detecção de ferramentas de servidor.
2. Inicie ou reinicie o cliente AI após adicionar o arquivo de configuração. Peça a ele para ligar para `get_my_board`.
3. **Cliente conectado** aparecerá no AgentBoard e uma nova sessão aparecerá na lista. Só isso confirma a conexão do seu cliente. Se houver ferramentas, mas o cliente não estiver conectado, verifique o caminho do programa, o nome do arquivo de configuração e sua sintaxe JSON/TOML.

Comece sua sessão de trabalho com `get_my_board`. AI não recebe tarefas de `Backlog` e `Complete`, mesmo que conheça seu ID. A IA não pode fingir ser humana e não possui um comando geral de “mover para qualquer lugar”. Trabalhar com o sistema de arquivos fora do AgentBoard depende dos recursos e da área restrita do cliente selecionado.

Instruções oficiais do cliente: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
