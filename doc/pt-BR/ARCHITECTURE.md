# Arquitetura e API

AgentBoard é um processo Go local executando SQLite. A API HTTP para humanos, o adaptador MCP para IA e a interface da web compartilham uma lógica de serviço comum. SQLite – fonte de estado; frontend não ignora o servidor. As ferramentas de IA têm uma superfície de direitos separada e o status `Backlog`/`Complete` no MCP não está disponível. O registro histórico do `System` é criado pelo próprio aplicativo.

## Componentes

-`cmd/agentboard` -CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` – modelos de domínio e erros.
- `internal/service` – transições de tarefas, permissões, importação e configurações.
- `internal/persistence` - SQLite e migrações.
- `internal/httpapi` - API HTTP humana.
- `internal/mcpserver` - Ferramentas MCP para um trabalhador separado.
- `web` - UI React/TypeScript; a montagem termina em `internal/webui/dist` e é incluída no EXE.

## Rotas HTTP básicas

| Método e caminho | Destino |
| --- | --- |
| `GET /api/health` | Verificação do servidor. |
| `GET /api/projects` | Projetos cadastrados. |
| `GET /api/projects/{id}/board` | Quadro do projeto. |
| `GET/PUT /api/projects/{id}/settings` | Configurações, ordem e nomes das colunas. |
| `POST /api/projects/{id}/tasks` | Crie uma tarefa. |
| `POST /api/projects/{id}/tasks/import` | Importar atomicamente JSON versão 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Cartão, alterar, excluir. |
| `POST /api/tasks/{id}/move` | Movendo-se por uma pessoa. |
| `GET/POST /api/projects/{id}/workers` | Lista e criação de um trabalhador. |
| `GET /api/workers/{id}/mcp/check` | Testes internos de handshake e ferramentas do MCP. |
| `GET /api/workers/{id}/sessions` | Diagnóstico de sessão. |
| `GET/POST /api/projects/{id}/properties` | Propriedades personalizadas. |

HTTP é para o usuário confiável local. Não publique uma porta web na Internet sem sua própria autenticação, restrições de rede e HTTPS. O servidor MCP é iniciado usando `stdio` para um projeto e trabalhador específico; comece com `get_my_board`. Suas permissões restringem ações dentro do AgentBoard, mas não substituem a sandbox do sistema de arquivos do cliente AI.

As colunas do usuário são armazenadas separadamente de `tasks.state`: o núcleo mantém a tarefa em tal coluna em `backlog` e `tasks.board_column` determina o local no quadro Humano. Isso mantém o mesmo modelo de transição de IA. A remoção de uma coluna limpa as tarefas `board_column`.
