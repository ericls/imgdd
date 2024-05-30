export const BREAK_POINT_MIN_WIDTH = {
  "": 0,
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
SORTED_BREAKPOINT_TUPLE.sort((a, b) => a[0] - b[0]);

export function getBreakpointName(width: number): BreakPoints {
  for (const [w, name] of SORTED_BREAKPOINT_TUPLE) {
    if (width > w) return name as BreakPoints;
  }
  return "";
}
