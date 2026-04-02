import type { Artifact, BlobStore, CleanupPolicy, Repository } from "@/types/api";

function pick<T>(value: Record<string, unknown>, ...keys: string[]) {
  for (const key of keys) {
    if (key in value) {
      return value[key] as T;
    }
  }
  return undefined;
}

export function normalizeArtifact(value: unknown): Artifact {
  const raw = (value || {}) as Record<string, unknown>;
  return {
    id: pick<string>(raw, "id", "ID"),
    repositoryID: pick<string>(raw, "repositoryID", "RepositoryID"),
    path: pick<string>(raw, "path", "Path"),
    size: pick<number>(raw, "size", "Size"),
    contentType: pick<string>(raw, "contentType", "ContentType"),
    checksums: pick(raw, "checksums", "Checksums"),
    createdAt: pick<string>(raw, "createdAt", "CreatedAt"),
    lastDownloadedAt: pick<string | null>(raw, "lastDownloadedAt", "LastDownloadedAt") ?? null,
  };
}

export function normalizeCleanupPolicy(value: unknown): CleanupPolicy {
  const raw = (value || {}) as Record<string, unknown>;
  return {
    id: pick<string>(raw, "id", "ID"),
    name: pick<string>(raw, "name", "Name"),
    criteria: pick(raw, "criteria", "Criteria"),
    createdAt: pick<string>(raw, "createdAt", "CreatedAt"),
    updatedAt: pick<string>(raw, "updatedAt", "UpdatedAt"),
  };
}

export function normalizeBlobStore(value: unknown): BlobStore {
  const raw = (value || {}) as Record<string, unknown>;
  return {
    id: pick<string>(raw, "id", "ID"),
    name: pick<string>(raw, "name", "Name"),
    type: pick<string>(raw, "type", "Type"),
    attributes: pick(raw, "attributes", "Attributes"),
    createdAt: pick<string>(raw, "createdAt", "CreatedAt"),
    updatedAt: pick<string>(raw, "updatedAt", "UpdatedAt"),
  };
}

export function normalizeRepository(value: unknown): Repository {
  const raw = (value || {}) as Record<string, unknown>;
  const artifacts = (pick<unknown[]>(raw, "artifacts", "Artifacts") || []).map(normalizeArtifact);
  const cleanupPolicyRaw = pick(raw, "cleanupPolicy", "CleanupPolicy");
  const blobStoreRaw = pick(raw, "blobStore", "BlobStore");

  return {
    id: pick<string>(raw, "id", "ID"),
    name: pick<string>(raw, "name", "Name"),
    format: pick<string>(raw, "format", "Format"),
    type: pick<string>(raw, "type", "Type"),
    attributes: pick(raw, "attributes", "Attributes"),
    cleanupPolicyID: pick<string | null>(raw, "cleanupPolicyID", "CleanupPolicyID") ?? null,
    blobStoreID: pick<string | null>(raw, "blobStoreID", "BlobStoreID") ?? null,
    createdAt: pick<string>(raw, "createdAt", "CreatedAt"),
    updatedAt: pick<string>(raw, "updatedAt", "UpdatedAt"),
    artifacts,
    cleanupPolicy: cleanupPolicyRaw ? normalizeCleanupPolicy(cleanupPolicyRaw) : null,
    blobStore: blobStoreRaw ? normalizeBlobStore(blobStoreRaw) : null,
  };
}
