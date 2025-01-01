//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package image

import (
	"bytes"
	"database/sql"
	"imgdd/db"
	"imgdd/db/.gen/imgdd/public/model"
	. "imgdd/db/.gen/imgdd/public/table"
	dm "imgdd/domainmodels"
	"imgdd/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type DBImageRepo struct {
	db.RepoConn
}

func (repo *DBImageRepo) WithTransaction(tx *sql.Tx) db.DBRepo {
	return &DBImageRepo{
		RepoConn: repo.RepoConn.WithTransaction(tx),
	}
}

func NewDBImageRepo(conn *sql.DB) *DBImageRepo {
	return &DBImageRepo{
		RepoConn: db.NewRepoConn(conn),
	}
}

func (repo *DBImageRepo) GetImageById(id string) (*dm.Image, error) {
	stmt := ImageTable.
		SELECT(
			ImageTable.AllColumns,
		).
		FROM(
			ImageTable,
		).
		WHERE(
			ImageTable.ID.EQ(UUID(uuid.MustParse(id))),
		)

	dest := model.ImageTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}

	var parentId string
	if dest.ParentID != nil {
		parentId = dest.ParentID.String()
	}
	var rootId string
	if dest.RootID != nil {
		rootId = dest.RootID.String()
	}

	return &dm.Image{
		Id:              dest.ID.String(),
		Identifier:      dest.Identifier,
		CreatedAt:       dest.CreatedAt,
		Name:            dest.Name,
		ParentId:        parentId,
		RootId:          rootId,
		UploaderIP:      utils.SafeDeref(dest.UploaderIP),
		MIMEType:        dest.MimeType,
		NominalWidth:    dest.NominalWidth,
		NominalHeight:   dest.NominalHeight,
		NominalByteSize: dest.NominalByteSize,
	}, nil
}

func (repo *DBImageRepo) CreateImage(image *dm.Image) (*dm.Image, error) {
	var parentId *string
	var rootId *string
	var createdById *string

	if image.ParentId != "" {
		parent, err := repo.GetImageById(image.ParentId)
		if err != nil || parent == nil {
			return nil, err
		}
		parentId = &parent.Id
		if parent.RootId != "" {
			rootId = &parent.RootId
		} else {
			rootId = &parent.Id
		}
	}

	if image.CreatedById != "" {
		createdById = &image.CreatedById
	} else {
		createdById = nil
	}

	stmt := ImageTable.INSERT(
		ImageTable.Identifier,
		ImageTable.Name,
		ImageTable.ParentID,
		ImageTable.RootID,
		ImageTable.UploaderIP,
		ImageTable.CreatedByID,
		ImageTable.MimeType,
		ImageTable.NominalByteSize,
		ImageTable.NominalWidth,
		ImageTable.NominalHeight,
	).VALUES(
		image.Identifier,
		image.Name,
		parentId,
		rootId,
		image.UploaderIP,
		createdById,
		image.MIMEType,
		image.NominalByteSize,
		image.NominalWidth,
		image.NominalHeight,
	).RETURNING(
		ImageTable.AllColumns,
	)

	dest := model.ImageTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return repo.GetImageById(dest.ID.String())
}

func (repo *DBImageRepo) CreateStoredImage(imageId string, storageDefinitionId string, fileIdentifier string, copiedFromId *string) (*dm.StoredImage, error) {
	var copiedFromIDValue StringExpression
	if copiedFromId != nil {
		copiedFromIDValue = UUID(uuid.MustParse(*copiedFromId))
	}
	stmt := StoredImageTable.INSERT(
		StoredImageTable.ImageID,
		StoredImageTable.FileIdentifier,
		StoredImageTable.StorageDefinitionID,
		StoredImageTable.CopiedFromID,
	).VALUES(
		UUID(uuid.MustParse(imageId)),
		fileIdentifier,
		UUID(uuid.MustParse(storageDefinitionId)),
		copiedFromIDValue,
	).RETURNING(
		StoredImageTable.AllColumns,
	)

	dest := model.StoredImageTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	image, err := repo.GetImageById(dest.ImageID.String())
	if err != nil {
		return nil, err
	}
	return &dm.StoredImage{
		Id:    dest.ID.String(),
		Image: image,
	}, nil
}

func (repo *DBImageRepo) CreateAndSaveUploadedImage(image *dm.Image, fileBytes []byte, storageDefinitionId string, saveFn SaveFunc) (*dm.StoredImage, error) {
	return db.RunInTransaction(repo, func(txRepo *DBImageRepo) (*dm.StoredImage, error) {
		image, err := txRepo.CreateImage(image)
		if err != nil {
			return nil, err
		}
		fileIdentifier := uuid.New().String()
		reader := bytes.NewReader(fileBytes)
		err = saveFn(reader, fileIdentifier, image.MIMEType)
		if err != nil {
			return nil, err
		}
		storedImage, err := txRepo.CreateStoredImage(image.Id, storageDefinitionId, fileIdentifier, nil)
		return storedImage, err
	})
}

func (repo *DBImageRepo) ListImages(filters ListImagesFilters, ordering ListImagesOrdering) ([]*dm.Image, error) {
	return []*dm.Image{}, nil
}
