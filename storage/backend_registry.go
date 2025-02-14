package storage

import lru "github.com/hashicorp/golang-lru/v2"

var BackendRegistry = make(map[string]StorageBackend)

func RegisterBackend(name string, backend StorageBackend) {
	BackendRegistry[name] = backend
}

func GetBackend(name string) StorageBackend {
	return BackendRegistry[name]
}

func init() {
	s3StorageInstanceCache, err := lru.New2Q[uint32, *S3Storage](0x20)
	if err != nil {
		panic(err)
	}
	RegisterBackend("s3", &S3StorageBackend{
		cache: s3StorageInstanceCache,
	})
	RegisterBackend("fs", &FSStorageBackend{})
}
