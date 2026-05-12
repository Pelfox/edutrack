import { createFileRoute } from "@tanstack/react-router";

import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/")({
  component: DashboardIndexPage,
});

function DashboardIndexPage() {
  return <DashboardRouteView page="home" />;
}
