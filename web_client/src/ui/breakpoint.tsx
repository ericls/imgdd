import React from "react";

export const BREAK_POINT_MIN_WIDTH = {
  "2xs": 0,
  xs: 480,
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
  "2xl": 1536,
} as const;

export type BreakPoints = keyof typeof BREAK_POINT_MIN_WIDTH;

const SORTED_BREAKPOINT_TUPLE: [width: number, name: string][] = Object.entries(
  BREAK_POINT_MIN_WIDTH
).map(([k, v]) => [v, k]);
SORTED_BREAKPOINT_TUPLE.sort((a, b) => b[0] - a[0]);

export function getBreakpointName(width?: number): BreakPoints {
  width = width || window.innerWidth;
  for (const [w, name] of SORTED_BREAKPOINT_TUPLE) {
    if (width > w) return name as BreakPoints;
  }
  return "2xs";
}

export function useBreakpointName() {
  const [breakpointName, setBreakpointName] = React.useState<BreakPoints>(
    getBreakpointName()
  );
  React.useEffect(() => {
    const onResize = () => {
      setBreakpointName(getBreakpointName());
    };
    window.addEventListener("resize", onResize);
    return () => {
      window.removeEventListener("resize", onResize);
    };
  }, []);
  return breakpointName;
}
