import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "primereact/button";
import { InputText } from "primereact/inputtext";
import { InputSwitch } from "primereact/inputswitch";
import { Message } from "primereact/message";
import { ProgressSpinner } from "primereact/progressspinner";
import { api, withVisibleLoading } from "./api";
import { useNotify } from "./notifications";
import { t } from "./i18n";
import { useGlobalSettings, type GlobalSettings } from "./useGlobalSettings";

export function GlobalSettingsPanel() {
  const settings = useGlobalSettings();
  if (settings.error) return <Message severity="error" text={settings.error.message} />;
  return settings.data ? <GlobalSettingsForm value={settings.data} key={JSON.stringify(settings.data)} /> : <ProgressSpinner />;
}

function GlobalSettingsForm({ value }: { value: GlobalSettings }) {
  const [form, setForm] = useState(value);
  const qc = useQueryClient();
  const notify = useNotify();
  const save = useMutation({
    mutationFn: () => withVisibleLoading(api<GlobalSettings>("/settings/global", { method: "PUT", body: JSON.stringify(form) })),
    onSuccess: (saved) => { qc.setQueryData(["global-settings"], saved); notify("success", t("Settings saved")); },
    onError: (error) => notify("error", t("Save failed"), error.message),
  });
  return <section className="settings-panel global-settings-panel">
    <h2>{t("Global settings")}</h2>
    <p>{t("These settings apply to all project boards and MCP workers using this AgentBoard database.")}</p>
    <h3>{t("Task drafts and autosave")}</h3>
    <div className="settings-grid">
      <label htmlFor="draft-autosave">{t("Timed autosave")}<InputSwitch inputId="draft-autosave" checked={form.draftAutosaveEnabled} onChange={(e) => setForm({ ...form, draftAutosaveEnabled: !!e.value })} /></label>
      <label htmlFor="draft-interval">{t("Autosave interval (seconds)")}<InputText id="draft-interval" type="number" min={1} max={3600} value={String(form.draftAutosaveSeconds)} onChange={(e) => setForm({ ...form, draftAutosaveSeconds: Number(e.target.value) })} /></label>
      <label htmlFor="draft-blur">{t("Save when leaving a field")}<InputSwitch inputId="draft-blur" checked={form.draftSaveOnBlur} onChange={(e) => setForm({ ...form, draftSaveOnBlur: !!e.value })} /></label>
    </div>
    <h3>{t("MCP project context")}</h3>
    <label className="global-toggle" htmlFor="mcp-context-check"><InputSwitch inputId="mcp-context-check" checked={form.checkProjectContext} onChange={(e) => setForm({ ...form, checkProjectContext: !!e.value })} />{t("Check the worker project against its launch directory and client workspace")}</label>
    <p>{t("When enabled, a mismatch asks for confirmation; clients without question support receive an error. When disabled, workers use their configured project without this check. Task workflow permissions always remain enforced. Changes apply to existing MCP sessions on their next tool call.")}</p>
    {save.error && <Message severity="error" text={save.error.message} />}
    <Button className="primary" label={t("Save global settings")} loading={save.isPending} disabled={save.isPending || form.draftAutosaveSeconds < 1 || form.draftAutosaveSeconds > 3600} onClick={() => save.mutate()} />
  </section>;
}
