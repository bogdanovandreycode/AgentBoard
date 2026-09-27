# Windows에서 설치 및 실행

## 설치 프로그램(권장)

[Releases](https://github.com/bogdanovandreycode/AgentBoard/releases)] 페이지에서 `agentboard-VERSION-windows-amd64-setup.exe`를 다운로드하세요. 설치 마법사가 폴더를 제안합니다. 기본값은 `C:\AI\AgentBoard`입니다. `agentboard.exe` 및 문서를 복사하고 바로 가기를 생성하며 선택한 폴더를 `PATH` 시스템에 추가합니다. 설치 후 `agentboard` 명령을 사용할 수 있도록 새 터미널을 엽니다.

프로젝트 폴더에서 다음을 실행합니다.

```powershell
agentboard init
agentboard open
```

인터페이스는 `agentboard.exe`에 내장되어 있습니다. Go 또는 Node.js를 별도로 설치할 필요가 없습니다. 데이터는 `%AppData%\AgentBoard`에 저장되며 프로그램이 업데이트되거나 제거될 때 유지됩니다. "설치된 응용 프로그램"을 통해 제거하면 바로가기와 `PATH`의 항목이 제거됩니다.

## 스쿠프 설치하기

Scoop이 아직 설치되지 않은 경우 일반 사용자로 PowerShell을 열고 [공식 Scoop](https://scoop.sh/) 지침을 따르세요. 회사 컴퓨터에 제한 사항이 있는 경우 관리자에게 문의하세요. AgentBoard는 Scoop 없이 ZIP에서 시작할 수도 있습니다.

또는 Scoop을 통해 AgentBoard를 설치하십시오.

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

GitHub 릴리스의 매니페스트 사용 가능 여부는 [페이지 릴리스](https://github.com/bogdanovandreycode/AgentBoard/releases). Scoop에 매니페스트를 추가한 후 버킷 이름으로 버킷을 설치하고 `scoop update agentboard` 명령을 사용하여 업데이트할 수 있습니다.

## 스쿠프 없이 ZIP

릴리스에서 `agentboard-VERSION-windows-amd64.zip`를 다운로드하고, 예를 들어 `C:\Tools\AgentBoard`에 압축을 풉니다. PowerShell의 프로젝트 폴더:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

AI 클라이언트의 경우 프로그램이 `agentboard.exe`에 없으면 MCP 구성에서 `PATH`에 대한 전체 경로를 지정합니다.

## 소스에서 빌드

`go.mod` 및 Node.js 22 이상의 Go 버전을 설치하세요. 저장소 루트에 있는 PowerShell에서:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1`는 프런트엔드를 `internal/webui/dist`로 조립하고 Go 테스트를 실행한 후 하나의 `agentboard.exe`를 조립합니다. 순서가 중요합니다. 인터페이스는 Go가 빌드될 때 바이너리에 내장됩니다. 이전 `agentboard.exe open`가 실행 중인 경우 다시 빌드하기 전에 중지하십시오(Ctrl+C). 그렇지 않으면 Windows에서 파일 교체를 허용하지 않습니다.

## 데이터는 어디에 있나요?

- `%AppData%\AgentBoard\agentboard.db` - 작업, 프로젝트, 작업자 및 설정. `--db` 플래그를 사용하여 다른 파일을 지정할 수 있지만 `init`, `open`/`serve` 및 `mcp`는 **동일한 경로**를 가져야 합니다.
- `<ваш проект>\.agentboard\project.json` — 프로젝트 식별자. 이 파일에는 작업이 포함되어 있지 않습니다.
- 서버는 기본적으로 `127.0.0.1:7337`만 수신합니다. 프로젝트 경로 `--addr` 앞에 다른 주소 `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`를 지정합니다.

`open`는 프로젝트 페이지를 열고 해당 주소에서 이미 실행 중인 서버를 재사용합니다. 하나의 브라우저에 여러 프로젝트가 열려 있는 경우 왼쪽 목록을 통해 선택하세요.
