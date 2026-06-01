import React from "react";
import { Navigate, Outlet } from "react-router";
import { PiImages as ImagesIcon } from "react-icons/pi";
import { DashboardLayout } from "~src/common/layout/dashboardLayout";
import { useTranslation } from "react-i18next";
import { useAuth } from "~src/lib/auth";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";

export function ProfileLayout() {
  const { t } = useTranslation();
  const { data: authData, isLoading } = useAuth();

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

  if (isLoading) return <FullScreenLoader />;
  if (!authData?.viewer.organizationUser) return <Navigate to="/" replace />;

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
