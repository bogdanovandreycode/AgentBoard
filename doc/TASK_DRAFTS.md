# Project management and task drafts

The **All projects** page shows the registered folder, published task count and
latest task activity. Activity includes changes to tasks, history, and usage
records. Use the pencil button to rename a project or select an existing folder
with the native folder dialog. A path change updates the registration and marker;
it preserves the project's ID, workers, tasks and settings, and does not move files.
Folders already registered to another project cannot be selected.

## Global settings

Open **Global settings** from the project launcher or from the **Global settings**
tab inside any board's settings. These preferences are shared by all projects in
the same AgentBoard database:

- Timed draft autosave is enabled, every **10 seconds**, by default.
- **Save when leaving a field** is enabled by default and can be changed separately.
- **MCP project context** checks can be enabled or disabled globally. See
  [Workers and MCP](WORKERS_MCP.md) for their scope and the individual CLI override.

## Draft lifecycle

The new-task dialog saves incomplete work as a draft. Its status line distinguishes
unsaved changes, an active save, a saved draft and a save error. A draft stays off
the task board and is unavailable to MCP workers until it is published.

**Title, State, Responsible and Description** must be filled before **Create task**
publishes it. Before then, **Save draft** retains it without creating a task.
The default empty form is discarded when closed. Clearing an existing draft back
to empty removes it at the next save. Closing a nonempty dialog also saves it.

Use **Drafts** on the board to resume or delete saved drafts. Multiple drafts are
supported. A reload restores the draft that was open. Drafts are stored in SQLite;
a browser-local recovery copy is updated on each change so unsaved edits can also
survive a reload or a temporary server failure. If browser storage is unavailable,
the dialog warns that the server draft must be saved before closing.

Draft writes use revisions to prevent one window from overwriting another's work.
A lost save response is reconciled when the server contains the identical input.
If the content differs, the conflict is shown and **Recover as a new draft** keeps
the local content separately. Publishing is atomic and idempotent: retrying the
same draft after a lost response returns the same task rather than creating a copy.
