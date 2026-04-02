import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

export function Card({
  className,
  children,
}: {
  className?: string;
  children: ReactNode;
}) {
  return (
    <div
      className={cn(
        "rounded-3xl border border-white/10 bg-white/5 p-6 shadow-glow backdrop-blur",
        className,
      )}
    >
      {children}
    </div>
  );
}
