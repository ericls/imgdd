package config

import (
	"fmt"
	"imgdd/db"
	"imgdd/httpserver"
)

type ConfigDef struct {
	Db         db.DBConfigDef
	HttpServer httpserver.HttpServerConfigDef
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
			Bind:               configFile.HTTPServer.Bind,
			WriteTimeout:       configFile.HTTPServer.WriteTimeout,
			ReadTimeout:        configFile.HTTPServer.ReadTimeout,
			SessionKey:         configFile.HTTPServer.SessionKey,
			RedisURIForSession: configFile.Redis.GetCacheRedisURI(),
			SiteName:           configFile.HTTPServer.SiteName,
		},
	}, nil
}

func mergeConfigs(configs ...*ConfigDef) *ConfigDef {
	merged := &ConfigDef{}
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
