// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type CreateUserWithOrganizationInput struct {
	UserEmail        string `json:"userEmail"`
	UserPassword     string `json:"userPassword"`
	OrganizationName string `json:"organizationName"`
}

type Mutation struct {
}

type Query struct {
}

type ViewerResult struct {
	Viewer *Viewer `json:"viewer"`
}

type PaginationDirection string

const (
	PaginationDirectionAsc  PaginationDirection = "asc"
	PaginationDirectionDesc PaginationDirection = "desc"
)

var AllPaginationDirection = []PaginationDirection{
	PaginationDirectionAsc,
	PaginationDirectionDesc,
}

func (e PaginationDirection) IsValid() bool {
	switch e {
	case PaginationDirectionAsc, PaginationDirectionDesc:
		return true
	}
	return false
}

func (e PaginationDirection) String() string {
	return string(e)
}

func (e *PaginationDirection) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PaginationDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PaginationDirection", str)
	}
	return nil
}

func (e PaginationDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
