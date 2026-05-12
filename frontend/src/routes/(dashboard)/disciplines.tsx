import { createFileRoute } from "@tanstack/react-router";

import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/disciplines")({
  component: DisciplinesPage,
});

function DisciplinesPage() {
  return <DashboardRouteView page="disciplines" />;
}
