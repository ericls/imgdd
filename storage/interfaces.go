package storage

import (
	dm "imgdd/domainmodels"
	"io"
)

type FileMeta struct {
	ByteSize    int64
	ContentType string
	ETag        string
}

type SeekerReader interface {
	io.Seeker
	io.Reader
}

type Storage interface {
	GetReader(filename string) io.Reader
	Save(file SeekerReader, filename string, mimeType string) error
	GetMeta(filename string) FileMeta
	Delete(filename string) error
	CheckConnection() error
}

type StorageBackend interface {
	FromJSON(jsonConfig []byte) (Storage, error)
}

type StorageRepo interface {
	GetStorageDefinitionByID(id string) (*dm.StorageDefinition, error)
	ListStorageDefinitions(organizationUserId string) ([]*dm.StorageDefinition, error)
	CreateStorageDefinition(organizationUserId string, storage_type string, config string, identifier string, isEnabled bool, priority int32) (*dm.StorageDefinition, error)
}
