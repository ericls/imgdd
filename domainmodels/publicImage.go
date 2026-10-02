package domainmodels

import (
	"time"
)

// PublicImage is the part of an Image that can be shown to anyone, including anonymous visitors.
//
// Any live image can be fetched by its link without authentication, so this view only exposes
// what a link holder can already see (content, size, type), plus its name and upload time.
// It must not expose the owner or related images (lineage, edit changes, stored copies):
// e.g. a blurred image's base parent is the unblurred original.
type PublicImage struct {
	Id              string
	CreatedAt       time.Time
	Name            string
	Identifier      string
	MIMEType        string
	NominalWidth    int32
	NominalHeight   int32
	NominalByteSize int32
	// Not for display; lets callers tell whether the viewer owns the image.
	CreatedById string
}

// ToPublic returns the public view of the image.
func (i *Image) ToPublic() *PublicImage {
	return &PublicImage{
		Id:              i.Id,
		CreatedAt:       i.CreatedAt,
		Name:            i.Name,
		Identifier:      i.Identifier,
		MIMEType:        i.MIMEType,
		NominalWidth:    i.NominalWidth,
		NominalHeight:   i.NominalHeight,
		NominalByteSize: i.NominalByteSize,
		CreatedById:     i.CreatedById,
	}
}

// IsOwnedBy reports whether the image was uploaded by the given organization user.
func (i *PublicImage) IsOwnedBy(organizationUserId string) bool {
	return organizationUserId != "" && i.CreatedById == organizationUserId
}

func (i *PublicImage) GetURL(imageDomain string, isSecure bool, storedImages []*ExternalImageIdentifier, format ImageURLFormat) string {
	image := Image{Id: i.Id, Identifier: i.Identifier, MIMEType: i.MIMEType}
	return image.GetURL(imageDomain, isSecure, storedImages, format)
}
