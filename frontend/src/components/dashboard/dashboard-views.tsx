import {
  BarChart3,
  BookOpen,
  CalendarDays,
  GraduationCap,
  House,
  Layers3,
  LibraryBig,
  Users,
} from "lucide-react";

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
import type { AuthProfile, AuthUser } from "@/lib/context/auth";

export function DashboardView({
  role,
  administratorPage,
  onLogout,
  profile,
  studentPage,
  teacherPage,
  user,
}: {
  role: DashboardRole;
  administratorPage: AdministratorPage;
  onLogout?: () => void;
  profile: AuthProfile | null;
  studentPage: StudentPage;
  teacherPage: TeacherPage;
  user: AuthUser | null;
}) {
  return (
    <DashboardShell
      config={getDashboardConfig(
        role,
        administratorPage,
        teacherPage,
        studentPage,
        user,
        profile,
        onLogout,
      )}
    >
      {role === "administrator" && <AdministratorDashboard page={administratorPage} />}
      {role === "teacher" && <TeacherDashboard page={teacherPage} profile={profile} user={user} />}
      {role === "student" && <StudentDashboard page={studentPage} profile={profile} />}
    </DashboardShell>
  );
}

function getDashboardConfig(
  role: DashboardRole,
  administratorPage: AdministratorPage,
  teacherPage: TeacherPage,
  studentPage: StudentPage,
  user: AuthUser | null,
  profile: AuthProfile | null,
  onLogout?: () => void,
): DashboardConfig {
  if (role === "administrator") {
    return {
      navItems: getAdministratorNavItems(administratorPage),
      onLogout,
      searchPlaceholder: "Поиск...",
      user: getDashboardUser(role, user, profile),
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
      onLogout,
      searchPlaceholder: "Поиск студентов, курсов...",
      user: getDashboardUser(role, user, profile),
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
    onLogout,
    searchPlaceholder: "Поиск курсов, материалов...",
    user: getDashboardUser(role, user, profile),
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
      icon: GraduationCap,
      label: "Преподаватели",
      active: activePage === "teachers",
      href: "/teachers",
    },
    {
      icon: Layers3,
      label: "Группы",
      active: activePage === "groups",
      href: "/groups",
    },
    {
      icon: LibraryBig,
      label: "Специальности",
      active: activePage === "specialties",
      href: "/specialties",
    },
    {
      icon: BookOpen,
      label: "Дисциплины",
      active: activePage === "disciplines",
      href: "/disciplines",
    },
    {
      icon: CalendarDays,
      label: "Учебные планы",
      active: activePage === "curriculums",
      href: "/curriculums",
    },
    {
      icon: BarChart3,
      label: "Аналитика",
      active: activePage === "analytics",
      href: "/analytics",
    },
  ];
}

function getDashboardUser(
  role: DashboardRole,
  user: AuthUser | null,
  profile: AuthProfile | null,
): DashboardConfig["user"] {
  const email = profile?.email ?? user?.email ?? "unknown@example.com";
  const fullName = getProfileFullName(profile);
  const initials = getProfileInitials(profile);

  if (role === "administrator") {
    return {
      initials: initials ?? "АД",
      name: fullName ?? "Администратор",
      detail: email,
      nameWeight: "bold",
    };
  }

  if (role === "teacher") {
    return {
      initials: initials ?? "ПР",
      name: fullName ?? "Преподаватель",
      detail: email,
    };
  }

  return {
    initials: initials ?? "СТ",
    name: fullName ?? "Студент",
    detail: email,
  };
}

function getProfileFullName(profile: AuthProfile | null) {
  if (!profile) {
    return null;
  }

  return [profile.last_name, profile.first_name, profile.middle_name].filter(isFilled).join(" ");
}

function getProfileInitials(profile: AuthProfile | null) {
  if (!profile) {
    return null;
  }

  const initials = [profile.last_name, profile.first_name]
    .filter(isFilled)
    .map((name) => name.at(0)?.toLocaleUpperCase("ru-RU"))
    .join("");

  return initials || null;
}

function isFilled(value: string | undefined): value is string {
  return typeof value === "string" && value.length > 0;
}
