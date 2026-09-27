# JSON에서 작업 가져오기

**가져오기 작업** 탭에서는 UTF-8 인코딩의 JSON 파일을 허용합니다. 다른 서비스에서 작업을 전송하는 경우 사용하세요. **JSON 템플릿 다운로드**를 클릭합니다. 템플릿에는 이 프로젝트에 사용 가능한 작업자의 현재 사용자 정의 속성 이름과 `Slug`가 포함됩니다.

## 단계별로

1. 가져오기 전에 필요한 워커(**Workers**)와 속성(**Properties**)을 생성합니다.
2. 템플릿을 다운로드하고 텍스트 편집기에서 열고 예제를 자신의 작업으로 바꾸고 파일을 `.json` 확장자로 저장합니다.
3. **가져오기 작업** 탭에서 파일을 선택합니다. 인터페이스에 작업 수가 표시됩니다.
4. **작업 가져오기**를 클릭합니다. 서버는 전체 파일을 확인합니다. 오류가 감지되면 해당 파일의 단일 작업이 기록되지 않습니다. 오류 메시지를 수정하고 다시 시도하십시오.

최소 파일:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

작업자, 속성 및 종속성의 예:

```json
{
  "version": 1,
  "tasks": [
    {
      "key": "design",
      "title": "Prepare the design",
      "description": "## Goal\nPrepare the home page mockup.",
      "state": "features",
      "priority": "high",
      "testing_mode": "hybrid",
      "assignee": { "type": "worker", "worker": "codex" },
      "ai_test_instructions": "Check the build.",
      "human_test_instructions": "Review the page in a browser.",
      "properties": { "Department": "Design" }
    },
    {
      "title": "Approve the design",
      "depends_on": ["design"],
      "assignee": { "type": "human" }
    }
  ]
}
```

`version`는 `1`, 배열 `tasks`여야 합니다(1~1000개 작업). `title`가 필요합니다. 유효한 단계: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. 우선순위: `critical`, `high`, `medium`, `low`; 테스트 모드: `ai`, `human`, `hybrid`. 담당: 기존 `unassigned`(슬러그 또는 ID)가 있는 `human`, `worker` 또는 `worker`. `properties`는 이미 생성된 속성의 이름이나 ID를 사용합니다. `key`는 파일 내에서 고유합니다. `depends_on`는 이러한 키를 나타냅니다. 서버는 작업 ID 자체를 생성합니다.

동일한 파일을 다시 가져오면 새로운 문제가 발생하므로 다시 클릭하기 전에 게시판을 확인하세요. 가져오면 작업이 사람이 만든 것으로 기록됩니다. 기록에 시스템 항목이 나타납니다.
