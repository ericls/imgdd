import classNames from "classnames";
import React from "react";
import { isChinaTimezone, isChinese } from "~src/lib/locale";
import { SECONDARY_TEXT_COLOR_DIM } from "~src/ui/classNames";

export function Footer({ center = false }: { center?: boolean }) {
  return (
    <div
      className={classNames(
        "p-2 text-sm flex justify-end",
        SECONDARY_TEXT_COLOR_DIM,
        {
          ["text-end"]: !center,
          ["text-center"]: center,
        }
      )}
    >
      <div>
        <span className="opacity-50">IMGDD.COM 🇨🇦</span>
        <br />
        <div className="flex gap-x-1">
          {isChinese() && isChinaTimezone() ? (
            <a
              className="px-2"
              href="https://www.npmjs.com/package/picgo-plugin-imgdd"
              target="_blank noopener noreferrer nofollow external"
            >
              PicGo 插件
            </a>
          ) : null}
          <a className="px-2" href={`mailto:${window.SUPPORT_EMAIL}`}>
            Report abuse
          </a>
        </div>
      </div>
    </div>
  );
}
