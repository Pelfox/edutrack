import { BarChart3, BookOpen, CalendarDays, GraduationCap, House, Users } from "lucide-react";

import type {
  DashboardConfig,
  DashboardNavItem,
  DashboardRole,
} from "@/components/dashboard/dashboard-layout";
import { DashboardShell } from "@/components/dashboard/dashboard-layout";
import type { AdministratorPage } from "@/components/dashboard/roles/administrator";
import { AdministratorDashboard } from "@/components/dashboard/roles/administrator";
import type { StudentPage } from "@/components/dashboard/roles/student";
import { StudentDashboard } from "@/components/dashboard/roles/student";
import type { TeacherPage } from "@/components/dashboard/roles/teacher";
import { TeacherDashboard } from "@/components/dashboard/roles/teacher";

export function DashboardView({
  role,
  administratorPage,
  studentPage,
  teacherPage,
}: {
  role: DashboardRole;
  administratorPage: AdministratorPage;
  studentPage: StudentPage;
  teacherPage: TeacherPage;
}) {
  return (
    <DashboardShell config={getDashboardConfig(role, administratorPage, teacherPage, studentPage)}>
      {role === "administrator" && <AdministratorDashboard page={administratorPage} />}
      {role === "teacher" && <TeacherDashboard page={teacherPage} />}
      {role === "student" && <StudentDashboard page={studentPage} />}
    </DashboardShell>
  );
}

function getDashboardConfig(
  role: DashboardRole,
  administratorPage: AdministratorPage,
  teacherPage: TeacherPage,
  studentPage: StudentPage,
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
        {
          icon: House,
          label: "Главная",
          active: teacherPage === "home",
          href: "/",
        },
        {
          icon: BookOpen,
          label: "Мои дисциплины",
          active: teacherPage === "disciplines",
          href: "/disciplines",
        },
        {
          icon: GraduationCap,
          label: "Оценки",
          active: teacherPage === "grades",
          href: "/grades",
        },
        {
          icon: CalendarDays,
          label: "Расписание",
          active: teacherPage === "schedule",
          href: "/schedule",
        },
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
      {
        icon: House,
        label: "Главная",
        active: studentPage === "home",
        href: "/",
      },
      {
        icon: BookOpen,
        label: "Дисциплины",
        active: studentPage === "disciplines",
        href: "/disciplines",
      },
      {
        icon: GraduationCap,
        label: "Оценки",
        active: studentPage === "grades",
        href: "/grades",
      },
      {
        icon: CalendarDays,
        label: "Расписание",
        active: studentPage === "schedule",
        href: "/schedule",
      },
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
      href: "/",
    },
    {
      icon: Users,
      label: "Студенты",
      active: activePage === "students",
      href: "/students",
    },
    {
      icon: BookOpen,
      label: "Дисциплины",
      active: activePage === "disciplines",
      href: "/disciplines",
    },
    {
      icon: BarChart3,
      label: "Аналитика",
      active: activePage === "analytics",
      href: "/analytics",
    },
  ];
}
