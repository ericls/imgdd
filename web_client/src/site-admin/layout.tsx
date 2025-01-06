import classnames from "classnames";
import classNames from "classnames";
import React from "react";
import { Link, Outlet, useLocation } from "react-router-dom";
import {
  HiOutlineUsers as UsersIcon,
  HiOutlineServerStack as SotrageIcon,
} from "react-icons/hi2";
import { PiImages as ImagesIcon } from "react-icons/pi";
import { TbDeviceAnalytics as AnalyticsIcon } from "react-icons/tb";
import { Footer } from "~src/common/Footer";
import { LazyRouteFallback } from "~src/common/LazyRouteFallback";
import { TopNav } from "~src/common/TopNav";
import {
  PRIMARY_BORDER_COLOR,
  PRIMARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIMMER,
} from "~src/ui/classNames";

export function SiteAdminLayout() {
  const location = useLocation();
  const sideBarMenuGroups = React.useMemo(() => {
    return [
      {
        title: "General",
        items: [
          {
            title: "Users",
            to: "/site-admin/users",
            icon: <UsersIcon />,
            active: location.pathname.startsWith("/site-admin/users"),
          },
          {
            title: "Images",
            to: "/site-admin/images",
            icon: <ImagesIcon />,
            active: location.pathname.startsWith("/site-admin/images"),
          },
        ],
      },
      {
        title: "Storage",
        items: [
          {
            title: "Storage Backend",
            to: "/site-admin/storage",
            active: location.pathname.startsWith("/site-admin/storage"),
            icon: <SotrageIcon />,
          },
        ],
      },
      {
        title: "Analytics",
        items: [
          {
            title: "Access log",
            to: "/site-admin/analytics",
            active: location.pathname.startsWith("/site-admin/analytics"),
            icon: <AnalyticsIcon />,
          },
        ],
      },
    ];
  }, [location]);
  return (
    <div className="main min-h-full flex flex-col">
      <div className={classnames("mb-0 px-4")}>
        <TopNav />
      </div>
      <div className="grow relative z-0 flex">
        <div className="site-admin-sidebar basis-56">
          <div className={classNames("h-full overflow-y-auto sticky buttom-0")}>
            <div className="site-admin-sidebar-menu py-4 flex flex-col gap-4">
              {sideBarMenuGroups.map((group) => (
                <div key={group.title}>
                  <div
                    className={classnames(
                      SECONDARY_TEXT_COLOR_DIMMER,
                      "px-4 pb-1 text-xs font-poppins font-medium select-none border-l-4 border-transparent"
                    )}
                  >
                    {group.title}
                  </div>
                  <div className="">
                    {group.items.map((item) => (
                      <div
                        key={item.title}
                        className={classnames("w-full", {
                          [PRIMARY_TEXT_COLOR]: item.active,
                          [SECONDARY_TEXT_COLOR]: !item.active,
                        })}
                      >
                        <Link
                          to={item.to}
                          aria-label={item.title}
                          className={classnames(
                            "w-full flex py-1 pl-4 border-l-4 items-center transition-colors duration-200 ease-in-out",
                            "hover:text-neutral-800 hover:dark:text-neutral-100",
                            {
                              [PRIMARY_BORDER_COLOR]: item.active,
                              "border-transparent": !item.active,
                              "hover:border-neutral-400 hover:dark:border-neutral-500":
                                !item.active,
                            }
                          )}
                        >
                          <span className="inline-flex w-6 h-6 mr-2 items-center justify-center text-lg">
                            {item.icon}
                          </span>
                          {item.title}
                        </Link>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
        <div className="site-admin-main grow min-h-full flex flex-col mt-4 mr-4 ml-4">
          <div className="grow">
            <LazyRouteFallback />
            <Outlet />
          </div>
          <div>
            <Footer />
          </div>
        </div>
      </div>
    </div>
  );
}

export default SiteAdminLayout;
