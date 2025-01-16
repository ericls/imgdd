import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import "dayjs/locale/en"; // or import other locales as needed
import { useTranslation } from "react-i18next";

dayjs.extend(relativeTime);

export function useHumanizeDateTime({ datetimeStr }: { datetimeStr: string }) {
  const { i18n } = useTranslation();
  dayjs.locale(i18n.language);

  // Convert the date string to a Day.js object and return relative time
  return dayjs(datetimeStr).fromNow();
}
