import type { DashboardRole } from "@/components/dashboard/dashboard-layout";
import type { AdministratorPage } from "@/components/dashboard/roles/administrator";
import type { StudentPage } from "@/components/dashboard/roles/student";
import type { TeacherPage } from "@/components/dashboard/roles/teacher";
import type { AuthUser } from "@/lib/context/auth";

export type DashboardPage =
  | "home"
  | "students"
  | "teachers"
  | "groups"
  | "specialties"
  | "disciplines"
  | "curriculums"
  | "analytics"
  | "grades"
  | "schedule";

export function getDashboardRole(role: AuthUser["role"] | undefined): DashboardRole {
  if (role === "administrator" || role === "teacher" || role === "student") {
    return role;
  }

  return "student";
}

export function isDashboardPageAvailable(role: DashboardRole | undefined, page: DashboardPage) {
  if (page === "home") {
    return true;
  }

  if (role === "administrator") {
    return (
      page === "students" ||
      page === "teachers" ||
      page === "groups" ||
      page === "specialties" ||
      page === "disciplines" ||
      page === "curriculums" ||
      page === "analytics"
    );
  }

  if (role === "teacher") {
    return page === "disciplines" || page === "grades" || page === "schedule";
  }

  return (
    page === "disciplines" || page === "curriculums" || page === "grades" || page === "schedule"
  );
}

export function getAdministratorPage(page: DashboardPage): AdministratorPage {
  if (
    page === "students" ||
    page === "teachers" ||
    page === "groups" ||
    page === "specialties" ||
    page === "disciplines" ||
    page === "curriculums" ||
    page === "analytics"
  ) {
    return page;
  }

  return "home";
}

export function getTeacherPage(page: DashboardPage): TeacherPage {
  if (page === "disciplines" || page === "grades" || page === "schedule") {
    return page;
  }

  return "home";
}

export function getStudentPage(page: DashboardPage): StudentPage {
  if (
    page === "disciplines" ||
    page === "curriculums" ||
    page === "grades" ||
    page === "schedule"
  ) {
    return page;
  }

  return "home";
}
