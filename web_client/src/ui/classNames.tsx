import cx from "classnames";

export const LINK_COLOR =
  "text-indigo-600 hover:text-indigo-500 dark:text-indigo-500 dark:hover:text-indigo-400";
export const TEXT_COLOR = "text-neutral-600 dark:text-neutral-100";
export const SECONDARY_TEXT_COLOR = "text-neutral-700 dark:text-neutral-100";
export const SECONDARY_TEXT_COLOR_DIM =
  "text-neutral-500 dark:text-neutral-400";
export const SECONDARY_TEXT_COLOR_HILIGHT =
  "text-neutral-900 dark:text-neutral-100";
export const HEADING_1 = cx("text-3xl font-bold tracking-tight", TEXT_COLOR);
export const BASE_LAYER = "bg-neutral-100 dark:bg-neutral-900";
export const BASE_LAYER_HOVER =
  "hover:bg-neutral-200 hover:dark:bg-neutral-700";
export const SECOND_LAYER = "bg-neutral-50 dark:bg-neutral-800";
export const SECOND_LAYER_HOVER =
  "hover:bg-neutral-200 hover:dark:bg-neutral-700";

export const LOGO_TEXT_1 = cx("font-poppins font-bold select-none", TEXT_COLOR);
export const LOGO_TEXT_2 = "text-indigo-600 dark:text-indigo-500";

// list of classes that we want tailwind to include
["rounded-full", "bg-primary-600"];
