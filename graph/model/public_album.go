package model

import "time"

type PublicAlbum struct {
	// ID is the album's public identifier, never its internal id.
	ID        string              `json:"id"`
	Title     string              `json:"title"`
	CreatedAt time.Time           `json:"createdAt"`
	Images    []*PublicAlbumImage `json:"images"`
	Truncated bool                `json:"truncated"`
}

// PublicAlbumImage is an image's place in an album.
type PublicAlbumImage struct {
	Image   *PublicImage `json:"image"`
	Caption *string      `json:"caption"`
	AddedAt time.Time    `json:"addedAt"`
}
