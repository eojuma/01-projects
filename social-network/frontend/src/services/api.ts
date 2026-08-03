const API_BASE = "";

async function apiFetch(path: string, options: RequestInit = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: { "Content-Type": "application/json", ...(options.headers || {}) },
    ...options,
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  // No body on 204
  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  get: (path: string) => apiFetch(path),
  post: (path: string, body?: unknown) =>
    apiFetch(path, { method: "POST", body: body != null ? JSON.stringify(body) : undefined }),
  patch: (path: string, body?: unknown) =>
    apiFetch(path, { method: "PATCH", body: body != null ? JSON.stringify(body) : undefined }),
  delete: (path: string) => apiFetch(path, { method: "DELETE" }),
};

export default api;
