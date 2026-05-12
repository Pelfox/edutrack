import { createFileRoute } from "@tanstack/react-router";
import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/teachers")({
  component: TeachersRoute,
});

function TeachersRoute() {
  return <DashboardRouteView page="teachers" />;
}
