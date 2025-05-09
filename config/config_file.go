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
	BIND                      string `toml:"BIND" comment:"HTTP server bind address"`
	WRITE_TIMEOUT             int    `toml:"WRITE_TIMEOUT" comment:"HTTP server write timeout"`
	READ_TIMEOUT              int    `toml:"READ_TIMEOUT" comment:"HTTP server read timeout"`
	SESSION_KEY               string `toml:"SESSION_KEY" comment:"Session key"`
	SITE_NAME                 string `toml:"SITE_NAME" comment:"Site name"`
	SITE_TITLE                string `toml:"SITE_TITLE" comment:"Site title"`
	IMAGE_DOMAIN              string `toml:"IMAGE_DOMAIN" comment:"Image domain"`
	DEFAULT_URL_FORMAT        string `toml:"DEFAULT_URL_FORMAT" comment:"Default URL format. Choices are \n1. 'canonical' - Chooses best backend, and proxies content from that backend. \n2. 'direct' - A backend identifier is included in the URL and directly proxies that storage backend. \n3. 'backend_direct' - URL directly links to the backend (not supported yet)"`
	ENABLE_SAFE_IMAGE_CHECK   string `toml:"ENABLE_SAFE_IMAGE_CHECK" comment:"Enable safe image check. 'true', '1' or 'yes' to enable"`
	SAFE_IMAGE_CHECK_ENDPOINT string `toml:"SAFE_IMAGE_CHECK_ENDPOINT" comment:"Safe image check endpoint. Used if ENABLE_SAFE_IMAGE_CHECK is true"`
	CAPTCHA_PROVIDER          string `toml:"CAPTCHA_PROVIDER" comment:"Captcha provider. Choices are 'off', 'recaptcha', 'turnstile'"`
	RECAPTCHA_CLIENT_KEY      string `toml:"RECAPTCHA_CLIENT_KEY" comment:"Recaptcha client key. Used if CAPTCHA_PROVIDER is 'recaptcha'"`
	TURNSTILE_SITE_KEY        string `toml:"TURNSTILE_SITE_KEY" comment:"Turnstile site key. Used if CAPTCHA_PROVIDER is 'turnstile'"`
	RECAPTCHA_SERVER_KEY      string `toml:"RECAPTCHA_SERVER_KEY" comment:"Recaptcha server key. Used if CAPTCHA_PROVIDER is 'recaptcha'"`
	TURNSTILE_SECRET_KEY      string `toml:"TURNSTILE_SECRET_KEY" comment:"Turnstile secret key. Used if CAPTCHA_PROVIDER is 'turnstile'"`
	CUSTOM_CSS                string `toml:"CUSTOM_CSS" comment:"Custom CSS"`
	CUSTOM_JS                 string `toml:"CUSTOM_JS" comment:"Custom JS"`
	ALLOW_UPLOAD              string `toml:"ALLOW_UPLOAD" comment:"Allow uploading new images. 'true', '1' or 'yes' to enable"`
	ALLOW_NEW_USER            string `toml:"ALLOW_NEW_USER" comment:"Allow new user creation. 'true', '1' or 'yes' to enable"`
}

type StorageBackendItem struct {
	ID           string `toml:"ID" comment:"ID for internal references, must be a valid and unique uuid"`
	IDENTIFIER   string `toml:"IDENTIFIER" comment:"Storage identifier. \nThis value can be exposed to public. \nThis value must be unique. \nThis value is used in the URL format 'direct'"`
	STORAGE_TYPE string `toml:"STORAGE_TYPE" comment:"Storage type"`
	CONFIG       string `toml:"CONFIG" comment:"Storage configuration. \nFormat is dependent on STORAGE_TYPE"`
	IS_ENABLED   bool   `toml:"IS_ENABLED" comment:"Is storage enabled"`
	PRIORITY     int32  `toml:"PRIORITY" comment:"Storage priority. Lower value means higher priority"`
}

type StorageConfigFileDef struct {
	STORAGE_BACKEND_SOURCE string               `toml:"STORAGE_BACKEND_SOURCE" comment:"Storage backend source. \nCan be 'db' or 'conf'. \nIf 'db', the storage backends are read from the database. \nIf 'conf', the storage backends are read from the configuration file."`
	STORAGE_BACKENDS       []StorageBackendItem `toml:"STORAGE_BACKENDS" comment:"Storage backends. \nOnly used if STORAGE_BACKEND_SOURCE is 'conf'"`
}

type SMTPConfigFileDef struct {
	HOST     string `toml:"HOST" comment:"SMTP host"`
	PORT     string `toml:"PORT" comment:"SMTP port"`
	USERNAME string `toml:"USERNAME" comment:"SMTP username"`
	PASSOWRD string `toml:"PASSWORD" comment:"SMTP password"`
	FROM     string `toml:"FROM" comment:"SMTP from"`
}

type EmailConfigFileDef struct {
	TYPE string             `toml:"TYPE" comment:"Email backend. Choices are 'dummy' or 'smtp'"`
	SMTP *SMTPConfigFileDef `toml:"SMTP" comment:"SMTP configuration. Used if TYPE is 'smtp'"`
}

type CleanupTaskConfigFileDef struct {
	ENABLED  bool `toml:"ENABLED" comment:"Enable cleanup task"`
	INTERVAL int  `toml:"INTERVAL" comment:"Cleanup task interval in seconds"`
}

type ConfigFileDef struct {
	DB         *DBConfigFileDef          `toml:"DBConfig" comment:"Database configuration"`
	Redis      *RedisConfigFileDef       `toml:"RedisConfig" comment:"Redis configuration"`
	HTTPServer *HTTPServerConfigFileDef  `toml:"HTTPServerConfig" comment:"HTTP server configuration"`
	Storage    *StorageConfigFileDef     `toml:"StorageConfig" comment:"Storage configuration"`
	Email      *EmailConfigFileDef       `toml:"EmailConfig" comment:"Email configuration"`
	Cleanup    *CleanupTaskConfigFileDef `toml:"CleanupTaskConfig" comment:"Cleanup task configuration"`
}

func (cfd *ConfigFileDef) Clone() ConfigFileDef {
	tomlData, err := toml.Marshal(EmptyConfig)
	if err != nil {
		panic(err)
	}
	var newConfig ConfigFileDef
	err = toml.Unmarshal(tomlData, &newConfig)
	if err != nil {
		panic(err)
	}
	return newConfig
}

var EmptyConfig = ConfigFileDef{
	DB: &DBConfigFileDef{
		LOG_QUERIES: new(bool),
	},
	HTTPServer: &HTTPServerConfigFileDef{
		DEFAULT_URL_FORMAT: "canonical",
	},
	Redis: &RedisConfigFileDef{},
	Storage: &StorageConfigFileDef{
		STORAGE_BACKEND_SOURCE: "db",
		STORAGE_BACKENDS: []StorageBackendItem{
			{
				ID:           "00000000-0000-0000-0000-000000000000",
				IDENTIFIER:   "default",
				STORAGE_TYPE: "s3",
				CONFIG:       `{"endpoint":"http://s3.ca-central-1.amazonaws.com","bucket":"foo","access":"access","secret":"secret!"}`,
				IS_ENABLED:   false,
				PRIORITY:     0,
			},
			{
				ID:           "00000000-0000-0000-0000-000000000001",
				IDENTIFIER:   "default-fs",
				STORAGE_TYPE: "fs",
				CONFIG:       `{"mediaRoot":"/var/www/imgdd/media"}`,
				IS_ENABLED:   false,
				PRIORITY:     0,
			},
			{
				ID:           "00000000-0000-0000-0000-000000000002",
				IDENTIFIER:   "default-webdav",
				STORAGE_TYPE: "webdav",
				CONFIG:       `{"endpoint":"http://webdav.home.arpa:8000","username":"user","password":"password", "pathPrefix":"/foo"}`,
				IS_ENABLED:   false,
				PRIORITY:     0,
			},
		},
	},
	Email: &EmailConfigFileDef{
		TYPE: "dummy",
		SMTP: &SMTPConfigFileDef{
			HOST:     "smtp.home.arpa",
			PORT:     "587",
			USERNAME: "user",
			PASSOWRD: "password",
			FROM:     "noreply@imgdd.com",
		},
	},
	Cleanup: &CleanupTaskConfigFileDef{
		ENABLED:  true,
		INTERVAL: 300,
	},
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
