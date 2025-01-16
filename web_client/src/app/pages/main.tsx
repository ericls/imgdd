import classnames from "classnames";
import React from "react";
import { SECONDARY_TEXT_COLOR_DIM, TEXT_COLOR } from "~src/ui/classNames";
import { Uplodaer } from "~src/uploader/uploader";
import { useTranslation } from "react-i18next";
import { useAuth } from "~src/lib/auth";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";
import { TextLogo } from "~src/common/TextLogo";

export function AppMainPage() {
  const { t } = useTranslation();
  const { isLoading: isAuthLoading } = useAuth();
  if (isAuthLoading) {
    return (
      <div className="absolute w-full h-full">
        <FullScreenLoader withLogo />
      </div>
    );
  }
  return (
    <div className="max-w-screen-sm flex flex-col grow mx-auto items-center">
      <h1
        className={classnames(
          TEXT_COLOR,
          "font-poppins mb-4 mt-20 font-bold select-none",
        )}
      ></h1>
      <TextLogo className="text-4xl" />
      <p
        className={classnames(
          SECONDARY_TEXT_COLOR_DIM,
          "font-poppins text-2xl mb-10 text-center",
        )}
      >
        {t("home.tagLine")}
      </p>
      <Uplodaer />
    </div>
  );
}
