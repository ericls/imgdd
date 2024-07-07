package image

import (
	dm "imgdd/domainmodels"
	"imgdd/utils"
)

type SaveFunc func(file utils.SeekerReader, filename string, mimeType string) error

type ImageRepo interface {
	CreateAndSaveUploadedImage(image *dm.Image, fileBytes []byte, storageDefinitionId string, saveFn SaveFunc) (*dm.StoredImage, error)
}
