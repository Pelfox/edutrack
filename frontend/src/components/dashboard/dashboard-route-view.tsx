import { Navigate, useRouter } from "@tanstack/react-router";

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
  const router = useRouter();
  const role = getDashboardRole(auth.user?.role);

  function handleLogout() {
    auth.logout();
    router.history.push("/login");
  }

  if (!isDashboardPageAvailable(role, page)) {
    return <Navigate replace={true} to="/" />;
  }

  return (
    <DashboardView
      administratorPage={getAdministratorPage(page)}
      onLogout={handleLogout}
      role={role}
      studentPage={getStudentPage(page)}
      teacherPage={getTeacherPage(page)}
      profile={auth.profile}
      user={auth.user}
    />
  );
}
