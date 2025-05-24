import React from "react";
import { Route, Routes } from "react-router-dom";
import { UsersList } from "./usersList";
import { UserImageGallery } from "./UserImageGallery";

export function Users() {
  return (
    <Routes>
      <Route path="/" element={<UsersList />} />
      <Route path=":userId/images" element={<UserImageGallery />} />
    </Routes>
  );
}
