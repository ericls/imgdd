import React from "react";
import cx from "classnames";
import { useTranslation } from "react-i18next";
import { HiOutlineLanguage as LanguageSwitchIcon } from "react-icons/hi2";
import { noop } from "lodash-es";

import { Button } from "./button";
import { SECONDARY_TEXT_COLOR_DIM } from "./classNames";
import { MenuWithTrigger } from "./menuWithTrigger";
import { loadLanguage } from "~src/i18n";
import { setStoredLanguage, SupportedLanguage } from "~src/lib/locale";

const LANGUAGES: { id: SupportedLanguage; label: string }[] = [
  { id: "en", label: "English" },
  { id: "zh_hans", label: "简体中文" },
];

export function LanguageSettings() {
  const { i18n } = useTranslation();
  const current = (i18n.resolvedLanguage ?? i18n.language) as SupportedLanguage;
  const onPick = async (lang: SupportedLanguage) => {
    setStoredLanguage(lang);
    await loadLanguage(lang);
    await i18n.changeLanguage(lang);
  };
  return (
    <MenuWithTrigger
      placement="top-start"
      containerClassName="block text-center"
      trigger={
        <Button
          variant="transparent"
          className={cx(SECONDARY_TEXT_COLOR_DIM, "block text-center")}
        >
          <LanguageSwitchIcon size={24} className="w-4" />
        </Button>
      }
      menuSections={{
        children: [
          {
            id: "main",
            items: LANGUAGES.map((lang) => ({
              id: lang.id,
              children: (
                <div className="flex items-center">
                  <input
                    checked={current === lang.id}
                    type="radio"
                    className="w-4 h-4 mr-1 text-blue-600 bg-gray-100 rounded border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                    onChange={noop}
                  />
                  {lang.label}
                </div>
              ),
              action: () => onPick(lang.id),
            })),
          },
        ],
      }}
    />
  );
}
