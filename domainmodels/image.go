package domainmodels

import (
	"imgdd/utils/pagination"
	"time"
)

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

type ListImageResult struct {
	Images  []*Image
	HasNext bool
	HasPrev bool
}

type ImageOrderField struct {
	pagination.BaseOrderField
}

func (f *ImageOrderField) GetValue(object interface{}) string {
	image := object.(*Image)
	switch f.Name() {
	case "id":
		return image.Id
	case "name":
		return image.Name
	case "createdAt":
		return image.CreatedAt.Format(time.RFC3339)
	}
	return ""
}

func NewImageOrderField(name string, asc bool) *ImageOrderField {
	return &ImageOrderField{
		pagination.BaseOrderField{
			FieldName: name,
			FieldAsc:  asc,
		},
	}
}
