import { t } from "i18next";
import React from "react";
import classNames from "~node_modules/classnames";
import { ImageGallery } from "~src/common/ImageGallery/render";
import { HEADING_2 } from "~src/ui/classNames";

export function ListImages() {
  return (
    <div className="m-auto max-w-full mx-8">
      <h1 className={classNames(HEADING_2, "font-poppins")}>
        {t("siteadmin.images.list.title")}
      </h1>
      <div className="mt-6">
        <ImageGallery />
      </div>
    </div>
  );
}
