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
  if (!Number.isFinite(seconds)) return "unknown";
  if (seconds < 60) return `${seconds}s ago`;
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  return `${Math.floor(seconds / 86400)}d ago`;
}
