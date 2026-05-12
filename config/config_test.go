package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ericls/imgdd/captcha"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/storage"
)

func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	return configPath
}

func TestReadFromBytesDoesNotInjectSampleConfig(t *testing.T) {
	conf, err := ReadFromBytes([]byte(""))
	if err != nil {
		t.Fatalf("expected empty config to parse: %v", err)
	}

	if conf.Storage != nil {
		t.Fatal("expected raw empty config not to include sample storage config")
	}
	if conf.HTTPServer != nil {
		t.Fatal("expected raw empty config not to include sample HTTP server config")
	}
	if conf.Email != nil {
		t.Fatal("expected raw empty config not to include sample email config")
	}
}

func TestConfigFromFileUsesRuntimeDefaultsForMissingValues(t *testing.T) {
	configPath := writeTestConfig(t, "")

	conf, err := ConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("expected empty config to parse: %v", err)
	}

	if conf.Storage.StorageDefSource != storage.StorageDefSourceDB {
		t.Fatalf("expected default storage backend source db, got %q", conf.Storage.StorageDefSource)
	}
	if len(conf.Storage.StorageDefs) != 0 {
		t.Fatalf("expected empty config not to include sample storage backends, got %d entries", len(conf.Storage.StorageDefs))
	}
	if conf.HttpServer.DefaultURLFormat != domainmodels.ImageURLFormat_CANONICAL {
		t.Fatalf("expected default URL format canonical, got %q", conf.HttpServer.DefaultURLFormat)
	}
	if conf.Email.Type != email.EmailBackendDummy {
		t.Fatalf("expected default email backend dummy, got %q", conf.Email.Type)
	}
	if conf.HttpServer.CaptchaProvider != captcha.CaptchaProviderOff {
		t.Fatalf("expected default captcha provider off, got %q", conf.HttpServer.CaptchaProvider)
	}
}

func TestConfigFromFileAcceptsEmptyStorageBackendsWithMissingSource(t *testing.T) {
	configPath := writeTestConfig(t, `
[StorageConfig]
STORAGE_BACKENDS = []
`)

	conf, err := ConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("expected config to parse: %v", err)
	}

	if conf.Storage.StorageDefSource != storage.StorageDefSourceDB {
		t.Fatalf("expected default storage backend source db, got %q", conf.Storage.StorageDefSource)
	}
	if len(conf.Storage.StorageDefs) != 0 {
		t.Fatalf("expected explicit empty storage backend list to be preserved, got %d entries", len(conf.Storage.StorageDefs))
	}
}

func TestConfigFromFileAcceptsSparseHTTPAndEmailConfig(t *testing.T) {
	configPath := writeTestConfig(t, `
[HTTPServerConfig]
SITE_NAME = "custom"

[EmailConfig]
`)

	conf, err := ConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("expected sparse config to parse: %v", err)
	}

	if conf.HttpServer.SiteName != "custom" {
		t.Fatalf("expected site name override, got %q", conf.HttpServer.SiteName)
	}
	if conf.HttpServer.DefaultURLFormat != "canonical" {
		t.Fatalf("expected missing default URL format to use canonical, got %q", conf.HttpServer.DefaultURLFormat)
	}
	if conf.Email.Type != "dummy" {
		t.Fatalf("expected missing email type to use dummy, got %q", conf.Email.Type)
	}
}

func TestConfigFromFileReadsImageResponseCacheConfig(t *testing.T) {
	configPath := writeTestConfig(t, `
[HTTPServerConfig]
IMAGE_CACHE_MAX_BYTES = 1048576
IMAGE_CACHE_MAX_FILE_BYTES = 262144
`)

	conf, err := ConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("expected config to parse: %v", err)
	}

	if conf.HttpServer.ImageCacheMaxBytes != 1048576 {
		t.Fatalf("expected image cache max bytes 1048576, got %d", conf.HttpServer.ImageCacheMaxBytes)
	}
	if conf.HttpServer.ImageCacheMaxFileBytes != 262144 {
		t.Fatalf("expected image cache max file bytes 262144, got %d", conf.HttpServer.ImageCacheMaxFileBytes)
	}
}

func TestConfigFromFileDefaultsCleanupIntervalWhenEnabled(t *testing.T) {
	configPath := writeTestConfig(t, `
[CleanupTaskConfig]
ENABLED = true
`)

	conf, err := ConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("expected cleanup config to parse: %v", err)
	}

	if conf.CleanupConfig == nil {
		t.Fatal("expected cleanup config")
	}
	if conf.CleanupConfig.Interval != 600*time.Second {
		t.Fatalf("expected default cleanup interval 600s, got %s", conf.CleanupConfig.Interval)
	}
}

func TestGetConfigDoesNotLetFileRuntimeDefaultsOverrideEnv(t *testing.T) {
	t.Setenv("IMGDD_DEFAULT_URL_FORMAT", "direct")
	t.Setenv("EMAIL_BACKEND_TYPE", "smtp")
	t.Setenv("IMGDD_CAPTCHA_PROVIDER", "recaptcha")
	configPath := writeTestConfig(t, `
[HTTPServerConfig]
SITE_NAME = "custom"

[EmailConfig]
`)

	conf, err := GetConfig(configPath)
	if err != nil {
		t.Fatalf("expected merged config to parse: %v", err)
	}

	if conf.HttpServer.SiteName != "custom" {
		t.Fatalf("expected file override for site name, got %q", conf.HttpServer.SiteName)
	}
	if conf.HttpServer.DefaultURLFormat != domainmodels.ImageURLFormat_DIRECT {
		t.Fatalf("expected env URL format direct to be preserved, got %q", conf.HttpServer.DefaultURLFormat)
	}
	if conf.Email.Type != email.EmailBackendTypeSMTP {
		t.Fatalf("expected env email backend smtp to be preserved, got %q", conf.Email.Type)
	}
	if conf.HttpServer.CaptchaProvider != captcha.CaptchaProviderGoogleRecaptcha {
		t.Fatalf("expected env captcha provider recaptcha to be preserved, got %q", conf.HttpServer.CaptchaProvider)
	}
}
