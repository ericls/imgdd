import React from "react";
import { Outlet } from "react-router-dom";
import { Footer } from "~src/common/Footer";
import { LazyRouteFallback } from "~src/common/LazyRouteFallback";
import { TopNav } from "~src/common/TopNav";

export function AdminLayout() {
  return (
    <div className="main min-h-full flex flex-col mx-2">
      <TopNav />
      <div className="grow relative z-0">
        <LazyRouteFallback />
        <Outlet />
      </div>
      <div>
        <Footer />
      </div>
    </div>
  );
}
