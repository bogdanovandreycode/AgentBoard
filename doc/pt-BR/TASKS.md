# Trabalhando com tarefas

## Adicionar tarefa

Na guia **Quadro**, clique em **Nova tarefa**. Preencha o título e a descrição. A descrição e as instruções de teste suportam Markdown: destaque o texto e use a barra de formatação para negrito, itálico, link e lista. Clique em **Salvar**.

Campos:

| Campo | O que faz |
| --- | --- |
| Título | Nome abreviado da tarefa. |
| Descrição | O que precisa ser feito e como entender que a obra está pronta. |
| Estado | Estágio de trabalho; uma nova tarefa geralmente começa em `Backlog`. |
| Prioridade | `critical`, `high`, `medium` ou `low`. |
| Responsável | Uma pessoa, um trabalhador específico ou sem marcação. |
| Modo de teste | `AI` - verifica IA; `Human` – verificações humanas; `Hybrid` - ambos. |
| Dependências | Tarefas que devem ser concluídas mais cedo. |
| Instruções de teste de IA/humano | Instruções para o revisor apropriado. |
| Propriedades personalizadas | Campos adicionais criados na guia **Propriedades**. |

Para uma tarefa de IA, primeiro crie um trabalhador, selecione-o em Responsável e depois transfira o cartão para `Features`. A IA só vê as tarefas atribuídas a ela em quatro estágios, de `Features` a `Verification`.

## Mover tarefa

Arraste o cartão entre as colunas. Você pode alternar entre uma visualização de todas as colunas e colunas largas com rolagem horizontal. `Backlog` e `Complete` são controlados por humanos. A IA só pode avançar na tarefa `Features → In progress → Testing → Verification` por meio de ferramentas MCP especiais. Uma tarefa com teste manual incompleto não deve passar na verificação de IA.

## Cartão de tarefas

Clique em um cartão para ver a descrição, o proprietário, as instruções do teste e as guias de histórico, testes, artefatos e custos de IA. **Editar tarefa** altera o conteúdo. O comentário da pessoa é adicionado à história geral. A busca no topo filtra as buscas por título, descrição, ID e trabalhador; um botão separado abre uma grande pesquisa.

## Colunas e propriedades

Na aba **Configurações** você pode alterar a ordem das colunas disponíveis para uma pessoa e adicionar a sua própria. Os quatro estágios da IA ​​são fixos e prosseguem na mesma ordem. A coluna do usuário é um local para tarefas adiadas por uma pessoa: para IA, tal tarefa tem o status `Backlog`. Quando uma coluna é excluída, suas tarefas voltam ao normal `Backlog`.

Na aba **Propriedades** você pode adicionar campos como texto, número, bandeira, data, seleção e URL. `Human only` esconde o campo da IA; `Agent read` permite leitura, `Agent read/write` também permite gravação por meio de ferramentas suportadas. Isto não altera os direitos da IA ​​nas etapas da tarefa.
