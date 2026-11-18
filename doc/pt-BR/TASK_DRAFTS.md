# Gerenciamento de projetos e rascunhos de tarefas

A página **Todos os projetos** mostra a pasta cadastrada, o número de projetos publicados
tarefas e atividades mais recentes. A atividade leva em consideração mudanças nas tarefas, histórico e
registros de despesas. O botão de lápis permite renomear o projeto ou selecionar
pasta existente na caixa de diálogo do sistema. A alteração do caminho atualiza o cadastro e o marcador;
IDs de projetos, trabalhadores, tarefas e configurações são salvos, os arquivos não são movidos.
Pastas que já estão cadastradas em outro projeto não podem ser selecionadas.

## Configurações globais

Abra **Configurações globais** na página ou guia de seleção de projetos
**Configurações globais** nas configurações de qualquer placa. Esses parâmetros são comuns a todos os projetos
uma base AgentBoard:

- O salvamento automático de rascunhos por cronômetro está ativado por padrão, a cada **10 segundos**.
- **Salvar ao sair do campo** está habilitado por padrão e pode ser configurado separadamente.
- As verificações de **contexto do projeto MCP** são ativadas e desativadas globalmente. A área deles
  ações e um sinalizador CLI separado são descritos na seção [Workers e MCP](WORKERS_MCP.md).

## Rascunho do ciclo de vida

A caixa de diálogo Nova tarefa salva o trabalho inacabado como rascunho. Linha de status
distingue entre alterações não salvas, salvamento em andamento, rascunho salvo e erro
conservação. A minuta não aparece no quadro e não está disponível para os trabalhadores do MCP até
publicado.

Antes de publicar, use o botão **Criar tarefa** para preencher **Título, Estado,
Responsável e Descrição**. Antes disso, **Salvar rascunho** salva um rascunho sem criar
tarefas. Um formulário vazio é descartado por padrão quando fechado. Se desmarcado
rascunho existente para um estado vazio, ele será excluído na próxima vez que for salvo.
Fechar uma caixa de diálogo não vazia também a salva.

O botão **Rascunhos** no quadro permite continuar ou excluir rascunhos salvos.
Vários rascunhos suportados. Após uma reinicialização, o aberto
rascunho. Os rascunhos são armazenados no SQLite; A cada mudança, o local é atualizado
copie no navegador, para que as edições não salvas sobrevivam às reinicializações ou
falha temporária do servidor. Se o armazenamento do navegador não estiver disponível, uma caixa de diálogo avisará
que o rascunho deve ser salvo no servidor antes de fechar.

Escrever rascunhos usa revisões para evitar que uma janela substitua o trabalho
outro. A resposta de salvamento perdido é consistente se já houver uma no servidor
dados idênticos. Se o conteúdo for diferente, um conflito será mostrado e
**Recuperar como um novo rascunho** salva o conteúdo local separadamente. Publicação
atômico e idempotente: republicar o mesmo rascunho depois de perdido
resposta retorna a mesma tarefa em vez de criar uma cópia.
