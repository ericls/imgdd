import React from "react";

type SlotsFiller = string | ((r: typeof React) => JSX.Element | string);

const SlotsContext = React.createContext<{
  slots: Record<string, SlotsFiller | undefined>;
}>({
  slots: {},
});

export function SlotsProvider({ children }: { children: React.ReactNode }) {
  const slots: Record<string, SlotsFiller> = React.useMemo(() => {
    return (window.IMGDD_PLUGINS || []).reduce(
      (acc: Record<string, SlotsFiller>, plugin) => {
        return {
          ...acc,
          ...plugin.textSlots,
        };
      },
      {},
    );
  }, []);
  return (
    <SlotsContext.Provider value={{ slots }}>{children}</SlotsContext.Provider>
  );
}

export function Slot({
  id,
  fallback,
}: {
  id: string;
  fallback: JSX.Element | string;
}) {
  const { slots } = React.useContext(SlotsContext);
  const slot = slots[id];
  if (!slot) {
    return <>{fallback}</>;
  }
  const content = typeof slot === "function" ? slot(React) : slot;
  return <React.Fragment key={id}>{content}</React.Fragment>;
}
