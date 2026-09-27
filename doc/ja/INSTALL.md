# Windows でのインストールと起動

## インストーラー (推奨)

[](https://github.com/bogdanovandreycode/AgentBoard/releases) のリリース] ページから `agentboard-VERSION-windows-amd64-setup.exe` をダウンロードします。インストール ウィザードはフォルダーを提案します。デフォルトは`C:\AI\AgentBoard`です。 `agentboard.exe` とドキュメントをコピーし、ショートカットを作成し、選択したフォルダーをシステム `PATH` に追加します。インストール後、新しいターミナルを開くと、`agentboard` コマンドが使用可能になります。

プロジェクト フォルダーで次を実行します。

```powershell
agentboard init
agentboard open
```

インターフェースは `agentboard.exe` に組み込まれています。 Go や Node.js を個別にインストールする必要はありません。データは `%AppData%\AgentBoard` に保存され、プログラムの更新またはアンインストール時に保持されます。 「インストールされているアプリケーション」からアンインストールすると、ショートカットと `PATH` からのエントリが削除されます。

## スコップのインストール

Scoop がまだインストールされていない場合は、通常のユーザーとして PowerShell を開き、[Scoop](https://scoop.sh/) の公式手順に従ってください。会社のコンピュータに制限がある場合は、管理者に問い合わせてください。 AgentBoard は、Scoop を使用せずに ZIP から起動することもできます。

あるいは、Scoop 経由で AgentBoard をインストールします。

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

GitHub リリースでのマニフェストの可用性は、[リリース](https://github.com/bogdanovandreycode/AgentBoard/releases)ページ]で確認できます。マニフェストを Scoop に追加した後、バケット名を使用してバケットをインストールし、`scoop update agentboard` コマンドで更新できます。

## スクープなしの ZIP

リリースから `agentboard-VERSION-windows-amd64.zip` をダウンロードし、たとえば `C:\Tools\AgentBoard` に解凍します。 PowerShell のプロジェクト フォルダーで次のようにします。

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

AI クライアントの場合、プログラムが `agentboard.exe` にない場合は、MCP 構成で `PATH` へのフル パスを指定します。

## ソースからビルドする

`go.mod` の Go バージョンと Node.js 22 以降をインストールします。リポジトリのルートにある PowerShell で次のようにします。

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` はフロントエンドを `internal/webui/dist` にアセンブルし、Go テストを実行して 1 つの `agentboard.exe` をアセンブルします。順序は重要です。インターフェイスは Go のビルド時にバイナリに組み込まれます。古い `agentboard.exe open` が実行中の場合は、再構築する前に停止してください (Ctrl+C)。そうしないと、Windows でファイルを置き換えることができません。

## データはどこにありますか

- `%AppData%\AgentBoard\agentboard.db` - タスク、プロジェクト、ワーカー、および設定。 `--db` フラグを使用して別のファイルを指定できますが、`init`、`open`/`serve`、および `mcp` は **同じパス**である必要があります。
- `<ваш проект>\.agentboard\project.json` — プロジェクト識別子。このファイルにはタスクが含まれていません。
- サーバーはデフォルトで `127.0.0.1:7337` のみをリッスンします。プロジェクト パス `--addr` の前に別のアドレス `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject` を指定します。

`open` はプロジェクト ページを開き、そのアドレスですでに実行されているサーバーを再利用します。 1 つのブラウザで複数のプロジェクトが開いている場合は、左側のリストからそれらを選択します。
