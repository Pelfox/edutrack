import { createFileRoute } from "@tanstack/react-router";

import type { DashboardRole } from "@/components/dashboard/dashboard-layout";
import { DashboardView } from "@/components/dashboard/dashboard-views";
import type { AdministratorPage } from "@/components/dashboard/roles/administrator";

export const Route = createFileRoute("/(dashboard)/")({
  validateSearch: (search): { page: AdministratorPage; view: DashboardRole } => ({
    page: parseAdministratorPage(search.page),
    view: parseDashboardRole(search.view),
  }),
  component: DashboardIndexPage,
});

function DashboardIndexPage() {
  const { page, view } = Route.useSearch();

  return <DashboardView administratorPage={page} role={view} />;
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
