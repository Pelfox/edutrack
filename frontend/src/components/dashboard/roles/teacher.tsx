import { BookOpen, CalendarDays, Clock3, Users } from "lucide-react";

import type { Metric, ScheduleItem, WorkItem } from "@/components/dashboard/dashboard-widgets";
import {
  DashboardSection,
  MetricsGrid,
  PageHeading,
  ReviewList,
  ScheduleList,
} from "@/components/dashboard/dashboard-widgets";

const scheduleItems: ScheduleItem[] = [
  {
    id: "teacher-finished-class",
    time: "09:00",
    title: "Программирование на Python",
    detail: "ИТ-301 • Ауд. 215",
    status: "Завершено",
    statusVariant: "secondary",
  },
  {
    id: "teacher-current-class",
    time: "09:00",
    title: "Программирование на Python",
    detail: "ИТ-301 • Ауд. 215",
    status: "Идёт",
    statusVariant: "primary",
  },
  {
    id: "teacher-upcoming-class",
    time: "09:00",
    title: "Программирование на Python",
    detail: "ИТ-301 • Ауд. 215",
    status: "Предстоит",
  },
];

const workItems: WorkItem[] = [
  {
    id: "python-lab-1",
    student: "Иванов Иван",
    task: "Лабораторная работа №5",
    subject: "Python",
    time: "2 часа назад",
  },
  {
    id: "python-lab-2",
    student: "Иванов Иван",
    task: "Лабораторная работа №5",
    subject: "Python",
    time: "2 часа назад",
  },
  {
    id: "python-lab-3",
    student: "Иванов Иван",
    task: "Лабораторная работа №5",
    subject: "Python",
    time: "2 часа назад",
  },
];

export function TeacherDashboard() {
  const metrics: Metric[] = [
    {
      title: "Мои курсы",
      value: "5",
      description: "Активных курсов",
      icon: BookOpen,
    },
    {
      title: "Студенты",
      value: "150",
      description: "Всего студентов",
      icon: Users,
    },
    {
      title: "Занятия сегодня",
      value: "3",
      description: "Осталось провести",
      icon: Clock3,
    },
    {
      title: "Непроверенные",
      value: "10",
      description: "Работ на проверке",
      icon: CalendarDays,
    },
  ];

  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="Сегодня четверг, 23 апреля 2026 г." title="Добро пожаловать!" />
      <MetricsGrid metrics={metrics} />
      <div className="grid grid-cols-2 gap-4">
        <DashboardSection
          className="h-[414px]"
          description="Ваши занятия на 23 апреля"
          title="Расписание на сегодня"
        >
          <ScheduleList items={scheduleItems} />
        </DashboardSection>
        <DashboardSection
          className="h-[414px]"
          description="Недавно отправленные студентами"
          title="Непроверенные работы"
        >
          <ReviewList items={workItems} />
        </DashboardSection>
      </div>
    </div>
  );
}
