# Resolução de problemas

| Sintoma | O que verificar |
| --- | --- |
| `agentboard` não encontrado | Reinicie o PowerShell após o Scoop. Ao instalar a partir de um ZIP, use o caminho completo para `agentboard.exe`. |
| Página da Web mostrando interface antiga após compilação | Pare o servidor em execução Ctrl+C. Execute `./scripts/build.ps1`, inicie um novo binário. O Vite deve ser compilado antes do Go porque a interface está incorporada no EXE. Atualize a página Ctrl+F5. |
| Porta 7337 ocupada | AgentBoard pode já estar em execução. Abra `http://127.0.0.1:7337` ou encerre o processo antigo. Para outra porta use `--addr`. |
| Projeto não encontrado | Na pasta desejada, execute `agentboard init`. Então `agentboard open` dele ou `agentboard open C:\путь\к\проекту`. |
| O trabalhador não vê a tarefa | A tarefa deve ser atribuída a este trabalhador específico e estar localizada em `Features`, `In progress`, `Testing` ou `Verification`. AI não vê `Backlog`, `Complete` e colunas personalizadas. |
| `NO_PROJECTS` | Não há projetos na base de dados do MCP. Execute `agentboard init` na pasta do projeto e crie um trabalhador; o agente pode repetir a chamada sem reiniciar. Verifique se o MCP e a interface da web estão usando o mesmo `--db`. |
| `PROJECT_REQUIRED` / `PROJECT_AMBIGUOUS` | Sem `--project`, nem a pasta de trabalho do cliente nem a pasta de inicialização do MCP estão exatamente no mesmo projeto registrado. Execute o agente na pasta do projeto, registre-o via `agentboard init` ou adicione `--project`. Não há necessidade de iniciar a interface da web. |
| `PROJECT_NOT_FOUND`/`WORKER_NOT_FOUND`/`WORKER_REQUIRED`/`NO_WORKERS` | Confira `--project`, a lesma de trabalho neste projeto e `--db`. A mensagem lista projetos registrados ou trabalhadores habilitados. Uma configuração MCP customizada sem `--project` precisa de um trabalhador com o mesmo slug em cada projeto. |
| `PROJECT_CONTEXT_MISMATCH` | A pasta de inicialização do MCP ou a pasta de trabalho do cliente não corresponde ao `--project` especificado. Por favor, confirme a pergunta se for intencional, caso contrário, remova `--project` (detecção automática) ou indique o correto. |
| A verificação do servidor MCP existe, mas o cliente não está conectado | A verificação do servidor não verifica as configurações do cliente externo. Reinicie o cliente, verifique seu arquivo de configuração, caminho para `agentboard.exe`, `--project`, `--worker` e `--db` geral. Peça para ligar para `get_my_board`. |
| Trabalhador off-line | O cliente pode ter encerrado ou ainda não ter iniciado o MCP. Após 90 segundos sem pulsação, a sessão é considerada desconectada. |
| Arquivo JSON não importando | Verifique `version: 1`, `title` necessário, trabalhadores `Slug` existentes e nomes de propriedades. JSON não permite comentários ou vírgulas finais. |
| Não é possível reconstruir `agentboard.exe` | O Windows não pode substituir um EXE em execução. Pare o servidor Ctrl+C e tente a compilação novamente. |
| As tarefas desapareceram após a atualização | Verifique se `--db` não aponta para outro arquivo e se você está conectado como o mesmo usuário do Windows. A base padrão está localizada em `%AppData%\AgentBoard`. |

Se o erro não for descrito, colete o texto exato da mensagem, as versões `agentboard version` e Windows e as etapas para tentar novamente. Não publique dados de projetos privados ou conteúdos de bancos de dados em uma edição aberta.
