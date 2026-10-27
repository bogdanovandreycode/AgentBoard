export type TaskFormData = {
  Title: string; Description: string; State: string; Priority: string; TestingMode: string;
  AssigneeType: string; AssigneeWorkerID: string; AITestInstructions: string; HumanTestInstructions: string;
  DependencyIDs: string[]; Properties: Record<string, string>;
};
export type DraftInput = Omit<TaskFormData, "AssigneeType" | "AssigneeWorkerID"> & { BoardColumn: string; Assignee: { type: string; worker_id?: string } };
export type TaskDraft = {
  ID: string; ProjectID: string; Input: DraftInput; Revision: number; Saved: boolean; Ready: boolean;
  MissingFields: string[]; CreatedAt: string; UpdatedAt: string; LastSavedInput?: string;
};
export const emptyTaskForm = (): TaskFormData => ({ Title: "", Description: "", State: "", Priority: "medium", TestingMode: "ai", AssigneeType: "unassigned", AssigneeWorkerID: "", AITestInstructions: "", HumanTestInstructions: "", DependencyIDs: [], Properties: {} });
export function draftInput(form: TaskFormData): DraftInput {
  const { AssigneeType, AssigneeWorkerID, ...fields } = form;
  return { ...fields, State: form.State.startsWith("custom-") ? "backlog" : form.State, BoardColumn: form.State.startsWith("custom-") ? form.State : "", Assignee: { type: AssigneeType, ...(AssigneeType === "worker" ? { worker_id: AssigneeWorkerID } : {}) } };
}
export function draftForm(input: DraftInput): TaskFormData {
  const { Assignee, BoardColumn, ...fields } = input;
  const form = { ...emptyTaskForm(), ...fields, State: BoardColumn || fields.State || "", AssigneeType: Assignee?.type || "unassigned", AssigneeWorkerID: Assignee?.worker_id || "", DependencyIDs: fields.DependencyIDs || [], Properties: fields.Properties || {} };
  return Object.fromEntries(Object.keys(emptyTaskForm()).map(key => [key, form[key as keyof TaskFormData]])) as TaskFormData;
}
export function meaningfulForm(form: TaskFormData) {
  return Boolean((form.Title + form.Description + form.AITestInstructions + form.HumanTestInstructions).trim() || form.State || form.Priority !== "medium" || form.TestingMode !== "ai" || form.AssigneeType !== "unassigned" || form.DependencyIDs.length || Object.values(form.Properties).some(value => value.trim()));
}
export function missingDraftFields(form: TaskFormData) {
  return [!form.Title.trim() && "Title", !form.State && "State", (!form.AssigneeType || form.AssigneeType === "unassigned" || (form.AssigneeType === "worker" && !form.AssigneeWorkerID)) && "Responsible", !form.Description.trim() && "Description"].filter(Boolean) as string[];
}
const cachePrefix = (projectID: string) => `agentboard-draft:${projectID}:`;
export const activeDraftKey = (projectID: string) => `agentboard-active-draft:${projectID}`;
export function readLocalDrafts(projectID: string): TaskDraft[] {
  const result: TaskDraft[] = [];
  try {
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (!key?.startsWith(cachePrefix(projectID))) continue;
      try {
        const value = JSON.parse(localStorage.getItem(key) || "null") as TaskDraft;
        if (value?.ProjectID === projectID && value.ID && value.Input && meaningfulForm(draftForm(value.Input))) result.push(value);
      } catch { /* Ignore malformed recovery entries. */ }
    }
  } catch { /* The server draft remains available when browser storage is blocked. */ }
  return result;
}
export function readActiveDraft(projectID: string) {
  try { const id = localStorage.getItem(activeDraftKey(projectID)); return readLocalDrafts(projectID).find(draft => draft.ID === id); } catch { return undefined; }
}
export function cacheDraft(draft: TaskDraft) { localStorage.setItem(cachePrefix(draft.ProjectID) + draft.ID, JSON.stringify(draft)); }
export function clearDraftCache(projectID: string, id: string) {
  try {
    localStorage.removeItem(cachePrefix(projectID) + id);
    if (localStorage.getItem(activeDraftKey(projectID)) === id) localStorage.removeItem(activeDraftKey(projectID));
  } catch { /* Server-side removal remains authoritative. */ }
}
