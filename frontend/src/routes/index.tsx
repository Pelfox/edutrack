import { createFileRoute } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="p-8">
      <p>Hello, world!</p>
      <Button type="button">Hello!</Button>
    </div>
  );
}
