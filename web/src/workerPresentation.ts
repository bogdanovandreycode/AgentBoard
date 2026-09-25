export const capabilityPresets: Record<string, Record<string, boolean>> = {
  Generic: {
    coding: false,
    testing: false,
    git: false,
    backend: false,
    frontend: false,
    image_generation: false,
    comfyui: false,
  },
  Coding: {
    coding: true,
    testing: true,
    git: true,
    backend: true,
    frontend: true,
    image_generation: false,
    comfyui: false,
  },
  Graphics: {
    coding: false,
    testing: false,
    git: false,
    backend: false,
    frontend: false,
    image_generation: true,
    comfyui: true,
  },
  Reviewer: {
    coding: false,
    testing: true,
    git: true,
    backend: false,
    frontend: false,
    image_generation: false,
    comfyui: false,
  },
};

export function relativeTime(value: string) {
  const seconds = Math.max(
    0,
    Math.floor((Date.now() - Date.parse(value)) / 1000),
  );
  if (!Number.isFinite(seconds)) return "—";
  const locale = document.documentElement.lang || "en";
  const format = new Intl.RelativeTimeFormat(locale, { numeric: "auto", style: "short" });
  if (seconds < 60) return format.format(-seconds, "second");
  if (seconds < 3600) return format.format(-Math.floor(seconds / 60), "minute");
  if (seconds < 86400) return format.format(-Math.floor(seconds / 3600), "hour");
  return format.format(-Math.floor(seconds / 86400), "day");
}
