interface Window {
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
}
