import { request } from "@/lib/api/client";
import { normalizeCleanupPolicy } from "@/lib/api/normalize";
import type { CleanupPolicy } from "@/types/api";

export async function listCleanupPolicies(token?: string | null) {
  const response = await request<unknown[]>("/cleanup-policies", {
    method: "GET",
    authToken: token,
  });
  return response.map(normalizeCleanupPolicy);
}

export async function createCleanupPolicy(
  payload: { name: string; criteria: unknown },
  token?: string | null,
) {
  const response = await request<unknown>("/cleanup-policies", {
    method: "POST",
    body: payload,
    authToken: token,
  });
  return normalizeCleanupPolicy(response);
}

export async function deleteCleanupPolicy(id: string, token?: string | null) {
  return request<void>(`/cleanup-policies/${encodeURIComponent(id)}`, {
    method: "DELETE",
    authToken: token,
  });
}
