"use client";

import { useState } from "react";

import { AppShell } from "@/components/app-shell";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { EmptyState } from "@/components/ui/empty-state";
import { Input } from "@/components/ui/input";
import { useRequireSession } from "@/hooks/use-require-session";
import { searchArtifacts } from "@/lib/api/repositories";
import { formatDate } from "@/lib/utils";
import type { Artifact } from "@/types/api";

export default function SearchPage() {
  const session = useRequireSession();
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Artifact[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async () => {
    if (!session.token) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      setResults(await searchArtifacts(query, session.token));
    } catch (searchError) {
      setError(searchError instanceof Error ? searchError.message : "Search failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AppShell>
      <div>
        <p className="subtle-label">Search</p>
        <h1 className="mt-2 text-3xl font-semibold text-white">Artifact search</h1>
      </div>

      <Card className="mt-6">
        <div className="flex flex-col gap-3 md:flex-row">
          <Input
            placeholder="package, file name, content type"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
          <Button onClick={handleSearch} disabled={loading}>
            {loading ? "Searching..." : "Search"}
          </Button>
        </div>
      </Card>

      {error ? (
        <Card className="mt-6 border-rose-400/20 bg-rose-500/10">
          <p className="text-sm text-rose-100">{error}</p>
        </Card>
      ) : null}

      <div className="mt-6">
        {results.length === 0 ? (
          <EmptyState
            title="No search results yet"
            description="Run a query after uploading a few artifacts to explore what is indexed."
          />
        ) : (
          <div className="grid gap-4">
            {results.map((artifact) => (
              <Card key={artifact.id || artifact.path}>
                <h2 className="text-lg font-semibold text-white">{artifact.path || "Artifact"}</h2>
                <p className="mt-2 text-sm text-slate-400">
                  {artifact.contentType || "application/octet-stream"} · {artifact.size ?? 0} bytes
                </p>
                <p className="mt-1 text-xs text-slate-500">
                  Created {formatDate(artifact.createdAt)}
                </p>
              </Card>
            ))}
          </div>
        )}
      </div>
    </AppShell>
  );
}
