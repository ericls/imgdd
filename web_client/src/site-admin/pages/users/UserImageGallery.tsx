import React from "react";
import { useParams } from "react-router-dom";
import { DEFAULT_MENU_CONFIG } from "~src/common/ImageGallery/menu";
import { ImageGallery } from "~src/common/ImageGallery/render";
import { useTranslation } from "react-i18next";

export function UserImageGallery() {
  const { userId } = useParams();
  const { t } = useTranslation();

  return (
    <>
      <h2 className="text-2xl font-bold mb-4">{t("userImageGallery.title")}</h2>
      <div className="p-4">
        <ImageGallery createdById={userId} menuConfig={DEFAULT_MENU_CONFIG} />
      </div>
    </>
  );
}
