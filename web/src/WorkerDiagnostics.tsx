import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "primereact/button";
import { InputTextarea } from "primereact/inputtextarea";
import { Panel } from "primereact/panel";
import { Tag } from "primereact/tag";
import { Message } from "primereact/message";
import { DataTable } from "primereact/datatable";
import { Column } from "primereact/column";
import { Divider } from "primereact/divider";

import { relativeTime } from "./workerPresentation";

type Session = {
  ID: string;
  StartedAt: string;
  EndedAt?: string;
  LastActivityAt?: string;
  LastSeenAt?: string;
  MCPCalls: number;
  ClientInfo: string;
};
type WorkerInfo = {
  ID: string;
  Slug: string;
  Enabled: boolean;
  ActiveSessionCount: number;
  SessionCount: number;
  MCPCalls: number;
  LastActivityAt?: string;
};
async function read<T>(path: string): Promise<T> {
  const response = await fetch(`/api${path}`);
  if (!response.ok)
    throw new Error(`Diagnostics unavailable (${response.status})`);
  return response.json();
}

export function WorkerDiagnostics({
  worker,
  projectPath,
}: {
  worker: WorkerInfo;
  projectPath: string;
}) {
  const [copyStatus, setCopyStatus] = useState("");
  const info = useQuery({
    queryKey: ["worker", worker.ID],
    queryFn: () => read<WorkerInfo>(`/workers/${worker.ID}`),
    refetchInterval: 4000,
  });
  const sessions = useQuery({
    queryKey: ["sessions", worker.ID],
    queryFn: () => read<Session[]>(`/workers/${worker.ID}/sessions`),
    refetchInterval: 4000,
  });
  const current = info.data || worker;
  const config = `[mcp_servers.agentboard]
command = "agentboard"
args = ["mcp", "--project", ${JSON.stringify(projectPath)}, "--worker", ${JSON.stringify(current.Slug)}]
`;
  return (
    <Panel className="worker-diagnostics" header="MCP diagnostics" toggleable>
      <div className="diagnostic-status">
        <Tag value={current.Enabled ? "Enabled" : "Disabled"} severity={current.Enabled ? "success" : "secondary"} />
        <Tag value={current.ActiveSessionCount > 0 ? "Active session" : "Offline"}
          severity={current.ActiveSessionCount > 0 ? "success" : "secondary"} />
      </div>
      <p>
        {current.SessionCount} sessions · {current.MCPCalls} calls · Last MCP:{" "}
        {current.LastActivityAt
          ? relativeTime(current.LastActivityAt)
          : "Never"}
      </p>
      <p>Project: {projectPath}</p>
      <Message severity="info" text="Open sessions are active while the MCP process sends a heartbeat. After 90 seconds without a heartbeat, they become offline." />
      {(info.error || sessions.error) && (
        <Message severity="error" text={(info.error || sessions.error)?.message} />
      )}
      <label>
        Codex MCP config
        <InputTextarea className="mono" readOnly value={config} rows={5} />
      </label>
      <Button
        type="button"
        onClick={async () => {
          try {
            await navigator.clipboard.writeText(config);
            setCopyStatus("Copied");
          } catch {
            setCopyStatus(
              "Clipboard unavailable. Select and copy the config above.",
            );
          }
        }}
      >
        Copy MCP config
      </Button>
      {copyStatus && <Message severity={copyStatus === "Copied" ? "success" : "warn"} text={copyStatus} />}
      <Divider />
      <h4>Recent sessions (up to 50)</h4>
      <DataTable value={sessions.data || []} dataKey="ID" scrollable scrollHeight="240px"
        className="sessions-table" loading={sessions.isLoading}
        emptyMessage="No MCP sessions yet. Configure your client, connect, then call get_my_board.">
        <Column header="Started" body={(session: Session) => new Date(session.StartedAt).toLocaleString()} />
        <Column header="Status" body={(session: Session) => {
          const active = !session.EndedAt && !!session.LastSeenAt && sessions.dataUpdatedAt - Date.parse(session.LastSeenAt) < 90000;
          return <Tag value={session.EndedAt ? "Ended" : active ? "Active" : "Offline"}
            severity={active ? "success" : "secondary"} />;
        }} />
        <Column header="Calls" body={(session: Session) => <span>{session.MCPCalls}<small>Last: {session.LastActivityAt ? relativeTime(session.LastActivityAt) : "Never"}</small></span>} />
        <Column field="ClientInfo" header="Client" />
      </DataTable>
    </Panel>
  );
}
