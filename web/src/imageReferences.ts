export type ImageReference = { url: string; label: string };

export function isManagedImageURL(url: string | undefined): url is string {
  return !!url && /^\/api\/images\/[A-Za-z0-9_-]+$/.test(url);
}

export function isManagedFileURL(url: string | undefined): url is string {
  return !!url && /^\/api\/files\/[A-Za-z0-9_-]+$/.test(url);
}

// Ignore examples in code when finding the attachments belonging to a message.
function withoutCode(value: string) {
  return value.replace(/(^|\n)( {0,3})(`{3,}|~{3,})[^\n]*\n[\s\S]*?(?:\n\2\3[^\n]*(?=\n|$)|$)/g,
    (code) => " ".repeat(code.length)).replace(/(`+)[\s\S]*?\1/g, (code) => " ".repeat(code.length));
}

function imageLinkPattern(kind = "images") {
  return new RegExp(String.raw`!?\[((?:\\.|[^\]\\])*)\]\(\s*<?(/api/(?:${kind})/[A-Za-z0-9_-]+)>?(?:\s+["'][^"']*["'])?\s*\)`, "g");
}

export function imageReferences(value: string, kind = "images"): ImageReference[] {
  const found = new Map<string, ImageReference>();
  for (const match of withoutCode(value).matchAll(imageLinkPattern(kind))) {
    if (!found.has(match[2])) {
      found.set(match[2], { url: match[2], label: match[1].replace(/\\([[\]\\])/g, "$1") || "🖼" });
    }
  }
  return [...found.values()];
}

export function fileReferences(value: string): ImageReference[] {
  return imageReferences(value, "files");
}

export function imageReferenceSummary(value: string): string {
  for (const match of [...withoutCode(value).matchAll(imageLinkPattern("images|files"))].reverse()) {
    const label = match[1].replace(/\\([[\]\\])/g, "$1");
    const icon = isManagedFileURL(match[2]) ? "📄" : "🖼";
    const text = label.startsWith(icon) ? label : `${icon}${label ? ` ${label}` : ""}`;
    value = value.slice(0, match.index) + text + value.slice(match.index! + match[0].length);
  }
  return value;
}

export function nextImageNumber(value: string): number {
  const images = imageReferences(value);
  return images.reduce((largest, image) => {
    const number = Number(image.label.match(/(\d+)\s*$/)?.[1] || 0);
    return Number.isSafeInteger(number) && number < Number.MAX_SAFE_INTEGER ? Math.max(largest, number) : largest;
  }, images.length) + 1;
}

export function removeImageReference(value: string, url: string): string {
  const matches = [...withoutCode(value).matchAll(imageLinkPattern("images|files"))].filter((match) => match[2] === url);
  for (const match of matches.reverse()) {
    value = value.slice(0, match.index) + value.slice(match.index! + match[0].length);
  }
  return value;
}

export async function uploadImage(projectID: string, file: File, signal: AbortSignal, kind = "images"): Promise<string> {
  const form = new FormData();
  form.append("file", file);
  const response = await fetch(`/api/projects/${encodeURIComponent(projectID)}/${kind}`, { method: "POST", body: form, signal });
  const result = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(result.error || response.statusText);
  if (!(kind === "images" ? isManagedImageURL(result.URL) : isManagedFileURL(result.URL))) throw new Error("Invalid attachment upload response");
  return result.URL;
}
