package graph

// DefaultPublicAlbumMaxImages caps how many images the public album query
// returns, unless Resolver.PublicAlbumMaxImages says otherwise. The album page
// polls this query, so the cap also bounds the cost of a poll. When an album
// has more, the newest images are kept, since the page is a live view of what
// the agent is doing now.
const DefaultPublicAlbumMaxImages = 500

func (r *Resolver) publicAlbumMaxImages() int {
	if r.PublicAlbumMaxImages > 0 {
		return r.PublicAlbumMaxImages
	}
	return DefaultPublicAlbumMaxImages
}
