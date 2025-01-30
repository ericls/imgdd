type GetCaptchaResult = {
  token?: string;
  cleanup?: () => void;
};

const EMPTY_RESULT: GetCaptchaResult = {
  token: undefined,
  cleanup: undefined,
};

function getTokenTurnstile(
  elementId: string,
  action: string = "generic",
): Promise<GetCaptchaResult> {
  return new Promise((resolve) => {
    const key = window.TURNSTILE_SITE_KEY;
    const turnstile = window.turnstile;
    if (!key || !turnstile) {
      resolve(EMPTY_RESULT);
      return;
    }
    turnstile.ready(() => {
      const id = turnstile.render(elementId, {
        action: action,
        sitekey: key,
        callback: (token) => {
          resolve({
            token,
            cleanup,
          });
        },
        ["error-callback"]: (error) => {
          console.error(error);
          resolve(EMPTY_RESULT);
        },
      });
      if (!id) {
        resolve(EMPTY_RESULT);
        return;
      }
      const cleanup = () => {
        turnstile.remove(id);
      };
    });
  });
}

function getTokenRecaptcha(
  _elementId: string,
  action: string = "generic",
): Promise<GetCaptchaResult> {
  return new Promise((resolve) => {
    const key = window.RECAPTCHA_CLIENT_KEY;
    const grecaptcha = window.grecaptcha;
    if (!grecaptcha || !key) {
      resolve(EMPTY_RESULT);
    } else {
      grecaptcha.ready(function () {
        grecaptcha.execute(key, { action }).then(function (token) {
          resolve({
            token,
            cleanup: undefined,
          });
        });
      });
    }
  });
}

export function captchaProtected<T>(
  elementId: string,
  action: string,
  cb: (token?: string, cleanup?: () => void) => T,
): Promise<T> {
  const provider = window.CAPTCHA_PROVIDER;
  if (!provider) {
    return Promise.resolve(cb());
  }
  const getToken = (() => {
    if (provider === "recaptcha") {
      return getTokenRecaptcha;
    } else if (provider === "turnstile") {
      return getTokenTurnstile;
    } else {
      return undefined;
    }
  })();
  if (!getToken) {
    return Promise.resolve(cb());
  }
  return new Promise((resolve) => {
    getToken(
      elementId.startsWith("#") ? elementId : `#${elementId}`,
      action,
    ).then((result) => {
      resolve(cb(result.token, result.cleanup));
    });
  });
}
