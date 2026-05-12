import { createFileRoute } from "@tanstack/react-router";

import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/students")({
  component: StudentsPage,
});

function StudentsPage() {
  return <DashboardRouteView page="students" />;
}
