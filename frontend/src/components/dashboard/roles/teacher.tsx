import {
  BookOpen,
  CalendarDays,
  Clock3,
  Ellipsis,
  Filter,
  GraduationCap,
  Plus,
  Search,
  Users,
} from "lucide-react";
import type { FormEvent, ReactNode } from "react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import type { components } from "@/api";
import { apiClient } from "@/api";
import type { Metric } from "@/components/dashboard/dashboard-widgets";
import {
  DashboardSection,
  MetricsGrid,
  PageHeading,
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
import type { AuthProfile, AuthUser } from "@/lib/context/auth";

export type TeacherPage = "home" | "disciplines" | "grades" | "schedule";

type ApiCurriculum = components["schemas"]["dto.Curriculum"];
type ApiGrade = components["schemas"]["dto.Grade"];
type ApiGroup = components["schemas"]["dto.Group"];
type ApiStudent = components["schemas"]["dto.Student"];
type ApiSubject = components["schemas"]["dto.Subject"];

type TeacherWorkspace = {
  courses: TeacherCourse[];
  grades: TeacherGradeRecord[];
  groups: ApiGroup[];
  students: ApiStudent[];
  curriculums: ApiCurriculum[];
  isLoading: boolean;
  reload: () => Promise<void>;
};

type TeacherCourse = {
  id: string;
  curriculumID: string;
  title: string;
  group: string;
  groupID: string;
  reportType: string;
  semester: number;
  hours: number;
  students: number;
  grades: number;
  average: string;
};

type TeacherGradeRecord = {
  id: string;
  curriculumID: string;
  studentID: string;
  studentName: string;
  courseTitle: string;
  groupName: string;
  value: number;
  comment: string;
};

type SimpleStat = {
  label: string;
  value: string;
  tone?: "default" | "warning";
};

export function TeacherDashboard({
  page,
  profile,
  user,
}: {
  page: TeacherPage;
  profile: AuthProfile | null;
  user: AuthUser | null;
}) {
  const workspace = useTeacherWorkspace(user?.id);

  if (page === "disciplines") {
    return <TeacherDisciplinesPage workspace={workspace} />;
  }

  if (page === "grades") {
    return <TeacherGradesPage workspace={workspace} />;
  }

  if (page === "schedule") {
    return <TeacherSchedulePage workspace={workspace} />;
  }

  return <TeacherHomePage profile={profile} workspace={workspace} />;
}

function TeacherHomePage({
  profile,
  workspace,
}: {
  profile: AuthProfile | null;
  workspace: TeacherWorkspace;
}) {
  const metrics = useMemo(() => getTeacherMetrics(workspace), [workspace]);
  const todayTitle = profile ? `Добро пожаловать, ${profile.first_name}!` : "Добро пожаловать!";

  return (
    <div className="flex flex-col gap-6">
      <PageHeading description="Сводка по вашим дисциплинам и оценкам" title={todayTitle} />
      <MetricsGrid metrics={metrics} />
      <DashboardSection
        className="h-[414px]"
        description="Дисциплины из назначенных учебных планов"
        title="Мои дисциплины"
      >
        <CompactCourseList
          courses={workspace.courses.slice(0, 5)}
          isLoading={workspace.isLoading}
        />
      </DashboardSection>
    </div>
  );
}

function TeacherDisciplinesPage({ workspace }: { workspace: TeacherWorkspace }) {
  return (
    <div className="flex flex-col gap-6 pb-4">
      <TeacherPageHeader
        description="Дисциплины и группы, назначенные вам в учебных планах"
        title="Мои дисциплины"
      />
      <div className="flex flex-col gap-4">
        {workspace.isLoading && <EmptyState text="Загружаем дисциплины..." />}
        {!workspace.isLoading && workspace.courses.length === 0 && (
          <EmptyState text="Вам пока не назначены учебные планы." />
        )}
        {workspace.courses.map((course) => (
          <TeacherCourseCard course={course} key={course.id} />
        ))}
      </div>
    </div>
  );
}

function TeacherGradesPage({ workspace }: { workspace: TeacherWorkspace }) {
  const [search, setSearch] = useState("");
  const [courseID, setCourseID] = useState("all");

  const filteredGrades = useMemo(() => {
    const query = search.trim().toLowerCase();

    return workspace.grades.filter((grade) => {
      const matchesCourse = courseID === "all" || grade.curriculumID === courseID;
      const matchesQuery =
        !query ||
        [grade.studentName, grade.courseTitle, grade.groupName].some((value) =>
          value.toLowerCase().includes(query),
        );

      return matchesCourse && matchesQuery;
    });
  }, [courseID, search, workspace.grades]);

  const stats = useMemo(() => getTeacherGradeStats(workspace), [workspace]);

  return (
    <div className="flex flex-col gap-6">
      <TeacherPageHeader
        action={
          <GradeDialog
            curriculums={workspace.curriculums}
            courses={workspace.courses}
            onSaved={workspace.reload}
            students={workspace.students}
          />
        }
        description="Выставление и просмотр оценок по вашим дисциплинам"
        title="Оценки"
      />
      <TeacherGradeFilters
        courseID={courseID}
        courses={workspace.courses}
        onCourseChange={setCourseID}
        onSearchChange={setSearch}
        search={search}
      />
      <SimpleStatsGrid columns={4} stats={stats} />
      <GradesTable
        grades={filteredGrades}
        isLoading={workspace.isLoading}
        onChanged={workspace.reload}
      />
    </div>
  );
}

function TeacherSchedulePage({ workspace }: { workspace: TeacherWorkspace }) {
  const stats = useMemo(() => getTeacherScheduleStats(workspace), [workspace]);

  return (
    <div className="flex flex-col gap-6 pb-4">
      <TeacherPageHeader description="Расписание строится из учебных планов" title="Расписание" />
      <SimpleStatsGrid cardClassName="h-36" stats={stats} />
      <SchedulePlaceholder courses={workspace.courses} isLoading={workspace.isLoading} />
    </div>
  );
}

function TeacherPageHeader({
  title,
  description,
  action,
}: {
  title: string;
  description: string;
  action?: ReactNode;
}) {
  return (
    <div className="flex h-16 items-center justify-between">
      <div className="flex flex-col gap-1">
        <h1 className="h-9 text-[30px] font-semibold leading-9 text-foreground">{title}</h1>
        <p className="h-6 text-base leading-6 text-muted-foreground">{description}</p>
      </div>
      {action}
    </div>
  );
}

function TeacherCourseCard({ course }: { course: TeacherCourse }) {
  return (
    <Card className="gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="gap-2 px-6 pb-0 pt-4">
        <div className="flex items-center gap-3">
          <CardTitle className="text-xl font-medium leading-7 text-card-foreground">
            {course.title}
          </CardTitle>
          <Badge className="h-[22px] rounded-lg px-[9px]" variant="outline">
            {course.group}
          </Badge>
        </div>
        <CardDescription className="flex items-center gap-3 text-base leading-6 text-muted-foreground">
          <span>Семестр {course.semester}</span>
          <span className="flex items-center gap-1">
            <Clock3 className="size-4" />
            {course.hours} часов
          </span>
          <span>{course.reportType}</span>
        </CardDescription>
        <CardAction>
          <TeacherActionsMenu label={`Действия для курса ${course.title}`} />
        </CardAction>
      </CardHeader>
      <CardContent className="px-6 pb-4 pt-4">
        <div className="grid grid-cols-3 gap-4">
          <CourseMetric label="Студентов" value={course.students} />
          <CourseMetric label="Оценок" value={course.grades} />
          <CourseMetric label="Средний балл" value={course.average} />
        </div>
      </CardContent>
    </Card>
  );
}

function CourseMetric({ label, value }: { label: string; value: number | string }) {
  return (
    <div className="flex flex-col gap-1">
      <div className="text-2xl font-bold leading-8 text-card-foreground">{value}</div>
      <div className="text-xs leading-4 text-muted-foreground">{label}</div>
    </div>
  );
}

function TeacherGradeFilters({
  courses,
  courseID,
  search,
  onCourseChange,
  onSearchChange,
}: {
  courses: TeacherCourse[];
  courseID: string;
  search: string;
  onCourseChange: (value: string) => void;
  onSearchChange: (value: string) => void;
}) {
  return (
    <div className="flex h-9 gap-3">
      <Select onValueChange={onCourseChange} value={courseID}>
        <SelectTrigger className="h-9 w-[300px] rounded-lg bg-muted">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Все дисциплины</SelectItem>
          {courses.map((course) => (
            <SelectItem key={course.curriculumID} value={course.curriculumID}>
              {course.title} &middot; {course.group}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <div className="relative min-w-0 flex-1">
        <Search className="-translate-y-1/2 absolute left-3 top-1/2 size-4 text-muted-foreground" />
        <Input
          aria-label="Поиск студентов"
          className="h-9 border-transparent bg-muted pl-9 text-sm shadow-none placeholder:text-muted-foreground focus-visible:ring-0"
          onChange={(event) => onSearchChange(event.target.value)}
          placeholder="Поиск студентов..."
          value={search}
        />
      </div>
      <Button className="h-9 rounded-lg px-[13px]" disabled={true} variant="outline">
        <Filter className="size-4" />
        Фильтры
      </Button>
    </div>
  );
}

function GradeDialog({
  courses,
  curriculums,
  students,
  onSaved,
}: {
  courses: TeacherCourse[];
  curriculums: ApiCurriculum[];
  students: ApiStudent[];
  onSaved: () => Promise<void>;
}) {
  const [open, setOpen] = useState(false);
  const [curriculumID, setCurriculumID] = useState("");
  const [studentID, setStudentID] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const valueID = "teacher-grade-value";
  const commentID = "teacher-grade-comment";

  const selectedCurriculum = curriculums.find((curriculum) => curriculum.id === curriculumID);
  const availableStudents = selectedCurriculum
    ? students.filter((student) => student.group_id === selectedCurriculum.group_id)
    : [];

  useEffect(() => {
    if (!selectedCurriculum || availableStudents.some((student) => student.id === studentID)) {
      return;
    }

    setStudentID("");
  }, [availableStudents, selectedCurriculum, studentID]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const value = Number(formData.get("value"));
    const comment = String(formData.get("comment") ?? "").trim();

    if (!curriculumID || !studentID || !value) {
      toast.error("Выберите дисциплину, студента и оценку.");
      return;
    }

    setIsSubmitting(true);
    const { error } = await apiClient.POST("/grades", {
      body: {
        curriculum_id: curriculumID,
        student_id: studentID,
        value,
        ...(comment ? { comment } : {}),
      },
    });
    setIsSubmitting(false);

    if (error) {
      toast.error("Не удалось сохранить оценку.");
      return;
    }

    toast.success("Оценка сохранена.");
    setOpen(false);
    setCurriculumID("");
    setStudentID("");
    await onSaved();
  }

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild={true}>
        <Button className="h-9 rounded-lg px-3" disabled={courses.length === 0}>
          <Plus className="size-4" />
          Новая оценка
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[520px]">
        <DialogHeader>
          <DialogTitle>Новая оценка</DialogTitle>
          <DialogDescription>Выставьте оценку студенту по своей дисциплине.</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <FieldGroup className="gap-4">
            <Field>
              <FieldLabel>Дисциплина</FieldLabel>
              <Select onValueChange={setCurriculumID} value={curriculumID}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Выберите дисциплину" />
                </SelectTrigger>
                <SelectContent>
                  {courses.map((course) => (
                    <SelectItem key={course.curriculumID} value={course.curriculumID}>
                      {course.title} &middot; {course.group}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel>Студент</FieldLabel>
              <Select disabled={!selectedCurriculum} onValueChange={setStudentID} value={studentID}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Выберите студента" />
                </SelectTrigger>
                <SelectContent>
                  {availableStudents.map((student) => (
                    <SelectItem key={student.id} value={student.id ?? ""}>
                      {getStudentName(student)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel htmlFor={valueID}>Оценка</FieldLabel>
              <Input id={valueID} max={5} min={2} name="value" type="number" />
            </Field>
            <Field>
              <FieldLabel htmlFor={commentID}>Комментарий</FieldLabel>
              <Input id={commentID} name="comment" placeholder="Необязательно" />
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
      </DialogContent>
    </Dialog>
  );
}

function GradesTable({
  grades,
  isLoading,
  onChanged,
}: {
  grades: TeacherGradeRecord[];
  isLoading: boolean;
  onChanged: () => Promise<void>;
}) {
  return (
    <div className="overflow-hidden rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Студент</TableHead>
            <TableHead>Дисциплина</TableHead>
            <TableHead>Группа</TableHead>
            <TableHead>Оценка</TableHead>
            <TableHead>Комментарий</TableHead>
            <TableHead className="w-9" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading && <TableEmpty colSpan={6} text="Загружаем оценки..." />}
          {!isLoading && grades.length === 0 && <TableEmpty colSpan={6} text="Оценок пока нет." />}
          {grades.map((grade) => (
            <TableRow key={grade.id}>
              <TableCell>{grade.studentName}</TableCell>
              <TableCell>{grade.courseTitle}</TableCell>
              <TableCell>{grade.groupName}</TableCell>
              <TableCell>
                <Badge className="rounded-lg">{grade.value}</Badge>
              </TableCell>
              <TableCell>{grade.comment || "—"}</TableCell>
              <TableCell className="text-right">
                <TeacherActionsMenu
                  label={`Действия для оценки ${grade.studentName}`}
                  onDelete={async () => {
                    const { error } = await apiClient.DELETE("/grades/{id}", {
                      params: { path: { id: grade.id } },
                    });
                    if (error) {
                      toast.error("Не удалось удалить оценку.");
                      return;
                    }
                    toast.success("Оценка удалена.");
                    await onChanged();
                  }}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
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
          <CardContent>
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

function CompactCourseList({
  courses,
  isLoading,
}: {
  courses: TeacherCourse[];
  isLoading: boolean;
}) {
  if (isLoading) {
    return <EmptyState text="Загружаем дисциплины..." />;
  }
  if (courses.length === 0) {
    return <EmptyState text="Назначенных дисциплин пока нет." />;
  }

  return (
    <div className="flex flex-col gap-3">
      {courses.map((course) => (
        <div
          className="flex items-center justify-between rounded-lg border border-border p-3"
          key={course.id}
        >
          <div>
            <div className="font-medium text-foreground">{course.title}</div>
            <div className="text-sm text-muted-foreground">{course.group}</div>
          </div>
          <Badge variant="outline">{course.average}</Badge>
        </div>
      ))}
    </div>
  );
}

function SchedulePlaceholder({
  courses,
  isLoading,
}: {
  courses: TeacherCourse[];
  isLoading: boolean;
}) {
  if (isLoading) {
    return <EmptyState text="Загружаем учебные планы..." />;
  }

  return (
    <Card className="gap-0 rounded-[14px] py-0 ring-border">
      <CardHeader className="px-6 pb-0 pt-6">
        <CardTitle className="text-base font-medium text-card-foreground">
          Нагрузка по учебным планам
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-3 px-6 pb-6 pt-6">
        {courses.length === 0 && <EmptyState text="Учебные планы пока не назначены." />}
        {courses.map((course) => (
          <div
            className="flex items-center justify-between rounded-lg border border-border p-4"
            key={course.id}
          >
            <div>
              <div className="font-medium">{course.title}</div>
              <div className="text-sm text-muted-foreground">{course.group}</div>
            </div>
            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <span>{course.hours} ч.</span>
              <span>{course.semester} семестр</span>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
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

function EmptyState({ text }: { text: string }) {
  return (
    <div className="rounded-lg border border-border p-8 text-center text-muted-foreground">
      {text}
    </div>
  );
}

function TeacherActionsMenu({ label, onDelete }: { label: string; onDelete?: () => void }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild={true}>
        <Button aria-label={label} className="size-9 rounded-lg" size="icon" variant="ghost">
          <Ellipsis className="size-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-44">
        <DropdownMenuItem disabled={true}>Открыть</DropdownMenuItem>
        <DropdownMenuItem disabled={true}>Редактировать</DropdownMenuItem>
        <DropdownMenuItem disabled={true}>Материалы</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem disabled={!onDelete} onClick={onDelete} variant="destructive">
          Удалить
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function useTeacherWorkspace(teacherUserID: string | undefined): TeacherWorkspace {
  const [courses, setCourses] = useState<TeacherCourse[]>([]);
  const [grades, setGrades] = useState<TeacherGradeRecord[]>([]);
  const [groups, setGroups] = useState<ApiGroup[]>([]);
  const [students, setStudents] = useState<ApiStudent[]>([]);
  const [curriculums, setCurriculums] = useState<ApiCurriculum[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadWorkspace = useCallback(async () => {
    if (!teacherUserID) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    const [
      curriculumsResponse,
      subjectsResponse,
      groupsResponse,
      studentsResponse,
      gradesResponse,
    ] = await Promise.all([
      apiClient.GET("/curriculums"),
      apiClient.GET("/subjects"),
      apiClient.GET("/groups"),
      apiClient.GET("/students"),
      apiClient.GET("/grades"),
    ]);

    if (
      curriculumsResponse.error ||
      subjectsResponse.error ||
      groupsResponse.error ||
      studentsResponse.error ||
      gradesResponse.error
    ) {
      toast.error("Не удалось загрузить данные преподавателя.");
      setIsLoading(false);
      return;
    }

    const teacherCurriculums = (curriculumsResponse.data ?? []).filter(
      (curriculum) => curriculum.lead_by === teacherUserID,
    );
    const loadedStudents = studentsResponse.data ?? [];
    const loadedGrades = (gradesResponse.data ?? []).filter((grade) =>
      teacherCurriculums.some((curriculum) => curriculum.id === grade.curriculum_id),
    );

    setCurriculums(teacherCurriculums);
    setGroups(groupsResponse.data ?? []);
    setStudents(loadedStudents);
    setCourses(
      toTeacherCourses(
        teacherCurriculums,
        subjectsResponse.data ?? [],
        groupsResponse.data ?? [],
        loadedStudents,
        loadedGrades,
      ),
    );
    setGrades(
      toTeacherGradeRecords(
        loadedGrades,
        teacherCurriculums,
        subjectsResponse.data ?? [],
        groupsResponse.data ?? [],
        loadedStudents,
      ),
    );
    setIsLoading(false);
  }, [teacherUserID]);

  useEffect(() => {
    void loadWorkspace();
  }, [loadWorkspace]);

  return {
    courses,
    curriculums,
    grades,
    groups,
    isLoading,
    reload: loadWorkspace,
    students,
  };
}

function toTeacherCourses(
  curriculums: ApiCurriculum[],
  subjects: ApiSubject[],
  groups: ApiGroup[],
  students: ApiStudent[],
  grades: ApiGrade[],
) {
  const subjectMap = new Map(subjects.map((subject) => [subject.id, subject.title]));
  const groupMap = new Map(groups.map((group) => [group.id, group.name]));

  return curriculums.map((curriculum) => {
    const courseStudents = students.filter((student) => student.group_id === curriculum.group_id);
    const courseGrades = grades.filter((grade) => grade.curriculum_id === curriculum.id);

    return {
      average: getAverage(courseGrades),
      curriculumID: curriculum.id ?? "",
      grades: courseGrades.length,
      group: groupMap.get(curriculum.group_id) ?? "Группа не указана",
      groupID: curriculum.group_id ?? "",
      hours: curriculum.hours ?? 0,
      id: curriculum.id ?? "",
      reportType: getReportTypeLabel(curriculum.report_type),
      semester: curriculum.semester ?? 0,
      students: courseStudents.length,
      title: subjectMap.get(curriculum.subject_id) ?? "Дисциплина не указана",
    };
  });
}

function toTeacherGradeRecords(
  grades: ApiGrade[],
  curriculums: ApiCurriculum[],
  subjects: ApiSubject[],
  groups: ApiGroup[],
  students: ApiStudent[],
) {
  const curriculumMap = new Map(curriculums.map((curriculum) => [curriculum.id, curriculum]));
  const subjectMap = new Map(subjects.map((subject) => [subject.id, subject.title]));
  const groupMap = new Map(groups.map((group) => [group.id, group.name]));
  const studentMap = new Map(students.map((student) => [student.id, student]));

  return grades.map((grade) => {
    const curriculum = curriculumMap.get(grade.curriculum_id);
    const student = studentMap.get(grade.student_id);

    return {
      comment: grade.comment ?? "",
      courseTitle: subjectMap.get(curriculum?.subject_id) ?? "Дисциплина не указана",
      curriculumID: grade.curriculum_id ?? "",
      groupName: groupMap.get(curriculum?.group_id) ?? "Группа не указана",
      id: grade.id ?? "",
      studentID: grade.student_id ?? "",
      studentName: student ? getStudentName(student) : "Студент не указан",
      value: grade.value ?? 0,
    };
  });
}

function getTeacherMetrics(workspace: TeacherWorkspace): Metric[] {
  const groupIDs = new Set(workspace.courses.map((course) => course.groupID));
  const average = getAverage(workspace.grades.map((grade) => ({ value: grade.value })));

  return [
    {
      title: "Мои курсы",
      value: String(workspace.courses.length),
      description: "Назначенных учебных планов",
      icon: BookOpen,
    },
    {
      title: "Студенты",
      value: String(workspace.students.length),
      description: `В ${groupIDs.size} группах`,
      icon: Users,
    },
    {
      title: "Оценки",
      value: String(workspace.grades.length),
      description: "Выставлено по курсам",
      icon: GraduationCap,
    },
    {
      title: "Средний балл",
      value: average,
      description: "По вашим дисциплинам",
      icon: CalendarDays,
    },
  ];
}

function getTeacherGradeStats(workspace: TeacherWorkspace): SimpleStat[] {
  return [
    { label: "Всего студентов", value: String(workspace.students.length) },
    {
      label: "Средний балл",
      value: getAverage(workspace.grades.map((grade) => ({ value: grade.value }))),
    },
    { label: "Оценок", value: String(workspace.grades.length) },
    { label: "Курсов", value: String(workspace.courses.length) },
  ];
}

function getTeacherScheduleStats(workspace: TeacherWorkspace): SimpleStat[] {
  const hours = workspace.courses.reduce((sum, course) => sum + course.hours, 0);

  return [
    { label: "Часов в планах", value: String(hours) },
    { label: "Учебных планов", value: String(workspace.courses.length) },
    {
      label: "Групп",
      value: String(new Set(workspace.courses.map((course) => course.groupID)).size),
    },
  ];
}

function getStudentName(student: ApiStudent) {
  return [student.last_name, student.first_name, student.middle_name].filter(isFilled).join(" ");
}

function getAverage(grades: { value?: number }[]) {
  if (grades.length === 0) {
    return "-";
  }

  const average = grades.reduce((sum, grade) => sum + (grade.value ?? 0), 0) / grades.length;
  return average.toFixed(1);
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

function isFilled(value: string | undefined): value is string {
  return typeof value === "string" && value.length > 0;
}
