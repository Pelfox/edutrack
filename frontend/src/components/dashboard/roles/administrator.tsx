import {
  BarChart3,
  BookOpen,
  Clock3,
  Mail,
  MoreHorizontal,
  Phone,
  Plus,
  Search,
  Users,
} from "lucide-react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Line,
  LineChart,
  Pie,
  PieChart,
  XAxis,
  YAxis,
} from "recharts";
import type { DashboardDialogConfig } from "@/components/dashboard/dashboard-dialogs";
import { DashboardActionDialog } from "@/components/dashboard/dashboard-dialogs";
import type { ActivityItem, Metric, TopStudent } from "@/components/dashboard/dashboard-widgets";
import {
  ActivityList,
  DashboardSection,
  MetricsGrid,
  PageHeading,
  TopStudentsList,
} from "@/components/dashboard/dashboard-widgets";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
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
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export type AdministratorPage = "home" | "students" | "disciplines" | "analytics";

type StudentRecord = {
  id: string;
  initials: string;
  name: string;
  group: string;
  email: string;
  phone: string;
  average: string;
  status: string;
};

type CourseRecord = {
  id: string;
  title: string;
  teacher: string;
  category: string;
  term: string;
  students: number;
  hours: number;
};

const activityItems: ActivityItem[] = [
  {
    id: "math-grade",
    student: "Иванов Иван",
    subject: "Математический анализ",
    grade: "5",
    time: "2 часа назад",
  },
  {
    id: "programming-grade",
    student: "Иванов Иван",
    subject: "Программирование",
    grade: "4",
    time: "5 часов назад",
  },
  {
    id: "physics-grade",
    student: "Иванов Иван",
    subject: "Физика",
    grade: "5",
    time: "1 день назад",
  },
  {
    id: "english-grade",
    student: "Иванов Иван",
    subject: "Английский язык",
    grade: "3",
    time: "1 день назад",
    tone: "muted",
  },
];

const topStudents: TopStudent[] = [
  { name: "Иванов Иван", group: "ЭФБО-01-24", score: "4.9", progress: 98 },
  { name: "Иванов Иван", group: "ЭФБО-01-24", score: "4.8", progress: 96 },
  { name: "Иванов Иван", group: "ЭФБО-01-24", score: "4.7", progress: 94 },
  { name: "Иванов Иван", group: "ЭФБО-01-24", score: "4.6", progress: 92 },
];

const students: StudentRecord[] = Array.from({ length: 6 }, (_, index) => ({
  id: `student-${index + 1}`,
  initials: "ИИ",
  name: "Иванов Иван Иванович",
  group: "ЭФБО-01-24",
  email: "ivanov@example.com",
  phone: "+7 (999) 123-45-67",
  average: "4.8",
  status: "Активен",
}));

const courses: CourseRecord[] = Array.from({ length: 7 }, (_, index) => ({
  id: `course-${index + 1}`,
  title: "Математический анализ",
  teacher: "Иванов И.И.",
  category: "Математика",
  term: "Весенний 2026",
  students: 45,
  hours: 128,
}));

const performanceTrendData = [
  { average: 4.1, attendance: 82, month: "Окт" },
  { average: 4.2, attendance: 84, month: "Ноя" },
  { average: 4.1, attendance: 81, month: "Дек" },
  { average: 4.3, attendance: 86, month: "Янв" },
  { average: 4.3, attendance: 83, month: "Фев" },
  { average: 4.4, attendance: 88, month: "Мар" },
  { average: 4.5, attendance: 80, month: "Апр" },
];

const gradeDistributionData = [
  { fill: "var(--color-excellent)", grade: "5", value: 45 },
  { fill: "var(--color-good)", grade: "4", value: 35 },
  { fill: "var(--color-satisfactory)", grade: "3", value: 15 },
  { fill: "var(--color-poor)", grade: "2", value: 5 },
];

const subjectPerformanceData = [
  { average: 4.5, subject: "Мат. анализ" },
  { average: 4.3, subject: "Программ." },
  { average: 4.6, subject: "БД" },
  { average: 4.1, subject: "Алгоритмы" },
  { average: 4.4, subject: "Веб-разр." },
  { average: 3.9, subject: "Физика" },
];

const trendChartConfig = {
  attendance: {
    color: "var(--chart-2)",
    label: "Посещаемость",
  },
  average: {
    color: "var(--chart-1)",
    label: "Средний балл",
  },
} satisfies ChartConfig;

const gradeChartConfig = {
  excellent: {
    color: "var(--chart-2)",
    label: "Отлично",
  },
  good: {
    color: "var(--chart-1)",
    label: "Хорошо",
  },
  satisfactory: {
    color: "var(--chart-4)",
    label: "Удовлетворительно",
  },
  poor: {
    color: "var(--destructive)",
    label: "Неудовлетворительно",
  },
} satisfies ChartConfig;

const subjectChartConfig = {
  average: {
    color: "var(--chart-1)",
    label: "Средний балл",
  },
} satisfies ChartConfig;

const addStudentDialog: DashboardDialogConfig = {
  title: "Добавить студента",
  description: "Заполните основные данные студента для добавления в систему.",
  submitLabel: "Добавить",
  fields: [
    { label: "ФИО", name: "student-name", placeholder: "Иванов Иван Иванович" },
    { label: "Группа", name: "student-group", placeholder: "ЭФБО-01-24" },
    { label: "Email", name: "student-email", placeholder: "ivanov@example.com", type: "email" },
    { label: "Телефон", name: "student-phone", placeholder: "+7 (999) 123-45-67" },
  ],
};

const addCourseDialog: DashboardDialogConfig = {
  title: "Добавить курс",
  description: "Создайте новую дисциплину и назначьте преподавателя.",
  submitLabel: "Добавить",
  fields: [
    { label: "Название", name: "course-title", placeholder: "Математический анализ" },
    { label: "Преподаватель", name: "course-teacher", placeholder: "Иванов И.И." },
    { label: "Категория", name: "course-category", placeholder: "Математика" },
    { label: "Семестр", name: "course-term", placeholder: "Весенний 2026" },
  ],
};

export function AdministratorDashboard({ page }: { page: AdministratorPage }) {
  if (page === "students") {
    return <AdministratorStudentsPage />;
  }

  if (page === "disciplines") {
    return <AdministratorDisciplinesPage />;
  }

  if (page === "analytics") {
    return <AdministratorAnalyticsPage />;
  }

  return <AdministratorHomePage />;
}

function AdministratorHomePage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeading
        description="Обзор успеваемости и активности студентов"
        title="Панель управления"
      />
      <MetricsGrid metrics={administratorOverviewMetrics} />
      <div className="grid grid-cols-[1.35fr_1fr] gap-4">
        <DashboardSection
          className="h-[360px]"
          description="Недавно выставленные оценки"
          title="Последняя активность"
        >
          <ActivityList items={activityItems} />
        </DashboardSection>
        <DashboardSection
          className="h-[360px]"
          description="По среднему баллу"
          title="Лучшие студенты"
        >
          <TopStudentsList items={topStudents} />
        </DashboardSection>
      </div>
    </div>
  );
}

function AdministratorStudentsPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        actionLabel="Добавить студента"
        dialog={addStudentDialog}
        description="Управление списком студентов и их данными"
        title="Студенты"
      />
      <SearchAndFilter placeholder="Поиск студентов по имени, email или группе..." />
      <StudentsTable />
    </div>
  );
}

function AdministratorDisciplinesPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        actionLabel="Добавить курс"
        dialog={addCourseDialog}
        description="Управление учебными курсами и программами"
        title="Дисциплины"
      />
      <DisciplineFilters />
      <div className="grid grid-cols-3 gap-2.5">
        {courses.map((course) => (
          <CourseCard course={course} key={course.id} />
        ))}
      </div>
    </div>
  );
}

function AdministratorAnalyticsPage() {
  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="Статистика и анализ успеваемости студентов" title="Аналитика" />
      <MetricsGrid metrics={administratorAnalyticsMetrics} />
      <div className="grid grid-cols-2 gap-4">
        <DashboardSection
          className="h-[420px]"
          description="Средний балл и посещаемость по месяцам"
          title="Динамика успеваемости"
        >
          <LineChartPreview />
        </DashboardSection>
        <DashboardSection
          className="h-[420px]"
          description="Общая статистика по всем дисциплинам"
          title="Распределение оценок"
        >
          <PieChartPreview />
        </DashboardSection>
      </div>
      <DashboardSection
        className="h-[420px]"
        description="Средний балл по каждой дисциплине"
        title="Успеваемость по дисциплинам"
      >
        <BarChartPreview />
      </DashboardSection>
    </div>
  );
}

const administratorOverviewMetrics: Metric[] = [
  {
    title: "Всего студентов",
    value: "250",
    description: "+10% от прошлого месяца",
    icon: Users,
    tone: "positive",
  },
  {
    title: "Средний балл",
    value: "4.5",
    description: "+0.5 от прошлого месяца",
    icon: BarChart3,
    tone: "positive",
  },
  {
    title: "Общая посещаемость",
    value: "80%",
    description: "-10% от прошлого месяца",
    icon: BarChart3,
    tone: "negative",
  },
];

const administratorAnalyticsMetrics: Metric[] = [
  {
    title: "Всего оценок",
    value: "250",
    description: "+10% от прошлого месяца",
    icon: BookOpen,
    tone: "positive",
  },
  {
    title: "Средний балл",
    value: "4.5",
    description: "+0.5 от прошлого месяца",
    icon: BarChart3,
    tone: "positive",
  },
  {
    title: "Общая посещаемость",
    value: "80%",
    description: "-10% от прошлого месяца",
    icon: Users,
    tone: "negative",
  },
];

function PageHeaderWithAction({
  title,
  description,
  actionLabel,
  dialog,
}: {
  title: string;
  description: string;
  actionLabel: string;
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
          <Button className="h-9 rounded-lg px-3">
            <Plus className="size-4" />
            {actionLabel}
          </Button>
        }
      />
    </div>
  );
}

function SearchAndFilter({ placeholder }: { placeholder: string }) {
  return (
    <div className="flex h-9 gap-3">
      <div className="relative min-w-0 flex-1">
        <Search className="-translate-y-1/2 absolute left-3 top-1/2 size-4 text-muted-foreground" />
        <Input
          aria-label="Поиск"
          className="h-9 border-transparent bg-muted pl-9 text-sm shadow-none placeholder:text-muted-foreground focus-visible:ring-0"
          placeholder={placeholder}
        />
      </div>
      <Button className="h-9 rounded-lg px-[17px]" variant="outline">
        Фильтры
      </Button>
    </div>
  );
}

function StudentsTable() {
  return (
    <div className="overflow-hidden rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow className="h-10 hover:bg-transparent">
            <TableHead className="w-[37%]">Студент</TableHead>
            <TableHead className="w-[34%]">Контакты</TableHead>
            <TableHead className="w-[14%]">Средний балл</TableHead>
            <TableHead className="w-[12%]">Статус</TableHead>
            <TableHead className="w-9" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {students.map((student) => (
            <TableRow className="h-[61px]" key={student.id}>
              <TableCell>
                <div className="flex items-center gap-3">
                  <Avatar className="size-9" size="default">
                    <AvatarFallback className="bg-primary text-xs text-primary-foreground">
                      {student.initials}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="text-sm font-medium leading-[14px] text-foreground">
                      {student.name}
                    </div>
                    <div className="mt-1 text-sm leading-5 text-muted-foreground">
                      {student.group}
                    </div>
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <div className="flex flex-col gap-1">
                  <div className="flex items-center gap-2 text-sm leading-5 text-foreground">
                    <Mail className="size-3.5 text-muted-foreground" />
                    {student.email}
                  </div>
                  <div className="flex items-center gap-2 text-sm leading-5 text-muted-foreground">
                    <Phone className="size-3.5" />
                    {student.phone}
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <div className="flex size-9 items-center justify-center rounded-lg bg-green-100 text-sm font-bold text-green-700">
                  {student.average}
                </div>
              </TableCell>
              <TableCell>
                <Badge className="h-[22px] rounded-lg bg-primary px-2 text-primary-foreground">
                  {student.status}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <AdminActionsMenu label={`Действия для ${student.name}`} />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

function DisciplineFilters() {
  const filters = ["Все", "Математика", "Программирование", "Информатика", "Физика"];

  return (
    <div className="flex h-8 items-center gap-2">
      {filters.map((filter, index) => (
        <Button
          className="h-8 rounded-lg px-3"
          key={filter}
          variant={index === 0 ? "default" : "outline"}
        >
          {filter}
        </Button>
      ))}
    </div>
  );
}

function CourseCard({ course }: { course: CourseRecord }) {
  return (
    <Card className="h-[241px] rounded-[14px] py-0 ring-border">
      <CardHeader className="px-6 pb-0 pt-6">
        <CardTitle className="text-lg font-medium leading-7 text-foreground">
          {course.title}
        </CardTitle>
        <CardDescription className="text-base leading-6 text-muted-foreground">
          {course.teacher}
        </CardDescription>
        <CardAction>
          <AdminActionsMenu label={`Действия для курса ${course.title}`} />
        </CardAction>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col justify-between px-6 pb-4">
        <div className="flex items-center gap-2">
          <Badge className="h-[22px] rounded-lg px-[9px]" variant="outline">
            {course.category}
          </Badge>
          <Badge className="h-[22px] rounded-lg px-[9px]" variant="secondary">
            {course.term}
          </Badge>
        </div>
        <div className="flex h-5 items-center justify-between text-sm text-muted-foreground">
          <div className="flex items-center gap-2">
            <Users className="size-4" />
            <span>{course.students} студентов</span>
          </div>
          <div className="flex items-center gap-2">
            <Clock3 className="size-4" />
            <span>{course.hours} часов</span>
          </div>
        </div>
        <Button className="h-9 rounded-lg" variant="outline">
          <BookOpen className="size-4" />
          Открыть курс
        </Button>
      </CardContent>
    </Card>
  );
}

function AdminActionsMenu({ label }: { label: string }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild={true}>
        <Button aria-label={label} className="size-9" size="icon" variant="ghost">
          <MoreHorizontal className="size-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-40">
        <DropdownMenuItem>Открыть</DropdownMenuItem>
        <DropdownMenuItem>Редактировать</DropdownMenuItem>
        <DropdownMenuItem>Дублировать</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem variant="destructive">Удалить</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function LineChartPreview() {
  return (
    <ChartContainer className="h-[300px] w-full" config={trendChartConfig}>
      <LineChart
        accessibilityLayer={true}
        data={performanceTrendData}
        margin={{ left: 8, right: 8 }}
      >
        <CartesianGrid strokeDasharray="3 3" vertical={true} />
        <XAxis dataKey="month" tickLine={false} />
        <YAxis domain={[0, 100]} tickLine={false} width={32} />
        <ChartTooltip content={<ChartTooltipContent />} />
        <Line
          dataKey="attendance"
          dot={true}
          stroke="var(--color-attendance)"
          strokeWidth={2}
          type="monotone"
        />
        <Line
          dataKey="average"
          dot={true}
          stroke="var(--color-average)"
          strokeWidth={2}
          type="monotone"
        />
      </LineChart>
    </ChartContainer>
  );
}

function PieChartPreview() {
  return (
    <ChartContainer className="mx-auto h-[300px] w-full" config={gradeChartConfig}>
      <PieChart accessibilityLayer={true}>
        <ChartTooltip content={<ChartTooltipContent hideLabel={true} />} />
        <Pie data={gradeDistributionData} dataKey="value" nameKey="grade">
          {gradeDistributionData.map((entry) => (
            <Cell fill={entry.fill} key={entry.grade} />
          ))}
        </Pie>
        <ChartLegend content={<ChartLegendContent nameKey="grade" />} />
      </PieChart>
    </ChartContainer>
  );
}

function BarChartPreview() {
  return (
    <ChartContainer className="h-[300px] w-full" config={subjectChartConfig}>
      <BarChart
        accessibilityLayer={true}
        data={subjectPerformanceData}
        margin={{ left: 8, right: 8 }}
      >
        <CartesianGrid vertical={false} />
        <XAxis dataKey="subject" tickLine={false} />
        <YAxis domain={[0, 5]} tickCount={3} tickLine={false} width={32} />
        <ChartTooltip content={<ChartTooltipContent />} />
        <Bar dataKey="average" fill="var(--color-average)" radius={[4, 4, 0, 0]} />
      </BarChart>
    </ChartContainer>
  );
}
