const backendBase = (process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080").replace(
  /\/$/,
  "",
);
const apiBase = `${backendBase}/api/v1`;

type RequestOptions = Omit<RequestInit, "body"> & {
  body?: BodyInit | object | null;
  authToken?: string | null;
};

async function readResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get("content-type") ?? "";

  if (response.status === 204) {
    return undefined as T;
  }

  if (contentType.includes("application/json")) {
    return (await response.json()) as T;
  }

  const text = await response.text();
  if (!text.trim()) {
    return undefined as T;
  }

  return text as T;
}

export async function request<T>(path: string, options: RequestOptions = {}) {
  const headers = new Headers(options.headers);

  if (options.authToken) {
    headers.set("Authorization", `Bearer ${options.authToken}`);
  }

  let body = options.body;
  if (body && !(body instanceof FormData) && !(body instanceof Blob) && typeof body === "object") {
    headers.set("Content-Type", "application/json");
    body = JSON.stringify(body);
  }

  const response = await fetch(`${apiBase}${path}`, {
    ...options,
    headers,
    body,
    cache: "no-store",
  });

  if (!response.ok) {
    const message = await readResponse<string | { message?: string; error?: string }>(response);
    if (typeof message === "string") {
      throw new Error(message || `Request failed with status ${response.status}`);
    }

    throw new Error(message?.error || message?.message || `Request failed with status ${response.status}`);
  }

  return readResponse<T>(response);
}

export async function download(path: string, authToken?: string | null) {
  const response = await fetch(`${backendBase}${path}`, {
    headers: authToken ? { Authorization: `Bearer ${authToken}` } : undefined,
    cache: "no-store",
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Download failed with status ${response.status}`);
  }

  return response.blob();
}

export function toBackendUrl(path: string) {
  return `${backendBase}${path}`;
}
