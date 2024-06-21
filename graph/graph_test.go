package graph_test

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

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "alpine", TEST_DB_CONF.EnvLines())
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	TEST_DB_CONF.POSTGRES_PORT = resource.GetPort("5432/tcp")
	println("Settingup db", TEST_DB_CONF.POSTGRES_PORT, "5432/tcp")

	// TEST_DB_CONF.SetLogQueries()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		conn := db.GetConnection(&TEST_DB_CONF)
		return conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	println("Migrating db")
	db.RunMigrationUp(&TEST_DB_CONF)
	db.PopulateBuiltInRoles(&TEST_DB_CONF)

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
