package storage_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/ericls/imgdd/storage"
)

func TestFSStoraage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_fs_storage_*")
	if err != nil {
		t.Fatal(err)
	}
	// defer os.RemoveAll(tempDir)

	configJSON := []byte(`{"mediaRoot": "` + tempDir + `"}`)
	backend, err := storage.GetBackend("fs").FromJSONConfig(configJSON)
	if err != nil {
		t.Fatal(err)
	}
	if err := backend.CheckConnection(); err != nil {
		t.Fatal(err)
	}

	// Save the file
	data := []byte("test data")
	r := bytes.NewReader(data)
	err = backend.Save(r, "test.txt", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// Check Meta
	meta := backend.GetMeta("test.txt")
	if meta.ByteSize != int64(len(data)) {
		t.Fatal("file size mismatch")
	}

	// Read the file
	buf := new(bytes.Buffer)
	reader := backend.GetReader("test.txt")
	if reader == nil {
		t.Fatal("file not found")
	}
	_, err = buf.ReadFrom(reader)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Fatal("file content mismatch")
	}

	// Delete the file
	err = backend.Delete("test.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Check Meta
	meta = backend.GetMeta("test.txt")
	if meta.ByteSize != 0 {
		t.Fatal("file not deleted")
	}
}
