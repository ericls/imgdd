import classnames from "classnames";
import classNames from "classnames";
import React from "react";
import { Link, Outlet, useLocation } from "react-router-dom";
import {
  HiOutlineUsers as UsersIcon,
  HiOutlineServerStack as SotrageIcon,
} from "react-icons/hi2";
import { PiImages as ImagesIcon } from "react-icons/pi";
// import { TbDeviceAnalytics as AnalyticsIcon } from "react-icons/tb";
import { Footer } from "~src/common/Footer";
import { LazyRouteFallback } from "~src/common/LazyRouteFallback";
import { TopNav } from "~src/common/TopNav";
import {
  PRIMARY_BORDER_COLOR,
  PRIMARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIMMER,
} from "~src/ui/classNames";
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
