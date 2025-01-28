package httpserver

type captchaProvider string

const (
	captchaProviderGoogleRecaptcha     captchaProvider = "recaptcha"
	captchaProviderCloudflareTurnstile captchaProvider = "turnstile"
	captchaProviderOff                 captchaProvider = "off"
)

func (p captchaProvider) isValid() bool {
	switch p {
	case captchaProviderGoogleRecaptcha, captchaProviderCloudflareTurnstile, captchaProviderOff:
		return true
	}
	return false
}
