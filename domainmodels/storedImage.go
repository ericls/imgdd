package domainmodels

type StoredImage struct {
	Id                  string
	Image               *Image
	StorageDefinitionId string
	FileIdentifier      string
	CopiedFrom          *StoredImage
}
