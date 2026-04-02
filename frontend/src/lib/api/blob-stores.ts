import { request } from "@/lib/api/client";
import { normalizeBlobStore } from "@/lib/api/normalize";
import type { BlobStore } from "@/types/api";

export async function listBlobStores(token?: string | null) {
  const response = await request<unknown[]>("/blob-stores", {
    method: "GET",
    authToken: token,
  });
  return response.map(normalizeBlobStore);
}

export async function createBlobStore(
  payload: { name: string; type: string; attributes: unknown },
  token?: string | null,
) {
  const response = await request<unknown>("/blob-stores", {
    method: "POST",
    body: payload,
    authToken: token,
  });
  return normalizeBlobStore(response);
}

export async function deleteBlobStore(id: string, token?: string | null) {
  return request<void>(`/blob-stores/${encodeURIComponent(id)}`, {
    method: "DELETE",
    authToken: token,
  });
}
