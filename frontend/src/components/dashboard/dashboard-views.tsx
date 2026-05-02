import { BarChart3, BookOpen, CalendarDays, GraduationCap, House, Users } from "lucide-react";

import type {
  DashboardConfig,
  DashboardNavItem,
  DashboardRole,
} from "@/components/dashboard/dashboard-layout";
import { DashboardShell } from "@/components/dashboard/dashboard-layout";
import type { AdministratorPage } from "@/components/dashboard/roles/administrator";
import { AdministratorDashboard } from "@/components/dashboard/roles/administrator";
import { StudentDashboard } from "@/components/dashboard/roles/student";
import { TeacherDashboard } from "@/components/dashboard/roles/teacher";

export function DashboardView({
  role,
  administratorPage,
}: {
  role: DashboardRole;
  administratorPage: AdministratorPage;
}) {
  return (
    <DashboardShell config={getDashboardConfig(role, administratorPage)}>
      {role === "administrator" && <AdministratorDashboard page={administratorPage} />}
      {role === "teacher" && <TeacherDashboard />}
      {role === "student" && <StudentDashboard />}
    </DashboardShell>
  );
}

function getDashboardConfig(
  role: DashboardRole,
  administratorPage: AdministratorPage,
): DashboardConfig {
  if (role === "administrator") {
    return {
      navItems: getAdministratorNavItems(administratorPage),
      searchPlaceholder: "Поиск...",
      user: {
        initials: "АД",
        name: "Администратор",
        detail: "someone@example.com",
        nameWeight: "bold",
      },
    };
  }

  if (role === "teacher") {
    return {
      navItems: [
        { icon: House, label: "Главная", active: true, href: "/?view=teacher" },
        { icon: BookOpen, label: "Мои дисциплины" },
        { icon: GraduationCap, label: "Оценки" },
        { icon: CalendarDays, label: "Расписание" },
      ],
      searchPlaceholder: "Поиск студентов, курсов...",
      user: {
        initials: "ИП",
        name: "Иванов П.С.",
        detail: "someone@example.com",
      },
    };
  }

  return {
    navItems: [
      { icon: House, label: "Главная", active: true, href: "/?view=student" },
      { icon: BookOpen, label: "Дисциплины" },
      { icon: GraduationCap, label: "Оценки" },
      { icon: CalendarDays, label: "Расписание" },
    ],
    searchPlaceholder: "Поиск курсов, материалов...",
    user: {
      initials: "ИИ",
      name: "Иванов Иван",
      detail: "ИТ-301",
    },
  };
}

function getAdministratorNavItems(activePage: AdministratorPage): DashboardNavItem[] {
  return [
    {
      icon: House,
      label: "Главная",
      active: activePage === "home",
      href: "/?view=administrator&page=home",
    },
    {
      icon: Users,
      label: "Студенты",
      active: activePage === "students",
      href: "/?view=administrator&page=students",
    },
    {
      icon: BookOpen,
      label: "Дисциплины",
      active: activePage === "disciplines",
      href: "/?view=administrator&page=disciplines",
    },
    {
      icon: BarChart3,
      label: "Аналитика",
      active: activePage === "analytics",
      href: "/?view=administrator&page=analytics",
    },
  ];
}
