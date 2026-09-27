# アーキテクチャと API

AgentBoard は、SQLite を実行するローカル Go プロセスです。人間用の HTTP API、AI 用の MCP アダプター、および Web インターフェイスは、共通のサービス ロジックを共有します。 SQLite - 状態ソース。フロントエンドはサーバーをバイパスしません。 AI ツールには別の権利サーフェスがあり、MCP の `Backlog`/`Complete` ステータスは利用できません。 `System` からの履歴レコードは、アプリケーション自体によって作成されます。

## コンポーネント

- `cmd/agentboard` - CLI `init`、`open`、`serve`、`mcp`、`version`。
- `internal/core` - ドメイン モデルとエラー。
- `internal/service` - タスクの遷移、権限、インポートおよび設定。
- `internal/persistence` - SQLite と移行。
- `internal/httpapi` - 人間の HTTP API。
- `internal/mcpserver` - 個別のワーカー用の MCP ツール。
- `web` - React/TypeScript UI;アセンブリは最終的に `internal/webui/dist` になり、EXE に含まれます。

## 基本的な HTTP ルート

|メソッドとパス |目的地 |
| --- | --- |
| `GET /api/health` |サーバーチェック。 |
| `GET /api/projects` |登録されたプロジェクト。 |
| `GET /api/projects/{id}/board` |プロジェクトボード。 |
| `GET/PUT /api/projects/{id}/settings` |列の設定、順序、名前。 |
| `POST /api/projects/{id}/tasks` |タスクを作成します。 |
| `POST /api/projects/{id}/tasks/import` | JSON バージョン 1 をアトミックにインポートします。
| `GET/PATCH/DELETE /api/tasks/{id}` |カード、変更、削除。 |
| `POST /api/tasks/{id}/move` |人による移動。 |
| `GET/POST /api/projects/{id}/workers` |ワーカーのリストと作成。 |
| `GET /api/workers/{id}/mcp/check` | MCP ハンドシェイクとツールの内部テスト。 |
| `GET /api/workers/{id}/sessions` |セッション診断。 |
| `GET/POST /api/projects/{id}/properties` |カスタムプロパティ。 |

HTTP はローカルの信頼できるユーザー用です。独自の認証、ネットワーク制限、および HTTPS を使用せずに、Web ポートをインターネットに公開しないでください。 MCP サーバーは、特定のプロジェクトとワーカーに対して `stdio` を使用して起動されます。 `get_my_board` から始まります。その権限は AgentBoard 内のアクションを制限しますが、AI クライアントのファイル システム サンドボックスを置き換えるものではありません。

ユーザー列は `tasks.state` とは別に保存されます。コアは `backlog` のそのような列にタスクを保持し、`tasks.board_column` がヒューマン ボード上の場所を決定します。これにより、同じ AI 移行モデルが維持されます。列を削除すると、`board_column` タスクがクリアされます。
