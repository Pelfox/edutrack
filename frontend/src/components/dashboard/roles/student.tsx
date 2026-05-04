import { BarChart3, BookOpen, Clock3, GraduationCap, MapPin } from "lucide-react";

import type {
  CourseItem,
  GradeItem,
  Metric,
  ScheduleItem,
} from "@/components/dashboard/dashboard-widgets";
import {
  CourseGrid,
  DashboardSection,
  GradesList,
  MetricsGrid,
  PageHeading,
  ScheduleList,
} from "@/components/dashboard/dashboard-widgets";

const scheduleItems: ScheduleItem[] = [
  {
    id: "student-finished-class",
    time: "09:00",
    title: "Математический анализ",
    meta: [
      { icon: GraduationCap, label: "Проф. Смирнова А.В." },
      { icon: MapPin, label: "Ауд. 215" },
    ],
    status: "Завершено",
    statusVariant: "secondary",
  },
  {
    id: "student-upcoming-class-1",
    time: "09:00",
    title: "Математический анализ",
    meta: [
      { icon: GraduationCap, label: "Проф. Смирнова А.В." },
      { icon: MapPin, label: "Ауд. 215" },
    ],
    status: "Предстоит",
    statusVariant: "primary",
  },
  {
    id: "student-upcoming-class-2",
    time: "09:00",
    title: "Математический анализ",
    meta: [
      { icon: GraduationCap, label: "Проф. Смирнова А.В." },
      { icon: MapPin, label: "Ауд. 215" },
    ],
    status: "Предстоит",
    statusVariant: "primary",
  },
];

const recentGrades: GradeItem[] = [
  {
    title: "Лабораторная работа №3",
    subject: "Базы данных",
    time: "2 дня назад",
    grade: "5",
  },
  {
    title: "Домашнее задание №5",
    subject: "Python",
    time: "3 дня назад",
    grade: "4",
    tone: "blue",
  },
  {
    title: "Контрольная работа",
    subject: "Мат. анализ",
    time: "5 дней назад",
    grade: "5",
  },
];

const courseItems: CourseItem[] = [
  {
    title: "Программирование на Python",
    teacher: "Доц. Иванов П.С.",
    score: "4.5",
    progress: 45,
    next: "Сегодня, 11:00",
  },
  {
    title: "Математический анализ",
    teacher: "Проф. Смирнова А.В.",
    score: "4.8",
    progress: 67,
    next: "Завтра, 09:00",
  },
  {
    title: "Базы данных",
    teacher: "Доц. Петрова М.И.",
    score: "5.0",
    progress: 72,
    next: "Завтра, 14:00",
  },
  {
    title: "Английский язык",
    teacher: "Преп. Морозова О.И.",
    score: "4.2",
    progress: 55,
    next: "Сегодня, 14:00",
    tone: "blue",
  },
];

export function StudentDashboard() {
  const metrics: Metric[] = [
    {
      title: "Средний балл",
      value: "4.5",
      description: "+0.5 от прошлого месяца",
      icon: GraduationCap,
      tone: "positive",
    },
    {
      title: "Мои курсы",
      value: "10",
      description: "Активных курсов",
      icon: BookOpen,
    },
    {
      title: "Занятия сегодня",
      value: "5",
      description: "Осталось 2",
      icon: Clock3,
    },
    {
      title: "Посещаемость",
      value: "90%",
      description: "Отличный результат",
      icon: BarChart3,
      tone: "positive",
    },
  ];

  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="четверг, 23 апреля 2026 г." title="Привет, Иван!" />
      <MetricsGrid metrics={metrics} />
      <div className="grid grid-cols-2 gap-4">
        <DashboardSection
          className="h-[360px]"
          description="Ваши занятия"
          title="Расписание на сегодня"
        >
          <ScheduleList items={scheduleItems} />
        </DashboardSection>
        <DashboardSection
          className="h-[360px]"
          description="Недавно выставленные"
          title="Последние оценки"
        >
          <GradesList items={recentGrades} />
        </DashboardSection>
      </div>
      <DashboardSection className="min-h-[454px]" description="Прогресс обучения" title="Мои курсы">
        <CourseGrid items={courseItems} />
      </DashboardSection>
    </div>
  );
}
