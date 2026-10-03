package storage_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ericls/imgdd/storage"
)

func TestFSStoraage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_fs_storage_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

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

func fsConfig(mediaRoot string) []byte {
	return []byte(`{"mediaRoot": "` + mediaRoot + `"}`)
}

func TestFSValidateJSONConfig(t *testing.T) {
	backend := storage.GetBackend("fs")
	tempDir := t.TempDir()

	if err := backend.ValidateJSONConfig(fsConfig(tempDir)); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatal("validation left files behind")
	}

	missing := filepath.Join(tempDir, "missing")
	if err := backend.ValidateJSONConfig(fsConfig(missing)); err == nil {
		t.Fatal("expected error for missing mediaRoot")
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Fatal("validation must not create mediaRoot")
	}

	regularFile := filepath.Join(tempDir, "file")
	if err := os.WriteFile(regularFile, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := backend.ValidateJSONConfig(fsConfig(regularFile)); err == nil {
		t.Fatal("expected error for non-directory mediaRoot")
	}

	for _, root := range []string{"", "relative/media", "./media"} {
		if err := backend.ValidateJSONConfig(fsConfig(root)); err == nil {
			t.Fatalf("expected ValidateJSONConfig error for %q", root)
		}
		if _, err := backend.FromJSONConfig(fsConfig(root)); err == nil {
			t.Fatalf("expected FromJSONConfig error for %q", root)
		}
	}
}

func TestFSValidateJSONConfigNotWritable(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to root")
	}
	tempDir := t.TempDir()
	if err := os.Chmod(tempDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(tempDir, 0755)
	err := storage.GetBackend("fs").ValidateJSONConfig(fsConfig(tempDir))
	if err == nil || !strings.Contains(err.Error(), "not writable") {
		t.Fatalf("expected not writable error, got %v", err)
	}
}

func TestFSStorageConfinedToRoot(t *testing.T) {
	parent := t.TempDir()
	mediaRoot := filepath.Join(parent, "media")
	if err := os.Mkdir(mediaRoot, 0755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(parent, "outside.txt")
	if err := os.WriteFile(outside, []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(mediaRoot, "link.txt")); err != nil {
		t.Fatal(err)
	}

	backend, err := storage.GetBackend("fs").FromJSONConfig(fsConfig(mediaRoot))
	if err != nil {
		t.Fatal(err)
	}

	badNames := []string{"../outside.txt", "..", ".", "", outside, "sub/file.txt", `..\outside.txt`}
	for _, name := range badNames {
		if err := backend.Save(bytes.NewReader([]byte("x")), name, "text/plain"); err == nil {
			t.Errorf("Save(%q) should fail", name)
		}
		if r := backend.GetReader(name); r != nil {
			r.Close()
			t.Errorf("GetReader(%q) should fail", name)
		}
		if meta := backend.GetMeta(name); meta.ByteSize != 0 {
			t.Errorf("GetMeta(%q) should be empty", name)
		}
		if err := backend.Delete(name); err == nil {
			t.Errorf("Delete(%q) should fail", name)
		}
	}

	// symlink escaping the root must not be followed
	if r := backend.GetReader("link.txt"); r != nil {
		r.Close()
		t.Error("GetReader followed symlink outside of mediaRoot")
	}
	if meta := backend.GetMeta("link.txt"); meta.ByteSize != 0 {
		t.Error("GetMeta followed symlink outside of mediaRoot")
	}
	if err := backend.Save(bytes.NewReader([]byte("overwritten")), "link.txt", "text/plain"); err == nil {
		t.Error("Save followed symlink outside of mediaRoot")
	}

	data, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "secret" {
		t.Fatal("file outside of mediaRoot was modified")
	}
}
