# 에이전트보드

![미리보기 AgentBoard](../../assets/social-preview.png)

**[Windows용 다운로드](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [문서 사이트](https://bogdanovandreycode.github.io/AgentBoard/) · [라이센스 MIT](../../LICENSE)**

AgentBoard는 인간과 AI 작업자가 공통 작업을 수행하지만 서로 다른 권한을 갖는 로컬 작업 보드입니다. 이 애플리케이션은 하나의 파일 `agentboard.exe`로 시작되고 브라우저에서 웹 인터페이스를 열고 `stdio`를 통해 작업자에게 별도의 MCP 서버를 제공합니다. 데이터는 컴퓨터에 남아 있습니다.

**[처음부터 시작](START_HERE.md) · [tasks](TASKS.md) 작업 · [MCP](WORKERS_MCP.md)를 통해 AI 연결 · [JSON](IMPORT.md) 가져오기 · [설정](SETTINGS.md) · [문제 해결](TROUBLESHOOTING.md)**

## MVP 기능

- 작업 검색, JSON 가져오기, 사용자 정의 열 및 속성이 포함된 로컬 프로젝트 보드.
- AI 전환이 제한적이고 최종 인간이 작업을 수락하는 특정 작업자에 대한 MCP 액세스입니다.
- 일반 작업 내역, 확인 지침, 아티팩트, AI 비용 및 작업자 연결 진단.
- Windows 설치 프로그램, ZIP 및 Scoop 매니페스트; 웹 인터페이스는 실행 파일에 내장되어 있습니다.

AgentBoard는 신뢰할 수 있는 로컬 사용자를 위해 설계되었습니다. 애플리케이션은 클라우드에서 프로젝트를 호스팅하지 않으며 AI 클라이언트 자체를 시작하지 않습니다. 필요한 경우 MCP 호환 클라이언트를 작업자에 연결합니다.

## 5분 후

1. [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases).dll]에서 `agentboard-VERSION-windows-amd64-setup.exe` 설치 프로그램을 다운로드합니다. 폴더(기본적으로 `C:\AI\AgentBoard`)를 제안하고 이를 `PATH`에 추가합니다. Scoop과 ZIP도 가능합니다.
2. 프로젝트 폴더(예: `C:\Projects\MyApp`)에서 PowerShell을 엽니다.
3. `agentboard init`를 실행합니다(ZIP의 경우: `agentboard.exe` 및 `init`에 대한 전체 경로).
4. `agentboard open`를 실행합니다. `http://127.0.0.1:7337`가 열립니다.
5. **새 작업** 버튼을 사용하여 작업을 추가합니다. AI 작업자의 경우 **작업자 → 작업자 추가**를 열고 클라이언트 프로필을 선택한 후 MCP 구성을 복사합니다.

프로젝트 폴더가 아직 없으면 Windows 탐색기에서 만듭니다. Git과 코드가 없어도 어떤 폴더든 프로젝트가 될 수 있습니다.

## Scoop을 통한 설치

릴리스 이후에 [Scoop](https://scoop.sh/)]가 이미 설치된 PowerShell에서:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

개발자에게는 [소스 ](INSTALL.md)에서 빌드]가 있습니다. 릴리스 워크플로는 동일한 아티팩트에서 SHA-256을 사용하여 설치 프로그램, ZIP 및 Scoop 매니페스트를 생성합니다. 버킷에 매니페스트를 추가한 후 Scoop: `scoop update agentboard`를 통해 설치된 버전 업데이트 세부정보 - [출시 준비](SCOOP_RELEASE.md).

## 보드는 어떻게 구성되어 있나요?

`Backlog → Features → In progress → Testing → Verification → Complete`

사람은 보드 주위에서 작업을 이동할 수 있습니다. AI는 `Features → In progress → Testing → Verification`만 이동할 수 있습니다. `Complete`에 대한 최종 승인은 인간에 의해 수행됩니다. 사용자 열은 사람을 위한 것입니다. 해당 열의 작업은 MCP에 대해 `Backlog` 상태로 유지됩니다. 테스트 모드: AI, 인간 및 하이브리드. 기록, 테스트, 아티팩트 및 AI 비용이 작업에 첨부됩니다.

## 팀

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init`는 폴더를 등록하고 거기에 `.agentboard/project.json`만 씁니다. SQLite 프로덕션 데이터는 Windows 사용자 구성 디렉터리(`%AppData%\AgentBoard\agentboard.db`), 프로젝트 외부 및 Scoop 설치 외부에 있습니다. 패키지를 제거하거나 업데이트해도 이 데이터는 제거되지 않습니다. 다른 컴퓨터로 전송하기 전에 AgentBoard가 중지된 동안 데이터베이스의 복사본을 만드십시오.

작업자 - 논리 계정 AgentBoard 자체는 Codex, Claude 또는 기타 AI 클라이언트를 실행하지 않습니다. 클라이언트는 특정 작업자에 대한 로컬 MCP 프로세스를 시작합니다. MCP는 애플리케이션 권한 경계이며, 파일 격리를 위해 AI 클라이언트 샌드박스를 사용합니다.

## 개발자를 위한

스택: Go, SQLite, 공식 MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. 먼저 프런트엔드를 빌드한 다음 Go: `./scripts/build.ps1`를 작성하세요. 웹 파일은 `go:embed`를 통해 바이너리에 포함됩니다. 아키텍처와 API는 [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## 라이센스

AgentBoard는 [라이센스 MIT](../../LICENSE)]에 따른 무료 오픈 소스 프로젝트입니다. 저작권 표시와 라이센스가 유지되는 경우 상업적 사용, 수정, 포크 및 재배포가 허용됩니다.
