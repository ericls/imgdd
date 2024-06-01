import React from "react";
import { TextLogoSmall } from "./TextLogo";
import { DarkModeSettings } from "~src/ui/darkModeToggle";

export function TopNav() {
  return (
    <div className="top-nav sticky top-0 backdrop-blur-sm p-2 flex justify-between text-neutral-800 dark:text-neutral-200 text-end">
      <div> {/* left */}
        <TextLogoSmall className="text-2xl" />
      </div>
      <div> {/* right */}
        <DarkModeSettings />
      </div>
    </div>
  );
}
