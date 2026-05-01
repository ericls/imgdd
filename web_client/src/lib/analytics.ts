declare global {
  interface Window {
    GAID?: string;
    gtag?: (...args: unknown[]) => void;
  }
}

export type AnalyticsProps = Record<
  string,
  string | number | boolean | undefined
>;

export type AnalyticsEvent = {
  name: string;
  props?: AnalyticsProps;
};

export interface AnalyticsProvider {
  track(event: AnalyticsEvent): void;
}

const providers: AnalyticsProvider[] = [];

export function registerProvider(provider: AnalyticsProvider) {
  providers.push(provider);
}

export function track(name: string, props?: AnalyticsProps) {
  const event: AnalyticsEvent = { name, props };
  for (const provider of providers) {
    try {
      provider.track(event);
    } catch {
      // swallow — analytics must never break the app
    }
  }
}

export const GoogleAnalyticsProvider: AnalyticsProvider = {
  track({ name, props }) {
    if (typeof window === "undefined" || typeof window.gtag !== "function") {
      return;
    }
    window.gtag("event", name, props ?? {});
  },
};

if (typeof window !== "undefined" && window.GAID) {
  registerProvider(GoogleAnalyticsProvider);
}
