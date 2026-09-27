import { useState } from "react";
import { Button } from "primereact/button";
import { InputText } from "primereact/inputtext";
import { Paginator } from "primereact/paginator";
import { Message } from "primereact/message";
import { t } from "./i18n";

type Project = { ID: string; Name: string; Path: string };
const pageSize = 12;

export function ProjectsPage({ projects, invalidProject, onOpen }: {
  projects: Project[];
  invalidProject: boolean;
  onOpen: (id: string) => void;
}) {
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(0);
  const filtered = projects.filter((project) =>
    `${project.Name} ${project.Path}`.toLocaleLowerCase().includes(search.trim().toLocaleLowerCase()),
  ).sort((a, b) => a.Name.localeCompare(b.Name, undefined, { numeric: true }));
  const first = Math.min(page * pageSize, Math.max(0, Math.ceil(filtered.length / pageSize) - 1) * pageSize);
  const visible = filtered.slice(first, first + pageSize);

  return <main className="project-launcher">
    <header className="launcher-header">
      <div className="launcher-brand"><img src="/agentboard-icon.png" alt="" /><span>AgentBoard</span></div>
      <span className="launcher-count">{t("Projects")}: {projects.length}</span>
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
        <div className="project-tile-top"><span className="project-avatar">{project.Name.trim().charAt(0).toLocaleUpperCase() || "P"}</span><i className="pi pi-arrow-up-right" aria-hidden="true" /></div>
        <h2>{project.Name}</h2>
        <p title={project.Path}>{project.Path}</p>
        <Button label={t("Open project")} icon="pi pi-arrow-right" iconPos="right" text onClick={() => onOpen(project.ID)} />
      </article>)}</div> : <div className="launcher-no-results"><i className="pi pi-search" aria-hidden="true" /><h2>{t("No matching projects")}</h2><p>{t("Try another name or path.")}</p><Button label={t("Clear search")} text onClick={() => setSearch("")} /></div>}
      {filtered.length > pageSize && <Paginator className="launcher-paginator" first={first} rows={pageSize} totalRecords={filtered.length} onPageChange={(event) => setPage(event.page)} />}
    </>}
  </main>;
}
