import { createFileRoute } from "@tanstack/react-router";
import { DashboardRouteView } from "@/components/dashboard/dashboard-route-view";

export const Route = createFileRoute("/(dashboard)/groups")({
  component: GroupsRoute,
});

function GroupsRoute() {
  return <DashboardRouteView page="groups" />;
}
