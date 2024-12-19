package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type DBConfigFileDef struct {
	POSTGRES_DB       string `toml:"POSTGRES_DB" comment:"Postgres database name"`
	POSTGRES_PASSWORD string `toml:"POSTGRES_PASSWORD" comment:"Postgres password"`
	POSTGRES_USER     string `toml:"POSTGRES_USER" comment:"Postgres user"`
	POSTGRES_HOST     string `toml:"POSTGRES_HOST" comment:"Postgres host"`
	POSTGRES_PORT     string `toml:"POSTGRES_PORT" comment:"Postgres port"`
	LOG_QUERIES       bool   `toml:"LOG_QUERIES" comment:"Log queries, used for debugging"`
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

func resolveFilePath(userInput string, checkExist bool) (string, error) {
	expandedPath := userInput
	if userInput[:1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		expandedPath = filepath.Join(homeDir, userInput[1:])
	}
	absolutePath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	if checkExist {
		_, err = os.Stat(absolutePath)
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file does not exist: %s", absolutePath)
		} else if err != nil {
			return "", fmt.Errorf("error checking file existence: %w", err)
		}
	}

	return absolutePath, nil
}

func ReadFromTomlFile(filePath string) (*ConfigFileDef, error) {
	resolvedPath, err := resolveFilePath(filePath, false)
	println(resolvedPath, "here")
	return nil, err
}

func ReadFromBytes(data []byte) (*ConfigFileDef, error) {
	var x ConfigFileDef
	toml.Unmarshal(data, &x)
	return nil, nil
}
