import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { getInitialLanguage, SupportedLanguage } from "./lib/locale";
import EN from "./localization/en.json";

type AsyncLoader = () => Promise<{ default: Record<string, unknown> }>;

const LAZY_LOADERS: Partial<Record<SupportedLanguage, AsyncLoader>> = {
  zh_hans: () => import("./localization/zh_hans.json"),
};

export async function loadLanguage(lang: SupportedLanguage) {
  if (i18n.hasResourceBundle(lang, "translation")) return;
  const loader = LAZY_LOADERS[lang];
  if (!loader) return;
  const mod = await loader();
  i18n.addResourceBundle(lang, "translation", mod);
}

const initialLanguage = getInitialLanguage();

export const i18nReady = i18n
  .use(initReactI18next)
  .init({
    resources: { en: { translation: EN } },
    lng: "en",
    fallbackLng: "en",
    interpolation: { escapeValue: false },
    react: {
      transKeepBasicHtmlNodesFor: ["br", "strong", "i", "p", "em", "b"],
    },
  })
  .then(async () => {
    if (initialLanguage !== "en") {
      await loadLanguage(initialLanguage);
      await i18n.changeLanguage(initialLanguage);
    }
  });

export default i18n;
