package storage_test

import (
	"bytes"
	"testing"

	"github.com/ericls/imgdd/storage"
)

func TestWebDAVStorage(t *testing.T) {
	TestServiceMan.StartWebDav()
	data := []byte("test data")
	config := TestServiceMan.GetWebDavConfig()
	backend := storage.GetBackend("webdav").(*storage.WebDAVBackend)
	if backend == nil {
		t.Fatal("webdav backend not found")
	}
	storeWithPrefix, err := backend.FromJSONConfig([]byte(`{"url":"` + "http://localhost:" + config.Port + `","username":"` + config.Username + `","password":"` + config.Password + `", "pathPrefix":"/foo"}`))
	if err != nil {
		t.Fatal(err)
	}
	if storeWithPrefix == nil {
		t.Fatal("store is nil")
	}
	storeWithoutPrefix, err := backend.FromJSONConfig([]byte(`{"url":"` + "http://localhost:" + config.Port + `","username":"` + config.Username + `","password":"` + config.Password + `"}`))
	testStore := func(store storage.Storage) {
		err = store.Save(bytes.NewReader(data), "test.txt", "text/plain")
		if err != nil {
			t.Fatal(err)
		}
		meta := store.GetMeta("test.txt")
		if meta.ByteSize != int64(len(data)) {
			t.Fatal("file size mismatch")
		}
		if meta.ContentType != "text/plain" {
			t.Fatal("file content type mismatch")
		}
		if meta.ETag == "" {
			t.Fatal("file ETag is empty")
		}
		reader := store.GetReader("test.txt")
		if reader == nil {
			t.Fatal("file not found")
		}
		content := make([]byte, len(data))
		reader.Read(content)
		if !bytes.Equal(content, data) {
			t.Fatal("file content mismatch")
		}
		err = store.Delete("test.txt")
		if err != nil {
			t.Fatal(err)
		}
		meta = store.GetMeta("test.txt")
		if meta.ByteSize != 0 {
			t.Fatal("file not deleted")
		}
	}
	testStore(storeWithPrefix)
	testStore(storeWithoutPrefix)
	TestServiceMan.StopWebDav()
}
