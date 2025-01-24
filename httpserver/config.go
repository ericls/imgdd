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
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
func ReadServerConfigFromEnv() HttpServerConfigDef {
	return HttpServerConfigDef{
		Bind:                   os.Getenv("imgdd_HTTP_BIND"),
		WriteTimeout:           10,
		ReadTimeout:            10,
		SessionKey:             getenv("IMGDD_SESSION_KEY", "NOT_SECURE_KEY"),
		RedisURIForSession:     getenv("IMGDD_REDIS_URI_FOR_SESSION", "redis://localhost:30102"),
		SiteName:               getenv("IMGDD_SITE_NAME", "imgdd"),
		ImageDomain:            getenv("IMGDD_IMAGE_DOMAIN", ""),
		DefaultURLFormat:       domainmodels.ImageURLFormat(getenv("IMGDD_DEFAULT_URL_FORMAT", "canonical")),
		EnableSafeImageCheck:   utils.IsStrTruthy(getenv("IMGDD_ENABLE_SAFE_IMAGE_CHECK", "false")),
		SafeImageCheckEndpoint: getenv("IMGDD_SAFE_IMAGE_CHECK_ENDPOINT", ""),
	}
}
