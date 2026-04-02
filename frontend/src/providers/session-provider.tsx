"use client";

import type { ReactNode } from "react";
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

import { getMe, login as loginRequest, logout as logoutRequest } from "@/lib/api/auth";
import type { LoginRequest, User } from "@/types/api";

type SessionState = {
  token: string | null;
  user: User | null;
  isReady: boolean;
  isLoading: boolean;
  error: string | null;
};

type SessionContextValue = SessionState & {
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
  setToken: (token: string | null) => void;
};

const storageKey = "dvorfs.session.token";
const SessionContext = createContext<SessionContextValue | null>(null);

export function SessionProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(null);
  const [user, setUser] = useState<User | null>(null);
  const [isReady, setIsReady] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const persistToken = useCallback((value: string | null) => {
    setTokenState(value);
    if (typeof window === "undefined") {
      return;
    }

    if (value) {
      window.localStorage.setItem(storageKey, value);
    } else {
      window.localStorage.removeItem(storageKey);
    }
  }, []);

  const loadUser = useCallback(async (nextToken: string) => {
    const loadedUser = await getMe(nextToken);
    setUser(loadedUser);
  }, []);

  const refreshUser = useCallback(async () => {
    if (!token) {
      setUser(null);
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      await loadUser(token);
    } catch (fetchError) {
      setUser(null);
      persistToken(null);
      setError(fetchError instanceof Error ? fetchError.message : "Failed to load session");
    } finally {
      setIsLoading(false);
      setIsReady(true);
    }
  }, [persistToken, token]);

  useEffect(() => {
    const storedToken =
      typeof window === "undefined" ? null : window.localStorage.getItem(storageKey);

    if (!storedToken) {
      setIsLoading(false);
      setIsReady(true);
      return;
    }

    persistToken(storedToken);
    loadUser(storedToken)
      .catch(() => {
        persistToken(null);
        setUser(null);
      })
      .finally(() => {
        setIsLoading(false);
        setIsReady(true);
      });
  }, []);

  const login = useCallback(async (credentials: LoginRequest) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await loginRequest(credentials);
      const nextToken = response.token ?? null;
      if (!nextToken) {
        throw new Error("Login response did not include a token");
      }

      persistToken(nextToken);
      await loadUser(nextToken);
    } catch (loginError) {
      persistToken(null);
      setUser(null);
      throw loginError;
    } finally {
      setIsLoading(false);
      setIsReady(true);
    }
  }, [loadUser, persistToken]);

  const logout = useCallback(async () => {
    setIsLoading(true);
    try {
      await logoutRequest(token);
    } catch {
      // Local logout should still succeed even if the API is unavailable.
    } finally {
      persistToken(null);
      setUser(null);
      setIsLoading(false);
      setIsReady(true);
    }
  }, [persistToken, token]);

  const value = useMemo<SessionContextValue>(
    () => ({
      token,
      user,
      isReady,
      isLoading,
      error,
      login,
      logout,
      refreshUser,
      setToken: persistToken,
    }),
    [error, isLoading, isReady, login, logout, persistToken, refreshUser, token, user],
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession() {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error("useSession must be used inside SessionProvider");
  }

  return context;
}
