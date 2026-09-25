import { useState } from "react";
import { Button } from "primereact/button";
import { InputText } from "primereact/inputtext";
import { Message } from "primereact/message";
import { SelectField } from "./SelectField";
import { t } from "./i18n";
import { LanguagePicker } from "./LanguagePicker";
import { builtInColumns, lockedColumns, type ProjectSettings } from "./settings";
import moment from "moment-timezone";

const themes = ["dark", "light", "black", "ubuntu", "windows"] as const;
const timezones = ["local", ...moment.tz.names()];

export function SettingsPage({ value, onSave, saving, error, projectPath }: {
  value: ProjectSettings;
  onSave: (value: ProjectSettings) => void;
  saving: boolean;
  error?: string;
  projectPath: string;
}) {
  const [form, setForm] = useState(value);
  const [newColumn, setNewColumn] = useState("");
  const changeColumn = (from: number, to: number) => {
    if (to < 0 || to >= form.columns.length || lockedColumns.has(form.columns[from].id) || lockedColumns.has(form.columns[to].id)) return;
    const columns = [...form.columns];
    const [column] = columns.splice(from, 1);
    columns.splice(to, 0, column);
    setForm({ ...form, columns });
  };
  const addColumn = () => {
    const name = newColumn.trim();
    if (!name || name.length > 40 || form.columns.length >= 24) return;
    const id = `custom-${crypto.randomUUID().slice(0, 12)}`;
    setForm({ ...form, columns: [...form.columns, { id, name }] });
    setNewColumn("");
  };
  return <div className="settings-page">
    <section className="settings-panel">
      <h2>{t("Appearance and language")}</h2>
      <div className="settings-grid">
        <label>{t("Language")}<LanguagePicker value={form.language} onChange={(language) => setForm({ ...form, language })} /></label>
        <label>{t("Color scheme")}<SelectField value={form.theme} onChange={(e) => setForm({ ...form, theme: e.target.value as ProjectSettings["theme"] })}>
          {themes.map((theme) => <option key={theme} value={theme}>{t(theme[0].toUpperCase() + theme.slice(1))}</option>)}
        </SelectField></label>
        <label>{t("Time zone")}<SelectField filter value={form.timezone} onChange={(e) => setForm({ ...form, timezone: e.target.value })}>
          {timezones.map((zone) => <option key={zone} value={zone}>{zone === "local" ? t("System time zone") : zone}</option>)}
          </SelectField>
        </label>
      </div>
    </section>
    <section className="settings-panel">
      <h2>{t("Board columns")}</h2>
      <p>{t("AI workflow columns are fixed. Custom columns are human-only; moving a task there places it in Backlog for the AI workflow.")}</p>
      <div className="settings-columns">
        {form.columns.map((column, index) => <div className="settings-column" key={column.id}>
          <span>{t(column.name)}</span><small>{lockedColumns.has(column.id) ? t("AI workflow") : column.id.startsWith("custom-") ? t("Custom") : t("Human workflow")}</small>
          {!lockedColumns.has(column.id) && <div className="settings-column-actions">
            <Button icon="pi pi-arrow-up" text aria-label={t("Move up")} disabled={index === 0 || lockedColumns.has(form.columns[index - 1].id)} onClick={() => changeColumn(index, index - 1)} />
            <Button icon="pi pi-arrow-down" text aria-label={t("Move down")} disabled={index === form.columns.length - 1 || lockedColumns.has(form.columns[index + 1].id)} onClick={() => changeColumn(index, index + 1)} />
            {column.id.startsWith("custom-") && <Button icon="pi pi-trash" text severity="danger" aria-label={t("Delete column")} onClick={() => setForm({ ...form, columns: form.columns.filter((c) => c.id !== column.id) })} />}
          </div>}
        </div>)}
      </div>
      <div className="settings-add"><InputText placeholder={t("New column name")} value={newColumn} onChange={(e) => setNewColumn(e.target.value)} onKeyDown={(e) => { if (e.key === "Enter") addColumn(); }} /><Button label={t("Add column")} onClick={addColumn} disabled={!newColumn.trim()} /></div>
      <p>{t("Deleting a custom column returns its tasks to Backlog when you save.")}</p>
    </section>
    <section className="settings-panel">
      <h2>{t("Web and MCP")}</h2>
      <div className="settings-grid">
        <label>{t("Board refresh (seconds)")}<InputText type="number" min={1} max={60} value={form.boardRefreshSeconds.toString()} onChange={(e) => setForm({ ...form, boardRefreshSeconds: Number(e.target.value) })} /></label>
        <label>{t("Worker refresh (seconds)")}<InputText type="number" min={2} max={120} value={form.workerRefreshSeconds.toString()} onChange={(e) => setForm({ ...form, workerRefreshSeconds: Number(e.target.value) })} /></label>
      </div>
      <p>{t("The web server listens on 127.0.0.1:7337 by default. Change the address with --addr when starting agentboard open or serve; restart is required.")}</p>
      <p>{t("Each worker uses its own MCP command. Open a worker to copy its client configuration and check the connection.")}</p>
      <code className="settings-command">agentboard mcp --project "{projectPath}" --worker WORKER_SLUG</code>
    </section>
    {error && <Message severity="error" text={error} />}
    <Button className="primary" label={saving ? t("Saving…") : t("Save settings")} loading={saving} disabled={saving || form.columns.length < builtInColumns.length} onClick={() => onSave(form)} />
  </div>;
}
