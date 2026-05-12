import createClient, { type Middleware } from "openapi-fetch";

import type { paths } from "@/api/schema";

const defaultApiBaseUrl = "http://localhost:8000";
const authTokenCookieMaxAge = 60 * 60 * 24 * 7;

export const authTokenCookieName = "edutrack_auth_token";

function getApiBaseUrl() {
  return import.meta.env.VITE_API_BASE_URL?.trim() || defaultApiBaseUrl;
}

export function getAuthToken() {
  if (typeof document === "undefined" || document.cookie.length === 0) {
    return null;
  }

  const cookieName = `${encodeURIComponent(authTokenCookieName)}=`;
  const authTokenCookie = document.cookie
    .split("; ")
    .find((cookie) => cookie.startsWith(cookieName));

  if (!authTokenCookie) {
    return null;
  }

  return decodeURIComponent(authTokenCookie.slice(cookieName.length));
}

export function setAuthToken(token: string) {
  const secureCookieAttribute = window.location.protocol === "https:" ? "; Secure" : "";

  // biome-ignore lint/suspicious/noDocumentCookie: frontend sets the token until backend issues HttpOnly cookies.
  document.cookie = `${encodeURIComponent(authTokenCookieName)}=${encodeURIComponent(
    token,
  )}; Path=/; SameSite=Lax; Max-Age=${authTokenCookieMaxAge}${secureCookieAttribute}`;
}

export function clearAuthToken() {
  // biome-ignore lint/suspicious/noDocumentCookie: frontend clears the token until backend issues HttpOnly cookies.
  document.cookie = `${encodeURIComponent(authTokenCookieName)}=; Path=/; SameSite=Lax; Max-Age=0`;
}

const authMiddleware: Middleware = {
  onRequest({ request }) {
    const authToken = getAuthToken();

    if (authToken) {
      request.headers.set("Authorization", `Bearer ${authToken}`);
    }

    return request;
  },
};

export const apiClient = createClient<paths>({
  baseUrl: getApiBaseUrl(),
  credentials: "include",
});

apiClient.use(authMiddleware);
