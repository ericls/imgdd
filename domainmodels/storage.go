package domainmodels

type StorageTypeName string

const (
	S3StorageType     StorageTypeName = "s3"
	FSStorageType     StorageTypeName = "fs"
	WebDavStorageType StorageTypeName = "webdav"
)

func (s StorageTypeName) IsValid() bool {
	switch s {
	case S3StorageType, FSStorageType, WebDavStorageType:
		return true
	}
	return false
}

type StorageDefinition struct {
	Id          string
	Identifier  string
	StorageType StorageTypeName
	Config      string
	IsEnabled   bool
	Priority    int32
}
