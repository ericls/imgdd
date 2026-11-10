package image_test

import (
	"strings"
	"testing"
	"time"

	"github.com/ericls/imgdd/db"
	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/test_support"
	"github.com/google/uuid"
)

func TestAlbumRepo(t *testing.T) {
	dbConf := TestServiceMan.GetDBConfig()
	test_support.ResetDatabase(dbConf)
	conn := db.GetConnection(dbConf)
	imageRepo := image.NewDBImageRepo(conn)
	identityRepo := identity.NewDBIdentityRepo(conn)

	owner, err := identityRepo.CreateUserWithOrganization("album-owner@home.arpa", "album-org", "password")
	if err != nil {
		t.Fatal(err)
	}

	createImage := func(albumId string, caption string) *dm.Image {
		t.Helper()
		id := uuid.New().String()
		img, err := imageRepo.CreateImage(&dm.Image{
			Identifier:  id,
			Name:        id + ".png",
			MIMEType:    "image/png",
			CreatedById: owner.Id,
			AlbumId:     albumId,
			Caption:     caption,
		})
		if err != nil {
			t.Fatal(err)
		}
		// created_at has microsecond resolution; keep insertion order unambiguous.
		time.Sleep(2 * time.Millisecond)
		return img
	}

	t.Run("create and get", func(t *testing.T) {
		album, err := imageRepo.CreateAlbum("Fix login layout", owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		if len(album.Identifier) != 22 || strings.ContainsAny(album.Identifier, "+/=") {
			t.Errorf("identifier should be 22 URL-safe chars, got %q", album.Identifier)
		}
		if album.Title != "Fix login layout" || album.CreatedById != owner.Id {
			t.Errorf("unexpected album %+v", album)
		}
		got, err := imageRepo.GetAlbumByIdentifier(album.Identifier)
		if err != nil || got == nil || got.Id != album.Id {
			t.Fatalf("GetAlbumByIdentifier = %+v, %v", got, err)
		}

		other, err := imageRepo.CreateAlbum("", owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		if other.Identifier == album.Identifier {
			t.Error("identifiers should be unique")
		}
	})

	t.Run("unknown or deleted album is nil", func(t *testing.T) {
		for _, identifier := range []string{"", "does-not-exist"} {
			got, err := imageRepo.GetAlbumByIdentifier(identifier)
			if err != nil || got != nil {
				t.Errorf("GetAlbumByIdentifier(%q) = %+v, %v; want nil, nil", identifier, got, err)
			}
		}
		album, err := imageRepo.CreateAlbum("deleted", owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Exec(`UPDATE album_table SET deleted_at = now() WHERE id = $1`, album.Id); err != nil {
			t.Fatal(err)
		}
		got, err := imageRepo.GetAlbumByIdentifier(album.Identifier)
		if err != nil || got != nil {
			t.Errorf("deleted album should not be returned, got %+v, %v", got, err)
		}
	})

	t.Run("validation", func(t *testing.T) {
		if _, err := imageRepo.CreateAlbum("no owner", ""); err != image.ErrAlbumOwnerRequired {
			t.Errorf("expected ErrAlbumOwnerRequired, got %v", err)
		}
		if _, err := imageRepo.CreateAlbum(strings.Repeat("a", dm.AlbumMaxTitleLength+1), owner.Id); err != image.ErrAlbumTitleTooLong {
			t.Errorf("expected ErrAlbumTitleTooLong, got %v", err)
		}
		// The limit is in characters, not bytes.
		if _, err := imageRepo.CreateAlbum(strings.Repeat("界", dm.AlbumMaxTitleLength), owner.Id); err != nil {
			t.Errorf("max-length multi-byte title should be accepted, got %v", err)
		}
	})

	t.Run("images carry album and caption", func(t *testing.T) {
		album, err := imageRepo.CreateAlbum("captions", owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		withAlbum := createImage(album.Id, "step 1: before fix")
		got, err := imageRepo.GetImageById(withAlbum.Id)
		if err != nil {
			t.Fatal(err)
		}
		if got.AlbumId != album.Id || got.Caption != "step 1: before fix" {
			t.Errorf("album/caption not persisted: %+v", got)
		}
		plain := createImage("", "")
		got, err = imageRepo.GetImageById(plain.Id)
		if err != nil {
			t.Fatal(err)
		}
		if got.AlbumId != "" || got.Caption != "" {
			t.Errorf("expected no album or caption, got %+v", got)
		}
	})

	t.Run("list album images", func(t *testing.T) {
		album, err := imageRepo.CreateAlbum("listing", owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		otherAlbum, err := imageRepo.CreateAlbum("other", owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		first := createImage(album.Id, "first")
		deleted := createImage(album.Id, "deleted")
		expired := createImage(album.Id, "expired")
		createImage(otherAlbum.Id, "other album")
		createImage("", "no album")
		expiresLater := createImage(album.Id, "expires later")
		last := createImage(album.Id, "last")

		if err := imageRepo.DeleteImageById(deleted.Id); err != nil {
			t.Fatal(err)
		}
		past := time.Now().Add(-time.Minute)
		if _, err := imageRepo.SetImagesExpiration([]string{expired.Id}, &past); err != nil {
			t.Fatal(err)
		}
		future := time.Now().Add(time.Hour)
		if _, err := imageRepo.SetImagesExpiration([]string{expiresLater.Id}, &future); err != nil {
			t.Fatal(err)
		}

		images, err := imageRepo.ListAlbumImages(album.Id, 0)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{first.Id, expiresLater.Id, last.Id}
		if len(images) != len(want) {
			t.Fatalf("got %d images, want %d", len(images), len(want))
		}
		for i, img := range images {
			if img.Id != want[i] {
				t.Errorf("image %d: got %s (%q), want %s", i, img.Id, img.Caption, want[i])
			}
		}

		limited, err := imageRepo.ListAlbumImages(album.Id, 2)
		if err != nil {
			t.Fatal(err)
		}
		if len(limited) != 2 || limited[0].Id != first.Id {
			t.Errorf("limit should keep the oldest images, got %d", len(limited))
		}
	})
}
