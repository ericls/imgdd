package httpserver

import (
	"io/fs"
	"os"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/utils"
)

type HttpServerConfigDef struct {
	Bind                   string
	WriteTimeout           int
	ReadTimeout            int
	TemplatesFS            fs.FS
	StaticFS               fs.FS
	SessionKey             string
	RedisURIForSession     string
	SiteName               string
	ImageDomain            string
	DefaultURLFormat       domainmodels.ImageURLFormat
	EnableGqlPlayground    bool
	EnableSafeImageCheck   bool
	SafeImageCheckEndpoint string
	CaptchaProvider        captchaProvider
	RecaptchaClientKey     string
	TurnstileSiteKey       string
	RecaptchaServerKey     string
	TurnstileSecretKey     string
}

func ReadServerConfigFromEnv() HttpServerConfigDef {
	return HttpServerConfigDef{
		Bind:                   os.Getenv("imgdd_HTTP_BIND"),
		WriteTimeout:           10,
		ReadTimeout:            10,
		SessionKey:             utils.GetEnv("IMGDD_SESSION_KEY", "NOT_SECURE_KEY"),
		RedisURIForSession:     utils.GetEnv("IMGDD_REDIS_URI_FOR_SESSION", "redis://localhost:30102"),
		SiteName:               utils.GetEnv("IMGDD_SITE_NAME", "imgdd"),
		ImageDomain:            utils.GetEnv("IMGDD_IMAGE_DOMAIN", ""),
		DefaultURLFormat:       domainmodels.ImageURLFormat(utils.GetEnv("IMGDD_DEFAULT_URL_FORMAT", "canonical")),
		EnableSafeImageCheck:   utils.IsStrTruthy(utils.GetEnv("IMGDD_ENABLE_SAFE_IMAGE_CHECK", "false")),
		SafeImageCheckEndpoint: utils.GetEnv("IMGDD_SAFE_IMAGE_CHECK_ENDPOINT", ""),
		CaptchaProvider:        captchaProvider(utils.GetEnv("IMGDD_CAPTCHA_PROVIDER", "off")),
		RecaptchaClientKey:     utils.GetEnv("IMGDD_RECAPTCHA_CLIENT_KEY", ""),
		TurnstileSiteKey:       utils.GetEnv("IMGDD_TURNSTILE_SITE_KEY", ""),
		RecaptchaServerKey:     utils.GetEnv("IMGDD_RECAPTCHA_SERVER_KEY", ""),
		TurnstileSecretKey:     utils.GetEnv("IMGDD_TURNSTILE_SECRET_KEY", ""),
	}
}
