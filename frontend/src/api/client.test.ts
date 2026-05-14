import { beforeEach, describe, expect, it } from "vitest";

import { authTokenCookieName, clearAuthToken, getAuthToken, setAuthToken } from "./client";

describe("auth token cookie helpers", () => {
  beforeEach(() => {
    clearAuthToken();
  });

  it("returns null when auth token cookie is missing", () => {
    expect(getAuthToken()).toBeNull();
  });

  it("sets and reads auth token from cookies", () => {
    setAuthToken("token.with.payload");

    expect(getAuthToken()).toBe("token.with.payload");
    expect(document.cookie).toContain(`${authTokenCookieName}=token.with.payload`);
  });

  it("encodes token values before writing cookies", () => {
    setAuthToken("token with spaces");

    expect(getAuthToken()).toBe("token with spaces");
    expect(document.cookie).toContain(`${authTokenCookieName}=token%20with%20spaces`);
  });

  it("clears auth token cookie", () => {
    setAuthToken("token.with.payload");

    clearAuthToken();

    expect(getAuthToken()).toBeNull();
  });
});
