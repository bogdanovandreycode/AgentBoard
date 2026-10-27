import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button } from "primereact/button";
import { Dialog } from "primereact/dialog";
import { Message } from "primereact/message";
import { ProgressSpinner } from "primereact/progressspinner";
import { api, withVisibleLoading } from "./api";
import { useNotify } from "./notifications";
import { clearDraftCache, readLocalDrafts, type TaskDraft } from "./taskDrafts";
import { formatDate } from "./settings";
import { t } from "./i18n";

export function TaskDraftsDialog({ projectID, close, open }: { projectID: string; close: () => void; open: (draft: TaskDraft) => void }) {
  const qc = useQueryClient();
  const notify = useNotify();
  const drafts = useQuery({ queryKey: ["task-drafts", projectID], queryFn: () => api<TaskDraft[]>(`/projects/${projectID}/drafts`) });
  const merged = new Map((drafts.data || []).map(draft => [draft.ID, draft]));
  for (const draft of readLocalDrafts(projectID)) if (!merged.has(draft.ID) || !draft.Saved) merged.set(draft.ID, draft);
  const discard = useMutation({
    mutationFn: (draft: TaskDraft) => withVisibleLoading(api(`/projects/${projectID}/drafts/${draft.ID}`, { method: "DELETE" })),
    onSuccess: (_, draft) => { clearDraftCache(projectID, draft.ID); qc.invalidateQueries({ queryKey: ["task-drafts", projectID] }); notify("success", t("Draft deleted")); },
    onError: (error) => notify("error", t("Save failed"), error.message),
  });
  return <Dialog header={t("Task drafts")} visible modal blockScroll onHide={close} style={{ width: "min(760px, calc(100vw - 24px))" }}>
    {drafts.isLoading && <ProgressSpinner />}
    {drafts.error && <Message severity="warn" text={drafts.error.message} />}
    {!drafts.isLoading && !merged.size && <p>{t("No drafts yet")}</p>}
    <div className="task-drafts-list">{[...merged.values()].sort((a,b) => b.UpdatedAt.localeCompare(a.UpdatedAt)).map(draft => <article key={draft.ID}>
      <div><strong>{draft.Input.Title || t("Untitled draft")}</strong><small>{draft.Saved ? t("Saved draft") : t("Local recovery copy")} · {formatDate(draft.UpdatedAt)}</small></div>
      <Button label={t("Resume draft")} onClick={() => open(draft)} /><Button icon="pi pi-trash" severity="danger" text aria-label={`${t("Delete draft")} ${draft.Input.Title || t("Untitled draft")}`} loading={discard.isPending && discard.variables?.ID === draft.ID} disabled={discard.isPending} onClick={() => discard.mutate(draft)} />
    </article>)}</div>
  </Dialog>;
}
