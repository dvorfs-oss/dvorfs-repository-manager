"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { cn } from "@/lib/utils";

const docsHref = `${(process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080").replace(/\/$/, "")}/swagger/index.html`;

const links = [
  { href: "/repositories", label: "Repositories" },
  { href: "/repositories/new", label: "Create repository" },
  { href: "/search", label: "Artifact search" },
  { href: "/admin/users", label: "Users" },
  { href: "/admin/roles", label: "Roles" },
  { href: "/admin/blob-stores", label: "Blob stores" },
  { href: "/admin/cleanup-policies", label: "Cleanup policies" },
  { href: docsHref, label: "API docs" },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="border-b border-white/10 bg-ink-950/80 px-5 py-6 lg:min-h-screen lg:border-b-0 lg:border-r">
      <div className="mb-8">
        <p className="text-xs uppercase tracking-[0.35em] text-accent-300/80">Dvorfs</p>
        <h1 className="mt-2 text-2xl font-semibold text-white">Repository Manager</h1>
      </div>

      <nav className="space-y-2">
        {links.map((link) => {
          const active = pathname === link.href || pathname.startsWith(`${link.href}/`);
          return (
            <Link
              key={link.href}
              href={link.href}
              className={cn(
                "block rounded-2xl px-4 py-3 text-sm transition",
                active
                  ? "bg-accent-400/15 text-accent-200 ring-1 ring-accent-400/20"
                  : "text-slate-300 hover:bg-white/5 hover:text-white",
              )}
            >
              {link.label}
            </Link>
          );
        })}
      </nav>

      <div className="mt-8 rounded-3xl border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
        Surface area now includes security, storage settings, cleanup policies, search, and repository operations.
      </div>
    </aside>
  );
}
