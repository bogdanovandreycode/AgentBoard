# Trabalhadores e conexão MCP

## Etapa 1. Crie um trabalhador

No AgentBoard, abra **Trabalhadores → Adicionar trabalhador**. Selecione um perfil de cliente: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot ou outro cliente MCP. O perfil preenche recursos típicos (código, testes, Git); você pode alterá-los. O nome fica visível na placa e `Slug` é um identificador curto sem espaços para o comando MCP. Clique em **Salvar**.

Um trabalhador corresponde a uma personalidade de IA. Crie diferentes trabalhadores para diferentes clientes ou equipes. O perfil de capacidade descreve a especialização, mas não estende os direitos da IA ​​aos estágios da tarefa.

## Etapa 2. Copie a configuração

Abra o trabalhador criado. O bloco **MCP diagnostics** mostra o fragmento de configuração e o arquivo onde adicioná-lo. Clique em **Copiar configuração do MCP**. Se o arquivo já existir, adicione o servidor sugerido ao objeto `mcpServers`/`servers`/`mcp` existente sem apagar os outros servidores.

O comando principal é assim:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` do trabalhador criado. O próprio cliente MCP executa esse comando quando precisa de ferramentas.

### MCP funciona sem interface web

`agentboard mcp` é um processo stdio separado que abre o banco de dados AgentBoard. Ele não precisa de `agentboard open`/`serve`, não depende da pasta a partir da qual a interface da web é iniciada e de quantas instâncias da interface da web estão em execução. Vários agentes podem ser executados simultaneamente, cada um em seu próprio projeto: cada cliente executa seu próprio processo MCP e todos usam com segurança o mesmo banco de dados SQLite.

### Como um projeto é selecionado

O servidor MCP é executado em qualquer pasta. Na primeira vez que você chama a ferramenta, ela fica vinculada ao projeto:

1. `--project PATH`, se especificado (qualquer pasta dentro do projeto cadastrado é adequada);
2. caso contrário, as pastas de trabalho (raízes) do cliente, se o cliente as reportar;
3. caso contrário, a pasta na qual o cliente iniciou o servidor MCP.

`get_my_board` retorna o projeto e o trabalhador selecionados para `session` junto com `binding_source`. Se as pastas de trabalho do cliente forem alteradas, a próxima chamada à ferramenta executará a ligação novamente. Sem `--worker`, é utilizado o único trabalhador habilitado do projeto.

Caso a vinculação não seja possível, cada ferramenta retorna um erro que o agente deve passar ao usuário; os dados da tarefa não são lidos neste caso:

| Código | Significado |
| --- | --- |
| `NO_PROJECTS` | Não há projetos no banco de dados. Execute `agentboard init` na pasta do projeto e crie um trabalhador. Não é necessário reiniciar. |
| `PROJECT_REQUIRED` | A pasta de trabalho está fora dos projetos cadastrados. A mensagem lista projetos. |
| `PROJECT_NOT_FOUND` | `--project` não está dentro de um projeto cadastrado. |
| `PROJECT_AMBIGUOUS` | As pastas de trabalho de um cliente abrangem vários projetos. Adicione `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Crie, habilite ou selecione um trabalhador por meio de `--worker`. |

As configurações para **Codex** e **Gemini CLI** são personalizadas e compartilhadas em todas as pastas, portanto, não possuem `--project`: Crie um trabalhador com o mesmo slug em cada projeto em que este cliente está sendo executado. Arquivos de nível de projeto (Claude Code, Cursor, OpenCode, VS Code) contêm `--project` de seu projeto.

### Design explícito e verificação de contexto

Se `--project` for especificado e a pasta de inicialização ou pasta de trabalho do cliente estiver em um projeto diferente, as ferramentas solicitarão confirmação por meio da elicitação do MCP. Clientes sem elicitação de formulário recebem `PROJECT_CONTEXT_MISMATCH` antes de acessar as tarefas. Corrija a configuração, remova `--project` para detecção automática ou confirme explicitamente que pretende trabalhar com outro projeto com o sinalizador `--allow-project-mismatch`. A alteração das pastas de trabalho redefine as confirmações anteriores. Projetos definidos automaticamente não requerem confirmação.

**Configurações globais → Contexto do projeto MCP** controla essas verificações para todos os projetos e trabalhadores do mesmo banco de dados. Habilitado por padrão. Se desativado, o MCP usa o projeto especificado sem confirmações ou erros de incompatibilidade. As sessões já abertas aplicam a configuração na próxima vez que a ferramenta for chamada; visibilidade de tarefas, transições e regras de testes humanos continuam a ser seguidas.

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

Exemplos para Codex (configuração do usuário, projeto detectado automaticamente) e Claude Code (arquivo de projeto) com projeto `C:\Projects\MyFirstProject` e trabalhador `codex`:

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

Em JSON, a barra invertida do Windows é duplicada; um fragmento pronto da interface faz isso automaticamente. As configurações geradas contêm o `--db` real do servidor web, portanto HTTP e MCP utilizam a mesma base, mesmo que o ambiente do cliente seja diferente.

**Ollama** fornece um modelo local, mas não substitui o cliente MCP. O perfil **Ollama via OpenCode** gera uma configuração MCP para OpenCode; configure separadamente o OpenCode no modelo Ollama. Outro cliente compatível com MCP com Ollama também é adequado.

## Etapa 4: verifique sua conexão

1. Abra o cartão de trabalhador e clique em **Verificar novamente**. **Server check** executa um executável stdio real com este projeto, trabalhador e base, verifica o handshake, a lista de ferramentas e se o projeto e o trabalhador estão localizados.
2. Inicie ou reinicie o cliente AI após adicionar o arquivo de configuração. Peça a ele para ligar para `get_my_board`.
3. **Cliente conectado** aparecerá no AgentBoard e uma nova sessão aparecerá na lista. Só isso confirma a conexão do seu cliente. Se houver ferramentas, mas o cliente não estiver conectado, verifique o caminho do programa, o nome do arquivo de configuração e sua sintaxe JSON/TOML.

Comece sua sessão de trabalho com `get_my_board`. AI não recebe tarefas de `Backlog` e `Complete`, mesmo que conheça seu ID. A IA não pode fingir ser humana e não tem um comando geral de “mover para qualquer lugar”. Trabalhar com o sistema de arquivos fora do AgentBoard depende dos recursos e da área restrita do cliente selecionado.

Cada ferramenta declara `title` e todas as quatro dicas de comportamento do MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) como valores booleanos explícitos. As ferramentas de leitura são somente leitura e idempotentes; as ferramentas de gravação apenas adicionam registros ou avançam tarefas, portanto, nenhuma perturba ou vai além da base local do AgentBoard.

Instruções oficiais do cliente: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
