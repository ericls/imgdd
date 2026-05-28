function getUserLocaleLower() {
  return navigator.language.toLowerCase();
}

export function isChinese() {
  const lang = getUserLocaleLower();
  return lang === "zh" || lang.startsWith("zh-");
}

export function isThai() {
  const lang = getUserLocaleLower();
  return lang === "th" || lang.startsWith("th-");
}

export function isKorean() {
  const lang = getUserLocaleLower();
  return lang === "ko" || lang.startsWith("ko-");
}

export function isRussian() {
  const lang = getUserLocaleLower();
  return lang === "ru" || lang.startsWith("ru-");
}

export type SupportedLanguage = "en" | "zh_hans" | "th" | "ko" | "ru";

const LANGUAGE_STORAGE_KEY = "imgdd.language";

export function getStoredLanguage(): SupportedLanguage | null {
  try {
    const v = localStorage.getItem(LANGUAGE_STORAGE_KEY);
    if (v === "en" || v === "zh_hans" || v === "th" || v === "ko" || v === "ru")
      return v;
  } catch {
    // localStorage may be unavailable (SSR, privacy mode)
  }
  return null;
}

export function setStoredLanguage(lang: SupportedLanguage | null) {
  try {
    if (lang === null) {
      localStorage.removeItem(LANGUAGE_STORAGE_KEY);
    } else {
      localStorage.setItem(LANGUAGE_STORAGE_KEY, lang);
    }
  } catch {
    // ignore
  }
}

export function getInitialLanguage(): SupportedLanguage {
  const stored = getStoredLanguage();
  if (stored) return stored;
  if (isChinese()) return "zh_hans";
  if (isThai()) return "th";
  if (isKorean()) return "ko";
  if (isRussian()) return "ru";
  return "en";
}

export function isChinaTimezone() {
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  return [
    "Asia/Shanghai",
    "Asia/Chongqing",
    "Asia/Harbin",
    "Asia/Urumqi",
  ].includes(timezone);
}
