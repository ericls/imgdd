package storage

import lru "github.com/hashicorp/golang-lru/v2"

var metaCache, _ = lru.New2Q[string, FileMeta](2048)

func GetMetaCached(storage Storage, filename string) FileMeta {
	if meta, ok := metaCache.Get(filename); ok {
		return meta
	}

	meta := storage.GetMeta(filename)
	metaCache.Add(filename, meta)
	return meta
}
