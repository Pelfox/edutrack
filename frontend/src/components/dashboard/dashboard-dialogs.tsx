import type { ReactNode } from "react";

import { Button } from "@/components/ui/button";
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
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";

export type DashboardDialogField = {
  label: string;
  max?: number;
  min?: number;
  name: string;
  placeholder: string;
  type?: string;
};

export type DashboardDialogConfig = {
  title: string;
  description: string;
  submitLabel: string;
  fields: DashboardDialogField[];
};

export function DashboardActionDialog({
  config,
  trigger,
}: {
  config: DashboardDialogConfig;
  trigger: ReactNode;
}) {
  return (
    <Dialog>
      <DialogTrigger asChild={true}>{trigger}</DialogTrigger>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>{config.title}</DialogTitle>
          <DialogDescription>{config.description}</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-5" onSubmit={(event) => event.preventDefault()}>
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
            <DialogClose asChild={true}>
              <Button type="submit">{config.submitLabel}</Button>
            </DialogClose>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
