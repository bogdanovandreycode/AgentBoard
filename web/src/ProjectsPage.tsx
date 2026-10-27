import { useState } from "react";
import { Button } from "primereact/button";
import { InputText } from "primereact/inputtext";
import { Paginator } from "primereact/paginator";
import { Message } from "primereact/message";
import { Dialog } from "primereact/dialog";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, withVisibleLoading } from "./api";
import { useNotify } from "./notifications";
import { GlobalSettingsPanel } from "./GlobalSettings";
import { relativeTime } from "./workerPresentation";
import { formatDate } from "./settings";
import { t } from "./i18n";

export type Project = { ID: string; Name: string; Path: string; TaskCount: number; LastActivityAt: string };
const pageSize = 12;

export function ProjectsPage({ projects, invalidProject, onOpen }: {
  projects: Project[];
  invalidProject: boolean;
  onOpen: (id: string) => void;
}) {
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(0);
  const [editing, setEditing] = useState<Project>();
  const [globalOpen, setGlobalOpen] = useState(false);
  const filtered = projects.filter((project) =>
    `${project.Name} ${project.Path}`.toLocaleLowerCase().includes(search.trim().toLocaleLowerCase()),
  ).sort((a, b) => a.Name.localeCompare(b.Name, undefined, { numeric: true }));
  const first = Math.min(page * pageSize, Math.max(0, Math.ceil(filtered.length / pageSize) - 1) * pageSize);
  const visible = filtered.slice(first, first + pageSize);

  return <main className="project-launcher">
    <header className="launcher-header">
      <div className="launcher-brand"><img src="/agentboard-icon.png" alt="" /><span>AgentBoard</span></div>
      <div className="launcher-actions"><span className="launcher-count">{t("Projects")}: {projects.length}</span><Button text icon="pi pi-cog" label={t("Global settings")} onClick={() => setGlobalOpen(true)} /></div>
    </header>
    <section className="launcher-intro">
      <p className="launcher-eyebrow">{t("Workspace")}</p>
      <h1>{t("Projects")}</h1>
      <p>{t("Choose a project to open its board and workers.")}</p>
    </section>
    {invalidProject && <Message severity="warn" text={t("The requested project is not registered on this server. Select a project below.")} />}
    {projects.length === 0 ? <div className="launcher-empty">
      <img src="/agentboard-icon.png" alt="" />
      <h2>{t("No projects connected yet")}</h2>
      <p>{t("Open a terminal in your project folder and run:")}</p>
      <code>agentboard init</code>
      <p>{t("Then run agentboard open. The project will appear here automatically.")}</p>
    </div> : <>
      <div className="launcher-toolbar">
        <div className="launcher-search"><i className="pi pi-search" aria-hidden="true" /><InputText value={search} onChange={(event) => { setSearch(event.target.value); setPage(0); }} placeholder={t("Search projects by name or path…")} aria-label={t("Search projects")} /></div>
        <span>{t("Projects")}: {filtered.length}</span>
      </div>
      {visible.length ? <div className="project-grid">{visible.map((project) => <article className="project-tile" key={project.ID}>
        <div className="project-tile-top"><span className="project-avatar">{project.Name.trim().charAt(0).toLocaleUpperCase() || "P"}</span><Button icon="pi pi-pencil" text aria-label={`${t("Edit project")} ${project.Name}`} onClick={() => setEditing(project)} /></div>
        <h2>{project.Name}</h2>
        <p className="project-path" title={project.Path}>{project.Path}</p>
        <div className="project-statistics"><span><i className="pi pi-list" /> {t("Tasks")}: <strong>{project.TaskCount ?? 0}</strong></span><span title={project.LastActivityAt ? formatDate(project.LastActivityAt) : undefined}><i className="pi pi-clock" /> {t("Last task activity")}: {project.TaskCount ? relativeTime(project.LastActivityAt) : t("No task activity yet")}</span></div>
        <Button label={t("Open project")} icon="pi pi-arrow-right" iconPos="right" text onClick={() => onOpen(project.ID)} />
      </article>)}</div> : <div className="launcher-no-results"><i className="pi pi-search" aria-hidden="true" /><h2>{t("No matching projects")}</h2><p>{t("Try another name or path.")}</p><Button label={t("Clear search")} text onClick={() => setSearch("")} /></div>}
      {filtered.length > pageSize && <Paginator className="launcher-paginator" first={first} rows={pageSize} totalRecords={filtered.length} onPageChange={(event) => setPage(event.page)} />}
    </>}
    {editing && <ProjectEditor project={editing} close={() => setEditing(undefined)} />}
    <Dialog header={t("Global settings")} visible={globalOpen} onHide={() => setGlobalOpen(false)} modal blockScroll style={{ width: "min(850px, calc(100vw - 24px))" }}><GlobalSettingsPanel /></Dialog>
  </main>;
}

function ProjectEditor({ project, close }: { project: Project; close: () => void }) {
  const [name, setName] = useState(project.Name);
  const [path, setPath] = useState(project.Path);
  const qc = useQueryClient();
  const notify = useNotify();
  const pick = useMutation({
    mutationFn: () => api<{ path: string }>("/folders/pick", { method: "POST", body: JSON.stringify({ Path: path }) }),
    onSuccess: (result) => { if (result.path) setPath(result.path); },
    onError: (error) => notify("error", t("Folder selection failed"), error.message),
  });
  const save = useMutation({
    mutationFn: () => withVisibleLoading(api<Project>(`/projects/${project.ID}`, { method: "PATCH", body: JSON.stringify({ Name: name, Path: path }) })),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["projects"] }); notify("success", t("Project saved")); close(); },
    onError: (error) => notify("error", t("Save failed"), error.message),
  });
  return <Dialog header={t("Edit project")} visible modal blockScroll onHide={close} style={{ width: "min(620px, calc(100vw - 24px))" }}>
    <div className="form-grid">
      <label className="wide" htmlFor="project-name">{t("Project name")}<InputText id="project-name" value={name} maxLength={120} onChange={(event) => setName(event.target.value)} autoFocus /></label>
      <label className="wide" htmlFor="project-path">{t("Project folder")}<div className="project-folder-field"><InputText id="project-path" value={path} onChange={(event) => setPath(event.target.value)} /><Button icon="pi pi-folder-open" aria-label={t("Choose project folder")} loading={pick.isPending} disabled={save.isPending} onClick={() => pick.mutate()} /></div></label>
    </div>
    <p>{t("Changing the folder updates the registered project path. Your files are not moved.")}</p>
    {(save.error || pick.error) && <Message severity="error" text={(save.error || pick.error)?.message} />}
    <div className="project-editor-actions"><Button label={t("Cancel")} text onClick={close} disabled={save.isPending} /><Button label={t("Save")} loading={save.isPending} disabled={save.isPending || pick.isPending || !name.trim() || !path.trim()} onClick={() => save.mutate()} /></div>
  </Dialog>;
}
