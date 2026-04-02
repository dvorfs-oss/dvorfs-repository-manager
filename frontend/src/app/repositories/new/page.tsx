"use client";

import { useRouter } from "next/navigation";
import type { FormEvent } from "react";
import { useState } from "react";

import { AppShell } from "@/components/app-shell";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { useRequireSession } from "@/hooks/use-require-session";
import { createRepository, parseRepositoryAttributes } from "@/lib/api/repositories";
import type { RepositoryFormValues } from "@/types/api";

const initialValues: RepositoryFormValues = {
  name: "",
  format: "raw",
  type: "hosted",
  blobStoreID: "",
  cleanupPolicyID: "",
  attributes: "{}",
};

export default function NewRepositoryPage() {
  const router = useRouter();
  const session = useRequireSession();
  const [values, setValues] = useState(initialValues);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const setField = <K extends keyof RepositoryFormValues>(key: K, value: RepositoryFormValues[K]) =>
    setValues((current) => ({ ...current, [key]: value }));

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setSaving(true);

    try {
      const payload = parseRepositoryAttributes(values);
      const created = await createRepository(payload, session.token);
      const createdName = (created as { name?: string } | undefined)?.name;
      const repositoryName =
        typeof createdName === "string" && createdName.trim() ? createdName : values.name;
      router.push(`/repositories/${encodeURIComponent(repositoryName)}`);
    } catch (createError) {
      setError(createError instanceof Error ? createError.message : "Failed to create repository");
    } finally {
      setSaving(false);
    }
  };

  return (
    <AppShell>
      <div className="max-w-3xl">
        <p className="subtle-label">Create repository</p>
        <h1 className="mt-2 text-3xl font-semibold text-white">Hosted repository setup</h1>
        <p className="mt-2 text-sm text-slate-400">
          The MVP only enables hosted RAW and Maven repositories. Proxy and group support can be
          layered in later without changing this flow.
        </p>
      </div>

      <Card className="mt-6 max-w-3xl">
        <form className="grid gap-5" onSubmit={onSubmit}>
          <div className="grid gap-5 md:grid-cols-2">
            <label className="block space-y-2">
              <span className="text-sm font-medium text-slate-200">Name</span>
              <Input
                placeholder="maven-hosted"
                value={values.name}
                onChange={(event) => setField("name", event.target.value)}
              />
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-slate-200">Format</span>
              <Select
                value={values.format}
                onChange={(event) =>
                  setField("format", event.target.value as RepositoryFormValues["format"])
                }
              >
                <option value="raw">raw</option>
                <option value="maven">maven</option>
              </Select>
            </label>
          </div>

          <label className="block space-y-2">
            <span className="text-sm font-medium text-slate-200">Type</span>
            <Input value="hosted" disabled />
          </label>

          <div className="grid gap-5 md:grid-cols-2">
            <label className="block space-y-2">
              <span className="text-sm font-medium text-slate-200">Blob store ID</span>
              <Input
                placeholder="optional"
                value={values.blobStoreID}
                onChange={(event) => setField("blobStoreID", event.target.value)}
              />
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-slate-200">Cleanup policy ID</span>
              <Input
                placeholder="optional"
                value={values.cleanupPolicyID}
                onChange={(event) => setField("cleanupPolicyID", event.target.value)}
              />
            </label>
          </div>

          <label className="block space-y-2">
            <span className="text-sm font-medium text-slate-200">Attributes JSON</span>
            <Textarea
              rows={8}
              value={values.attributes}
              onChange={(event) => setField("attributes", event.target.value)}
            />
          </label>

          {error ? (
            <div className="rounded-2xl border border-rose-400/20 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
              {error}
            </div>
          ) : null}

          <div className="flex flex-wrap gap-3">
            <Button type="submit" disabled={saving}>
              {saving ? "Creating..." : "Create repository"}
            </Button>
            <Button type="button" variant="secondary" onClick={() => router.push("/repositories")}>
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </AppShell>
  );
}
