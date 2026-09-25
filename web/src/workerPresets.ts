import { capabilityPresets } from "./workerPresentation";

export const harnessPresets = [
  { id: "codex", name: "Codex", capabilities: "Coding", description: "OpenAI Codex CLI / desktop", docs: "https://developers.openai.com/learn/docs-mcp" },
  { id: "claude", name: "Claude Code", capabilities: "Coding", description: "Anthropic Claude Code", docs: "https://docs.anthropic.com/en/docs/claude-code/mcp" },
  { id: "gemini", name: "Gemini CLI", capabilities: "Coding", description: "Google Gemini CLI", docs: "https://geminicli.com/docs/tools/mcp-server/" },
  { id: "cursor", name: "Cursor", capabilities: "Coding", description: "Cursor editor / CLI", docs: "https://prod.cursor.com/docs/cli/mcp" },
  { id: "opencode", name: "OpenCode", capabilities: "Coding", description: "OpenCode terminal agent", docs: "https://opencode.ai/docs/mcp-servers/" },
  { id: "ollama", name: "Ollama via OpenCode", capabilities: "Coding", description: "Ollama supplies the model; OpenCode runs the MCP client", docs: "https://opencode.ai/docs/mcp-servers/" },
  { id: "vscode", name: "VS Code Copilot", capabilities: "Coding", description: "VS Code Agent mode", docs: "https://code.visualstudio.com/docs/agent-customization/mcp-servers" },
  { id: "custom", name: "Other MCP client", capabilities: "Generic", description: "Any client that starts a local stdio MCP server", docs: "https://modelcontextprotocol.io/docs/develop/connect-local-servers" },
] as const;
export type HarnessID = typeof harnessPresets[number]["id"];

export function presetCapabilities(harness: string) {
  const capability = harnessPresets.find((item) => item.id === harness)?.capabilities || "Generic";
  return JSON.stringify(capabilityPresets[capability], null, 2);
}

export function mcpConfiguration(harness: string, projectPath: string, slug: string, executable = "agentboard") {
  const args = ["mcp", "--project", projectPath, "--worker", slug];
  const server = { command: executable, args };
  if (harness === "codex") return {
    location: "~/.codex/config.toml", text: `[mcp_servers.agentboard_${slug.replace(/[^a-zA-Z0-9_-]/g, "_")}]\ncommand = ${JSON.stringify(executable)}\nargs = [${args.map((arg) => JSON.stringify(arg)).join(", ")}]\n`,
  };
  if (harness === "opencode" || harness === "ollama") return {
    location: "opencode.json", text: JSON.stringify({ $schema: "https://opencode.ai/config.json", mcp: { [`agentboard_${slug}`]: { type: "local", command: [executable, ...args], enabled: true } } }, null, 2),
  };
  if (harness === "vscode") return { location: ".vscode/mcp.json", text: JSON.stringify({ servers: { [`agentboard_${slug}`]: server } }, null, 2) };
  if (harness === "gemini") return { location: "~/.gemini/settings.json (merge mcpServers)", text: JSON.stringify({ mcpServers: { [`agentboard_${slug}`]: server } }, null, 2) };
  if (harness === "cursor") return { location: ".cursor/mcp.json", text: JSON.stringify({ mcpServers: { [`agentboard_${slug}`]: server } }, null, 2) };
  return { location: harness === "claude" ? ".mcp.json (Claude Code project)" : "Client MCP configuration", text: JSON.stringify({ mcpServers: { [`agentboard_${slug}`]: { type: "stdio", ...server } } }, null, 2) };
}
