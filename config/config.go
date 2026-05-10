package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ericls/imgdd/captcha"
	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/httpserver"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/utils"

	dm "github.com/ericls/imgdd/domainmodels"
)

const (
	defaultConfigURLFormat     = dm.ImageURLFormat_CANONICAL
	defaultConfigEmailBackend  = email.EmailBackendDummy
	defaultConfigStorageSource = storage.StorageDefSourceDB
	defaultConfigCaptcha       = captcha.CaptchaProviderOff
	defaultCleanupInterval     = 600 * time.Second
)

type ConfigDef struct {
	Db            db.DBConfigDef
	HttpServer    httpserver.HttpServerConfigDef
	Storage       storage.StorageConfigDef
	Email         email.EmailConfigDef
	CleanupConfig *storage.CleanupConfig
	configFileDef *ConfigFileDef
}

func ConfigFromEnv() (*ConfigDef, error) {
	return &ConfigDef{
		Db:            db.ReadConfigFromEnv(),
		HttpServer:    httpserver.ReadServerConfigFromEnv(),
		Email:         email.ReadEmailConfigFromEnv(),
		CleanupConfig: storage.ReadCleanupConfigFromEnv(),
	}, nil
}

func ConfigFromFile(filePath string) (*ConfigDef, error) {
	if filePath == "" {
		return nil, nil
	}
	configFile, err := ReadFromTomlFile(filePath)
	if err != nil {
		return nil, err
	}

	dbConfig := configFile.DB
	if dbConfig == nil {
		dbConfig = &DBConfigFileDef{}
	}
	redisConfig := configFile.Redis
	if redisConfig == nil {
		redisConfig = &RedisConfigFileDef{}
	}
	httpServerConfig := configFile.HTTPServer
	if httpServerConfig == nil {
		httpServerConfig = &HTTPServerConfigFileDef{}
	}
	storageConfig := configFile.Storage
	if storageConfig == nil {
		storageConfig = &StorageConfigFileDef{}
	}
	emailConfig := configFile.Email
	if emailConfig == nil {
		emailConfig = &EmailConfigFileDef{}
	}

	storageDefs := make([]dm.StorageDefinition, len(storageConfig.STORAGE_BACKENDS))
	for i, storageDef := range storageConfig.STORAGE_BACKENDS {
		storageDefs[i] = dm.StorageDefinition{
			Id:          storageDef.ID,
			Identifier:  storageDef.IDENTIFIER,
			StorageType: dm.StorageTypeName(storageDef.STORAGE_TYPE),
			Config:      storageDef.CONFIG,
			IsEnabled:   storageDef.IS_ENABLED,
			Priority:    storageDef.PRIORITY,
		}
	}
	defaultURLFormat := dm.ImageURLFormat(httpServerConfig.DEFAULT_URL_FORMAT)
	if defaultURLFormat == "" {
		defaultURLFormat = defaultConfigURLFormat
	}
	if !defaultURLFormat.IsValid() {
		return nil, fmt.Errorf("invalid default URL format: %s", httpServerConfig.DEFAULT_URL_FORMAT)
	}

	emailBackendType := email.EmailBackendType(emailConfig.TYPE)
	if emailBackendType == "" {
		emailBackendType = defaultConfigEmailBackend
	}
	if !emailBackendType.IsValid() {
		return nil, fmt.Errorf("invalid email backend type: %s", emailConfig.TYPE)
	}
	storageDefSource := storage.StorageDefSource(storageConfig.STORAGE_BACKEND_SOURCE)
	if storageDefSource == "" {
		storageDefSource = defaultConfigStorageSource
	}
	if storageDefSource != storage.StorageDefSourceDB && storageDefSource != storage.StorageDefSourceConf {
		return nil, fmt.Errorf("invalid storage backend source: %s", storageConfig.STORAGE_BACKEND_SOURCE)
	}
	captchaProvider := captcha.CaptchaProvider(httpServerConfig.CAPTCHA_PROVIDER)
	if captchaProvider == "" {
		captchaProvider = defaultConfigCaptcha
	}
	if !captchaProvider.IsValid() {
		return nil, fmt.Errorf("invalid captcha provider: %s", httpServerConfig.CAPTCHA_PROVIDER)
	}
	SMTPConfigFormFile := emailConfig.SMTP
	var SMTPConfig *email.SMTPConfigDef
	if SMTPConfigFormFile != nil {
		SMTPConfig = &email.SMTPConfigDef{
			Host:     SMTPConfigFormFile.HOST,
			Port:     SMTPConfigFormFile.PORT,
			Username: SMTPConfigFormFile.USERNAME,
			Password: SMTPConfigFormFile.PASSOWRD,
			From:     SMTPConfigFormFile.FROM,
		}
	}
	var cleanupConfig *storage.CleanupConfig
	if configFile.Cleanup != nil {
		cleanupInterval := time.Duration(configFile.Cleanup.INTERVAL) * time.Second
		if configFile.Cleanup.ENABLED && cleanupInterval == 0 {
			cleanupInterval = defaultCleanupInterval
		}
		cleanupConfig = &storage.CleanupConfig{
			Enabled:  configFile.Cleanup.ENABLED,
			Interval: cleanupInterval,
		}
	}

	return &ConfigDef{
		Db: db.DBConfigDef{
			POSTGRES_DB:       dbConfig.POSTGRES_DB,
			POSTGRES_USER:     dbConfig.POSTGRES_USER,
			POSTGRES_PASSWORD: dbConfig.POSTGRES_PASSWORD,
			POSTGRES_HOST:     dbConfig.POSTGRES_HOST,
			POSTGRES_PORT:     dbConfig.POSTGRES_PORT,
			LOG_QUERIES:       dbConfig.LOG_QUERIES,
		},
		HttpServer: httpserver.HttpServerConfigDef{
			Bind:                   httpServerConfig.BIND,
			WriteTimeout:           httpServerConfig.WRITE_TIMEOUT,
			ReadTimeout:            httpServerConfig.READ_TIMEOUT,
			SessionKey:             httpServerConfig.SESSION_KEY,
			RedisURIForSession:     redisConfig.GetSessionRedisURI(),
			RedisURI:               redisConfig.REDIS_URI,
			SiteName:               httpServerConfig.SITE_NAME,
			SiteTitle:              httpServerConfig.SITE_TITLE,
			ImageDomain:            httpServerConfig.IMAGE_DOMAIN,
			DefaultURLFormat:       defaultURLFormat,
			EnableSafeImageCheck:   utils.IsStrTruthy(httpServerConfig.ENABLE_SAFE_IMAGE_CHECK),
			SafeImageCheckEndpoint: httpServerConfig.SAFE_IMAGE_CHECK_ENDPOINT,
			CaptchaProvider:        captchaProvider,
			RecaptchaClientKey:     httpServerConfig.RECAPTCHA_CLIENT_KEY,
			TurnstileSiteKey:       httpServerConfig.TURNSTILE_SITE_KEY,
			RecaptchaServerKey:     httpServerConfig.RECAPTCHA_SERVER_KEY,
			TurnstileSecretKey:     httpServerConfig.TURNSTILE_SECRET_KEY,
			CustomCSS:              httpServerConfig.CUSTOM_CSS,
			CustomJS:               httpServerConfig.CUSTOM_JS,
			GoogleAnalyticsID:      httpServerConfig.GOOGLE_ANALYTICS_ID,
			AllowUpload:            utils.IsStrTruthy(httpServerConfig.ALLOW_UPLOAD),
			AllowNewUser:           utils.IsStrTruthy(httpServerConfig.ALLOW_NEW_USER),
		},
		Storage: storage.StorageConfigDef{
			StorageDefSource: storageDefSource,
			StorageDefs:      storageDefs,
		},
		Email: email.EmailConfigDef{
			Type: emailBackendType,
			SMTP: SMTPConfig,
		},
		CleanupConfig: cleanupConfig,
		configFileDef: configFile,
	}, nil
}

func mergeConfigs(configs ...*ConfigDef) *ConfigDef {
	merged := &ConfigDef{
		Storage: storage.StorageConfigDef{
			StorageDefSource: storage.StorageDefSourceDB,
		},
		HttpServer: configs[0].HttpServer,
	}
	for _, config := range configs {
		if config == nil {
			continue
		}
		fileConfig := config.configFileDef
		var fileHTTPServerConfig *HTTPServerConfigFileDef
		var fileStorageConfig *StorageConfigFileDef
		var fileEmailConfig *EmailConfigFileDef
		if fileConfig != nil {
			fileHTTPServerConfig = fileConfig.HTTPServer
			fileStorageConfig = fileConfig.Storage
			fileEmailConfig = fileConfig.Email
		}
		if config.Db.POSTGRES_DB != "" {
			merged.Db.POSTGRES_DB = config.Db.POSTGRES_DB
		}
		if config.Db.POSTGRES_USER != "" {
			merged.Db.POSTGRES_USER = config.Db.POSTGRES_USER
		}
		if config.Db.POSTGRES_PASSWORD != "" {
			merged.Db.POSTGRES_PASSWORD = config.Db.POSTGRES_PASSWORD
		}
		if config.Db.POSTGRES_HOST != "" {
			merged.Db.POSTGRES_HOST = config.Db.POSTGRES_HOST
		}
		if config.Db.POSTGRES_PORT != "" {
			merged.Db.POSTGRES_PORT = config.Db.POSTGRES_PORT
		}
		if config.Db.LOG_QUERIES != nil {
			merged.Db.LOG_QUERIES = config.Db.LOG_QUERIES
		}
		if config.HttpServer.Bind != "" {
			merged.HttpServer.Bind = config.HttpServer.Bind
		}
		if config.HttpServer.WriteTimeout != 0 {
			merged.HttpServer.WriteTimeout = config.HttpServer.WriteTimeout
		}
		if config.HttpServer.ReadTimeout != 0 {
			merged.HttpServer.ReadTimeout = config.HttpServer.ReadTimeout
		}
		if config.HttpServer.SessionKey != "" {
			merged.HttpServer.SessionKey = config.HttpServer.SessionKey
		}
		if config.HttpServer.RedisURIForSession != "" {
			merged.HttpServer.RedisURIForSession = config.HttpServer.RedisURIForSession
		}
		if config.HttpServer.RedisURI != "" {
			merged.HttpServer.RedisURI = config.HttpServer.RedisURI
		}
		if config.HttpServer.SiteName != "" {
			merged.HttpServer.SiteName = config.HttpServer.SiteName
		}
		if config.HttpServer.SiteTitle != "" {
			merged.HttpServer.SiteTitle = config.HttpServer.SiteTitle
		}
		if config.HttpServer.ImageDomain != "" {
			merged.HttpServer.ImageDomain = config.HttpServer.ImageDomain
		}
		if config.HttpServer.DefaultURLFormat != "" && (fileConfig == nil || fileHTTPServerConfig != nil && fileHTTPServerConfig.DEFAULT_URL_FORMAT != "") {
			merged.HttpServer.DefaultURLFormat = config.HttpServer.DefaultURLFormat
		}
		if fileHTTPServerConfig != nil && fileHTTPServerConfig.ENABLE_SAFE_IMAGE_CHECK != "" {
			merged.HttpServer.EnableSafeImageCheck = config.HttpServer.EnableSafeImageCheck
		}
		if config.HttpServer.SafeImageCheckEndpoint != "" {
			merged.HttpServer.SafeImageCheckEndpoint = config.HttpServer.SafeImageCheckEndpoint
		}
		if config.HttpServer.CaptchaProvider != "" && (fileConfig == nil || fileHTTPServerConfig != nil && fileHTTPServerConfig.CAPTCHA_PROVIDER != "") {
			merged.HttpServer.CaptchaProvider = config.HttpServer.CaptchaProvider
			if config.HttpServer.TurnstileSiteKey != "" {
				merged.HttpServer.TurnstileSiteKey = config.HttpServer.TurnstileSiteKey
			}
			if config.HttpServer.RecaptchaClientKey != "" {
				merged.HttpServer.RecaptchaClientKey = config.HttpServer.RecaptchaClientKey
			}
			if config.HttpServer.RecaptchaServerKey != "" {
				merged.HttpServer.RecaptchaServerKey = config.HttpServer.RecaptchaServerKey
			}
			if config.HttpServer.TurnstileSecretKey != "" {
				merged.HttpServer.TurnstileSecretKey = config.HttpServer.TurnstileSecretKey
			}
		}
		if config.HttpServer.CustomCSS != "" {
			merged.HttpServer.CustomCSS = config.HttpServer.CustomCSS
		}
		if config.HttpServer.CustomJS != "" {
			merged.HttpServer.CustomJS = config.HttpServer.CustomJS
		}
		if config.HttpServer.GoogleAnalyticsID != "" {
			merged.HttpServer.GoogleAnalyticsID = config.HttpServer.GoogleAnalyticsID
		}
		if fileHTTPServerConfig != nil && fileHTTPServerConfig.ALLOW_UPLOAD != "" {
			merged.HttpServer.AllowUpload = config.HttpServer.AllowUpload
		}
		if fileHTTPServerConfig != nil && fileHTTPServerConfig.ALLOW_NEW_USER != "" {
			merged.HttpServer.AllowNewUser = config.HttpServer.AllowNewUser
		}
		if config.Storage.StorageDefSource != "" && (fileConfig == nil || fileStorageConfig != nil && (fileStorageConfig.STORAGE_BACKEND_SOURCE != "" || fileStorageConfig.STORAGE_BACKENDS != nil)) {
			merged.Storage.StorageDefSource = config.Storage.StorageDefSource
		}
		if config.Storage.StorageDefSource == storage.StorageDefSourceConf {
			merged.Storage.StorageDefs = config.Storage.StorageDefs
		}
		if config.Email.Type != "" && (fileConfig == nil || fileEmailConfig != nil && fileEmailConfig.TYPE != "") {
			merged.Email.Type = config.Email.Type
		}
		if config.Email.SMTP != nil && config.Email.SMTP.Host != "" {
			merged.Email.SMTP = config.Email.SMTP
		}
		if config.CleanupConfig != nil {
			merged.CleanupConfig = config.CleanupConfig
		}
	}
	if merged.Storage.StorageDefSource == storage.StorageDefSourceDB {
		merged.Storage.Conn = utils.NewLazy(func() *sql.DB { return db.GetConnection(&merged.Db) })
	}
	return merged
}

func GetConfig(maybeConfigFile string) (*ConfigDef, error) {
	envConf, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}
	fileConf, err := ConfigFromFile(maybeConfigFile)
	return mergeConfigs(envConf, fileConf), err
}

func (c *ConfigDef) PrintConfig() {
	fmt.Printf("%#v\n", c)
}
