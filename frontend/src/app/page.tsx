"use client";

import Link from "next/link";
import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useSession } from "@/hooks/use-session";

export default function HomePage() {
  const router = useRouter();
  const session = useSession();

  useEffect(() => {
    if (session.isReady && !session.isLoading) {
      router.replace(session.user ? "/repositories" : "/login");
    }
  }, [router, session.isLoading, session.isReady, session.user]);

  return (
    <main className="flex min-h-screen items-center justify-center px-4 py-12">
      <div className="relative w-full max-w-4xl overflow-hidden rounded-[2rem] border border-white/10 bg-white/[0.04] p-8 shadow-glow">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_right,rgba(63,191,147,0.18),transparent_30%),radial-gradient(circle_at_bottom_left,rgba(82,108,163,0.14),transparent_28%)]" />
        <div className="relative grid gap-8 md:grid-cols-[1.2fr_0.8fr] md:items-center">
          <div className="space-y-5">
            <p className="subtle-label">Dvorfs Repository Manager</p>
            <h1 className="max-w-xl text-4xl font-semibold tracking-tight text-white md:text-6xl">
              A focused control plane for repositories and artifacts.
            </h1>
            <p className="max-w-2xl text-base leading-7 text-slate-300">
              Operational console for auth, repository CRUD, artifact movement, search,
              users, roles, blob stores, and cleanup policy management.
            </p>
            <div className="flex flex-wrap gap-3">
              <Button asChild>
                <Link href="/login">Open login</Link>
              </Button>
              <Button variant="secondary" asChild>
                <Link href="/repositories">View repositories</Link>
              </Button>
            </div>
          </div>

          <Card className="relative border-white/15 bg-ink-950/60">
            <div className="space-y-4">
              <div>
                <p className="subtle-label">MVP coverage</p>
                <h2 className="mt-2 text-xl font-semibold text-white">Ready for the first release</h2>
              </div>
              <ul className="space-y-3 text-sm text-slate-300">
                <li>Login and session persistence</li>
                <li>Repository list and creation form</li>
                <li>Repository details with artifact upload/download</li>
                <li>Search, users, roles, blob stores, cleanup policies</li>
              </ul>
            </div>
          </Card>
        </div>
      </div>
    </main>
  );
}
