import type { ReactNode } from "react";
import { createContext, useCallback, useContext, useMemo, useState } from "react";

import type { components } from "@/api";
import { apiClient, clearAuthToken, getAuthToken, setAuthToken } from "@/api";

export type AuthUser = components["schemas"]["dto.User"];
export type LoginCredentials = components["schemas"]["dto.Login"];

export type AuthContextValue = {
  isAuthenticated: boolean;
  login: (credentials: LoginCredentials) => Promise<AuthUser | null>;
  logout: () => void;
  user: AuthUser | null;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(() => getInitialAuthUser());

  const login = useCallback(async (credentials: LoginCredentials) => {
    const { data, error } = await apiClient.POST("/auth/login", {
      body: credentials,
    });

    if (error) {
      throw new Error(getAuthErrorMessage(error));
    }

    if (!data?.token) {
      throw new Error("Сервер не вернул токен авторизации.");
    }

    const authenticatedUser = data.user ?? getAuthUserFromToken(data.token);
    if (!authenticatedUser?.role) {
      throw new Error("Сервер не вернул роль пользователя.");
    }

    setAuthToken(data.token);
    setUser(authenticatedUser);

    return authenticatedUser;
  }, []);

  const logout = useCallback(() => {
    clearAuthToken();
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({
      isAuthenticated: user !== null,
      login,
      logout,
      user,
    }),
    [login, logout, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// biome-ignore lint/style/useComponentExportOnlyModules: Auth module intentionally keeps provider and hook together.
export function useAuth() {
  const auth = useContext(AuthContext);

  if (!auth) {
    throw new Error("useAuth must be used within AuthProvider.");
  }

  return auth;
}

function getAuthErrorMessage(error: components["schemas"]["dto.Error"]) {
  return error.message ?? error.error ?? "Не удалось войти в аккаунт.";
}

function getInitialAuthUser() {
  const token = getAuthToken();

  if (!token) {
    return null;
  }

  const tokenUser = getAuthUserFromToken(token);

  if (!tokenUser) {
    clearAuthToken();
    return null;
  }

  return tokenUser;
}

function getAuthUserFromToken(token: string): AuthUser | null {
  const tokenPayload = decodeTokenPayload(token);

  if (!tokenPayload || isExpired(tokenPayload.exp) || !isUserRole(tokenPayload.role)) {
    return null;
  }

  if (typeof tokenPayload.user_id === "string") {
    return {
      id: tokenPayload.user_id,
      role: tokenPayload.role,
    };
  }

  return {
    role: tokenPayload.role,
  };
}

function decodeTokenPayload(token: string): Record<string, unknown> | null {
  const [, payload] = token.split(".");

  if (!payload) {
    return null;
  }

  try {
    const normalizedPayload = payload.replaceAll("-", "+").replaceAll("_", "/");
    const padding = "=".repeat((4 - (normalizedPayload.length % 4)) % 4);

    return JSON.parse(atob(`${normalizedPayload}${padding}`));
  } catch {
    return null;
  }
}

function isExpired(expiresAt: unknown) {
  return typeof expiresAt === "number" && expiresAt * 1000 <= Date.now();
}

function isUserRole(role: unknown): role is NonNullable<AuthUser["role"]> {
  return role === "administrator" || role === "teacher" || role === "student";
}
