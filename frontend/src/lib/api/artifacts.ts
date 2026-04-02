import { download, toBackendUrl } from "@/lib/api/client";
import { encodeArtifactPath } from "@/lib/utils";

export async function uploadArtifact(
  repositoryName: string,
  artifactPath: string,
  file: File,
  token?: string | null,
) {
  const response = await fetch(
    toBackendUrl(`/repository/${encodeURIComponent(repositoryName)}/${encodeArtifactPath(artifactPath)}`),
    {
      method: "PUT",
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      body: file,
    },
  );

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Upload failed with status ${response.status}`);
  }

  return response.text();
}

export async function downloadArtifact(repositoryName: string, artifactPath: string, token?: string | null) {
  return download(
    `/repository/${encodeURIComponent(repositoryName)}/${encodeArtifactPath(artifactPath)}`,
    token,
  );
}

export async function removeArtifact(
  repositoryName: string,
  artifactPath: string,
  token?: string | null,
) {
  const response = await fetch(
    toBackendUrl(`/repository/${encodeURIComponent(repositoryName)}/${encodeArtifactPath(artifactPath)}`),
    {
      method: "DELETE",
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    },
  );

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Delete failed with status ${response.status}`);
  }
}
