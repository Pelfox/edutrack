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
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(() => getAuthToken() !== null);

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

    setAuthToken(data.token);
    setUser(data.user ?? null);
    setIsAuthenticated(true);

    return data.user ?? null;
  }, []);

  const logout = useCallback(() => {
    clearAuthToken();
    setUser(null);
    setIsAuthenticated(false);
  }, []);

  const value = useMemo(
    () => ({
      isAuthenticated,
      login,
      logout,
      user,
    }),
    [isAuthenticated, login, logout, user],
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
