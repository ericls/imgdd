import React from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { TextLogoSmall } from "./TextLogo";
import { DarkModeSettings } from "~src/ui/darkModeToggle";
import { useAuth } from "~src/lib/auth";
import { HiOutlineUser } from "react-icons/hi2";
import { Button } from "~src/ui/button";
import { SECONDARY_TEXT_COLOR_DIM } from "~src/ui/classNames";
import { MenuWithTrigger } from "~src/ui/menuWithTrigger";
import { MenuSection } from "~src/ui/menu";
import classNames from "classnames";

function TopNavAuthInfo() {
  const { t } = useTranslation();
  const { data: authData, isLoading: isAuthLoading, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const onSignInClick = React.useCallback(() => {
    if (location.pathname === "/auth") return;
    navigate("/auth");
  }, [location.pathname, navigate]);
  if (isAuthLoading) return null;
  const orgUser = authData?.viewer.organizationUser;
  const hasSiteOwnerAccess = authData?.viewer.hasSiteOwnerAccess;
  if (!orgUser?.id) {
    return (
      <Button
        variant="transparent"
        onClick={onSignInClick}
        className={SECONDARY_TEXT_COLOR_DIM}
      >
        <div className="flex items-center">
          <HiOutlineUser className="mr-2" size={"1.25rem"} />
          {t("topNav.signIn")}
        </div>
      </Button>
    );
  }
  return (
    <div className="flex items-center">
      <MenuWithTrigger
        trigger={
          <Button variant="transparent">
            <HiOutlineUser
              className={SECONDARY_TEXT_COLOR_DIM}
              size={"1.25rem"}
            />
          </Button>
        }
        menuSections={{
          children: [
            {
              id: "main",
              items: [
                {
                  id: "profile",
                  children: t("topNav.profile"),
                  action: () => navigate("/profile"),
                },
                {
                  id: "signOut",
                  children: t("topNav.signOut"),
                  action: logout,
                },
              ],
            },
            hasSiteOwnerAccess && {
              id: "site_owner",
              items: [
                {
                  id: "siteSettings",
                  children: t("topNav.siteSettings"),
                  action: () => navigate("/site-admin"),
                },
              ],
            },
          ].filter(Boolean) as MenuSection[],
        }}
      />
    </div>
  );
}

type TopNavProps = {
  hideLogo?: boolean;
  leftContent?: React.ReactNode;
};

export function TopNav({ hideLogo, leftContent }: TopNavProps) {
  return (
    <div className="top-nav sticky top-0 p-2 flex justify-between text-neutral-800 dark:text-neutral-200 text-end z-50">
      <div className="flex items-center">
        {/* left */}
        <Link to="/" className={classNames({ hidden: hideLogo })}>
          <TextLogoSmall className="text-2xl" />
        </Link>
        {leftContent}
      </div>
      <div>
        {/* right */}
        <div className="flex items-center">
          <TopNavAuthInfo />
          <DarkModeSettings />
        </div>
      </div>
    </div>
  );
}
