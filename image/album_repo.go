//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package image

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"slices"
	"unicode/utf8"

	"github.com/ericls/imgdd/db/.gen/imgdd/public/model"
	. "github.com/ericls/imgdd/db/.gen/imgdd/public/table"
	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

var ErrAlbumTitleTooLong = errors.New("album title is too long")
var ErrAlbumOwnerRequired = errors.New("album must have an owner")

// newAlbumIdentifier returns a random, URL-safe album identifier carrying 128
// bits of randomness, so album links cannot be guessed.
func newAlbumIdentifier() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func albumFromRow(row *model.AlbumTable) *dm.Album {
	return &dm.Album{
		Id:          row.ID.String(),
		Identifier:  row.Identifier,
		Title:       row.Title,
		CreatedById: utils.SafeDerefWithDefault(row.CreatedByID, ZeroUUID).String(),
		CreatedAt:   row.CreatedAt,
	}
}

// CreateAlbum creates an album owned by createdById. Albums always have an
// owner, since only the owner may add images to them.
func (repo *DBImageRepo) CreateAlbum(title string, createdById string) (*dm.Album, error) {
	if utf8.RuneCountInString(title) > dm.AlbumMaxTitleLength {
		return nil, ErrAlbumTitleTooLong
	}
	ownerId, err := uuid.Parse(createdById)
	if err != nil {
		return nil, ErrAlbumOwnerRequired
	}
	identifier, err := newAlbumIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := AlbumTable.INSERT(
		AlbumTable.Identifier,
		AlbumTable.Title,
		AlbumTable.CreatedByID,
	).VALUES(
		identifier,
		title,
		UUID(ownerId),
	).RETURNING(AlbumTable.AllColumns)
	dest := model.AlbumTable{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil, err
	}
	return albumFromRow(&dest), nil
}

// GetAlbumByIdentifier returns the live album with the given public
// identifier, or nil if there is none.
func (repo *DBImageRepo) GetAlbumByIdentifier(identifier string) (*dm.Album, error) {
	if identifier == "" {
		return nil, nil
	}
	stmt := AlbumTable.
		SELECT(AlbumTable.AllColumns).
		FROM(AlbumTable).
		WHERE(AlbumTable.Identifier.EQ(String(identifier)).AND(AlbumTable.DeletedAt.IS_NULL())).
		LIMIT(1)
	dest := []model.AlbumTable{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil, err
	}
	if len(dest) == 0 {
		return nil, nil
	}
	return albumFromRow(&dest[0]), nil
}

// ListAlbumImages returns the album's live images, oldest first. Deleted and
// expired images are left out, including expired ones the cleanup task has
// not deleted yet. A positive limit keeps the newest limit images (still
// returned oldest first), so a live view keeps showing new uploads; 0 or less
// means no limit.
func (repo *DBImageRepo) ListAlbumImages(albumId string, limit int) ([]*dm.Image, error) {
	parsed, err := uuid.Parse(albumId)
	if err != nil {
		return nil, err
	}
	stmt := ImageTable.
		SELECT(ImageTable.AllColumns).
		FROM(ImageTable).
		WHERE(
			ImageTable.AlbumID.EQ(UUID(parsed)).
				AND(ImageTable.DeletedAt.IS_NULL()).
				AND(ImageTable.ExpiresAt.IS_NULL().OR(ImageTable.ExpiresAt.GT(NOW()))),
		)
	if limit > 0 {
		stmt = stmt.ORDER_BY(ImageTable.CreatedAt.DESC(), ImageTable.ID.DESC()).LIMIT(int64(limit))
	} else {
		stmt = stmt.ORDER_BY(ImageTable.CreatedAt.ASC(), ImageTable.ID.ASC())
	}
	dest := []model.ImageTable{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil, err
	}
	images := make([]*dm.Image, len(dest))
	for i := range dest {
		images[i] = imageFromRow(&dest[i])
	}
	if limit > 0 {
		slices.Reverse(images)
	}
	return images, nil
}
