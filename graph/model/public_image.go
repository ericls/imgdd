package model

import (
	"time"

	"github.com/ericls/imgdd/domainmodels"
)

type PublicImage struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	NominalWidth    int       `json:"nominalWidth"`
	NominalHeight   int       `json:"nominalHeight"`
	NominalByteSize int       `json:"nominalByteSize"`
	CreatedAt       time.Time `json:"createdAt"`
	MIMEType        string    `json:"MIMEType"`
	URL             string    `json:"url"`
	// Not exposed in the schema; used to build the URL.
	Identifier string `json:"-"`
	// Not exposed in the schema; used to resolve viewerIsOwner.
	CreatedById string `json:"-"`
}

func FromPublicImage(i *domainmodels.PublicImage) *PublicImage {
	return &PublicImage{
		ID:              i.Id,
		Name:            i.Name,
		NominalWidth:    int(i.NominalWidth),
		NominalHeight:   int(i.NominalHeight),
		NominalByteSize: int(i.NominalByteSize),
		CreatedAt:       i.CreatedAt,
		MIMEType:        i.MIMEType,
		Identifier:      i.Identifier,
		CreatedById:     i.CreatedById,
	}
}
