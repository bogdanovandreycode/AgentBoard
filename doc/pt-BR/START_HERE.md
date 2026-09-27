# Comece aqui

## O que é AgentBoard

Imagine um quadro normal com cartões de tarefas. Você cria tarefas, atribui alguém responsável e supervisiona o trabalho. Os trabalhadores de IA recebem apenas as tarefas atribuídas por meio do MCP e relatam o progresso. Você decide quando a tarefa está finalmente pronta.

**Projeto** - uma pasta no computador e um quadro separado. **Tarefa** - cartão com descrição, responsável e etapa. **Worker** é o nome lógico do cliente de IA. **MCP** é a forma como o cliente AI se conecta ao AgentBoard. **Scoop** é um gerenciador de instalação de software para Windows.

Nenhum código é necessário. Você precisa do Windows, de um navegador, do PowerShell e, para que a IA funcione, de um cliente de IA instalado com suporte para servidores MCP locais.

## Primeiro lançamento

1. Instale o aplicativo de acordo com [instruções](INSTALL.md).
2. Crie uma pasta de projeto no Explorer, por exemplo `C:\Projects\MyFirstProject`.
3. Abra esta pasta no Explorer. Clique na barra de endereço, digite `powershell` e pressione Enter.
4. Na janela que se abre, faça:

```powershell
agentboard init
agentboard open
```

5. Um navegador será aberto com o endereço `http://127.0.0.1:7337`. Deixe a janela do PowerShell aberta enquanto usa o quadro branco. Fechar a janela irá parar o servidor local, mas as tarefas permanecerão.

Se o comando `agentboard` não for encontrado, feche o PowerShell e reabra-o após instalar o Scoop. Ao instalar a partir de um ZIP, use o caminho completo para `agentboard.exe`.

## Primeira tarefa

Clique em **Nova tarefa**, preencha **Título**, se necessário **Descrição** e **Salvar**. Uma nova tarefa em `Backlog` está disponível para humanos. Para que a IA comece a funcionar, atribua um trabalhador e transfira a tarefa para `Features`. Descrição passo a passo dos campos - [TASKS.md](TASKS.md).

## Primeiro trabalhador

Abra **Trabalhadores → Adicionar trabalhador**. Selecione o cliente de IA que você está usando (por exemplo Codex ou Claude Code), verifique o nome e o ID abreviado `Slug`, clique em **Salvar**. Abra o trabalhador criado: há uma configuração do MCP, um botão de cópia e uma verificação do servidor. Copie a configuração para o cliente AI de acordo com [WORKERS_MCP.md](WORKERS_MCP.md). Uma vez conectado, peça ao cliente para ligar para `get_my_board`.

## O que ler a seguir

- [Tarefas, testes, histórias e colunas](TASKS.md)
- [Conectando trabalhadores e verificando MCP](WORKERS_MCP.md)
- [Importar tarefas de JSON](IMPORT.md)
- [Idioma, tema, fuso horário e atualização](SETTINGS.md)
- [Problemas típicos](TROUBLESHOOTING.md)
