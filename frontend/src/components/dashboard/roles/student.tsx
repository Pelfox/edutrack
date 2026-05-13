import type { LucideIcon } from "lucide-react";
import {
  Award,
  BarChart3,
  BookOpen,
  Clock3,
  GraduationCap,
  MapPin,
  TrendingUp,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import type { components } from "@/api";
import { apiClient } from "@/api";
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
import type { AuthProfile } from "@/lib/context/auth";

export type StudentPage = "home" | "disciplines" | "curriculums" | "grades" | "schedule";

type ApiCurriculum = components["schemas"]["dto.Curriculum"];
type ApiGrade = components["schemas"]["dto.Grade"];
type ApiGroup = components["schemas"]["dto.Group"];
type ApiProfile = components["schemas"]["dto.Profile"];
type ApiStudent = components["schemas"]["dto.Student"];
type ApiSubject = components["schemas"]["dto.Subject"];

type StudentWorkspace = {
  courses: StudentCourse[];
  grades: StudentGradeRecord[];
  isLoading: boolean;
  student: ApiStudent | null;
};

type StudentCourse = {
  id: string;
  title: string;
  group: string;
  teacher: string;
  average: string;
  gradesCount: number;
  hours: number;
  reportType: string;
  semester: number;
};

type StudentGradeRecord = {
  id: string;
  courseTitle: string;
  groupName: string;
  value: number;
  comment: string;
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
  average: string;
  rows: StudentGradeRow[];
};

type StudentGradeRow = {
  assignment: string;
  comment: string;
  grade: string;
  tone: "positive" | "blue";
};

export function StudentDashboard({
  page,
  profile,
}: {
  page: StudentPage;
  profile: AuthProfile | null;
}) {
  const workspace = useStudentWorkspace();

  if (page === "disciplines") {
    return <StudentCoursesPage workspace={workspace} />;
  }

  if (page === "grades") {
    return <StudentGradesPage workspace={workspace} />;
  }

  if (page === "schedule") {
    return <StudentCurriculumsPage workspace={workspace} title="Расписание" />;
  }

  if (page === "curriculums") {
    return <StudentCurriculumsPage workspace={workspace} title="Учебные планы" />;
  }

  return <StudentHomePage profile={profile} workspace={workspace} />;
}

function StudentHomePage({
  profile,
  workspace,
}: {
  profile: AuthProfile | null;
  workspace: StudentWorkspace;
}) {
  const metrics = useMemo(() => getStudentMetrics(workspace), [workspace]);
  const recentGrades = useMemo(() => toRecentGradeItems(workspace.grades), [workspace.grades]);
  const courseItems = useMemo(() => toCourseItems(workspace.courses), [workspace.courses]);
  const title = profile ? `Привет, ${profile.first_name}!` : "Привет!";

  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="Сводка по вашим дисциплинам и оценкам" title={title} />
      <MetricsGrid metrics={metrics} />
      <div className="grid grid-cols-2 items-start gap-4">
        <DashboardSection
          className="h-[390px]"
          description="Учебная нагрузка из ваших учебных планов"
          title="Расписание"
        >
          <ScheduleList items={toScheduleItems(workspace.courses)} />
        </DashboardSection>
        <DashboardSection description="Недавно выставленные" title="Последние оценки">
          <GradesList items={recentGrades} />
        </DashboardSection>
      </div>
      <DashboardSection className="min-h-[454px]" description="Прогресс обучения" title="Мои курсы">
        {workspace.isLoading ? (
          <EmptyState text="Загружаем курсы..." />
        ) : (
          <CourseGrid items={courseItems} />
        )}
      </DashboardSection>
    </div>
  );
}

function StudentCoursesPage({ workspace }: { workspace: StudentWorkspace }) {
  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="Ваши дисциплины из учебного плана группы" title="Мои курсы" />
      <div className="flex flex-col gap-4">
        {workspace.isLoading && <EmptyState text="Загружаем курсы..." />}
        {!workspace.isLoading && workspace.courses.length === 0 && (
          <EmptyState text="Для вашей группы пока нет учебных планов." />
        )}
        {workspace.courses.map((course) => (
          <StudentCourseCard course={course} key={course.id} />
        ))}
      </div>
    </div>
  );
}

function StudentGradesPage({ workspace }: { workspace: StudentWorkspace }) {
  const summary = useMemo(() => getStudentGradeSummary(workspace.grades), [workspace.grades]);
  const gradeCourses = useMemo(
    () => toStudentGradeCourses(workspace.courses, workspace.grades),
    [workspace.courses, workspace.grades],
  );

  return (
    <div className="flex flex-col gap-6 pb-4">
      <PageHeading description="Ваша успеваемость по всем курсам" title="Мои оценки" />
      <div className="grid grid-cols-4 gap-4">
        {summary.map((item) => (
          <StudentGradeSummaryCard item={item} key={item.title} />
        ))}
      </div>
      <div className="flex flex-col gap-4">
        {workspace.isLoading && <EmptyState text="Загружаем оценки..." />}
        {!workspace.isLoading && gradeCourses.length === 0 && (
          <EmptyState text="Оценок пока нет." />
        )}
        {gradeCourses.map((course) => (
          <StudentGradeCourseCard course={course} key={course.title} />
        ))}
      </div>
    </div>
  );
}

function StudentCurriculumsPage({
  workspace,
  title,
}: {
  workspace: StudentWorkspace;
  title: string;
}) {
  return (
    <div className="flex flex-col gap-6 pb-4">
      <PageHeading description="Учебная нагрузка из ваших учебных планов" title={title} />
      <StudentScheduleCard courses={workspace.courses} isLoading={workspace.isLoading} />
    </div>
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
            {course.teacher} · {course.group} · {course.semester} семестр · {course.reportType}
          </CardDescription>
        </div>
        <CardAction>
          <div className="flex size-14 items-center justify-center rounded-[10px] bg-green-100 text-xl font-bold leading-7 text-green-700">
            {course.average}
          </div>
        </CardAction>
      </CardHeader>
      <CardContent className="flex flex-col gap-6 px-6 pb-6 pt-6">
        <div className="grid grid-cols-3 gap-4">
          <CourseMetric label="Часов" value={course.hours} />
          <CourseMetric label="Оценок" value={course.gradesCount} />
          <CourseMetric label="Средний балл" value={course.average} />
        </div>
      </CardContent>
    </Card>
  );
}

function CourseMetric({ label, value }: { label: string; value: number | string }) {
  return (
    <div className="flex flex-col gap-1 rounded-lg border border-border p-3">
      <div className="text-2xl font-bold leading-8 text-card-foreground">{value}</div>
      <div className="text-xs leading-4 text-muted-foreground">{label}</div>
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
            {course.rows.length} оценок
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
                <TableHead className="w-1/2">Запись</TableHead>
                <TableHead className="w-[24%] text-center">Комментарий</TableHead>
                <TableHead className="w-[26%] text-center">Оценка</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {course.rows.map((row) => (
                <TableRow className="h-[57px]" key={row.assignment}>
                  <TableCell className="font-medium">{row.assignment}</TableCell>
                  <TableCell className="text-center text-muted-foreground">
                    {row.comment || "—"}
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

function StudentScheduleCard({
  courses,
  isLoading,
}: {
  courses: StudentCourse[];
  isLoading: boolean;
}) {
  return (
    <Card className="gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="px-6 pb-0 pt-6">
        <div className="flex flex-col">
          <CardTitle className="text-base font-medium leading-6 text-card-foreground">
            Учебная нагрузка
          </CardTitle>
          <CardDescription className="text-base leading-6 text-muted-foreground">
            По дисциплинам вашей группы
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent className="flex flex-col gap-3 px-6 pb-6 pt-6">
        {isLoading && <EmptyState text="Загружаем учебные планы..." />}
        {!isLoading && courses.length === 0 && (
          <EmptyState text="Учебные планы пока не назначены." />
        )}
        {courses.map((course) => (
          <div
            className="flex items-center justify-between rounded-lg border border-border p-4"
            key={course.id}
          >
            <div>
              <div className="font-medium">{course.title}</div>
              <div className="text-sm text-muted-foreground">
                {course.teacher} · {course.group}
              </div>
            </div>
            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <span>{course.hours} ч.</span>
              <span>{course.semester} семестр</span>
              <span>{course.reportType}</span>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}

function EmptyState({ text }: { text: string }) {
  return (
    <div className="rounded-lg border border-border p-8 text-center text-muted-foreground">
      {text}
    </div>
  );
}

function useStudentWorkspace(): StudentWorkspace {
  const [student, setStudent] = useState<ApiStudent | null>(null);
  const [courses, setCourses] = useState<StudentCourse[]>([]);
  const [grades, setGrades] = useState<StudentGradeRecord[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadWorkspace = useCallback(async () => {
    setIsLoading(true);
    const [
      studentResponse,
      curriculumsResponse,
      subjectsResponse,
      groupsResponse,
      gradesResponse,
      teachersResponse,
    ] = await Promise.all([
      apiClient.GET("/students/me"),
      apiClient.GET("/curriculums"),
      apiClient.GET("/subjects"),
      apiClient.GET("/groups"),
      apiClient.GET("/grades"),
      apiClient.GET("/teachers"),
    ]);

    if (
      studentResponse.error ||
      curriculumsResponse.error ||
      subjectsResponse.error ||
      groupsResponse.error ||
      gradesResponse.error ||
      teachersResponse.error
    ) {
      toast.error("Не удалось загрузить данные студента.");
      setIsLoading(false);
      return;
    }

    const currentStudent = studentResponse.data ?? null;
    const studentCurriculums = (curriculumsResponse.data ?? []).filter(
      (curriculum) => curriculum.group_id === currentStudent?.group_id,
    );
    const studentGrades = gradesResponse.data ?? [];
    const nextCourses = toStudentCourses(
      studentCurriculums,
      subjectsResponse.data ?? [],
      groupsResponse.data ?? [],
      teachersResponse.data ?? [],
      studentGrades,
    );

    setStudent(currentStudent);
    setCourses(nextCourses);
    setGrades(
      toStudentGradeRecords(
        studentGrades,
        studentCurriculums,
        subjectsResponse.data ?? [],
        groupsResponse.data ?? [],
      ),
    );
    setIsLoading(false);
  }, []);

  useEffect(() => {
    void loadWorkspace();
  }, [loadWorkspace]);

  return { courses, grades, isLoading, student };
}

function toStudentCourses(
  curriculums: ApiCurriculum[],
  subjects: ApiSubject[],
  groups: ApiGroup[],
  teachers: ApiProfile[],
  grades: ApiGrade[],
) {
  const subjectMap = new Map(subjects.map((subject) => [subject.id, subject.title]));
  const groupMap = new Map(groups.map((group) => [group.id, group.name]));
  const teacherMap = new Map(teachers.map((teacher) => [teacher.user_id, getProfileName(teacher)]));

  return curriculums.map((curriculum) => {
    const courseGrades = grades.filter((grade) => grade.curriculum_id === curriculum.id);

    return {
      average: getAverage(courseGrades),
      gradesCount: courseGrades.length,
      group: groupMap.get(curriculum.group_id) ?? "Группа не указана",
      hours: curriculum.hours ?? 0,
      id: curriculum.id ?? "",
      reportType: getReportTypeLabel(curriculum.report_type),
      semester: curriculum.semester ?? 0,
      teacher: teacherMap.get(curriculum.lead_by) ?? "Преподаватель не указан",
      title: subjectMap.get(curriculum.subject_id) ?? "Дисциплина не указана",
    };
  });
}

function toStudentGradeRecords(
  grades: ApiGrade[],
  curriculums: ApiCurriculum[],
  subjects: ApiSubject[],
  groups: ApiGroup[],
) {
  const curriculumMap = new Map(curriculums.map((curriculum) => [curriculum.id, curriculum]));
  const subjectMap = new Map(subjects.map((subject) => [subject.id, subject.title]));
  const groupMap = new Map(groups.map((group) => [group.id, group.name]));

  return grades.map((grade) => {
    const curriculum = curriculumMap.get(grade.curriculum_id);

    return {
      comment: grade.comment ?? "",
      courseTitle: subjectMap.get(curriculum?.subject_id) ?? "Дисциплина не указана",
      groupName: groupMap.get(curriculum?.group_id) ?? "Группа не указана",
      id: grade.id ?? "",
      value: grade.value ?? 0,
    };
  });
}

function getStudentMetrics(workspace: StudentWorkspace): Metric[] {
  return [
    {
      title: "Средний балл",
      value: getAverage(workspace.grades.map((grade) => ({ value: grade.value }))),
      description: "По всем оценкам",
      icon: GraduationCap,
      tone: "positive",
    },
    {
      title: "Мои курсы",
      value: String(workspace.courses.length),
      description: "По учебному плану группы",
      icon: BookOpen,
    },
    {
      title: "Оценки",
      value: String(workspace.grades.length),
      description: "Всего выставлено",
      icon: BarChart3,
    },
    {
      title: "Часы",
      value: String(workspace.courses.reduce((sum, course) => sum + course.hours, 0)),
      description: "В учебных планах",
      icon: Clock3,
    },
  ];
}

function getStudentGradeSummary(grades: StudentGradeRecord[]): StudentGradeSummary[] {
  const excellent = grades.filter((grade) => grade.value === 5).length;
  const good = grades.filter((grade) => grade.value === 4).length;

  return [
    {
      title: "Средний балл",
      value: getAverage(grades.map((grade) => ({ value: grade.value }))),
      detail: grades.length > 0 ? "По всем оценкам" : "Оценок пока нет",
      icon: Award,
      tone: "positive",
    },
    {
      title: "Пятёрок",
      value: String(excellent),
      detail: `${getPercent(excellent, grades.length)} от всех оценок`,
      icon: TrendingUp,
    },
    {
      title: "Четвёрок",
      value: String(good),
      detail: `${getPercent(good, grades.length)} от всех оценок`,
      icon: BookOpen,
    },
    { title: "Всего оценок", value: String(grades.length), detail: "За всё время" },
  ];
}

function toStudentGradeCourses(courses: StudentCourse[], grades: StudentGradeRecord[]) {
  return courses
    .map((course) => {
      const rows = grades
        .filter((grade) => grade.courseTitle === course.title)
        .map((grade, index) => ({
          assignment: `Оценка №${index + 1}`,
          comment: grade.comment,
          grade: String(grade.value),
          tone: grade.value >= 5 ? ("positive" as const) : ("blue" as const),
        }));

      return {
        average: course.average,
        rows,
        title: course.title,
      };
    })
    .filter((course) => course.rows.length > 0);
}

function toRecentGradeItems(grades: StudentGradeRecord[]): GradeItem[] {
  if (grades.length === 0) {
    return [];
  }

  return grades.slice(0, 5).map((grade) => ({
    grade: String(grade.value),
    subject: grade.courseTitle,
    time: "Недавно",
    title: grade.comment || "Оценка",
    ...(grade.value >= 5 ? {} : { tone: "blue" as const }),
  }));
}

function toCourseItems(courses: StudentCourse[]): CourseItem[] {
  return courses.map((course) => ({
    href: "/curriculums",
    next: `${course.semester} семестр`,
    score: course.average,
    teacher: course.teacher,
    title: course.title,
  }));
}

function toScheduleItems(courses: StudentCourse[]): ScheduleItem[] {
  if (courses.length === 0) {
    return [
      {
        id: "student-schedule-placeholder",
        meta: [
          { icon: GraduationCap, label: "Учебный план не назначен" },
          { icon: MapPin, label: "Группа не указана" },
        ],
        status: "План",
        time: "—",
        title: "Расписание пока недоступно",
      },
    ];
  }

  return courses.slice(0, 3).map((course) => ({
    id: course.id,
    meta: [
      { icon: GraduationCap, label: course.group },
      { icon: MapPin, label: course.teacher },
    ],
    status: course.reportType,
    time: `${course.semester} сем.`,
    title: course.title,
  }));
}

function getAverage(grades: { value?: number }[]) {
  if (grades.length === 0) {
    return "-";
  }

  const average = grades.reduce((sum, grade) => sum + (grade.value ?? 0), 0) / grades.length;
  return average.toFixed(1);
}

function getPercent(count: number, total: number) {
  if (total === 0) {
    return "0%";
  }

  return `${Math.round((count / total) * 100)}%`;
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

  return "Отчётность не указана";
}

function getProfileName(profile: ApiProfile) {
  return [profile.last_name, profile.first_name, profile.middle_name].filter(isFilled).join(" ");
}

function isFilled(value: string | undefined): value is string {
  return typeof value === "string" && value.length > 0;
}
