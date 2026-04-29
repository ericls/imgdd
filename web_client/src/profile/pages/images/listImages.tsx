import { t } from "i18next";
import React from "react";
import classNames from "classnames";
import { DEFAULT_MENU_CONFIG } from "~src/common/ImageGallery/menu";
import { ImageGallery } from "~src/common/ImageGallery/render";
import { useAuth } from "~src/lib/auth";
import { useDebounce } from "~src/lib/hooks";
import { HEADING_2 } from "~src/ui/classNames";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";
import { Input } from "~src/ui/input";

export function ListImages() {
  const { data: authData, isLoading: isAuthLoading } = useAuth();
  const user = React.useMemo(
    () => authData?.viewer.organizationUser,
    [authData],
  );
  const [searchTerm, setSearchTerm] = React.useState("");
  const debouncedSearchTerm = useDebounce(searchTerm, 500);
  if (isAuthLoading || !user) {
    return <FullScreenLoader />;
  }
  return (
    <div className="m-auto max-w-full mx-8">
      <h1 className={classNames(HEADING_2, "font-poppins")}>
        {t("profile.images.list.title")}
      </h1>
      <div className="my-4 w-full max-w-md mt-2">
        <Input
          type="text"
          placeholder={t("profile.images.list.searchPlaceholder")}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>
      <div className="mt-6">
        <ImageGallery
          menuConfig={DEFAULT_MENU_CONFIG}
          createdById={user?.id}
          nameContains={debouncedSearchTerm || undefined}
        />
      </div>
    </div>
  );
}
