interface Window {
  // Server side rendered data
  SITE_NAME: string;
  DEBUG?: boolean;
  TERMS_AND_CONDITIONS_URL: string;
  PRIVACY_POLICY_URL: string;
  GOOGLE_LOGIN_URL: string;
  SUPPORT_EMAIL: string;
  HOME_PAGE: string;
  SENTRY_DSN?: string;
  SENTRY_ENVIRONMENT?: string;
  STRIPE_PUB_KEY?: string;
  LOGO?: string;
  FLAGS?: {
    linkTagsEnabled?: boolean;
    deepLinkEnabled?: boolean;
  };
  VERSION: string;
  SESSION_HEADER_NAME: string;
  CAPTCHA_PROVIDER?: "recaptcha" | "turnstile" | string;
  RECAPTCHA_CLIENT_KEY?: string;
  TURNSTILE_SITE_KEY?: string;
  // plugins
  IMGDD_PLUGINS?: IMGDDPlugin[];
}
