import React from "react";
import { RiLoader4Line as LoaderIcon } from "react-icons/ri";

export function Loader({ size = 24 }: { size?: number }) {
  return (
    <div className="animate-spin">
      <LoaderIcon size={size} />
    </div>
  );
}

export function BlockLoader() {
  return (
    <div className="h-12 w-full flex justify-center items-center">
      <Loader size={36} />
    </div>
  );
}
