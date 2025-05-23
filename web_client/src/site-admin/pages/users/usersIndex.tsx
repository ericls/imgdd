import React from "react";
import { Route, Routes } from "react-router-dom";
import { UsersList } from "./usersList";

export function Users() {
  return (
    <Routes>
      <Route path="/" element={<UsersList />} />
    </Routes>
  );
}
