# エージェントボード

[🌐 Languages](../LANGUAGES.md)

AgentBoard は、人間と AI ワーカーが共通のタスクに取り組むローカル タスク ボードですが、異なる権限を持っています。アプリケーションは 1 つのファイル `agentboard.exe` で起動され、ブラウザで Web インターフェイスが開き、`stdio` 経由でワーカーに別の MCP サーバーを提供します。データはコンピュータ上に残ります。

**[ゼロから始める](START_HERE.md) · [タスクの操作](TASKS.md) · [MCP](WORKERS_MCP.md)経由でAIに接続する] · [JSON](IMPORT.md)をインポートする · [設定](SETTINGS.md) · [問題を解決する](TROUBLESHOOTING.md)**

## 5分以内に

1. [リリース ](https://github.com/bogdanovandreycode/AgentBoard/releases)] から `agentboard-VERSION-windows-amd64-setup.exe` インストーラーをダウンロードします。フォルダー (デフォルトでは `C:\AI\AgentBoard`) が提案され、それを `PATH` に追加します。スクープやZIPもございます。
2. プロジェクト フォルダー (`C:\Projects\MyApp` など) で PowerShell を開きます。
3. `agentboard init` (ZIP の場合: `agentboard.exe` および `init` へのフルパス) を実行します。
4. `agentboard open`を実行します。 `http://127.0.0.1:7337` が開きます。
5. [**新しいタスク**] ボタンを使用してタスクを追加します。 AI ワーカーの場合は、**ワーカー → ワーカーの追加** を開き、クライアント プロファイルを選択して MCP 構成をコピーします。

プロジェクト フォルダーがまだない場合は、Windows エクスプローラーで作成します。プロジェクトには、Git やコードがなくても、任意のフォルダーを使用できます。

## Scoop によるインストール

リリース後に [Scoop](https://scoop.sh/)] がすでにインストールされている PowerShell の場合:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

開発者向けには [ソース ](INSTALL.md) からのビルド] があります。リリース ワークフローでは、同じアーティファクトから SHA-256 を使用したインストーラー、ZIP、および Scoop マニフェストが作成されます。バケットにマニフェストを追加した後、Scoop: `scoop update agentboard` 経由でインストールされたバージョンを更新します。詳細 - 【発売準備中](SCOOP_RELEASE.md)。

## ボードの構造

`Backlog → Features → In progress → Testing → Verification → Complete`

ユーザーはボード上でタスクを移動できます。 AI は `Features → In progress → Testing → Verification` のみを移動できます。 `Complete` への最終的な受け入れは人間によって行われます。ユーザー列は人間を対象としています。ユーザー列内のタスクは MCP の `Backlog` 状態のままです。テストモード: AI、ヒューマン、ハイブリッド。履歴、テスト、成果物、AI コストがタスクに関連付けられます。

## チーム

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` はフォルダーを登録し、そこに `.agentboard/project.json` のみを書き込みます。 SQLite 実稼働データは、プロジェクトの外部および Scoop インストールの外部の Windows ユーザー構成ディレクトリ (`%AppData%\AgentBoard\agentboard.db`) にあります。パッケージを削除または更新しても、このデータは削除されません。別のコンピュータに転送する前に、AgentBoard を停止した状態でデータベースのコピーを作成してください。

労働者 - 論理アカウント。 AgentBoard 自体は、Codex、Claude、またはその他の AI クライアントを実行しません。クライアントは、特定のワーカーに対してローカル MCP プロセスを開始します。 MCP はアプリケーションのアクセス許可境界であり、ファイルの分離には AI クライアント サンドボックスを使用します。

## 開発者向け

スタック: Go、SQLite、公式 MCP Go SDK、React、TypeScript、Vite、PrimeReact、TanStack Query、dnd-kit。最初にフロントエンドを構築してから、`./scripts/build.ps1` に進みます。 Web ファイルは、`go:embed` 経由でバイナリに含まれます。アーキテクチャと API については、[doc/ARCHITECTURE.md](ARCHITECTURE.md).
