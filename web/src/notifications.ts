import { createContext, useContext } from "react";

export type Notify = (severity: "success" | "error", summary: string, detail?: string) => void;
export const NotificationContext = createContext<Notify>(() => {});
export const useNotify = () => useContext(NotificationContext);
