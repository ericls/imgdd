function getUserLocaleLower() {
  return navigator.language.toLowerCase();
}

export function isChinese() {
  const lang = getUserLocaleLower();
  return lang === "zh" || lang.startsWith("zh-");
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
