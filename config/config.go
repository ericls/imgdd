package config

import (
	"database/sql"
	"fmt"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/httpserver"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/utils"

	dm "github.com/ericls/imgdd/domainmodels"
)

type ConfigDef struct {
	Db         db.DBConfigDef
	HttpServer httpserver.HttpServerConfigDef
	Storage    storage.StorageConfigDef
}

func ConfigFromEnv() (*ConfigDef, error) {
	return &ConfigDef{
		Db:         db.ReadConfigFromEnv(),
		HttpServer: httpserver.ReadServerConfigFromEnv(),
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
			StorageType: storageDef.STORAGE_TYPE,
			Config:      storageDef.CONFIG,
			IsEnabled:   storageDef.IS_ENABLED,
			Priority:    storageDef.PRIORITY,
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
			Bind:               configFile.HTTPServer.BIND,
			WriteTimeout:       configFile.HTTPServer.WRITE_TIMEOUT,
			ReadTimeout:        configFile.HTTPServer.READ_TIMEOUT,
			SessionKey:         configFile.HTTPServer.SESSION_KEY,
			RedisURIForSession: configFile.Redis.GetSessionRedisURI(),
			SiteName:           configFile.HTTPServer.SITE_NAME,
		},
		Storage: storage.StorageConfigDef{
			StorageDefSource: storage.StorageDefSource(configFile.Storage.STORAGE_BACKEND_SOURCE),
			StorageDefs:      storageDefs,
		},
	}, nil
}

func mergeConfigs(configs ...*ConfigDef) *ConfigDef {
	merged := &ConfigDef{
		Storage: storage.StorageConfigDef{
			StorageDefSource: storage.StorageDefSourceDB,
		},
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
		if config.HttpServer.ImageDomain != "" {
			merged.HttpServer.ImageDomain = config.HttpServer.ImageDomain
		}
		if config.HttpServer.DefaultURLFormat != "" {
			merged.HttpServer.DefaultURLFormat = config.HttpServer.DefaultURLFormat
		}
		if config.Storage.StorageDefSource != "" {
			merged.Storage.StorageDefSource = config.Storage.StorageDefSource
		}
		if config.Storage.StorageDefSource == storage.StorageDefSourceConf {
			merged.Storage.StorageDefs = config.Storage.StorageDefs
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
	fmt.Printf("%#v", c)
}
