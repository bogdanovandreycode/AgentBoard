import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

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
    <section className="worker-diagnostics">
      <h3>MCP diagnostics</h3>
      <p>
        {current.Enabled ? "Worker enabled" : "Worker disabled"} ·{" "}
        {current.ActiveSessionCount > 0 ? "Active session" : "Offline"}
      </p>
      <p>
        {current.SessionCount} sessions · {current.MCPCalls} calls · Last MCP:{" "}
        {current.LastActivityAt
          ? relativeTime(current.LastActivityAt)
          : "Never"}
      </p>
      <p>Project: {projectPath}</p>
      <p>
        Open sessions are active while the MCP process sends a heartbeat. After
        90 seconds without a heartbeat, they become offline.
      </p>
      {(info.error || sessions.error) && (
        <p role="alert">{(info.error || sessions.error)?.message}</p>
      )}
      <label>
        Codex MCP config
        <textarea className="mono" readOnly value={config} rows={5} />
      </label>
      <button
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
      </button>
      <span role="status">{copyStatus}</span>
      <h4>Recent sessions (up to 50)</h4>
      {sessions.isLoading && <p>Loading sessions…</p>}
      {sessions.data?.length === 0 && (
        <p>
          No MCP sessions yet. Configure your client, connect, then call
          get_my_board.
        </p>
      )}
      <div className="session-list">
        {sessions.data?.map((session) => (
          <div key={session.ID}>
            <strong>{new Date(session.StartedAt).toLocaleString()}</strong>
            <span>
              {session.EndedAt
                ? `Ended ${new Date(session.EndedAt).toLocaleString()}`
                : session.LastSeenAt &&
                    sessions.dataUpdatedAt - Date.parse(session.LastSeenAt) <
                      90000
                  ? "Active"
                  : "Offline (no recent heartbeat)"}
            </span>
            <span>
              {session.MCPCalls} calls · Last MCP:{" "}
              {session.LastActivityAt
                ? relativeTime(session.LastActivityAt)
                : "Never"}
            </span>
            <small>{session.ClientInfo}</small>
          </div>
        ))}
      </div>
    </section>
  );
}
