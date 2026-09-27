# 从 JSON 导入任务

**导入任务**选项卡接受 UTF-8 编码的 JSON 文件。如果您要从其他服务转移任务，请使用它。单击 **下载 JSON 模板**：该模板包含该项目可用工作人员的当前自定义属性名称和 `Slug`。

## 一步一步

1. 在导入之前创建必要的工作人员（**Workers**）和属性（**Properties**）。
2. 下载模板，在文本编辑器中打开它，将示例替换为您自己的任务，并使用 `.json` 扩展名保存文件。
3. 在 **导入任务** 选项卡上，选择一个文件。界面会显示任务数量。
4. 单击“**导入任务**”。服务器将检查整个文件：如果检测到错误，则不会写入其中的任何任务。更正错误消息并重试。

最小文件：

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

工人、财产和依赖性的示例：

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

`version` 应为 `1`，数组 `tasks` - 从 1 到 1000 个任务。需要 `title`。有效阶段：`backlog`、`features`、`in_progress`、`testing`、`verification`、`complete`。优先级：`critical`、`high`、`medium`、`low`；测试模式：`ai`、`human`、`hybrid`。负责：`unassigned`、`human` 或 `worker` 以及现有的 `worker`（Slug 或 ID）。 `properties` 使用已创建属性的名称或 ID。 `key` 在文件内是唯一的； `depends_on` 指的是此类键。服务器自行创建任务 ID。

再次导入相同的文件会产生新的问题，因此请在再次单击之前检查面板。导入时，任务被记录为人工创建的；系统条目出现在历史记录中。
