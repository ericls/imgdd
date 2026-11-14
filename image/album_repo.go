//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package image

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"slices"
	"unicode/utf8"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/db/.gen/imgdd/public/model"
	. "github.com/ericls/imgdd/db/.gen/imgdd/public/table"
	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

var ErrAlbumTitleTooLong = errors.New("album title is too long")
var ErrAlbumOwnerRequired = errors.New("album must have an owner")
var ErrCaptionTooLong = errors.New("caption is too long")

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

// AddImagesToAlbum adds images to an album with their captions there, as
// addedById. Images already in the album are left as they are, so adding is
// idempotent. New images go to the end of the album in the order given. It
// returns how many images were newly added. Permission checks are the
// caller's job; see dm.Album.
func (repo *DBImageRepo) AddImagesToAlbum(albumId string, entries []dm.AlbumImageEntry, addedById string) (int, error) {
	if len(entries) == 0 {
		return 0, nil
	}
	parsedAlbumId, err := uuid.Parse(albumId)
	if err != nil {
		return 0, err
	}
	var addedBy *uuid.UUID
	if parsed, err := uuid.Parse(addedById); err == nil {
		addedBy = &parsed
	}
	rows := make([]model.AlbumImageTable, len(entries))
	for i, entry := range entries {
		if utf8.RuneCountInString(entry.Caption) > dm.ImageMaxCaptionLength {
			return 0, ErrCaptionTooLong
		}
		imageId, err := uuid.Parse(entry.ImageId)
		if err != nil {
			return 0, err
		}
		rows[i] = model.AlbumImageTable{AlbumID: parsedAlbumId, ImageID: imageId, AddedByID: addedBy}
		if entry.Caption != "" {
			caption := entry.Caption
			rows[i].Caption = &caption
		}
	}
	stmt := AlbumImageTable.INSERT(
		AlbumImageTable.AlbumID,
		AlbumImageTable.ImageID,
		AlbumImageTable.Caption,
		AlbumImageTable.AddedByID,
	).MODELS(rows).ON_CONFLICT(AlbumImageTable.AlbumID, AlbumImageTable.ImageID).DO_NOTHING()
	result, err := stmt.Exec(repo.DB)
	if err != nil {
		return 0, err
	}
	added, err := result.RowsAffected()
	return int(added), err
}

// RemoveImagesFromAlbum removes images from an album. The images themselves
// are not touched. It returns how many images were in the album and removed.
func (repo *DBImageRepo) RemoveImagesFromAlbum(albumId string, imageIds []string) (int, error) {
	if len(imageIds) == 0 {
		return 0, nil
	}
	parsedAlbumId, err := uuid.Parse(albumId)
	if err != nil {
		return 0, err
	}
	ids := make([]Expression, len(imageIds))
	for i, id := range imageIds {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return 0, err
		}
		ids[i] = UUID(parsed)
	}
	stmt := AlbumImageTable.DELETE().WHERE(
		AlbumImageTable.AlbumID.EQ(UUID(parsedAlbumId)).AND(AlbumImageTable.ImageID.IN(ids...)),
	)
	result, err := stmt.Exec(repo.DB)
	if err != nil {
		return 0, err
	}
	removed, err := result.RowsAffected()
	return int(removed), err
}

// CreateAndSaveUploadedImageInAlbum stores a new upload and adds it to an
// album in one transaction, so an upload never ends up outside the album it
// was meant for.
func (repo *DBImageRepo) CreateAndSaveUploadedImageInAlbum(image *dm.Image, mimeType string, fileBytes []byte, storageDefinitionId string, saveFn SaveFunc, albumId string, caption string) (*dm.StoredImage, error) {
	return db.RunInTransaction(repo, func(txRepo *DBImageRepo) (*dm.StoredImage, error) {
		storedImage, err := txRepo.CreateAndSaveUploadedImage(image, mimeType, fileBytes, storageDefinitionId, saveFn)
		if err != nil {
			return nil, err
		}
		entry := dm.AlbumImageEntry{ImageId: storedImage.Image.Id, Caption: caption}
		if _, err := txRepo.AddImagesToAlbum(albumId, []dm.AlbumImageEntry{entry}, image.CreatedById); err != nil {
			return nil, err
		}
		return storedImage, nil
	})
}

type albumImageRow struct {
	model.AlbumImageTable
	Image model.ImageTable
}

// ListAlbumImages returns the album's live images in the order they were
// added, oldest first. Deleted and expired images are left out, including
// expired ones the cleanup task has not deleted yet. A positive limit keeps
// the most recently added limit images (still returned oldest first), so a
// live view keeps showing new additions; 0 or less means no limit.
func (repo *DBImageRepo) ListAlbumImages(albumId string, limit int) ([]*dm.AlbumImage, error) {
	parsed, err := uuid.Parse(albumId)
	if err != nil {
		return nil, err
	}
	stmt := SELECT(AlbumImageTable.AllColumns, ImageTable.AllColumns).
		FROM(AlbumImageTable.INNER_JOIN(ImageTable, ImageTable.ID.EQ(AlbumImageTable.ImageID))).
		WHERE(
			AlbumImageTable.AlbumID.EQ(UUID(parsed)).
				AND(ImageTable.DeletedAt.IS_NULL()).
				AND(ImageTable.ExpiresAt.IS_NULL().OR(ImageTable.ExpiresAt.GT(NOW()))),
		)
	if limit > 0 {
		stmt = stmt.ORDER_BY(AlbumImageTable.Position.DESC()).LIMIT(int64(limit))
	} else {
		stmt = stmt.ORDER_BY(AlbumImageTable.Position.ASC())
	}
	dest := []albumImageRow{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil, err
	}
	images := make([]*dm.AlbumImage, len(dest))
	for i := range dest {
		images[i] = &dm.AlbumImage{
			Image:   imageFromRow(&dest[i].Image),
			Caption: utils.SafeDeref(dest[i].Caption),
			AddedAt: dest[i].AddedAt,
		}
	}
	if limit > 0 {
		slices.Reverse(images)
	}
	return images, nil
}
