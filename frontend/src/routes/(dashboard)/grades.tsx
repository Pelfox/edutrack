import { createFileRoute } from "@tanstack/react-router";

import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/grades")({
  component: GradesPage,
});

function GradesPage() {
  return <DashboardRouteView page="grades" />;
}
