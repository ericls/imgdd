package domainmodels

import "time"

// APIKey is a credential that lets an organization user call the HTTP API
// (e.g. /api/v1/createImages) without a browser session. The secret itself is
// never stored; only its hash and a short, displayable prefix are kept.
type APIKey struct {
	Id                 string
	OrganizationUserId string
	Name               string
	KeyPrefix          string
	LastUsedAt         *time.Time
	RevokedAt          *time.Time
	CreatedAt          time.Time
}
