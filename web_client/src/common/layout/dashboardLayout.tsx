import React from "react";
import { Link, Outlet, useLocation } from "react-router-dom";
import classNames from "classnames";
import { TopNav } from "~src/common/TopNav";
import {
  PRIMARY_BORDER_COLOR,
  PRIMARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIMMER,
} from "~src/ui/classNames";
import { Footer } from "~src/common/Footer";
import { LazyRouteFallback } from "~src/common/LazyRouteFallback";

export type DashboardLayoutMenuItem = {
  title: string;
  to: string;
  icon: React.ReactNode;
  active: (location: { pathname: string }) => boolean;
};
export type DashboardLayoutMenuGroup = {
  title: string;
  items: DashboardLayoutMenuItem[];
};
export type DashboardLayoutMenuGroups = DashboardLayoutMenuGroup[];

export function DashboardLayout({
  menuGroups,
  mainAreaClassName,
  children,
}: {
  menuGroups: DashboardLayoutMenuGroups;
  children?: React.ReactNode;
  mainAreaClassName?: string;
}) {
  const location = useLocation();
  return (
    <div className="main min-h-full flex flex-col">
      <div className={classNames("mb-0 px-2")}>
        <TopNav />
      </div>
      <div className="grow relative z-0 flex">
        <div className="site-admin-sidebar basis-56">
          <div className={classNames("h-full overflow-y-auto sticky buttom-0")}>
            <div className="site-admin-sidebar-menu py-4 flex flex-col gap-4">
              {menuGroups.map((group) => (
                <div key={group.title}>
                  <div
                    className={classNames(
                      SECONDARY_TEXT_COLOR_DIMMER,
                      "px-2 pb-1 text-xs font-poppins font-medium select-none border-l-4 border-transparent"
                    )}
                  >
                    {group.title}
                  </div>
                  <div className="">
                    {group.items.map((item) => {
                      const isActive = item.active(location);
                      return (
                        <div
                          key={item.title}
                          className={classNames("w-full", {
                            [PRIMARY_TEXT_COLOR]: isActive,
                            [SECONDARY_TEXT_COLOR]: !isActive,
                          })}
                        >
                          <Link
                            to={item.to}
                            aria-label={item.title}
                            className={classNames(
                              "w-full flex py-1 pl-2 border-l-4 items-center transition-colors duration-200 ease-in-out",
                              "hover:text-neutral-800 hover:dark:text-neutral-100",
                              {
                                [PRIMARY_BORDER_COLOR]: isActive,
                                "border-transparent": !isActive,
                                "hover:border-neutral-400 hover:dark:border-neutral-500":
                                  !isActive,
                              }
                            )}
                          >
                            <span className="inline-flex w-6 h-6 mr-2 items-center justify-center text-lg">
                              {item.icon}
                            </span>
                            {item.title}
                          </Link>
                        </div>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
        <div
          className={classNames(
            mainAreaClassName,
            "grow min-h-full flex flex-col mt-4 mr-4 ml-4"
          )}
        >
          <div className="grow">
            <LazyRouteFallback />
            {children || <Outlet />}
          </div>
          <div>
            <Footer />
          </div>
        </div>
      </div>
    </div>
  );
}
