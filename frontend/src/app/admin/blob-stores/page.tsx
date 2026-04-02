"use client";

import type { FormEvent } from "react";
import { useCallback, useEffect, useState } from "react";

import { AppShell } from "@/components/app-shell";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useRequireSession } from "@/hooks/use-require-session";
import { createBlobStore, deleteBlobStore, listBlobStores } from "@/lib/api/blob-stores";
import { formatJson } from "@/lib/utils";
import type { BlobStore, BlobStoreFormValues } from "@/types/api";

const initialValues: BlobStoreFormValues = {
  name: "",
  type: "file",
  attributes: "{}",
};

export default function BlobStoresAdminPage() {
  const session = useRequireSession();
  const [blobStores, setBlobStores] = useState<BlobStore[]>([]);
  const [values, setValues] = useState(initialValues);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadBlobStores = useCallback(async () => {
    if (!session.token) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      setBlobStores(await listBlobStores(session.token));
    } catch (fetchError) {
      setError(fetchError instanceof Error ? fetchError.message : "Failed to load blob stores");
    } finally {
      setLoading(false);
    }
  }, [session.token]);

  useEffect(() => {
    if (session.isReady && session.user) {
      void loadBlobStores();
    }
  }, [loadBlobStores, session.isReady, session.user]);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!session.token) {
      return;
    }

    setSaving(true);
    setError(null);
    try {
      await createBlobStore(
        {
          name: values.name.trim(),
          type: values.type.trim(),
          attributes: JSON.parse(values.attributes),
        },
        session.token,
      );
      setValues(initialValues);
      await loadBlobStores();
    } catch (saveError) {
      setError(saveError instanceof Error ? saveError.message : "Failed to create blob store");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id?: string) => {
    if (!session.token || !id) {
      return;
    }

    try {
      await deleteBlobStore(id, session.token);
      await loadBlobStores();
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "Failed to delete blob store");
    }
  };

  return (
    <AppShell>
      <div>
        <p className="subtle-label">Administration</p>
        <h1 className="mt-2 text-3xl font-semibold text-white">Blob stores</h1>
      </div>

      <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_360px]">
        <Card>
          {loading ? (
            <p className="text-sm text-slate-400">Loading blob stores...</p>
          ) : blobStores.length === 0 ? (
            <EmptyState
              title="No blob stores"
              description="Create a file-backed blob store to attach repositories."
            />
          ) : (
            <div className="space-y-3">
              {blobStores.map((entry) => (
                <div
                  key={entry.id || entry.name}
                  className="rounded-2xl border border-white/10 bg-white/[0.04] p-4"
                >
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h2 className="text-lg font-semibold text-white">{entry.name}</h2>
                      <p className="mt-1 text-sm text-slate-400">{entry.type || "file"}</p>
                      <pre className="mt-3 overflow-x-auto rounded-2xl bg-black/20 p-4 text-xs text-slate-300">
                        {formatJson(entry.attributes)}
                      </pre>
                    </div>
                    <Button variant="danger" onClick={() => handleDelete(entry.id)}>
                      Delete
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>

        <Card>
          <h2 className="text-xl font-semibold text-white">Create blob store</h2>
          <form className="mt-6 space-y-4" onSubmit={onSubmit}>
            <Input
              placeholder="blob store name"
              value={values.name}
              onChange={(event) => setValues((current) => ({ ...current, name: event.target.value }))}
            />
            <Input
              placeholder="file"
              value={values.type}
              onChange={(event) => setValues((current) => ({ ...current, type: event.target.value }))}
            />
            <Textarea
              rows={8}
              value={values.attributes}
              onChange={(event) =>
                setValues((current) => ({ ...current, attributes: event.target.value }))
              }
            />

            {error ? (
              <div className="rounded-2xl border border-rose-400/20 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
                {error}
              </div>
            ) : null}

            <Button type="submit" className="w-full" disabled={saving}>
              {saving ? "Creating..." : "Create blob store"}
            </Button>
          </form>
        </Card>
      </div>
    </AppShell>
  );
}
