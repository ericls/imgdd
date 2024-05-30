import React from "react";
import { Menu, MenuProps } from "./menu";
import {
  useFloating,
  autoUpdate,
  flip,
  Placement,
} from "@floating-ui/react-dom";
import { Transition } from "@headlessui/react";
import classNames from "classnames";

type MenuWithTriggerProps = Omit<MenuProps, "open" | "onClose" | "ref"> & {
  trigger: React.ReactElement;
  placement?: Placement;
  containerClassName?: string;
};

export function MenuWithTrigger({
  trigger,
  style,
  placement,
  containerClassName,
  menuSections,
  ...props
}: MenuWithTriggerProps) {
  const [open, setOpen] = React.useState(false);
  const toggleOpen = React.useCallback(() => {
    setOpen((v) => !v);
  }, [setOpen]);
  const onClose = React.useCallback(() => {
    setOpen(false);
  }, [setOpen]);
  const { x, y, reference, floating, strategy } = useFloating({
    whileElementsMounted: autoUpdate,
    placement: placement || "bottom-end",
    middleware: [flip()],
  });
  const finalStyle = {
    ...style,
    position: strategy,
    top: y ?? 0,
    left: x ?? 0,
    width: "max-content",
  };
  const totalItemsCount = React.useMemo(() => {
    let res = 0;
    for (const section of menuSections.children) {
      res += section.items.length;
    }
    return res;
  }, [menuSections]);
  return (
    <>
      <div
        ref={reference}
        className={classNames("trigger-container w-min", containerClassName)}
        onClick={toggleOpen}
      >
        {React.cloneElement(trigger, {
          disabled: !totalItemsCount,
          className:
            trigger.props.className +
            " " +
            classNames({
              ["pointer-events-none opacity-50"]: !totalItemsCount,
            }),
        })}
      </div>
      <Transition
        as={React.Fragment}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
        show={open}
      >
        <Menu
          ref={floating}
          menuSections={menuSections}
          {...props}
          style={finalStyle}
          onClose={onClose}
        />
      </Transition>
    </>
  );
}
