package test_support

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/logging"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type TestS3Config struct {
	Bucket string
	Access string
	Secret string
	Port   string
}

type TestWebDAVConfig struct {
	Username string
	Password string
	Port     string
}

type TestExternalServiceManager struct {
	Pool           *dockertest.Pool
	dbResource     *dockertest.Resource
	redisResource  *dockertest.Resource
	minioResource  *dockertest.Resource
	webDavResource *dockertest.Resource

	dbConfig     *db.DBConfigDef
	redisURI     string
	s3Config     *TestS3Config
	webDavConfig *TestWebDAVConfig

	logger zerolog.Logger
	lock   sync.Mutex
}

func NewTestExternalServiceManager() *TestExternalServiceManager {
	pool, err := dockertest.NewPool("")
	pool.MaxWait = 10 * time.Second
	if err != nil {
		panic(err)
	}
	err = pool.Client.Ping()
	if err != nil {
		panic(err)
	}
	s3Config := &TestS3Config{
		Bucket: "test-bucket",
		Access: "minio",
		Secret: "minio123",
		Port:   "",
	}
	dbConfig := &db.DBConfigDef{
		POSTGRES_DB:       "imgdd_test",
		POSTGRES_PASSWORD: "imgdd_test",
		POSTGRES_USER:     "imgdd_test",
		POSTGRES_HOST:     "localhost",
		POSTGRES_PORT:     "",
	}
	webDavConfig := &TestWebDAVConfig{
		Username: "test",
		Password: "test",
		Port:     "",
	}
	return &TestExternalServiceManager{
		Pool:          pool,
		dbResource:    nil,
		redisResource: nil,
		minioResource: nil,
		s3Config:      s3Config,
		dbConfig:      dbConfig,
		webDavConfig:  webDavConfig,
		logger:        logging.GetLogger("TESM"),
	}
}

func (ts *TestExternalServiceManager) StartPostgres() {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.logger.Info().Msg("Starting Postgres")
	var err error
	ts.dbResource, err = ts.Pool.Run("postgres", "alpine", ts.dbConfig.EnvLines())
	if err != nil {
		panic(err)
	}
	ts.dbConfig.POSTGRES_PORT = ts.dbResource.GetPort("5432/tcp")
	if err := ts.Pool.Retry(func() error {
		conn := db.GetConnection(ts.dbConfig)
		return conn.Ping()
	}); err != nil {
		ts.Pool.Purge(ts.dbResource)
		panic(err)
	}
	ts.logger.Info().Msg("Running migrations")
	db.RunMigrationUp(ts.GetDBConfig())
	db.PopulateBuiltInRoles(ts.GetDBConfig())
}

func (ts *TestExternalServiceManager) StartRedis() {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.logger.Info().Msg("Starting Redis")
	var err error
	ts.redisResource, err = ts.Pool.Run("redis", "alpine", nil)
	if err != nil {
		panic(err)
	}
	testRedisURI := "redis://" + ts.redisResource.GetHostPort("6379/tcp")
	if err := ts.Pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{
			Addr: strings.TrimPrefix(testRedisURI, "redis://"),
		})
		return client.Ping(context.Background()).Err()
	}); err != nil {
		ts.Pool.Purge(ts.redisResource)
		panic(err)
	}
	ts.redisURI = testRedisURI
}

func (ts *TestExternalServiceManager) StartMinio() {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.logger.Info().Msg("Starting Minio")
	var err error
	minioContainer, err := ts.Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "RELEASE.2021-04-22T15-44-28Z",
		Env: []string{
			"MINIO_ROOT_USER=" + ts.s3Config.Access,
			"MINIO_ROOT_PASSWORD=" + ts.s3Config.Secret,
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9000/tcp": {{HostIP: "0.0.0.0", HostPort: "0"}},
		},
		Cmd: []string{"server", "/data"},
	})
	if err != nil {
		panic(err)
	}
	port := minioContainer.GetPort("9000/tcp")
	ts.minioResource = minioContainer
	ts.s3Config.Port = port
}

func (ts *TestExternalServiceManager) StartWebDav() {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.logger.Info().Msg("Starting WebDav")
	var err error
	webDavContainer, err := ts.Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "bytemark/webdav",
		Tag:        "2.4",
		Env: []string{
			"USERNAME=" + ts.webDavConfig.Username,
			"PASSWORD=" + ts.webDavConfig.Password,
			"AUTH_TYPE=Digest",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"80/tcp": {{HostIP: "0.0.0.0", HostPort: "0"}},
		},
	})
	if err != nil {
		panic(err)
	}
	port := webDavContainer.GetPort("80/tcp")
	ts.webDavResource = webDavContainer
	ts.webDavConfig.Port = port
}

func (ts *TestExternalServiceManager) Purge() {
	if ts.dbResource != nil {
		ts.Pool.Purge(ts.dbResource)
	}
	if ts.redisResource != nil {
		ts.Pool.Purge(ts.redisResource)
	}
	if ts.minioResource != nil {
		ts.Pool.Purge(ts.minioResource)
	}
	if ts.webDavResource != nil {
		ts.Pool.Purge(ts.webDavResource)
	}
}

func (ts *TestExternalServiceManager) GetDBConfig() *db.DBConfigDef {
	return ts.dbConfig
}

func (ts *TestExternalServiceManager) GetRedisURI() string {
	return ts.redisURI
}

func (ts *TestExternalServiceManager) GetS3Config() *TestS3Config {
	return ts.s3Config
}

func (ts *TestExternalServiceManager) GetS3ConfigJSON() string {
	s3Config := ts.GetS3Config()
	config := fmt.Sprintf(`{"endpoint":"http://localhost:%s","bucket":"%s","access":"%s","secret":"%s"}`,
		s3Config.Port, s3Config.Bucket, s3Config.Access, s3Config.Secret,
	)
	return config
}

func (ts *TestExternalServiceManager) GetWebDavConfig() *TestWebDAVConfig {
	return ts.webDavConfig
}
