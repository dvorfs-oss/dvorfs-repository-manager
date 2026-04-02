"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";

import { AppShell } from "@/components/app-shell";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { useRequireSession } from "@/hooks/use-require-session";
import { listRepositories } from "@/lib/api/repositories";
import { formatDate } from "@/lib/utils";
import type { Repository } from "@/types/api";

export default function RepositoriesPage() {
  const session = useRequireSession();
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const loadRepositories = useCallback(async () => {
    if (!session.token) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      setRepositories(await listRepositories(session.token));
    } catch (fetchError) {
      setError(fetchError instanceof Error ? fetchError.message : "Failed to load repositories");
    } finally {
      setLoading(false);
    }
  }, [session.token]);

  useEffect(() => {
    if (session.isReady && session.user) {
      void loadRepositories();
    }
  }, [loadRepositories, session.isReady, session.user]);

  return (
    <AppShell>
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="subtle-label">Repositories</p>
          <h1 className="mt-2 text-3xl font-semibold text-white">Storage inventory</h1>
          <p className="mt-2 max-w-2xl text-sm text-slate-400">
            Hosted repositories for RAW and Maven are the primary MVP surface. This view keeps
            the operational picture in one place.
          </p>
        </div>
        <Button asChild>
          <Link href="/repositories/new">Create repository</Link>
        </Button>
      </div>

      <div className="mt-6 grid gap-4 md:grid-cols-3">
        <Card>
          <p className="text-sm text-slate-400">Total repositories</p>
          <p className="mt-3 text-3xl font-semibold text-white">{repositories.length}</p>
        </Card>
        <Card>
          <p className="text-sm text-slate-400">Formats in scope</p>
          <p className="mt-3 text-3xl font-semibold text-white">RAW, Maven</p>
        </Card>
        <Card>
          <p className="text-sm text-slate-400">Session</p>
          <p className="mt-3 text-3xl font-semibold text-white">
            {session.user?.username || "guest"}
          </p>
        </Card>
      </div>

      <div className="mt-6">
        {loading ? (
          <Card>
            <p className="text-sm text-slate-400">Loading repositories...</p>
          </Card>
        ) : error ? (
          <Card className="border-rose-400/20 bg-rose-500/10">
            <p className="text-sm text-rose-100">{error}</p>
            <div className="mt-4">
              <Button variant="secondary" onClick={loadRepositories}>
                Retry
              </Button>
            </div>
          </Card>
        ) : repositories.length === 0 ? (
          <EmptyState
            title="No repositories yet"
            description="Create the first hosted repository to start storing artifacts."
          />
        ) : (
          <div className="grid gap-4">
            {repositories.map((repository) => (
              <Card key={repository.id || repository.name}>
                <div className="flex flex-wrap items-start justify-between gap-4">
                  <div>
                    <div className="flex flex-wrap items-center gap-2">
                      <h2 className="text-xl font-semibold text-white">{repository.name}</h2>
                      <Badge>{repository.format || "unknown"}</Badge>
                      <Badge>{repository.type || "unknown"}</Badge>
                    </div>
                    <p className="mt-2 text-sm text-slate-400">
                      Updated {formatDate(repository.updatedAt)}.
                    </p>
                  </div>

                  <Button variant="secondary" asChild>
                    <Link href={`/repositories/${encodeURIComponent(repository.name || "")}`}>
                      Open
                    </Link>
                  </Button>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </AppShell>
  );
}
