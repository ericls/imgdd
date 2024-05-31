import classnames from "classnames";
import React from "react";
import { Footer } from "./common/Footer";
import {
  LINK_COLOR,
  SECONDARY_TEXT_COLOR_DIM,
  TEXT_COLOR,
} from "./ui/classNames";
import { Uplodaer } from "./uploader/uploader";
import { TopNav } from "./common/TopNav";
import { useTranslation } from "react-i18next";

export function App() {
  const { t } = useTranslation();
  return (
    <div className="main h-full flex flex-col mx-2">
      <TopNav />
      <div className="max-w-screen-sm flex flex-col grow mx-auto items-center">
        <h1
          className={classnames(
            TEXT_COLOR,
            "font-poppins text-4xl mb-4 mt-20 font-bold select-none"
          )}
        >
          IMG<span className={LINK_COLOR}>DD</span>
        </h1>
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
      <div>
        <Footer />
      </div>
    </div>
  );
}
