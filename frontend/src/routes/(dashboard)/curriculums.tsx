import { createFileRoute } from "@tanstack/react-router";
import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/curriculums")({
  component: CurriculumsRoute,
});

function CurriculumsRoute() {
  return <DashboardRouteView page="curriculums" />;
}
