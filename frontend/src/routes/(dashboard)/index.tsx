import { createFileRoute } from "@tanstack/react-router";

import type { DashboardRole } from "@/components/dashboard/dashboard-layout";
import { DashboardView } from "@/components/dashboard/dashboard-views";
import type { AdministratorPage } from "@/components/dashboard/roles/administrator";
import type { StudentPage } from "@/components/dashboard/roles/student";
import type { TeacherPage } from "@/components/dashboard/roles/teacher";

export const Route = createFileRoute("/(dashboard)/")({
  validateSearch: (
    search,
  ): {
    administratorPage: AdministratorPage;
    studentPage: StudentPage;
    teacherPage: TeacherPage;
    view: DashboardRole;
  } => {
    const view = parseDashboardRole(search.view);

    return {
      administratorPage: parseAdministratorPage(search.page),
      studentPage: parseStudentPage(search.page),
      teacherPage: parseTeacherPage(search.page),
      view,
    };
  },
  component: DashboardIndexPage,
});

function DashboardIndexPage() {
  const { administratorPage, studentPage, teacherPage, view } = Route.useSearch();

  return (
    <DashboardView
      administratorPage={administratorPage}
      role={view}
      studentPage={studentPage}
      teacherPage={teacherPage}
    />
  );
}

function parseDashboardRole(view: unknown): DashboardRole {
  if (view === "administrator" || view === "teacher" || view === "student") {
    return view;
  }

  return "student";
}

function parseAdministratorPage(page: unknown): AdministratorPage {
  if (page === "home" || page === "students" || page === "disciplines" || page === "analytics") {
    return page;
  }

  return "home";
}

function parseTeacherPage(page: unknown): TeacherPage {
  if (page === "home" || page === "disciplines" || page === "grades" || page === "schedule") {
    return page;
  }

  return "home";
}

function parseStudentPage(page: unknown): StudentPage {
  if (page === "home" || page === "disciplines" || page === "grades" || page === "schedule") {
    return page;
  }

  return "home";
}
