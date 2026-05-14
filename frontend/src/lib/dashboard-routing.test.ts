import { describe, expect, it } from "vitest";

import {
  getAdministratorPage,
  getDashboardRole,
  getStudentPage,
  getTeacherPage,
  isDashboardPageAvailable,
} from "./dashboard-routing";

describe("dashboard routing", () => {
  it("normalizes unsupported roles to student", () => {
    expect(getDashboardRole("administrator")).toBe("administrator");
    expect(getDashboardRole("teacher")).toBe("teacher");
    expect(getDashboardRole("student")).toBe("student");
    expect(getDashboardRole(undefined)).toBe("student");
  });

  it("allows administrator pages only for administrators", () => {
    expect(isDashboardPageAvailable("administrator", "analytics")).toBe(true);
    expect(isDashboardPageAvailable("administrator", "students")).toBe(true);
    expect(isDashboardPageAvailable("teacher", "analytics")).toBe(false);
    expect(isDashboardPageAvailable("student", "students")).toBe(false);
  });

  it("allows teacher pages for teachers", () => {
    expect(isDashboardPageAvailable("teacher", "disciplines")).toBe(true);
    expect(isDashboardPageAvailable("teacher", "grades")).toBe(true);
    expect(isDashboardPageAvailable("teacher", "schedule")).toBe(true);
    expect(isDashboardPageAvailable("teacher", "students")).toBe(false);
  });

  it("allows student pages for students", () => {
    expect(isDashboardPageAvailable("student", "disciplines")).toBe(true);
    expect(isDashboardPageAvailable("student", "curriculums")).toBe(true);
    expect(isDashboardPageAvailable("student", "grades")).toBe(true);
    expect(isDashboardPageAvailable("student", "analytics")).toBe(false);
  });

  it("falls back to home when a role-specific page is unavailable", () => {
    expect(getAdministratorPage("analytics")).toBe("analytics");
    expect(getAdministratorPage("grades")).toBe("home");
    expect(getTeacherPage("grades")).toBe("grades");
    expect(getTeacherPage("analytics")).toBe("home");
    expect(getStudentPage("curriculums")).toBe("curriculums");
    expect(getStudentPage("students")).toBe("home");
  });
});
