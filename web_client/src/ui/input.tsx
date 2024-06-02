import React from "react";
import cx from "classnames";

export const Input = React.forwardRef<
  HTMLInputElement,
  JSX.IntrinsicElements["input"]
>(({ className, ...props }, ref) => {
  return (
    <input
      ref={ref}
      className={cx(
        className,
        `appearance-none rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 dark:bg-neutral-800 dark:border-neutral-900 dark:placeholder-neutral-400 dark:text-neutral-100`
      )}
      {...props}
    />
  );
});
