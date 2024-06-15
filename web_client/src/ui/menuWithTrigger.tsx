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
import { useBreakpointName } from "./breakpoint";

type MenuWithTriggerProps = Omit<MenuProps, "open" | "onClose" | "ref"> & {
  trigger: React.ReactElement;
  placement?: Placement;
  containerClassName?: string;
  bottomFixedOnMobile?: boolean;
};

export function MenuWithTrigger({
  trigger,
  style,
  placement,
  containerClassName,
  menuSections,
  bottomFixedOnMobile = true,
  ...props
}: MenuWithTriggerProps) {
  const breakpointName = useBreakpointName();
  const [open, setOpen] = React.useState(false);
  const toggleOpen = React.useCallback(() => {
    setOpen((v) => !v);
  }, [setOpen]);
  const onClose = React.useCallback(() => {
    setOpen(false);
  }, [setOpen]);
  const { x, y, refs, strategy } = useFloating({
    whileElementsMounted: autoUpdate,
    placement: placement || "bottom-end",
    middleware: [flip()],
    strategy:
      breakpointName == "2xs" && bottomFixedOnMobile ? "fixed" : "absolute",
  });
  let finalStyle: React.CSSProperties = {
    ...style,
    position: strategy,
    top: y ?? 0,
    left: x ?? 0,
    width: "max-content",
  };
  if (breakpointName === "2xs" && bottomFixedOnMobile) {
    finalStyle = {
      ...style,
      position: strategy,
      width: "100%",
      minHeight: "220px",
      maxHeight: "calc(100vh - 2rem)",
      overflowY: "auto",
      bottom: 0,
      left: 0,
      borderTopLeftRadius: "1rem",
      borderTopRightRadius: "1rem",
    };
  }
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
        ref={refs.setReference}
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
      {/* Full screen backdrop on mobile */}
      <div
        className={classNames(
          "fixed inset-0 bg-black bg-opacity-80",
          "z-9",
          "h-full",
          {
            hidden: !(open && breakpointName === "2xs"),
          }
        )}
        onClick={onClose}
      />
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
          ref={refs.setFloating}
          menuSections={menuSections}
          {...props}
          style={finalStyle}
          onClose={onClose}
        />
      </Transition>
    </>
  );
}
