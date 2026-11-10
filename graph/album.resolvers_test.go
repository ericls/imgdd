package graph_test

import (
	"testing"
	"time"

	"github.com/ericls/imgdd/domainmodels"

	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type publicAlbumResp struct {
	PublicAlbum *struct {
		ID        string
		Title     string
		Truncated bool
		Images    []struct {
			Caption *string
			Image   struct {
				ID            string
				URL           string
				ViewerIsOwner bool
			}
		}
	}
}

func queryPublicAlbum(tc *TestContext, id string) (*publicAlbumResp, error) {
	var resp publicAlbumResp
	err := tc.client.Post(`
	query publicAlbum($id: ID!) {
		publicAlbum(id: $id) {
			id
			title
			truncated
			images {
				caption
				image { id url viewerIsOwner }
			}
		}
	}`, &resp, client.Var("id", id))
	return &resp, err
}

func createAlbumImage(t *testing.T, tc *TestContext, ownerId string, storageDefId string, albumId string, caption string) *domainmodels.Image {
	t.Helper()
	identifier := uuid.New().String()
	img := domainmodels.Image{
		CreatedById:     ownerId,
		MIMEType:        "image/png",
		Name:            identifier + ".png",
		Identifier:      identifier,
		NominalByteSize: 100,
		NominalWidth:    100,
		NominalHeight:   100,
		AlbumId:         albumId,
		Caption:         caption,
	}
	si, err := tc.imageRepo.CreateAndSaveUploadedImage(&img, "image/png", []byte(""), storageDefId, mockSaveFunc)
	require.NoError(t, err)
	// created_at has microsecond resolution; keep insertion order unambiguous.
	time.Sleep(2 * time.Millisecond)
	return si.Image
}

func tPublicAlbumIsReadableByLinkHolders(t *testing.T, tc *TestContext) {
	sd := createStorageDefinition(t, tc)
	owner := tc.forceAuthenticate()
	album, err := tc.imageRepo.CreateAlbum("Fix login layout", owner.Id)
	require.NoError(t, err)
	first := createAlbumImage(t, tc, owner.Id, sd.Id, album.Id, "before fix")
	expired := createAlbumImage(t, tc, owner.Id, sd.Id, album.Id, "expired")
	deleted := createAlbumImage(t, tc, owner.Id, sd.Id, album.Id, "deleted")
	second := createAlbumImage(t, tc, owner.Id, sd.Id, album.Id, "")
	createAlbumImage(t, tc, owner.Id, sd.Id, "", "not in album")

	past := time.Now().Add(-time.Minute)
	_, err = tc.imageRepo.SetImagesExpiration([]string{expired.Id}, &past)
	require.NoError(t, err)
	require.NoError(t, tc.imageRepo.DeleteImageById(deleted.Id))

	// Anonymous visitors can read the album.
	tc.clearAuthenticationInfo()
	resp, err := queryPublicAlbum(tc, album.Identifier)
	require.NoError(t, err)
	got := resp.PublicAlbum
	require.NotNil(t, got)
	require.Equal(t, album.Identifier, got.ID)
	require.Equal(t, "Fix login layout", got.Title)
	require.False(t, got.Truncated)
	require.Len(t, got.Images, 2)
	require.Equal(t, first.Id, got.Images[0].Image.ID)
	require.Equal(t, "before fix", *got.Images[0].Caption)
	require.NotEmpty(t, got.Images[0].Image.URL)
	require.False(t, got.Images[0].Image.ViewerIsOwner)
	require.Equal(t, second.Id, got.Images[1].Image.ID)
	require.Nil(t, got.Images[1].Caption)

	// The owner sees the same album, marked as theirs.
	tc.setAuthenticatedUser(owner)
	resp, err = queryPublicAlbum(tc, album.Identifier)
	require.NoError(t, err)
	require.True(t, resp.PublicAlbum.Images[0].Image.ViewerIsOwner)
}

func tPublicAlbumDoesNotExposeOwner(t *testing.T, tc *TestContext) {
	owner := tc.forceAuthenticate()
	album, err := tc.imageRepo.CreateAlbum("private owner", owner.Id)
	require.NoError(t, err)
	for _, field := range []string{"createdBy { id }", "createdById", "owner { id }"} {
		var resp publicAlbumResp
		err := tc.client.Post(`
		query publicAlbum($id: ID!) {
			publicAlbum(id: $id) {
				id
				`+field+`
			}
		}`, &resp, client.Var("id", album.Identifier))
		require.Error(t, err, field)
	}
}

func tPublicAlbumNotFound(t *testing.T, tc *TestContext) {
	owner := tc.forceAuthenticate()
	album, err := tc.imageRepo.CreateAlbum("album", owner.Id)
	require.NoError(t, err)

	// The internal id is not a valid handle.
	for _, id := range []string{"", "does-not-exist", album.Id} {
		resp, err := queryPublicAlbum(tc, id)
		require.NoError(t, err, id)
		require.Nil(t, resp.PublicAlbum, id)
	}
}

func TestAlbumResolvers(t *testing.T) {
	tc := newTestContext(t)
	tc.runTestCases(
		tPublicAlbumIsReadableByLinkHolders,
		tPublicAlbumDoesNotExposeOwner,
		tPublicAlbumNotFound,
	)
}
