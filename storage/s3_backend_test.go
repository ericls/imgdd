package storage_test

import (
	"bytes"
	"imgdd/storage"
	"io"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
)

var testBucket = "test-bucket"
var testAccess = "minio"
var testSecret = "minio123"
var data = []byte("test data")

func TestS3Storage(t *testing.T) {

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
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "RELEASE.2021-04-22T15-44-28Z",
		Env: []string{
			"MINIO_ROOT_USER=" + testAccess,
			"MINIO_ROOT_PASSWORD=" + testSecret,
		},
		ExposedPorts: []string{"9000"},
		Cmd:          []string{"server", "/data"},
	})
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	port := resource.GetPort("9000/tcp")

	conf := storage.S3StorageConfig{
		Endpoint: "http://localhost:" + port,
		Bucket:   testBucket,
		Access:   testAccess,
		Secret:   testSecret,
	}

	// Create a new S3 storage backend
	backend := storage.GetBackend("s3")
	if backend == nil {
		t.Fatal("s3 backend not found")
	}

	// Create a new S3 storage
	s3Storage := backend.FromJSON(conf.ToJSON()).(*storage.S3Storage)

	if err := pool.Retry(func() error {
		return s3Storage.CheckConnection()
	}); err != nil {
		log.Fatalf("Could not connect to minio: %s", err)
	}

	// create bucket
	err = s3Storage.CreateBucket(testBucket)
	if err != nil {
		t.Fatal(err)
	}

	// Save the file
	r := bytes.NewReader(data)
	err = s3Storage.Save(r, "test.txt", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// Check Meta
	meta := s3Storage.GetMeta("test.txt")
	if meta.ByteSize != int64(len(data)) {
		t.Fatal("file size mismatch")
	}

	// Get the file
	reader := s3Storage.GetReader("test.txt")
	if reader == nil {
		t.Fatal("file not found")
	}

	// Read the file
	data, err = io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	// Check the file content
	if string(data) != "test data" {
		t.Fatal("file content mismatch")
	}

	// Delete the file
	err = s3Storage.Delete("test.txt")
	if err != nil {
		t.Fatal(err)
	}
}
