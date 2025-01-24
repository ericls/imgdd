package storage

import lru "github.com/hashicorp/golang-lru/v2"

var metaCache, _ = lru.New2Q[cacheKey, FileMeta](2048)

type cacheKey struct {
	storageDefId string
	filename     string
}

func GetMetaCached(storage Storage, storageDefId string, filename string) FileMeta {
	key := cacheKey{storageDefId, filename}
	if meta, ok := metaCache.Get(key); ok {
		return meta
	}

	meta := storage.GetMeta(filename)
	metaCache.Add(key, meta)
	return meta
}
