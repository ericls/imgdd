import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { getInitialLanguage, SupportedLanguage } from "./lib/locale";
import EN from "./localization/en.json";

type AsyncLoader = () => Promise<{ default: Record<string, unknown> }>;

const LAZY_LOADERS: Partial<Record<SupportedLanguage, AsyncLoader>> = {
  zh_hans: () => import("./localization/zh_hans.json"),
  th: () => import("./localization/th.json"),
};

export async function loadLanguage(lang: SupportedLanguage) {
  if (i18n.hasResourceBundle(lang, "translation")) return;
  const loader = LAZY_LOADERS[lang];
  if (!loader) return;
  const mod = await loader();
  i18n.addResourceBundle(lang, "translation", mod);
}

const initialLanguage = getInitialLanguage();

const INITIAL_LANGUAGE_TIMEOUT_MS = 1_000;

const init = i18n
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
    if (initialLanguage === "en") return;
    await loadLanguage(initialLanguage);
    await i18n.changeLanguage(initialLanguage);
  });

const timeout = new Promise<void>((_, reject) =>
  setTimeout(
    () =>
      reject(
        new Error(
          `i18n initial language load timed out after ${INITIAL_LANGUAGE_TIMEOUT_MS}ms`,
        ),
      ),
    INITIAL_LANGUAGE_TIMEOUT_MS,
  ),
);

export const i18nReady = Promise.race([init, timeout]).catch((e) => {
  console.warn("i18n init failed, falling back to en", e);
});

export default i18n;
