import { createRootRouteWithContext, Outlet, redirect } from "@tanstack/react-router";
import { Toaster } from "sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import type { AuthContextValue } from "@/lib/context/auth";

export const Route = createRootRouteWithContext<{
  auth: AuthContextValue;
}>()({
  beforeLoad: ({ context, location }) => {
    if (isDashboardPath(location.pathname) && !context.auth.isAuthenticated) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
        replace: true,
      });
    }
  },
  component: RootComponent,
});

function RootComponent() {
  return (
    <div className="w-screen h-screen bg-background text-foreground font-sans">
      <TooltipProvider>
        <Outlet />
      </TooltipProvider>
      <Toaster closeButton={true} position="top-right" richColors={true} />
    </div>
  );
}

function isDashboardPath(pathname: string) {
  return pathname !== "/login";
}
