package storage

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io"

	"github.com/ericls/imgdd/utils"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/studio-b12/gowebdav"
)

type WebDAVStorageConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
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

func (s *WebDAVStorage) GetReader(filename string) io.ReadCloser {
	reader, err := s.client.ReadStream("/" + filename)
	if err != nil {
		return nil
	}
	return reader
}

func (s *WebDAVStorage) Save(file utils.SeekerReader, filename string, mimeType string) error {
	return s.client.WriteStream("/"+filename, file, 0644)
}

func (s *WebDAVStorage) GetMeta(filename string) FileMeta {
	i, err := s.client.Stat("/" + filename)
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
	return s.client.Remove("/" + filename)
}

func (s *WebDAVStorage) CheckConnection() error {
	info, err := s.client.Stat("/")
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("not a directory")
	}
	return nil
}
