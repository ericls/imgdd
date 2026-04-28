import type { JSX } from "react";

export type IMGDDPlugin = {
  textSlots: {
    [key: string]: JSX.Element | (() => JSX.Element);
  };
};
