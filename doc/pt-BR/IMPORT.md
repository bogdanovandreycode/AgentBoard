# Importar tarefas do JSON

A guia **Importar tarefas** aceita um arquivo JSON na codificação UTF-8. Use-o se estiver transferindo tarefas de outro serviço. Clique em **Baixar modelo JSON**: o modelo inclui os nomes de propriedades customizadas atuais e `Slug` do trabalhador disponível para este projeto.

## Passo a passo

1. Crie os trabalhadores (**Workers**) e propriedades (**Properties**) necessários antes de importar.
2. Baixe o modelo, abra-o em um editor de texto, substitua os exemplos pelas suas próprias tarefas e salve o arquivo com a extensão `.json`.
3. Na guia **Importar tarefas**, selecione um arquivo. A interface mostrará o número de tarefas.
4. Clique em **Importar tarefas**. O servidor verificará o arquivo inteiro: se for detectado um erro, nenhuma tarefa dele será gravada. Corrija a mensagem de erro e tente novamente.

Arquivo mínimo:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Exemplo com trabalhador, propriedade e dependência:

```json
{
  "version": 1,
  "tasks": [
    {
      "key": "design",
      "title": "Prepare the design",
      "description": "## Goal\nPrepare the home page mockup.",
      "state": "features",
      "priority": "high",
      "testing_mode": "hybrid",
      "assignee": { "type": "worker", "worker": "codex" },
      "ai_test_instructions": "Check the build.",
      "human_test_instructions": "Review the page in a browser.",
      "properties": { "Department": "Design" }
    },
    {
      "title": "Approve the design",
      "depends_on": ["design"],
      "assignee": { "type": "human" }
    }
  ]
}
```

`version` deve ser `1`, matriz `tasks` - de 1 a 1.000 tarefas. `title` é obrigatório. Estágios válidos: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioridades: `critical`, `high`, `medium`, `low`; modos de teste: `ai`, `human`, `hybrid`. Responsável: `unassigned`, `human` ou `worker` com `worker` existente (Slug ou ID). `properties` usa nomes ou IDs de propriedades já criadas. `key` é exclusivo no arquivo; `depends_on` refere-se a essas chaves. O próprio servidor cria IDs de tarefa.

Importar o mesmo arquivo novamente criará novos problemas, portanto verifique o quadro antes de clicar novamente. Quando importadas, as tarefas são registradas como criadas por humanos; Uma entrada do sistema aparece no histórico.
