import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import "dayjs/locale/en"; // or import other locales as needed
import i18n from "~src/i18n";

dayjs.extend(relativeTime);

export function useHumanizeDateTime({ datetimeStr }: { datetimeStr: string }) {
  dayjs.locale(i18n.language);

  // Convert the date string to a Day.js object and return relative time
  return dayjs(datetimeStr).fromNow();
}
