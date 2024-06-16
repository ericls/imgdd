import classnames from "classnames";
import classNames from "classnames";
import React from "react";
import { Link, Outlet, useLocation } from "react-router-dom";
import { Footer } from "~src/common/Footer";
import { LazyRouteFallback } from "~src/common/LazyRouteFallback";
import { TopNav } from "~src/common/TopNav";
import {
  LINK_COLOR,
  PRIMARY_BORDER_COLOR,
  PRIMARY_BORDER_COLOR_ON_HOVER,
  PRIMARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIM,
  SECONDARY_TEXT_COLOR_DIMMER,
  SECOND_LAYER,
  SECOND_LAYER_HOVER,
} from "~src/ui/classNames";

export function SiteAdminLayout() {
  const location = useLocation();
  const sideBarMenuGroups = React.useMemo(() => {
    console.log(location.pathname);
    return [
      {
        title: "General",
        items: [
          {
            title: "Overview",
            to: "/site-admin",
          },
          {
            title: "Configuration",
            to: "/site-admin/configuration",
          },
          {
            title: "Users",
            to: "/site-admin/users",
          },
          {
            title: "Groups",
            to: "/site-admin/groups",
          },
          {
            title: "Permissions",
            to: "/site-admin/permissions",
          },
          {
            title: "Audit Log",
            to: "/site-admin/audit-log",
          },
          {
            title: "Site Admin",
            to: "/site-admin/site-admin",
          },
        ],
      },
      {
        title: "Storage",
        items: [
          {
            title: "Storage Config",
            to: "/site-admin/storage",
            active: location.pathname.startsWith("/site-admin/storage"),
          },
          {
            title: "Storage Buckets",
            to: "/site-admin/storage-buckets",
          },
          {
            title: "Storage Files",
            to: "/site-admin/storage-files",
          },
        ],
      },
    ];
  }, [location]);
  console.log(sideBarMenuGroups);
  return (
    <div className="main min-h-full flex flex-col mx-2">
      <TopNav />
      <div className="grow relative z-0 flex gap-6">
        <div className="site-admin-sidebar basis-56 pb-2">
          <div
            className={classNames(
              SECOND_LAYER,
              "rounded-md h-full overflow-y-auto sticky buttom-0"
            )}
          >
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
                          [SECONDARY_TEXT_COLOR_DIM]: !item.active,
                        })}
                      >
                        <Link
                          to={item.to}
                          aria-label={item.title}
                          className={classnames(
                            "w-full block py-1 pl-4 border-l-4",
                            "hover:text-neutral-800 hover:dark:text-neutral-100",
                            {
                              [PRIMARY_BORDER_COLOR]: item.active,
                              "border-transparent": !item.active,
                              "hover:border-neutral-400 hover:dark:border-neutral-500":
                                !item.active,
                            }
                          )}
                        >
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
        <div className="site-admin-main grow min-h-full flex flex-col">
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
