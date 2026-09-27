# Scoop 출시 준비

## 이미 자동화된 것

`scripts/package-scoop.ps1 -Version X.Y.Z`는 버전 번호가 있는 `npm ci`, 프런트엔드 빌드, `go test ./...`, Windows 빌드 `agentboard.exe`를 실행하고 ZIP을 생성하고 해당 ZIP의 SHA-256을 계산합니다. **동일** 파일에서 URL, 해시, MIT 라이센스, CLI-shim, 바로가기, `release/agentboard.json` 및 `checkver`를 사용하여 `autoupdate`를 생성합니다. ZIP 및 설치 프로그램에는 `LICENSE` 파일이 포함되어 있습니다. 데이터베이스는 설치 디렉터리 외부에 있으므로 매니페스트에 `persist`가 필요하지 않습니다.

`.github/workflows/release.yml` 태그의 `vX.Y.Z`는 Windows 러너에서 동일한 패키지를 실행하고 Inno 설치 프로그램을 빌드하며 ZIP, 설치 프로그램 및 매니페스트를 GitHub 릴리스에 연결합니다. 워크플로를 수동으로 실행하면 릴리스를 게시하지 않고 테스트용 아티팩트만 생성됩니다.

## 관리자를 위한 게시 명령

1. 코드, 문서, 버전 번호 및 파일 `LICENSE`가 준비되었는지 확인합니다.
2. `./scripts/package-scoop.ps1 -Version X.Y.Z`를 로컬에서 실행합니다. 포장을 풀고 `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` 및 `agentboard version` 출력을 확인하십시오. 해시가 계산된 후에는 ZIP을 변경하지 마세요.
3. `vX.Y.Z` 태그를 생성하고 제출합니다. GitHub Actions는 ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` 및 매니페스트가 포함된 릴리스를 게시합니다. 릴리스 페이지에서 세 파일을 모두 확인하고 매니페스트에서 SHA-256 ZIP을 확인하세요.
4. Scoop이 있는 깨끗한 Windows 시스템에서 `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`를 실행한 다음 테스트 폴더에서 `agentboard version`, `agentboard init` 및 `agentboard open`를 실행합니다.
5. 영구 업데이트 피드의 경우 생성된 `agentboard.json`를 자체 Scoop 버킷에 넣거나 적합한 공개 버킷에 제공하세요. 다음 릴리스 이후에 `scoop update agentboard`를 다시 확인하세요. URL을 통한 설치 명령은 처음 아는 사람에게 적합하고 버킷은 업데이트에 더 편리합니다.

임의의 값 `hash`를 수동으로 대체하지 마십시오. Scoop은 다운로드된 ZIP의 내용을 확인합니다.

매니페스트 확인의 경우 [형식 Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [매니페스트](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) 생성 및 [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
