import { TextLogoSmall } from "./TextLogo";

export function TopNav() {
  return (
    <div className="top-nav sticky top-0 backdrop-blur-sm p-2 flex justify-between text-neutral-800 dark:text-neutral-200 text-end">
      <TextLogoSmall className="text-2xl" />
      <h1>TopNav</h1>
    </div>
  );
}
