/// <reference types="vite/client" />

import { createRouter, RouterProvider } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";
import type { AuthContextValue } from "@/lib/context/auth";
import { AuthProvider, useAuth } from "@/lib/context/auth";
import { routeTree } from "./routeTree.gen";
import "./globals.css";

const authPlaceholder = null as unknown as AuthContextValue;

const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
  context: {
    auth: authPlaceholder,
  },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

// biome-ignore lint/style/noNonNullAssertion: Application will automatically fail without root element.
const rootElement = document.getElementById("root")!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <AuthProvider>
      <RouterWithContext />
    </AuthProvider>,
  );
}

// biome-ignore lint/style/useComponentExportOnlyModules: Entry point keeps router initialization colocated.
function RouterWithContext() {
  const auth = useAuth();

  return <RouterProvider context={{ auth }} router={router} />;
}
