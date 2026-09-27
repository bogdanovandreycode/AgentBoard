# 아키텍처와 API

AgentBoard는 SQLite를 실행하는 로컬 Go 프로세스입니다. 인간을 위한 HTTP API, AI를 위한 MCP 어댑터, 웹 인터페이스는 공통 서비스 로직을 공유합니다. SQLite - 상태 소스; 프런트엔드는 서버를 우회하지 않습니다. AI 도구에는 별도의 권한 표면이 있으며 MCP의 `Backlog`/`Complete` 상태를 사용할 수 없습니다. `System`의 기록 레코드는 애플리케이션 자체에서 생성됩니다.

## 구성 요소

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - 도메인 모델 및 오류.
- `internal/service` - 작업 전환, 권한, 가져오기 및 설정.
- `internal/persistence` - SQLite 및 마이그레이션.
- `internal/httpapi` - 휴먼 HTTP API.
- `internal/mcpserver` - 별도의 작업자를 위한 MCP 도구입니다.
- `web` - 반응/타입스크립트 UI; 어셈블리는 `internal/webui/dist`로 끝나고 EXE에 포함됩니다.

## 기본 HTTP 경로

| 방법 및 경로 | 목적지 |
| --- | --- |
| `GET /api/health` | 서버 점검. |
| `GET /api/projects` | 등록된 프로젝트입니다. |
| `GET /api/projects/{id}/board` | 프로젝트 보드. |
| `GET/PUT /api/projects/{id}/settings` | 열의 설정, 순서 및 이름입니다. |
| `POST /api/projects/{id}/tasks` | 작업을 만듭니다. |
| `POST /api/projects/{id}/tasks/import` | JSON 버전 1을 원자적으로 가져옵니다. |
| `GET/PATCH/DELETE /api/tasks/{id}` | 카드, 변경, 삭제. |
| `POST /api/tasks/{id}/move` | 사람이 이동합니다. |
| `GET/POST /api/projects/{id}/workers` | 작업자 목록 및 생성. |
| `GET /api/workers/{id}/mcp/check` | MCP 핸드셰이크 및 도구의 내부 테스트입니다. |
| `GET /api/workers/{id}/sessions` | 세션 진단. |
| `GET/POST /api/projects/{id}/properties` | 사용자 정의 속성. |

HTTP는 신뢰할 수 있는 로컬 사용자를 위한 것입니다. 자체 인증, 네트워크 제한 및 HTTPS 없이 웹 포트를 인터넷에 게시하지 마십시오. MCP 서버는 특정 프로젝트 및 작업자에 대해 `stdio`를 사용하여 시작됩니다. `get_my_board`로 시작하세요. 해당 권한은 AgentBoard 내의 작업을 제한하지만 AI 클라이언트의 파일 시스템 샌드박스를 대체하지는 않습니다.

사용자 열은 `tasks.state`와 별도로 저장됩니다. 코어는 `backlog`의 해당 열에 작업을 유지하고 `tasks.board_column`는 휴먼 보드의 위치를 ​​결정합니다. 이는 동일한 AI 전환 모델을 유지합니다. 열을 제거하면 `board_column` 작업이 지워집니다.
