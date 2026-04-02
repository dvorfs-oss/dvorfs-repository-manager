"use client";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useSession } from "@/hooks/use-session";
import { formatDate } from "@/lib/utils";

export function Header() {
  const session = useSession();

  return (
    <header className="flex flex-col gap-4 border-b border-white/10 bg-ink-950/55 px-6 py-5 backdrop-blur md:flex-row md:items-center md:justify-between">
      <div>
        <p className="text-xs uppercase tracking-[0.3em] text-slate-500">Connected session</p>
        <div className="mt-2 flex flex-wrap items-center gap-3">
          <h2 className="text-lg font-semibold text-white">
            {session.user?.username || "Guest"}
          </h2>
          <Badge>{session.token ? "Authenticated" : "Anonymous"}</Badge>
        </div>
        {session.user?.updatedAt ? (
          <p className="mt-1 text-sm text-slate-400">
            Last update {formatDate(session.user.updatedAt)}
          </p>
        ) : null}
      </div>

      <div className="flex items-center gap-3">
        {session.error ? (
          <span className="rounded-full bg-rose-500/10 px-3 py-2 text-xs text-rose-100 ring-1 ring-rose-400/20">
            {session.error}
          </span>
        ) : null}
        <Button variant="secondary" onClick={() => session.logout()}>
          Logout
        </Button>
      </div>
    </header>
  );
}
