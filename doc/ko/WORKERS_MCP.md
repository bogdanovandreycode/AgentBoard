# 작업자 및 MCP 연결

## 1단계. 작업자 생성

AgentBoard에서 **작업자 → 작업자 추가**를 엽니다. 클라이언트 프로필을 선택합니다: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, OpenCode를 통한 Ollama, VS Code Copilot 또는 다른 MCP 클라이언트. 프로필은 일반적인 기능(코드, 테스트, Git)을 채웁니다. 변경할 수 있습니다. 이름은 보드에 표시되며 `Slug`는 MCP 명령을 위한 공백이 없는 짧은 식별자입니다. **저장**을 클릭합니다.

작업자 한 명은 AI 성격 한 명에 해당합니다. 다양한 클라이언트나 팀을 위해 다양한 작업자를 만듭니다. 기능 프로필은 전문화를 설명하지만 AI의 권한을 작업 단계로 확장하지는 않습니다.

## 2단계. 구성 복사

생성된 작업자를 엽니다. **MCP 진단** 블록에는 구성 조각과 이를 추가할 파일이 표시됩니다. **MCP 구성 복사**를 클릭합니다. 파일이 이미 존재하는 경우 다른 서버를 삭제하지 않고 기존 `mcpServers`/`servers`/`mcp` 개체에 제안된 서버를 추가합니다.

주요 명령은 다음과 같습니다.

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project`는 `agentboard init`를 통해 등록한 폴더를 가리켜야 합니다. `--worker` — 생성된 작업자의 `Slug`입니다. MCP 클라이언트는 도구가 필요할 때 이 명령을 자체적으로 실행합니다. AgentBoard 브라우저에서는 웹 서버를 별도로 실행할 수 있습니다.

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

Codex 및 Claude Code의 경우 프로젝트 `C:\Projects\MyFirstProject` 및 작업자 `codex`의 예:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

JSON에서는 Windows 백슬래시가 두 배가 됩니다. 인터페이스에서 미리 만들어진 조각이 이 작업을 자동으로 수행합니다. `--db`를 커스텀 베이스와 함께 사용하는 경우 `args` MCP 구성에 추가하고 `open`를 시작할 때와 동일한 경로를 지정합니다.

**Ollama**는 로컬 모델을 제공하지만 MCP 클라이언트를 대체하지는 않습니다. **OpenCode를 통한 Ollama** 프로필은 OpenCode에 대한 MCP 구성을 생성합니다. Ollama 모델에 OpenCode를 별도로 구성합니다. Ollama를 사용하는 또 다른 MCP 호환 클라이언트도 적합합니다.

## 4단계: 연결 확인

1. 근로자 카드를 열고 **다시 확인**을 클릭하세요. **서버 확인**에는 MCP 도구 수가 표시되어야 합니다. 이는 내부 프로토콜 검사 및 서버 도구 탐지입니다.
2. 구성 파일을 추가한 후 AI 클라이언트를 시작하거나 다시 시작합니다. 그에게 `get_my_board`에 전화하라고 요청하세요.
3. **클라이언트 연결**이 AgentBoard에 나타나고 새 세션이 목록에 나타납니다. 이것으로만 클라이언트 연결이 확인됩니다. 도구는 있지만 클라이언트가 연결되지 않은 경우 프로그램 경로, 구성 파일 이름 및 해당 JSON/TOML 구문을 확인하세요.

`get_my_board`로 작업 세션을 시작하세요. AI는 ID를 알고 있더라도 `Backlog` 및 `Complete`로부터 작업을 수신하지 않습니다. AI는 인간인 척할 수 없으며 일반적인 "어디로 이동" 명령도 없습니다. AgentBoard 외부의 파일 시스템 작업은 선택한 클라이언트의 기능과 샌드박스에 따라 다릅니다.

공식 고객 지침: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
