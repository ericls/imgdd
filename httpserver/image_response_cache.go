package httpserver

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/hashicorp/golang-lru/v2/simplelru"
)

type imageResponseCacheEntry struct {
	key         string
	contentType string
	etag        string
	xImgddSI    string
	body        []byte
}

// imageResponseCache is an opt-in origin-side byte cache for the image HTTP
// paths. It is intended for deployments where a very small set of immutable
// images suddenly accounts for most traffic and imgdd is repeatedly proxying
// those same bytes from a storage backend. It should stay small and bounded by
// bytes; broader caching belongs in a CDN or reverse proxy.
type imageResponseCache struct {
	maxBytes     int64
	maxFileBytes int64

	mu           sync.Mutex
	currentBytes int64
	items        *simplelru.LRU[string, *imageResponseCacheEntry]
}

func newImageResponseCache(maxBytes, maxFileBytes int64) *imageResponseCache {
	if maxBytes <= 0 {
		return nil
	}
	if maxFileBytes <= 0 || maxFileBytes > maxBytes {
		maxFileBytes = maxBytes
	}

	var cache *imageResponseCache
	items, err := simplelru.NewLRU(
		maxEntriesForByteBudget(maxBytes),
		func(_ string, entry *imageResponseCacheEntry) {
			cache.currentBytes -= int64(len(entry.body))
		},
	)
	if err != nil {
		return nil
	}
	cache = &imageResponseCache{
		maxBytes:     maxBytes,
		maxFileBytes: maxFileBytes,
		items:        items,
	}
	return cache
}

func maxEntriesForByteBudget(maxBytes int64) int {
	maxInt := int(^uint(0) >> 1)
	if maxBytes > int64(maxInt) {
		return maxInt
	}
	return int(maxBytes)
}

func (c *imageResponseCache) get(key string) (*imageResponseCacheEntry, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.items.Get(key)
}

func (c *imageResponseCache) put(key string, contentType string, etag string, xImgddSI string, body []byte) {
	if c == nil || int64(len(body)) > c.maxFileBytes {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if existing, ok := c.items.Get(key); ok {
		c.currentBytes -= int64(len(existing.body))
	}
	entry := &imageResponseCacheEntry{
		key:         key,
		contentType: contentType,
		etag:        etag,
		xImgddSI:    xImgddSI,
		body:        append([]byte(nil), body...),
	}
	c.items.Add(key, entry)
	c.currentBytes += int64(len(entry.body))

	for c.currentBytes > c.maxBytes {
		if _, _, ok := c.items.RemoveOldest(); !ok {
			break
		}
	}
}

func writeCachedImageResponse(w http.ResponseWriter, r *http.Request, entry *imageResponseCacheEntry) {
	if r.Header.Get("If-None-Match") == entry.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", entry.contentType)
	w.Header().Set("Content-Length", stringInt(len(entry.body)))
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("ETag", entry.etag)
	w.Header().Set("X-imgdd-si", entry.xImgddSI)
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	w.Write(entry.body)
}

func stringInt(v int) string {
	return strconv.Itoa(v)
}
