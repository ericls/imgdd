import classNames from "classnames";
import React from "react";
import { Slot } from "~src/lib/slot";
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
        },
      )}
    >
      <div>
        <span slot-id="footer-content" className="opacity-50">
          <Slot id="footer-content" fallback="IMGDD.COM 🇨🇦" />
          {/* IMGDD.COM 🇨🇦 */}
        </span>
      </div>
    </div>
  );
}
