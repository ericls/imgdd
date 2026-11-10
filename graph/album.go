package graph

// PublicAlbumMaxImages caps how many images the public album query returns.
// The album page polls this query, so the cap also bounds the cost of a poll.
// When an album has more, the newest images are kept, since the page is a
// live view of what the agent is doing now. A variable so tests can lower it.
var PublicAlbumMaxImages = 500
