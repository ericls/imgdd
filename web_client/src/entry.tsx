import React from "react";
import { createRoot } from "react-dom/client";
import { ToastContainer } from "react-toastify";
import { ApolloProvider } from "@apollo/client";
import "./entry.css";
import "../node_modules/react-toastify/dist/ReactToastify.css";
import { DarkModeContext, DarkModeProvider } from "./lib/darkMode";

import "./i18n";
import { apolloClient } from "./apollo";
import { AuthProvider } from "./lib/auth";
import { Layout } from "./layout";

function AutoDarkToastContainer() {
  const darkContext = React.useContext(DarkModeContext);
  return (
    <ToastContainer
      position="bottom-center"
      theme={darkContext.isDarkMode ? "dark" : "light"}
    />
  );
}

function Root() {
  return (
    <ApolloProvider client={apolloClient}>
      <AuthProvider>
        <DarkModeProvider>
          <Layout />
          <AutoDarkToastContainer />
        </DarkModeProvider>
      </AuthProvider>
    </ApolloProvider>
  );
}

const appRoot = document.getElementById("app");
if (appRoot) {
  const domRoot = createRoot(appRoot);
  domRoot.render(<Root />);
} else {
  console.error("Failed to load dom");
}
