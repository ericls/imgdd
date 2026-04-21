package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ericls/imgdd/utils"
	lru "github.com/hashicorp/golang-lru/v2"
)

type IPFSMFSStorageConfig struct {
	ApiUrl     string `json:"apiUrl"`
	PathPrefix string `json:"pathPrefix"`
	Pin        bool   `json:"pin"`
}

func (c *IPFSMFSStorageConfig) Hash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(c.ApiUrl + "|" + c.PathPrefix + "|" + strconv.FormatBool(c.Pin)))
	return h.Sum32()
}

func (c *IPFSMFSStorageConfig) normalizePrefix() {
	prefix := c.PathPrefix
	if prefix == "" {
		prefix = "/"
	} else {
		prefix = "/" + strings.Trim(prefix, "/") + "/"
	}
	c.PathPrefix = prefix
}

type IPFSMFSStorageBackend struct {
	cache *lru.TwoQueueCache[uint32, *IPFSMFSStorage]
}

func (b *IPFSMFSStorageBackend) FromJSONConfig(jsonConfig []byte) (Storage, error) {
	var config IPFSMFSStorageConfig
	if err := json.Unmarshal(jsonConfig, &config); err != nil {
		return nil, err
	}
	config.normalizePrefix()

	hash := config.Hash()
	if storage, ok := b.cache.Get(hash); ok {
		return storage, nil
	}
	storage := &IPFSMFSStorage{
		config: config,
		client: &http.Client{Timeout: 30 * time.Second},
	}
	if err := storage.ensurePrefix(); err != nil {
		return nil, err
	}
	b.cache.Add(hash, storage)
	return storage, nil
}

func (b *IPFSMFSStorageBackend) ValidateJSONConfig(jsonConfig []byte) error {
	var config IPFSMFSStorageConfig
	if err := json.Unmarshal(jsonConfig, &config); err != nil {
		return err
	}
	if config.ApiUrl == "" {
		return errors.New("invalid IPFS MFS storage config. Missing apiUrl")
	}
	return nil
}

type IPFSMFSStorage struct {
	config IPFSMFSStorageConfig
	client *http.Client
}

type ipfsFilesStatResp struct {
	Hash string `json:"Hash"`
	Size uint64 `json:"Size"`
	Type string `json:"Type"`
}

type ipfsIDResp struct {
	ID string `json:"ID"`
}

type ipfsErrorResp struct {
	Message string `json:"Message"`
	Code    int    `json:"Code"`
}

func (s *IPFSMFSStorage) nameWithPrefix(filename string) string {
	return s.config.PathPrefix + filename
}

func (s *IPFSMFSStorage) doRPC(path string, query url.Values, body io.Reader, contentType string) (*http.Response, error) {
	reqURL := strings.TrimRight(s.config.ApiUrl, "/") + "/api/v0/" + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}
	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var kuboErr ipfsErrorResp
		if jsonErr := json.Unmarshal(data, &kuboErr); jsonErr == nil && kuboErr.Message != "" {
			return nil, fmt.Errorf("ipfs rpc %s: %s", path, kuboErr.Message)
		}
		return nil, fmt.Errorf("ipfs rpc %s: status %d: %s", path, resp.StatusCode, string(data))
	}
	return resp, nil
}

func (s *IPFSMFSStorage) ensurePrefix() error {
	if s.config.PathPrefix == "/" {
		return nil
	}
	q := url.Values{}
	q.Set("arg", strings.TrimRight(s.config.PathPrefix, "/"))
	q.Set("parents", "true")
	resp, err := s.doRPC("files/mkdir", q, nil, "")
	if err != nil {
		if strings.Contains(err.Error(), "file already exists") {
			return nil
		}
		return err
	}
	resp.Body.Close()
	return nil
}

func (s *IPFSMFSStorage) statEntry(filename string) (*ipfsFilesStatResp, error) {
	q := url.Values{}
	q.Set("arg", s.nameWithPrefix(filename))
	resp, err := s.doRPC("files/stat", q, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var stat ipfsFilesStatResp
	if err := json.NewDecoder(resp.Body).Decode(&stat); err != nil {
		return nil, err
	}
	return &stat, nil
}

func (s *IPFSMFSStorage) GetReader(filename string) io.ReadCloser {
	q := url.Values{}
	q.Set("arg", s.nameWithPrefix(filename))
	resp, err := s.doRPC("files/read", q, nil, "")
	if err != nil {
		return nil
	}
	return resp.Body
}

func (s *IPFSMFSStorage) Save(file utils.SeekerReader, filename string, mimeType string) error {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	go func() {
		var closeErr error
		defer func() {
			pw.CloseWithError(closeErr)
		}()
		part, err := mw.CreateFormFile("file", filename)
		if err != nil {
			closeErr = err
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			closeErr = err
			return
		}
		closeErr = mw.Close()
	}()

	q := url.Values{}
	q.Set("arg", s.nameWithPrefix(filename))
	q.Set("create", "true")
	q.Set("truncate", "true")
	q.Set("parents", "true")
	resp, err := s.doRPC("files/write", q, pr, mw.FormDataContentType())
	if err != nil {
		return err
	}
	resp.Body.Close()

	if s.config.Pin {
		stat, err := s.statEntry(filename)
		if err != nil {
			return err
		}
		pinQ := url.Values{}
		pinQ.Set("arg", stat.Hash)
		pinResp, err := s.doRPC("pin/add", pinQ, nil, "")
		if err != nil {
			return err
		}
		pinResp.Body.Close()
	}
	return nil
}

func (s *IPFSMFSStorage) GetMeta(filename string) FileMeta {
	stat, err := s.statEntry(filename)
	if err != nil {
		return FileMeta{}
	}
	filenameParts := strings.Split(filename, ".")
	var contentType string
	if len(filenameParts) > 1 {
		contentType = mime.TypeByExtension("." + filenameParts[len(filenameParts)-1])
	}
	return FileMeta{
		ByteSize:    int64(stat.Size),
		ContentType: contentType,
		ETag:        "\"" + stat.Hash + "\"",
	}
}

func (s *IPFSMFSStorage) Delete(filename string) error {
	if s.config.Pin {
		if stat, err := s.statEntry(filename); err == nil {
			pinQ := url.Values{}
			pinQ.Set("arg", stat.Hash)
			if pinResp, err := s.doRPC("pin/rm", pinQ, nil, ""); err == nil {
				pinResp.Body.Close()
			} else if !strings.Contains(err.Error(), "not pinned") {
				return err
			}
		}
	}
	q := url.Values{}
	q.Set("arg", s.nameWithPrefix(filename))
	q.Set("force", "true")
	resp, err := s.doRPC("files/rm", q, nil, "")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (s *IPFSMFSStorage) CheckConnection() error {
	resp, err := s.doRPC("id", nil, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var id ipfsIDResp
	if err := json.NewDecoder(resp.Body).Decode(&id); err != nil {
		return err
	}
	if id.ID == "" {
		return errors.New("ipfs node reported empty ID")
	}
	return nil
}
