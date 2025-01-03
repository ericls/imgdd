package image

import (
	dm "imgdd/domainmodels"
	"imgdd/utils"
	"imgdd/utils/pagination"
	"time"
)

type SaveFunc func(file utils.SeekerReader, filename string, mimeType string) error
type PaginationDirection string

var (
	PaginationDirectionAsc  PaginationDirection = "asc"
	PaginationDirectionDesc PaginationDirection = "desc"
)

type ListImagesFilters struct {
	NameContains string
	CreatedAtLt  *time.Time
	CreatedAtGt  *time.Time
	CreatedBy    *string
	Limit        int
}

func FromPaginationFilter(pf *pagination.Filter) *ListImagesFilters {
	if pf == nil {
		return nil
	}
	f := &ListImagesFilters{}
	for _, ff := range pf.Fields {
		if ff.Name == "name" && ff.Operator == pagination.FilterOperatorContains {
			f.NameContains = ff.Value
		} else if ff.Name == "createdAt" && ff.Operator == pagination.FilterOperatorLt {
			t, err := time.Parse(time.RFC3339, ff.Value)
			if err == nil {
				f.CreatedAtLt = &t
			} else {
				panic(err)
			}
		} else if ff.Name == "createdAt" && ff.Operator == pagination.FilterOperatorGt {
			t, err := time.Parse(time.RFC3339, ff.Value)
			if err == nil {
				f.CreatedAtGt = &t
			} else {
				panic(err)
			}
		} else if ff.Name == "createdBy" && ff.Operator == pagination.FilterOperatorEq {
			f.CreatedBy = &ff.Value
		}
	}
	f.Limit = 24
	return f
}

type ListImagesOrdering struct {
	ID        *PaginationDirection `json:"id,omitempty"`
	Name      *PaginationDirection `json:"name,omitempty"`
	CreatedAt *PaginationDirection `json:"createdAt,omitempty"`
	Checksum  string
}

func FromPaginationOrder(po *pagination.Order) *ListImagesOrdering {
	if po == nil {
		return nil
	}
	o := &ListImagesOrdering{}
	for _, of := range po.Fields {
		if of.Name() == "id" {
			if of.Asc() {
				o.ID = &PaginationDirectionAsc
			} else {
				o.ID = &PaginationDirectionDesc
			}
		} else if of.Name() == "name" {
			if of.Asc() {
				o.Name = &PaginationDirectionAsc
			} else {
				o.Name = &PaginationDirectionDesc
			}
		} else if of.Name() == "createdAt" {
			if of.Asc() {
				o.CreatedAt = &PaginationDirectionAsc
			} else {
				o.CreatedAt = &PaginationDirectionDesc
			}
		}
	}
	return o
}

type ImageRepo interface {
	CreateAndSaveUploadedImage(image *dm.Image, fileBytes []byte, storageDefinitionId string, saveFn SaveFunc) (*dm.StoredImage, error)
	ListImages(filters *ListImagesFilters, ordering *ListImagesOrdering) (dm.ListImageResult, error)
	CountImages(filters *ListImagesFilters) (int, error)
}
