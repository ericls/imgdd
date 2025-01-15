import React from "react";
import classNames from "classnames";
import { TEXT_COLOR } from "~src/ui/classNames";

export function Users() {
  return (
    <div className="flex w-full items-center justify-center">
      <div className="max-w-md text-center px-4">
        <h1 className={classNames("mb-4 text-4xl font-bold", TEXT_COLOR)}>
          Under Construction 🚧
        </h1>
        <p
          className={classNames("text-base sm:text-lg TEXT_COLOR", TEXT_COLOR)}
        >
          We&apos;re working hard to finish this page. Please check back soon!
        </p>
      </div>
    </div>
  );
}
