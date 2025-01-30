package captcha

import (
	"context"
	"net/http"
)

type CaptchaProvider string
type captchaTokenKey string

const (
	CaptchaProviderGoogleRecaptcha     CaptchaProvider = "recaptcha"
	CaptchaProviderCloudflareTurnstile CaptchaProvider = "turnstile"
	CaptchaProviderOff                 CaptchaProvider = "off"

	captchaTokenKeyContextKey = captchaTokenKey("captchaToken")
)

func (p CaptchaProvider) IsValid() bool {
	switch p {
	case CaptchaProviderGoogleRecaptcha, CaptchaProviderCloudflareTurnstile, CaptchaProviderOff:
		return true
	}
	return false
}

func MakeClient(provider CaptchaProvider, recaptchaKey, turnstileKey string) CaptchaClient {
	switch provider {
	case CaptchaProviderGoogleRecaptcha:
		return NewRecaptchaClient(recaptchaKey)
	case CaptchaProviderCloudflareTurnstile:
		return NewTurnStileClient(turnstileKey)
	case CaptchaProviderOff:
		return nil
	}
	return nil
}

func GetToken(ctx context.Context) string {
	v, ok := ctx.Value(captchaTokenKeyContextKey).(string)
	if !ok {
		return ""
	}
	return v
}

func MakeHttpMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captchaToken := r.URL.Query().Get("captchaToken")
			ctx := context.WithValue(r.Context(), captchaTokenKeyContextKey, captchaToken)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
