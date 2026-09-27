# Preparing a release for Scoop

## What is already automated

`scripts/package-scoop.ps1 -Version 0.2.0` runs `npm ci`, frontend build, `go test ./...`, Windows build `agentboard.exe` with version number, creates a ZIP and calculates the SHA-256 of that ZIP. From **the same** file it creates `release/agentboard.json` with URL, hash, CLI-shim, shortcut, `checkver` and `autoupdate`. The database lives outside the installation directory, so `persist` in manifest is not needed.

`.github/workflows/release.yml` on the `vX.Y.Z` tag runs the same packaging in the Windows runner, builds the Inno Setup installer and attaches the ZIP, installer and manifest to the GitHub Release. Manually running a workflow only creates an artifact for testing, without publishing a release.

## Posting order for maintainer

1. Check that the code, documentation and version number are ready. Determine the project's license: manifest currently indicates `Unknown` because there is no license file defined in the repository. If you choose a license, add `LICENSE` and update `license` in the script before release.
2. Run `./scripts/package-scoop.ps1 -Version X.Y.Z` locally. Check `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` and the output of `agentboard version` after extracting. Do not change the ZIP after the hash has been calculated.
3. Create and submit the `vX.Y.Z` tag. GitHub Actions will publish a Release with a ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` and manifest. Check all three files on the Release page and the SHA-256 ZIP from manifest.
4. On a clean Windows machine with Scoop, run `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, then `agentboard version`, `agentboard init` and `agentboard open` in the test folder.
5. For a permanent update feed, place the generated `agentboard.json` in your own Scoop bucket or offer it in a suitable public bucket. Check `scoop update agentboard` after the next release. The installation command from a URL is suitable for the first acquaintance, while the bucket is more convenient for updates.

Do not manually substitute a random `hash` value: Scoop checks the contents of the downloaded ZIP.

For manifest checks, refer to [Scoop format](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [creating manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) and [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
