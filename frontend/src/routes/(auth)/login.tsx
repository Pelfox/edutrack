import { createFileRoute } from "@tanstack/react-router";
import { LogInIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Field } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export const Route = createFileRoute("/(auth)/login")({
  component: LoginPage,
});

function LoginPage() {
  return (
    <div className="w-full h-full flex flex-col gap-12 items-center justify-center">
      {/* Приветствие пользователя */}
      <div className="text-center space-y-2">
        <h1 className="text-2xl font-medium">Добро пожаловать</h1>
        <p className="text-muted-foreground">Для продолжения необходимо войти в аккаунт.</p>
      </div>

      {/* Форма для ввода данных аккаунта */}
      <form onSubmit={(event) => event.preventDefault()} className="w-xs space-y-6">
        <Field>
          <Label>E-mail адрес</Label>
          <Input type="email" placeholder="someone@example.com" />
        </Field>

        <Field>
          <Label>Пароль</Label>
          <Input type="password" />
        </Field>

        <Button type="submit" className="w-full">
          <LogInIcon />
          Войти в аккаунт
        </Button>
      </form>
    </div>
  );
}
