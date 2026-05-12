import { createFileRoute } from "@tanstack/react-router";

import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/analytics")({
  component: AnalyticsPage,
});

function AnalyticsPage() {
  return <DashboardRouteView page="analytics" />;
}
