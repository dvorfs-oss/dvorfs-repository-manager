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
import { createCleanupPolicy, deleteCleanupPolicy, listCleanupPolicies } from "@/lib/api/cleanup";
import { formatJson } from "@/lib/utils";
import type { CleanupPolicy, CleanupPolicyFormValues } from "@/types/api";

const initialValues: CleanupPolicyFormValues = {
  name: "",
  criteria: "{\"maxAgeDays\":30,\"keepLastN\":10}",
};

export default function CleanupPoliciesPage() {
  const session = useRequireSession();
  const [policies, setPolicies] = useState<CleanupPolicy[]>([]);
  const [values, setValues] = useState(initialValues);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadPolicies = useCallback(async () => {
    if (!session.token) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      setPolicies(await listCleanupPolicies(session.token));
    } catch (fetchError) {
      setError(
        fetchError instanceof Error ? fetchError.message : "Failed to load cleanup policies",
      );
    } finally {
      setLoading(false);
    }
  }, [session.token]);

  useEffect(() => {
    if (session.isReady && session.user) {
      void loadPolicies();
    }
  }, [loadPolicies, session.isReady, session.user]);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!session.token) {
      return;
    }

    setSaving(true);
    setError(null);
    try {
      await createCleanupPolicy(
        {
          name: values.name.trim(),
          criteria: JSON.parse(values.criteria),
        },
        session.token,
      );
      setValues(initialValues);
      await loadPolicies();
    } catch (saveError) {
      setError(
        saveError instanceof Error ? saveError.message : "Failed to create cleanup policy",
      );
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id?: string) => {
    if (!session.token || !id) {
      return;
    }

    try {
      await deleteCleanupPolicy(id, session.token);
      await loadPolicies();
    } catch (deleteError) {
      setError(
        deleteError instanceof Error ? deleteError.message : "Failed to delete cleanup policy",
      );
    }
  };

  return (
    <AppShell>
      <div>
        <p className="subtle-label">Administration</p>
        <h1 className="mt-2 text-3xl font-semibold text-white">Cleanup policies</h1>
      </div>

      <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_360px]">
        <Card>
          {loading ? (
            <p className="text-sm text-slate-400">Loading cleanup policies...</p>
          ) : policies.length === 0 ? (
            <EmptyState
              title="No cleanup policies"
              description="Create retention rules for repository storage hygiene."
            />
          ) : (
            <div className="space-y-3">
              {policies.map((entry) => (
                <div
                  key={entry.id || entry.name}
                  className="rounded-2xl border border-white/10 bg-white/[0.04] p-4"
                >
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h2 className="text-lg font-semibold text-white">{entry.name}</h2>
                      <pre className="mt-3 overflow-x-auto rounded-2xl bg-black/20 p-4 text-xs text-slate-300">
                        {formatJson(entry.criteria)}
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
          <h2 className="text-xl font-semibold text-white">Create cleanup policy</h2>
          <form className="mt-6 space-y-4" onSubmit={onSubmit}>
            <Input
              placeholder="policy name"
              value={values.name}
              onChange={(event) => setValues((current) => ({ ...current, name: event.target.value }))}
            />
            <Textarea
              rows={8}
              value={values.criteria}
              onChange={(event) =>
                setValues((current) => ({ ...current, criteria: event.target.value }))
              }
            />

            {error ? (
              <div className="rounded-2xl border border-rose-400/20 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
                {error}
              </div>
            ) : null}

            <Button type="submit" className="w-full" disabled={saving}>
              {saving ? "Creating..." : "Create cleanup policy"}
            </Button>
          </form>
        </Card>
      </div>
    </AppShell>
  );
}
