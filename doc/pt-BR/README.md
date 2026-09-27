# AgentBoard

[🌐 Languages](../LANGUAGES.md)

AgentBoard é um quadro de tarefas local onde trabalhadores humanos e de IA trabalham em tarefas comuns, mas têm direitos diferentes. O aplicativo é iniciado com um arquivo `agentboard.exe`, abre a interface web no navegador e fornece aos trabalhadores um servidor MCP separado via `stdio`. Os dados permanecem no seu computador.

**[Começar do zero](START_HERE.md) · [Trabalhar com tarefas](TASKS.md) · [Conectar IA via MCP](WORKERS_MCP.md) · [Importar JSON](IMPORT.md) · [Configurações](SETTINGS.md) · [Resolver problemas](TROUBLESHOOTING.md)**

## Em cinco minutos

1. Baixe o instalador `agentboard-VERSION-windows-amd64-setup.exe` em [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Ele irá sugerir uma pasta (por padrão `C:\AI\AgentBoard`) e adicioná-la ao `PATH`. Scoop e ZIP também estão disponíveis.
2. Abra o PowerShell na pasta do seu projeto, por exemplo `C:\Projects\MyApp`.
3. Execute `agentboard init` (para ZIP: caminho completo para `agentboard.exe` e `init`).
4. Execute `agentboard open`. `http://127.0.0.1:7337` será aberto.
5. Adicione uma tarefa usando o botão **Nova tarefa**. Para um trabalhador de IA, abra **Trabalhadores → Adicionar trabalhador**, selecione o perfil do cliente e copie a configuração do MCP.

Se você ainda não possui uma pasta de projeto, crie uma no Windows Explorer. Um projeto pode ser qualquer pasta, mesmo sem Git e código.

## Instalação via Scoop

No PowerShell com [Scoop](https://scoop.sh/)] já instalado após o lançamento:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Para desenvolvedores existe [construir a partir da fonte ](INSTALL.md). O fluxo de trabalho de lançamento cria um instalador, manifesto ZIP e Scoop com SHA-256 a partir do mesmo artefato. Atualizando a versão instalada via Scoop: `scoop update agentboard` após adicionar manifesto ao bucket; detalhes - [preparação para lançamento](SCOOP_RELEASE.md).

## Como o conselho está estruturado

`Backlog → Features → In progress → Testing → Verification → Complete`

Uma pessoa pode mover tarefas pelo quadro. AI só pode mover `Features → In progress → Testing → Verification`; a aceitação final no `Complete` é realizada por um ser humano. As colunas do usuário são destinadas a humanos: a tarefa nelas permanece no estado `Backlog` para MCP. Modos de teste: IA, Humano e Híbrido. Histórico, testes, artefatos e custos de IA estão associados à tarefa.

## Equipes

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registra a pasta e grava apenas `.agentboard/project.json` nela. Os dados de produção do SQLite estão localizados no diretório de configuração do usuário do Windows (`%AppData%\AgentBoard\agentboard.db`), fora do projeto e fora da instalação do Scoop. A remoção ou atualização do pacote não deve remover esses dados. Antes de transferir para outro computador, faça uma cópia do banco de dados enquanto o AgentBoard estiver parado.

Trabalhadores - contas lógicas; O próprio AgentBoard não executa Codex, Claude ou qualquer outro cliente de IA. O cliente inicia um processo MCP local para um trabalhador específico. MCP é o limite de permissão do aplicativo e, para isolamento de arquivos, use a sandbox do cliente AI.

## Para desenvolvedores

Stack: Go, SQLite, SDK oficial do MCP Go, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Crie o frontend primeiro e depois vá: `./scripts/build.ps1`. Os arquivos da Web são incluídos no binário por meio de `go:embed`. A arquitetura e a API estão descritas em [doc/ARCHITECTURE.md](ARCHITECTURE.md).
