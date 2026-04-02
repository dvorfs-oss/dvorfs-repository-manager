import { cloneElement, isValidElement, type ButtonHTMLAttributes, type ReactElement } from "react";

import { cn } from "@/lib/utils";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  asChild?: boolean;
};

const variants = {
  primary:
    "bg-accent-400 text-ink-950 shadow-lg shadow-accent-500/20 hover:bg-accent-300",
  secondary:
    "bg-white/5 text-slate-100 ring-1 ring-white/10 hover:bg-white/10 hover:ring-white/20",
  ghost: "text-slate-200 hover:bg-white/5",
  danger: "bg-rose-500/15 text-rose-100 ring-1 ring-rose-400/30 hover:bg-rose-500/25",
};

export function Button({
  className,
  variant = "primary",
  asChild,
  children,
  ...props
}: ButtonProps) {
  const shared = cn(
    "inline-flex items-center justify-center rounded-xl px-4 py-2 text-sm font-semibold transition",
    "disabled:cursor-not-allowed disabled:opacity-50",
    variants[variant],
    className,
  );

  if (asChild && isValidElement(children)) {
    return cloneElement(children as ReactElement<{ className?: string }>, {
      className: cn(shared, (children as ReactElement<{ className?: string }>).props.className),
    });
  }

  return (
    <button className={shared} {...props}>
      {children}
    </button>
  );
}
