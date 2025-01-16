import React from "react";
import cx from "classnames";
import { DarkModeContext, DarkModeTheme } from "~src/lib/darkMode";
import { Button } from "./button";

import {
  HiOutlineMoon as DarkModeIcon,
  HiOutlineSun as LightModeIcon,
} from "react-icons/hi2";
import { SECONDARY_TEXT_COLOR_DIM } from "./classNames";
import { MenuWithTrigger } from "./menuWithTrigger";
import { noop } from "lodash-es";

export function DarkModeSettings() {
  const { theme, setTheme, isDarkMode } = React.useContext(DarkModeContext);
  const [checkbox, setCheckbox] = React.useState<HTMLInputElement | null>(null);
  const onToggleUseSystem = React.useCallback(() => {
    let newValue: DarkModeTheme = "system";
    if (theme === "system") {
      newValue = isDarkMode ? "light" : "dark";
    }
    setTheme(newValue);
  }, [isDarkMode, setTheme, theme]);
  const onToggleDarkMode = React.useCallback(() => {
    setTheme(isDarkMode ? "light" : "dark");
  }, [isDarkMode, setTheme]);
  React.useEffect(() => {
    const current = checkbox;
    if (!current) return;
    if (theme === "system") {
      current.indeterminate = true;
    } else {
      current.indeterminate = false;
    }
  }, [theme, checkbox]);
  return (
    <MenuWithTrigger
      placement="top-start"
      containerClassName="block text-center"
      trigger={
        <Button
          variant="transparent"
          className={cx(SECONDARY_TEXT_COLOR_DIM, "block text-center")}
        >
          {isDarkMode ? (
            <LightModeIcon size={24} className="w-4" />
          ) : (
            <DarkModeIcon size={24} className="w-4" />
          )}
        </Button>
      }
      menuSections={{
        children: [
          {
            id: "main",
            items: [
              {
                id: "darkModelToggle",
                children: (
                  <div className="flex items-center">
                    <input
                      ref={setCheckbox}
                      checked={isDarkMode}
                      disabled={theme === "system"}
                      type="checkbox"
                      className={cx(
                        "w-4 h-4 mr-1 text-blue-600 bg-gray-100 rounded border-gray-300 focus:ring-indigo-500 dark:focus:ring-indigo-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600",
                        {
                          "opacity-50 focus:ring-0 ring-0 bg-gray-300":
                            theme === "system",
                        },
                      )}
                      onChange={noop}
                    />
                    Dark
                  </div>
                ),
                action: onToggleDarkMode,
                disabled: theme === "system",
              },
              {
                id: "useSystem",
                children: (
                  <div className="flex items-center">
                    <input
                      checked={theme === "system"}
                      type="checkbox"
                      className="w-4 h-4 mr-1 text-blue-600 bg-gray-100 rounded border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                      onChange={noop}
                    />
                    Auto
                  </div>
                ),
                action: onToggleUseSystem,
              },
            ],
          },
        ],
      }}
    />
  );
}
