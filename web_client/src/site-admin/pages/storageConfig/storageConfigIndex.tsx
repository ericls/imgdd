import React from "react";
import { ListStorageDef } from "./listStorageDef";
import { Navigate, Outlet, Route, Routes } from "react-router-dom";
import { UpdateOrCreateStorageDef } from "./updateOrCreateStorageDef";

export function StorageConfig() {
  return (
    <>
      <Routes>
        <Route
          path="/"
          element={
            <Navigate
              to={"/site-admin/storage/storage-def/list"}
              relative="route"
              replace
            />
          }
        />
        <Route path="/storage-def/list" element={<ListStorageDef />} />
        <Route path="/storage-def/:id" element={<UpdateOrCreateStorageDef />} />
      </Routes>
      <Outlet />
    </>
  );
}
