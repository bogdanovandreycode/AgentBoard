# Configurações

Abra **Configurações** no menu lateral do projeto selecionado. As alterações são salvas com o botão **Salvar configurações** e armazenadas para o projeto no banco de dados AgentBoard.

## Linguagem e aparência

O campo **Idioma** abre uma lista suspensa de pesquisa. Por padrão, a opção **Seguir sistema** está selecionada, que leva o idioma do navegador. Idiomas disponíveis nas capturas de tela: árabe, português (Brasil), chinês simplificado, tcheco, dinamarquês, holandês, inglês, finlandês, francês, alemão, italiano, japonês, coreano, bokmål norueguês, polonês, russo, espanhol, sueco, turco, ucraniano, vietnamita; adicionalmente bielorrusso, romeno e búlgaro. **Seguir sistema** utiliza o idioma do navegador. As assinaturas principais são traduzidas manualmente, as demais linhas da UI possuem uma tradução automática preliminar. Antes do lançamento público, é aconselhável revisar as traduções feitas por falantes nativos; nomes de clientes, comandos, campos JSON e dados do usuário permanecem sem tradução.

**Esquema de cores**: Escuro (tema original), Claro, Preto, Ubuntu e Windows. **Fuso horário** controla a exibição de datas; os dados continuam a ser armazenados em UTC. **O fuso horário do sistema** usa as configurações do computador.

## Colunas

Os quatro estágios `Features`, `In progress`, `Testing`, `Verification` são fixados nesta ordem: eles não podem ser renomeados ou excluídos. Outras colunas podem ser reorganizadas usando as setas disponíveis. Insira um nome e clique em **Adicionar coluna** para criar uma coluna personalizada. Ele foi projetado para tarefas humanas: a IA vê uma tarefa como `Backlog` e não a recebe por meio do MCP. Excluir uma coluna ao salvar transfere suas tarefas para o `Backlog` normal.

## Web e MCP

**Atualização do tabuleiro** define a taxa de atualização do tabuleiro e dos cartões em segundos (1–60). **Atualização do trabalhador** atualiza o status dos trabalhadores (2 a 120 segundos). Esta é uma pesquisa da interface da web, não uma frequência de disparo de IA. Os trabalhadores não iniciam automaticamente.

O endereço do servidor web é definido ao iniciar a CLI, por exemplo `agentboard open --addr 127.0.0.1:7444`. Para alterar o endereço, o servidor deve ser reiniciado. O padrão é `127.0.0.1:7337`. O MCP opera por meio de um comando local separado `agentboard mcp --project ... --worker ...` e é independente da porta da web. Para uma base fora do padrão, especifique o mesmo `--db` em todos os comandos. Copie a configuração de cada cliente do cartão de trabalho.
