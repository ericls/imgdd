package model

import (
	"imgdd/domainmodels"
	"imgdd/utils/pagination"
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
	CreatedAtLt  *time.Time `json:"createdAtLte,omitempty"`
	CreatedAtGt  *time.Time `json:"createdAtGte,omitempty"`
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
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
	TotalCount      *int    `json:"count,omitempty"`
	CurrentCount    *int    `json:"currentCount,omitempty"`
}

type ImagesResult struct {
	Edges    []*ImageEdge   `json:"edges"`
	PageInfo *ImagePageInfo `json:"pageInfo"`
}

type CursorEncoder func(i *domainmodels.Image) string

func FromListImageResult(r *domainmodels.ListImageResult, totalCount int, genCursor CursorEncoder) *ImagesResult {
	currentCount := len(r.Images)
	edges := make([]*ImageEdge, currentCount)
	var lastCursor string
	var firstCursor string
	for i, image := range r.Images {
		cursor := genCursor(image)
		edges[i] = &ImageEdge{
			Node:   FromImage(image),
			Cursor: cursor,
		}
		lastCursor = cursor
		if i == 0 {
			firstCursor = cursor
		}
	}
	return &ImagesResult{
		Edges: edges,
		PageInfo: &ImagePageInfo{
			HasNextPage:     r.HasNext,
			HasPreviousPage: r.HasPrev,
			TotalCount:      &totalCount,
			CurrentCount:    &currentCount,
			StartCursor:     &firstCursor,
			EndCursor:       &lastCursor,
		},
	}
}

func MakeImagePaginator(orderInput *ImageOrderByInput, filterInput *ImageFilterInput) *pagination.Paginator {
	order := pagination.Order{}
	if orderInput != nil {
		if orderInput.CreatedAt != nil {
			order.AddField(domainmodels.NewImageOrderField("createdAt", *orderInput.CreatedAt == PaginationDirectionAsc))
		}
		if orderInput.ID != nil {
			order.AddField(domainmodels.NewImageOrderField("id", *orderInput.ID == PaginationDirectionAsc))
		}
		if orderInput.Name != nil {
			order.AddField(domainmodels.NewImageOrderField("name", *orderInput.Name == PaginationDirectionAsc))
		}
	}
	if len(order.Fields) == 0 {
		order.AddField(domainmodels.NewImageOrderField("createdAt", false))
	}
	filter := pagination.Filter{}
	if filterInput != nil {
		if filterInput.CreatedAtGt != nil {
			filter.AddFilterField("createdAt", pagination.FilterOperatorGt, filterInput.CreatedAtGt.Format(time.RFC3339))
		}
		if filterInput.CreatedAtLt != nil {
			filter.AddFilterField("createdAt", pagination.FilterOperatorLt, filterInput.CreatedAtLt.Format(time.RFC3339))
		}
		if filterInput.NameContains != nil {
			filter.AddFilterField("name", pagination.FilterOperatorContains, *filterInput.NameContains)
		}
		if filterInput.CreatedBy != nil {
			filter.AddFilterField("createdBy", pagination.FilterOperatorEq, *filterInput.CreatedBy)
		}
	}
	return &pagination.Paginator{
		Order:  &order,
		Filter: &filter,
	}
}
