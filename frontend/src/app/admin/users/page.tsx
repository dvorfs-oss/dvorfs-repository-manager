"use client";

import type { FormEvent } from "react";
import { useCallback, useEffect, useState } from "react";

import { AppShell } from "@/components/app-shell";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { Input } from "@/components/ui/input";
import { useRequireSession } from "@/hooks/use-require-session";
import { createUser, deleteUser, listRoles, listUsers } from "@/lib/api/security";
import type { Role, User, UserFormValues } from "@/types/api";

const initialValues: UserFormValues = {
  username: "",
  email: "",
  password: "",
  roleIds: "",
};

export default function UsersAdminPage() {
  const session = useRequireSession();
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [values, setValues] = useState(initialValues);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadData = useCallback(async () => {
    if (!session.token) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const [nextUsers, nextRoles] = await Promise.all([
        listUsers(session.token),
        listRoles(session.token),
      ]);
      setUsers(nextUsers);
      setRoles(nextRoles);
    } catch (fetchError) {
      setError(fetchError instanceof Error ? fetchError.message : "Failed to load security data");
    } finally {
      setLoading(false);
    }
  }, [session.token]);

  useEffect(() => {
    if (session.isReady && session.user) {
      void loadData();
    }
  }, [loadData, session.isReady, session.user]);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!session.token) {
      return;
    }

    setSaving(true);
    setError(null);
    try {
      await createUser(
        {
          username: values.username.trim(),
          email: values.email.trim(),
          password: values.password,
          roleIds: values.roleIds
            .split(",")
            .map((value) => value.trim())
            .filter(Boolean),
        },
        session.token,
      );
      setValues(initialValues);
      await loadData();
    } catch (saveError) {
      setError(saveError instanceof Error ? saveError.message : "Failed to create user");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (username?: string) => {
    if (!session.token || !username) {
      return;
    }

    try {
      await deleteUser(username, session.token);
      await loadData();
    } catch (deleteError) {
      setError(deleteError instanceof Error ? deleteError.message : "Failed to delete user");
    }
  };

  return (
    <AppShell>
      <div>
        <p className="subtle-label">Administration</p>
        <h1 className="mt-2 text-3xl font-semibold text-white">Users</h1>
      </div>

      <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_360px]">
        <Card>
          {loading ? (
            <p className="text-sm text-slate-400">Loading users...</p>
          ) : users.length === 0 ? (
            <EmptyState title="No users" description="Create the first account to populate RBAC." />
          ) : (
            <div className="space-y-3">
              {users.map((entry) => (
                <div
                  key={entry.id || entry.username}
                  className="rounded-2xl border border-white/10 bg-white/[0.04] p-4"
                >
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h2 className="text-lg font-semibold text-white">{entry.username}</h2>
                      <p className="mt-1 text-sm text-slate-400">{entry.email || "No email"}</p>
                      <div className="mt-3 flex flex-wrap gap-2">
                        {(entry.roles || []).map((role) => (
                          <Badge key={role.id || role.name}>{role.name || "role"}</Badge>
                        ))}
                      </div>
                    </div>
                    <Button variant="danger" onClick={() => handleDelete(entry.username)}>
                      Delete
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>

        <Card>
          <h2 className="text-xl font-semibold text-white">Create user</h2>
          <p className="mt-2 text-sm text-slate-400">
            Role IDs can be copied from the roles section and entered as a comma-separated list.
          </p>

          <form className="mt-6 space-y-4" onSubmit={onSubmit}>
            <Input
              placeholder="username"
              value={values.username}
              onChange={(event) =>
                setValues((current) => ({ ...current, username: event.target.value }))
              }
            />
            <Input
              placeholder="email"
              value={values.email}
              onChange={(event) =>
                setValues((current) => ({ ...current, email: event.target.value }))
              }
            />
            <Input
              placeholder="password"
              type="password"
              value={values.password}
              onChange={(event) =>
                setValues((current) => ({ ...current, password: event.target.value }))
              }
            />
            <Input
              placeholder="role id, role id"
              value={values.roleIds}
              onChange={(event) =>
                setValues((current) => ({ ...current, roleIds: event.target.value }))
              }
            />

            {roles.length > 0 ? (
              <div className="rounded-2xl bg-black/20 p-4 text-xs text-slate-300">
                <p className="mb-2 font-semibold text-white">Available roles</p>
                {roles.map((role) => (
                  <p key={role.id || role.name}>
                    {role.name}: {role.id}
                  </p>
                ))}
              </div>
            ) : null}

            {error ? (
              <div className="rounded-2xl border border-rose-400/20 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
                {error}
              </div>
            ) : null}

            <Button type="submit" className="w-full" disabled={saving}>
              {saving ? "Creating..." : "Create user"}
            </Button>
          </form>
        </Card>
      </div>
    </AppShell>
  );
}
