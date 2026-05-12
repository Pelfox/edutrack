import { Navigate } from "@tanstack/react-router";

import { DashboardView } from "@/components/dashboard/dashboard-views";
import { useAuth } from "@/lib/context/auth";
import type { DashboardPage } from "@/lib/dashboard-routing";
import {
  getAdministratorPage,
  getDashboardRole,
  getStudentPage,
  getTeacherPage,
  isDashboardPageAvailable,
} from "@/lib/dashboard-routing";

export function DashboardRouteView({ page }: { page: DashboardPage }) {
  const auth = useAuth();
  const role = getDashboardRole(auth.user?.role);

  if (!isDashboardPageAvailable(role, page)) {
    return <Navigate replace={true} to="/" />;
  }

  return (
    <DashboardView
      administratorPage={getAdministratorPage(page)}
      role={role}
      studentPage={getStudentPage(page)}
      teacherPage={getTeacherPage(page)}
    />
  );
}
