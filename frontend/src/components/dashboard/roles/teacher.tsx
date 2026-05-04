import {
  BookOpen,
  CalendarDays,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  Clock3,
  Ellipsis,
  Filter,
  MapPin,
  Plus,
  Search,
  Users,
} from "lucide-react";
import type { DashboardDialogConfig } from "@/components/dashboard/dashboard-dialogs";
import { DashboardActionDialog } from "@/components/dashboard/dashboard-dialogs";
import type { Metric, ScheduleItem, WorkItem } from "@/components/dashboard/dashboard-widgets";
import {
  DashboardSection,
  MetricsGrid,
  PageHeading,
  ReviewList,
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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";

export type TeacherPage = "home" | "disciplines" | "grades" | "schedule";

type TeacherCourse = {
  id: string;
  title: string;
  group: string;
  nextTime: string;
  nextRoom: string;
  students: number;
  materials: number;
  assignments: number;
};

type SimpleStat = {
  label: string;
  value: string;
  detail?: string;
  tone?: "default" | "warning";
};

type WeeklyScheduleDay = {
  day: string;
  count: string;
  lessons: WeeklyLesson[];
};

type WeeklyLesson = {
  id: string;
  title: string;
  group: string;
  kind: string;
  time: string;
  room: string;
};

const scheduleItems: ScheduleItem[] = [
  {
    id: "teacher-finished-class",
    time: "09:00",
    title: "Программирование на Python",
    meta: [
      { icon: Users, label: "ИТ-301" },
      { icon: MapPin, label: "Ауд. 215" },
    ],
    status: "Завершено",
    statusVariant: "secondary",
  },
  {
    id: "teacher-current-class",
    time: "09:00",
    title: "Программирование на Python",
    meta: [
      { icon: Users, label: "ИТ-301" },
      { icon: MapPin, label: "Ауд. 215" },
    ],
    status: "Идёт",
    statusVariant: "primary",
  },
  {
    id: "teacher-upcoming-class",
    time: "09:00",
    title: "Программирование на Python",
    meta: [
      { icon: Users, label: "ИТ-301" },
      { icon: MapPin, label: "Ауд. 215" },
    ],
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

const teacherCourses: TeacherCourse[] = [
  {
    id: "web-development",
    title: "Веб-разработка",
    group: "ИТ-302",
    nextTime: "Сегодня, 11:00",
    nextRoom: "Ауд. 318",
    students: 40,
    materials: 20,
    assignments: 10,
  },
];

const teacherGradeStats: SimpleStat[] = [
  { label: "Всего студентов", value: "150" },
  { label: "Средний балл", value: "4.5" },
  { label: "Непроверенных работ", value: "10", tone: "warning" },
];

const teacherScheduleStats: SimpleStat[] = [
  { label: "Часов в неделю", value: "20" },
  { label: "Занятий в неделю", value: "10" },
  { label: "Групп", value: "5" },
  { label: "Следующее занятие", value: "Через 2 ч" },
];

const createMaterialDialog: DashboardDialogConfig = {
  title: "Создать материал",
  description: "Добавьте учебный материал к одному из ваших курсов.",
  submitLabel: "Создать",
  fields: [
    { label: "Название", name: "material-title", placeholder: "Лекция 6. Компоненты React" },
    { label: "Курс", name: "material-course", placeholder: "Веб-разработка" },
    { label: "Тип материала", name: "material-type", placeholder: "Лекция, презентация, файл" },
    { label: "Ссылка", name: "material-url", placeholder: "https://example.com/material" },
  ],
};

const createGradeDialog: DashboardDialogConfig = {
  title: "Новая оценка",
  description: "Выставьте оценку студенту по выбранному курсу.",
  submitLabel: "Сохранить",
  fields: [
    { label: "Студент", name: "grade-student", placeholder: "Иванов Иван" },
    { label: "Курс", name: "grade-course", placeholder: "Программирование на Python" },
    { label: "Работа", name: "grade-task", placeholder: "Лабораторная работа №5" },
    { label: "Оценка", max: 5, min: 2, name: "grade-value", placeholder: "5", type: "number" },
  ],
};

const createLessonDialog: DashboardDialogConfig = {
  title: "Добавить занятие",
  description: "Запланируйте занятие в расписании.",
  submitLabel: "Добавить",
  fields: [
    { label: "Дисциплина", name: "lesson-title", placeholder: "Программирование на Python" },
    { label: "Группа", name: "lesson-group", placeholder: "ИТ-301" },
    { label: "Дата", name: "lesson-date", placeholder: "2026-04-23", type: "date" },
    { label: "Время", name: "lesson-time", placeholder: "09:00", type: "time" },
    { label: "Аудитория", name: "lesson-room", placeholder: "Ауд. 215" },
  ],
};

const weeklySchedule: WeeklyScheduleDay[] = [
  {
    day: "Понедельник",
    count: "2 занятий",
    lessons: [
      {
        id: "monday-python",
        title: "Программирование на Python",
        group: "ИТ-301",
        kind: "Лекция",
        time: "09:00 - 10:30",
        room: "Ауд. 215",
      },
      {
        id: "monday-algorithms",
        title: "Алгоритмы и структуры данных",
        group: "ИТ-302",
        kind: "Практика",
        time: "11:00 - 12:30",
        room: "Ауд. 205",
      },
    ],
  },
  {
    day: "Вторник",
    count: "2 занятий",
    lessons: [
      {
        id: "tuesday-databases",
        title: "Базы данных",
        group: "ИТ-303",
        kind: "Лекция",
        time: "09:00 - 10:30",
        room: "Ауд. 412",
      },
      {
        id: "tuesday-web",
        title: "Веб-разработка",
        group: "ИТ-302",
        kind: "Лабораторная",
        time: "14:00 - 15:30",
        room: "Ауд. 318",
      },
    ],
  },
  {
    day: "Среда",
    count: "2 занятий",
    lessons: [
      {
        id: "wednesday-python",
        title: "Программирование на Python",
        group: "ИТ-303",
        kind: "Лекция",
        time: "09:00 - 10:30",
        room: "Ауд. 215",
      },
      {
        id: "wednesday-web",
        title: "Веб-разработка",
        group: "ИТ-302",
        kind: "Практика",
        time: "11:00 - 12:30",
        room: "Ауд. 318",
      },
    ],
  },
  {
    day: "Четверг",
    count: "2 занятий",
    lessons: [
      {
        id: "thursday-algorithms",
        title: "Алгоритмы и структуры данных",
        group: "ИТ-301",
        kind: "Лекция",
        time: "11:00 - 12:30",
        room: "Ауд. 205",
      },
      {
        id: "thursday-databases",
        title: "Базы данных",
        group: "ИТ-303",
        kind: "Практика",
        time: "14:00 - 15:30",
        room: "Ауд. 412",
      },
    ],
  },
  {
    day: "Пятница",
    count: "1 занятие",
    lessons: [
      {
        id: "friday-python",
        title: "Программирование на Python",
        group: "ИТ-301",
        kind: "Практика",
        time: "09:00 - 10:30",
        room: "Ауд. 215",
      },
    ],
  },
  {
    day: "Суббота",
    count: "0 занятий",
    lessons: [],
  },
];

export function TeacherDashboard({ page }: { page: TeacherPage }) {
  if (page === "disciplines") {
    return <TeacherDisciplinesPage />;
  }

  if (page === "grades") {
    return <TeacherGradesPage />;
  }

  if (page === "schedule") {
    return <TeacherSchedulePage />;
  }

  return <TeacherHomePage />;
}

function TeacherHomePage() {
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

function TeacherDisciplinesPage() {
  return (
    <div className="flex flex-col gap-6 pb-4">
      <TeacherPageHeader
        actionLabel="Создать материал"
        dialog={createMaterialDialog}
        description="Управление вашими курсами и учебными материалами"
        title="Мои дисциплины"
      />
      <TeacherTabs
        items={[
          { label: "Все курсы", value: "all" },
          { label: "Активные", value: "active" },
          { label: "Материалы", value: "materials" },
          { label: "Задания", value: "tasks" },
        ]}
        value="all"
      />
      <div className="flex flex-col gap-4">
        {teacherCourses.map((course) => (
          <TeacherCourseCard course={course} key={course.id} />
        ))}
      </div>
    </div>
  );
}

function TeacherGradesPage() {
  return (
    <div className="flex flex-col gap-6">
      <TeacherPageHeader
        actionLabel="Новая оценка"
        actionVariant="outline"
        dialog={createGradeDialog}
        description="Управление оценками студентов по вашим курсам"
        title="Оценки"
      />
      <TeacherGradeFilters />
      <SimpleStatsGrid columns={4} stats={teacherGradeStats} />
      <p className="text-base font-medium leading-5 text-muted-foreground">
        Воспользуйтесь поиском для получения сведений о студенте.
      </p>
    </div>
  );
}

function TeacherSchedulePage() {
  return (
    <div className="flex flex-col gap-6 pb-4">
      <TeacherPageHeader
        actionLabel="Добавить занятие"
        dialog={createLessonDialog}
        description="Ваше расписание занятий и мероприятий"
        title="Расписание"
      />
      <SimpleStatsGrid cardClassName="h-36" stats={teacherScheduleStats} />
      <TeacherTabs
        items={[
          { label: "Неделя", value: "week" },
          { label: "Предстоящие", value: "upcoming" },
          { label: "Календарь", value: "calendar" },
        ]}
        value="week"
      />
      <WeeklyScheduleCard />
    </div>
  );
}

function TeacherPageHeader({
  title,
  description,
  actionLabel,
  actionVariant = "default",
  dialog,
}: {
  title: string;
  description: string;
  actionLabel: string;
  actionVariant?: "default" | "outline";
  dialog: DashboardDialogConfig;
}) {
  return (
    <div className="flex h-16 items-center justify-between">
      <div className="flex flex-col gap-1">
        <h1 className="h-9 text-[30px] font-semibold leading-9 text-foreground">{title}</h1>
        <p className="h-6 text-base leading-6 text-muted-foreground">{description}</p>
      </div>
      <DashboardActionDialog
        config={dialog}
        trigger={
          <Button className="h-9 rounded-lg px-3" variant={actionVariant}>
            {actionVariant === "default" && <Plus className="size-4" />}
            {actionLabel}
          </Button>
        }
      />
    </div>
  );
}

function TeacherTabs({
  items,
  value,
}: {
  items: { label: string; value: string }[];
  value: string;
}) {
  return (
    <Tabs value={value}>
      <TabsList className="h-9 rounded-[14px] bg-muted p-[3px]">
        {items.map((item) => (
          <TabsTrigger
            className="h-[29px] rounded-[14px] px-[9px] text-sm font-medium"
            key={item.value}
            value={item.value}
          >
            {item.label}
          </TabsTrigger>
        ))}
      </TabsList>
    </Tabs>
  );
}

function TeacherCourseCard({ course }: { course: TeacherCourse }) {
  return (
    <Card className="h-[272px] gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="h-[90px] gap-2 px-6 pb-0 pt-6">
        <div className="flex items-center gap-3">
          <CardTitle className="text-xl font-medium leading-7 text-card-foreground">
            {course.title}
          </CardTitle>
          <Badge className="h-[22px] rounded-lg px-[9px]" variant="outline">
            {course.group}
          </Badge>
        </div>
        <CardDescription className="flex items-center gap-3 text-base leading-6 text-muted-foreground">
          <span>Следующее занятие:</span>
          <span className="flex items-center gap-1">
            <CalendarDays className="size-4" />
            {course.nextTime}
          </span>
          <span className="flex items-center gap-1">
            <MapPin className="size-4" />
            {course.nextRoom}
          </span>
        </CardDescription>
        <CardAction>
          <TeacherActionsMenu label={`Действия для курса ${course.title}`} />
        </CardAction>
      </CardHeader>
      <CardContent className="grid grid-cols-[1fr_1fr] gap-6 px-6 pb-6 pt-6">
        <div className="flex flex-col gap-9">
          <div className="grid grid-cols-3 gap-4">
            <CourseMetric label="Студентов" value={course.students} />
            <CourseMetric label="Материалов" value={course.materials} />
            <CourseMetric label="Заданий" value={course.assignments} />
          </div>
          <div className="grid grid-cols-2 gap-2">
            <Button className="h-9 rounded-lg">Открыть курс</Button>
            <Button className="h-9 rounded-lg" variant="outline">
              Оценки
            </Button>
          </div>
        </div>
        <div />
      </CardContent>
    </Card>
  );
}

function CourseMetric({ label, value }: { label: string; value: number }) {
  return (
    <div className="flex flex-col gap-1">
      <div className="text-2xl font-bold leading-8 text-card-foreground">{value}</div>
      <div className="text-xs leading-4 text-muted-foreground">{label}</div>
    </div>
  );
}

function TeacherGradeFilters() {
  return (
    <div className="flex h-9 gap-3">
      <Button
        className="h-9 w-[300px] justify-between rounded-lg bg-muted px-[13px]"
        variant="ghost"
      >
        <span>Программирование на Python</span>
        <ChevronDown className="size-4" />
      </Button>
      <div className="relative min-w-0 flex-1">
        <Search className="-translate-y-1/2 absolute left-3 top-1/2 size-4 text-muted-foreground" />
        <Input
          aria-label="Поиск студентов"
          className="h-9 border-transparent bg-muted pl-9 text-sm shadow-none placeholder:text-muted-foreground focus-visible:ring-0"
          placeholder="Поиск студентов..."
        />
      </div>
      <Button className="h-9 rounded-lg px-[13px]" variant="outline">
        <Filter className="size-4" />
        Фильтры
      </Button>
    </div>
  );
}

function SimpleStatsGrid({
  stats,
  cardClassName = "h-40",
  columns = stats.length,
}: {
  stats: SimpleStat[];
  cardClassName?: string;
  columns?: number;
}) {
  return (
    <div
      className="grid w-full gap-4"
      style={{ gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))` }}
    >
      {stats.map((stat) => (
        <Card
          className={`${cardClassName} flex items-start justify-between rounded-[14px] ring-border`}
          key={stat.label}
        >
          <CardHeader className="h-[62px] w-full">
            <CardTitle className="text-sm font-medium leading-5 text-muted-foreground">
              {stat.label}
            </CardTitle>
          </CardHeader>
          <CardContent className="">
            <div
              className={`text-2xl font-bold leading-8 ${
                stat.tone === "warning" ? "text-orange-600" : "text-card-foreground"
              }`}
            >
              {stat.value}
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

function WeeklyScheduleCard() {
  return (
    <Card className="gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="h-[70px] px-6 pb-0 pt-6">
        <div className="flex h-10 items-start justify-between">
          <div className="flex flex-col">
            <CardTitle className="h-4 text-base font-medium leading-4 text-card-foreground">
              Недельное расписание
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
        {weeklySchedule.map((day) => (
          <ScheduleDaySection day={day} key={day.day} />
        ))}
      </CardContent>
    </Card>
  );
}

function ScheduleDaySection({ day }: { day: WeeklyScheduleDay }) {
  return (
    <section className="flex flex-col gap-3">
      <div className="flex h-[27px] items-center gap-3">
        <h3 className="text-lg font-semibold leading-[27px] text-foreground">{day.day}</h3>
        <div className="h-px flex-1 bg-border" />
        <span className="text-sm leading-5 text-muted-foreground">{day.count}</span>
      </div>
      {day.lessons.length === 0 ? (
        <div className="flex h-[54px] items-center justify-center rounded-[10px] border border-dashed border-border text-sm text-muted-foreground">
          Нет занятий
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-3">
          {day.lessons.map((lesson) => (
            <LessonCard key={lesson.id} lesson={lesson} />
          ))}
        </div>
      )}
    </section>
  );
}

function LessonCard({ lesson }: { lesson: WeeklyLesson }) {
  return (
    <div className="flex h-[106px] flex-col gap-2 rounded-[10px] border border-border px-[17px] pb-px pt-[17px]">
      <div className="flex h-11 items-start justify-between gap-4">
        <div className="min-w-0">
          <h4 className="truncate text-base font-medium leading-6 text-foreground">
            {lesson.title}
          </h4>
          <p className="text-sm leading-5 text-muted-foreground">{lesson.group}</p>
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
  );
}

function TeacherActionsMenu({ label }: { label: string }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild={true}>
        <Button aria-label={label} className="size-9 rounded-lg" size="icon" variant="ghost">
          <Ellipsis className="size-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-44">
        <DropdownMenuItem>Открыть</DropdownMenuItem>
        <DropdownMenuItem>Редактировать</DropdownMenuItem>
        <DropdownMenuItem>Материалы</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem variant="destructive">Удалить</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
