import { request } from "@/lib/api/client";
import type { LoginRequest, LoginResponse, User } from "@/types/api";

export async function login(credentials: LoginRequest) {
  return request<LoginResponse>("/auth/login", {
    method: "POST",
    body: credentials,
  });
}

export async function logout(token?: string | null) {
  return request<void>("/auth/logout", {
    method: "POST",
    authToken: token,
  });
}

export async function getMe(token?: string | null) {
  return request<User>("/auth/me", {
    method: "GET",
    authToken: token,
  });
}
