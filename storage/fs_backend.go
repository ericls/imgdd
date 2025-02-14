package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"strings"

	"github.com/ericls/imgdd/utils"
)

type RootFS interface {
	fs.ReadFileFS
	fs.ReadDirFS
	fs.StatFS
}

type WritableRootFs struct {
	RootFS
	rootPath string
}

func (w *WritableRootFs) getWriter(filename string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	fullPath := w.rootPath + "/" + filename
	fmt.Fprintln(os.Stderr, "fullPath: ", fullPath)
	return os.OpenFile(fullPath, flag, perm)
}

func newRootFs(rootPath string) (*WritableRootFs, error) {
	rootFS, ok := os.DirFS(rootPath).(RootFS)
	if !ok {
		return nil, fmt.Errorf("mediaRoot %s is not a valid root filesystem", rootPath)
	}
	return &WritableRootFs{
		RootFS:   rootFS,
		rootPath: rootPath,
	}, nil
}

type FSStorageConfig struct {
	MediaRoot string `json:"mediaRoot"`
}

type FSStorageBackend struct {
}

func (s *FSStorageBackend) FromJSONConfig(config []byte) (Storage, error) {
	var conf FSStorageConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}
	rootFS, err := newRootFs(conf.MediaRoot)
	if err != nil {
		return nil, err
	}
	return &FSStorage{
		root: rootFS,
	}, nil
}

func (s *FSStorageBackend) ValidateJSONConfig(config []byte) error {
	var conf FSStorageConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return err
	}
	if conf.MediaRoot == "" {
		return fmt.Errorf("mediaRoot is required")
	}
	if err = os.MkdirAll(conf.MediaRoot, 0755); err != nil {
		return err
	}
	return nil
}

type FSStorage struct {
	root *WritableRootFs
}

func (s *FSStorage) GetReader(filename string) io.ReadCloser {
	f, err := s.root.Open(filename)
	if err != nil {
		return nil
	}
	return f
}

func (s *FSStorage) Save(file utils.SeekerReader, filename string, mimeType string) error {
	f, err := s.root.getWriter(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return err
}

func (s *FSStorage) GetMeta(filename string) FileMeta {
	stat, err := s.root.Stat(filename)
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
	return os.Remove(s.root.rootPath + "/" + filename)
}

func (s *FSStorage) CheckConnection() error {
	state, err := os.Stat(s.root.rootPath)
	if err != nil {
		return err
	}
	if !state.IsDir() {
		return fmt.Errorf("mediaRoot %s is not a directory", s.root.rootPath)
	}
	return nil
}
