package domainmodels

import (
	"time"

	"github.com/ericls/imgdd/utils/pagination"
)

type ListImageResult struct {
	Images  []*Image
	HasNext bool
	HasPrev bool
}

const ImageResultPerPage = 24

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
		return image.CreatedAt.Format(time.RFC3339Nano)
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
