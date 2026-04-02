"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { useSession } from "@/hooks/use-session";

export function useRequireSession() {
  const router = useRouter();
  const session = useSession();

  useEffect(() => {
    if (session.isReady && !session.isLoading && !session.user) {
      router.replace("/login");
    }
  }, [router, session.isLoading, session.isReady, session.user]);

  return session;
}
