package storage

import (
	"io"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/utils"
)

type FileMeta struct {
	ByteSize    int64
	ContentType string
	ETag        string // This should be quoted
}

type Storage interface {
	GetReader(filename string) io.ReadCloser
	Save(file utils.SeekerReader, filename string, mimeType string) error
	GetMeta(filename string) FileMeta
	Delete(filename string) error
	CheckConnection() error
}

type StorageBackend interface {
	FromJSONConfig(jsonConfig []byte) (Storage, error)
	ValidateJSONConfig(jsonConfig []byte) error
}

type StorageRepo interface {
	GetStorageDefinitionById(id string) (*dm.StorageDefinition, error)
	GetStorageDefinitionByIdentifier(id string) (*dm.StorageDefinition, error)
	// order by priority
	ListStorageDefinitions() ([]*dm.StorageDefinition, error)
	CreateStorageDefinition(storage_type string, config string, identifier string, isEnabled bool, priority int64) (*dm.StorageDefinition, error)
	UpdateStorageDefinition(identifier string, storage_type *string, config *string, isEnabled *bool, priority *int64) (*dm.StorageDefinition, error)

	GetStoredImageByIdentifierAndMimeType(identifier, mime string) (*dm.StoredImage, error)
	GetStoredImagesByIds(ids []string) ([]*dm.StoredImage, error)
	GetStoredImageIdsByImageIds(imageIds []string) (map[string][]string, error)
}
