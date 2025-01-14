package domainmodels

import (
	"imgdd/utils/pagination"
	"mime"
	"strings"
	"time"
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

func getFileExtFromMIMEType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		exts, err := mime.ExtensionsByType(mimeType)
		if err == nil && len(exts) > 0 {
			return exts[0]
		}
	}
	return ""
}

func getImageURL(maybeImageDomain string, identifier string, mimeType string, isSecure bool) string {
	suffix := getFileExtFromMIMEType(mimeType)
	if suffix == "" {
		return ""
	}
	filename := identifier + suffix
	if maybeImageDomain != "" {
		if isSecure {
			return "https://" + maybeImageDomain + "/image/" + filename
		}
		return "http://" + maybeImageDomain + "/image/" + filename
	}
	return "/image/" + filename
}

func (i *Image) GetURL(imageDomain string, isSecure bool) string {
	return getImageURL(imageDomain, i.Identifier, i.MIMEType, isSecure)
}
