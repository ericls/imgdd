package config

import "github.com/pelletier/go-toml/v2"

type DBConfigFileDef struct {
	POSTGRES_DB       string `toml:"POSTGRES_DB" comment:"Postgres database name"`
	POSTGRES_PASSWORD string `toml:"POSTGRES_PASSWORD" comment:"Postgres password"`
	POSTGRES_USER     string `toml:"POSTGRES_USER" comment:"Postgres user"`
	POSTGRES_HOST     string `toml:"POSTGRES_HOST" comment:"Postgres host"`
	POSTGRES_PORT     string `toml:"POSTGRES_PORT" comment:"Postgres port"`
	LOG_QUERIES       bool   `toml:"LOG_QUERIES" comment:"Log queries used for debugging"`
}

type RedisConfigFileDef struct {
	REDIS_URI         string `toml:"REDIS_URI" comment:"Redis URI"`
	CACHE_REDIS_URI   string `toml:"CACHE_REDIS_URI" comment:"Redis URI for caching, if different from main Redis"`
	SESSION_REDIS_URI string `toml:"SESSION_REDIS_URI" comment:"Redis URI for session storage, if different from main Redis"`
}

type HTTPServerConfigFileDef struct {
	Bind         string `toml:"Bind" comment:"HTTP server bind address"`
	WriteTimeout int    `toml:"WriteTimeout" comment:"HTTP server write timeout"`
	ReadTimeout  int    `toml:"ReadTimeout" comment:"HTTP server read timeout"`
	SessionKey   string `toml:"SessionKey" comment:"Session key"`
	SiteName     string `toml:"SiteName" comment:"Site name"`
	ImageDomain  string `toml:"ImageDomain" comment:"Image domain"`
}

type ConfigFileDef struct {
	DB         *DBConfigFileDef         `toml:"DBConfig" comment:"Database configuration"`
	HTTPServer *HTTPServerConfigFileDef `toml:"HTTPServerConfig" comment:"HTTP server configuration"`
}

func ReadFromTomlFile(filename string) (*ConfigFileDef, error) {
	var x ConfigFileDef
	data := []byte(`value  = "42"`)
	toml.Unmarshal(data, &x)
	return nil, nil
}
