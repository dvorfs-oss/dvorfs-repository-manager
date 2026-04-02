import { request } from "@/lib/api/client";
import { normalizeArtifact, normalizeRepository } from "@/lib/api/normalize";
import type { Artifact, Repository, RepositoryFormValues } from "@/types/api";

export function parseRepositoryAttributes(values: RepositoryFormValues) {
  let attributes: unknown = {};

  if (values.attributes.trim()) {
    attributes = JSON.parse(values.attributes);
  }

  return {
    name: values.name.trim(),
    format: values.format,
    type: values.type,
    attributes,
    blobStoreID: values.blobStoreID.trim() || null,
    cleanupPolicyID: values.cleanupPolicyID.trim() || null,
  };
}

export async function listRepositories(token?: string | null) {
  const response = await request<unknown[]>("/repositories", {
    method: "GET",
    authToken: token,
  });
  return response.map(normalizeRepository);
}

export async function getRepository(name: string, token?: string | null) {
  const response = await request<unknown>(`/repositories/${encodeURIComponent(name)}`, {
    method: "GET",
    authToken: token,
  });
  return normalizeRepository(response);
}

export async function createRepository(payload: unknown, token?: string | null) {
  const response = await request<unknown | void>("/repositories", {
    method: "POST",
    body: payload as Record<string, unknown>,
    authToken: token,
  });
  return response ? normalizeRepository(response) : response;
}

export async function updateRepository(name: string, payload: unknown, token?: string | null) {
  return request<void>(`/repositories/${encodeURIComponent(name)}`, {
    method: "PUT",
    body: payload as Record<string, unknown>,
    authToken: token,
  });
}

export async function deleteRepository(name: string, token?: string | null) {
  return request<void>(`/repositories/${encodeURIComponent(name)}`, {
    method: "DELETE",
    authToken: token,
  });
}

export async function searchArtifacts(query: string, token?: string | null) {
  const encoded = new URLSearchParams({ q: query });
  const response = await request<unknown[]>(`/search/artifacts?${encoded.toString()}`, {
    method: "GET",
    authToken: token,
  });
  return response.map(normalizeArtifact);
}
