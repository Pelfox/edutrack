/// <reference types="vite/client" />

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";
import type { AuthContextValue } from "@/lib/context/auth";
import { AuthProvider, useAuth } from "@/lib/context/auth";
import { routeTree } from "./routeTree.gen";
import "./globals.css";

const authPlaceholder = null as unknown as AuthContextValue;

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 30_000,
    },
  },
});

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
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <RouterWithContext />
      </AuthProvider>
    </QueryClientProvider>,
  );
}

// biome-ignore lint/style/useComponentExportOnlyModules: Entry point keeps router initialization colocated.
function RouterWithContext() {
  const auth = useAuth();

  return <RouterProvider context={{ auth }} router={router} />;
}
