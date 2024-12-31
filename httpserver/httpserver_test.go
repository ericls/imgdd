package httpserver_test

import (
	"context"
	"imgdd/db"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
)

var TEST_DB_CONF = db.DBConfigDef{
	POSTGRES_DB:       "imgdd_test",
	POSTGRES_PASSWORD: "imgdd_test",
	POSTGRES_USER:     "imgdd_test",
	POSTGRES_HOST:     "localhost",
	POSTGRES_PORT:     "0", // this is set in TestMain
}

var TEST_REDIS_URI = "" // this is set in TestMain

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	pool.MaxWait = 10 * time.Second
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	db_resource, err := pool.Run("postgres", "alpine", TEST_DB_CONF.EnvLines())
	if err != nil {
		log.Fatalf("Could not start db for test: %s", err)
	}
	TEST_DB_CONF.POSTGRES_PORT = db_resource.GetPort("5432/tcp")
	println("Settingup db", TEST_DB_CONF.POSTGRES_PORT, "5432/tcp")
	if err := pool.Retry(func() error {
		conn := db.GetConnection(&TEST_DB_CONF)
		return conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	redis_resource, err := pool.Run("redis", "alpine", nil)
	if err != nil {
		log.Fatalf("Could not start redis for test: %s", err)
	}
	TEST_REDIS_URI = "redis://" + redis_resource.GetHostPort("6379/tcp")
	println("Settingup redis", TEST_REDIS_URI)
	if err := pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{
			Addr: strings.TrimPrefix(TEST_REDIS_URI, "redis://"),
		})
		return client.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	println("Migrating db")
	db.RunMigrationUp(&TEST_DB_CONF)
	db.PopulateBuiltInRoles(&TEST_DB_CONF)

	code := m.Run()
	if err := pool.Purge(db_resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(redis_resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
