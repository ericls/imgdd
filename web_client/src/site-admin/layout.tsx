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
import { useTranslation } from "react-i18next";

export function SiteAdminLayout() {
  const { t } = useTranslation();
  const sideBarMenuGroups: DashboardLayoutMenuGroups = React.useMemo(
    () => [
      {
        title: t("siteadmin.nav.general"),
        items: [
          {
            title: t("siteadmin.nav.users"),
            to: "/site-admin/users",
            icon: <UsersIcon />,
            active: (location) =>
              location.pathname.startsWith("/site-admin/users"),
          },
          {
            title: t("siteadmin.nav.images"),
            to: "/site-admin/images",
            icon: <ImagesIcon />,
            active: (location) =>
              location.pathname.startsWith("/site-admin/images"),
          },
        ],
      },
      {
        title: t("siteadmin.nav.storage"),
        items: [
          {
            title: t("siteadmin.nav.storageBackend"),
            to: "/site-admin/storage",
            active: (location) =>
              location.pathname.startsWith("/site-admin/storage"),
            icon: <SotrageIcon />,
          },
        ],
      },
    ],
    [t],
  );
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
