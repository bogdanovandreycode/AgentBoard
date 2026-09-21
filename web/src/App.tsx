import { useEffect, useState } from "react";
import {
  DndContext,
  useDraggable,
  useDroppable,
  type DragEndEvent,
} from "@dnd-kit/core";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Bot,
  CheckCircle2,
  ClipboardCheck,
  Columns3,
  ExternalLink,
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
import { capabilityPresets, relativeTime } from "./workerPresentation";
import "./App.css";
import "./Properties.css";

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
        {projectID && (
          <p role="alert">
            The requested project is not registered on this server. Select a
            project below.
          </p>
        )}
        {!projects.data?.length && (
          <p>
            Run <code>agentboard init</code> in a project folder. This list
            updates automatically.
          </p>
        )}
        {projects.data?.map((p) => (
          <button
            className="project-tile"
            key={p.ID}
            onClick={() => selectProject(p.ID)}
          >
            <strong>{p.Name}</strong>
            <span>{p.Path}</span>
          </button>
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
        <select
          className="project-select"
          value={active}
          onChange={(e) => selectProject(e.target.value)}
        >
          {projects.data?.map((p) => (
            <option value={p.ID} key={p.ID}>
              {p.Name}
            </option>
          ))}
        </select>
        <button onClick={() => selectProject("")}>All projects</button>
        <nav>
          <button
            className={tab === "board" ? "active" : ""}
            onClick={() => setTab("board")}
          >
            <Columns3 />
            Board
          </button>
          <button
            className={tab === "workers" ? "active" : ""}
            onClick={() => setTab("workers")}
          >
            <Users />
            Workers
          </button>
          <button
            className={tab === "properties" ? "active" : ""}
            onClick={() => setTab("properties")}
          >
            <Settings2 />
            Properties
          </button>
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
            <button className="primary" onClick={() => setCreate(true)}>
              <Plus size={17} />
              New task
            </button>
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
    <select
      value={filters[key]}
      onChange={(e) => setFilters({ ...filters, [key]: e.target.value })}
    >
      {children}
    </select>
  );
  return (
    <div className="filters">
      <span>
        <Filter size={15} />
        Filters
      </span>
      {field(
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
      )}
    </div>
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
      <div className="column-title">
        <span className={"state-dot " + state} />
        <h2>{label}</h2>
        <b>{tasks.length}</b>
      </div>
      <div className="card-list">
        {tasks.map((t) => (
          <TaskCard task={t} key={t.ID} onOpen={onOpen} />
        ))}
        {!tasks.length && <div className="empty-column">Drop tasks here</div>}
      </div>
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
      onClick={() => onOpen(task.ID)}
    >
      <button
        className="drag"
        {...listeners}
        {...attributes}
        onClick={(e) => e.stopPropagation()}
      >
        <GripVertical size={15} />
      </button>
      <div className="task-key">{short(task.ID)}</div>
      <h3>{task.Title}</h3>
      {task.Description && <p>{task.Description}</p>}
      <div className="badges">
        <span className={"badge priority " + task.Priority}>
          {task.Priority}
        </span>
        <span className="badge test">
          <ClipboardCheck size={12} />
          {task.TestingMode}
        </span>
        {task.Dependencies?.length > 0 && (
          <span className="badge">↳ {task.Dependencies.length}</span>
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
          <input
            autoFocus
            value={form.Title}
            onChange={(e) => setForm({ ...form, Title: e.target.value })}
          />
        </label>
        <label className="wide">
          Description
          <textarea
            value={form.Description}
            onChange={(e) => setForm({ ...form, Description: e.target.value })}
          />
        </label>
        <label>
          State
          <select
            value={form.State}
            onChange={(e) => setForm({ ...form, State: e.target.value })}
          >
            {states.map((s) => (
              <option key={s[0]} value={s[0]}>
                {s[1]}
              </option>
            ))}
          </select>
        </label>
        <label>
          Priority
          <select
            value={form.Priority}
            onChange={(e) => setForm({ ...form, Priority: e.target.value })}
          >
            <option>critical</option>
            <option>high</option>
            <option>medium</option>
            <option>low</option>
          </select>
        </label>
        <label>
          Responsible
          <select
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
          </select>
        </label>
        <label>
          Testing mode
          <select
            value={form.TestingMode}
            onChange={(e) => setForm({ ...form, TestingMode: e.target.value })}
          >
            <option value="ai">AI</option>
            <option value="human">Human</option>
            <option value="hybrid">Hybrid</option>
          </select>
        </label>
        <label className="wide">
          Dependencies
          <select
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
          </select>
        </label>
        {definitions.data?.map((d) => (
          <label key={d.ID}>
            {d.Name}
            <input
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
          <textarea
            value={form.AITestInstructions}
            onChange={(e) =>
              setForm({ ...form, AITestInstructions: e.target.value })
            }
          />
        </label>
        <label className="wide">
          Human test instructions
          <textarea
            value={form.HumanTestInstructions}
            onChange={(e) =>
              setForm({ ...form, HumanTestInstructions: e.target.value })
            }
          />
        </label>
      </div>
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
      <div className="drawer">
        <button className="icon-button close" onClick={close}>
          <X />
        </button>
        <p>{q.error ? q.error.message : "Loading…"}</p>
      </div>
    );
  const t = q.data;
  return (
    <div className="drawer-backdrop" onMouseDown={close}>
      <aside className="drawer" onMouseDown={(e) => e.stopPropagation()}>
        <div className="drawer-header">
          <div className="drawer-key">
            {short(t.ID)} · {t.State.replaceAll("_", " ")}
          </div>
          <div className="drawer-actions">
            <button onClick={() => setEditing(true)}>Edit task</button>
            <button
              className="icon-button"
              aria-label="Close task"
              onClick={close}
            >
              <X />
            </button>
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
            <select
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
            </select>
          </label>
          <Info label="Priority" value={t.Priority} />
          <Info label="Testing" value={t.TestingMode} />
        </div>
        <section className="detail-section">
          <h3>Test instructions</h3>
          <div className="instruction">
            <b>AI</b>
            <span>{t.AITestInstructions || "None"}</span>
          </div>
          <div className="instruction">
            <b>Human</b>
            <span>{t.HumanTestInstructions || "None"}</span>
          </div>
        </section>
        {t.Dependencies?.length > 0 && (
          <section className="detail-section">
            <h3>Dependencies</h3>
            {t.Dependencies.map((d) => (
              <div className="row-item" key={d.DependsOnTaskID}>
                <span>
                  <b>{short(d.DependsOnTaskID)}</b> {d.Title}
                </span>
                <span>
                  {d.State} · {d.Artifacts?.length || 0} artifacts
                </span>
              </div>
            ))}
          </section>
        )}
        {t.SpawnedTasks?.length > 0 && (
          <section className="detail-section">
            <h3>Spawned tasks</h3>
            {t.SpawnedTasks.map((s) => (
              <div className="row-item" key={s.ID}>
                <span>
                  {short(s.ID)} · {s.Title}
                </span>
                <span>{s.State}</span>
              </div>
            ))}
          </section>
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
    </div>
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
          <input
            value={form.Title}
            onChange={(e) => setForm({ ...form, Title: e.target.value })}
          />
        </label>
        <label className="wide">
          Description
          <textarea
            value={form.Description}
            onChange={(e) => setForm({ ...form, Description: e.target.value })}
          />
        </label>
        <label>
          Priority
          <select
            value={form.Priority}
            onChange={(e) => setForm({ ...form, Priority: e.target.value })}
          >
            <option>critical</option>
            <option>high</option>
            <option>medium</option>
            <option>low</option>
          </select>
        </label>
        <label>
          Testing mode
          <select
            value={form.TestingMode}
            onChange={(e) => setForm({ ...form, TestingMode: e.target.value })}
          >
            <option value="ai">AI</option>
            <option value="human">Human</option>
            <option value="hybrid">Hybrid</option>
          </select>
        </label>
        <label className="wide">
          Dependencies
          <select
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
          </select>
        </label>
        {definitions.data?.map((d) => (
          <label key={d.ID}>
            {d.Name}
            <input
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
          <textarea
            value={form.AITestInstructions}
            onChange={(e) =>
              setForm({ ...form, AITestInstructions: e.target.value })
            }
          />
        </label>
        <label className="wide">
          Human test instructions
          <textarea
            value={form.HumanTestInstructions}
            onChange={(e) =>
              setForm({ ...form, HumanTestInstructions: e.target.value })
            }
          />
        </label>
      </div>
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
      <div className="tabs">
        {tabs.map((x) => (
          <button
            key={x}
            className={tab === x ? "active" : ""}
            onClick={() => setTab(x)}
          >
            {x}
            <small>
              {x === "History"
                ? details.History?.length
                : x === "Testing"
                  ? details.TestRuns?.length
                  : x === "Artifacts"
                    ? details.Artifacts?.length
                    : details.Usage?.length}
            </small>
          </button>
        ))}
      </div>
      {tab === "History" && (
        <>
          <div className="timeline">
            {details.History?.map((h) => (
              <div className="event" key={h.ID}>
                <span className={"avatar " + h.ActorType}>
                  {h.ActorType === "agent" ? (
                    <Bot />
                  ) : h.ActorType === "human" ? (
                    <User />
                  ) : (
                    <Settings2 />
                  )}
                </span>
                <div>
                  <b>{h.ActorName}</b>
                  <time>{date(h.CreatedAt)}</time>
                  <p>{h.Content}</p>
                  <small>{h.EntryType}</small>
                </div>
              </div>
            ))}
          </div>
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
            <input
              placeholder="Add a Human comment…"
              value={text}
              onChange={(e) => setText(e.target.value)}
            />
            <button>Send</button>
          </form>
        </>
      )}
      {tab === "Testing" && (
        <div className="timeline">
          {details.TestRuns?.map((r) => (
            <div className="event" key={r.ID}>
              <span className={"status-icon " + r.Status}>
                <CheckCircle2 />
              </span>
              <div>
                <b>
                  {r.Status} · {r.Runner}
                </b>
                <time>{date(r.CreatedAt)}</time>
                <p>{r.Summary}</p>
                {r.Command && <code>{r.Command}</code>}
              </div>
            </div>
          ))}
        </div>
      )}
      {tab === "Artifacts" && (
        <div>
          {details.Artifacts?.map((a) => (
            <div className="artifact" key={a.ID}>
              <ExternalLink />
              <div>
                <b>{a.Name}</b>
                <p>{a.Path || a.URL}</p>
                <small>
                  {a.Kind} · {a.CreatorName || a.CreatedByType} ·{" "}
                  {date(a.CreatedAt)}
                </small>
              </div>
              {a.Path && (
                <button onClick={() => navigator.clipboard.writeText(a.Path!)}>
                  Copy path
                </button>
              )}
            </div>
          ))}
        </div>
      )}
      {tab === "AI usage" && (
        <div>
          {details.Usage?.map((u) => (
            <div className="usage" key={u.ID}>
              <b>{u.WorkerName}</b>
              <span>
                {u.Provider || "—"} / {u.Model || "—"}
              </span>
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
            </div>
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
      <div className="section-head">
        <div>
          <h2>Registered workers</h2>
          <p>
            Logical AI identities. AgentBoard never launches them automatically.
          </p>
        </div>
        <button className="primary" onClick={() => setAdding(true)}>
          <Plus size={16} />
          Add worker
        </button>
      </div>
      <div className="worker-table">
        <div className="worker-row header">
          <span>Name</span>
          <span>Kind</span>
          <span>Status</span>
          <span>Capabilities</span>
          <span>Assigned</span>
        </div>
        {workers.map((w) => (
          <button
            key={w.ID}
            className="worker-row"
            onClick={() => setEditing(w)}
          >
            <span>
              <b>{w.Name}</b>
              <small>{w.Slug}</small>
            </span>
            <span>{w.Kind}</span>
            <span>
              <i
                className={
                  w.Enabled && w.ActiveSessionCount > 0 ? "enabled" : "disabled"
                }
              />
              {w.Enabled
                ? w.ActiveSessionCount > 0
                  ? "Active"
                  : "Offline"
                : "Disabled"}
              <small>
                {w.SessionCount} sessions · {w.MCPCalls} calls
              </small>
              <small>
                Last MCP:{" "}
                {w.LastActivityAt ? relativeTime(w.LastActivityAt) : "Never"}
              </small>
            </span>
            <span className="capabilities">
              {capabilities(w.Capabilities).join(", ") || "—"}
            </span>
            <span>{w.AssignedTaskCount}</span>
          </button>
        ))}
      </div>
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
      <div className="section-head">
        <div>
          <h2>Task properties</h2>
          <p>Trello-like fields with explicit Worker visibility.</p>
        </div>
        <button className="primary" onClick={() => setAdding(true)}>
          <Plus size={16} />
          Add property
        </button>
      </div>
      <div className="worker-table">
        <div className="property-row header">
          <span>Name</span>
          <span>Type</span>
          <span>Visibility</span>
          <span>Options</span>
        </div>
        {q.data?.map((p) => (
          <button
            key={p.ID}
            className="property-row"
            onClick={() => setEditing(p)}
          >
            <b>{p.Name}</b>
            <span>{p.Type}</span>
            <span>{p.Visibility.replaceAll("_", " ")}</span>
            <code>{p.Options}</code>
          </button>
        ))}
      </div>
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
          <input
            value={form.Name}
            onChange={(e) => setForm({ ...form, Name: e.target.value })}
          />
        </label>
        <label>
          Type
          <select
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
          </select>
        </label>
        <label className="wide">
          Visibility
          <select
            value={form.Visibility}
            onChange={(e) => setForm({ ...form, Visibility: e.target.value })}
          >
            <option value="human_only">Human only</option>
            <option value="agent_read">Agent read</option>
            <option value="agent_read_write">Agent read/write</option>
          </select>
        </label>
        {["select", "multi_select"].includes(form.Type) && (
          <label className="wide">
            Options JSON
            <textarea
              className="mono"
              value={form.Options}
              onChange={(e) => setForm({ ...form, Options: e.target.value })}
            />
          </label>
        )}
      </div>
      {property && (
        <button className="danger-link" onClick={() => remove.mutate()}>
          Delete property
        </button>
      )}
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
          <input
            value={form.Name}
            onChange={(e) => setForm({ ...form, Name: e.target.value })}
          />
        </label>
        <label>
          Slug
          <input
            value={form.Slug}
            onChange={(e) => setForm({ ...form, Slug: e.target.value })}
          />
        </label>
        <label className="wide">
          Description
          <textarea
            value={form.Description}
            onChange={(e) => setForm({ ...form, Description: e.target.value })}
          />
        </label>
        <label>
          Kind
          <select
            value={form.Kind}
            onChange={(e) => setForm({ ...form, Kind: e.target.value })}
          >
            <option value="external">external</option>
            <option value="custom">custom</option>
          </select>
        </label>
        <label className="check">
          <input
            type="checkbox"
            checked={form.Enabled}
            onChange={(e) => setForm({ ...form, Enabled: e.target.checked })}
          />
          Enabled
        </label>
        <label className="wide">
          Capability preset
          <select
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
          </select>
        </label>
        <label className="wide">
          Advanced: capabilities JSON
          <textarea
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
      {save.error && <p role="alert">{save.error.message}</p>}
      {archive.error && <p role="alert">{archive.error.message}</p>}
      {worker && (
        <WorkerDiagnostics worker={worker} projectPath={projectPath} />
      )}
      {worker && (
        <button className="danger-link" onClick={() => archive.mutate()}>
          Archive worker
        </button>
      )}
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
    <div className="modal-backdrop" onMouseDown={close}>
      <div className="modal" onMouseDown={(e) => e.stopPropagation()}>
        <div className="modal-head">
          <h2>{title}</h2>
          <button className="icon-button" onClick={close}>
            <X />
          </button>
        </div>
        {children}
      </div>
    </div>
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
    <div className="dialog-actions">
      <button onClick={close}>Cancel</button>
      <button className="primary" disabled={busy || disabled} onClick={save}>
        {busy ? "Saving…" : "Save"}
      </button>
    </div>
  );
}
function Info({ label, value }: { label: string; value: string }) {
  return (
    <div className="info">
      <small>{label}</small>
      <span>{value}</span>
    </div>
  );
}
function Empty({ title, subtitle }: { title: string; subtitle?: string }) {
  return (
    <div className="empty-screen">
      <span className="brand-mark">
        <Sparkles />
      </span>
      <h1>{title}</h1>
      {subtitle && <p>{subtitle}</p>}
    </div>
  );
}
export default App;
