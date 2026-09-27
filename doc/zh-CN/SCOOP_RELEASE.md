# 为 Scoop 准备一个版本

## 哪些已经自动化

`scripts/package-scoop.ps1 -Version X.Y.Z` 执行 `npm ci`、前端版本、`go test ./...`、带有版本号的 Windows 版本 `agentboard.exe`，创建 ZIP 并计算该 ZIP 的 SHA-256。它从**同一个**文件创建带有 URL、哈希、MIT 许可证、CLI-shim、快捷方式、`release/agentboard.json` 和 `checkver` 的 `autoupdate`。 ZIP 和安装程序包含文件 `LICENSE`。数据库位于安装目录之外，因此清单中不需要 `persist`。

标签 `.github/workflows/release.yml` 上的 `vX.Y.Z` 在 Windows 运行程序中运行相同的包，构建 Inno Setup 安装程序并将 ZIP、安装程序和清单附加到 GitHub 版本。手动运行工作流只会创建用于测试的工件，而不会发布版本。

## 维护者发布订单

1. 验证代码、文档、版本号和文件 `LICENSE` 是否准备就绪。
2. 本地运行`./scripts/package-scoop.ps1 -Version X.Y.Z`。开箱后检查 `release/agentboard-X.Y.Z-windows-amd64.zip`、`release/agentboard.json` 和 `agentboard version` 输出。计算哈希值后请勿更改 ZIP。
3. 创建并提交标签`vX.Y.Z`。 GitHub Actions 将发布带有 ZIP、`agentboard-X.Y.Z-windows-amd64-setup.exe` 和清单的版本。检查“发布”页面上的所有三个文件以及清单中的 SHA-256 ZIP。
4. 在装有 Scoop 的干净 Windows 计算机上，运行 `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`，然后运行测试文件夹中的 `agentboard version`、`agentboard init` 和 `agentboard open`。
5. 对于永久更新源，请将生成的 `agentboard.json` 放入您自己的 Scoop 存储桶中或将其提供到合适的公共存储桶中。下一个版本后请回来查看 `scoop update agentboard`。 URL安装命令适合初次接触，而bucket则更方便更新。

不要手动替换随机值 `hash`：Scoop 检查下载的 ZIP 的内容。

对于清单检查，请重点关注[格式Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests)、[创建清单](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest)]和[自动更新](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate)。
