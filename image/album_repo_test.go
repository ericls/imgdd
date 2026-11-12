package image_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ericls/imgdd/db"
	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/test_support"
	"github.com/ericls/imgdd/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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

	createImage := func() *dm.Image {
		t.Helper()
		id := uuid.New().String()
		img, err := imageRepo.CreateImage(&dm.Image{
			Identifier:  id,
			Name:        id + ".png",
			MIMEType:    "image/png",
			CreatedById: owner.Id,
		})
		if err != nil {
			t.Fatal(err)
		}
		return img
	}
	addToAlbum := func(album *dm.Album, img *dm.Image, caption string) {
		t.Helper()
		added, err := imageRepo.AddImagesToAlbum(album.Id, []dm.AlbumImageEntry{{ImageId: img.Id, Caption: caption}}, owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		if added != 1 {
			t.Fatalf("expected 1 image added, got %d", added)
		}
		// added_at has microsecond resolution; keep the order unambiguous.
		time.Sleep(2 * time.Millisecond)
	}
	listIds := func(album *dm.Album, limit int) ([]string, []string) {
		t.Helper()
		images, err := imageRepo.ListAlbumImages(album.Id, limit)
		if err != nil {
			t.Fatal(err)
		}
		ids := make([]string, len(images))
		captions := make([]string, len(images))
		for i, img := range images {
			ids[i] = img.Image.Id
			captions[i] = img.Caption
		}
		return ids, captions
	}
	newAlbum := func(title string) *dm.Album {
		t.Helper()
		album, err := imageRepo.CreateAlbum(title, owner.Id)
		if err != nil {
			t.Fatal(err)
		}
		return album
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

	t.Run("an image can be in several albums with a caption in each", func(t *testing.T) {
		run := newAlbum("run")
		showcase := newAlbum("showcase")
		shot := createImage()
		addToAlbum(run, shot, "step 3: button still wraps")
		addToAlbum(showcase, shot, "")

		ids, captions := listIds(run, 0)
		require.Equal(t, []string{shot.Id}, ids)
		require.Equal(t, []string{"step 3: button still wraps"}, captions)
		ids, captions = listIds(showcase, 0)
		require.Equal(t, []string{shot.Id}, ids)
		require.Equal(t, []string{""}, captions)
	})

	t.Run("adding is idempotent and keeps the first caption", func(t *testing.T) {
		album := newAlbum("idempotent")
		shot := createImage()
		addToAlbum(album, shot, "first")
		added, err := imageRepo.AddImagesToAlbum(album.Id, []dm.AlbumImageEntry{{ImageId: shot.Id, Caption: "second"}}, owner.Id)
		require.NoError(t, err)
		require.Equal(t, 0, added)
		_, captions := listIds(album, 0)
		require.Equal(t, []string{"first"}, captions)

		_, err = imageRepo.AddImagesToAlbum(album.Id, []dm.AlbumImageEntry{{ImageId: shot.Id, Caption: strings.Repeat("a", dm.ImageMaxCaptionLength+1)}}, owner.Id)
		require.ErrorIs(t, err, image.ErrCaptionTooLong)
	})

	t.Run("albums list images in the order they were added", func(t *testing.T) {
		album := newAlbum("order")
		older := createImage()
		newer := createImage()
		addToAlbum(album, newer, "")
		addToAlbum(album, older, "")
		ids, _ := listIds(album, 0)
		require.Equal(t, []string{newer.Id, older.Id}, ids)
	})

	t.Run("removing leaves the image and its other albums alone", func(t *testing.T) {
		album := newAlbum("remove")
		other := newAlbum("keep")
		shot := createImage()
		kept := createImage()
		addToAlbum(album, shot, "")
		addToAlbum(album, kept, "")
		addToAlbum(other, shot, "")

		removed, err := imageRepo.RemoveImagesFromAlbum(album.Id, []string{shot.Id, uuid.New().String()})
		require.NoError(t, err)
		require.Equal(t, 1, removed)
		ids, _ := listIds(album, 0)
		require.Equal(t, []string{kept.Id}, ids)
		ids, _ = listIds(other, 0)
		require.Equal(t, []string{shot.Id}, ids)
		stillThere, err := imageRepo.GetImageById(shot.Id)
		require.NoError(t, err)
		require.NotNil(t, stillThere)
	})

	t.Run("uploading into an album is atomic", func(t *testing.T) {
		album := newAlbum("upload")
		storageDef, err := storage.NewDBStorageDefRepo(conn).CreateStorageDefinition("fs", `{"mediaRoot":"/nonexistent"}`, "album-upload", true, 1)
		require.NoError(t, err)
		storageDefId := storageDef.Id
		newUpload := func() *dm.Image {
			id := uuid.New().String()
			return &dm.Image{Identifier: id, Name: id + ".png", MIMEType: "image/png", CreatedById: owner.Id}
		}
		ok := func(utils.SeekerReader, string, string) error { return nil }
		stored, err := imageRepo.CreateAndSaveUploadedImageInAlbum(newUpload(), "image/png", []byte("png"), storageDefId, ok, album.Id, "uploaded")
		require.NoError(t, err)
		ids, captions := listIds(album, 0)
		require.Equal(t, []string{stored.Image.Id}, ids)
		require.Equal(t, []string{"uploaded"}, captions)

		failing := func(utils.SeekerReader, string, string) error { return errors.New("storage is down") }
		failed := newUpload()
		_, err = imageRepo.CreateAndSaveUploadedImageInAlbum(failed, "image/png", []byte("png"), storageDefId, failing, album.Id, "")
		require.Error(t, err)
		ids, _ = listIds(album, 0)
		require.Equal(t, []string{stored.Image.Id}, ids)

		// A failure adding to the album rolls back the image too.
		before, err := imageRepo.CountImages(&image.ListImagesFilters{})
		require.NoError(t, err)
		_, err = imageRepo.CreateAndSaveUploadedImageInAlbum(newUpload(), "image/png", []byte("png"), storageDefId, ok, uuid.New().String(), "")
		require.Error(t, err)
		after, err := imageRepo.CountImages(&image.ListImagesFilters{})
		require.NoError(t, err)
		require.Equal(t, before, after)
	})

	t.Run("list album images", func(t *testing.T) {
		album := newAlbum("listing")
		otherAlbum := newAlbum("other")
		first := createImage()
		deleted := createImage()
		expired := createImage()
		expiresLater := createImage()
		last := createImage()
		for _, img := range []*dm.Image{first, deleted, expired, expiresLater, last} {
			addToAlbum(album, img, "")
		}
		addToAlbum(otherAlbum, createImage(), "")
		createImage()

		require.NoError(t, imageRepo.DeleteImageById(deleted.Id))
		past := time.Now().Add(-time.Minute)
		_, err := imageRepo.SetImagesExpiration([]string{expired.Id}, &past)
		require.NoError(t, err)
		future := time.Now().Add(time.Hour)
		_, err = imageRepo.SetImagesExpiration([]string{expiresLater.Id}, &future)
		require.NoError(t, err)

		ids, _ := listIds(album, 0)
		require.Equal(t, []string{first.Id, expiresLater.Id, last.Id}, ids)
		// A limit keeps the most recently added images, oldest first.
		ids, _ = listIds(album, 2)
		require.Equal(t, []string{expiresLater.Id, last.Id}, ids)
	})
}
