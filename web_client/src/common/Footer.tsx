import classNames from "classnames";
import React from "react";
import { Slot } from "~src/lib/slot";
import { SECONDARY_TEXT_COLOR_DIM } from "~src/ui/classNames";

const DefaultFooterContent = (
  <ul className="flex space-x-1">
    <li>
      <a href="https://hub.docker.com/r/ericls/imgdd">Docker</a>
    </li>
    <li className="hidden md:inline select-none">|</li>
    <li>
      <a href="https://github.com/ericls/imgdd">Github</a>
    </li>
    <li className="hidden md:inline select-none">|</li>
    <li>
      <a href="http://imgdd.com">IMGDD.COM 🇨🇦</a>
    </li>
  </ul>
);

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
          <Slot id="footer-content" fallback={DefaultFooterContent} />
          {/* IMGDD.COM 🇨🇦 */}
        </span>
      </div>
    </div>
  );
}
