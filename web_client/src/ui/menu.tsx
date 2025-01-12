import React from "react";
import cx from "classnames";

export { TbDots as DefaultMenuIcon } from "react-icons/tb";

export type MenuItemVariant = "danger" | "normal";
export type MenuItemActionArgs = {
  event: React.MouseEvent<HTMLElement>;
  close?: () => void;
};
export type MenuItem = {
  id: string;
  children: React.ReactNode;
  action?: (args: MenuItemActionArgs) => void;
  variant?: MenuItemVariant;
  disabled?: boolean;
  className?: string;
};

export type MenuSection = {
  id: string;
  items: MenuItem[];
};

export type MenuSections = {
  children: MenuSection[];
};

export type MenuProps = {
  onClose?: () => void;
  closeOnOutside?: boolean;
  menuSections: MenuSections;
} & JSX.IntrinsicElements["div"];
export const Menu = React.forwardRef<HTMLDivElement, MenuProps>(
  (
    { onClose, menuSections, closeOnOutside = true, className, ...props },
    containerRef
  ) => {
    const localContainerRef = React.useRef<HTMLDivElement>();
    const menuIdToElement = React.useRef<Map<string, HTMLElement>>(new Map());
    const finalContainerRef = React.useCallback(
      (el: HTMLDivElement) => {
        if (containerRef && "current" in containerRef) {
          containerRef.current = el;
        } else {
          containerRef?.(el);
        }
        localContainerRef.current = el;
      },
      [containerRef]
    );
    const itemRef = React.useCallback((el: HTMLButtonElement) => {
      const id = el?.getAttribute("data-item-id");
      if (id) {
        menuIdToElement.current.set(id, el);
      } else {
        menuIdToElement.current.delete(id || "");
      }
    }, []);
    const [menuIdToItem, menuItemIds]: [{ [key: string]: MenuItem }, string[]] =
      React.useMemo(() => {
        const items: { [key: string]: MenuItem } = {};
        const ids: string[] = [];
        for (const section of menuSections.children) {
          for (const item of section.items) {
            if (item.disabled) continue;
            if (item.id in items) {
              console.warn("Duplicated menu item id", item.id, item);
            }
            items[item.id] = item;
            ids.push(item.id);
          }
        }
        return [items, ids];
      }, [menuSections]);
    const [activeId, setActiveId] = React.useState<undefined | string>();
    const moveSelection = React.useCallback(
      (direction: -1 | 1) => {
        setActiveId((currentActiveId) => {
          const currentIndex = menuItemIds.indexOf(currentActiveId || "");
          return menuItemIds.at(
            (currentIndex + direction) % menuItemIds.length
          );
        });
      },
      [setActiveId, menuItemIds]
    );
    React.useLayoutEffect(() => {
      localContainerRef.current?.focus();
    }, []);
    React.useEffect(() => {
      const el = menuIdToElement.current.get(activeId || "");
      el?.focus();
    }, [activeId]);
    const onMouseEnterItem = React.useCallback(
      (e: React.MouseEvent<HTMLElement>) => {
        const div = e.currentTarget;
        const id = div.getAttribute("data-item-id");
        setActiveId(id || undefined);
      },
      [setActiveId]
    );
    const onClickItem = React.useCallback(
      (e: React.MouseEvent<HTMLElement>) => {
        const div = e.currentTarget;
        const id = div.getAttribute("data-item-id");
        const item = menuIdToItem[id || ""];
        if (item === undefined) return;
        item.action?.({ close: onClose, event: e });
      },
      [menuIdToItem, onClose]
    );
    const onKeypress = React.useCallback(
      (e: React.KeyboardEvent<HTMLElement>) => {
        if (e.key === "Tab" || e.key === "ArrowDown") {
          moveSelection(1);
          e.preventDefault();
          e.stopPropagation();
        } else if ((e.key === "Tab" && e.shiftKey) || e.key === "ArrowUp") {
          moveSelection(-1);
          e.preventDefault();
          e.stopPropagation();
        } else if (e.key == "Escape") {
          onClose?.();
        }
      },
      [moveSelection, onClose]
    );
    const onMouseLeaveContaienr = React.useCallback(() => {
      // TODO: only clear selection if the last selection is made by mouse events
      setActiveId(undefined);
    }, [setActiveId]);
    React.useEffect(() => {
      if (!closeOnOutside) return;
      const listener = (e: MouseEvent) => {
        if (
          e.target instanceof HTMLElement &&
          localContainerRef.current &&
          localContainerRef.current.contains(e.target)
        ) {
          return;
        } else {
          onClose?.();
        }
      };
      document.addEventListener("click", listener);
      return () => {
        document.removeEventListener("click", listener);
      };
    }, [onClose, closeOnOutside]);
    return (
      <div
        ref={finalContainerRef}
        {...props}
        tabIndex={0}
        onKeyDown={onKeypress}
        onMouseLeave={onMouseLeaveContaienr}
        className={cx(
          "menu",
          "mt-2 rounded shadow-md ring-1 ring-black ring-opacity-5",
          "bg-neutral-50 dark:bg-neutral-700 dark:shadow-neutral-800",
          "divide-y dark:divide-neutral-800",
          "z-20",
          className
        )}
      >
        {menuSections.children.map((section) => {
          return section.items.length ? (
            <div className="menu-section p-2" key={section.id}>
              {section.items.map((item, i) => {
                const itemElement = (
                  <button
                    tabIndex={i + 1}
                    ref={itemRef}
                    onMouseEnter={onMouseEnterItem}
                    data-disabled={item.disabled}
                    key={item.id}
                    data-item-id={item.id}
                    type="button"
                    disabled={item.disabled}
                    className={cx(
                      item.className,
                      "cursor-pointer block rounded py-1 px-2 focus:outline-none w-full text-start",
                      {
                        active: activeId === item.id,
                        "bg-neutral-200 dark:bg-neutral-800":
                          activeId === item.id,
                        "text-red-600 dark:text-red-500":
                          item.variant === "danger",
                      }
                    )}
                    onClick={onClickItem}
                  >
                    {item.children}
                  </button>
                );
                return itemElement;
              })}
            </div>
          ) : null;
        })}
      </div>
    );
  }
);
