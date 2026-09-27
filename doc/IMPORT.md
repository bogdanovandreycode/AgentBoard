# Import tasks from JSON

The **Import tasks** tab accepts UTF-8 JSON files. To move tasks from another service, click **Download JSON template**. The template includes the project's custom property names and available worker slugs.

## Step by step

1. Create the necessary workers (**Workers**) and properties (**Properties**) before importing.
2. Download the template, open it in a text editor, replace the examples with your tasks and save the file with the `.json` extension.
3. On the **Import tasks** tab, select a file. The interface will show the number of tasks.
4. Click **Import tasks**. The server will check the entire file: if an error is detected, not a single task from it will be written. Correct the error message and try again.

Minimum file:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Example with worker, property and dependency:

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

`version` should be `1`, the `tasks` array should be from 1 to 1000 tasks. `title` is required. Valid stages: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Priorities: `critical`, `high`, `medium`, `low`; verification modes: `ai`, `human`, `hybrid`. Responsible: `unassigned`, `human` or `worker` with an existing `worker` (Slug or ID). `properties` uses the names or IDs of already created properties. `key` is unique within the file; `depends_on` refers to such keys. The server creates task IDs itself.

Importing the same file again creates new tasks, so check the board before repeating an import. Imported tasks are recorded as human-created, with a system entry in their history.
