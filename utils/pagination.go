package utils

import (
	"encoding/base64"
	"encoding/json"
)

type PaginationDirection string

const (
	PaginationAsc  PaginationDirection = "asc"
	PaginationDesc PaginationDirection = "desc"
)

type PaginationOrderBy struct {
	Field     string              `json:"field"`
	Direction PaginationDirection `json:"direction"`
	Value     string              `json:"value"`
}

type PaginationCursor struct {
	OrderByItems []PaginationOrderBy `json:"orderByItems"`
}

func NewPaginationCursor() *PaginationCursor {
	return &PaginationCursor{
		OrderByItems: make([]PaginationOrderBy, 0),
	}
}

func (c *PaginationCursor) AddOrderByItem(field string, direction PaginationDirection, value string) {
	c.OrderByItems = append(c.OrderByItems, PaginationOrderBy{
		Field:     field,
		Direction: direction,
		Value:     value,
	})
}

func (c *PaginationCursor) Stringify() string {
	// Base64 encoded JSON string
	jsonBytes, _ := json.Marshal(c)
	base64String := base64.URLEncoding.EncodeToString(jsonBytes)
	return base64String
}

func ParsePaginationCursor(cursor string) (*PaginationCursor, error) {
	// Decode base64 string
	jsonBytes, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	paginationCursor := &PaginationCursor{}
	err = json.Unmarshal(jsonBytes, paginationCursor)
	if err != nil {
		return nil, err
	}

	return paginationCursor, nil
}
