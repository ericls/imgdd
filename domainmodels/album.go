package domainmodels

import "time"

// AlbumMaxTitleLength is the maximum length of an album title, in characters.
const AlbumMaxTitleLength = 255

// ImageMaxCaptionLength is the maximum length of an image caption, in characters.
const ImageMaxCaptionLength = 2000

// Album groups images in upload order, e.g. the screenshots of one agent run.
//
// Permission model: anyone holding the Identifier can view the album (it is a
// capability URL, like an image link), but only the creator can add images to
// it, and images join an album only when they are uploaded. Id is internal and
// must never be exposed; Identifier is the only external handle.
type Album struct {
	Id          string
	Identifier  string
	Title       string
	CreatedById string
	CreatedAt   time.Time
}
