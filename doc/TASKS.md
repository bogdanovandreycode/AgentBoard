# Working with tasks

## Add task

On the **Board** tab, click **New task**. Fill in the title and description. The description and testing instructions support Markdown: highlight text and use the formatting bar for bold, italic, link, and list. Click **Save**.

Fields:

| Field | Purpose |
| --- | --- |
| Title | Short name of the task. |
| Description | What needs to be done and how to understand that the work is ready. |
| State | Work stage; a new task usually starts in `Backlog`. |
| Priority | `critical`, `high`, `medium` or `low`. |
| Responsible | A human, a specific worker or unassigned. |
| Testing mode | `AI` - AI testing; `Human` - human testing; `Hybrid` - both. |
| Dependencies | Tasks that should be completed earlier. |
| AI/Human test instructions | Instructions to the appropriate reviewer. |
| Custom properties | Additional fields created in the **Properties** tab. |

For an AI task, first create a worker, select it in Responsible, then transfer the card to `Features`. AI sees only the tasks assigned to it at four stages from `Features` to `Verification`.

## Move task

Drag the card between the columns. You can switch between an all-columns view and wide columns with horizontal scrolling. `Backlog` and `Complete` are human managed. AI can only advance a task `Features → In progress → Testing → Verification` through special MCP tools. A task with uncompleted manual testing should not pass AI verification.

## Task card

Click a card to see the description, assignee, test instructions, history, tests, artifacts and AI costs. **Edit task** changes its content. Human comments are added to the shared history. The search field matches titles, descriptions, IDs and workers; a separate button opens expanded search.

## Columns and properties

In the **Settings** tab you can change the order of columns available to a person and add your own. The four AI stages are fixed and proceed in the same order. The user column is a place for tasks postponed by a person: for AI, such a task has the `Backlog` state. When a column is deleted, its tasks return to the normal `Backlog`.

In the **Properties** tab you can add fields such as text, number, flag, date, selection and URL. `Human only` hides the field from AI; `Agent read` allows reading, `Agent read/write` also allows writing through supported tools. This does not change the AI's rights to the task steps.
