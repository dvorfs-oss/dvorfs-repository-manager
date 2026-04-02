"use client";

import { useParams, useRouter } from "next/navigation";
import type { FormEvent } from "react";
import { useCallback, useEffect, useMemo, useState } from "react";

import { AppShell } from "@/components/app-shell";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { Input } from "@/components/ui/input";
import { useRequireSession } from "@/hooks/use-require-session";
import { downloadArtifact, removeArtifact, uploadArtifact } from "@/lib/api/artifacts";
import { deleteRepository, getRepository } from "@/lib/api/repositories";
import { formatDate, formatJson } from "@/lib/utils";
import type { Artifact, Repository } from "@/types/api";

export default function RepositoryDetailsPage() {
  const router = useRouter();
  const params = useParams<{ name: string }>();
  const repositoryName = useMemo(() => {
    const rawName = Array.isArray(params.name) ? params.name[0] : params.name;
    return decodeURIComponent(rawName || "");
  }, [params.name]);
  const session = useRequireSession();
  const [repository, setRepository] = useState<Repository | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [artifactPath, setArtifactPath] = useState("");
  const [artifactFile, setArtifactFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);

  const loadRepository = useCallback(async () => {
    if (!session.token) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      setRepository(await getRepository(repositoryName, session.token));
    } catch (fetchError) {
      setError(fetchError instanceof Error ? fetchError.message : "Failed to load repository");
    } finally {
      setLoading(false);
    }
  }, [repositoryName, session.token]);

  useEffect(() => {
    if (session.isReady && session.user) {
      void loadRepository();
    }
  }, [loadRepository, session.isReady, session.user]);

  const artifacts = repository?.artifacts || [];

  const handleUpload = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!artifactFile) {
      setError("Choose a file to upload");
      return;
    }

    const nextPath = artifactPath.trim() || artifactFile.name;
    setUploading(true);
    setError(null);

    try {
      await uploadArtifact(repositoryName, nextPath, artifactFile, session.token);
      setArtifactFile(null);
      setArtifactPath("");
      await loadRepository();
    } catch (uploadError) {
      setError(uploadError instanceof Error ? uploadError.message : "Upload failed");
    } finally {
      setUploading(false);
    }
  };

  const handleDownload = async (artifact: Artifact) => {
    if (!artifact.path) {
      return;
    }

    const blob = await downloadArtifact(repositoryName, artifact.path, session.token);
    const url = window.URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = artifact.path.split("/").pop() || artifact.path;
    anchor.click();
    window.URL.revokeObjectURL(url);
  };

  const handleDeleteArtifact = async (artifact: Artifact) => {
    if (!artifact.path) {
      return;
    }

    try {
      await removeArtifact(repositoryName, artifact.path, session.token);
      await loadRepository();
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "Artifact deletion failed");
    }
  };

  const handleDeleteRepository = async () => {
    try {
      await deleteRepository(repositoryName, session.token);
      router.push("/repositories");
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "Repository deletion failed");
    }
  };

  return (
    <AppShell>
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p className="subtle-label">Repository details</p>
          <h1 className="mt-2 text-3xl font-semibold text-white">{repositoryName}</h1>
          <p className="mt-2 text-sm text-slate-400">
            Browse metadata, upload artifacts, and download stored assets from the repository.
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Badge>{repository?.format || "unknown format"}</Badge>
          <Badge>{repository?.type || "unknown type"}</Badge>
          <Button variant="danger" onClick={handleDeleteRepository}>
            Delete repository
          </Button>
        </div>
      </div>

      {loading ? (
        <Card className="mt-6">
          <p className="text-sm text-slate-400">Loading repository...</p>
        </Card>
      ) : error && !repository ? (
        <Card className="mt-6 border-rose-400/20 bg-rose-500/10">
          <p className="text-sm text-rose-100">{error}</p>
          <div className="mt-4">
            <Button variant="secondary" onClick={loadRepository}>
              Retry
            </Button>
          </div>
        </Card>
      ) : (
        <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_380px]">
          <div className="space-y-6">
            <Card>
              <div className="grid gap-4 md:grid-cols-2">
                <div>
                  <p className="text-sm text-slate-400">Repository ID</p>
                  <p className="mt-1 break-all text-sm text-slate-100">{repository?.id || "N/A"}</p>
                </div>
                <div>
                  <p className="text-sm text-slate-400">Updated</p>
                  <p className="mt-1 text-sm text-slate-100">{formatDate(repository?.updatedAt)}</p>
                </div>
                <div>
                  <p className="text-sm text-slate-400">Blob store</p>
                  <p className="mt-1 text-sm text-slate-100">{repository?.blobStoreID || "None"}</p>
                </div>
                <div>
                  <p className="text-sm text-slate-400">Cleanup policy</p>
                  <p className="mt-1 text-sm text-slate-100">
                    {repository?.cleanupPolicyID || "None"}
                  </p>
                </div>
              </div>

              <div className="mt-6">
                <p className="text-sm text-slate-400">Attributes</p>
                <pre className="mt-2 overflow-x-auto rounded-2xl bg-black/20 p-4 text-xs text-slate-200">
                  {formatJson(repository?.attributes)}
                </pre>
              </div>
            </Card>

            <Card>
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-sm text-slate-400">Artifacts</p>
                  <h2 className="text-lg font-semibold text-white">{artifacts.length} stored item(s)</h2>
                </div>
              </div>

              <div className="mt-4">
                {artifacts.length === 0 ? (
                  <EmptyState
                    title="No artifacts yet"
                    description="Upload the first file to populate this repository."
                  />
                ) : (
                  <div className="space-y-3">
                    {artifacts.map((artifact) => (
                      <div
                        key={artifact.id || artifact.path}
                        className="rounded-2xl border border-white/10 bg-white/[0.04] p-4"
                      >
                        <div className="flex flex-wrap items-center justify-between gap-3">
                          <div>
                            <p className="font-medium text-white">{artifact.path}</p>
                            <p className="mt-1 text-xs text-slate-400">
                              {artifact.contentType || "application/octet-stream"} ·{" "}
                              {artifact.size ?? 0} bytes · {formatDate(artifact.createdAt)}
                            </p>
                          </div>
                          <Button variant="secondary" onClick={() => handleDownload(artifact)}>
                            Download
                          </Button>
                          <Button variant="danger" onClick={() => handleDeleteArtifact(artifact)}>
                            Delete
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </Card>
          </div>

          <Card>
            <div className="space-y-2">
              <p className="subtle-label">Upload</p>
              <h2 className="text-xl font-semibold text-white">Add artifact</h2>
              <p className="text-sm text-slate-400">
                The file body is sent directly to the repository upload endpoint.
              </p>
            </div>

            <form className="mt-6 space-y-4" onSubmit={handleUpload}>
              <label className="block space-y-2">
                <span className="text-sm font-medium text-slate-200">Artifact path</span>
                <Input
                  placeholder="path/in/repo/file.jar"
                  value={artifactPath}
                  onChange={(event) => setArtifactPath(event.target.value)}
                />
              </label>

              <label className="block space-y-2">
                <span className="text-sm font-medium text-slate-200">File</span>
                <Input
                  type="file"
                  onChange={(event) => setArtifactFile(event.target.files?.[0] ?? null)}
                />
              </label>

              {artifactFile ? (
                <p className="text-xs text-slate-400">Selected: {artifactFile.name}</p>
              ) : null}

              {error ? (
                <div className="rounded-2xl border border-rose-400/20 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
                  {error}
                </div>
              ) : null}

              <Button type="submit" className="w-full" disabled={uploading}>
                {uploading ? "Uploading..." : "Upload artifact"}
              </Button>
            </form>
          </Card>
        </div>
      )}
    </AppShell>
  );
}
