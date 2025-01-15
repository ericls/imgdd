import React from "react";
import { Outlet } from "react-router-dom";
import {
  HiOutlineUsers as UsersIcon,
  HiOutlineServerStack as SotrageIcon,
} from "react-icons/hi2";
import { PiImages as ImagesIcon } from "react-icons/pi";
import {
  DashboardLayout,
  DashboardLayoutMenuGroups,
} from "~src/common/layout/dashboardLayout";

const sideBarMenuGroups: DashboardLayoutMenuGroups = [
  {
    title: "General",
    items: [
      {
        title: "Users",
        to: "/site-admin/users",
        icon: <UsersIcon />,
        active: (location) => location.pathname.startsWith("/site-admin/users"),
      },
      {
        title: "Images",
        to: "/site-admin/images",
        icon: <ImagesIcon />,
        active: (location) =>
          location.pathname.startsWith("/site-admin/images"),
      },
    ],
  },
  {
    title: "Storage",
    items: [
      {
        title: "Storage Backend",
        to: "/site-admin/storage",
        active: (location) =>
          location.pathname.startsWith("/site-admin/storage"),
        icon: <SotrageIcon />,
      },
    ],
  },
];

export function SiteAdminLayout() {
  return (
    <DashboardLayout
      menuGroups={sideBarMenuGroups}
      mainAreaClassName="site-admin-main"
    >
      <Outlet />
    </DashboardLayout>
  );
}

export default SiteAdminLayout;
