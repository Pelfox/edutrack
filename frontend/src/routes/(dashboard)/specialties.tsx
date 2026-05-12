import { createFileRoute } from "@tanstack/react-router";
import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/specialties")({
  component: SpecialtiesRoute,
});

function SpecialtiesRoute() {
  return <DashboardRouteView page="specialties" />;
}
