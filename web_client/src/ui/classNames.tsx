import cx from "classnames";

export const LINK_COLOR =
  "text-indigo-600 hover:text-indigo-500 dark:text-indigo-500 dark:hover:text-indigo-400";
export const TEXT_COLOR = "text-neutral-600 dark:text-neutral-100";
export const SECONDARY_TEXT_COLOR = "text-neutral-700 dark:text-neutral-100";
export const SECONDARY_TEXT_COLOR_DIM =
  "text-neutral-500 dark:text-neutral-400";
export const SECONDARY_TEXT_COLOR_DIMMER =
  "text-neutral-400 dark:text-neutral-500";
export const SECONDARY_TEXT_COLOR_HILIGHT =
  "text-neutral-900 dark:text-neutral-100";
export const HEADING_1 = cx("text-3xl font-bold tracking-tight", TEXT_COLOR);
export const HEADING_2 = cx("text-2xl font-bold tracking-tight", TEXT_COLOR);
export const BASE_LAYER = "bg-neutral-100 dark:bg-neutral-900";
export const BASE_LAYER_HOVER =
  "hover:bg-neutral-200 hover:dark:bg-neutral-700";
export const SECOND_LAYER = "bg-white dark:bg-neutral-800";
export const SECOND_LAYER_HOVER =
  "hover:bg-neutral-200 hover:dark:bg-neutral-700";

export const PRIMARY_TEXT_COLOR = "text-primary-600 dark:text-primary-400";
export const PRIMARY_BORDER_COLOR =
  "border-primary-600 dark:border-primary-500";
export const PRIMARY_BORDER_COLOR_ON_HOVER =
  "hover:border-primary-600 dark:hover:border-primary-500";

export const LOGO_TEXT_1 = cx("font-poppins font-bold select-none", TEXT_COLOR);
export const LOGO_TEXT_2 = "text-indigo-600 dark:text-indigo-500";

export const DEFAULT_INPUT_LABEL = "block mb-1";
export const DEFAULT_INPUT =
  "appearance-none rounded-md border border-gray-300 px-3 py-2 " +
  "text-gray-900 placeholder-gray-500 " +
  "focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 " +
  "dark:bg-neutral-800 dark:border-neutral-900 dark:placeholder-neutral-400 dark:text-neutral-100 " +
  "dark:focus:border-indigo-500 " +
  "disabled:opacity-50 disabled:cursor-not-allowed ";
export const DEFAULT_INPUT_CHECKBOX =
  "accent-pink-700 rounded-md px-2 py-2 " +
  "disabled:opacity-50 disabled:cursor-not-allowed ";

// list of classes that we want tailwind to include
// eslint-disable-next-line @typescript-eslint/no-unused-expressions
["rounded-full", "bg-primary-600"];
