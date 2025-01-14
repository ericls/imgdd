package domainmodels

type StoredImage struct {
	Id                string
	Image             *Image
	StorageDefinition *StorageDefinition
	FileIdentifier    string
	CopiedFrom        *StoredImage
}
