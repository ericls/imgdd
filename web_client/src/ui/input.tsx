import React, { useId } from "react";
import cx from "classnames";
import { DEFAULT_INPUT, DEFAULT_INPUT_LABEL } from "./classNames";

export const Input = React.forwardRef<
  HTMLInputElement,
  JSX.IntrinsicElements["input"]
>(({ className, ...props }, ref) => {
  return (
    <input ref={ref} className={cx(className, DEFAULT_INPUT)} {...props} />
  );
});

export const InputWithLabel = React.forwardRef<
  HTMLInputElement,
  React.ComponentProps<typeof Input> & {
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
      <Input id={id} {...props} ref={ref} />
    </div>
  );
});
