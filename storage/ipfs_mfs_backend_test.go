package storage_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/ericls/imgdd/storage"
)

func TestIPFSMFSStorage(t *testing.T) {
	TestServiceMan.StartIPFS()
	defer TestServiceMan.StopIPFS()

	data := []byte("test data")
	cfg := TestServiceMan.GetIPFSMFSConfig()

	backend := storage.GetBackend("ipfs_mfs").(*storage.IPFSMFSStorageBackend)
	if backend == nil {
		t.Fatal("ipfs_mfs backend not found")
	}

	storeWithPrefix, err := backend.FromJSONConfig(fmt.Appendf(nil,
		`{"apiUrl":"http://localhost:%s","pathPrefix":"/imgdd-test","pin":true}`, cfg.ApiPort,
	))
	if err != nil {
		t.Fatal(err)
	}
	storeDefault, err := backend.FromJSONConfig(fmt.Appendf(nil,
		`{"apiUrl":"http://localhost:%s","pin":false}`, cfg.ApiPort,
	))
	if err != nil {
		t.Fatal(err)
	}

	if err := storeWithPrefix.CheckConnection(); err != nil {
		t.Fatal(err)
	}

	testStore := func(label string, store storage.Storage) {
		t.Helper()
		if err := store.Save(bytes.NewReader(data), "test.txt", "text/plain"); err != nil {
			t.Fatalf("%s: save: %v", label, err)
		}

		meta := store.GetMeta("test.txt")
		if meta.ByteSize != int64(len(data)) {
			t.Fatalf("%s: size mismatch: got %d, want %d", label, meta.ByteSize, len(data))
		}
		if meta.ETag == "" || meta.ETag == "\"\"" {
			t.Fatalf("%s: expected non-empty ETag (CID)", label)
		}
		if meta.ContentType != "text/plain; charset=utf-8" {
			t.Fatalf("%s: unexpected content-type: %q", label, meta.ContentType)
		}

		reader := store.GetReader("test.txt")
		if reader == nil {
			t.Fatalf("%s: reader nil", label)
		}
		got, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			t.Fatalf("%s: read: %v", label, err)
		}
		if !bytes.Equal(got, data) {
			t.Fatalf("%s: content mismatch: got %q, want %q", label, got, data)
		}

		if err := store.Delete("test.txt"); err != nil {
			t.Fatalf("%s: delete: %v", label, err)
		}

		meta = store.GetMeta("test.txt")
		if meta.ByteSize != 0 {
			t.Fatalf("%s: file not deleted, size=%d", label, meta.ByteSize)
		}
	}

	testStore("with-prefix+pin", storeWithPrefix)
	testStore("default+nopin", storeDefault)
}
