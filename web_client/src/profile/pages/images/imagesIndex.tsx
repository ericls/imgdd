import React from "react";
import { Navigate, Outlet, Route, Routes } from "react-router-dom";
import { ListImages } from "./listImages";

export function Images() {
  return (
    <>
      <Routes>
        <Route
          path="/"
          element={
            <Navigate to={"/profile/images/list"} relative="route" replace />
          }
        />
        <Route path="/list" element={<ListImages />} />
      </Routes>
      <Outlet />
    </>
  );
}
