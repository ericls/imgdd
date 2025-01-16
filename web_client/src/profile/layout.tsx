import React from "react";
import { Outlet } from "react-router-dom";
import { PiImages as ImagesIcon } from "react-icons/pi";
import { DashboardLayout } from "~src/common/layout/dashboardLayout";
import { useTranslation } from "react-i18next";

export function ProfileLayout() {
  const { t } = useTranslation();
  const sideBarMenuGroups = React.useMemo(
    () => [
      {
        title: t("profile.nav.general"),
        items: [
          {
            title: t("profile.nav.images"),
            to: "/profile/images",
            icon: <ImagesIcon />,
            active: (location: { pathname: string }) =>
              location.pathname.startsWith("/profile/images"),
          },
        ],
      },
    ],
    [t],
  );

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
