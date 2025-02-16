package storage_test

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/ericls/imgdd/storage"
)

func TestS3Storage(t *testing.T) {
	TestServiceMan.StartMinio()
	defer TestServiceMan.StopMinio()
	data := []byte("test data")

	s3Config := TestServiceMan.GetS3Config()
	conf := storage.S3StorageConfig{
		Endpoint: "http://localhost:" + s3Config.Port,
		Bucket:   s3Config.Bucket,
		Access:   s3Config.Access,
		Secret:   s3Config.Secret,
	}

	// Create a new S3 storage backend
	backend := storage.GetBackend("s3").(*storage.S3StorageBackend)
	if backend == nil {
		t.Fatal("s3 backend not found")
	}

	// Create a new S3 storage
	store, err := backend.FromJSONConfig(conf.ToJSON())
	s3Storage := store.(*storage.S3Storage)
	if err != nil {
		t.Fatal(err)
	}

	if err := TestServiceMan.Pool.Retry(func() error {
		return s3Storage.CheckConnection()
	}); err != nil {
		log.Fatalf("Could not connect to minio: %s", err)
	}

	// create bucket
	err = s3Storage.CreateBucket(s3Config.Bucket)
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
