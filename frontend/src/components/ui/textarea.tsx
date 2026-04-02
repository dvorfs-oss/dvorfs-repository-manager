import type { TextareaHTMLAttributes } from "react";

import { cn } from "@/lib/utils";

export function Textarea({ className, ...props }: TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      className={cn(
        "w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-100",
        "placeholder:text-slate-500 focus:border-accent-400/60 focus:outline-none focus:ring-2 focus:ring-accent-400/20",
        className,
      )}
      {...props}
    />
  );
}
