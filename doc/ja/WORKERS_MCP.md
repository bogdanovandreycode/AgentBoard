# ワーカーと MCP 接続

## ステップ 1. ワーカーを作成する

AgentBoard で、**「ワーカー」→「ワーカーの追加」** を開きます。クライアント プロファイルを選択します: Codex、Claude Code、Gemini CLI、Cursor、OpenCode、Ollama via OpenCode、VS Code Copilot、または別の MCP クライアント。プロファイルには、典型的な機能 (コード、テスト、Git) が入力されます。変更することができます。この名前はボード上に表示され、`Slug` は MCP コマンドのスペースを含まない短い識別子です。 [**保存**] をクリックします。

1 人のワーカーが 1 つの AI 人格に対応します。異なるクライアントまたはチームに異なるワーカーを作成します。機能プロファイルは専門化を説明しますが、AI の権限をタスク段階まで拡張するものではありません。

## ステップ 2. 構成をコピーする

作成したワーカーを開きます。 **MCP 診断** ブロックには、構成フラグメントとそれを追加するファイルが表示されます。 [**MCP 構成をコピー**] をクリックします。ファイルが既に存在する場合は、他のサーバーを削除せずに、提案されたサーバーを既存の `mcpServers`/`servers`/`mcp` オブジェクトに追加します。

メインコマンドは次のようになります。

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` は、`agentboard init` で登録したフォルダーを指す必要があります。 `--worker` — 作成されたワーカーの `Slug`。 MCP クライアントは、ツールが必要なときにこのコマンド自体を実行します。 AgentBoard ブラウザでは、Web サーバーを個別に実行できます。

確認後、診断ブロックには、実行中の `agentboard.exe` への絶対パスが表示されます。 ZIP からインストールする場合に特に便利です。 Scoop 経由でインストールする場合、クライアントが同じ `agentboard` を認識する場合は、`PATH` コマンドを使用できます。

## ステップ 3: サーバーをクライアントに追加する

ワーカー インターフェイスにはすでに既製のフラグメントがあります。以下に使用場所について説明します。

|クライアント |挿入する場所 |クライアント側で確認する方法 |
| --- | --- | --- |
|コーデックス | `%USERPROFILE%\.codex\config.toml`、セクション `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
|クロード・コード |プロジェクトフォルダー内の `.mcp.json` | `claude mcp list` |
|ジェミニ CLI | `%USERPROFILE%\.gemini\settings.json`、オブジェクト `mcpServers` | Gemini CLI の `/mcp list` |
|カーソル | `.cursor\mcp.json` プロジェクト |カーソル設定の MCP サーバーのリスト |
|オープンコード | `opencode.json` プロジェクト、オブジェクト `mcp` | OpenCode の MCP ツールのリスト |
| VS コードのコパイロット | `.vscode\mcp.json` プロジェクト、オブジェクト `servers` |コマンド **MCP: サーバーのリスト** |

Codex と Claude Code の場合、プロジェクト `C:\Projects\MyFirstProject` とワーカー `codex` の例:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

JSON では、Windows のバックスラッシュは 2 つになります。インターフェイスからの既製のフラグメントがこれを自動的に行います。カスタム ベースで `--db` を使用している場合は、それを `args` MCP 構成に追加し、`open` の起動時と同じパスを指定します。

**Ollama** はローカル モデルを提供しますが、MCP クライアントを置き換えるものではありません。 **Ollama via OpenCode** プロファイルは、OpenCode の MCP 構成を生成します。 Ollama モデルで OpenCode を別途設定します。 Ollama を備えた別の MCP 互換クライアントも適しています。

## ステップ 4: 接続を確認する

1. 従業員カードを開き、**もう一度確認** をクリックします。 **サーバー チェック** により、MCP ツールの数が表示されるはずです。これは、サーバー ツールの内部プロトコル チェックと検出です。
2. 構成ファイルを追加した後、AI クライアントを起動または再起動します。 `get_my_board` に電話してもらいます。
3. **クライアントが接続されました** が AgentBoard に表示され、新しいセッションがリストに表示されます。これだけでクライアントの接続が確認されます。ツールはあるがクライアントが接続されていない場合は、プログラムへのパス、構成ファイルの名前、およびその JSON/TOML 構文を確認してください。

`get_my_board` で作業セッションを開始します。 AI は、`Backlog` および `Complete` からのタスクは、ID がわかっていても受信しません。 AI は人間のふりをすることはできず、一般的な「どこにでも移動」コマンドはありません。 AgentBoard の外部でファイル システムを操作するかどうかは、選択したクライアントの機能とサンドボックスによって異なります。

公式の顧客向け指示: [Codex](https://developers.openai.com/learn/docs-mcp)、[クロード Code](https://docs.anthropic.com/en/docs/claude-code/mcp)、[Gemini CLI](https://geminicli.com/docs/tools/mcp-server/)、[Cursor](https://prod.cursor.com/docs/cli/mcp)、[OpenCode](https://opencode.ai/docs/mcp-servers/)、[VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers)]
