import { t } from "i18next";
import React from "react";
import classNames from "~node_modules/classnames";
import { DEFAULT_MENU_CONFIG } from "~src/common/ImageGallery/menu";
import { ImageGallery } from "~src/common/ImageGallery/render";
import { useAuth } from "~src/lib/auth";
import { HEADING_2 } from "~src/ui/classNames";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";

export function ListImages() {
  const { data: authData, isLoading: isAuthLoading } = useAuth();
  const user = React.useMemo(
    () => authData?.viewer.organizationUser,
    [authData],
  );
  if (isAuthLoading || !user) {
    return <FullScreenLoader />;
  }
  return (
    <div className="m-auto max-w-full mx-8">
      <h1 className={classNames(HEADING_2, "font-poppins")}>
        {t("profile.images.list.title")}
      </h1>
      <div className="mt-6">
        <ImageGallery menuConfig={DEFAULT_MENU_CONFIG} createdById={user?.id} />
      </div>
    </div>
  );
}
