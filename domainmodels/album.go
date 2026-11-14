package domainmodels

import "time"

// AlbumMaxTitleLength is the maximum length of an album title, in characters.
const AlbumMaxTitleLength = 255

// ImageMaxCaptionLength is the maximum length of an image's caption in an
// album, in characters.
const ImageMaxCaptionLength = 2000

// Album groups images in the order they were added, e.g. the screenshots of
// one agent run. Albums and images are many-to-many: an image can be in any
// number of albums, with a different caption in each.
//
// Permission model:
//   - View: anyone holding the Identifier (it is a capability URL, like an
//     image link). An album page never reveals the image's other albums.
//   - Create: authenticated users only, who become the owner.
//   - Add: only the owner. Uploading into the album adds the new image;
//     adding an existing image also requires owning that image, so nobody can
//     put someone else's image on their album page.
//   - Remove: only the owner. Removing an image from an album never deletes
//     the image, and deleting an album never deletes its images.
//
// Id is internal and must never be exposed; Identifier is the only external
// handle.
type Album struct {
	Id          string
	Identifier  string
	Title       string
	CreatedById string
	CreatedAt   time.Time
}

// AlbumImageEntry is an image to add to an album, with its caption there.
type AlbumImageEntry struct {
	ImageId string
	// Caption describes the image in this album, or "" for none.
	Caption string
}

// AlbumImage is an image as listed in an album.
type AlbumImage struct {
	Image *Image
	// Caption describes the image in this album, or "" for none.
	Caption string
	// AddedAt is when the image was added to this album, which sets its
	// place in the album's order.
	AddedAt time.Time
}
