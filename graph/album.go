package graph

// PublicAlbumMaxImages caps how many images the public album query returns.
// The album page polls this query, so the cap also bounds the cost of a poll.
const PublicAlbumMaxImages = 500
