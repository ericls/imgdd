package image

import (
	"time"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/utils"
	"github.com/ericls/imgdd/utils/pagination"
)

type SaveFunc func(file utils.SeekerReader, filename string, mimeType string) error
type PaginationDirection string

func (pd PaginationDirection) Reverse() PaginationDirection {
	if pd == PaginationDirectionAsc {
		return PaginationDirectionDesc
	}
	return PaginationDirectionAsc
}

var (
	PaginationDirectionAsc  PaginationDirection = "asc"
	PaginationDirectionDesc PaginationDirection = "desc"
)

type ListImagesFilters struct {
	NameContains string
	NameLt       string
	NameGt       string
	CreatedAtLt  *time.Time
	CreatedAtGt  *time.Time
	IdGt         string
	IdLt         string
	CreatedBy    *string
	Limit        int
}

func FromPaginationFilter(pf *pagination.Filter) *ListImagesFilters {
	if pf == nil {
		return nil
	}
	f := &ListImagesFilters{}
	for _, ff := range pf.Fields {
		if ff.Name == "name" {
			if ff.Operator == pagination.FilterOperatorContains {
				f.NameContains = ff.Value
			} else if ff.Operator == pagination.FilterOperatorLt {
				f.NameLt = ff.Value
			} else if ff.Operator == pagination.FilterOperatorGt {
				f.NameGt = ff.Value
			} else {
				panic("Invalid operator for name")
			}
		} else if ff.Name == "createdAt" {
			t, err := time.Parse(time.RFC3339Nano, ff.Value)
			if err != nil {
				panic(err)
			}
			if ff.Operator == pagination.FilterOperatorLt {
				f.CreatedAtLt = &t
			} else if ff.Operator == pagination.FilterOperatorGt {
				f.CreatedAtGt = &t
			} else {
				panic("Invalid operator for createdAt")
			}
		} else if ff.Name == "createdBy" && ff.Operator == pagination.FilterOperatorEq {
			f.CreatedBy = &ff.Value
		} else if ff.Name == "id" {
			if ff.Operator == pagination.FilterOperatorLt {
				f.IdLt = ff.Value
			} else if ff.Operator == pagination.FilterOperatorGt {
				f.IdGt = ff.Value
			} else {
				panic("Invalid operator for id")
			}
		}
	}
	f.Limit = dm.ImageResultPerPage
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
	CreateAndSaveUploadedImage(image *dm.Image, mimeType string, fileBytes []byte, storageDefinitionId string, saveFn SaveFunc) (*dm.StoredImage, error)
	ListImages(filtersWithoutCursor *ListImagesFilters, filtersWithCursor *ListImagesFilters, ordering *ListImagesOrdering, reverse bool) (dm.ListImageResult, error)
	CountImages(filters *ListImagesFilters) (int, error)
	GetImageById(id string) (*dm.Image, error)
	DeleteImageById(id string) error
}
