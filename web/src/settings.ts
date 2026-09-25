export type BoardColumn = { id: string; name: string };
export type ProjectSettings = {
  columns: BoardColumn[];
  timezone: string;
  language: string;
  theme: "dark" | "light" | "black" | "ubuntu" | "windows";
  boardRefreshSeconds: number;
  workerRefreshSeconds: number;
};

export const builtInColumns: BoardColumn[] = [
  { id: "backlog", name: "Backlog" },
  { id: "features", name: "Features" },
  { id: "in_progress", name: "In progress" },
  { id: "testing", name: "Testing" },
  { id: "verification", name: "Verification" },
  { id: "complete", name: "Complete" },
];
export const lockedColumns = new Set(["features", "in_progress", "testing", "verification"]);

export function formatDate(value: string, timezone = "local", language = "system") {
  try {
    return new Intl.DateTimeFormat(language === "system" ? undefined : language, {
      dateStyle: "medium", timeStyle: "short", ...(timezone === "local" ? {} : { timeZone: timezone }),
    }).format(new Date(value));
  } catch {
    return new Date(value).toLocaleString();
  }
}
