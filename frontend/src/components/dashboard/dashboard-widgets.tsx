import { Link } from "@tanstack/react-router";
import type { LucideIcon } from "lucide-react";
import { CalendarDays } from "lucide-react";
import type { ReactNode } from "react";

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
import { Progress } from "@/components/ui/progress";

type Tone = "positive" | "negative" | "muted" | "blue";

export type Metric = {
  title: string;
  value: string;
  description: string;
  icon: LucideIcon;
  tone?: Tone;
};

export type ScheduleItem = {
  id: string;
  time: string;
  title: string;
  meta: {
    icon: LucideIcon;
    label: string;
  }[];
  status: string;
  statusVariant?: "primary" | "secondary" | "outline";
};

export type WorkItem = {
  id: string;
  student: string;
  task: string;
  subject: string;
  time: string;
};

export type GradeItem = {
  title: string;
  subject: string;
  time: string;
  grade: string;
  tone?: Tone;
};

export type CourseItem = {
  title: string;
  teacher: string;
  score: string;
  next: string;
  href?: string;
  tone?: Tone;
};

export type TopStudent = {
  name: string;
  group: string;
  score: string;
  progress: number;
};

export type ActivityItem = {
  id: string;
  student: string;
  subject: string;
  grade: string;
  time: string;
  tone?: Tone;
};

export function PageHeading({ title, description }: { title: string; description: string }) {
  return (
    <div className="flex h-16 w-full flex-col gap-1">
      <h1 className="h-9 text-[30px] font-semibold leading-9 tracking-normal text-foreground">
        {title}
      </h1>
      <p className="h-6 text-base leading-6 text-muted-foreground">{description}</p>
    </div>
  );
}

export function MetricsGrid({ metrics }: { metrics: Metric[] }) {
  return (
    <div
      className="grid w-full gap-4"
      style={{ gridTemplateColumns: `repeat(${metrics.length}, 1fr)` }}
    >
      {metrics.map((metric) => (
        <MetricCard key={metric.title} metric={metric} />
      ))}
    </div>
  );
}

function MetricCard({ metric }: { metric: Metric }) {
  const Icon = metric.icon;

  return (
    <Card className="h-[150px] gap-6 rounded-[14px] py-0 ring-border">
      <CardHeader className="h-[52px] px-6 pb-2 pt-6">
        <CardTitle className="text-sm font-medium leading-5 text-card-foreground">
          {metric.title}
        </CardTitle>
        <CardAction>
          <Icon className="size-4 text-muted-foreground" />
        </CardAction>
      </CardHeader>
      <CardContent className="px-6">
        <div className="text-2xl font-bold leading-8 tracking-normal text-card-foreground">
          {metric.value}
        </div>
        <div className={`text-xs leading-4 ${toneTextClass(metric.tone)}`}>
          {metric.description}
        </div>
      </CardContent>
    </Card>
  );
}

export function DashboardSection({
  title,
  description,
  className,
  children,
}: {
  title: string;
  description: string;
  className?: string;
  children: ReactNode;
}) {
  return (
    <Card className={`rounded-[14px] py-0 ring-border ${className ?? ""}`}>
      <CardHeader className="h-[70px] px-6 pb-0 pt-6">
        <CardTitle className="text-base font-medium leading-4 text-card-foreground">
          {title}
        </CardTitle>
        <CardDescription className="text-base leading-6 text-muted-foreground">
          {description}
        </CardDescription>
      </CardHeader>
      <CardContent className="px-6 pb-6">{children}</CardContent>
    </Card>
  );
}

export function ScheduleList({ items }: { items: ScheduleItem[] }) {
  return (
    <div className="flex flex-col gap-3">
      {items.map((item) => (
        <div
          className="flex h-[90px] items-center justify-between rounded-[10px] border border-border p-[13px]"
          key={item.id}
        >
          <div className="flex h-16 items-center gap-3">
            <div className="flex size-16 shrink-0 items-center justify-center rounded-full bg-muted text-xs leading-4 text-muted-foreground">
              {item.time}
            </div>
            <div className="min-w-0">
              <div className="truncate text-base font-medium leading-6 text-foreground">
                {item.title}
              </div>
              <div className="flex flex-wrap items-center gap-3 text-sm leading-5 text-muted-foreground">
                {item.meta.map((metaItem) => (
                  <span className="flex min-w-0 items-center gap-1" key={metaItem.label}>
                    <metaItem.icon className="size-3.5 shrink-0" />
                    <span className="truncate">{metaItem.label}</span>
                  </span>
                ))}
              </div>
            </div>
          </div>
          <StatusBadge variant={item.statusVariant}>{item.status}</StatusBadge>
        </div>
      ))}
    </div>
  );
}

export function ReviewList({ items }: { items: WorkItem[] }) {
  return (
    <div className="flex flex-col gap-4">
      {items.map((item) => (
        <div
          className="flex min-h-[70px] items-center justify-between border-border pb-[13px] last:border-b-0 last:pb-0 [&:not(:last-child)]:border-b"
          key={item.id}
        >
          <div className="flex flex-col gap-1">
            <div className="text-sm font-medium leading-5 text-foreground">{item.student}</div>
            <div className="text-sm leading-5 text-muted-foreground">{item.task}</div>
            <div className="flex items-center gap-2">
              <Badge className="h-[22px] rounded-lg px-[9px] py-[3px]" variant="outline">
                {item.subject}
              </Badge>
              <span className="text-xs leading-4 text-muted-foreground">{item.time}</span>
            </div>
          </div>
          <Button className="h-8 rounded-lg px-[13px]" size="default" variant="outline">
            Проверить
          </Button>
        </div>
      ))}
    </div>
  );
}

export function GradesList({ items }: { items: GradeItem[] }) {
  return (
    <div className="flex flex-col gap-4">
      {items.map((item) => (
        <div className="flex items-center justify-between" key={`${item.title}-${item.grade}`}>
          <div className="flex flex-col gap-1">
            <div className="text-sm font-medium leading-5 text-foreground">{item.title}</div>
            <div className="text-sm leading-5 text-muted-foreground">{item.subject}</div>
            <div className="text-xs leading-4 text-muted-foreground">{item.time}</div>
          </div>
          <GradePill grade={item.grade} tone={item.tone} />
        </div>
      ))}
    </div>
  );
}

export function CourseGrid({ items }: { items: CourseItem[] }) {
  return (
    <div className="grid grid-cols-2 gap-4">
      {items.map((item) => (
        <div
          className="flex h-[132px] flex-col justify-between rounded-[10px] border border-border px-[17px] pb-[17px] pt-[17px]"
          key={item.title}
        >
          <div className="flex h-12 items-start justify-between">
            <div className="min-w-0">
              <h3 className="truncate text-base font-medium leading-6 text-foreground">
                {item.title}
              </h3>
              <p className="truncate text-sm leading-5 text-muted-foreground">{item.teacher}</p>
            </div>
            <GradePill grade={item.score} tone={item.tone} />
          </div>
          <div className="flex h-8 items-center justify-between">
            <div className="flex items-center gap-1 text-sm leading-5 text-muted-foreground">
              <CalendarDays className="size-3.5" />
              <span>{item.next}</span>
            </div>
            {item.href ? (
              <Button
                asChild={true}
                className="h-8 px-3 text-muted-foreground"
                size="default"
                variant="ghost"
              >
                <Link to={item.href}>Открыть</Link>
              </Button>
            ) : (
              <Button className="h-8 px-3 text-muted-foreground" size="default" variant="ghost">
                Открыть
              </Button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}

export function TopStudentsList({ items }: { items: TopStudent[] }) {
  return (
    <div className="flex flex-col gap-4">
      {items.map((item) => (
        <div className="flex flex-col gap-2" key={`${item.name}-${item.score}`}>
          <div className="flex h-[34px] items-center justify-between">
            <div className="flex items-center gap-3">
              <Avatar size="default">
                <AvatarFallback className="bg-primary text-sm font-medium text-primary-foreground">
                  {item.name[0]}
                </AvatarFallback>
              </Avatar>
              <div>
                <div className="text-sm font-medium leading-[14px] text-foreground">
                  {item.name}
                </div>
                <div className="mt-1 text-xs leading-4 text-muted-foreground">{item.group}</div>
              </div>
            </div>
            <div className="text-sm font-bold leading-5 text-foreground">{item.score}</div>
          </div>
          <Progress className="h-1.5 bg-primary/20" value={item.progress} />
        </div>
      ))}
    </div>
  );
}

export function ActivityList({ items }: { items: ActivityItem[] }) {
  return (
    <div className="flex flex-col gap-4">
      {items.map((item) => (
        <div
          className="flex min-h-[51px] items-center justify-between border-border pb-[13px] last:border-b-0 last:pb-0 [&:not(:last-child)]:border-b"
          key={item.id}
        >
          <div>
            <div className="text-sm font-medium leading-[14px] text-foreground">{item.student}</div>
            <div className="mt-1 text-sm leading-5 text-muted-foreground">{item.subject}</div>
          </div>
          <div className="flex w-[140px] items-center gap-3">
            <GradePill grade={item.grade} tone={item.tone} />
            <span className="flex-1 text-right text-xs leading-4 text-muted-foreground">
              {item.time}
            </span>
          </div>
        </div>
      ))}
    </div>
  );
}

function GradePill({ grade, tone = "positive" }: { grade: string; tone?: Tone | undefined }) {
  return (
    <div
      className={`flex size-10 shrink-0 items-center justify-center rounded-[10px] text-base font-bold leading-6 ${tonePillClass(
        tone,
      )}`}
    >
      {grade}
    </div>
  );
}

function StatusBadge({
  variant = "outline",
  children,
}: {
  variant?: "primary" | "secondary" | "outline" | undefined;
  children: ReactNode;
}) {
  if (variant === "primary") {
    return (
      <Badge className="h-[22px] rounded-lg bg-primary px-[9px] py-[3px] text-primary-foreground">
        {children}
      </Badge>
    );
  }

  if (variant === "secondary") {
    return (
      <Badge className="h-[22px] rounded-lg px-[9px] py-[3px]" variant="secondary">
        {children}
      </Badge>
    );
  }

  return (
    <Badge className="h-[22px] rounded-lg px-[9px] py-[3px]" variant="outline">
      {children}
    </Badge>
  );
}

function toneTextClass(tone: Tone = "muted") {
  if (tone === "positive") {
    return "text-green-600";
  }

  if (tone === "negative") {
    return "text-destructive";
  }

  if (tone === "blue") {
    return "text-blue-600";
  }

  return "text-muted-foreground";
}

function tonePillClass(tone: Tone) {
  if (tone === "blue") {
    return "bg-blue-100 text-blue-700";
  }

  if (tone === "negative") {
    return "bg-destructive/10 text-destructive";
  }

  if (tone === "muted") {
    return "bg-muted text-muted-foreground";
  }

  return "bg-green-100 text-green-700";
}
