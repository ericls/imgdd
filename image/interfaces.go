package image

import (
	dm "imgdd/domainmodels"
	"imgdd/utils"
	"time"
)

type SaveFunc func(file utils.SeekerReader, filename string, mimeType string) error
type PaginationDirection string

const (
	PaginationDirectionAsc  PaginationDirection = "asc"
	PaginationDirectionDesc PaginationDirection = "desc"
)

type ListImagesFilters struct {
	NameContains string
	CreatedAtLte time.Time
	CreatedAtGte time.Time
}

type ListImagesOrdering struct {
	ID        *PaginationDirection
	Name      *PaginationDirection
	CreatedAt *PaginationDirection
}

type ImageRepo interface {
	CreateAndSaveUploadedImage(image *dm.Image, fileBytes []byte, storageDefinitionId string, saveFn SaveFunc) (*dm.StoredImage, error)
	ListImages(filters ListImagesFilters, ordering ListImagesOrdering) ([]*dm.Image, error)
}
