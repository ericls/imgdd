import React from "react";
import { TextLogoSmall } from "./TextLogo";
import { DarkModeSettings } from "~src/ui/darkModeToggle";
import { useAuth } from "~src/lib/auth";
import { HiOutlineUser } from "react-icons/hi2";
import { Button } from "~src/ui/button";

function TopNavAuthInfo() {
  const { data: authData, isLoading: isAuthLoading } = useAuth();
  if (isAuthLoading) return null;
  if (!authData?.viewer.organizationUser?.id) return (
    <Button variant="transparent">
      <div className="flex items-center">
        <HiOutlineUser className="mr-2" size={"1.25rem"}/>
        Sign In
      </div>
    </Button>
  )
}

export function TopNav() {
  return (
    <div className="top-nav sticky top-0 backdrop-blur-sm p-2 flex justify-between text-neutral-800 dark:text-neutral-200 text-end">
      <div> {/* left */}
        <TextLogoSmall className="text-2xl" />
      </div>
      <div> {/* right */}
        <div className="flex items-center">
          <TopNavAuthInfo />
          <DarkModeSettings />
        </div>
      </div>
    </div>
  );
}
