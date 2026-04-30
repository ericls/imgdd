package storage

import (
	"io"
	"iter"

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

type StorageDefRepo interface {
	GetStorageDefinitionById(id string) (*dm.StorageDefinition, error)
	GetStorageDefinitionsByIds(ids []string) ([]*dm.StorageDefinition, error)
	GetStorageDefinitionByIdentifier(id string) (*dm.StorageDefinition, error)
	// order by priority
	ListStorageDefinitions() ([]*dm.StorageDefinition, error)
	CreateStorageDefinition(storage_type string, config string, identifier string, isEnabled bool, priority int64) (*dm.StorageDefinition, error)
	UpdateStorageDefinition(identifier string, storage_type *string, config *string, isEnabled *bool, priority *int64) (*dm.StorageDefinition, error)
}

type StoredImageRepo interface {
	GetStoredImageByIdentifierAndMimeType(identifier, mime string) ([]*dm.StoredImage, error)
	GetStoredImagesByIds(ids []string) ([]*dm.StoredImage, error)
	GetStoredImageIdsByImageIds(imageIds []string) (map[string][]string, error)
	GetStoredImagesToDelete() ([]*dm.StoredImage, error)
	MarkStoredImagesAsDeleted(ids []string) error
	// Returns all non-deleted stored images for a given image ID (with Image populated for MIMEType).
	GetStoredImagesByImageId(imageId string) ([]*dm.StoredImage, error)
	// Returns the total count of images to replicate and a lazy iterator that streams them
	// one row at a time. The count is fetched eagerly so callers can report progress.
	GetStoredImagesForReplication(sourceStorageDefId string, targetStorageDefId string) (int, iter.Seq2[*dm.StoredImage, error], error)
}
