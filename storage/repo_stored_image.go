//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package storage

import (
	"context"
	"database/sql"
	"iter"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/db/.gen/imgdd/public/model"
	. "github.com/ericls/imgdd/db/.gen/imgdd/public/table"
	dm "github.com/ericls/imgdd/domainmodels"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type DBStoredImageRepo struct {
	db.RepoConn
}

func (repo *DBStoredImageRepo) WithTransaction(tx *sql.Tx) DBStoredImageRepo {
	return DBStoredImageRepo{
		RepoConn: repo.RepoConn.WithTransaction(tx),
	}
}

func NewDBStoredImageRepo(conn *sql.DB) *DBStoredImageRepo {
	return &DBStoredImageRepo{
		RepoConn: db.NewRepoConn(conn),
	}
}

func (repo *DBStoredImageRepo) GetStoredImageByIdentifierAndMimeType(identifier, mime string) ([]*dm.StoredImage, error) {
	stmt := SELECT(
		StoredImageTable.AllColumns,
	).FROM(StoredImageTable.INNER_JOIN(
		ImageTable, ImageTable.ID.EQ(StoredImageTable.ImageID),
	)).WHERE(
		ImageTable.Identifier.EQ(String(identifier)).
			AND(ImageTable.DeletedAt.IS_NULL()).
			AND(
				StoredImageTable.IsFileDeleted.EQ(Bool(false)),
			).
			AND(
				ImageTable.MimeType.EQ(String(mime)),
			),
	)
	dest := []struct {
		StoredImageTable model.StoredImageTable
	}{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	result := make([]*dm.StoredImage, len(dest))
	for i, d := range dest {
		result[i] = &dm.StoredImage{
			Id:                  d.StoredImageTable.ID.String(),
			FileIdentifier:      d.StoredImageTable.FileIdentifier,
			StorageDefinitionId: d.StoredImageTable.StorageDefinitionID.String(),
			IsFileDeleted:       d.StoredImageTable.IsFileDeleted,
		}
	}
	return result, nil
}

func (repo *DBStoredImageRepo) GetStoredImagesByIds(ids []string) ([]*dm.StoredImage, error) {
	uuids := make([]Expression, len(ids))
	for i, id := range ids {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := SELECT(
		StoredImageTable.AllColumns,
	).FROM(StoredImageTable).WHERE(
		StoredImageTable.ID.IN(uuids...),
	)
	dest := []struct {
		StoredImageTable model.StoredImageTable
	}{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	result := make([]*dm.StoredImage, len(dest))
	for i, d := range dest {
		result[i] = &dm.StoredImage{
			Id:                  d.StoredImageTable.ID.String(),
			FileIdentifier:      d.StoredImageTable.FileIdentifier,
			StorageDefinitionId: d.StoredImageTable.StorageDefinitionID.String(),
			IsFileDeleted:       d.StoredImageTable.IsFileDeleted,
		}
	}
	return result, nil
}

func (repo *DBStoredImageRepo) GetStoredImageIdsByImageIds(imageIds []string) (map[string][]string, error) {
	uuids := make([]Expression, len(imageIds))
	for i, id := range imageIds {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := SELECT(
		StoredImageTable.ID,
		StoredImageTable.ImageID,
	).FROM(StoredImageTable).WHERE(
		StoredImageTable.ImageID.IN(uuids...),
	)
	dest := []struct {
		model.StoredImageTable
	}{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string)
	for _, d := range dest {
		if _, ok := result[d.ID.String()]; !ok {
			result[d.ImageID.String()] = []string{}
		}
		result[d.ImageID.String()] = append(result[d.ImageID.String()], d.ID.String())
	}
	return result, nil
}

func (repo *DBStoredImageRepo) GetStoredImagesToDelete() ([]*dm.StoredImage, error) {
	stmt := SELECT(
		StoredImageTable.AllColumns,
	).FROM(
		StoredImageTable.LEFT_JOIN(ImageTable, StoredImageTable.ImageID.EQ(ImageTable.ID)),
	).WHERE(
		StoredImageTable.IsFileDeleted.EQ(Bool(false)).
			AND(ImageTable.DeletedAt.IS_NOT_NULL()),
	)
	dest := []struct {
		StoredImageTable model.StoredImageTable
	}{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	result := make([]*dm.StoredImage, len(dest))
	for i, d := range dest {
		result[i] = &dm.StoredImage{
			Id:                  d.StoredImageTable.ID.String(),
			FileIdentifier:      d.StoredImageTable.FileIdentifier,
			StorageDefinitionId: d.StoredImageTable.StorageDefinitionID.String(),
			IsFileDeleted:       d.StoredImageTable.IsFileDeleted,
		}
	}
	return result, nil
}

func storedImageFromJoinRow(si model.StoredImageTable, img model.ImageTable) *dm.StoredImage {
	return &dm.StoredImage{
		Id:                  si.ID.String(),
		FileIdentifier:      si.FileIdentifier,
		StorageDefinitionId: si.StorageDefinitionID.String(),
		IsFileDeleted:       si.IsFileDeleted,
		Image: &dm.Image{
			Id:       img.ID.String(),
			MIMEType: img.MimeType,
		},
	}
}

func (repo *DBStoredImageRepo) GetStoredImagesByImageId(imageId string) ([]*dm.StoredImage, error) {
	stmt := SELECT(
		StoredImageTable.AllColumns,
		ImageTable.MimeType,
		ImageTable.ID,
	).FROM(
		StoredImageTable.INNER_JOIN(ImageTable, StoredImageTable.ImageID.EQ(ImageTable.ID)),
	).WHERE(
		StoredImageTable.ImageID.EQ(UUID(uuid.MustParse(imageId))).
			AND(StoredImageTable.IsFileDeleted.EQ(Bool(false))).
			AND(ImageTable.DeletedAt.IS_NULL()),
	)
	dest := []struct {
		StoredImageTable model.StoredImageTable
		ImageTable       model.ImageTable
	}{}
	if err := stmt.Query(repo.DB, &dest); err != nil {
		return nil, err
	}
	result := make([]*dm.StoredImage, len(dest))
	for i, d := range dest {
		result[i] = storedImageFromJoinRow(d.StoredImageTable, d.ImageTable)
	}
	return result, nil
}

func (repo *DBStoredImageRepo) GetStoredImagesForReplication(sourceStorageDefId string, targetStorageDefId string) (int, iter.Seq2[*dm.StoredImage, error], error) {
	sourceUUID := UUID(uuid.MustParse(sourceStorageDefId))
	targetUUID := UUID(uuid.MustParse(targetStorageDefId))

	targetSI := StoredImageTable.AS("target_si")
	antiJoinCondition := StoredImageTable.IsFileDeleted.EQ(Bool(false)).
		AND(StoredImageTable.StorageDefinitionID.EQ(sourceUUID)).
		AND(ImageTable.DeletedAt.IS_NULL()).
		AND(targetSI.ID.IS_NULL())
	antiJoinFrom := StoredImageTable.
		INNER_JOIN(ImageTable, StoredImageTable.ImageID.EQ(ImageTable.ID)).
		LEFT_JOIN(targetSI,
			targetSI.ImageID.EQ(StoredImageTable.ImageID).
				AND(targetSI.StorageDefinitionID.EQ(targetUUID)).
				AND(targetSI.IsFileDeleted.EQ(Bool(false))))

	countDest := []struct{ Count int }{}
	if err := SELECT(COUNT(StoredImageTable.ID).AS("count")).
		FROM(antiJoinFrom).
		WHERE(antiJoinCondition).
		Query(repo.DB, &countDest); err != nil {
		return 0, nil, err
	}
	total := 0
	if len(countDest) > 0 {
		total = countDest[0].Count
	}

	seq := func(yield func(*dm.StoredImage, error) bool) {
		query, args := SELECT(
			StoredImageTable.AllColumns,
			ImageTable.MimeType,
			ImageTable.ID,
		).FROM(antiJoinFrom).WHERE(antiJoinCondition).Sql()

		rows, err := repo.DB.QueryContext(context.Background(), query, args...)
		if err != nil {
			yield(nil, err)
			return
		}
		defer rows.Close()

		scanCtx, err := qrm.NewScanContext(rows)
		if err != nil {
			yield(nil, err)
			return
		}
		for rows.Next() {
			var dest struct {
				StoredImageTable model.StoredImageTable
				ImageTable       model.ImageTable
			}
			if err := qrm.ScanOneRowToDest(scanCtx, rows, &dest); err != nil {
				yield(nil, err)
				return
			}
			if !yield(storedImageFromJoinRow(dest.StoredImageTable, dest.ImageTable), nil) {
				return
			}
		}
		if err := rows.Err(); err != nil {
			yield(nil, err)
		}
	}
	return total, seq, nil
}

func (repo *DBStoredImageRepo) MarkStoredImagesAsDeleted(ids []string) error {
	uuids := make([]Expression, len(ids))
	for i, id := range ids {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := StoredImageTable.UPDATE().SET(
		StoredImageTable.IsFileDeleted.SET(Bool(true)),
		StoredImageTable.UpdatedAt.SET(TimestampzExp(Func("NOW"))),
	).WHERE(
		StoredImageTable.ID.IN(uuids...),
	)
	_, err := stmt.Exec(repo.DB)
	return err
}
