import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TooltipProvider } from "@/components/ui/tooltip";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <div className="w-screen h-screen bg-background text-foreground font-sans">
      <TooltipProvider>
        <Outlet />
      </TooltipProvider>
    </div>
  );
}
