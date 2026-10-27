export async function api<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch("/api" + url, { headers: { "Content-Type": "application/json" }, ...options });
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(error.error || response.statusText);
  }
  return response.status === 204 ? undefined as T : response.json();
}

export async function withVisibleLoading<T>(request: Promise<T>): Promise<T> {
  const [value] = await Promise.all([request, new Promise<void>((resolve) => window.setTimeout(resolve, 350))]);
  return value;
}
