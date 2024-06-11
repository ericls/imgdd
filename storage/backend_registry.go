package storage

var BackendRegistry = make(map[string]StorageBackend)

func RegisterBackend(name string, backend StorageBackend) {
	BackendRegistry[name] = backend
}

func GetBackend(name string) StorageBackend {
	return BackendRegistry[name]
}

func init() {
	RegisterBackend("s3", &S3StorageBackend{})
}
