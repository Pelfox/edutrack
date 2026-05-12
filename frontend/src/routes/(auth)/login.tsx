import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { LogInIcon } from "lucide-react";
import { useId } from "react";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod/v3";
import { Button } from "@/components/ui/button";
import { Field, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { AuthUser } from "@/lib/context/auth";
import { useAuth } from "@/lib/context/auth";
import type { DashboardPage } from "@/lib/dashboard-routing";
import { getDashboardRole, isDashboardPageAvailable } from "@/lib/dashboard-routing";

const loginSchema = z.object({
  email: z.string().email("Введите корректный e-mail адрес."),
  password: z.string().min(1, "Введите пароль."),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export const Route = createFileRoute("/(auth)/login")({
  validateSearch: (
    search,
  ): {
    redirect?: string;
  } => {
    if (typeof search.redirect === "string") {
      return {
        redirect: search.redirect,
      };
    }

    return {};
  },
  component: LoginPage,
});

function LoginPage() {
  const auth = useAuth();
  const router = useRouter();
  const search = Route.useSearch();
  const emailId = useId();
  const passwordId = useId();

  const form = useForm<LoginFormValues>({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (values: LoginFormValues) => {
    try {
      const user = await auth.login({
        email: values.email.trim(),
        password: values.password,
      });
      router.history.push(getRedirectPath(search.redirect, user));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Не удалось войти в аккаунт.");
    }
  };

  return (
    <div className="w-full h-full flex flex-col gap-12 items-center justify-center">
      {/* Приветствие пользователя */}
      <div className="text-center space-y-2">
        <h1 className="text-2xl font-medium">Добро пожаловать</h1>
        <p className="text-muted-foreground">Для продолжения необходимо войти в аккаунт.</p>
      </div>

      {/* Форма для ввода данных аккаунта */}
      <form className="w-xs space-y-6" noValidate={true} onSubmit={form.handleSubmit(onSubmit)}>
        <Controller
          name="email"
          control={form.control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <Label htmlFor={emailId}>E-mail адрес</Label>
              <Input
                aria-invalid={fieldState.invalid}
                autoComplete="email"
                disabled={form.formState.isSubmitting}
                id={emailId}
                placeholder="someone@example.com"
                type="email"
                {...field}
              />
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <Controller
          name="password"
          control={form.control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <Label htmlFor={passwordId}>Пароль</Label>
              <Input
                aria-invalid={fieldState.invalid}
                autoComplete="current-password"
                disabled={form.formState.isSubmitting}
                id={passwordId}
                type="password"
                {...field}
              />
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <Button className="w-full" disabled={form.formState.isSubmitting} type="submit">
          <LogInIcon />
          {form.formState.isSubmitting ? "Загрузка..." : "Войти в аккаунт"}
        </Button>
      </form>
    </div>
  );
}

function getRedirectPath(redirect: string | undefined, user: AuthUser | null) {
  if (redirect?.startsWith("/") !== true || redirect.startsWith("//")) {
    return "/";
  }

  const page = getDashboardPageByPath(redirect);

  if (!page) {
    return "/";
  }

  const role = getDashboardRole(user?.role);

  return isDashboardPageAvailable(role, page) ? redirect : "/";
}

function getDashboardPageByPath(path: string): DashboardPage | null {
  const pathname = path.split(/[?#]/)[0];

  if (pathname === "/") {
    return "home";
  }

  if (pathname === "/students") {
    return "students";
  }

  if (pathname === "/teachers") {
    return "teachers";
  }

  if (pathname === "/groups") {
    return "groups";
  }

  if (pathname === "/specialties") {
    return "specialties";
  }

  if (pathname === "/disciplines") {
    return "disciplines";
  }

  if (pathname === "/curriculums") {
    return "curriculums";
  }

  if (pathname === "/analytics") {
    return "analytics";
  }

  if (pathname === "/grades") {
    return "grades";
  }

  if (pathname === "/schedule") {
    return "schedule";
  }

  return null;
}
