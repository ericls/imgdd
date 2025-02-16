package storage

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io"
	"strings"

	"github.com/ericls/imgdd/utils"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/studio-b12/gowebdav"
)

type WebDAVStorageConfig struct {
	URL        string `json:"url"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PathPrefix string `json:"pathPrefix"`
}

func (c *WebDAVStorageConfig) Hash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(
		c.URL + "|" + c.Username + "|" + c.Password,
	))
	return h.Sum32()
}

type WebDAVStorage struct {
	config WebDAVStorageConfig
	client *gowebdav.Client
}

type WebDAVBackend struct {
	cache *lru.TwoQueueCache[uint32, *WebDAVStorage]
}

func (b *WebDAVBackend) FromJSONConfig(jsonConfig []byte) (Storage, error) {
	var config WebDAVStorageConfig
	err := json.Unmarshal(jsonConfig, &config)
	if err != nil {
		return nil, err
	}

	hash := config.Hash()
	storage, ok := b.cache.Get(hash)
	if !ok {
		storage = &WebDAVStorage{
			config: config,
			client: gowebdav.NewClient(config.URL, config.Username, config.Password),
		}
		err = storage.client.MkdirAll(storage.config.PathPrefix, 0755)
		if err != nil {
			return nil, err
		}
		b.cache.Add(hash, storage)
	}

	return storage, nil
}

func (b *WebDAVBackend) ValidateJSONConfig(jsonConfig []byte) error {
	var config WebDAVStorageConfig
	err := json.Unmarshal(jsonConfig, &config)
	if err != nil {
		return err
	}
	if config.URL == "" {
		return errors.New("invalid WebDAV storage config. Missing URL")
	}
	return nil
}

func (b *WebDAVStorage) nameWithPrefix(filename string) string {
	prefix := b.config.PathPrefix
	if prefix == "" {
		prefix = "/"
	} else {
		prefix = "/" + strings.Trim(prefix, "/") + "/"
	}
	return prefix + filename
}

func (s *WebDAVStorage) GetReader(filename string) io.ReadCloser {
	reader, err := s.client.ReadStream(s.nameWithPrefix(filename))
	if err != nil {
		return nil
	}
	return reader
}

func (s *WebDAVStorage) Save(file utils.SeekerReader, filename string, mimeType string) error {
	return s.client.WriteStream(s.nameWithPrefix(filename), file, 0644)
}

func (s *WebDAVStorage) GetMeta(filename string) FileMeta {
	i, err := s.client.Stat(s.nameWithPrefix(filename))
	info := i.(*gowebdav.File)
	if err != nil {
		return FileMeta{}
	}
	return FileMeta{
		ByteSize:    info.Size(),
		ETag:        info.ETag(),
		ContentType: info.ContentType(),
	}
}

func (s *WebDAVStorage) Delete(filename string) error {
	return s.client.Remove(s.nameWithPrefix(filename))
}

func (s *WebDAVStorage) CheckConnection() error {
	info, err := s.client.Stat(s.nameWithPrefix(""))
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("not a directory")
	}
	return nil
}
