import React from "react";

export function getIsDark() {
  return document.documentElement.getAttribute("data-theme") === "dark";
}

export type DarkModeTheme = "dark" | "light" | "system";

export function setThemeAttr() {
  if (
    localStorage["imgdd-theme"] === "dark" ||
    (!("imgdd-theme" in localStorage) &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)
  ) {
    document.documentElement.setAttribute("data-theme", "dark");
  } else {
    document.documentElement.setAttribute("data-theme", "light");
  }
}

export function setTheme(theme: DarkModeTheme) {
  if (theme === "system") {
    localStorage.removeItem("imgdd-theme");
  } else {
    localStorage["imgdd-theme"] = theme;
  }
  setThemeAttr();
}

export function getTheme(): DarkModeTheme {
  if (!("imgdd-theme" in localStorage)) return "system";
  if (localStorage["imgdd-theme"] === "dark") return "dark";
  return "light";
}

setThemeAttr();

const darkModePreference = window.matchMedia("(prefers-color-scheme: dark)");
darkModePreference.addEventListener("change", setThemeAttr);

type DarkModeContextType = {
  isDarkMode: boolean;
  theme: DarkModeTheme;
  setTheme: (v: DarkModeTheme) => void;
};

const lightContext = {
  isDarkMode: false,
  setTheme,
};

export const DarkModeContext = React.createContext<DarkModeContextType>({
  ...lightContext,
  theme: "system",
});

export function DarkModeProvider({ children }: { children?: React.ReactNode }) {
  const [flag, setFlag] = React.useState(true);
  const onChange = React.useCallback(() => {
    setFlag((v) => !v);
  }, []);
  const contextValue: DarkModeContextType = React.useMemo(() => {
    return {
      setTheme,
      isDarkMode: getIsDark(),
      theme: getTheme(),
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [flag]);
  React.useEffect(() => {
    const mutationObserver = new MutationObserver(onChange);
    mutationObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["data-theme"],
    });
    return () => {
      mutationObserver.disconnect();
    };
  }, [onChange]);
  return (
    <DarkModeContext.Provider value={contextValue}>
      {children}
    </DarkModeContext.Provider>
  );
}
