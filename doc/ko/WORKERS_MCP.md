# 작업자 및 MCP 연결

## 1단계. 작업자 생성

AgentBoard에서 **작업자 → 작업자 추가**를 엽니다. 클라이언트 프로필을 선택합니다: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, OpenCode를 통한 Ollama, VS Code Copilot 또는 다른 MCP 클라이언트. 프로필은 일반적인 기능(코드, 테스트, Git)을 채웁니다. 변경할 수 있습니다. 이름은 보드에 표시되며 `Slug`는 MCP 명령을 위한 공백이 없는 짧은 식별자입니다. **저장**을 클릭합니다.

작업자 한 명은 AI 성격 한 명에 해당합니다. 다양한 클라이언트나 팀을 위해 다양한 작업자를 만듭니다. 기능 프로필은 전문화를 설명하지만 AI의 권한을 작업 단계로 확장하지는 않습니다.

## 2단계. 구성 복사

생성된 작업자를 엽니다. **MCP 진단** 블록에는 구성 조각과 이를 추가할 파일이 표시됩니다. **MCP 구성 복사**를 클릭합니다. 파일이 이미 존재하는 경우 다른 서버를 삭제하지 않고 기존 `mcpServers`/`servers`/`mcp` 개체에 제안된 서버를 추가합니다.

주요 명령은 다음과 같습니다.

```text
agentboard mcp --worker codex
```

`--worker` — 생성된 작업자의 `Slug`입니다. MCP 클라이언트는 도구가 필요할 때 이 명령을 자체적으로 실행합니다.

### MCP는 웹 인터페이스 없이 작동합니다.

`agentboard mcp`는 자체적으로 AgentBoard 데이터베이스를 여는 별도의 stdio 프로세스입니다. `agentboard open`/`serve`가 필요하지 않으며 웹 인터페이스가 실행되는 폴더와 실행 중인 웹 인터페이스 인스턴스 수에 의존하지 않습니다. 여러 에이전트가 각각 자체 프로젝트에서 동시에 실행될 수 있습니다. 각 클라이언트는 자체 MCP 프로세스를 실행하며 모두 동일한 SQLite 데이터베이스를 안전하게 사용합니다.

### 프로젝트 선택 방법

MCP 서버는 모든 폴더에서 실행됩니다. 도구를 처음 호출하면 프로젝트에 바인딩됩니다.

1. `--project PATH`(지정된 경우)(등록된 프로젝트 내의 모든 폴더가 적합함)
2. 클라이언트가 보고하는 경우 클라이언트의 작업 폴더(루트)
3. 그렇지 않으면 클라이언트가 MCP 서버를 시작한 폴더입니다.

`get_my_board`는 선택한 프로젝트와 작업자를 `session`와 함께 `binding_source`에 반환합니다. 클라이언트의 작업 폴더가 변경되면 다음에 도구를 호출할 때 바인딩이 다시 수행됩니다. `--worker`가 없으면 프로젝트에서 활성화된 유일한 작업자가 사용됩니다.

바인딩이 불가능한 경우 각 도구는 에이전트가 사용자에게 전달해야 하는 오류를 반환합니다. 이 경우 작업 데이터를 읽을 수 없습니다.

| 코드 | 의미 |
| --- | --- |
| `NO_PROJECTS` | 데이터베이스에 프로젝트가 없습니다. 프로젝트 폴더에서 `agentboard init`를 실행하고 워커를 생성합니다. 다시 시작할 필요가 없습니다. |
| `PROJECT_REQUIRED` | 작업 폴더가 등록된 프로젝트 외부에 있습니다. 메시지에는 프로젝트가 나열됩니다. |
| `PROJECT_NOT_FOUND` | `--project`는 등록된 프로젝트 내에 없습니다. |
| `PROJECT_AMBIGUOUS` | 클라이언트의 작업 폴더는 여러 프로젝트에 걸쳐 있습니다. `--project`를 추가하세요. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | `--worker`를 통해 작업자를 생성, 활성화 또는 선택합니다. |

**Codex** 및 **Gemini CLI**에 대한 구성은 사용자 정의되고 모든 폴더에서 공유되므로 `--project`가 없습니다. 이 클라이언트가 실행 중인 모든 프로젝트에서 동일한 슬러그를 사용하여 작업자를 만듭니다. 프로젝트 레벨 파일(Claude Code, Cursor, OpenCode, VS Code)에는 해당 프로젝트의 `--project`가 포함되어 있습니다.

### 명시적 디자인 및 컨텍스트 확인

`--project`가 지정되고 클라이언트의 시작 폴더 또는 작업 폴더가 다른 프로젝트에 있는 경우 도구는 MCP 추출을 통해 확인을 요청합니다. 양식 추출이 없는 클라이언트는 작업에 액세스하기 전에 `PROJECT_CONTEXT_MISMATCH`를 받습니다. 구성을 수정하고, 자동 감지를 위해 `--project`를 제거하거나, `--allow-project-mismatch` 플래그를 사용하여 다른 프로젝트에서 작업할 계획임을 명시적으로 확인하세요. 작업 폴더를 변경하면 이전 확인이 재설정됩니다. 자동으로 정의된 프로젝트에는 확인이 필요하지 않습니다.

**전역 설정 → MCP 프로젝트 컨텍스트**는 동일한 데이터베이스의 모든 프로젝트 및 작업자에 대한 이러한 검사를 제어합니다. 기본적으로 활성화되어 있습니다. 비활성화된 경우 MCP는 확인이나 불일치 오류 없이 지정된 프로젝트를 사용합니다. 이미 열려 있는 세션은 다음에 도구를 호출할 때 설정을 적용합니다. 작업 가시성, 전환 및 인간 테스트 규칙은 계속해서 준수됩니다.

확인한 후 진단 블록은 실행 중인 `agentboard.exe`에 대한 절대 경로를 표시합니다. ZIP으로 설치할 때 특히 유용합니다. Scoop을 통해 설치할 때 클라이언트에 동일한 `agentboard`가 표시되면 `PATH` 명령을 사용할 수 있습니다.

## 3단계: 클라이언트에 서버 추가

작업자 인터페이스에는 이미 준비된 조각이 있습니다. 아래는 사용처에 대한 설명입니다.

| 클라이언트 | 삽입할 위치 | 클라이언트 측에서 확인하는 방법 |
| --- | --- | --- |
| 코덱스 | `%USERPROFILE%\.codex\config.toml`, 섹션 `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| 클로드 코드 | 프로젝트 폴더의 `.mcp.json` | `claude mcp list` |
| 제미니 CLI | `%USERPROFILE%\.gemini\settings.json`, 개체 `mcpServers` | Gemini CLI의 `/mcp list` |
| 커서 | `.cursor\mcp.json` 프로젝트 | 커서 설정의 MCP 서버 목록 |
| 오픈코드 | `opencode.json` 프로젝트, 개체 `mcp` | OpenCode의 MCP 도구 목록 |
| VS 코드 부조종사 | `.vscode\mcp.json` 프로젝트, 개체 `servers` | 명령 **MCP: 서버 나열** |

프로젝트 `C:\Projects\MyFirstProject` 및 작업자 `codex`가 포함된 Codex(사용자 구성, 자동으로 감지된 프로젝트) 및 Claude Code(프로젝트 파일)의 예:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
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

JSON에서는 Windows 백슬래시가 두 배가 됩니다. 인터페이스에서 미리 만들어진 조각이 이 작업을 자동으로 수행합니다. 생성된 구성에는 웹 서버의 실제 `--db`가 포함되므로 클라이언트 환경이 다르더라도 HTTP와 MCP는 동일한 기반을 사용합니다.

**Ollama**는 로컬 모델을 제공하지만 MCP 클라이언트를 대체하지는 않습니다. **OpenCode를 통한 Ollama** 프로필은 OpenCode에 대한 MCP 구성을 생성합니다. Ollama 모델에 OpenCode를 별도로 구성합니다. Ollama를 사용하는 또 다른 MCP 호환 클라이언트도 적합합니다.

## 4단계: 연결 확인

1. 근로자 카드를 열고 **다시 확인**을 클릭하세요. **서버 확인**은 이 프로젝트, 작업자 및 베이스를 사용하여 실제 stdio 실행 파일을 실행하고 핸드셰이크, 도구 목록, 프로젝트 및 작업자 위치를 확인합니다.
2. 구성 파일을 추가한 후 AI 클라이언트를 시작하거나 다시 시작합니다. 그에게 `get_my_board`로 전화해달라고 요청하세요.
3. **클라이언트 연결**이 AgentBoard에 나타나고 새 세션이 목록에 나타납니다. 이것으로만 클라이언트 연결이 확인됩니다. 도구는 있지만 클라이언트가 연결되지 않은 경우 프로그램 경로, 구성 파일 이름 및 해당 JSON/TOML 구문을 확인하세요.

`get_my_board`로 작업 세션을 시작하세요. AI는 ID를 알고 있더라도 `Backlog` 및 `Complete`로부터 작업을 수신하지 않습니다. AI는 인간인 척할 수 없으며 일반적인 "어디로 이동" 명령도 없습니다. AgentBoard 외부의 파일 시스템 작업은 선택한 클라이언트의 기능과 샌드박스에 따라 다릅니다.

각 도구는 `title` 및 네 가지 MCP 동작 힌트(`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`)를 모두 명시적인 부울 값으로 선언합니다. 읽기 도구는 읽기 전용이며 멱등성을 갖습니다. 녹음 도구는 레코드를 추가하거나 작업을 진행하기만 하므로 방해가 되거나 로컬 AgentBoard 기반을 벗어나는 도구는 없습니다.

공식 고객 지침: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
