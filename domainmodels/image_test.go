package domainmodels

import (
	"testing"
	"time"
)

func TestImageIsExpiredAt(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Second)
	after := now.Add(time.Second)

	cases := []struct {
		name      string
		expiresAt *time.Time
		want      bool
	}{
		{"never expires", nil, false},
		{"expiry in the past", &before, true},
		{"expiry exactly now", &now, true},
		{"expiry in the future", &after, false},
	}
	for _, c := range cases {
		img := &Image{ExpiresAt: c.expiresAt}
		if got := img.IsExpiredAt(now); got != c.want {
			t.Errorf("%s: IsExpiredAt = %v, want %v", c.name, got, c.want)
		}
	}
}
