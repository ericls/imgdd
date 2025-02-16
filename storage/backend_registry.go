package storage

import (
	"github.com/ericls/imgdd/domainmodels"
	lru "github.com/hashicorp/golang-lru/v2"
)

var BackendRegistry = make(map[domainmodels.StorageTypeName]StorageBackend)

func RegisterBackend(name domainmodels.StorageTypeName, backend StorageBackend) {
	BackendRegistry[name] = backend
}

func GetBackend(name domainmodels.StorageTypeName) StorageBackend {
	return BackendRegistry[name]
}

func init() {
	s3StorageInstanceCache, err := lru.New2Q[uint32, *S3Storage](0x20)
	if err != nil {
		panic(err)
	}
	webDavStorageInstanceCache, err := lru.New2Q[uint32, *WebDAVStorage](0x20)
	if err != nil {
		panic(err)
	}
	RegisterBackend(domainmodels.S3StorageType, &S3StorageBackend{
		cache: s3StorageInstanceCache,
	})
	RegisterBackend(domainmodels.FSStorageType, &FSStorageBackend{})
	RegisterBackend(domainmodels.WebDavStorageType, &WebDAVBackend{
		cache: webDavStorageInstanceCache,
	})
}
