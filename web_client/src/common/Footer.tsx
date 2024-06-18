import classNames from "classnames";
import React from "react";
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
      </div>
    </div>
  );
}
