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

## Screenshots

Use **Attach images** in a description or comment, drop image files into the editor, or paste a screenshot from the clipboard. Multiple images are supported. The editor inserts numbered Markdown links such as `[🖼 Screenshot 1](/api/images/IMAGE_ID)` at the cursor and shows small attachments below the text. You can repeat a link when referring to the same screenshot in several places. Removing an attachment from the editor removes its references from that text.

Click a link or thumbnail to open the image viewer. The viewer can download the WebP file or show it in the local file manager. Saving waits until uploads finish, and a failed comment submission keeps your text and attachments.

PNG, JPEG, WebP and GIF images are converted to lossless WebP; animated GIFs use the first frame. Uploads are limited to 20 MiB and 20 million pixels. Files live in `images` beside the AgentBoard database; back up this folder together with the database. Image content is available to AI through `get_task` and `get_image`, and agents can use `upload_image` followed by `add_history` to attach their screenshots. The same task visibility rules apply to images.

## File attachments

The editor has separate **Images** and **Files** buttons with image and document icons. **Files** accepts any file type, up to 50 MiB per file, including archives, documents, source files and empty files. It inserts a document Markdown link with the filename and lists the attachment below the description or comment. Click its link or attachment to download the original file or show it in the local file manager. Dragging files or pasting attachments also works; supported raster images use WebP, while other formats remain unchanged.

File bytes are stored in the `files` folder beside the selected database; include that folder in backups. AI workers attach files with `upload_file` and its returned Markdown link in `add_history`. `get_task` exposes attachment metadata in `Files` and comment `Files`; `get_file` reads bounded UTF-8 or base64 chunks using `offset`, `length`, `next_offset` and `eof`. Task permissions apply to every upload and read.

## Columns and properties

In the **Settings** tab you can change the order of columns available to a person and add your own. The four AI stages are fixed and proceed in the same order. The user column is a place for tasks postponed by a person: for AI, such a task has the `Backlog` state. When a column is deleted, its tasks return to the normal `Backlog`.

In the **Properties** tab you can add fields such as text, number, flag, date, selection and URL. `Human only` hides the field from AI; `Agent read` allows reading, `Agent read/write` also allows writing through supported tools. This does not change the AI's rights to the task steps.
