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
	LOG_QUERIES       *bool  `toml:"LOG_QUERIES" comment:"Log queries, used for debugging"`
}

type RedisConfigFileDef struct {
	REDIS_URI         string `toml:"REDIS_URI" comment:"Redis URI"`
	CACHE_REDIS_URI   string `toml:"CACHE_REDIS_URI" comment:"Redis URI for caching, if different from main Redis"`
	SESSION_REDIS_URI string `toml:"SESSION_REDIS_URI" comment:"Redis URI for session storage, if different from main Redis"`
}

func (r *RedisConfigFileDef) GetSessionRedisURI() string {
	if r.SESSION_REDIS_URI != "" {
		return r.SESSION_REDIS_URI
	}
	return r.REDIS_URI
}

func (r *RedisConfigFileDef) GetCacheRedisURI() string {
	if r.CACHE_REDIS_URI != "" {
		return r.CACHE_REDIS_URI
	}
	return r.REDIS_URI
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
	Redis      *RedisConfigFileDef      `toml:"RedisConfig" comment:"Redis configuration"`
}

var EmptyConfig = ConfigFileDef{
	DB: &DBConfigFileDef{
		LOG_QUERIES: new(bool),
	},
	HTTPServer: &HTTPServerConfigFileDef{},
	Redis:      &RedisConfigFileDef{},
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
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config file path: %w", err)
	}
	file, err := os.Open(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	bytes, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return ReadFromBytes(bytes)
}

func ReadFromBytes(bytes []byte) (*ConfigFileDef, error) {
	var conf ConfigFileDef
	err := toml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &conf, nil
}

func GenerateEmptyConfigFile(filePath string) error {
	resolvedPath, err := resolveFilePath(filePath, false)
	if err != nil {
		return fmt.Errorf("failed to resolve config file path: %w", err)
	}
	// check the file is empty or the file does not exist
	file, err := os.OpenFile(resolvedPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
	}
	if stat.Size() > 0 {
		return fmt.Errorf("file is not empty: %s", resolvedPath)
	}
	tomlData, err := toml.Marshal(EmptyConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal empty config: %w", err)
	}
	// write empty config to file
	_, err = file.Write(tomlData)
	if err != nil {
		return fmt.Errorf("failed to write empty config to file: %w", err)
	}
	return nil
}

func PrintEmptyConfig() error {
	tomlData, err := toml.Marshal(EmptyConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal empty config: %w", err)
	}
	_, err = fmt.Print(string(tomlData))
	if err != nil {
		return fmt.Errorf("failed to print empty config: %w", err)
	}
	return nil
}
