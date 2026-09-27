# Settings

Open **Settings** in the side menu of the selected project. Changes are saved with the **Save settings** button and stored for the project in the AgentBoard database.

## Language and appearance

The **Language** field is a searchable menu. Its default, **Follow system**, uses the browser's language. The interface supports the languages listed on the [documentation languages](LANGUAGES.md) page. Some translations are automatic and may need review by native speakers. Client names, commands, JSON fields and user data remain unchanged.

**Color scheme**: Dark (original theme), Light, Black, Ubuntu and Windows. **Time zone** controls the display of dates; the data continues to be stored in UTC. **System time zone** uses computer settings.

## Columns

The four stages `Features`, `In progress`, `Testing`, `Verification` are fixed in this order: they cannot be renamed or deleted. Other columns can be rearranged using the available arrows. Enter a name and click **Add column** to create a custom column. It is intended for human tasks: AI sees such a task as `Backlog` and does not receive it through MCP. Removing a column while saving transfers its tasks to the regular `Backlog`.

## Web and MCP

**Board refresh** sets the refresh rate of the board and cards in seconds (1–60). **Worker refresh** updates the status of workers (2–120 seconds). This is web interface polling, not AI trigger frequency. Workers do not start automatically.

The web server address is set when starting the CLI, for example `agentboard open --addr 127.0.0.1:7444`. To change the address, the server must be restarted. The default is `127.0.0.1:7337`. MCP works through a separate local command `agentboard mcp --project ... --worker ...` and is independent of the web port. For a non-standard database, specify the same `--db` in all commands. Copy the configuration of each client from the worker card.
