package model

import (
	"imgdd/domainmodels"
	"time"
)

type Image struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Identifier      string    `json:"identifier"`
	NominalWidth    int       `json:"nominalWidth"`
	NominalHeight   int       `json:"nominalHeight"`
	NominalByteSize int       `json:"nominalByteSize"`
	CreatedAt       time.Time `json:"createdAt"`
}

func FromImage(i *domainmodels.Image) *Image {
	return &Image{
		ID:              i.Id,
		Name:            i.Name,
		Identifier:      i.Identifier,
		NominalWidth:    int(i.NominalWidth),
		NominalHeight:   int(i.NominalHeight),
		NominalByteSize: int(i.NominalByteSize),
		CreatedAt:       i.CreatedAt,
	}
}

type ImageEdge struct {
	Node   *Image `json:"node"`
	Cursor string `json:"cursor"`
}

type ImageFilterInput struct {
	NameContains *string    `json:"nameContains,omitempty"`
	CreatedAtLte *time.Time `json:"createdAtLte,omitempty"`
	CreatedAtGte *time.Time `json:"createdAtGte,omitempty"`
	CreatedBy    *string    `json:"createdBy,omitempty"`
}

type ImageOrderByInput struct {
	ID        *PaginationDirection `json:"id,omitempty"`
	Name      *PaginationDirection `json:"name,omitempty"`
	CreatedAt *PaginationDirection `json:"createdAt,omitempty"`
}

type ImagePageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	EndCursor       *string `json:"endCursor,omitempty"`
}

type ImagesResult struct {
	Edges    []*ImageEdge   `json:"edges"`
	PageInfo *ImagePageInfo `json:"pageInfo"`
}

type cursorEncoder func(i *domainmodels.Image) string

func FromListImageResult(r *domainmodels.ListImageResult, genCursor cursorEncoder) *ImagesResult {
	edges := make([]*ImageEdge, len(r.Images))
	for i, image := range r.Images {
		edges[i] = &ImageEdge{
			Node:   FromImage(image),
			Cursor: genCursor(image),
		}
	}
	return &ImagesResult{
		Edges: edges,
		PageInfo: &ImagePageInfo{
			HasNextPage:     r.HasNext,
			HasPreviousPage: r.HasPrev,
		},
	}
}
