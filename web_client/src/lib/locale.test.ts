import { afterEach, describe, expect, it, vi } from "vitest";

import {
  getInitialLanguage,
  getStoredLanguage,
  isChinaTimezone,
  setStoredLanguage,
} from "./locale";

function setNavigatorLanguage(language: string) {
  vi.spyOn(navigator, "language", "get").mockReturnValue(language);
}

afterEach(() => {
  localStorage.clear();
  vi.restoreAllMocks();
});

describe("stored language", () => {
  it("stores and reads supported language values", () => {
    setStoredLanguage("ko");

    expect(getStoredLanguage()).toBe("ko");
  });

  it("ignores unsupported stored language values", () => {
    localStorage.setItem("imgdd.language", "fr");

    expect(getStoredLanguage()).toBeNull();
  });

  it("removes the stored language when set to null", () => {
    setStoredLanguage("th");
    setStoredLanguage(null);

    expect(getStoredLanguage()).toBeNull();
  });
});

describe("getInitialLanguage", () => {
  it("prefers a supported stored language over navigator language", () => {
    setNavigatorLanguage("zh-TW");
    setStoredLanguage("ru");

    expect(getInitialLanguage()).toBe("ru");
  });

  it("detects traditional Chinese locales", () => {
    setNavigatorLanguage("zh-Hant");

    expect(getInitialLanguage()).toBe("zh_hant");
  });

  it("detects simplified Chinese locales", () => {
    setNavigatorLanguage("zh-CN");

    expect(getInitialLanguage()).toBe("zh_hans");
  });

  it("falls back to English for unsupported locales", () => {
    setNavigatorLanguage("fr-FR");

    expect(getInitialLanguage()).toBe("en");
  });
});

describe("isChinaTimezone", () => {
  it("returns true for configured China timezones", () => {
    vi.spyOn(Intl, "DateTimeFormat").mockReturnValue({
      resolvedOptions: () => ({ timeZone: "Asia/Shanghai" }),
    } as Intl.DateTimeFormat);

    expect(isChinaTimezone()).toBe(true);
  });

  it("returns false for other timezones", () => {
    vi.spyOn(Intl, "DateTimeFormat").mockReturnValue({
      resolvedOptions: () => ({ timeZone: "Etc/UTC" }),
    } as Intl.DateTimeFormat);

    expect(isChinaTimezone()).toBe(false);
  });
});
