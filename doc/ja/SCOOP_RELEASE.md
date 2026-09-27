# Scoop 用のリリースの準備

## すでに自動化されているもの

`scripts/package-scoop.ps1 -Version X.Y.Z` は、バージョン番号を使用して `npm ci`、フロントエンド ビルド、`go test ./...`、Windows ビルド `agentboard.exe` を実行し、ZIP を作成し、その ZIP の SHA-256 を計算します。 **同じ** ファイルから、URL、ハッシュ、MIT ライセンス、CLI シム、ショートカット、`release/agentboard.json` および `checkver` を含む `autoupdate` を作成します。 ZIP とインストーラーには、`LICENSE` ファイルが含まれています。データベースはインストール ディレクトリの外にあるため、マニフェストに `persist` は必要ありません。

タグ `.github/workflows/release.yml` の `vX.Y.Z` は、Windows ランナーで同じパッケージを実行し、Inno Setup インストーラーを構築し、ZIP、インストーラー、マニフェストを GitHub リリースに添付します。ワークフローを手動で実行すると、リリースは公開されず、テスト用のアーティファクトが作成されるだけです。

## メンテナへの投稿オーダー

1. コード、ドキュメント、バージョン番号、およびファイル `LICENSE` の準備ができていることを確認します。
2. `./scripts/package-scoop.ps1 -Version X.Y.Z` をローカルで実行します。開梱後、`release/agentboard-X.Y.Z-windows-amd64.zip`、`release/agentboard.json`、`agentboard version` の出力を確認してください。ハッシュの計算後は ZIP を変更しないでください。
3. タグ `vX.Y.Z` を作成して送信します。 GitHub Actions は、ZIP、`agentboard-X.Y.Z-windows-amd64-setup.exe`、およびマニフェストを含むリリースを公開します。リリース ページの 3 つのファイルすべてとマニフェストの SHA-256 ZIP を確認してください。
4. Scoop を備えたクリーンな Windows マシンで、`scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json` を実行し、次にテスト フォルダー内の `agentboard version`、`agentboard init`、および `agentboard open` を実行します。
5. 永続的な更新フィードの場合は、生成された `agentboard.json` を独自の Scoop バケットに配置するか、適切な公開バケットに提供します。次のリリース後に `scoop update agentboard` を確認してください。 URL からのインストール コマンドは最初の知人に適していますが、バケットは更新の場合により便利です。

ランダム値 `hash` を手動で置き換えないでください: Scoop は、ダウンロードされた ZIP の内容をチェックします。

マニフェスト チェックの場合は、[format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests)、[マニフェスト](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest)の作成]、[autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate)。
