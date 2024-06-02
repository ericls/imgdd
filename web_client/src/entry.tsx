import React from "react";
import { createRoot } from "react-dom/client";
import { ToastContainer } from "react-toastify";
import { ApolloProvider } from "@apollo/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import "./entry.css";
import "../node_modules/react-toastify/dist/ReactToastify.css";
import { DarkModeContext, DarkModeProvider } from "./lib/darkMode";

import "./i18n";
import { apolloClient } from "./apollo";
import { AuthProvider } from "./lib/auth";
import { AdminLayout } from "./admin/layout";
import { AppLayout } from "./app/layout";
import { AppMainPage } from "./app/pages/main";

function AutoDarkToastContainer() {
  const darkContext = React.useContext(DarkModeContext);
  return (
    <ToastContainer
      position="bottom-center"
      theme={darkContext.isDarkMode ? "dark" : "light"}
    />
  );
}

function ErrorPage() {
  return <div>Error</div>;
}

const router = createBrowserRouter([
  {
    path: "admin",
    element: <AdminLayout />,
    errorElement: <ErrorPage />,
  },
  {
    path: "/",
    element: <AppLayout />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "/",
        index: true,
        element: <AppMainPage />,
      },
      {
        path: "auth",
        lazy: async () => {
          const { AuthPage } = await import("./app/pages/auth");
          return {
            element: <AuthPage />,
          };
        },
      },
    ],
  },
]);

function Root() {
  return (
    <React.StrictMode>
      <ApolloProvider client={apolloClient}>
        <AuthProvider>
          <DarkModeProvider>
            <RouterProvider router={router} />
            <AutoDarkToastContainer />
          </DarkModeProvider>
        </AuthProvider>
      </ApolloProvider>
    </React.StrictMode>
  );
}

const appRoot = document.getElementById("app");
if (appRoot) {
  const domRoot = createRoot(appRoot);
  domRoot.render(<Root />);
} else {
  console.error("Failed to load dom");
}
