import type { LucideIcon } from "lucide-react";
import {
  Award,
  BarChart3,
  BookOpen,
  ChevronLeft,
  ChevronRight,
  Clock3,
  Download,
  FileText,
  GraduationCap,
  MapPin,
  NotebookTabs,
  TrendingUp,
} from "lucide-react";

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
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";

export type StudentPage = "home" | "disciplines" | "grades" | "schedule";

type StudentCourse = {
  id: string;
  title: string;
  teacher: string;
  average: string;
  materials: StudentMaterial[];
};

type StudentMaterial = {
  id: string;
  title: string;
  kind: string;
  date: string;
};

type StudentGradeSummary = {
  title: string;
  value: string;
  detail: string;
  icon?: LucideIcon;
  tone?: "positive" | "default";
};

type StudentGradeCourse = {
  title: string;
  teacher: string;
  average: string;
  rows: StudentGradeRow[];
};

type StudentGradeRow = {
  assignment: string;
  weight: string;
  grade: string;
  tone?: "positive" | "blue";
};

type StudentScheduleDay = {
  day: string;
  count: string;
  lessons: StudentLesson[];
};

type StudentLesson = {
  id: string;
  start: string;
  title: string;
  teacher: string;
  kind: string;
  time: string;
  room: string;
};

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

const studentCourses: StudentCourse[] = [
  {
    id: "python",
    title: "Программирование на Python",
    teacher: "Доц. Иванов П.С.",
    average: "4.5",
    materials: [
      {
        id: "python-lecture-1",
        title: "Лекция 1: Введение в Python",
        kind: "PDF",
        date: "15.02.2026",
      },
      { id: "python-lecture-2", title: "Лекция 2: Типы данных", kind: "PDF", date: "20.02.2026" },
      {
        id: "python-practice-1",
        title: "Практика 1: Основы синтаксиса",
        kind: "DOCX",
        date: "22.02.2026",
      },
    ],
  },
];

const studentGradeSummary: StudentGradeSummary[] = [
  {
    title: "Средний балл",
    value: "4.5",
    detail: "Отличный результат",
    icon: Award,
    tone: "positive",
  },
  { title: "Пятёрок", value: "10", detail: "65% от всех оценок", icon: TrendingUp },
  { title: "Четвёрок", value: "5", detail: "35% от всех оценок", icon: BookOpen },
  { title: "Всего оценок", value: "15", detail: "За семестр" },
];

const studentGradeCourses: StudentGradeCourse[] = [
  {
    title: "Программирование на Python",
    teacher: "Доц. Иванов П.С.",
    average: "4.5",
    rows: [
      { assignment: "ЛР №1", weight: "10%", grade: "5", tone: "positive" },
      { assignment: "ЛР №2", weight: "10%", grade: "4", tone: "blue" },
      { assignment: "ДЗ №1", weight: "5%", grade: "5", tone: "positive" },
      { assignment: "Промежуточная", weight: "25%", grade: "4", tone: "blue" },
    ],
  },
];

const studentSchedule: StudentScheduleDay[] = [
  {
    day: "Понедельник",
    count: "3 пары",
    lessons: [
      {
        id: "monday-python-1",
        start: "11:00",
        title: "Программирование на Python",
        teacher: "Доц. Иванов П.С.",
        kind: "Практика",
        time: "11:00 - 12:30",
        room: "Ауд. 318",
      },
      {
        id: "monday-python-2",
        start: "11:00",
        title: "Программирование на Python",
        teacher: "Доц. Иванов П.С.",
        kind: "Практика",
        time: "11:00 - 12:30",
        room: "Ауд. 318",
      },
      {
        id: "monday-python-3",
        start: "11:00",
        title: "Программирование на Python",
        teacher: "Доц. Иванов П.С.",
        kind: "Практика",
        time: "11:00 - 12:30",
        room: "Ауд. 318",
      },
    ],
  },
];

export function StudentDashboard({ page }: { page: StudentPage }) {
  if (page === "disciplines") {
    return <StudentCoursesPage />;
  }

  if (page === "grades") {
    return <StudentGradesPage />;
  }

  if (page === "schedule") {
    return <StudentSchedulePage />;
  }

  return <StudentHomePage />;
}

function StudentHomePage() {
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
      <div className="grid grid-cols-2 items-start gap-4">
        <DashboardSection
          className="h-[390px]"
          description="Ваши занятия"
          title="Расписание на сегодня"
        >
          <ScheduleList items={scheduleItems} />
        </DashboardSection>
        <DashboardSection description="Недавно выставленные" title="Последние оценки">
          <GradesList items={recentGrades} />
        </DashboardSection>
      </div>
      <DashboardSection className="min-h-[454px]" description="Прогресс обучения" title="Мои курсы">
        <CourseGrid items={courseItems} />
      </DashboardSection>
    </div>
  );
}

function StudentCoursesPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="Ваши курсы и учебные материалы" title="Мои курсы" />
      <div className="flex flex-col gap-4">
        {studentCourses.map((course) => (
          <StudentCourseCard course={course} key={course.id} />
        ))}
      </div>
    </div>
  );
}

function StudentGradesPage() {
  return (
    <div className="flex flex-col gap-6 pb-4">
      <PageHeading description="Ваша успеваемость по всем курсам" title="Мои оценки" />
      <div className="grid grid-cols-4 gap-4">
        {studentGradeSummary.map((item) => (
          <StudentGradeSummaryCard item={item} key={item.title} />
        ))}
      </div>
      <div className="flex flex-col gap-4">
        {studentGradeCourses.map((course) => (
          <StudentGradeCourseCard course={course} key={course.title} />
        ))}
      </div>
    </div>
  );
}

function StudentSchedulePage() {
  return (
    <div className="flex flex-col gap-6 pb-4">
      <PageHeading description="Ваше расписание занятий на неделю" title="Расписание" />
      <StudentScheduleCard />
    </div>
  );
}

function StudentTabs() {
  return (
    <Tabs defaultValue="materials">
      <TabsList className="h-9 rounded-[14px] bg-muted p-[3px]">
        <TabsTrigger className="h-[29px] rounded-[14px] px-2 text-sm font-medium" value="materials">
          <FileText className="size-4" />
          Материалы
        </TabsTrigger>
        <TabsTrigger className="h-[29px] rounded-[14px] px-2 text-sm font-medium" value="tasks">
          <NotebookTabs className="size-4" />
          Задания
        </TabsTrigger>
      </TabsList>
    </Tabs>
  );
}

function StudentCourseCard({ course }: { course: StudentCourse }) {
  return (
    <Card className="gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="h-[86px] px-6 pb-0 pt-6">
        <div className="flex min-w-0 flex-col gap-1">
          <CardTitle className="text-xl font-medium leading-7 text-card-foreground">
            {course.title}
          </CardTitle>
          <CardDescription className="text-base leading-6 text-muted-foreground">
            {course.teacher}
          </CardDescription>
        </div>
        <CardAction>
          <div className="flex size-14 items-center justify-center rounded-[10px] bg-green-100 text-xl font-bold leading-7 text-green-700">
            {course.average}
          </div>
        </CardAction>
      </CardHeader>
      <CardContent className="flex flex-col gap-6 px-6 pb-6 pt-6">
        <StudentTabs />
        <div className="flex flex-col gap-2">
          {course.materials.map((material) => (
            <MaterialRow key={material.id} material={material} />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function MaterialRow({ material }: { material: StudentMaterial }) {
  return (
    <div className="flex h-[66px] items-center justify-between rounded-[10px] border border-border p-[13px]">
      <div className="flex min-w-0 items-center gap-3">
        <div className="flex size-10 shrink-0 items-center justify-center rounded-[10px] bg-blue-100 text-blue-700">
          <FileText className="size-5" />
        </div>
        <div className="min-w-0">
          <div className="truncate text-sm font-medium leading-5 text-foreground">
            {material.title}
          </div>
          <div className="flex items-center gap-2 text-xs leading-4 text-muted-foreground">
            <span>{material.kind}</span>
            <span>{material.date}</span>
          </div>
        </div>
      </div>
      <Button
        aria-label={`Скачать ${material.title}`}
        className="h-8 w-9 rounded-lg"
        variant="ghost"
      >
        <Download className="size-4" />
      </Button>
    </div>
  );
}

function StudentGradeSummaryCard({ item }: { item: StudentGradeSummary }) {
  const Icon = item.icon;

  return (
    <Card className="h-[150px] gap-6 rounded-[14px] py-0 ring-border">
      <CardHeader className="h-[52px] px-6 pb-2 pt-6">
        <CardTitle className="text-sm font-medium leading-5 text-card-foreground">
          {item.title}
        </CardTitle>
        {Icon !== undefined && (
          <CardAction>
            <Icon className="size-4 text-muted-foreground" />
          </CardAction>
        )}
      </CardHeader>
      <CardContent className="px-6">
        <div className="text-2xl font-bold leading-8 text-card-foreground">{item.value}</div>
        <div
          className={`text-xs leading-4 ${
            item.tone === "positive" ? "text-green-600" : "text-muted-foreground"
          }`}
        >
          {item.detail}
        </div>
      </CardContent>
    </Card>
  );
}

function StudentGradeCourseCard({ course }: { course: StudentGradeCourse }) {
  return (
    <Card className="gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="h-[94px] px-6 pb-0 pt-6">
        <div className="flex min-w-0 flex-col gap-1">
          <CardTitle className="text-lg font-medium leading-7 text-card-foreground">
            {course.title}
          </CardTitle>
          <CardDescription className="text-base leading-6 text-muted-foreground">
            {course.teacher}
          </CardDescription>
        </div>
        <CardAction>
          <div className="flex size-16 items-center justify-center rounded-[10px] bg-green-100 text-xl font-bold leading-7 text-green-700">
            {course.average}
          </div>
        </CardAction>
      </CardHeader>
      <CardContent className="px-6 pb-6">
        <div className="overflow-hidden rounded-lg border border-border">
          <Table>
            <TableHeader>
              <TableRow className="h-10 hover:bg-transparent">
                <TableHead className="w-1/2">Задание</TableHead>
                <TableHead className="w-[24%] text-center">Вес</TableHead>
                <TableHead className="w-[26%] text-center">Оценка</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {course.rows.map((row) => (
                <TableRow className="h-[57px]" key={row.assignment}>
                  <TableCell className="font-medium">{row.assignment}</TableCell>
                  <TableCell className="text-center">
                    <Badge className="h-[22px] rounded-lg px-[9px]" variant="outline">
                      {row.weight}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex justify-center">
                      <GradeBox grade={row.grade} tone={row.tone} />
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}

function GradeBox({
  grade,
  tone = "positive",
}: {
  grade: string;
  tone?: "positive" | "blue" | undefined;
}) {
  return (
    <div
      className={`flex size-10 items-center justify-center rounded-lg text-sm font-bold ${
        tone === "blue" ? "bg-blue-100 text-blue-700" : "bg-green-100 text-green-700"
      }`}
    >
      {grade}
    </div>
  );
}

function StudentScheduleCard() {
  return (
    <Card className="min-h-[840px] gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="h-[70px] px-6 pb-0 pt-6">
        <div className="flex h-10 items-start justify-between">
          <div className="flex flex-col">
            <CardTitle className="h-4 text-base font-medium leading-4 text-card-foreground">
              Неделя
            </CardTitle>
            <CardDescription className="h-6 text-base leading-6 text-muted-foreground">
              23 - 29 апреля 2026
            </CardDescription>
          </div>
          <div className="flex gap-2">
            <Button
              aria-label="Предыдущая неделя"
              className="size-9 rounded-lg"
              size="icon"
              variant="outline"
            >
              <ChevronLeft className="size-4" />
            </Button>
            <Button
              aria-label="Следующая неделя"
              className="size-9 rounded-lg"
              size="icon"
              variant="outline"
            >
              <ChevronRight className="size-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex flex-col gap-6 px-6 pb-6 pt-6">
        {studentSchedule.map((day) => (
          <StudentScheduleDaySection day={day} key={day.day} />
        ))}
      </CardContent>
    </Card>
  );
}

function StudentScheduleDaySection({ day }: { day: StudentScheduleDay }) {
  return (
    <section className="flex flex-col gap-3">
      <div className="flex h-[27px] items-center gap-3">
        <h3 className="text-lg font-semibold leading-[27px] text-foreground">{day.day}</h3>
        <div className="h-px flex-1 bg-border" />
        <span className="text-sm leading-5 text-muted-foreground">{day.count}</span>
      </div>
      <div className="flex flex-col gap-3">
        {day.lessons.map((lesson) => (
          <StudentLessonCard key={lesson.id} lesson={lesson} />
        ))}
      </div>
    </section>
  );
}

function StudentLessonCard({ lesson }: { lesson: StudentLesson }) {
  return (
    <div className="flex h-[110px] rounded-[10px] border border-border px-[17px] pb-px pt-[17px]">
      <div className="flex h-[76px] w-full items-start gap-4">
        <div className="flex h-16 w-20 shrink-0 flex-col items-center justify-center rounded-[10px] bg-muted py-2.5">
          <span className="text-xs leading-4 text-muted-foreground">Начало</span>
          <span className="text-lg font-bold leading-7 text-foreground">{lesson.start}</span>
        </div>
        <div className="flex h-[76px] min-w-0 flex-1 flex-col gap-2">
          <div className="flex h-12 items-start justify-between gap-4">
            <div className="min-w-0">
              <h4 className="truncate text-base font-medium leading-6 text-foreground">
                {lesson.title}
              </h4>
              <p className="text-sm leading-5 text-muted-foreground">{lesson.teacher}</p>
            </div>
            <Badge className="h-[22px] shrink-0 rounded-lg px-[9px]" variant="outline">
              {lesson.kind}
            </Badge>
          </div>
          <div className="flex gap-4 text-sm leading-5 text-muted-foreground">
            <span className="flex items-center gap-1">
              <Clock3 className="size-3.5" />
              {lesson.time}
            </span>
            <span className="flex items-center gap-1">
              <MapPin className="size-3.5" />
              {lesson.room}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
