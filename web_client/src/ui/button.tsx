import React from "react";
import cx from "classnames";

const BUTTON_VARIANTS = {
  indigo:
    "bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-700 dark:hover:bg-indigo-600 focus:ring-indigo-500 text-white",
  transparent:
    "bg-transparent shadow-transparent focus:ring-0 focus:ring-offset-0 outline-none focus:ring-neutral-300 dark:focus:ring-neutral-500",
  secondary: `bg-neutral-50 dark:bg-neutral-700 hover:bg-neutral-100 hover:dark:bg-neutral-600`,
  green: `bg-green-600 hover:bg-green-700 dark:bg-green-700 dark:hover:bg-green-600 focus:ring-green-600 text-white`,
  destructive:
    "bg-red-600 hover:bg-red-700 dark:bg-red-700 dark:hover:bg-red-600 focus:ring-red-500 text-white",
} as const;

export const Button = React.forwardRef<
  HTMLButtonElement,
  JSX.IntrinsicElements["button"] & {
    variant?: keyof typeof BUTTON_VARIANTS;
    roundLevel?: "md" | "full" | "lg" | "";
    noPadding?: boolean;
  }
>(
  (
    {
      className,
      variant = "indigo",
      roundLevel = "md",
      noPadding = false,
      ...props
    },
    ref,
  ) => {
    const variantClassName = BUTTON_VARIANTS[variant];
    let roundedClassName = "";
    if (roundLevel) {
      roundedClassName = `rounded-${roundLevel}`;
    } else if (roundLevel === "") {
      roundedClassName = "rounded";
    }
    return (
      <button
        ref={ref}
        className={cx(
          className,
          variantClassName,
          roundedClassName,
          `justify-center border border-transparent shadow`,
          `text-sm font-medium`,
          `focus:outline-none focus:ring-2`,
          { "py-2 px-4": !noPadding },
        )}
        {...props}
      />
    );
  },
);
