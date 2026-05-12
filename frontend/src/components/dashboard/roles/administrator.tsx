import {
  BarChart3,
  BookOpen,
  Clock3,
  Mail,
  MoreHorizontal,
  Plus,
  Search,
  Users,
} from "lucide-react";
import type { FormEvent, ReactNode } from "react";
import { useCallback, useEffect, useId, useMemo, useState } from "react";
import { Bar, BarChart, CartesianGrid, Cell, Pie, PieChart, XAxis, YAxis } from "recharts";
import { toast } from "sonner";
import type { components } from "@/api";
import { apiClient } from "@/api";
import type { DashboardDialogConfig } from "@/components/dashboard/dashboard-dialogs";
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
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export type AdministratorPage =
  | "home"
  | "students"
  | "teachers"
  | "groups"
  | "specialties"
  | "disciplines"
  | "curriculums"
  | "analytics";

type ApiStudent = components["schemas"]["dto.Student"];
type ApiUser = components["schemas"]["dto.User"];
type ApiGroup = components["schemas"]["dto.Group"];
type ApiGrade = components["schemas"]["dto.Grade"];
type ApiProfile = components["schemas"]["dto.Profile"];
type ApiSpecialty = components["schemas"]["dto.Specialty"];
type ApiSubject = components["schemas"]["dto.Subject"];
type ApiCurriculum = components["schemas"]["dto.Curriculum"];
type ApiAnalyticsOverview = components["schemas"]["dto.AnalyticsOverview"];

type StudentRecord = {
  id: string;
  userId: string;
  groupId: string;
  lastName: string;
  firstName: string;
  middleName: string | undefined;
  initials: string;
  name: string;
  group: string;
  email: string;
  average: string;
  status: string;
};

type DisciplineRecord = {
  id: string;
  title: string;
  groups: string;
  reportTypes: string;
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

const addSubjectDialog: DashboardDialogConfig = {
  title: "Добавить дисциплину",
  description: "Создайте новую дисциплину в справочнике предметов.",
  submitLabel: "Добавить",
  fields: [{ label: "Название", name: "subject-title", placeholder: "Математический анализ" }],
};

export function AdministratorDashboard({ page }: { page: AdministratorPage }) {
  if (page === "students") {
    return <AdministratorStudentsPage />;
  }

  if (page === "teachers") {
    return <AdministratorTeachersPage />;
  }

  if (page === "groups") {
    return <AdministratorGroupsPage />;
  }

  if (page === "specialties") {
    return <AdministratorSpecialtiesPage />;
  }

  if (page === "disciplines") {
    return <AdministratorDisciplinesPage />;
  }

  if (page === "curriculums") {
    return <AdministratorCurriculumsPage />;
  }

  if (page === "analytics") {
    return <AdministratorAnalyticsPage />;
  }

  return <AdministratorHomePage />;
}

function AdministratorHomePage() {
  const [metrics, setMetrics] = useState<Metric[]>(administratorOverviewMetrics);

  useEffect(() => {
    let mounted = true;

    async function loadMetrics() {
      const { data, error } = await apiClient.GET("/analytics/overview");

      if (!mounted || error || !data) {
        return;
      }

      setMetrics([
        {
          title: "Всего студентов",
          value: String(data.students_count ?? 0),
          description: "В базе студентов",
          icon: Users,
          tone: "positive",
        },
        {
          title: "Средний балл",
          value: formatAverage(data.average_grade),
          description: "По всем выставленным оценкам",
          icon: BarChart3,
          tone: "positive",
        },
        {
          title: "Дисциплины",
          value: String(data.subjects_count ?? 0),
          description: "В справочнике предметов",
          icon: BookOpen,
          tone: "positive",
        },
      ]);
    }

    void loadMetrics();

    return () => {
      mounted = false;
    };
  }, []);

  return (
    <div className="flex flex-col gap-6">
      <PageHeading
        description="Обзор успеваемости и активности студентов"
        title="Панель управления"
      />
      <MetricsGrid metrics={metrics} />
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
  const [records, setRecords] = useState<StudentRecord[]>([]);
  const [groups, setGroups] = useState<ApiGroup[]>([]);
  const [editingStudent, setEditingStudent] = useState<StudentRecord | null>(null);
  const [search, setSearch] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const filteredRecords = useMemo(() => {
    const query = search.trim().toLowerCase();
    if (!query) {
      return records;
    }

    return records.filter((record) =>
      [record.name, record.email, record.group].some((value) =>
        value.toLowerCase().includes(query),
      ),
    );
  }, [records, search]);

  const loadStudents = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await loadStudentRecords();
      setRecords(data.records);
      setGroups(data.groups);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Не удалось загрузить студентов.");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadStudents();
  }, [loadStudents]);

  async function handleDeleteStudent(student: StudentRecord) {
    const { error: studentError } = await apiClient.DELETE("/students/{id}", {
      params: { path: { id: student.id } },
    });
    if (studentError) {
      toast.error("Не удалось удалить студента.");
      return;
    }

    const { error: userError } = await apiClient.DELETE("/users/{id}", {
      params: { path: { id: student.userId } },
    });
    if (userError) {
      toast.error("Студент удалён, но аккаунт пользователя остался.");
    } else {
      toast.success("Студент удалён.");
    }
    await loadStudents();
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        actionLabel="Добавить студента"
        description="Управление списком студентов и их данными"
        action={<AddStudentDialog groups={groups} onCreated={loadStudents} />}
        title="Студенты"
      />
      <SearchAndFilter
        onChange={setSearch}
        placeholder="Поиск студентов по имени, email или группе..."
        value={search}
      />
      <StudentsTable
        isLoading={isLoading}
        onEdit={setEditingStudent}
        onDelete={handleDeleteStudent}
        students={filteredRecords}
      />
      <EditStudentDialog
        groups={groups}
        onSaved={loadStudents}
        onStudentChange={setEditingStudent}
        student={editingStudent}
      />
    </div>
  );
}

function AdministratorTeachersPage() {
  const [teachers, setTeachers] = useState<ApiProfile[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadTeachers = useCallback(async () => {
    setIsLoading(true);
    const { data, error } = await apiClient.GET("/teachers");
    if (error) {
      toast.error("Не удалось загрузить преподавателей.");
    } else {
      setTeachers(data ?? []);
    }
    setIsLoading(false);
  }, []);

  useEffect(() => {
    void loadTeachers();
  }, [loadTeachers]);

  async function handleDelete(profile: ApiProfile) {
    if (!profile.id) {
      return;
    }

    const { error } = await apiClient.DELETE("/teachers/{id}", {
      params: { path: { id: profile.id } },
    });
    if (error) {
      toast.error("Не удалось удалить профиль преподавателя.");
      return;
    }

    if (profile.user_id) {
      await apiClient.DELETE("/users/{id}", {
        params: { path: { id: profile.user_id } },
      });
    }
    toast.success("Преподаватель удалён.");
    await loadTeachers();
  }

  return (
    <StaffPage
      isLoading={isLoading}
      onCreated={loadTeachers}
      onDelete={handleDelete}
      profiles={teachers}
      staffRole="teacher"
      title="Преподаватели"
    />
  );
}

function AdministratorSpecialtiesPage() {
  const [specialties, setSpecialties] = useState<ApiSpecialty[]>([]);
  const [editingSpecialty, setEditingSpecialty] = useState<ApiSpecialty | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const loadSpecialties = useCallback(async () => {
    setIsLoading(true);
    const { data, error } = await apiClient.GET("/specialties");
    if (error) {
      toast.error("Не удалось загрузить специальности.");
    } else {
      setSpecialties(data ?? []);
    }
    setIsLoading(false);
  }, []);

  useEffect(() => {
    void loadSpecialties();
  }, [loadSpecialties]);

  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        actionLabel="Добавить специальность"
        description="Управление направлениями подготовки"
        dialog={{
          description: "Добавьте новую специальность.",
          fields: [{ label: "Название", name: "title", placeholder: "Информационные системы" }],
          submitLabel: "Добавить",
          title: "Добавить специальность",
        }}
        onDialogSubmit={async (formData) => {
          const title = String(formData.get("title") ?? "").trim();
          if (!title) {
            toast.error("Введите название специальности.");
            return;
          }
          const { error } = await apiClient.POST("/specialties", { body: { title } });
          if (error) {
            toast.error("Не удалось добавить специальность.");
            return;
          }
          toast.success("Специальность добавлена.");
          await loadSpecialties();
        }}
        title="Специальности"
      />
      <SimpleEntityTable
        isLoading={isLoading}
        items={specialties}
        titleHeader="Специальность"
        getTitle={(item) => item.title ?? "Без названия"}
        onEdit={setEditingSpecialty}
        onDelete={async (item) => {
          if (!item.id) {
            return;
          }
          const { error } = await apiClient.DELETE("/specialties/{id}", {
            params: { path: { id: item.id } },
          });
          if (error) {
            toast.error("Не удалось удалить специальность. Возможно, есть связанные группы.");
            return;
          }
          toast.success("Специальность удалена.");
          await loadSpecialties();
        }}
      />
      <EditSpecialtyDialog
        onSaved={loadSpecialties}
        onSpecialtyChange={setEditingSpecialty}
        specialty={editingSpecialty}
      />
    </div>
  );
}

function AdministratorGroupsPage() {
  const [groups, setGroups] = useState<ApiGroup[]>([]);
  const [specialties, setSpecialties] = useState<ApiSpecialty[]>([]);
  const [editingGroup, setEditingGroup] = useState<ApiGroup | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const loadGroups = useCallback(async () => {
    setIsLoading(true);
    const [groupsResponse, specialtiesResponse] = await Promise.all([
      apiClient.GET("/groups"),
      apiClient.GET("/specialties"),
    ]);
    if (groupsResponse.error || specialtiesResponse.error) {
      toast.error("Не удалось загрузить группы.");
    } else {
      setGroups(groupsResponse.data ?? []);
      setSpecialties(specialtiesResponse.data ?? []);
    }
    setIsLoading(false);
  }, []);

  useEffect(() => {
    void loadGroups();
  }, [loadGroups]);

  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        action={<GroupDialog onSaved={loadGroups} specialties={specialties} />}
        actionLabel="Добавить группу"
        description="Управление учебными группами"
        title="Группы"
      />
      <GroupsTable
        groups={groups}
        isLoading={isLoading}
        onEdit={setEditingGroup}
        onDelete={async (group) => {
          if (!group.id) {
            return;
          }
          const { error } = await apiClient.DELETE("/groups/{id}", {
            params: { path: { id: group.id } },
          });
          if (error) {
            toast.error("Не удалось удалить группу. Возможно, есть связанные студенты.");
            return;
          }
          toast.success("Группа удалена.");
          await loadGroups();
        }}
        specialties={specialties}
      />
      <EditGroupDialog
        group={editingGroup}
        onGroupChange={setEditingGroup}
        onSaved={loadGroups}
        specialties={specialties}
      />
    </div>
  );
}

function AdministratorDisciplinesPage() {
  const [records, setRecords] = useState<DisciplineRecord[]>([]);
  const [editingDiscipline, setEditingDiscipline] = useState<DisciplineRecord | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const loadDisciplines = useCallback(async () => {
    setIsLoading(true);
    try {
      setRecords(await loadDisciplineRecords());
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Не удалось загрузить дисциплины.");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadDisciplines();
  }, [loadDisciplines]);

  async function handleDeleteDiscipline(id: string) {
    const { error } = await apiClient.DELETE("/subjects/{id}", {
      params: { path: { id } },
    });
    if (error) {
      toast.error("Не удалось удалить дисциплину. Возможно, она используется в учебном плане.");
      return;
    }

    toast.success("Дисциплина удалена.");
    await loadDisciplines();
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        actionLabel="Добавить дисциплину"
        dialog={addSubjectDialog}
        description="Управление учебными курсами и программами"
        onDialogSubmit={async (formData) => {
          const title = String(formData.get("subject-title") ?? "").trim();
          if (!title) {
            toast.error("Введите название дисциплины.");
            return;
          }

          const { error } = await apiClient.POST("/subjects", {
            body: { title },
          });
          if (error) {
            toast.error("Не удалось добавить дисциплину.");
            return;
          }

          toast.success("Дисциплина добавлена.");
          await loadDisciplines();
        }}
        title="Дисциплины"
      />
      <DisciplineFilters />
      <div className="grid grid-cols-3 gap-2.5">
        {isLoading && <EmptyState text="Загружаем дисциплины..." />}
        {!isLoading && records.length === 0 && <EmptyState text="Дисциплины пока не добавлены." />}
        {records.map((discipline) => (
          <CourseCard
            course={discipline}
            key={discipline.id}
            onEdit={() => setEditingDiscipline(discipline)}
            onDelete={() => handleDeleteDiscipline(discipline.id)}
          />
        ))}
      </div>
      <EditDisciplineDialog
        discipline={editingDiscipline}
        onDisciplineChange={setEditingDiscipline}
        onSaved={loadDisciplines}
      />
    </div>
  );
}

function AdministratorCurriculumsPage() {
  const [curriculums, setCurriculums] = useState<ApiCurriculum[]>([]);
  const [subjects, setSubjects] = useState<ApiSubject[]>([]);
  const [groups, setGroups] = useState<ApiGroup[]>([]);
  const [teachers, setTeachers] = useState<ApiProfile[]>([]);
  const [editingCurriculum, setEditingCurriculum] = useState<ApiCurriculum | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const loadCurriculums = useCallback(async () => {
    setIsLoading(true);
    const [curriculumsResponse, subjectsResponse, groupsResponse, teachersResponse] =
      await Promise.all([
        apiClient.GET("/curriculums"),
        apiClient.GET("/subjects"),
        apiClient.GET("/groups"),
        apiClient.GET("/teachers"),
      ]);

    if (
      curriculumsResponse.error ||
      subjectsResponse.error ||
      groupsResponse.error ||
      teachersResponse.error
    ) {
      toast.error("Не удалось загрузить учебные планы.");
    } else {
      setCurriculums(curriculumsResponse.data ?? []);
      setSubjects(subjectsResponse.data ?? []);
      setGroups(groupsResponse.data ?? []);
      setTeachers(teachersResponse.data ?? []);
    }
    setIsLoading(false);
  }, []);

  useEffect(() => {
    void loadCurriculums();
  }, [loadCurriculums]);

  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        action={
          <CurriculumDialog
            groups={groups}
            onSaved={loadCurriculums}
            subjects={subjects}
            teachers={teachers}
          />
        }
        actionLabel="Добавить учебный план"
        description="Связь дисциплин, групп, семестров и преподавателей"
        title="Учебные планы"
      />
      <CurriculumsTable
        curriculums={curriculums}
        groups={groups}
        isLoading={isLoading}
        onEdit={setEditingCurriculum}
        onDelete={async (curriculum) => {
          if (!curriculum.id) {
            return;
          }
          const { error } = await apiClient.DELETE("/curriculums/{id}", {
            params: { path: { id: curriculum.id } },
          });
          if (error) {
            toast.error("Не удалось удалить учебный план.");
            return;
          }
          toast.success("Учебный план удалён.");
          await loadCurriculums();
        }}
        subjects={subjects}
        teachers={teachers}
      />
      <EditCurriculumDialog
        curriculum={editingCurriculum}
        groups={groups}
        onCurriculumChange={setEditingCurriculum}
        onSaved={loadCurriculums}
        subjects={subjects}
        teachers={teachers}
      />
    </div>
  );
}

function AdministratorAnalyticsPage() {
  const [overview, setOverview] = useState<ApiAnalyticsOverview | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let mounted = true;

    async function loadOverview() {
      setIsLoading(true);
      const { data, error } = await apiClient.GET("/analytics/overview");
      if (!mounted) {
        return;
      }
      if (error) {
        toast.error("Не удалось загрузить аналитику.");
      } else {
        setOverview(data ?? null);
      }
      setIsLoading(false);
    }

    void loadOverview();

    return () => {
      mounted = false;
    };
  }, []);

  const metrics = useMemo(() => toAnalyticsMetrics(overview), [overview]);
  const distribution = useMemo(() => toGradeDistributionData(overview), [overview]);
  const subjectAverages = useMemo(() => toSubjectPerformanceData(overview), [overview]);

  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="Статистика и анализ успеваемости студентов" title="Аналитика" />
      <MetricsGrid metrics={metrics} />
      <div className="grid grid-cols-2 gap-4">
        <DashboardSection
          className="h-[420px]"
          description="Сводка по студентам, преподавателям и учебным планам"
          title="Наполнение системы"
        >
          <AnalyticsSummary overview={overview} isLoading={isLoading} />
        </DashboardSection>
        <DashboardSection
          className="h-[420px]"
          description="Общая статистика по всем дисциплинам"
          title="Распределение оценок"
        >
          <PieChartPreview data={distribution} />
        </DashboardSection>
      </div>
      <DashboardSection
        className="h-[420px]"
        description="Средний балл по каждой дисциплине"
        title="Успеваемость по дисциплинам"
      >
        <BarChartPreview data={subjectAverages} />
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
    value: "0",
    description: "В журнале оценок",
    icon: BookOpen,
    tone: "positive",
  },
  {
    title: "Средний балл",
    value: "-",
    description: "По всем оценкам",
    icon: BarChart3,
    tone: "positive",
  },
  {
    title: "Учебные планы",
    value: "0",
    description: "Связки групп и дисциплин",
    icon: Users,
    tone: "positive",
  },
];

function PageHeaderWithAction({
  title,
  description,
  actionLabel,
  action,
  dialog,
  onDialogSubmit,
}: {
  title: string;
  description: string;
  actionLabel: string;
  action?: ReactNode;
  dialog?: DashboardDialogConfig;
  onDialogSubmit?: (formData: FormData) => Promise<void> | void;
}) {
  return (
    <div className="flex h-16 items-center justify-between">
      <div className="flex flex-col gap-1">
        <h1 className="h-9 text-[30px] font-semibold leading-9 text-foreground">{title}</h1>
        <p className="h-6 text-base leading-6 text-muted-foreground">{description}</p>
      </div>
      {action ??
        (dialog ? (
          <SimpleActionDialog
            actionLabel={actionLabel}
            config={dialog}
            {...(onDialogSubmit ? { onSubmit: onDialogSubmit } : {})}
          />
        ) : null)}
    </div>
  );
}

function SearchAndFilter({
  placeholder,
  value,
  onChange,
}: {
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
}) {
  return (
    <div className="flex h-9 gap-3">
      <div className="relative min-w-0 flex-1">
        <Search className="-translate-y-1/2 absolute left-3 top-1/2 size-4 text-muted-foreground" />
        <Input
          aria-label="Поиск"
          className="h-9 border-transparent bg-muted pl-9 text-sm shadow-none placeholder:text-muted-foreground focus-visible:ring-0"
          onChange={(event) => onChange(event.target.value)}
          placeholder={placeholder}
          value={value}
        />
      </div>
      <Button className="h-9 rounded-lg px-[17px]" variant="outline">
        Фильтры
      </Button>
    </div>
  );
}

function StudentsTable({
  students,
  isLoading,
  onEdit,
  onDelete,
}: {
  students: StudentRecord[];
  isLoading: boolean;
  onEdit: (student: StudentRecord) => void;
  onDelete: (student: StudentRecord) => void;
}) {
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
          {isLoading && (
            <TableRow>
              <TableCell className="h-24 text-center text-muted-foreground" colSpan={5}>
                Загружаем студентов...
              </TableCell>
            </TableRow>
          )}
          {!isLoading && students.length === 0 && (
            <TableRow>
              <TableCell className="h-24 text-center text-muted-foreground" colSpan={5}>
                Студенты пока не добавлены.
              </TableCell>
            </TableRow>
          )}
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
                <AdminActionsMenu
                  label={`Действия для ${student.name}`}
                  onEdit={() => onEdit(student)}
                  onDelete={() => onDelete(student)}
                />
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

function CourseCard({
  course,
  onEdit,
  onDelete,
}: {
  course: DisciplineRecord;
  onEdit: () => void;
  onDelete: () => void;
}) {
  return (
    <Card className="h-[241px] rounded-[14px] py-0 ring-border">
      <CardHeader className="px-6 pb-0 pt-6">
        <CardTitle className="text-lg font-medium leading-7 text-foreground">
          {course.title}
        </CardTitle>
        <CardDescription className="text-base leading-6 text-muted-foreground">
          {course.groups}
        </CardDescription>
        <CardAction>
          <AdminActionsMenu
            label={`Действия для дисциплины ${course.title}`}
            onEdit={onEdit}
            onDelete={onDelete}
          />
        </CardAction>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col justify-between px-6 pb-4">
        <div className="flex items-center gap-2">
          <Badge className="h-[22px] rounded-lg px-[9px]" variant="outline">
            {course.reportTypes}
          </Badge>
          <Badge className="h-[22px] rounded-lg px-[9px]" variant="secondary">
            Учебные планы
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

function AnalyticsSummary({
  overview,
  isLoading,
}: {
  overview: ApiAnalyticsOverview | null;
  isLoading: boolean;
}) {
  if (isLoading) {
    return <EmptyState text="Загружаем аналитику..." />;
  }

  const items = [
    { label: "Студенты", value: overview?.students_count ?? 0 },
    { label: "Преподаватели", value: overview?.teachers_count ?? 0 },
    { label: "Группы", value: overview?.groups_count ?? 0 },
    { label: "Специальности", value: overview?.specialties_count ?? 0 },
    { label: "Дисциплины", value: overview?.subjects_count ?? 0 },
    { label: "Учебные планы", value: overview?.curriculums_count ?? 0 },
  ];

  return (
    <div className="grid h-full grid-cols-2 content-start gap-3">
      {items.map((item) => (
        <div className="rounded-lg border border-border p-4" key={item.label}>
          <div className="text-sm text-muted-foreground">{item.label}</div>
          <div className="mt-2 text-2xl font-semibold text-foreground">{item.value}</div>
        </div>
      ))}
    </div>
  );
}

function AdminActionsMenu({
  label,
  onEdit,
  onDelete,
}: {
  label: string;
  onEdit?: (() => void) | undefined;
  onDelete: () => void;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild={true}>
        <Button aria-label={label} className="size-9" size="icon" variant="ghost">
          <MoreHorizontal className="size-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-40">
        <DropdownMenuItem disabled={true}>Открыть</DropdownMenuItem>
        <DropdownMenuItem disabled={!onEdit} onClick={onEdit}>
          Редактировать
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={onDelete} variant="destructive">
          Удалить
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function SimpleActionDialog({
  actionLabel,
  config,
  onSubmit,
}: {
  actionLabel: string;
  config: DashboardDialogConfig;
  onSubmit?: (formData: FormData) => Promise<void> | void;
}) {
  const [open, setOpen] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    await onSubmit?.(new FormData(event.currentTarget));
    setOpen(false);
  }

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild={true}>
        <Button className="h-9 rounded-lg px-3">
          <Plus className="size-4" />
          {actionLabel}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>{config.title}</DialogTitle>
          <DialogDescription>{config.description}</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <FieldGroup className="gap-4">
            {config.fields.map((field) => (
              <Field key={field.name}>
                <FieldLabel htmlFor={field.name}>{field.label}</FieldLabel>
                <Input
                  id={field.name}
                  max={field.max}
                  min={field.min}
                  name={field.name}
                  placeholder={field.placeholder}
                  type={field.type ?? "text"}
                />
              </Field>
            ))}
          </FieldGroup>
          <DialogFooter>
            <DialogClose asChild={true}>
              <Button type="button" variant="outline">
                Отмена
              </Button>
            </DialogClose>
            <Button type="submit">{config.submitLabel}</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function AddStudentDialog({
  groups,
  onCreated,
}: {
  groups: ApiGroup[];
  onCreated: () => Promise<void>;
}) {
  const [open, setOpen] = useState(false);
  const [groupID, setGroupID] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const lastNameID = useId();
  const firstNameID = useId();
  const middleNameID = useId();
  const emailID = useId();
  const passwordID = useId();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const lastName = String(formData.get("last_name") ?? "").trim();
    const firstName = String(formData.get("first_name") ?? "").trim();
    const middleName = String(formData.get("middle_name") ?? "").trim();
    const email = String(formData.get("email") ?? "").trim();
    const password = String(formData.get("password") ?? "");

    if (!lastName || !firstName || !email || !password || !groupID) {
      toast.error("Заполните обязательные поля студента.");
      return;
    }

    setIsSubmitting(true);
    try {
      const { data: user, error: userError } = await apiClient.POST("/users", {
        body: { email, password, role: "student" },
      });
      if (userError || !user?.id) {
        throw new Error("Не удалось создать аккаунт студента.");
      }

      const studentBody = {
        first_name: firstName,
        group_id: groupID,
        last_name: lastName,
        user_id: user.id,
        ...(middleName ? { middle_name: middleName } : {}),
      };
      const { error: studentError } = await apiClient.POST("/students", {
        body: studentBody,
      });
      if (studentError) {
        await apiClient.DELETE("/users/{id}", {
          params: { path: { id: user.id } },
        });
        throw new Error("Не удалось создать профиль студента.");
      }

      toast.success("Студент добавлен.");
      setOpen(false);
      setGroupID("");
      await onCreated();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Не удалось добавить студента.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild={true}>
        <Button className="h-9 rounded-lg px-3" disabled={groups.length === 0}>
          <Plus className="size-4" />
          Добавить студента
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[520px]">
        <DialogHeader>
          <DialogTitle>Добавить студента</DialogTitle>
          <DialogDescription>Создайте аккаунт студента и привяжите его к группе.</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <FieldGroup className="grid grid-cols-2 gap-4">
            <Field>
              <FieldLabel htmlFor={lastNameID}>Фамилия</FieldLabel>
              <Input id={lastNameID} name="last_name" placeholder="Иванов" />
            </Field>
            <Field>
              <FieldLabel htmlFor={firstNameID}>Имя</FieldLabel>
              <Input id={firstNameID} name="first_name" placeholder="Иван" />
            </Field>
            <Field>
              <FieldLabel htmlFor={middleNameID}>Отчество</FieldLabel>
              <Input id={middleNameID} name="middle_name" placeholder="Иванович" />
            </Field>
            <Field>
              <FieldLabel>Группа</FieldLabel>
              <Select onValueChange={setGroupID} value={groupID}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Выберите группу" />
                </SelectTrigger>
                <SelectContent>
                  {groups.map((group) => (
                    <SelectItem key={group.id} value={group.id ?? ""}>
                      {group.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel htmlFor={emailID}>E-mail</FieldLabel>
              <Input id={emailID} name="email" placeholder="student@example.com" type="email" />
            </Field>
            <Field>
              <FieldLabel htmlFor={passwordID}>Пароль</FieldLabel>
              <Input id={passwordID} name="password" minLength={8} type="password" />
            </Field>
          </FieldGroup>
          <DialogFooter>
            <DialogClose asChild={true}>
              <Button type="button" variant="outline">
                Отмена
              </Button>
            </DialogClose>
            <Button disabled={isSubmitting} type="submit">
              {isSubmitting ? "Добавляем..." : "Добавить"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function EditStudentDialog({
  student,
  groups,
  onStudentChange,
  onSaved,
}: {
  student: StudentRecord | null;
  groups: ApiGroup[];
  onStudentChange: (student: StudentRecord | null) => void;
  onSaved: () => Promise<void>;
}) {
  const [groupID, setGroupID] = useState(student?.groupId ?? "");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const lastNameID = useId();
  const firstNameID = useId();
  const middleNameID = useId();

  useEffect(() => {
    setGroupID(student?.groupId ?? "");
  }, [student]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!student) {
      return;
    }

    const formData = new FormData(event.currentTarget);
    const lastName = String(formData.get("last_name") ?? "").trim();
    const firstName = String(formData.get("first_name") ?? "").trim();
    const middleName = String(formData.get("middle_name") ?? "").trim();

    if (!lastName || !firstName || !groupID) {
      toast.error("Заполните обязательные поля студента.");
      return;
    }

    setIsSubmitting(true);
    const { error } = await apiClient.PATCH("/students/{id}", {
      body: {
        first_name: firstName,
        group_id: groupID,
        last_name: lastName,
        middle_name: middleName,
      },
      params: { path: { id: student.id } },
    });
    setIsSubmitting(false);

    if (error) {
      toast.error("Не удалось обновить студента.");
      return;
    }

    toast.success("Студент обновлён.");
    onStudentChange(null);
    await onSaved();
  }

  return (
    <Dialog onOpenChange={(open) => !open && onStudentChange(null)} open={Boolean(student)}>
      <DialogContent className="sm:max-w-[520px]">
        <DialogHeader>
          <DialogTitle>Редактировать студента</DialogTitle>
          <DialogDescription>Обновите ФИО студента и привязку к группе.</DialogDescription>
        </DialogHeader>
        {student ? (
          <form className="flex flex-col gap-5" key={student.id} onSubmit={handleSubmit}>
            <FieldGroup className="grid grid-cols-2 gap-4">
              <Field>
                <FieldLabel htmlFor={lastNameID}>Фамилия</FieldLabel>
                <Input id={lastNameID} name="last_name" defaultValue={student.lastName} />
              </Field>
              <Field>
                <FieldLabel htmlFor={firstNameID}>Имя</FieldLabel>
                <Input id={firstNameID} name="first_name" defaultValue={student.firstName} />
              </Field>
              <Field>
                <FieldLabel htmlFor={middleNameID}>Отчество</FieldLabel>
                <Input
                  id={middleNameID}
                  name="middle_name"
                  defaultValue={student.middleName ?? ""}
                />
              </Field>
              <Field>
                <FieldLabel>Группа</FieldLabel>
                <Select onValueChange={setGroupID} value={groupID}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Выберите группу" />
                  </SelectTrigger>
                  <SelectContent>
                    {groups.map((group) => (
                      <SelectItem key={group.id} value={group.id ?? ""}>
                        {group.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </Field>
            </FieldGroup>
            <DialogFooter>
              <DialogClose asChild={true}>
                <Button type="button" variant="outline">
                  Отмена
                </Button>
              </DialogClose>
              <Button disabled={isSubmitting} type="submit">
                {isSubmitting ? "Сохраняем..." : "Сохранить"}
              </Button>
            </DialogFooter>
          </form>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function EditDisciplineDialog({
  discipline,
  onDisciplineChange,
  onSaved,
}: {
  discipline: DisciplineRecord | null;
  onDisciplineChange: (discipline: DisciplineRecord | null) => void;
  onSaved: () => Promise<void>;
}) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const titleID = useId();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!discipline) {
      return;
    }

    const formData = new FormData(event.currentTarget);
    const title = String(formData.get("title") ?? "").trim();
    if (!title) {
      toast.error("Введите название дисциплины.");
      return;
    }

    setIsSubmitting(true);
    const { error } = await apiClient.PATCH("/subjects/{id}", {
      body: { title },
      params: { path: { id: discipline.id } },
    });
    setIsSubmitting(false);

    if (error) {
      toast.error("Не удалось обновить дисциплину.");
      return;
    }

    toast.success("Дисциплина обновлена.");
    onDisciplineChange(null);
    await onSaved();
  }

  return (
    <Dialog onOpenChange={(open) => !open && onDisciplineChange(null)} open={Boolean(discipline)}>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>Редактировать дисциплину</DialogTitle>
          <DialogDescription>Измените название дисциплины.</DialogDescription>
        </DialogHeader>
        {discipline ? (
          <form className="flex flex-col gap-5" key={discipline.id} onSubmit={handleSubmit}>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor={titleID}>Название</FieldLabel>
                <Input id={titleID} name="title" defaultValue={discipline.title} />
              </Field>
            </FieldGroup>
            <DialogFooter>
              <DialogClose asChild={true}>
                <Button type="button" variant="outline">
                  Отмена
                </Button>
              </DialogClose>
              <Button disabled={isSubmitting} type="submit">
                {isSubmitting ? "Сохраняем..." : "Сохранить"}
              </Button>
            </DialogFooter>
          </form>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function EditSpecialtyDialog({
  specialty,
  onSpecialtyChange,
  onSaved,
}: {
  specialty: ApiSpecialty | null;
  onSpecialtyChange: (specialty: ApiSpecialty | null) => void;
  onSaved: () => Promise<void>;
}) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const titleID = useId();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!specialty?.id) {
      return;
    }

    const title = String(new FormData(event.currentTarget).get("title") ?? "").trim();
    if (!title) {
      toast.error("Введите название специальности.");
      return;
    }

    setIsSubmitting(true);
    const { error } = await apiClient.PATCH("/specialties/{id}", {
      body: { title },
      params: { path: { id: specialty.id } },
    });
    setIsSubmitting(false);

    if (error) {
      toast.error("Не удалось обновить специальность.");
      return;
    }

    toast.success("Специальность обновлена.");
    onSpecialtyChange(null);
    await onSaved();
  }

  return (
    <Dialog onOpenChange={(open) => !open && onSpecialtyChange(null)} open={Boolean(specialty)}>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>Редактировать специальность</DialogTitle>
          <DialogDescription>Измените название направления подготовки.</DialogDescription>
        </DialogHeader>
        {specialty ? (
          <form className="flex flex-col gap-5" key={specialty.id} onSubmit={handleSubmit}>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor={titleID}>Название</FieldLabel>
                <Input id={titleID} name="title" defaultValue={specialty.title} />
              </Field>
            </FieldGroup>
            <DialogFooter>
              <DialogClose asChild={true}>
                <Button type="button" variant="outline">
                  Отмена
                </Button>
              </DialogClose>
              <Button disabled={isSubmitting} type="submit">
                {isSubmitting ? "Сохраняем..." : "Сохранить"}
              </Button>
            </DialogFooter>
          </form>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function EditGroupDialog({
  group,
  specialties,
  onGroupChange,
  onSaved,
}: {
  group: ApiGroup | null;
  specialties: ApiSpecialty[];
  onGroupChange: (group: ApiGroup | null) => void;
  onSaved: () => Promise<void>;
}) {
  const [specialtyID, setSpecialtyID] = useState(group?.specialty_id ?? "");
  const [studyForm, setStudyForm] = useState(group?.study_form ?? "full_time");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const nameID = useId();
  const yearID = useId();

  useEffect(() => {
    setSpecialtyID(group?.specialty_id ?? "");
    setStudyForm(group?.study_form ?? "full_time");
  }, [group]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!group?.id) {
      return;
    }

    const formData = new FormData(event.currentTarget);
    const name = String(formData.get("name") ?? "").trim();
    const admissionYear = Number(formData.get("admission_year"));
    if (!name || !specialtyID || !admissionYear) {
      toast.error("Заполните обязательные поля группы.");
      return;
    }

    setIsSubmitting(true);
    const { error } = await apiClient.PATCH("/groups/{id}", {
      body: {
        admission_year: admissionYear,
        name,
        specialty_id: specialtyID,
        study_form: studyForm as components["schemas"]["repositories.StudyForm"],
      },
      params: { path: { id: group.id } },
    });
    setIsSubmitting(false);

    if (error) {
      toast.error("Не удалось обновить группу.");
      return;
    }

    toast.success("Группа обновлена.");
    onGroupChange(null);
    await onSaved();
  }

  return (
    <Dialog onOpenChange={(open) => !open && onGroupChange(null)} open={Boolean(group)}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Редактировать группу</DialogTitle>
          <DialogDescription>Измените параметры учебной группы.</DialogDescription>
        </DialogHeader>
        {group ? (
          <form className="flex flex-col gap-5" key={group.id} onSubmit={handleSubmit}>
            <FieldGroup className="gap-4">
              <Field>
                <FieldLabel htmlFor={nameID}>Название</FieldLabel>
                <Input id={nameID} name="name" defaultValue={group.name} />
              </Field>
              <Field>
                <FieldLabel htmlFor={yearID}>Год поступления</FieldLabel>
                <Input
                  id={yearID}
                  name="admission_year"
                  defaultValue={group.admission_year}
                  type="number"
                />
              </Field>
              <Field>
                <FieldLabel>Форма обучения</FieldLabel>
                <Select
                  onValueChange={(value) =>
                    setStudyForm(value as components["schemas"]["repositories.StudyForm"])
                  }
                  value={studyForm}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="full_time">Очная</SelectItem>
                    <SelectItem value="evening">Очно-заочная</SelectItem>
                    <SelectItem value="extramural">Заочная</SelectItem>
                  </SelectContent>
                </Select>
              </Field>
              <Field>
                <FieldLabel>Специальность</FieldLabel>
                <Select onValueChange={setSpecialtyID} value={specialtyID}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Выберите специальность" />
                  </SelectTrigger>
                  <SelectContent>
                    {specialties.map((specialty) => (
                      <SelectItem key={specialty.id} value={specialty.id ?? ""}>
                        {specialty.title}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </Field>
            </FieldGroup>
            <DialogFooter>
              <DialogClose asChild={true}>
                <Button type="button" variant="outline">
                  Отмена
                </Button>
              </DialogClose>
              <Button disabled={isSubmitting} type="submit">
                {isSubmitting ? "Сохраняем..." : "Сохранить"}
              </Button>
            </DialogFooter>
          </form>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function EditCurriculumDialog({
  curriculum,
  subjects,
  groups,
  teachers,
  onCurriculumChange,
  onSaved,
}: {
  curriculum: ApiCurriculum | null;
  subjects: ApiSubject[];
  groups: ApiGroup[];
  teachers: ApiProfile[];
  onCurriculumChange: (curriculum: ApiCurriculum | null) => void;
  onSaved: () => Promise<void>;
}) {
  const [subjectID, setSubjectID] = useState(curriculum?.subject_id ?? "");
  const [groupID, setGroupID] = useState(curriculum?.group_id ?? "");
  const [teacherUserID, setTeacherUserID] = useState(curriculum?.lead_by ?? "");
  const [reportType, setReportType] = useState(curriculum?.report_type ?? "exam");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const hoursID = useId();
  const semesterID = useId();

  useEffect(() => {
    setSubjectID(curriculum?.subject_id ?? "");
    setGroupID(curriculum?.group_id ?? "");
    setTeacherUserID(curriculum?.lead_by ?? "");
    setReportType(curriculum?.report_type ?? "exam");
  }, [curriculum]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!curriculum?.id) {
      return;
    }

    const formData = new FormData(event.currentTarget);
    const hours = Number(formData.get("hours"));
    const semester = Number(formData.get("semester"));
    if (!subjectID || !groupID || !teacherUserID || !hours || !semester) {
      toast.error("Заполните обязательные поля учебного плана.");
      return;
    }

    setIsSubmitting(true);
    const { error } = await apiClient.PATCH("/curriculums/{id}", {
      body: {
        group_id: groupID,
        hours,
        lead_by: teacherUserID,
        report_type: reportType as components["schemas"]["repositories.CurriculumReportType"],
        semester,
        subject_id: subjectID,
      },
      params: { path: { id: curriculum.id } },
    });
    setIsSubmitting(false);

    if (error) {
      toast.error("Не удалось обновить учебный план.");
      return;
    }

    toast.success("Учебный план обновлён.");
    onCurriculumChange(null);
    await onSaved();
  }

  return (
    <Dialog onOpenChange={(open) => !open && onCurriculumChange(null)} open={Boolean(curriculum)}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Редактировать учебный план</DialogTitle>
          <DialogDescription>
            Измените дисциплину, группу, преподавателя и нагрузку.
          </DialogDescription>
        </DialogHeader>
        {curriculum ? (
          <form className="flex flex-col gap-5" key={curriculum.id} onSubmit={handleSubmit}>
            <FieldGroup className="gap-4">
              <Field>
                <FieldLabel>Дисциплина</FieldLabel>
                <Select onValueChange={setSubjectID} value={subjectID}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Выберите дисциплину" />
                  </SelectTrigger>
                  <SelectContent>
                    {subjects.map((subject) => (
                      <SelectItem key={subject.id} value={subject.id ?? ""}>
                        {subject.title}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </Field>
              <Field>
                <FieldLabel>Группа</FieldLabel>
                <Select onValueChange={setGroupID} value={groupID}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Выберите группу" />
                  </SelectTrigger>
                  <SelectContent>
                    {groups.map((group) => (
                      <SelectItem key={group.id} value={group.id ?? ""}>
                        {group.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </Field>
              <Field>
                <FieldLabel>Преподаватель</FieldLabel>
                <Select onValueChange={setTeacherUserID} value={teacherUserID}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Выберите преподавателя" />
                  </SelectTrigger>
                  <SelectContent>
                    {teachers.map((teacher) => (
                      <SelectItem key={teacher.id} value={teacher.user_id ?? ""}>
                        {getProfileName(teacher)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </Field>
              <Field>
                <FieldLabel>Тип отчётности</FieldLabel>
                <Select
                  onValueChange={(value) =>
                    setReportType(
                      value as components["schemas"]["repositories.CurriculumReportType"],
                    )
                  }
                  value={reportType}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="exam">Экзамен</SelectItem>
                    <SelectItem value="test">Зачёт</SelectItem>
                    <SelectItem value="diff_test">Дифф. зачёт</SelectItem>
                  </SelectContent>
                </Select>
              </Field>
              <Field>
                <FieldLabel htmlFor={hoursID}>Часы</FieldLabel>
                <Input id={hoursID} name="hours" defaultValue={curriculum.hours} type="number" />
              </Field>
              <Field>
                <FieldLabel htmlFor={semesterID}>Семестр</FieldLabel>
                <Input
                  id={semesterID}
                  name="semester"
                  defaultValue={curriculum.semester}
                  type="number"
                />
              </Field>
            </FieldGroup>
            <DialogFooter>
              <DialogClose asChild={true}>
                <Button type="button" variant="outline">
                  Отмена
                </Button>
              </DialogClose>
              <Button disabled={isSubmitting} type="submit">
                {isSubmitting ? "Сохраняем..." : "Сохранить"}
              </Button>
            </DialogFooter>
          </form>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}

function EmptyState({ text }: { text: string }) {
  return (
    <div className="col-span-full rounded-lg border border-border p-8 text-center text-muted-foreground">
      {text}
    </div>
  );
}

function StaffPage({
  title,
  staffRole,
  profiles,
  isLoading,
  onCreated,
  onDelete,
}: {
  title: string;
  staffRole: "teacher" | "administrator";
  profiles: ApiProfile[];
  isLoading: boolean;
  onCreated: () => Promise<void>;
  onDelete: (profile: ApiProfile) => void;
}) {
  return (
    <div className="flex flex-col gap-6">
      <PageHeaderWithAction
        action={<StaffDialog onSaved={onCreated} role={staffRole} />}
        actionLabel={`Добавить ${staffRole === "teacher" ? "преподавателя" : "администратора"}`}
        description="Управление аккаунтами и профилями сотрудников"
        title={title}
      />
      <ProfilesTable isLoading={isLoading} onDelete={onDelete} profiles={profiles} />
    </div>
  );
}

function StaffDialog({
  role,
  onSaved,
}: {
  role: "teacher" | "administrator";
  onSaved: () => Promise<void>;
}) {
  const [open, setOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const lastNameID = useId();
  const firstNameID = useId();
  const middleNameID = useId();
  const emailID = useId();
  const passwordID = useId();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const lastName = String(formData.get("last_name") ?? "").trim();
    const firstName = String(formData.get("first_name") ?? "").trim();
    const middleName = String(formData.get("middle_name") ?? "").trim();
    const email = String(formData.get("email") ?? "").trim();
    const password = String(formData.get("password") ?? "");

    if (!lastName || !firstName || !email || !password) {
      toast.error("Заполните обязательные поля сотрудника.");
      return;
    }

    setIsSubmitting(true);
    try {
      const { data: user, error: userError } = await apiClient.POST("/users", {
        body: { email, password, role },
      });
      if (userError || !user?.id) {
        throw new Error("Не удалось создать аккаунт сотрудника.");
      }

      const body = {
        first_name: firstName,
        last_name: lastName,
        user_id: user.id,
        ...(middleName ? { middle_name: middleName } : {}),
      };
      const endpoint = role === "teacher" ? "/teachers" : "/administrators";
      const { error: profileError } = await apiClient.POST(endpoint, { body });
      if (profileError) {
        await apiClient.DELETE("/users/{id}", { params: { path: { id: user.id } } });
        throw new Error("Не удалось создать профиль сотрудника.");
      }

      toast.success("Сотрудник добавлен.");
      setOpen(false);
      await onSaved();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Не удалось добавить сотрудника.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild={true}>
        <Button className="h-9 rounded-lg px-3">
          <Plus className="size-4" />
          Добавить
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[520px]">
        <DialogHeader>
          <DialogTitle>Добавить сотрудника</DialogTitle>
          <DialogDescription>Создайте аккаунт и профиль сотрудника.</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <FieldGroup className="grid grid-cols-2 gap-4">
            <Field>
              <FieldLabel htmlFor={lastNameID}>Фамилия</FieldLabel>
              <Input id={lastNameID} name="last_name" placeholder="Иванов" />
            </Field>
            <Field>
              <FieldLabel htmlFor={firstNameID}>Имя</FieldLabel>
              <Input id={firstNameID} name="first_name" placeholder="Иван" />
            </Field>
            <Field>
              <FieldLabel htmlFor={middleNameID}>Отчество</FieldLabel>
              <Input id={middleNameID} name="middle_name" placeholder="Иванович" />
            </Field>
            <Field>
              <FieldLabel htmlFor={emailID}>E-mail</FieldLabel>
              <Input id={emailID} name="email" placeholder="teacher@example.com" type="email" />
            </Field>
            <Field>
              <FieldLabel htmlFor={passwordID}>Пароль</FieldLabel>
              <Input id={passwordID} minLength={8} name="password" type="password" />
            </Field>
          </FieldGroup>
          <DialogFooter>
            <DialogClose asChild={true}>
              <Button type="button" variant="outline">
                Отмена
              </Button>
            </DialogClose>
            <Button disabled={isSubmitting} type="submit">
              {isSubmitting ? "Добавляем..." : "Добавить"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function ProfilesTable({
  profiles,
  isLoading,
  onDelete,
}: {
  profiles: ApiProfile[];
  isLoading: boolean;
  onDelete: (profile: ApiProfile) => void;
}) {
  return (
    <div className="overflow-hidden rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ФИО</TableHead>
            <TableHead>E-mail</TableHead>
            <TableHead className="w-9" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading && <TableEmpty colSpan={3} text="Загружаем сотрудников..." />}
          {!isLoading && profiles.length === 0 && (
            <TableEmpty colSpan={3} text="Сотрудники пока не добавлены." />
          )}
          {profiles.map((profile) => (
            <TableRow key={profile.id}>
              <TableCell>{getProfileName(profile)}</TableCell>
              <TableCell>{profile.email}</TableCell>
              <TableCell className="text-right">
                <AdminActionsMenu
                  label={`Действия для ${getProfileName(profile)}`}
                  onDelete={() => onDelete(profile)}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

function SimpleEntityTable<T extends { id?: string }>({
  items,
  isLoading,
  titleHeader,
  getTitle,
  onEdit,
  onDelete,
}: {
  items: T[];
  isLoading: boolean;
  titleHeader: string;
  getTitle: (item: T) => string;
  onEdit?: (item: T) => void;
  onDelete: (item: T) => void;
}) {
  return (
    <div className="overflow-hidden rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{titleHeader}</TableHead>
            <TableHead className="w-9" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading && <TableEmpty colSpan={2} text="Загружаем данные..." />}
          {!isLoading && items.length === 0 && <TableEmpty colSpan={2} text="Записей пока нет." />}
          {items.map((item) => (
            <TableRow key={item.id}>
              <TableCell>{getTitle(item)}</TableCell>
              <TableCell className="text-right">
                <AdminActionsMenu
                  label={`Действия для ${getTitle(item)}`}
                  onEdit={onEdit ? () => onEdit(item) : undefined}
                  onDelete={() => onDelete(item)}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

function TableEmpty({ text, colSpan }: { text: string; colSpan: number }) {
  return (
    <TableRow>
      <TableCell className="h-24 text-center text-muted-foreground" colSpan={colSpan}>
        {text}
      </TableCell>
    </TableRow>
  );
}

function GroupDialog({
  specialties,
  onSaved,
}: {
  specialties: ApiSpecialty[];
  onSaved: () => Promise<void>;
}) {
  const [open, setOpen] = useState(false);
  const [specialtyID, setSpecialtyID] = useState("");
  const [studyForm, setStudyForm] = useState("full_time");
  const nameID = useId();
  const yearID = useId();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const name = String(formData.get("name") ?? "").trim();
    const admissionYear = Number(formData.get("admission_year"));

    if (!name || !specialtyID || !admissionYear) {
      toast.error("Заполните обязательные поля группы.");
      return;
    }

    const { error } = await apiClient.POST("/groups", {
      body: {
        admission_year: admissionYear,
        name,
        specialty_id: specialtyID,
        study_form: studyForm as components["schemas"]["repositories.StudyForm"],
      },
    });
    if (error) {
      toast.error("Не удалось добавить группу.");
      return;
    }

    toast.success("Группа добавлена.");
    setOpen(false);
    await onSaved();
  }

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild={true}>
        <Button className="h-9 rounded-lg px-3" disabled={specialties.length === 0}>
          <Plus className="size-4" />
          Добавить группу
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Добавить группу</DialogTitle>
          <DialogDescription>Создайте учебную группу внутри специальности.</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <FieldGroup className="gap-4">
            <Field>
              <FieldLabel htmlFor={nameID}>Название</FieldLabel>
              <Input id={nameID} name="name" placeholder="ЭФБО-01-24" />
            </Field>
            <Field>
              <FieldLabel htmlFor={yearID}>Год поступления</FieldLabel>
              <Input id={yearID} name="admission_year" placeholder="2024" type="number" />
            </Field>
            <Field>
              <FieldLabel>Форма обучения</FieldLabel>
              <Select onValueChange={setStudyForm} value={studyForm}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="full_time">Очная</SelectItem>
                  <SelectItem value="evening">Очно-заочная</SelectItem>
                  <SelectItem value="extramural">Заочная</SelectItem>
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel>Специальность</FieldLabel>
              <Select onValueChange={setSpecialtyID} value={specialtyID}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Выберите специальность" />
                </SelectTrigger>
                <SelectContent>
                  {specialties.map((specialty) => (
                    <SelectItem key={specialty.id} value={specialty.id ?? ""}>
                      {specialty.title}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
          </FieldGroup>
          <DialogFooter>
            <DialogClose asChild={true}>
              <Button type="button" variant="outline">
                Отмена
              </Button>
            </DialogClose>
            <Button type="submit">Добавить</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function GroupsTable({
  groups,
  specialties,
  isLoading,
  onEdit,
  onDelete,
}: {
  groups: ApiGroup[];
  specialties: ApiSpecialty[];
  isLoading: boolean;
  onEdit: (group: ApiGroup) => void;
  onDelete: (group: ApiGroup) => void;
}) {
  const specialtyMap = new Map(specialties.map((specialty) => [specialty.id, specialty.title]));

  return (
    <div className="overflow-hidden rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Группа</TableHead>
            <TableHead>Специальность</TableHead>
            <TableHead>Форма</TableHead>
            <TableHead>Год</TableHead>
            <TableHead className="w-9" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading && <TableEmpty colSpan={5} text="Загружаем группы..." />}
          {!isLoading && groups.length === 0 && (
            <TableEmpty colSpan={5} text="Группы пока не добавлены." />
          )}
          {groups.map((group) => (
            <TableRow key={group.id}>
              <TableCell>{group.name}</TableCell>
              <TableCell>{specialtyMap.get(group.specialty_id) ?? "Не указана"}</TableCell>
              <TableCell>{getStudyFormLabel(group.study_form)}</TableCell>
              <TableCell>{group.admission_year}</TableCell>
              <TableCell className="text-right">
                <AdminActionsMenu
                  label={`Действия для ${group.name}`}
                  onEdit={() => onEdit(group)}
                  onDelete={() => onDelete(group)}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

function CurriculumDialog({
  subjects,
  groups,
  teachers,
  onSaved,
}: {
  subjects: ApiSubject[];
  groups: ApiGroup[];
  teachers: ApiProfile[];
  onSaved: () => Promise<void>;
}) {
  const [open, setOpen] = useState(false);
  const [subjectID, setSubjectID] = useState("");
  const [groupID, setGroupID] = useState("");
  const [teacherUserID, setTeacherUserID] = useState("");
  const [reportType, setReportType] = useState("exam");
  const hoursID = useId();
  const semesterID = useId();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const hours = Number(formData.get("hours"));
    const semester = Number(formData.get("semester"));

    if (!subjectID || !groupID || !teacherUserID || !hours || !semester) {
      toast.error("Заполните обязательные поля учебного плана.");
      return;
    }

    const { error } = await apiClient.POST("/curriculums", {
      body: {
        group_id: groupID,
        hours,
        lead_by: teacherUserID,
        report_type: reportType as components["schemas"]["repositories.CurriculumReportType"],
        semester,
        subject_id: subjectID,
      },
    });
    if (error) {
      toast.error("Не удалось добавить учебный план.");
      return;
    }

    toast.success("Учебный план добавлен.");
    setOpen(false);
    await onSaved();
  }

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild={true}>
        <Button
          className="h-9 rounded-lg px-3"
          disabled={!subjects.length || !groups.length || !teachers.length}
        >
          <Plus className="size-4" />
          Добавить учебный план
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Добавить учебный план</DialogTitle>
          <DialogDescription>Назначьте дисциплину группе и преподавателю.</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <FieldGroup className="gap-4">
            <Field>
              <FieldLabel>Дисциплина</FieldLabel>
              <Select onValueChange={setSubjectID} value={subjectID}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Выберите дисциплину" />
                </SelectTrigger>
                <SelectContent>
                  {subjects.map((subject) => (
                    <SelectItem key={subject.id} value={subject.id ?? ""}>
                      {subject.title}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel>Группа</FieldLabel>
              <Select onValueChange={setGroupID} value={groupID}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Выберите группу" />
                </SelectTrigger>
                <SelectContent>
                  {groups.map((group) => (
                    <SelectItem key={group.id} value={group.id ?? ""}>
                      {group.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel>Преподаватель</FieldLabel>
              <Select onValueChange={setTeacherUserID} value={teacherUserID}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Выберите преподавателя" />
                </SelectTrigger>
                <SelectContent>
                  {teachers.map((teacher) => (
                    <SelectItem key={teacher.id} value={teacher.user_id ?? ""}>
                      {getProfileName(teacher)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel>Тип отчётности</FieldLabel>
              <Select onValueChange={setReportType} value={reportType}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="exam">Экзамен</SelectItem>
                  <SelectItem value="test">Зачёт</SelectItem>
                  <SelectItem value="diff_test">Дифф. зачёт</SelectItem>
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel htmlFor={hoursID}>Часы</FieldLabel>
              <Input id={hoursID} name="hours" type="number" />
            </Field>
            <Field>
              <FieldLabel htmlFor={semesterID}>Семестр</FieldLabel>
              <Input id={semesterID} name="semester" type="number" />
            </Field>
          </FieldGroup>
          <DialogFooter>
            <DialogClose asChild={true}>
              <Button type="button" variant="outline">
                Отмена
              </Button>
            </DialogClose>
            <Button type="submit">Добавить</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function CurriculumsTable({
  curriculums,
  subjects,
  groups,
  teachers,
  isLoading,
  onEdit,
  onDelete,
}: {
  curriculums: ApiCurriculum[];
  subjects: ApiSubject[];
  groups: ApiGroup[];
  teachers: ApiProfile[];
  isLoading: boolean;
  onEdit: (curriculum: ApiCurriculum) => void;
  onDelete: (curriculum: ApiCurriculum) => void;
}) {
  const subjectMap = new Map(subjects.map((subject) => [subject.id, subject.title]));
  const groupMap = new Map(groups.map((group) => [group.id, group.name]));
  const teacherMap = new Map(teachers.map((teacher) => [teacher.user_id, getProfileName(teacher)]));

  return (
    <div className="overflow-hidden rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Дисциплина</TableHead>
            <TableHead>Группа</TableHead>
            <TableHead>Преподаватель</TableHead>
            <TableHead>Семестр</TableHead>
            <TableHead>Часы</TableHead>
            <TableHead>Отчётность</TableHead>
            <TableHead className="w-9" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading && <TableEmpty colSpan={7} text="Загружаем учебные планы..." />}
          {!isLoading && curriculums.length === 0 && (
            <TableEmpty colSpan={7} text="Учебные планы пока не добавлены." />
          )}
          {curriculums.map((curriculum) => (
            <TableRow key={curriculum.id}>
              <TableCell>{subjectMap.get(curriculum.subject_id) ?? "Не указана"}</TableCell>
              <TableCell>{groupMap.get(curriculum.group_id) ?? "Не указана"}</TableCell>
              <TableCell>{teacherMap.get(curriculum.lead_by) ?? "Не указан"}</TableCell>
              <TableCell>{curriculum.semester}</TableCell>
              <TableCell>{curriculum.hours}</TableCell>
              <TableCell>{getReportTypeLabel(curriculum.report_type)}</TableCell>
              <TableCell className="text-right">
                <AdminActionsMenu
                  label="Действия для учебного плана"
                  onEdit={() => onEdit(curriculum)}
                  onDelete={() => onDelete(curriculum)}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

async function loadStudentRecords() {
  const [studentsResponse, groupsResponse, gradesResponse] = await Promise.all([
    apiClient.GET("/students"),
    apiClient.GET("/groups"),
    apiClient.GET("/grades"),
  ]);

  if (studentsResponse.error || groupsResponse.error || gradesResponse.error) {
    throw new Error("Не удалось загрузить данные студентов.");
  }

  const students = studentsResponse.data ?? [];
  const groups = groupsResponse.data ?? [];
  const grades = gradesResponse.data ?? [];
  const groupMap = new Map(groups.map((group) => [group.id, group.name]));
  const gradesByStudent = groupGradesByStudent(grades);
  const users = await loadUsersByID(students.map((student) => student.user_id).filter(isFilled));
  const records = students.map((student) =>
    toStudentRecord(
      student,
      groupMap,
      users.get(student.user_id ?? ""),
      gradesByStudent.get(student.id),
    ),
  );

  return { groups, records };
}

async function loadDisciplineRecords() {
  const [subjectsResponse, curriculumsResponse, groupsResponse, studentsResponse] =
    await Promise.all([
      apiClient.GET("/subjects"),
      apiClient.GET("/curriculums"),
      apiClient.GET("/groups"),
      apiClient.GET("/students"),
    ]);

  if (
    subjectsResponse.error ||
    curriculumsResponse.error ||
    groupsResponse.error ||
    studentsResponse.error
  ) {
    throw new Error("Не удалось загрузить данные дисциплин.");
  }

  const subjects = subjectsResponse.data ?? [];
  const curriculums = curriculumsResponse.data ?? [];
  const groups = groupsResponse.data ?? [];
  const students = studentsResponse.data ?? [];
  const groupMap = new Map(groups.map((group) => [group.id, group.name]));
  const groupSizes = new Map<string, number>();
  for (const student of students) {
    if (!student.group_id) {
      continue;
    }
    groupSizes.set(student.group_id, (groupSizes.get(student.group_id) ?? 0) + 1);
  }

  return subjects.map((subject) => {
    const subjectCurriculums = curriculums.filter(
      (curriculum) => curriculum.subject_id === subject.id,
    );
    const groupIDs = new Set(
      subjectCurriculums.map((curriculum) => curriculum.group_id).filter(isFilled),
    );
    const reportTypes = new Set(
      subjectCurriculums.map((curriculum) => curriculum.report_type).filter(isFilled),
    );
    const hours = subjectCurriculums.reduce((sum, curriculum) => sum + (curriculum.hours ?? 0), 0);
    const studentsCount = Array.from(groupIDs).reduce(
      (sum, groupID) => sum + (groupSizes.get(groupID) ?? 0),
      0,
    );
    const groupsLabel = Array.from(groupIDs)
      .map((groupID) => groupMap.get(groupID))
      .filter(isFilled)
      .join(", ");

    return {
      groups: groupsLabel || "Группы не назначены",
      hours,
      id: subject.id ?? "",
      reportTypes:
        Array.from(reportTypes).filter(isFilled).map(getReportTypeLabel).join(", ") ||
        "Без отчётности",
      students: studentsCount,
      title: subject.title ?? "Без названия",
    };
  });
}

async function loadUsersByID(userIDs: string[]) {
  const uniqueIDs = Array.from(new Set(userIDs));
  const entries = await Promise.all(
    uniqueIDs.map(async (id) => {
      const { data } = await apiClient.GET("/users/{id}", {
        params: { path: { id } },
      });

      return [id, data] as const;
    }),
  );

  return new Map(entries.filter((entry): entry is readonly [string, ApiUser] => Boolean(entry[1])));
}

function toStudentRecord(
  student: ApiStudent,
  groupMap: Map<string | undefined, string | undefined>,
  user: ApiUser | undefined,
  grades: ApiGrade[] | undefined,
): StudentRecord {
  const name = getFullName(student);

  return {
    average: grades?.length ? getAverageGrade(grades) : "-",
    email: user?.email ?? "Нет данных",
    firstName: student.first_name ?? "",
    group: groupMap.get(student.group_id) ?? "Без группы",
    groupId: student.group_id ?? "",
    id: student.id ?? "",
    initials: getInitials(student),
    lastName: student.last_name ?? "",
    middleName: student.middle_name,
    name,
    status: "Активен",
    userId: student.user_id ?? "",
  };
}

function groupGradesByStudent(grades: ApiGrade[]) {
  const groups = new Map<string | undefined, ApiGrade[]>();
  for (const grade of grades) {
    const studentGrades = groups.get(grade.student_id) ?? [];
    studentGrades.push(grade);
    groups.set(grade.student_id, studentGrades);
  }

  return groups;
}

function getAverageGrade(grades: ApiGrade[]) {
  if (grades.length === 0) {
    return "-";
  }

  const average = grades.reduce((sum, grade) => sum + (grade.value ?? 0), 0) / grades.length;
  return average.toFixed(1);
}

function formatAverage(value: number | undefined) {
  return typeof value === "number" ? value.toFixed(1) : "-";
}

function toAnalyticsMetrics(overview: ApiAnalyticsOverview | null): Metric[] {
  if (!overview) {
    return administratorAnalyticsMetrics;
  }

  return [
    {
      title: "Всего оценок",
      value: String(overview.grades_count ?? 0),
      description: "В журнале оценок",
      icon: BookOpen,
      tone: "positive",
    },
    {
      title: "Средний балл",
      value: formatAverage(overview.average_grade),
      description: "По всем оценкам",
      icon: BarChart3,
      tone: "positive",
    },
    {
      title: "Учебные планы",
      value: String(overview.curriculums_count ?? 0),
      description: "Связки групп и дисциплин",
      icon: Users,
      tone: "positive",
    },
  ];
}

function toGradeDistributionData(overview: ApiAnalyticsOverview | null) {
  if (!overview?.grade_distribution?.length) {
    return gradeDistributionData;
  }

  const fills = new Map([
    [5, "var(--color-excellent)"],
    [4, "var(--color-good)"],
    [3, "var(--color-satisfactory)"],
    [2, "var(--color-poor)"],
  ]);

  return overview.grade_distribution.map((item) => ({
    fill: fills.get(item.value ?? 0) ?? "var(--chart-5)",
    grade: String(item.value ?? "-"),
    value: item.count ?? 0,
  }));
}

function toSubjectPerformanceData(overview: ApiAnalyticsOverview | null) {
  if (!overview?.subject_averages?.length) {
    return subjectPerformanceData;
  }

  return overview.subject_averages.map((item) => ({
    average: Number((item.average_grade ?? 0).toFixed(1)),
    subject: item.subject_title ?? "Без названия",
  }));
}

function getFullName(student: ApiStudent) {
  return [student.last_name, student.first_name, student.middle_name].filter(isFilled).join(" ");
}

function getInitials(student: ApiStudent) {
  const initials = [student.last_name, student.first_name]
    .filter(isFilled)
    .map((name) => name.at(0)?.toLocaleUpperCase("ru-RU"))
    .join("");

  return initials || "??";
}

function getProfileName(profile: ApiProfile) {
  return [profile.last_name, profile.first_name, profile.middle_name].filter(isFilled).join(" ");
}

function getStudyFormLabel(studyForm: string | undefined) {
  if (studyForm === "full_time") {
    return "Очная";
  }
  if (studyForm === "evening") {
    return "Очно-заочная";
  }
  if (studyForm === "extramural") {
    return "Заочная";
  }

  return "Не указана";
}

function getReportTypeLabel(reportType: string | undefined) {
  if (reportType === "exam") {
    return "Экзамен";
  }
  if (reportType === "test") {
    return "Зачёт";
  }
  if (reportType === "diff_test") {
    return "Дифф. зачёт";
  }

  return reportType;
}

function isFilled(value: string | undefined): value is string {
  return typeof value === "string" && value.length > 0;
}

function PieChartPreview({
  data = gradeDistributionData,
}: {
  data?: typeof gradeDistributionData;
}) {
  return (
    <ChartContainer className="mx-auto h-[300px] w-full" config={gradeChartConfig}>
      <PieChart accessibilityLayer={true}>
        <ChartTooltip content={<ChartTooltipContent hideLabel={true} />} />
        <Pie data={data} dataKey="value" nameKey="grade">
          {data.map((entry) => (
            <Cell fill={entry.fill} key={entry.grade} />
          ))}
        </Pie>
        <ChartLegend content={<ChartLegendContent nameKey="grade" />} />
      </PieChart>
    </ChartContainer>
  );
}

function BarChartPreview({
  data = subjectPerformanceData,
}: {
  data?: typeof subjectPerformanceData;
}) {
  return (
    <ChartContainer className="h-[300px] w-full" config={subjectChartConfig}>
      <BarChart accessibilityLayer={true} data={data} margin={{ left: 8, right: 8 }}>
        <CartesianGrid vertical={false} />
        <XAxis dataKey="subject" tickLine={false} />
        <YAxis domain={[0, 5]} tickCount={3} tickLine={false} width={32} />
        <ChartTooltip content={<ChartTooltipContent />} />
        <Bar dataKey="average" fill="var(--color-average)" radius={[4, 4, 0, 0]} />
      </BarChart>
    </ChartContainer>
  );
}
