package storage_test

import (
	"imgdd/db"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
)

var TEST_DB_CONF = db.DBConfigDef{
	POSTGRES_DB:       "imgdd_test",
	POSTGRES_PASSWORD: "imgdd_test",
	POSTGRES_USER:     "imgdd_test",
	POSTGRES_HOST:     "localhost",
	POSTGRES_PORT:     "0", // this is set in TestMain
}

var testS3Bucket = "test-bucket"
var testS3Access = "minio"
var testS3Secret = "minio123"
var testS3Port = "" // this is set in TestMain

var dockerTestPool *dockertest.Pool

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	dockerTestPool = pool
	pool.MaxWait = 10 * time.Second
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	db_container, err := pool.Run("postgres", "alpine", TEST_DB_CONF.EnvLines())
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	TEST_DB_CONF.POSTGRES_PORT = db_container.GetPort("5432/tcp")

	// TEST_DB_CONF.SetLogQueries()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		conn := db.GetConnection(&TEST_DB_CONF)
		return conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	db.RunMigrationUp(TEST_DB_CONF)
	db.PopulateBuiltInRoles(TEST_DB_CONF)

	minio_container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "RELEASE.2021-04-22T15-44-28Z",
		Env: []string{
			"MINIO_ROOT_USER=" + testS3Access,
			"MINIO_ROOT_PASSWORD=" + testS3Secret,
		},
		ExposedPorts: []string{"9000"},
		Cmd:          []string{"server", "/data"},
	})
	if err != nil {
		log.Fatalf("Could not start minio: %s", err)
	}
	port := minio_container.GetPort("9000/tcp")
	testS3Port = port

	code := m.Run()
	if err := pool.Purge(db_container); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(minio_container); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
