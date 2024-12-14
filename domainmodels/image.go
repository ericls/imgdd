package domainmodels

import "time"

type Image struct {
	Id              string
	CreatedById     string
	CreatedAt       time.Time
	Name            string
	Identifier      string
	RootId          string
	ParentId        string
	UploaderIP      string
	MIMEType        string
	NominalWidth    int32
	NominalHeight   int32
	NominalByteSize int32
}

type StoredImage struct {
	Id                string
	Image             *Image
	StorageDefinition *StorageDefinition
	FileIdentifier    string
	CopiedFrom        *StoredImage
}
