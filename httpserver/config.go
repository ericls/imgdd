package httpserver

import (
	"io/fs"
	"os"
)

type HttpServerConfigDef struct {
	Bind               string
	WriteTimeout       int
	ReadTimeout        int
	TemplatesFS        fs.FS
	StaticFS           fs.FS
	SessionKey         string
	RedisURIForSession string
	SiteName           string
	ImageDomain        string
}

var Config HttpServerConfigDef

func init() {
	Config = readServerConfigFromEnv()
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
func readServerConfigFromEnv() HttpServerConfigDef {
	return HttpServerConfigDef{
		Bind:               os.Getenv("imgdd_HTTP_BIND"),
		WriteTimeout:       10,
		ReadTimeout:        10,
		SessionKey:         getenv("IMGDD_SESSION_KEY", "NOT_SECURE_KEY"),
		RedisURIForSession: getenv("IMGDD_REDIS_URI_FOR_SESSION", "redis://localhost:30102"),
		SiteName:           getenv("IMGDD_SITE_NAME", "imgdd"),
		ImageDomain:        getenv("IMGDD_IMAGE_DOMAIN", ""),
	}
}
