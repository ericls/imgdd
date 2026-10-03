package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ericls/imgdd/utils"
)

type FSStorageConfig struct {
	MediaRoot string `json:"mediaRoot"`
}

func checkMediaRoot(mediaRoot string) (string, error) {
	if mediaRoot == "" {
		return "", errors.New("mediaRoot is required")
	}
	if !filepath.IsAbs(mediaRoot) {
		return "", fmt.Errorf("mediaRoot %q must be an absolute path", mediaRoot)
	}
	return filepath.Clean(mediaRoot), nil
}

func parseFSStorageConfig(config []byte) (string, error) {
	var conf FSStorageConfig
	if err := json.Unmarshal(config, &conf); err != nil {
		return "", err
	}
	return checkMediaRoot(conf.MediaRoot)
}

// checkFilename only allows a single plain path component, file identifiers are never nested.
func checkFilename(filename string) error {
	if filename == "" || filename == "." || filename == ".." || strings.ContainsAny(filename, "/\\\x00") {
		return fmt.Errorf("invalid filename %q", filename)
	}
	return nil
}

type FSStorageBackend struct {
}

func (s *FSStorageBackend) FromJSONConfig(config []byte) (Storage, error) {
	rootPath, err := parseFSStorageConfig(config)
	if err != nil {
		return nil, err
	}
	return &FSStorage{
		rootPath: rootPath,
	}, nil
}

func (s *FSStorageBackend) ValidateJSONConfig(config []byte) error {
	rootPath, err := parseFSStorageConfig(config)
	if err != nil {
		return err
	}
	stat, err := os.Stat(rootPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("mediaRoot %s does not exist, it must be created by the operator", rootPath)
		}
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("mediaRoot %s is not a directory", rootPath)
	}
	storage := &FSStorage{rootPath: rootPath}
	return storage.withRoot(func(root *os.Root) error {
		suffix := make([]byte, 8)
		if _, err := rand.Read(suffix); err != nil {
			return err
		}
		probe := ".imgdd-write-check-" + hex.EncodeToString(suffix)
		f, err := root.OpenFile(probe, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("mediaRoot %s is not writable: %w", rootPath, err)
		}
		f.Close()
		return root.Remove(probe)
	})
}

// FSStorage confines all file access to rootPath via os.Root, which rejects
// ".." components, absolute paths and symlinks escaping the root.
// The root is opened per operation so no file descriptor outlives the call.
type FSStorage struct {
	rootPath string
}

func (s *FSStorage) withRoot(fn func(root *os.Root) error) error {
	root, err := os.OpenRoot(s.rootPath)
	if err != nil {
		return err
	}
	defer root.Close()
	return fn(root)
}

func (s *FSStorage) GetReader(filename string) io.ReadCloser {
	if checkFilename(filename) != nil {
		return nil
	}
	f, err := os.OpenInRoot(s.rootPath, filename)
	if err != nil {
		return nil
	}
	return f
}

func (s *FSStorage) Save(file utils.SeekerReader, filename string, mimeType string) error {
	if err := checkFilename(filename); err != nil {
		return err
	}
	return s.withRoot(func(root *os.Root) error {
		f, err := root.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		return err
	})
}

func (s *FSStorage) GetMeta(filename string) FileMeta {
	var stat os.FileInfo
	err := checkFilename(filename)
	if err == nil {
		err = s.withRoot(func(root *os.Root) error {
			var statErr error
			stat, statErr = root.Stat(filename)
			return statErr
		})
	}
	if err != nil {
		return FileMeta{
			ByteSize:    0,
			ContentType: "",
		}
	}
	if stat.IsDir() {
		return FileMeta{
			ByteSize: 0,
		}
	}
	filenameParts := strings.Split(stat.Name(), ".")
	ext := filenameParts[len(filenameParts)-1]
	if ext == "" {
		return FileMeta{
			ByteSize:    stat.Size(),
			ContentType: "",
			ETag:        stat.ModTime().String() + stat.Name(),
		}
	}
	mimeType := mime.TypeByExtension("." + ext)
	return FileMeta{
		ByteSize:    stat.Size(),
		ContentType: mimeType,
		ETag:        stat.ModTime().String() + stat.Name(),
	}
}

func (s *FSStorage) Delete(filename string) error {
	if err := checkFilename(filename); err != nil {
		return err
	}
	return s.withRoot(func(root *os.Root) error {
		return root.Remove(filename)
	})
}

func (s *FSStorage) CheckConnection() error {
	stat, err := os.Stat(s.rootPath)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("mediaRoot %s is not a directory", s.rootPath)
	}
	return nil
}
