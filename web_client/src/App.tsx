import classnames from "classnames";
import React from "react";
import {
  SECONDARY_TEXT_COLOR_DIM,
  TEXT_COLOR,
} from "./ui/classNames";
import { Uplodaer } from "./uploader/uploader";
import { useTranslation } from "react-i18next";
import { useAuth } from "./lib/auth";
import { FullScreenLoader } from "./ui/fullscreenLoader";
import { TextLogo } from "./common/TextLogo";

export function App() {
  const { t } = useTranslation();
  const { isLoading: isAuthLoading } = useAuth();
  if (isAuthLoading) {
    return <div className="grow flex items-center"><FullScreenLoader /></div>
  }
  return (
    <div className="max-w-screen-sm flex flex-col grow mx-auto items-center">
      <h1
        className={classnames(
          TEXT_COLOR,
          "font-poppins mb-4 mt-20 font-bold select-none"
        )}
      >
      </h1>
        <TextLogo className="text-4xl"/>
      <p
        className={classnames(
          SECONDARY_TEXT_COLOR_DIM,
          "font-poppins text-2xl mb-10 text-center"
        )}
      >
        {t("home.tagLine")}
      </p>
      <Uplodaer />
    </div>
  );
}
