import type { ReactNode } from "react";

import { Header } from "@/components/header";
import { Sidebar } from "@/components/sidebar";

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-mesh-radial text-slate-100">
      <div className="grid min-h-screen lg:grid-cols-[300px_1fr]">
        <Sidebar />
        <div className="flex min-h-screen flex-col">
          <Header />
          <main className="flex-1 px-4 py-6 md:px-6 lg:px-8">{children}</main>
        </div>
      </div>
    </div>
  );
}
