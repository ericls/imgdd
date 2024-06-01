import React from "react";
import { Footer } from "./common/Footer";
import { TopNav } from "./common/TopNav";
import { useTranslation } from "react-i18next";
import { App } from "./App";

export function Layout() {
  const { t } = useTranslation();
  return (
    <div className="main min-h-full flex flex-col mx-2">
      <TopNav />
      <App />
      <div>
        <Footer />
      </div>
    </div>
  );
}
