export type JsonPrimitive = string | number | boolean | null;
export type JsonValue = JsonPrimitive | JsonValue[] | { [key: string]: JsonValue };

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token?: string;
  user?: User;
}

export interface Role {
  id?: string;
  name?: string;
  privileges?: JsonValue;
  createdAt?: string;
}

export interface User {
  id?: string;
  username?: string;
  email?: string;
  passwordHash?: string;
  roles?: Role[];
  createdAt?: string;
  updatedAt?: string;
}

export interface Artifact {
  id?: string;
  repositoryID?: string;
  path?: string;
  size?: number;
  contentType?: string;
  checksums?: JsonValue;
  createdAt?: string;
  lastDownloadedAt?: string | null;
}

export interface CleanupPolicy {
  id?: string;
  name?: string;
  criteria?: JsonValue;
  createdAt?: string;
  updatedAt?: string;
}

export interface BlobStore {
  id?: string;
  name?: string;
  type?: string;
  attributes?: JsonValue;
  createdAt?: string;
  updatedAt?: string;
  repositories?: Repository[];
}

export interface Repository {
  id?: string;
  name?: string;
  format?: string;
  type?: string;
  attributes?: JsonValue;
  cleanupPolicyID?: string | null;
  blobStoreID?: string | null;
  createdAt?: string;
  updatedAt?: string;
  artifacts?: Artifact[];
  cleanupPolicy?: CleanupPolicy | null;
  blobStore?: BlobStore | null;
}

export interface RepositoryFormValues {
  name: string;
  format: "raw" | "maven";
  type: "hosted";
  blobStoreID: string;
  cleanupPolicyID: string;
  attributes: string;
}

export interface UserFormValues {
  username: string;
  email: string;
  password: string;
  roleIds: string;
}

export interface RoleFormValues {
  name: string;
  privileges: string;
}

export interface BlobStoreFormValues {
  name: string;
  type: string;
  attributes: string;
}

export interface CleanupPolicyFormValues {
  name: string;
  criteria: string;
}
