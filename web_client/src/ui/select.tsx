import React, { useId } from "react";
import cx from "classnames";
import { DEFAULT_INPUT, DEFAULT_INPUT_LABEL } from "./classNames";

export const Select = React.forwardRef<
  HTMLSelectElement,
  JSX.IntrinsicElements["select"]
>(({ className, ...props }, ref) => {
  return (
    <select ref={ref} className={cx(className, DEFAULT_INPUT)} {...props} />
  );
});

export const SelectWithLabel = React.forwardRef<
  HTMLSelectElement,
  React.ComponentProps<typeof Select> & {
    label: string;
    containerClassName?: string;
    labelClassName?: string;
  }
>(({ label, containerClassName, labelClassName, ...props }, ref) => {
  const id = useId();
  return (
    <div className={containerClassName}>
      <label className={labelClassName || DEFAULT_INPUT_LABEL} htmlFor={id}>
        {label}
      </label>
      <Select id={id} {...props} ref={ref} />
    </div>
  );
});
