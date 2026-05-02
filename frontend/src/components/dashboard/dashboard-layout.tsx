import type { LucideIcon } from "lucide-react";
import { Bell, Menu, Search } from "lucide-react";
import type { CSSProperties, ReactNode } from "react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarSeparator,
} from "@/components/ui/sidebar";

export type DashboardRole = "administrator" | "teacher" | "student";

export type DashboardNavItem = {
  icon: LucideIcon;
  label: string;
  active?: boolean;
  href?: string;
};

export type DashboardUser = {
  initials: string;
  name: string;
  detail: string;
  nameWeight?: "semibold" | "bold";
};

export type DashboardConfig = {
  navItems: DashboardNavItem[];
  searchPlaceholder: string;
  user: DashboardUser;
};

const sidebarWidth = "min(100vw, 256px)";

export function DashboardShell({
  config,
  children,
}: {
  config: DashboardConfig;
  children: ReactNode;
}) {
  return (
    <SidebarProvider
      className="h-svh min-h-0 w-full overflow-hidden bg-background"
      defaultOpen={true}
      style={{ "--sidebar-width": sidebarWidth } as CSSProperties}
    >
      <DashboardSidebar config={config} />
      <main className="flex h-svh min-h-0 min-w-0 flex-1 flex-col overflow-hidden bg-background">
        <DashboardTopBar placeholder={config.searchPlaceholder} />
        <div className="min-h-0 flex-1 overflow-auto px-6 pt-6">
          <div className="w-full min-w-[980px]">{children}</div>
        </div>
      </main>
    </SidebarProvider>
  );
}

function DashboardSidebar({ config }: { config: DashboardConfig }) {
  return (
    <Sidebar
      collapsible="none"
      className="h-svh border-r border-sidebar-border bg-sidebar pr-px text-sidebar-foreground"
    >
      <SidebarHeader className="h-16 w-[255px] flex-row items-center justify-between gap-3 border-b border-sidebar-border px-4 py-0">
        <div className="min-w-px flex-1 text-lg font-semibold leading-7 text-sidebar-foreground">
          EduTrack
        </div>
        <Button
          aria-label="Свернуть боковую панель"
          className="size-9 rounded-lg px-2.5 text-sidebar-foreground hover:bg-transparent"
          size="icon"
          variant="ghost"
        >
          <Menu className="size-4" strokeWidth={2} />
        </Button>
      </SidebarHeader>

      <SidebarContent className="w-[255px] bg-sidebar px-3 pt-3">
        <SidebarGroup className="p-0">
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              {config.navItems.map((item) => (
                <SidebarMenuItem key={item.label}>
                  <DashboardSidebarItem item={item} />
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarSeparator className="mx-0 w-[255px] bg-sidebar-border" />

      <SidebarFooter className="h-[81px] w-[255px] px-3 pb-0 pt-[13px]">
        <UserSummary user={config.user} />
      </SidebarFooter>
    </Sidebar>
  );
}

function DashboardSidebarItem({ item }: { item: DashboardNavItem }) {
  const content = (
    <>
      <item.icon className={item.active === true ? "size-6" : "size-5"} strokeWidth={2} />
      <span>{item.label}</span>
    </>
  );

  if (item.href !== undefined) {
    return (
      <SidebarMenuButton
        asChild={true}
        className="h-10 gap-3 rounded-[10px] px-3 py-2.5 text-sm font-medium leading-5 text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground data-active:bg-sidebar-primary data-active:text-sidebar-primary-foreground [&_svg]:size-5 [&_svg]:shrink-0"
        isActive={item.active === true}
        size="default"
      >
        <a href={item.href}>{content}</a>
      </SidebarMenuButton>
    );
  }

  return (
    <SidebarMenuButton
      className="h-10 gap-3 rounded-[10px] px-3 py-2.5 text-sm font-medium leading-5 text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground data-active:bg-sidebar-primary data-active:text-sidebar-primary-foreground [&_svg]:size-5 [&_svg]:shrink-0"
      isActive={item.active === true}
      size="default"
    >
      {content}
    </SidebarMenuButton>
  );
}

function UserSummary({ user }: { user: DashboardUser }) {
  return (
    <div className="relative h-14 w-full">
      <Avatar className="absolute left-3 top-3" size="default">
        <AvatarFallback className="bg-sidebar-primary text-xs font-medium leading-4 text-sidebar-primary-foreground">
          {user.initials}
        </AvatarFallback>
      </Avatar>
      <div className="absolute left-14 top-2.5 flex h-9 w-[163px] flex-col items-start">
        <div
          className={
            user.nameWeight === "bold"
              ? "h-5 w-full truncate text-sm font-bold leading-5 text-sidebar-foreground"
              : "h-5 w-full truncate text-sm font-semibold leading-5 text-sidebar-foreground"
          }
        >
          {user.name}
        </div>
        <div className="h-4 w-full truncate text-xs font-normal leading-4 text-sidebar-foreground/70">
          {user.detail}
        </div>
      </div>
    </div>
  );
}

function DashboardTopBar({ placeholder }: { placeholder: string }) {
  return (
    <header className="flex h-16 shrink-0 items-center justify-between border-b border-border px-6 py-3.5">
      <div className="relative h-9 w-[448px]">
        <Search className="-translate-y-1/2 absolute left-3 top-1/2 size-4 text-muted-foreground" />
        <Input
          aria-label="Поиск"
          className="h-9 border-transparent bg-muted/50 pl-9 text-sm shadow-none placeholder:text-muted-foreground focus-visible:ring-0"
          placeholder={placeholder}
        />
      </div>
      <Button
        aria-label="Уведомления"
        className="relative size-9 rounded-lg text-foreground"
        size="icon"
        variant="ghost"
      >
        <Bell className="size-4" />
        <span className="absolute left-[22px] top-1.5 size-2 rounded-full bg-destructive" />
      </Button>
    </header>
  );
}
