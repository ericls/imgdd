package domainmodels

import (
	"testing"

	"github.com/google/uuid"
)

func TestImageToPublic(t *testing.T) {
	ownerId := "00000000-0000-0000-0000-000000000001"
	cases := []struct {
		name        string
		createdById string
	}{
		{"no owner", ""},
		{"zero uuid owner", uuid.Nil.String()},
		{"owned", ownerId},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			img := &Image{Id: "id", Name: "a.png", Identifier: "abc", MIMEType: "image/png", CreatedById: tc.createdById}
			publicImage := img.ToPublic()
			if publicImage == nil || publicImage.Id != img.Id || publicImage.Name != img.Name {
				t.Fatalf("unexpected public image %+v", publicImage)
			}
			if url := publicImage.GetURL("", false, nil, ImageURLFormat_CANONICAL); url != img.GetURL("", false, nil, ImageURLFormat_CANONICAL) {
				t.Fatalf("expected public URL to match image URL, got %q", url)
			}
			if publicImage.IsOwnedBy(ownerId) != (tc.createdById == ownerId) {
				t.Fatalf("unexpected IsOwnedBy(%q) for owner %q", ownerId, tc.createdById)
			}
			if publicImage.IsOwnedBy("") {
				t.Fatal("expected no image to be owned by an empty user id")
			}
		})
	}
}
