# Instalação e lançamento no Windows

## Instalador (recomendado)

Baixe `agentboard-VERSION-windows-amd64-setup.exe` da página [Lançamentos](https://github.com/bogdanovandreycode/AgentBoard/releases). O assistente de instalação irá sugerir uma pasta; o padrão é `C:\AI\AgentBoard`. Ele copiará `agentboard.exe` e documentação, criará um atalho e adicionará a pasta selecionada ao sistema `PATH`. Após a instalação, abra um novo terminal para que o comando `agentboard` fique disponível.

Na pasta do seu projeto execute:

```powershell
agentboard init
agentboard open
```

A interface está integrada ao `agentboard.exe`; Nenhuma instalação separada do Go ou Node.js é necessária. Os dados são armazenados em `%AppData%\AgentBoard` e retidos quando o programa é atualizado ou desinstalado. A desinstalação por meio de “Aplicativos instalados” remove os atalhos e a entrada de `PATH`.

## Instalando o Scoop

Se o Scoop ainda não estiver instalado, abra o PowerShell como seu usuário regular e siga as [instruções oficiais do Scoop](https://scoop.sh/). Se houver restrições no seu computador corporativo, entre em contato com o administrador; AgentBoard também pode ser iniciado a partir de um ZIP sem Scoop.

Como alternativa, instale o AgentBoard via Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

A disponibilidade do manifesto no GitHub Release pode ser verificada na [página Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Depois de adicionar o manifesto ao Scoop, o bucket pode ser instalado pelo nome do bucket e atualizado com o comando `scoop update agentboard`.

## ZIP sem furo

Baixe `agentboard-VERSION-windows-amd64.zip` em Releases, descompacte, por exemplo, em `C:\Tools\AgentBoard`. No PowerShell, na pasta do projeto:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Para um cliente AI, especifique o caminho completo para `agentboard.exe` em sua configuração MCP se o programa não estiver localizado em `PATH`.

## Construa a partir da fonte

Instale a versão Go de `go.mod` e Node.js 22 ou posterior. No PowerShell na raiz do repositório:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` monta o frontend em `internal/webui/dist`, executa testes Go e monta um `agentboard.exe`. A ordem é importante: a interface é incorporada ao binário quando o Go é compilado. Se o `agentboard.exe open` antigo estiver em execução, pare-o antes de reconstruí-lo (Ctrl+C), caso contrário, o Windows não permitirá que você substitua o arquivo.

## Onde estão os dados

- `%AppData%\AgentBoard\agentboard.db` – tarefas, projetos, trabalhadores e configurações. Você pode especificar um arquivo diferente com o sinalizador `--db`, mas `init`, `open`/`serve` e `mcp` devem ter o **mesmo caminho**.
- `<ваш проект>\.agentboard\project.json` — identificador do projeto. Este arquivo não contém tarefas.
- O servidor escuta apenas `127.0.0.1:7337` por padrão. Especifique outro endereço `--addr` antes do caminho do projeto: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` abre a página do projeto e reutiliza um servidor já em execução nesse endereço. Se vários projetos estiverem abertos em um navegador, selecione-os na lista à esquerda.
