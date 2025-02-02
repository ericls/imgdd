import React from "react";
import classNames from "classnames";
import { LOGO_TEXT_1, LOGO_TEXT_2 } from "~src/ui/classNames";

export function TextLogoSmall({ className }: { className?: string }) {
  return (
    <h1
      className={classNames(LOGO_TEXT_1, className)}
      slot-id="text-logo-small"
    >
      {(() => {
        if (window.SITE_NAME) {
          return window.SITE_NAME;
        } else {
          return (
            <>
              I<span className={LOGO_TEXT_2}>D</span>
            </>
          );
        }
      })()}
    </h1>
  );
}

export function TextLogo({ className }: { className?: string }) {
  return (
    <div className={classNames(LOGO_TEXT_1, className)} slot-id="text-logo">
      {(() => {
        if (window.SITE_NAME) {
          return window.SITE_NAME;
        } else {
          return (
            <>
              IMG<span className={LOGO_TEXT_2}>DD</span>
            </>
          );
        }
      })()}
    </div>
  );
}
