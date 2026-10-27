import { useQuery } from "@tanstack/react-query";
import { api } from "./api";

export type GlobalSettings = { draftAutosaveEnabled: boolean; draftAutosaveSeconds: number; draftSaveOnBlur: boolean; checkProjectContext: boolean };
export function useGlobalSettings() {
  return useQuery({ queryKey: ["global-settings"], queryFn: () => api<GlobalSettings>("/settings/global"), refetchInterval: 5000 });
}
