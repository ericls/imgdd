import React from "react";
import { Route, Routes, useParams } from "react-router-dom";
import { UsersList } from "./usersList";
import { ImageGallery } from "~src/common/ImageGallery/render";

export function Users() {
  return (
    <Routes>
      {/* Main user list table */}
      <Route path="/" element={<UsersList />} />

      {/* Route for user-specific images */}
      <Route path=":userId/images" element={<ImageGalleryWrapper />} />
    </Routes>
  );
}

// Wrapper to map `userId` from params to `createdById` prop for ImageGallery
function ImageGalleryWrapper() {
  const { userId } = useParams();
  return <ImageGallery createdById={userId} />;
}
