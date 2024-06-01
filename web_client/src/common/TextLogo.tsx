import React from "react";
import classNames from "classnames";
import { LOGO_TEXT_1, LOGO_TEXT_2 } from "~src/ui/classNames";

export function TextLogoSmall({ className }: { className?: string }) {
  return (
    <h1 className={classNames(LOGO_TEXT_1, className)}>
      I<span className={LOGO_TEXT_2}>D</span>
    </h1>
  );
}

export function TextLogo({ className }: { className?: string }) {
  return (
    <div className={classNames(LOGO_TEXT_1, className)}>
      IMG<span className={LOGO_TEXT_2}>DD</span>
    </div>
  );
}
