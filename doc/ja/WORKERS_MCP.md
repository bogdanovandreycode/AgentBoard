# ワーカーと MCP 接続

## ステップ 1. ワーカーを作成する

AgentBoard で、**「ワーカー」→「ワーカーの追加」** を開きます。クライアント プロファイルを選択します: Codex、Claude Code、Gemini CLI、Cursor、OpenCode、Ollama via OpenCode、VS Code Copilot、または別の MCP クライアント。プロファイルには、典型的な機能 (コード、テスト、Git) が入力されます。変更することができます。この名前はボード上に表示され、`Slug` は MCP コマンドのスペースを含まない短い識別子です。 [**保存**] をクリックします。

1 人のワーカーが 1 つの AI 人格に対応します。異なるクライアントまたはチームに異なるワーカーを作成します。機能プロファイルは専門化を説明しますが、AI の権限をタスク段階まで拡張するものではありません。

## ステップ 2. 構成をコピーする

作成したワーカーを開きます。 **MCP 診断** ブロックには、構成フラグメントとそれを追加するファイルが表示されます。 [**MCP 構成をコピー**] をクリックします。ファイルが既に存在する場合は、他のサーバーを削除せずに、提案されたサーバーを既存の `mcpServers`/`servers`/`mcp` オブジェクトに追加します。

メインコマンドは次のようになります。

```text
agentboard mcp --worker codex
```

`--worker` — 作成されたワーカーの `Slug`。 MCP クライアントは、ツールが必要なときにこのコマンド自体を実行します。

### MCP は Web インターフェイスなしで動作します

`agentboard mcp` は、それ自体が AgentBoard データベースを開く別個の stdio プロセスです。 `agentboard open`/`serve` は必要ありません。Web インターフェイスが起動されるフォルダーや、実行されている Web インターフェイスのインスタンスの数には依存しません。複数のエージェントをそれぞれ独自のプロジェクトで同時に実行できます。各クライアントは独自の MCP プロセスを実行し、すべて同じ SQLite データベースを安全に使用します。

### プロジェクトの選択方法

MCP サーバーは任意のフォルダーで実行されます。初めてツールを呼び出すと、ツールはプロジェクトにバインドされます。

1. `--project PATH` (指定されている場合) (登録されたプロジェクト内の任意のフォルダーが適しています)。
2. それ以外の場合は、クライアントの作業フォルダー (ルート) (クライアントが報告した場合)。
3. それ以外の場合は、クライアントが MCP サーバーを起動したフォルダー。

`get_my_board` は、選択したプロジェクトとワーカーを `session` とともに `binding_source` に返します。クライアントの作業フォルダーが変更された場合、ツールへの次の呼び出しでバインドが再度実行されます。 `--worker` がない場合、プロジェクトの唯一の有効なワーカーが使用されます。

バインドが不可能な場合、各ツールはエラーを返し、エージェントはそれをユーザーに渡す必要があります。この場合、タスク データは読み取られません。

|コード |意味 |
| --- | --- |
| `NO_PROJECTS` |データベースにはプロジェクトがありません。プロジェクトフォルダー内の`agentboard init`を実行し、ワーカーを作成します。再起動は必要ありません。 |
| `PROJECT_REQUIRED` |作業フォルダーは登録されたプロジェクトの外にあります。メッセージにはプロジェクトがリストされます。 |
| `PROJECT_NOT_FOUND` | `--project` は登録されたプロジェクト内にありません。 |
| `PROJECT_AMBIGUOUS` |クライアントの作業フォルダーは複数のプロジェクトにまたがります。 `--project`を追加します。 |
| `NO_WORKERS`、`WORKER_REQUIRED`、`WORKER_NOT_FOUND`、`WORKER_DISABLED` | `--worker` を介してワーカーを作成、有効化、または選択します。 |

**Codex** と **Gemini CLI** の構成はカスタムであり、すべてのフォルダー間で共有されるため、`--project` がありません。このクライアントが実行されているすべてのプロジェクトで同じスラッグを持つワーカーを作成します。プロジェクト レベル ファイル (Claude Code、Cursor、OpenCode、VS Code) には、プロジェクトの `--project` が含まれています。

### 明示的な設計とコンテキストのチェック

`--project` が指定されており、クライアントの起動フォルダーまたは作業フォルダーが別のプロジェクトにある場合、ツールは MCP 誘発を通じて確認を求めます。フォーム引き出しのないクライアントは、タスクにアクセスする前に `PROJECT_CONTEXT_MISMATCH` を受け取ります。構成を修正するか、自動検出のために `--project` を削除するか、`--allow-project-mismatch` フラグを使用して別のプロジェクトを操作するつもりであることを明示的に確認してください。作業フォルダーを変更すると、以前の確認がリセットされます。自動的に定義されたプロジェクトには確認は必要ありません。

**グローバル設定 → MCP プロジェクト コンテキスト** は、同じデータベースのすべてのプロジェクトとワーカーに対するこれらのチェックを制御します。デフォルトで有効になっています。無効にすると、MCP は確認や不一致エラーなしで指定されたプロジェクトを使用します。すでに開いているセッションには、次回ツールが呼び出されたときに設定が適用されます。タスクの可視性、移行、人間によるテストのルールは引き続き遵守されます。

確認後、診断ブロックには、実行中の `agentboard.exe` への絶対パスが表示されます。 ZIP からインストールする場合に特に便利です。 Scoop 経由でインストールする場合、クライアントが同じ `agentboard` を認識する場合は、`PATH` コマンドを使用できます。

## ステップ 3: サーバーをクライアントに追加する

ワーカー インターフェイスにはすでに既製のフラグメントがあります。以下に使用場所について説明します。

|クライアント |挿入する場所 |クライアント側で確認する方法 |
| --- | --- | --- |
|コーデックス | `%USERPROFILE%\.codex\config.toml`、セクション `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
|クロード・コード |プロジェクト フォルダー内の `.mcp.json` | `claude mcp list` |
|ジェミニ CLI | `%USERPROFILE%\.gemini\settings.json`、オブジェクト `mcpServers` | Gemini CLI の `/mcp list` |
|カーソル | `.cursor\mcp.json` プロジェクト |カーソル設定の MCP サーバーのリスト |
|オープンコード | `opencode.json` プロジェクト、オブジェクト `mcp` | OpenCode の MCP ツールのリスト |
| VS コードのコパイロット | `.vscode\mcp.json` プロジェクト、オブジェクト `servers` |コマンド **MCP: サーバーのリスト** |

プロジェクト `C:\Projects\MyFirstProject` およびワーカー `codex` を使用した Codex (ユーザー構成、プロジェクトが自動的に検出される) および Claude Code (プロジェクト ファイル) の例:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

JSON では、Windows のバックスラッシュが 2 倍になります。インターフェイスからの既製のフラグメントがこれを自動的に行います。生成された構成には Web サーバーの実際の `--db` が含まれるため、クライアント環境が異なる場合でも、HTTP と MCP は同じベースを使用します。

**Ollama** はローカル モデルを提供しますが、MCP クライアントを置き換えるものではありません。 **Ollama via OpenCode** プロファイルは、OpenCode の MCP 構成を生成します。 Ollama モデルで OpenCode を別途設定します。 Ollama を備えた別の MCP 互換クライアントも適しています。

## ステップ 4: 接続を確認する

1. 従業員カードを開き、**もう一度確認** をクリックします。 **サーバー チェック** は、このプロジェクト、ワーカー、ベースで実際の stdio 実行可能ファイルを実行し、ハンドシェイク、ツールのリスト、およびプロジェクトとワーカーが配置されていることをチェックします。
2. 構成ファイルを追加した後、AI クライアントを起動または再起動します。 `get_my_board` に電話してもらいます。
3. **クライアントが接続されました** が AgentBoard に表示され、新しいセッションがリストに表示されます。これだけでクライアントの接続が確認されます。ツールはあるがクライアントが接続されていない場合は、プログラムへのパス、構成ファイルの名前、およびその JSON/TOML 構文を確認してください。

`get_my_board` で作業セッションを開始します。 AI は、`Backlog` と `Complete` の ID がわかっていても、タスクを受け取りません。 AI は人間のふりをすることはできず、一般的な「どこにでも移動」コマンドはありません。 AgentBoard の外部でファイル システムを操作するかどうかは、選択したクライアントの機能とサンドボックスによって異なります。

各ツールは、`title` と 4 つの MCP 動作ヒント (`readOnlyHint`、`destructiveHint`、`idempotentHint`、`openWorldHint`) を明示的なブール値として宣言します。読み取りツールは読み取り専用で冪等です。記録ツールはレコードを追加するか、タスクを前進させるだけなので、中断を引き起こしたり、ローカルの AgentBoard ベースを超えたりするものはありません。

公式の顧客向け指示: [Codex](https://developers.openai.com/learn/docs-mcp)、[クロード Code](https://docs.anthropic.com/en/docs/claude-code/mcp)、[Gemini CLI](https://geminicli.com/docs/tools/mcp-server/)、[Cursor](https://prod.cursor.com/docs/cli/mcp)、[OpenCode](https://opencode.ai/docs/mcp-servers/)、[VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers)]
