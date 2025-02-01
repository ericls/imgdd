import React from "react";
import { Link, Outlet, useLocation } from "react-router-dom";
import classNames from "classnames";
import { HiOutlineBars3 as MenuIcon } from "react-icons/hi2";
import { TopNav } from "~src/common/TopNav";
import {
  BASE_LAYER,
  PRIMARY_BORDER_COLOR,
  PRIMARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIMMER,
} from "~src/ui/classNames";
import { Footer } from "~src/common/Footer";
import { LazyRouteFallback } from "~src/common/LazyRouteFallback";
import { TextLogo, TextLogoSmall } from "../TextLogo";
import { Button } from "~src/ui/button";

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
  const [isSidebarOpen, setIsSidebarOpen] = React.useState(false);
  const toggleSidebar = React.useCallback(
    () => setIsSidebarOpen((prev) => !prev),
    [setIsSidebarOpen],
  );
  React.useEffect(() => {
    setIsSidebarOpen(false);
  }, [location, setIsSidebarOpen]);
  const sidebarToggle = React.useMemo(() => {
    return (
      <Button
        className={classNames("md:hidden", { hidden: isSidebarOpen })}
        onClick={toggleSidebar}
        aria-label="Toggle Sidebar"
        variant="transparent"
      >
        {<MenuIcon />}
      </Button>
    );
  }, [toggleSidebar, isSidebarOpen]);
  return (
    <div className="main min-h-full flex flex-col">
      <div className="grow relative z-0 flex max-h-screen overflow-hidden">
        <div
          className={classNames(
            "sidebar-backdrop",
            {
              hidden: !isSidebarOpen,
            },
            "md:hidden",
            "fixed top-0 left-0 bottom-0 right-0 z-10",
            "bg-black bg-opacity-50",
          )}
          onClick={toggleSidebar}
        ></div>
        <div
          className={classNames(
            "h-screen buttom-0 top-0 overscroll-y-none",
            "dashboard-sidebar w-60 absolute left-0 top-0 bottom-0 z-20 transition-all duration-200 ease-in-out",
            {
              "left-[-100%]": !isSidebarOpen,
            },
            BASE_LAYER,
            "md:block md:sticky md:basis-56 md:w-auto",
          )}
        >
          <div className={classNames("flex flex-col")}>
            <div className="flex items-center justify-start pl-4 p-2">
              <Link to="/" className="block">
                <TextLogoSmall className="text-2xl" />
              </Link>
              <Button className="invisible">
                <div className="flex items-center">
                  <MenuIcon className="mr-2" size={24} />
                </div>
              </Button>
            </div>
            <div className="dashboard-sidebar-menu py-4 flex flex-col gap-4">
              {menuGroups.map((group) => (
                <div key={group.title}>
                  <div
                    className={classNames(
                      SECONDARY_TEXT_COLOR_DIMMER,
                      "px-2 pb-1 text-xs font-poppins font-medium select-none border-l-4 border-transparent",
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
                              },
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
        <div className="screen-container overflow-y-auto grow">
          <div className={classNames("top-0 z-50", BASE_LAYER)}>
            <TopNav hideLogo leftContent={sidebarToggle} />
          </div>
          <div
            className={classNames(
              mainAreaClassName,
              "grow min-h-full flex flex-col mt-4 mr-4 ml-4",
            )}
          >
            <div className="grow w-full mx-auto xl:max-w-7xl">
              <LazyRouteFallback />
              {children || <Outlet />}
            </div>
            <div>
              <Footer />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
