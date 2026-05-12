import { createFileRoute } from "@tanstack/react-router";

import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/schedule")({
  component: SchedulePage,
});

function SchedulePage() {
  return <DashboardRouteView page="schedule" />;
}
