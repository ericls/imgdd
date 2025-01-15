import React from "react";
import { Outlet } from "react-router-dom";
import { PiImages as ImagesIcon } from "react-icons/pi";
import { DashboardLayout } from "~src/common/layout/dashboardLayout";

const sideBarMenuGroups = [
  {
    title: "General",
    items: [
      {
        title: "images",
        to: "/profile/images",
        icon: <ImagesIcon />,
        active: (location: { pathname: string }) =>
          location.pathname.startsWith("/profile/images"),
      },
    ],
  },
];

export function ProfileLayout() {
  return (
    <DashboardLayout
      menuGroups={sideBarMenuGroups}
      mainAreaClassName="profile-main"
    >
      <Outlet />
    </DashboardLayout>
  );
}

export default ProfileLayout;
