package httpserver

import (
	"io/fs"
	"os"
	"strconv"

	"github.com/ericls/imgdd/captcha"
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
	RedisURI               string
	SiteName               string
	SiteTitle              string
	ImageDomain            string
	DefaultURLFormat       domainmodels.ImageURLFormat
	EnableGqlPlayground    bool
	EnableSafeImageCheck   bool
	SafeImageCheckEndpoint string
	CaptchaProvider        captcha.CaptchaProvider
	RecaptchaClientKey     string
	TurnstileSiteKey       string
	RecaptchaServerKey     string
	TurnstileSecretKey     string
	CustomCSS              string
	CustomJS               string
	GoogleAnalyticsID      string
	AllowUpload            bool
	AllowNewUser           bool
	ImageCacheMaxBytes     int64
	ImageCacheMaxFileBytes int64
}

func ReadServerConfigFromEnv() HttpServerConfigDef {
	return HttpServerConfigDef{
		Bind:                   os.Getenv("IMGDD_HTTP_BIND"),
		WriteTimeout:           10,
		ReadTimeout:            10,
		SessionKey:             utils.GetEnv("IMGDD_SESSION_KEY", "NOT_SECURE_KEY"),
		RedisURIForSession:     utils.GetEnv("IMGDD_REDIS_URI_FOR_SESSION", "redis://localhost:30102"),
		RedisURI:               utils.GetEnv("IMGDD_REDIS_URI", "redis://localhost:30102"),
		SiteName:               utils.GetEnv("IMGDD_SITE_NAME", ""),
		SiteTitle:              utils.GetEnv("IMGDD_SITE_TITLE", "IMGDD - Image Direct Delivery"),
		ImageDomain:            utils.GetEnv("IMGDD_IMAGE_DOMAIN", ""),
		DefaultURLFormat:       domainmodels.ImageURLFormat(utils.GetEnv("IMGDD_DEFAULT_URL_FORMAT", "canonical")),
		EnableSafeImageCheck:   utils.IsStrTruthy(utils.GetEnv("IMGDD_ENABLE_SAFE_IMAGE_CHECK", "false")),
		SafeImageCheckEndpoint: utils.GetEnv("IMGDD_SAFE_IMAGE_CHECK_ENDPOINT", ""),
		CaptchaProvider:        captcha.CaptchaProvider(utils.GetEnv("IMGDD_CAPTCHA_PROVIDER", "off")),
		RecaptchaClientKey:     utils.GetEnv("IMGDD_RECAPTCHA_CLIENT_KEY", ""),
		TurnstileSiteKey:       utils.GetEnv("IMGDD_TURNSTILE_SITE_KEY", ""),
		RecaptchaServerKey:     utils.GetEnv("IMGDD_RECAPTCHA_SERVER_KEY", ""),
		TurnstileSecretKey:     utils.GetEnv("IMGDD_TURNSTILE_SECRET_KEY", ""),
		CustomCSS:              utils.GetEnv("IMGDD_CUSTOM_CSS", ""),
		CustomJS:               utils.GetEnv("IMGDD_CUSTOM_JS", ""),
		GoogleAnalyticsID:      utils.GetEnv("IMGDD_GOOGLE_ANALYTICS_ID", ""),
		AllowUpload:            utils.IsStrTruthy(utils.GetEnv("IMGDD_ALLOW_UPLOAD", "true")),
		AllowNewUser:           utils.IsStrTruthy(utils.GetEnv("IMGDD_ALLOW_NEW_USER", "true")),
		ImageCacheMaxBytes:     parseInt64Env("IMGDD_HTTP_IMAGE_CACHE_MAX_BYTES"),
		ImageCacheMaxFileBytes: parseInt64Env("IMGDD_HTTP_IMAGE_CACHE_MAX_FILE_BYTES"),
	}
}

func parseInt64Env(key string) int64 {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}
