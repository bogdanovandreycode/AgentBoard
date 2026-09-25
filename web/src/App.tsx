import { useEffect, useState } from "react";
import {
  DndContext,
  useDraggable,
  useDroppable,
  type DragEndEvent,
} from "@dnd-kit/core";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button } from "primereact/button";
import { Dialog } from "primereact/dialog";
import { Sidebar } from "primereact/sidebar";
import { InputText } from "primereact/inputtext";
import { InputTextarea } from "primereact/inputtextarea";
import { Tag } from "primereact/tag";
import { DataTable } from "primereact/datatable";
import { Column as DataColumn } from "primereact/column";
import { TabMenu } from "primereact/tabmenu";
import { InputSwitch } from "primereact/inputswitch";
import { Card } from "primereact/card";
import { Panel } from "primereact/panel";
import { Timeline } from "primereact/timeline";
import { Toolbar } from "primereact/toolbar";
import { Avatar } from "primereact/avatar";
import { Badge } from "primereact/badge";
import { Divider } from "primereact/divider";
import { Message } from "primereact/message";
import { ProgressSpinner } from "primereact/progressspinner";
import { ConfirmDialog } from "primereact/confirmdialog";
import {
  Bot,
  Columns3,
  Filter,
  GripVertical,
  Plus,
  Settings2,
  Sparkles,
  User,
  Users,
  X,
} from "lucide-react";
import { WorkerDiagnostics } from "./WorkerDiagnostics";
import { SelectField } from "./SelectField";
import { capabilityPresets, relativeTime } from "./workerPresentation";
import "./App.css";
import "./Properties.css";
import "./PrimeLayout.css";

type Project = { ID: string; Name: string; Path: string };
type Worker = {
  ID: string;
  Name: string;
  Slug: string;
  Description: string;
  Kind: string;
  Capabilities: string;
  Enabled: boolean;
  Archived: boolean;
  AssignedTaskCount: number;
  ActiveSessionCount: number;
  SessionCount: number;
  MCPCalls: number;
  LastActivityAt?: string;
};
type PropertyDef = {
  ID: string;
  Name: string;
  Type: string;
  Options: string;
  Visibility: string;
};
type Artifact = {
  ID: string;
  Name: string;
  Kind: string;
  Path?: string;
  URL?: string;
  Description: string;
  CreatedByType: string;
  CreatorName: string;
  CreatedAt: string;
};
type History = {
  ID: string;
  ActorType: string;
  ActorName: string;
  EntryType: string;
  Content: string;
  CreatedAt: string;
};
type TestRun = {
  ID: string;
  ActorName: string;
  Runner: string;
  Type: string;
  Command?: string;
  Status: string;
  Summary: string;
  OutputExcerpt: string;
  DurationMS?: number;
  CreatedAt: string;
};
type Usage = {
  ID: string;
  WorkerName: string;
  Provider?: string;
  Model?: string;
  InputTokens?: number;
  OutputTokens?: number;
  ReasoningTokens?: number;
  CachedTokens?: number;
  ModelCalls?: number;
  ToolCalls?: number;
  MCPCalls?: number;
  WallTimeMS?: number;
};
type Dependency = {
  DependsOnTaskID: string;
  Title: string;
  State: string;
  Artifacts: Artifact[];
};
type Task = {
  ID: string;
  ProjectID: string;
  Title: string;
  Description: string;
  State: string;
  Position: number;
  Priority: string;
  AssigneeType: string;
  AssigneeWorkerID?: string;
  AssigneeName: string;
  CreatedByType: string;
  CreatorName: string;
  CreatedAt: string;
  UpdatedAt: string;
  SourceTaskID?: string;
  TestingMode: string;
  AITestInstructions: string;
  HumanTestInstructions: string;
  Dependencies: Dependency[];
};
type Details = Task & {
  History: History[];
  TestRuns: TestRun[];
  Artifacts: Artifact[];
  Usage: Usage[];
  SpawnedTasks: Task[];
  Properties: { PropertyDefinitionID: string; Name: string; Value: string }[];
};
type Board = Record<string, Task[]>;
const states = [
  ["backlog", "Backlog"],
  ["features", "Features"],
  ["in_progress", "In progress"],
  ["testing", "Testing"],
  ["verification", "Verification"],
  ["complete", "Complete"],
] as const;
async function api<T>(url: string, options?: RequestInit): Promise<T> {
  const r = await fetch("/api" + url, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!r.ok) {
    const e = await r.json().catch(() => ({ error: r.statusText }));
    throw new Error(e.error || r.statusText);
  }
  return r.status === 204 ? (undefined as T) : r.json();
}
const short = (id: string) => id.slice(0, 8).toUpperCase();
const date = (v: string) => new Date(v).toLocaleString();

function App() {
  const [, tick] = useState(0);
  useEffect(() => {
    const timer = window.setInterval(() => tick((value) => value + 1), 1000);
    return () => window.clearInterval(timer);
  }, []);
  const qc = useQueryClient(),
    [tab, setTab] = useState<"board" | "workers" | "properties">("board"),
    [projectID, setProjectID] = useState(
      () => new URLSearchParams(location.search).get("project") || "",
    ),
    [selected, setSelected] = useState<string>(),
    [create, setCreate] = useState(false),
    [filters, setFilters] = useState({
      responsible: "all",
      priority: "all",
      origin: "all",
      testing: "all",
    });
  const projects = useQuery({
      queryKey: ["projects"],
      queryFn: () => api<Project[]>("/projects"),
      refetchInterval: 10000,
    }),
    active = projects.data?.some((p) => p.ID === projectID) ? projectID : "";
  const workers = useQuery({
      queryKey: ["workers", active],
      queryFn: () => api<Worker[]>(`/projects/${active}/workers`),
      refetchInterval: 4000,
      enabled: !!active,
    }),
    board = useQuery({
      queryKey: ["board", active],
      queryFn: () => api<Board>(`/projects/${active}/board`),
      refetchInterval: 1000,
      enabled: !!active,
    });
  const refresh = () => {
      qc.invalidateQueries({ queryKey: ["board", active] });
      qc.invalidateQueries({ queryKey: ["workers", active] });
    },
    move = useMutation({
      mutationFn: ({ id, state }: { id: string; state: string }) =>
        api(`/tasks/${id}/move`, {
          method: "POST",
          body: JSON.stringify({ State: state }),
        }),
      onSuccess: refresh,
    });
  const selectProject = (id: string) => {
    const url = new URL(location.href);
    if (id) url.searchParams.set("project", id);
    else url.searchParams.delete("project");
    history.pushState(null, "", url);
    setProjectID(id);
    setSelected(undefined);
    setCreate(false);
    setTab("board");
  };
  useEffect(() => {
    const onPop = () => {
      setProjectID(new URLSearchParams(location.search).get("project") || "");
      setSelected(undefined);
      setCreate(false);
      setTab("board");
    };
    window.addEventListener("popstate", onPop);
    return () => window.removeEventListener("popstate", onPop);
  }, []);
  const onDragEnd = ({ active: drag, over }: DragEndEvent) => {
    if (over && drag.data.current?.state !== over.id)
      move.mutate({ id: String(drag.id), state: String(over.id) });
  };
  if (projects.isLoading) return <Empty title="Loading AgentBoard…" />;
  if (projects.error)
    return (
      <Empty
        title="Could not load projects"
        subtitle={projects.error.message}
      />
    );
  if (!active)
    return (
      <main className="project-launcher">
        <h1>AgentBoard</h1>
        <h2>Projects</h2>
        {projectID && <Message severity="warn" text="The requested project is not registered on this server. Select a project below." />}
        {!projects.data?.length && <Message severity="info" text={<>Run <code>agentboard init</code> in a project folder. This list updates automatically.</>} />}
        {projects.data?.map((p) => (
          <Card className="project-tile" key={p.ID} title={p.Name} subTitle={p.Path}>
            <Button label="Open project" icon="pi pi-arrow-right" iconPos="right" text
              onClick={() => selectProject(p.ID)} />
          </Card>
        ))}
      </main>
    );
  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">
            <Sparkles size={18} />
          </span>
          <span>AgentBoard</span>
        </div>
        <div className="workspace-label">Workspace</div>
        <SelectField
          className="project-select"
          value={active}
          onChange={(e) => selectProject(e.target.value)}
        >
          {projects.data?.map((p) => (
            <option value={p.ID} key={p.ID}>
              {p.Name}
            </option>
          ))}
        </SelectField>
        <Button onClick={() => selectProject("")}>All projects</Button>
        <nav>
          <Button
            className={tab === "board" ? "active" : ""}
            onClick={() => setTab("board")}
          >
            <Columns3 />
            Board
          </Button>
          <Button
            className={tab === "workers" ? "active" : ""}
            onClick={() => setTab("workers")}
          >
            <Users />
            Workers
          </Button>
          <Button
            className={tab === "properties" ? "active" : ""}
            onClick={() => setTab("properties")}
          >
            <Settings2 />
            Properties
          </Button>
        </nav>
        <div className="sidebar-note">
          <Bot size={16} />
          <span>Workers run only when you launch them through MCP.</span>
        </div>
      </aside>
      <main>
        <header>
          <div>
            <p className="eyebrow">
              {projects.data?.find((p) => p.ID === active)?.Path}
            </p>
            <h1>
              {tab === "board"
                ? "Project board"
                : tab === "workers"
                  ? "Workers"
                  : "Custom properties"}
            </h1>
          </div>
          {tab === "board" && (
            <Button className="primary" onClick={() => setCreate(true)}>
              <Plus size={17} />
              New task
            </Button>
          )}
        </header>
        {tab === "board" ? (
          <>
            <FilterBar
              filters={filters}
              setFilters={setFilters}
              workers={workers.data || []}
            />
            <DndContext onDragEnd={onDragEnd}>
              <div className="board">
                {states.map(([key, label]) => (
                  <Column
                    key={key}
                    state={key}
                    label={label}
                    tasks={(board.data?.[key] || []).filter((t) =>
                      matches(t, filters),
                    )}
                    onOpen={setSelected}
                  />
                ))}
              </div>
            </DndContext>
            {move.error && <Message severity="error" text={move.error.message} className="action-error" />}
          </>
        ) : tab === "workers" ? (
          <Workers
            key={active}
            projectPath={
              projects.data?.find((p) => p.ID === active)?.Path || ""
            }
            projectID={active}
            workers={workers.data || []}
            refresh={refresh}
          />
        ) : (
          <Properties projectID={active} />
        )}
      </main>
      {selected && (
        <TaskDrawer
          id={selected}
          workers={workers.data || []}
          close={() => setSelected(undefined)}
          refresh={refresh}
        />
      )}{" "}
      {create && (
        <TaskForm
          projectID={active}
          workers={workers.data || []}
          close={() => setCreate(false)}
          refresh={refresh}
        />
      )}
    </div>
  );
}
function FilterBar({
  filters,
  setFilters,
  workers,
}: {
  filters: any;
  setFilters: (v: any) => void;
  workers: Worker[];
}) {
  const field = (key: string, children: React.ReactNode) => (
    <SelectField
      value={filters[key]}
      onChange={(e) => setFilters({ ...filters, [key]: e.target.value })}
    >
      {children}
    </SelectField>
  );
  return (
    <Toolbar className="filters" start={<span>
        <Filter size={15} />
        Filters
      </span>} end={<div className="filter-controls">{field(
        "responsible",
        <>
          <option value="all">All responsible</option>
          <option value="human">Human</option>
          <option value="unassigned">Unassigned</option>
          {workers.map((w) => (
            <option key={w.ID} value={w.ID}>
              {w.Name}
            </option>
          ))}
        </>,
      )}
      {field(
        "priority",
        <>
          <option value="all">All priority</option>
          <option>critical</option>
          <option>high</option>
          <option>medium</option>
          <option>low</option>
        </>,
      )}
      {field(
        "origin",
        <>
          <option value="all">All origins</option>
          <option value="human">Human-created</option>
          <option value="agent">AI-created</option>
        </>,
      )}
      {field(
        "testing",
        <>
          <option value="all">All testing</option>
          <option value="ai">AI</option>
          <option value="human">Human</option>
          <option value="hybrid">Hybrid</option>
        </>,
      )}</div>} />
  );
}
function matches(t: Task, f: any) {
  return (
    (f.responsible === "all" ||
      f.responsible === t.AssigneeType ||
      f.responsible === t.AssigneeWorkerID) &&
    (f.priority === "all" || f.priority === t.Priority) &&
    (f.origin === "all" || f.origin === t.CreatedByType) &&
    (f.testing === "all" || f.testing === t.TestingMode)
  );
}
function Column({
  state,
  label,
  tasks,
  onOpen,
}: {
  state: string;
  label: string;
  tasks: Task[];
  onOpen: (id: string) => void;
}) {
  const { setNodeRef, isOver } = useDroppable({ id: state });
  return (
    <section ref={setNodeRef} className={"column " + (isOver ? "over" : "")}>
      <Panel className="board-panel" header={<div className="column-title">
        <span className={"state-dot " + state} />
        <h2>{label}</h2>
        <Tag value={String(tasks.length)} rounded severity="secondary" />
      </div>}>
      <div className="card-list">
        {tasks.map((t) => (
          <TaskCard task={t} key={t.ID} onOpen={onOpen} />
        ))}
        {!tasks.length && <Card className="empty-column">Drop tasks here</Card>}
      </div>
      </Panel>
    </section>
  );
}
function TaskCard({
  task,
  onOpen,
}: {
  task: Task;
  onOpen: (id: string) => void;
}) {
  const { attributes, listeners, setNodeRef, transform, isDragging } =
      useDraggable({ id: task.ID, data: { state: task.State } }),
    style = transform
      ? { transform: `translate3d(${transform.x}px,${transform.y}px,0)` }
      : undefined;
  return (
    <article
      ref={setNodeRef}
      style={style}
      className={"task-card " + (isDragging ? "dragging" : "")}
    >
      <Button
        className="drag"
        aria-label={`Drag task ${task.Title}`}
        {...listeners}
        {...attributes}
        onClick={(e) => e.stopPropagation()}
      >
        <GripVertical size={15} />
      </Button>
      <Card className="task-surface" role="button" tabIndex={0}
        aria-label={`Open task ${task.Title}`} onClick={() => onOpen(task.ID)}
        onKeyDown={(event) => {
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            onOpen(task.ID);
          }
        }}>
      <div className="task-key">{short(task.ID)}</div>
      <h3>{task.Title}</h3>
      {task.Description && <p>{task.Description}</p>}
      <div className="badges">
        <Tag value={task.Priority} severity={task.Priority === "critical" ? "danger" : task.Priority === "high" ? "warning" : task.Priority === "low" ? "secondary" : "info"} />
        <Tag value={task.TestingMode} icon="pi pi-check-square" severity="secondary" />
        {task.Dependencies?.length > 0 && (
          <Tag value={`↳ ${task.Dependencies.length}`} severity="secondary" />
        )}
      </div>
      <div className="card-footer">
        <span className="assignee">
          {task.AssigneeType === "worker" ? (
            <Bot size={14} />
          ) : (
            <User size={14} />
          )}{" "}
          {task.AssigneeName || task.AssigneeType}
        </span>
        {task.CreatedByType === "agent" ? (
          <span className="origin ai">
            <Sparkles size={12} />
            AI · {task.CreatorName}
          </span>
        ) : (
          <span className="origin">Human</span>
        )}
      </div>
      <small title={date(task.UpdatedAt)}>
        Updated {relativeTime(task.UpdatedAt)}
      </small>
      </Card>
    </article>
  );
}

function TaskForm({
  projectID,
  workers,
  close,
  refresh,
}: {
  projectID: string;
  workers: Worker[];
  close: () => void;
  refresh: () => void;
}) {
  const definitions = useQuery({
      queryKey: ["properties", projectID],
      queryFn: () => api<PropertyDef[]>(`/projects/${projectID}/properties`),
    }),
    board = useQuery({
      queryKey: ["board", projectID],
      queryFn: () => api<Board>(`/projects/${projectID}/board`),
    }),
    allTasks = Object.values(board.data || {}).flat();
  const [form, setForm] = useState({
      Title: "",
      Description: "",
      State: "backlog",
      Priority: "medium",
      TestingMode: "ai",
      AssigneeType: "unassigned",
      AssigneeWorkerID: "",
      AITestInstructions: "",
      HumanTestInstructions: "",
      DependencyIDs: [] as string[],
      Properties: {} as Record<string, string>,
    }),
    save = useMutation({
      mutationFn: () =>
        api(`/projects/${projectID}/tasks`, {
          method: "POST",
          body: JSON.stringify({
            ...form,
            Assignee: {
              type: form.AssigneeType,
              worker_id:
                form.AssigneeType === "worker"
                  ? form.AssigneeWorkerID
                  : undefined,
            },
          }),
        }),
      onSuccess: () => {
        refresh();
        close();
      },
    });
  return (
    <Modal title="Create task" close={close}>
      <div className="form-grid">
        <label className="wide">
          Title
          <InputText
            autoFocus
            value={form.Title}
            onChange={(e) => setForm({ ...form, Title: e.target.value })}
          />
        </label>
        <label className="wide">
          Description
          <InputTextarea
            value={form.Description}
            onChange={(e) => setForm({ ...form, Description: e.target.value })}
          />
        </label>
        <label>
          State
          <SelectField
            value={form.State}
            onChange={(e) => setForm({ ...form, State: e.target.value })}
          >
            {states.map((s) => (
              <option key={s[0]} value={s[0]}>
                {s[1]}
              </option>
            ))}
          </SelectField>
        </label>
        <label>
          Priority
          <SelectField
            value={form.Priority}
            onChange={(e) => setForm({ ...form, Priority: e.target.value })}
          >
            <option>critical</option>
            <option>high</option>
            <option>medium</option>
            <option>low</option>
          </SelectField>
        </label>
        <label>
          Responsible
          <SelectField
            value={
              form.AssigneeType === "worker"
                ? form.AssigneeWorkerID
                : form.AssigneeType
            }
            onChange={(e) =>
              setForm({
                ...form,
                AssigneeType: workers.some((w) => w.ID === e.target.value)
                  ? "worker"
                  : e.target.value,
                AssigneeWorkerID: e.target.value,
              })
            }
          >
            <option value="unassigned">Unassigned</option>
            <option value="human">Human</option>
            {workers
              .filter((w) => w.Enabled)
              .map((w) => (
                <option key={w.ID} value={w.ID}>
                  {w.Name}
                </option>
              ))}
          </SelectField>
        </label>
        <label>
          Testing mode
          <SelectField
            value={form.TestingMode}
            onChange={(e) => setForm({ ...form, TestingMode: e.target.value })}
          >
            <option value="ai">AI</option>
            <option value="human">Human</option>
            <option value="hybrid">Hybrid</option>
          </SelectField>
        </label>
        <label className="wide">
          Dependencies
          <SelectField
            multiple
            value={form.DependencyIDs}
            onChange={(e) =>
              setForm({
                ...form,
                DependencyIDs: Array.from(
                  e.target.selectedOptions,
                  (o) => o.value,
                ),
              })
            }
          >
            {allTasks.map((t) => (
              <option key={t.ID} value={t.ID}>
                {short(t.ID)} · {t.Title}
              </option>
            ))}
          </SelectField>
        </label>
        {definitions.data?.map((d) => (
          <label key={d.ID}>
            {d.Name}
            <InputText
              value={form.Properties[d.ID] || ""}
              onChange={(e) =>
                setForm({
                  ...form,
                  Properties: { ...form.Properties, [d.ID]: e.target.value },
                })
              }
            />
            <small>{d.Visibility.replaceAll("_", " ")}</small>
          </label>
        ))}
        <label className="wide">
          AI test instructions
          <InputTextarea
            value={form.AITestInstructions}
            onChange={(e) =>
              setForm({ ...form, AITestInstructions: e.target.value })
            }
          />
        </label>
        <label className="wide">
          Human test instructions
          <InputTextarea
            value={form.HumanTestInstructions}
            onChange={(e) =>
              setForm({ ...form, HumanTestInstructions: e.target.value })
            }
          />
        </label>
      </div>
      {save.error && <Message severity="error" text={save.error.message} className="action-error" />}
      <DialogActions
        close={close}
        save={() => save.mutate()}
        busy={save.isPending}
        disabled={!form.Title}
      />
    </Modal>
  );
}

function TaskDrawer({
  id,
  workers,
  close,
  refresh,
}: {
  id: string;
  workers: Worker[];
  close: () => void;
  refresh: () => void;
}) {
  const [editing, setEditing] = useState(false),
    qc = useQueryClient(),
    q = useQuery({
      queryKey: ["task", id],
      queryFn: () => api<Details>(`/tasks/${id}`),
      refetchInterval: 1000,
    }),
    reload = () => {
      qc.invalidateQueries({ queryKey: ["task", id] });
      refresh();
    },
    assign = useMutation({
      mutationFn: (value: string) =>
        api(`/tasks/${id}/assign`, {
          method: "POST",
          body: JSON.stringify(
            workers.some((w) => w.ID === value)
              ? { type: "worker", worker_id: value }
              : { type: value },
          ),
        }),
      onSuccess: reload,
    }),
    comment = useMutation({
      mutationFn: (content: string) =>
        api(`/tasks/${id}/comments`, {
          method: "POST",
          body: JSON.stringify({ Content: content }),
        }),
      onSuccess: reload,
    });
  if (!q.data)
    return (
      <Sidebar visible onHide={close} position="right" className="task-sidebar">
        {q.error ? <Message severity="error" text={q.error.message} /> : <ProgressSpinner />}
      </Sidebar>
    );
  const t = q.data;
  return (
    <Sidebar visible onHide={close} position="right" className="task-sidebar" showCloseIcon={false} blockScroll>
      <aside className="drawer">
        <div className="drawer-header">
          <div className="drawer-key">
            {short(t.ID)} · {t.State.replaceAll("_", " ")}
          </div>
          <div className="drawer-actions">
            <Button onClick={() => setEditing(true)}>Edit task</Button>
            <Button
              className="icon-button"
              aria-label="Close task"
              onClick={close}
            >
              <X />
            </Button>
          </div>
        </div>
        <h2>{t.Title}</h2>
        <p className="description">{t.Description || "No description."}</p>
        <div className="detail-grid">
          <Info
            label="Created by"
            value={t.CreatedByType === "agent" ? t.CreatorName : "Human"}
          />
          <Info label="Created at" value={date(t.CreatedAt)} />
          <Info
            label="Source task"
            value={t.SourceTaskID ? short(t.SourceTaskID) : "—"}
          />
          <label>
            Responsible
            <SelectField
              value={
                t.AssigneeType === "worker"
                  ? t.AssigneeWorkerID
                  : t.AssigneeType
              }
              onChange={(e) => assign.mutate(e.target.value)}
            >
              <option value="unassigned">Unassigned</option>
              <option value="human">Human</option>
              {workers.map((w) => (
                <option key={w.ID} value={w.ID}>
                  {w.Name}
                  {!w.Enabled ? " (disabled)" : ""}
                </option>
              ))}
            </SelectField>
          </label>
          <Info label="Priority" value={t.Priority} />
          <Info label="Testing" value={t.TestingMode} />
        </div>
        {assign.error && <Message severity="error" text={assign.error.message} className="action-error" />}
        <Panel className="detail-panel" header="Test instructions">
          <div className="instruction">
            <Tag value="AI" severity="info" />
            <span>{t.AITestInstructions || "None"}</span>
          </div>
          <div className="instruction">
            <Tag value="Human" severity="secondary" />
            <span>{t.HumanTestInstructions || "None"}</span>
          </div>
        </Panel>
        {t.Dependencies?.length > 0 && (
          <Panel className="detail-panel" header="Dependencies">
            {t.Dependencies.map((d) => (
              <Card className="related-task" key={d.DependsOnTaskID}>
                <span>
                  <b>{short(d.DependsOnTaskID)}</b> {d.Title}
                </span>
                <Tag value={d.State.replaceAll("_", " ")} severity="secondary" />
                <Badge value={`${d.Artifacts?.length || 0} artifacts`} severity="info" />
              </Card>
            ))}
          </Panel>
        )}
        {t.SpawnedTasks?.length > 0 && (
          <Panel className="detail-panel" header="Spawned tasks">
            {t.SpawnedTasks.map((s) => (
              <Card className="related-task" key={s.ID}>
                <span>
                  {short(s.ID)} · {s.Title}
                </span>
                <Tag value={s.State.replaceAll("_", " ")} severity="secondary" />
              </Card>
            ))}
          </Panel>
        )}
        <Tabs details={t} addComment={(v) => comment.mutate(v)} />
        {editing && (
          <TaskEditForm
            task={t}
            close={() => setEditing(false)}
            refresh={reload}
          />
        )}
      </aside>
    </Sidebar>
  );
}

function TaskEditForm({
  task,
  close,
  refresh,
}: {
  task: Details;
  close: () => void;
  refresh: () => void;
}) {
  const definitions = useQuery({
      queryKey: ["properties", task.ProjectID],
      queryFn: () =>
        api<PropertyDef[]>(`/projects/${task.ProjectID}/properties`),
    }),
    board = useQuery({
      queryKey: ["board", task.ProjectID],
      queryFn: () => api<Board>(`/projects/${task.ProjectID}/board`),
    });
  const [form, setForm] = useState({
    Title: task.Title,
    Description: task.Description,
    Priority: task.Priority,
    TestingMode: task.TestingMode,
    AITestInstructions: task.AITestInstructions,
    HumanTestInstructions: task.HumanTestInstructions,
    DependencyIDs: task.Dependencies?.map((d) => d.DependsOnTaskID) || [],
    Properties: Object.fromEntries(
      (task.Properties || []).map((p) => [p.PropertyDefinitionID, p.Value]),
    ) as Record<string, string>,
  });
  const save = useMutation({
      mutationFn: () =>
        api(`/tasks/${task.ID}`, {
          method: "PATCH",
          body: JSON.stringify(form),
        }),
      onSuccess: () => {
        refresh();
        close();
      },
    }),
    allTasks = Object.values(board.data || {})
      .flat()
      .filter((t) => t.ID !== task.ID);
  return (
    <Modal title="Edit task" close={close}>
      <div className="form-grid">
        <label className="wide">
          Title
          <InputText
            value={form.Title}
            onChange={(e) => setForm({ ...form, Title: e.target.value })}
          />
        </label>
        <label className="wide">
          Description
          <InputTextarea
            value={form.Description}
            onChange={(e) => setForm({ ...form, Description: e.target.value })}
          />
        </label>
        <label>
          Priority
          <SelectField
            value={form.Priority}
            onChange={(e) => setForm({ ...form, Priority: e.target.value })}
          >
            <option>critical</option>
            <option>high</option>
            <option>medium</option>
            <option>low</option>
          </SelectField>
        </label>
        <label>
          Testing mode
          <SelectField
            value={form.TestingMode}
            onChange={(e) => setForm({ ...form, TestingMode: e.target.value })}
          >
            <option value="ai">AI</option>
            <option value="human">Human</option>
            <option value="hybrid">Hybrid</option>
          </SelectField>
        </label>
        <label className="wide">
          Dependencies
          <SelectField
            multiple
            value={form.DependencyIDs}
            onChange={(e) =>
              setForm({
                ...form,
                DependencyIDs: Array.from(
                  e.target.selectedOptions,
                  (o) => o.value,
                ),
              })
            }
          >
            {allTasks.map((t) => (
              <option key={t.ID} value={t.ID}>
                {short(t.ID)} · {t.Title}
              </option>
            ))}
          </SelectField>
        </label>
        {definitions.data?.map((d) => (
          <label key={d.ID}>
            {d.Name}
            <InputText
              value={form.Properties[d.ID] || ""}
              onChange={(e) =>
                setForm({
                  ...form,
                  Properties: { ...form.Properties, [d.ID]: e.target.value },
                })
              }
            />
            <small>{d.Visibility.replaceAll("_", " ")}</small>
          </label>
        ))}
        <label className="wide">
          AI test instructions
          <InputTextarea
            value={form.AITestInstructions}
            onChange={(e) =>
              setForm({ ...form, AITestInstructions: e.target.value })
            }
          />
        </label>
        <label className="wide">
          Human test instructions
          <InputTextarea
            value={form.HumanTestInstructions}
            onChange={(e) =>
              setForm({ ...form, HumanTestInstructions: e.target.value })
            }
          />
        </label>
      </div>
      {save.error && <Message severity="error" text={save.error.message} className="action-error" />}
      <DialogActions
        close={close}
        save={() => save.mutate()}
        busy={save.isPending}
        disabled={!form.Title}
      />
    </Modal>
  );
}

function Tabs({
  details,
  addComment,
}: {
  details: Details;
  addComment: (v: string) => void;
}) {
  const [tab, setTab] = useState("History"),
    [text, setText] = useState(""),
    tabs = ["History", "Testing", "Artifacts", "AI usage"];
  return (
    <section className="detail-section grow">
      <TabMenu className="detail-tabs" model={tabs.map((x) => ({
        label: `${x} · ${x === "History" ? details.History?.length || 0 : x === "Testing" ? details.TestRuns?.length || 0 : x === "Artifacts" ? details.Artifacts?.length || 0 : details.Usage?.length || 0}`,
        command: () => setTab(x),
      }))} activeIndex={tabs.indexOf(tab)} />
      {tab === "History" && (
        <>
          {details.History?.length ? <Timeline value={details.History} dataKey="ID" className="activity-timeline"
            marker={(h: History) => <Avatar icon={h.ActorType === "agent" ? "pi pi-sparkles" : h.ActorType === "human" ? "pi pi-user" : "pi pi-cog"}
              shape="circle" className={`actor-${h.ActorType}`} />}
            content={(h: History) => <Card className="activity-card">
              <div className="activity-heading"><b>{h.ActorName}</b><time>{date(h.CreatedAt)}</time></div>
              <p>{h.Content}</p><Tag value={h.EntryType.replaceAll("_", " ")} severity="secondary" />
            </Card>} /> : <Message severity="info" text="No history yet" />}
          <form
            className="comment"
            onSubmit={(e) => {
              e.preventDefault();
              if (text.trim()) {
                addComment(text);
                setText("");
              }
            }}
          >
            <InputText
              placeholder="Add a Human comment…"
              value={text}
              onChange={(e) => setText(e.target.value)}
            />
            <Button>Send</Button>
          </form>
        </>
      )}
      {tab === "Testing" && (
        details.TestRuns?.length ? <Timeline value={details.TestRuns} dataKey="ID" className="activity-timeline"
          marker={(r: TestRun) => <Avatar icon={r.Status === "passed" ? "pi pi-check" : "pi pi-exclamation-circle"}
            shape="circle" className={`run-${r.Status}`} />}
          content={(r: TestRun) => <Card className="activity-card">
            <div className="activity-heading"><b>{r.Runner}</b><time>{date(r.CreatedAt)}</time></div>
            <Tag value={r.Status} severity={r.Status === "passed" ? "success" : "warning"} />
            <p>{r.Summary}</p>{r.Command && <code>{r.Command}</code>}
          </Card>} /> : <Message severity="info" text="No test runs yet" />
      )}
      {tab === "Artifacts" && (
        <div className="artifact-list">
          {!details.Artifacts?.length && <Message severity="info" text="No artifacts yet" />}
          {details.Artifacts?.map((a) => (
            <Card className="artifact" key={a.ID} title={a.Name} subTitle={`${a.Kind} · ${a.CreatorName || a.CreatedByType} · ${date(a.CreatedAt)}`}>
              <p>{a.Path || a.URL}</p>
              {a.Path && (
                <Button label="Copy path" icon="pi pi-copy" text onClick={() => navigator.clipboard.writeText(a.Path!)} />
              )}
            </Card>
          ))}
        </div>
      )}
      {tab === "AI usage" && (
        <div className="usage-list">
          {!details.Usage?.length && <Message severity="info" text="No AI usage yet" />}
          {details.Usage?.map((u) => (
            <Card className="usage" key={u.ID} title={u.WorkerName} subTitle={`${u.Provider || "—"} / ${u.Model || "—"}`}>
              <Divider />
              <dl>
                <dt>Input</dt>
                <dd>{u.InputTokens ?? "—"}</dd>
                <dt>Output</dt>
                <dd>{u.OutputTokens ?? "—"}</dd>
                <dt>Reasoning</dt>
                <dd>{u.ReasoningTokens ?? "—"}</dd>
                <dt>Cached</dt>
                <dd>{u.CachedTokens ?? "—"}</dd>
                <dt>Model calls</dt>
                <dd>{u.ModelCalls ?? "—"}</dd>
                <dt>Tool calls</dt>
                <dd>{u.ToolCalls ?? "—"}</dd>
                <dt>MCP calls</dt>
                <dd>{u.MCPCalls ?? "—"}</dd>
              </dl>
            </Card>
          ))}
        </div>
      )}
    </section>
  );
}

function Workers({
  projectPath,
  projectID,
  workers,
  refresh,
}: {
  projectID: string;
  projectPath: string;
  workers: Worker[];
  refresh: () => void;
}) {
  const [editing, setEditing] = useState<Worker>(),
    [adding, setAdding] = useState(false);
  return (
    <div className="workers-page">
      <Toolbar className="section-head" start={<div>
          <h2>Registered workers</h2>
          <p>
            Logical AI identities. AgentBoard never launches them automatically.
          </p>
        </div>} end={<Button className="primary" onClick={() => setAdding(true)}>
          <Plus size={16} />
          Add worker
        </Button>} />
      <DataTable value={workers} dataKey="ID" className="board-table" stripedRows
        rowHover onRowClick={(event) => setEditing(event.data as Worker)} emptyMessage="No workers yet">
        <DataColumn header="Name" body={(w: Worker) => <span><b>{w.Name}</b><small>{w.Slug}</small></span>} />
        <DataColumn field="Kind" header="Kind" />
        <DataColumn header="Status" body={(w: Worker) => <span>
          <Tag value={w.Enabled ? w.ActiveSessionCount > 0 ? "Active" : "Offline" : "Disabled"}
            severity={w.Enabled && w.ActiveSessionCount > 0 ? "success" : "secondary"} />
          <small>{w.SessionCount} sessions · {w.MCPCalls} calls</small>
          <small>Last MCP: {w.LastActivityAt ? relativeTime(w.LastActivityAt) : "Never"}</small>
        </span>} />
        <DataColumn header="Capabilities" body={(w: Worker) => capabilities(w.Capabilities).join(", ") || "—"} />
        <DataColumn field="AssignedTaskCount" header="Assigned" />
      </DataTable>
      {(adding || editing) && (
        <WorkerForm
          projectPath={projectPath}
          projectID={projectID}
          worker={editing}
          close={() => {
            setAdding(false);
            setEditing(undefined);
          }}
          refresh={refresh}
        />
      )}
    </div>
  );
}
function Properties({ projectID }: { projectID: string }) {
  const qc = useQueryClient(),
    q = useQuery({
      queryKey: ["properties", projectID],
      queryFn: () => api<PropertyDef[]>(`/projects/${projectID}/properties`),
    }),
    [editing, setEditing] = useState<PropertyDef>(),
    [adding, setAdding] = useState(false),
    refresh = () =>
      qc.invalidateQueries({ queryKey: ["properties", projectID] });
  return (
    <div className="workers-page">
      <Toolbar className="section-head" start={<div>
          <h2>Task properties</h2>
          <p>Trello-like fields with explicit Worker visibility.</p>
        </div>} end={<Button className="primary" onClick={() => setAdding(true)}>
          <Plus size={16} />
          Add property
        </Button>} />
      <DataTable value={q.data || []} dataKey="ID" className="board-table" stripedRows
        rowHover onRowClick={(event) => setEditing(event.data as PropertyDef)} emptyMessage="No properties yet">
        <DataColumn field="Name" header="Name" />
        <DataColumn field="Type" header="Type" />
        <DataColumn header="Visibility" body={(p: PropertyDef) => p.Visibility.replaceAll("_", " ")} />
        <DataColumn field="Options" header="Options" />
      </DataTable>
      {(adding || editing) && (
        <PropertyForm
          projectID={projectID}
          property={editing}
          close={() => {
            setAdding(false);
            setEditing(undefined);
          }}
          refresh={refresh}
        />
      )}
    </div>
  );
}
function PropertyForm({
  projectID,
  property,
  close,
  refresh,
}: {
  projectID: string;
  property?: PropertyDef;
  close: () => void;
  refresh: () => void;
}) {
  const [confirmRemove, setConfirmRemove] = useState(false);
  const [form, setForm] = useState({
      Name: property?.Name || "",
      Type: property?.Type || "text",
      Options: property?.Options || "[]",
      Visibility: property?.Visibility || "agent_read",
    }),
    save = useMutation({
      mutationFn: () =>
        api(
          property
            ? `/properties/${property.ID}`
            : `/projects/${projectID}/properties`,
          { method: property ? "PATCH" : "POST", body: JSON.stringify(form) },
        ),
      onSuccess: () => {
        refresh();
        close();
      },
    }),
    remove = useMutation({
      mutationFn: () =>
        api(`/properties/${property!.ID}`, { method: "DELETE" }),
      onSuccess: () => {
        refresh();
        close();
      },
    });
  return (
    <Modal title={property ? "Edit property" : "Add property"} close={close}>
      <div className="form-grid">
        <label>
          Name
          <InputText
            value={form.Name}
            onChange={(e) => setForm({ ...form, Name: e.target.value })}
          />
        </label>
        <label>
          Type
          <SelectField
            value={form.Type}
            onChange={(e) => setForm({ ...form, Type: e.target.value })}
          >
            {[
              "text",
              "number",
              "boolean",
              "date",
              "datetime",
              "select",
              "multi_select",
              "url",
            ].map((x) => (
              <option key={x}>{x}</option>
            ))}
          </SelectField>
        </label>
        <label className="wide">
          Visibility
          <SelectField
            value={form.Visibility}
            onChange={(e) => setForm({ ...form, Visibility: e.target.value })}
          >
            <option value="human_only">Human only</option>
            <option value="agent_read">Agent read</option>
            <option value="agent_read_write">Agent read/write</option>
          </SelectField>
        </label>
        {["select", "multi_select"].includes(form.Type) && (
          <label className="wide">
            Options JSON
            <InputTextarea
              className="mono"
              value={form.Options}
              onChange={(e) => setForm({ ...form, Options: e.target.value })}
            />
          </label>
        )}
      </div>
      {save.error && <Message severity="error" text={save.error.message} className="action-error" />}
      {remove.error && <Message severity="error" text={remove.error.message} className="action-error" />}
      {property && (
        <Button className="danger-link" onClick={() => setConfirmRemove(true)}>
          Delete property
        </Button>
      )}
      <ConfirmDialog visible={confirmRemove} onHide={() => setConfirmRemove(false)}
        header="Delete property?" message={`Delete ${property?.Name}? Task values for this property will be removed.`}
        icon="pi pi-exclamation-triangle" acceptLabel="Delete" rejectLabel="Cancel"
        acceptClassName="p-button-danger" defaultFocus="reject" accept={() => remove.mutate()} />
      <DialogActions
        close={close}
        save={() => save.mutate()}
        busy={save.isPending}
        disabled={!form.Name}
      />
    </Modal>
  );
}
function capabilities(raw: string) {
  try {
    return Object.entries(JSON.parse(raw))
      .filter(([, v]) => v)
      .map(([k]) => k.replace("_", " "));
  } catch {
    return [];
  }
}
function WorkerForm({
  projectPath,
  projectID,
  worker,
  close,
  refresh,
}: {
  projectID: string;
  projectPath: string;
  worker?: Worker;
  close: () => void;
  refresh: () => void;
}) {
  const [confirmArchive, setConfirmArchive] = useState(false);
  const [preset, setPreset] = useState(worker ? "Custom" : "Generic");
  const [form, setForm] = useState({
      Name: worker?.Name || "",
      Slug: worker?.Slug || "",
      Description: worker?.Description || "",
      Kind: worker?.Kind || "external",
      Capabilities:
        worker?.Capabilities ||
        JSON.stringify(capabilityPresets.Generic, null, 2),
      Enabled: worker?.Enabled ?? true,
    }),
    save = useMutation({
      mutationFn: () =>
        api(
          worker ? `/workers/${worker.ID}` : `/projects/${projectID}/workers`,
          { method: worker ? "PATCH" : "POST", body: JSON.stringify(form) },
        ),
      onSuccess: () => {
        refresh();
        close();
      },
    }),
    archive = useMutation({
      mutationFn: () => api(`/workers/${worker!.ID}`, { method: "DELETE" }),
      onSuccess: () => {
        refresh();
        close();
      },
    });
  return (
    <Modal title={worker ? "Edit worker" : "Add worker"} close={close}>
      <div className="form-grid">
        <label>
          Name
          <InputText
            value={form.Name}
            onChange={(e) => setForm({ ...form, Name: e.target.value })}
          />
        </label>
        <label>
          Slug
          <InputText
            value={form.Slug}
            onChange={(e) => setForm({ ...form, Slug: e.target.value })}
          />
        </label>
        <label className="wide">
          Description
          <InputTextarea
            value={form.Description}
            onChange={(e) => setForm({ ...form, Description: e.target.value })}
          />
        </label>
        <label>
          Kind
          <SelectField
            value={form.Kind}
            onChange={(e) => setForm({ ...form, Kind: e.target.value })}
          >
            <option value="external">external</option>
            <option value="custom">custom</option>
          </SelectField>
        </label>
        <label className="check">
          <InputSwitch checked={form.Enabled}
            onChange={(e) => setForm({ ...form, Enabled: !!e.value })} />
          Enabled
        </label>
        <label className="wide">
          Capability preset
          <SelectField
            value={preset}
            onChange={(e) => {
              setPreset(e.target.value);
              if (e.target.value !== "Custom")
                setForm({
                  ...form,
                  Capabilities: JSON.stringify(
                    capabilityPresets[e.target.value],
                    null,
                    2,
                  ),
                });
            }}
          >
            {Object.keys(capabilityPresets).map((name) => (
              <option key={name}>{name}</option>
            ))}
            <option>Custom</option>
          </SelectField>
        </label>
        <label className="wide">
          Advanced: capabilities JSON
          <InputTextarea
            className="mono"
            rows={9}
            value={form.Capabilities}
            onChange={(e) => {
              setPreset("Custom");
              setForm({ ...form, Capabilities: e.target.value });
            }}
          />
        </label>
      </div>
      {save.error && <Message severity="error" text={save.error.message} className="action-error" />}
      {archive.error && <Message severity="error" text={archive.error.message} className="action-error" />}
      {worker && (
        <WorkerDiagnostics worker={worker} projectPath={projectPath} />
      )}
      {worker && (
        <Button className="danger-link" onClick={() => setConfirmArchive(true)}>
          Archive worker
        </Button>
      )}
      <ConfirmDialog visible={confirmArchive} onHide={() => setConfirmArchive(false)}
        header="Archive worker?" message={`Archive ${worker?.Name}? It will no longer be available for new task assignments.`}
        icon="pi pi-exclamation-triangle" acceptLabel="Archive" rejectLabel="Cancel"
        acceptClassName="p-button-danger" defaultFocus="reject" accept={() => archive.mutate()} />
      <DialogActions
        close={close}
        save={() => save.mutate()}
        busy={save.isPending}
        disabled={!form.Name || !form.Slug}
      />
    </Modal>
  );
}
function Modal({
  title,
  close,
  children,
}: {
  title: string;
  close: () => void;
  children: React.ReactNode;
}) {
  return (
    <Dialog header={title} visible onHide={close} modal blockScroll draggable={false}
      className="form-dialog" style={{ width: "min(760px, calc(100vw - 24px))" }}>
      {children}
    </Dialog>
  );
}
function DialogActions({
  close,
  save,
  busy,
  disabled,
}: {
  close: () => void;
  save: () => void;
  busy: boolean;
  disabled?: boolean;
}) {
  return (
    <Toolbar className="dialog-actions" end={<>
      <Button label="Cancel" text severity="secondary" onClick={close} />
      <Button className="primary" loading={busy} disabled={busy || disabled} onClick={save}>
        {busy ? "Saving…" : "Save"}
      </Button>
    </>} />
  );
}
function Info({ label, value }: { label: string; value: string }) {
  return (
    <Card className="info">
      <small>{label}</small>
      <span>{value}</span>
    </Card>
  );
}
function Empty({ title, subtitle }: { title: string; subtitle?: string }) {
  return (
    <div className="empty-screen">
      {subtitle ? <Message severity="error" text={subtitle} /> : <ProgressSpinner />}
      <h1>{title}</h1>
    </div>
  );
}
export default App;
