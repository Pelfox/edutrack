import { createRootRoute, Outlet } from "@tanstack/react-router";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <div className="w-screen h-screen bg-background text-foreground font-sans">
      <Outlet />
    </div>
  );
}
