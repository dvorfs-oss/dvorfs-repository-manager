import { request } from "@/lib/api/client";
import type { Role, User } from "@/types/api";

export async function listUsers(token?: string | null) {
  return request<User[]>("/security/users", {
    method: "GET",
    authToken: token,
  });
}

export async function createUser(
  payload: { username: string; email: string; password: string; roleIds: string[] },
  token?: string | null,
) {
  return request<User>("/security/users", {
    method: "POST",
    body: payload,
    authToken: token,
  });
}

export async function deleteUser(username: string, token?: string | null) {
  return request<void>(`/security/users/${encodeURIComponent(username)}`, {
    method: "DELETE",
    authToken: token,
  });
}

export async function changeUserPassword(
  username: string,
  password: string,
  token?: string | null,
) {
  return request<void>(`/security/users/${encodeURIComponent(username)}/password`, {
    method: "PUT",
    body: { password },
    authToken: token,
  });
}

export async function listRoles(token?: string | null) {
  return request<Role[]>("/security/roles", {
    method: "GET",
    authToken: token,
  });
}

export async function createRole(
  payload: { name: string; privileges: unknown },
  token?: string | null,
) {
  return request<Role>("/security/roles", {
    method: "POST",
    body: payload,
    authToken: token,
  });
}

export async function deleteRole(roleId: string, token?: string | null) {
  return request<void>(`/security/roles/${encodeURIComponent(roleId)}`, {
    method: "DELETE",
    authToken: token,
  });
}
