package config

import (
	"database/sql"
	"fmt"

	"github.com/ericls/imgdd/captcha"
	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/httpserver"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/utils"

	dm "github.com/ericls/imgdd/domainmodels"
)

type ConfigDef struct {
	Db            db.DBConfigDef
	HttpServer    httpserver.HttpServerConfigDef
	Storage       storage.StorageConfigDef
	Email         email.EmailConfigDef
	configFileDef *ConfigFileDef
}

func ConfigFromEnv() (*ConfigDef, error) {
	return &ConfigDef{
		Db:         db.ReadConfigFromEnv(),
		HttpServer: httpserver.ReadServerConfigFromEnv(),
		Email:      email.ReadEmailConfigFromEnv(),
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
	storageDefs := make([]dm.StorageDefinition, len(configFile.Storage.STORAGE_BACKENDS))
	for i, storageDef := range configFile.Storage.STORAGE_BACKENDS {
		storageDefs[i] = dm.StorageDefinition{
			Id:          storageDef.ID,
			Identifier:  storageDef.IDENTIFIER,
			StorageType: dm.StorageTypeName(storageDef.STORAGE_TYPE),
			Config:      storageDef.CONFIG,
			IsEnabled:   storageDef.IS_ENABLED,
			Priority:    storageDef.PRIORITY,
		}
	}
	defaultURLFormat := dm.ImageURLFormat(configFile.HTTPServer.DEFAULT_URL_FORMAT)
	if !defaultURLFormat.IsValid() {
		return nil, fmt.Errorf("invalid default URL format: %s", configFile.HTTPServer.DEFAULT_URL_FORMAT)
	}

	emailBackendType := email.EmailBackendType(configFile.Email.TYPE)
	if !emailBackendType.IsValid() {
		return nil, fmt.Errorf("invalid email backend type: %s", configFile.Email.TYPE)
	}
	SMTPConfigFormFile := configFile.Email.SMTP
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
	return &ConfigDef{
		Db: db.DBConfigDef{
			POSTGRES_DB:       configFile.DB.POSTGRES_DB,
			POSTGRES_USER:     configFile.DB.POSTGRES_USER,
			POSTGRES_PASSWORD: configFile.DB.POSTGRES_PASSWORD,
			POSTGRES_HOST:     configFile.DB.POSTGRES_HOST,
			POSTGRES_PORT:     configFile.DB.POSTGRES_PORT,
			LOG_QUERIES:       configFile.DB.LOG_QUERIES,
		},
		HttpServer: httpserver.HttpServerConfigDef{
			Bind:                   configFile.HTTPServer.BIND,
			WriteTimeout:           configFile.HTTPServer.WRITE_TIMEOUT,
			ReadTimeout:            configFile.HTTPServer.READ_TIMEOUT,
			SessionKey:             configFile.HTTPServer.SESSION_KEY,
			RedisURIForSession:     configFile.Redis.GetSessionRedisURI(),
			SiteName:               configFile.HTTPServer.SITE_NAME,
			SiteTitle:              configFile.HTTPServer.SITE_TITLE,
			ImageDomain:            configFile.HTTPServer.IMAGE_DOMAIN,
			DefaultURLFormat:       defaultURLFormat,
			EnableSafeImageCheck:   utils.IsStrTruthy(configFile.HTTPServer.ENABLE_SAFE_IMAGE_CHECK),
			SafeImageCheckEndpoint: configFile.HTTPServer.SAFE_IMAGE_CHECK_ENDPOINT,
			CaptchaProvider:        captcha.CaptchaProvider(configFile.HTTPServer.CAPTCHA_PROVIDER),
			RecaptchaClientKey:     configFile.HTTPServer.RECAPTCHA_CLIENT_KEY,
			TurnstileSiteKey:       configFile.HTTPServer.TURNSTILE_SITE_KEY,
			RecaptchaServerKey:     configFile.HTTPServer.RECAPTCHA_SERVER_KEY,
			TurnstileSecretKey:     configFile.HTTPServer.TURNSTILE_SECRET_KEY,
			CustomCSS:              configFile.HTTPServer.CUSTOM_CSS,
			CustomJS:               configFile.HTTPServer.CUSTOM_JS,
			AllowUpload:            utils.IsStrTruthy(configFile.HTTPServer.ALLOW_UPLOAD),
			AllowNewUser:           utils.IsStrTruthy(configFile.HTTPServer.ALLOW_NEW_USER),
		},
		Storage: storage.StorageConfigDef{
			StorageDefSource: storage.StorageDefSource(configFile.Storage.STORAGE_BACKEND_SOURCE),
			StorageDefs:      storageDefs,
		},
		Email: email.EmailConfigDef{
			Type: emailBackendType,
			SMTP: SMTPConfig,
		},
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
		if config.HttpServer.SiteName != "" {
			merged.HttpServer.SiteName = config.HttpServer.SiteName
		}
		if config.HttpServer.SiteTitle != "" {
			merged.HttpServer.SiteTitle = config.HttpServer.SiteTitle
		}
		if config.HttpServer.ImageDomain != "" {
			merged.HttpServer.ImageDomain = config.HttpServer.ImageDomain
		}
		if config.HttpServer.DefaultURLFormat != "" {
			merged.HttpServer.DefaultURLFormat = config.HttpServer.DefaultURLFormat
		}
		if config.configFileDef != nil && config.configFileDef.HTTPServer.ENABLE_SAFE_IMAGE_CHECK != "" {
			merged.HttpServer.EnableSafeImageCheck = config.HttpServer.EnableSafeImageCheck
		}
		if config.HttpServer.SafeImageCheckEndpoint != "" {
			merged.HttpServer.SafeImageCheckEndpoint = config.HttpServer.SafeImageCheckEndpoint
		}
		if config.HttpServer.CaptchaProvider != "" {
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
		if config.configFileDef != nil && config.configFileDef.HTTPServer.ALLOW_UPLOAD != "" {
			merged.HttpServer.AllowUpload = config.HttpServer.AllowUpload
		}
		if config.configFileDef != nil && config.configFileDef.HTTPServer.ALLOW_NEW_USER != "" {
			merged.HttpServer.AllowNewUser = config.HttpServer.AllowNewUser
		}
		if config.Storage.StorageDefSource != "" {
			merged.Storage.StorageDefSource = config.Storage.StorageDefSource
		}
		if config.Storage.StorageDefSource == storage.StorageDefSourceConf {
			merged.Storage.StorageDefs = config.Storage.StorageDefs
		}
		if config.Email.Type != "" {
			merged.Email.Type = config.Email.Type
		}
		if config.Email.SMTP != nil && config.Email.SMTP.Host != "" {
			merged.Email.SMTP = config.Email.SMTP
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
