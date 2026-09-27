# JSON からタスクをインポート

[**タスクのインポート**] タブでは、UTF-8 エンコードの JSON ファイルを受け入れます。別のサービスからタスクを転送する場合に使用します。 [**JSON テンプレートをダウンロード**] をクリックします。テンプレートには、現在のカスタム プロパティ名と、このプロジェクトで使用可能なワーカーの `Slug` が含まれています。

## ステップバイステップ

1. インポートする前に、必要なワーカー (**Workers**) とプロパティ (**Properties**) を作成します。
2. テンプレートをダウンロードしてテキスト エディターで開き、例を独自のタスクに置き換えて、`.json` 拡張子を付けてファイルを保存します。
3. **タスクのインポート** タブで、ファイルを選択します。インターフェイスにはタスクの数が表示されます。
4. [**タスクのインポート**] をクリックします。サーバーはファイル全体をチェックします。エラーが検出された場合、そのファイルからのタスクは 1 つも書き込まれません。エラー メッセージを修正して、再試行してください。

最小ファイル:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

ワーカー、プロパティ、依存関係の例:

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

`version` は `1`、配列 `tasks` - 1 ～ 1000 タスクである必要があります。 `title`は必須です。有効なステージ: `backlog`、`features`、`in_progress`、`testing`、`verification`、`complete`。優先順位: `critical`、`high`、`medium`、`low`;テストモード: `ai`、`human`、`hybrid`。責任者: `unassigned`、`human`、または `worker` と既存の `worker` (スラッグまたは ID)。 `properties` は、すでに作成されているプロパティの名前または ID を使用します。 `key` はファイル内で一意です。 `depends_on` はそのようなキーを指します。サーバー自体がタスク ID を作成します。

同じファイルを再度インポートすると新しい問題が発生するため、再度クリックする前にボードを確認してください。インポートされると、タスクは人間が作成したものとして記録されます。システム エントリが履歴に表示されます。
