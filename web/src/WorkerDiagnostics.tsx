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
import { harnessPresets, mcpConfiguration } from "./workerPresets";
import { t } from "./i18n";

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
  Harness: string;
};
type MCPCheck = { executable: string; projectPath: string; workerSlug: string; enabled: boolean; clientConnected: boolean; toolCount: number; serverOK: boolean; error?: string };
async function read<T>(path: string): Promise<T> {
  const response = await fetch(`/api${path}`);
  if (!response.ok)
    throw new Error(`Diagnostics unavailable (${response.status})`);
  return response.json();
}

export function WorkerDiagnostics({
  worker,
  projectPath,
  formatDate,
}: {
  worker: WorkerInfo;
  projectPath: string;
  formatDate: (value: string) => string;
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
  const check = useQuery({
    queryKey: ["mcp-check", worker.ID],
    queryFn: () => read<MCPCheck>(`/workers/${worker.ID}/mcp/check`),
    refetchOnWindowFocus: false,
  });
  const current = info.data || worker;
  const config = mcpConfiguration(current.Harness, projectPath, current.Slug, check.data?.executable || "agentboard");
  return (
    <Panel className="worker-diagnostics" header={t("MCP diagnostics")} toggleable>
      <div className="diagnostic-status">
        <Tag value={t(current.Enabled ? "Enabled" : "Disabled")} severity={current.Enabled ? "success" : "secondary"} />
        <Tag value={t(current.ActiveSessionCount > 0 ? "Active session" : "Offline")}
          severity={current.ActiveSessionCount > 0 ? "success" : "secondary"} />
      </div>
      <p>
        {current.SessionCount} {t("sessions")} · {current.MCPCalls} {t("calls")} · {t("Last MCP:")}{" "}
        {current.LastActivityAt
          ? relativeTime(current.LastActivityAt)
          : t("Never")}
      </p>
      <p>{t("Project:")} {projectPath}</p>
      <Message severity="info" text={t("Open sessions are active while the MCP process sends a heartbeat. After 90 seconds without a heartbeat, they become offline.")} />
      {(info.error || sessions.error) && (
        <Message severity="error" text={(info.error || sessions.error)?.message} />
      )}
      <div className="mcp-check">
        <strong>{t("Server check")}</strong>
        <Tag value={check.isLoading ? t("Checking…") : check.data?.serverOK ? `${check.data.toolCount} ${t("MCP tools available")}` : t("Unavailable")} severity={check.data?.serverOK ? "success" : "warning"} />
        <Tag value={t(check.data?.clientConnected ? "Client connected" : "No client connected")} severity={check.data?.clientConnected ? "success" : "secondary"} />
        <Button text label={t("Check again")} onClick={() => { info.refetch(); sessions.refetch(); check.refetch(); }} />
      </div>
      {check.data?.error && <Message severity="error" text={check.data.error} />}
      <p>{t("Server check verifies the MCP handshake and tool discovery. Client connected means a configured AI client has an active session. After adding this configuration, open the client and call")} <code>get_my_board</code>.</p>
      {current.Harness === "ollama" && <Message severity="info" text={t("Ollama provides the model. Configure OpenCode to use Ollama and add this MCP server to OpenCode; Ollama alone does not run MCP clients.")} />}
      <label>
        {harnessPresets.find((item) => item.id === current.Harness)?.name || "MCP"} {t("config")} · {config.location}
        <InputTextarea className="mono" readOnly value={config.text} rows={8} />
      </label>
      <Button
        type="button"
        onClick={async () => {
          try {
            await navigator.clipboard.writeText(config.text);
            setCopyStatus(t("Copied"));
          } catch {
            setCopyStatus(
              t("Clipboard unavailable. Select and copy the config above."),
            );
          }
        }}
      >
        {t("Copy MCP config")}
      </Button>
      {copyStatus && <Message severity="success" text={copyStatus} />}
      <Divider />
      <h4>{t("Recent sessions (up to 50)")}</h4>
      <DataTable value={sessions.data || []} dataKey="ID" scrollable scrollHeight="240px"
        className="sessions-table" loading={sessions.isLoading}
        emptyMessage={t("No MCP sessions yet. Configure your client, connect, then call get_my_board.")}>
        <Column header={t("Started")} body={(session: Session) => formatDate(session.StartedAt)} />
        <Column header={t("Status")} body={(session: Session) => {
          const active = !session.EndedAt && !!session.LastSeenAt && sessions.dataUpdatedAt - Date.parse(session.LastSeenAt) < 90000;
          return <Tag value={t(session.EndedAt ? "Ended" : active ? "Active" : "Offline")}
            severity={active ? "success" : "secondary"} />;
        }} />
        <Column header={t("Calls")} body={(session: Session) => <span>{session.MCPCalls}<small>{t("Last:")} {session.LastActivityAt ? relativeTime(session.LastActivityAt) : t("Never")}</small></span>} />
        <Column field="ClientInfo" header={t("Client")} />
      </DataTable>
    </Panel>
  );
}
