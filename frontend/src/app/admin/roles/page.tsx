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
import { createRole, deleteRole, listRoles } from "@/lib/api/security";
import { formatJson } from "@/lib/utils";
import type { Role, RoleFormValues } from "@/types/api";

const initialValues: RoleFormValues = {
  name: "",
  privileges: "[\"*\"]",
};

export default function RolesAdminPage() {
  const session = useRequireSession();
  const [roles, setRoles] = useState<Role[]>([]);
  const [values, setValues] = useState(initialValues);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadRoles = useCallback(async () => {
    if (!session.token) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      setRoles(await listRoles(session.token));
    } catch (fetchError) {
      setError(fetchError instanceof Error ? fetchError.message : "Failed to load roles");
    } finally {
      setLoading(false);
    }
  }, [session.token]);

  useEffect(() => {
    if (session.isReady && session.user) {
      void loadRoles();
    }
  }, [loadRoles, session.isReady, session.user]);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!session.token) {
      return;
    }

    setSaving(true);
    setError(null);
    try {
      await createRole(
        {
          name: values.name.trim(),
          privileges: JSON.parse(values.privileges),
        },
        session.token,
      );
      setValues(initialValues);
      await loadRoles();
    } catch (saveError) {
      setError(saveError instanceof Error ? saveError.message : "Failed to create role");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (roleId?: string) => {
    if (!session.token || !roleId) {
      return;
    }

    try {
      await deleteRole(roleId, session.token);
      await loadRoles();
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "Failed to delete role");
    }
  };

  return (
    <AppShell>
      <div>
        <p className="subtle-label">Administration</p>
        <h1 className="mt-2 text-3xl font-semibold text-white">Roles</h1>
      </div>

      <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_360px]">
        <Card>
          {loading ? (
            <p className="text-sm text-slate-400">Loading roles...</p>
          ) : roles.length === 0 ? (
            <EmptyState
              title="No roles"
              description="Create the first RBAC role to define privileges."
            />
          ) : (
            <div className="space-y-3">
              {roles.map((entry) => (
                <div
                  key={entry.id || entry.name}
                  className="rounded-2xl border border-white/10 bg-white/[0.04] p-4"
                >
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h2 className="text-lg font-semibold text-white">{entry.name}</h2>
                      <p className="mt-1 text-xs text-slate-500">{entry.id}</p>
                      <pre className="mt-3 overflow-x-auto rounded-2xl bg-black/20 p-4 text-xs text-slate-300">
                        {formatJson(entry.privileges)}
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
          <h2 className="text-xl font-semibold text-white">Create role</h2>
          <form className="mt-6 space-y-4" onSubmit={onSubmit}>
            <Input
              placeholder="role name"
              value={values.name}
              onChange={(event) => setValues((current) => ({ ...current, name: event.target.value }))}
            />
            <Textarea
              rows={8}
              value={values.privileges}
              onChange={(event) =>
                setValues((current) => ({ ...current, privileges: event.target.value }))
              }
            />

            {error ? (
              <div className="rounded-2xl border border-rose-400/20 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
                {error}
              </div>
            ) : null}

            <Button type="submit" className="w-full" disabled={saving}>
              {saving ? "Creating..." : "Create role"}
            </Button>
          </form>
        </Card>
      </div>
    </AppShell>
  );
}
