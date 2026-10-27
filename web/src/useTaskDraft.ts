import { useCallback, useEffect, useRef, useState, type Dispatch, type SetStateAction } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { api } from "./api";
import { useGlobalSettings } from "./useGlobalSettings";
import { activeDraftKey, cacheDraft, clearDraftCache, draftForm, draftInput, emptyTaskForm, meaningfulForm, missingDraftFields, type TaskDraft, type TaskFormData } from "./taskDrafts";

export function useTaskDraft(projectID: string, initial?: TaskDraft) {
  const settings = useGlobalSettings();
  const qc = useQueryClient();
  const [form, updateForm] = useState<TaskFormData>(() => initial ? draftForm(initial.Input) : emptyTaskForm());
  const formRef = useRef(form);
  const initialSavedInput = initial?.LastSavedInput || (initial?.Saved ? JSON.stringify(draftInput(form)) : "");
  const meta = useRef({ id: initial?.ID || crypto.randomUUID(), revision: initial?.Revision || 0, lastSaved: initialSavedInput, updatedAt: initial?.UpdatedAt || "", retired: false });
  const createdAt = initial?.CreatedAt || "";
  const [status, setStatus] = useState(() => !meaningfulForm(form) ? "empty" : JSON.stringify(draftInput(form)) === initialSavedInput ? "saved" : "dirty");
  const [error, setError] = useState("");
  const [storageError, setStorageError] = useState(false);
  const mounted = useRef(true);
  const queue = useRef<Promise<unknown>>(Promise.resolve());
  const persistLocal = useCallback((active = false) => {
    if (meta.current.retired) return true;
    const current = formRef.current;
    if (!meaningfulForm(current)) { clearDraftCache(projectID, meta.current.id); return true; }
    try {
      const input = draftInput(current);
      cacheDraft({ ID: meta.current.id, ProjectID: projectID, Input: input, Revision: meta.current.revision, Saved: meta.current.lastSaved === JSON.stringify(input), Ready: missingDraftFields(current).length === 0, MissingFields: missingDraftFields(current), LastSavedInput: meta.current.lastSaved, CreatedAt: createdAt, UpdatedAt: new Date().toISOString() });
      if (active) localStorage.setItem(activeDraftKey(projectID), meta.current.id);
      return true;
    } catch { return false; }
  }, [projectID, createdAt]);
  const setForm: Dispatch<SetStateAction<TaskFormData>> = useCallback(value => {
    const next = typeof value === "function" ? value(formRef.current) : value;
    formRef.current = next;
    updateForm(next);
    setStorageError(!persistLocal(true));
    setStatus(current => current === "saving" ? current : !meaningfulForm(next) ? "empty" : JSON.stringify(draftInput(next)) === meta.current.lastSaved ? "saved" : "dirty");
  }, [persistLocal]);
  const flush = useCallback((): Promise<TaskDraft | undefined> => {
    const request = queue.current.catch(() => {}).then(async () => {
      if (meta.current.retired) return undefined;
      const input = draftInput(formRef.current);
      const serialized = JSON.stringify(input);
      if ((!meaningfulForm(formRef.current) && meta.current.revision === 0) || serialized === meta.current.lastSaved) return undefined;
      persistLocal();
      if (mounted.current) { setStatus("saving"); setError(""); }
      try {
        let saved: TaskDraft;
        try {
          saved = await api<TaskDraft>(`/projects/${projectID}/drafts/${meta.current.id}`, { method: "PUT", body: JSON.stringify({ Input: input, Revision: meta.current.revision }), keepalive: new Blob([serialized]).size < 60000, signal: AbortSignal.timeout(15000) });
        } catch (caught) {
          if (!(caught instanceof Error) || !caught.message.includes("DRAFT_CONFLICT")) throw caught;
          // A previous save may have reached SQLite while its response was lost.
          // Reconcile an identical snapshot; never overwrite a different revision.
          const remote = (await api<TaskDraft[]>(`/projects/${projectID}/drafts`)).find(draft => draft.ID === meta.current.id);
          if (!remote || JSON.stringify(draftInput(draftForm(remote.Input))) !== serialized) throw caught;
          saved = remote;
        }
        meta.current.revision = saved.Revision;
        meta.current.lastSaved = serialized;
        meta.current.updatedAt = saved.UpdatedAt;
        if (!saved.Saved) {
          clearDraftCache(projectID, meta.current.id);
          if (saved.Revision > 0) meta.current = { id: crypto.randomUUID(), revision: 0, lastSaved: "", updatedAt: "", retired: false };
        } else { persistLocal(); }
        if (mounted.current) setStatus(saved.Saved ? (JSON.stringify(draftInput(formRef.current)) === serialized ? "saved" : "dirty") : "empty");
        qc.invalidateQueries({ queryKey: ["task-drafts", projectID] });
        return saved;
      } catch (caught) {
        const message = caught instanceof Error ? caught.message : String(caught);
        if (mounted.current) { setStatus("error"); setError(message); }
        persistLocal();
        throw caught;
      }
    });
    queue.current = request;
    return request;
  }, [persistLocal, projectID, qc]);
  useEffect(() => {
    if (!settings.data?.draftAutosaveEnabled) return;
    const timer = window.setInterval(() => { void flush().catch(() => {}); }, settings.data.draftAutosaveSeconds * 1000);
    return () => window.clearInterval(timer);
  }, [settings.data?.draftAutosaveEnabled, settings.data?.draftAutosaveSeconds, flush]);
  useEffect(() => {
    mounted.current = true;
    persistLocal(true);
    const pageHide = () => { persistLocal(true); void flush().catch(() => {}); };
    window.addEventListener("pagehide", pageHide);
    return () => {
      mounted.current = false;
      window.removeEventListener("pagehide", pageHide);
      persistLocal();
      void flush().catch(() => {});
    };
  }, [flush, persistLocal]);
  const deactivate = () => {
    persistLocal();
    try { localStorage.removeItem(activeDraftKey(projectID)); } catch { /* Recovery may be unavailable. */ }
  };
  const retire = () => { meta.current.retired = true; clearDraftCache(projectID, meta.current.id); qc.invalidateQueries({ queryKey: ["task-drafts", projectID] }); };
  const publish = async () => {
    await flush();
    const task = await api<{ ID: string }>(`/projects/${projectID}/drafts/${meta.current.id}/publish`, { method: "POST", body: JSON.stringify({ Revision: meta.current.revision }) });
    retire();
    return task;
  };
  const recoverAsNew = () => {
    clearDraftCache(projectID, meta.current.id);
    meta.current = { id: crypto.randomUUID(), revision: 0, lastSaved: "", updatedAt: "", retired: false };
    setError(""); setStatus("dirty"); persistLocal(true);
    return flush();
  };
  return { form, setForm, status, error, storageError, flush, publish, deactivate, recoverAsNew, meaningful: meaningfulForm(form), missingFields: missingDraftFields(form), settingsError: settings.error?.message, onBlur: () => { if (settings.data?.draftSaveOnBlur) void flush().catch(() => {}); } };
}
