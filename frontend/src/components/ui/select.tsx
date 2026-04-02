import type { SelectHTMLAttributes } from "react";

import { cn } from "@/lib/utils";

export function Select({ className, ...props }: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        "w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-100",
        "focus:border-accent-400/60 focus:outline-none focus:ring-2 focus:ring-accent-400/20",
        className,
      )}
      {...props}
    />
  );
}
