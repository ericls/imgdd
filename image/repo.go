//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package image

import (
	"bytes"
	"database/sql"
	"imgdd/db"
	"imgdd/db/.gen/imgdd/public/model"
	. "imgdd/db/.gen/imgdd/public/table"
	dm "imgdd/domainmodels"
	"imgdd/logging"
	"imgdd/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

var logger = logging.GetLogger("image-repo")

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

func hasPrevImageConditions(ordering ListImagesOrdering, image *dm.Image) []BoolExpression {
	var conditions []BoolExpression
	if ordering.ID != nil {
		if *ordering.ID == PaginationDirectionAsc {
			conditions = append(conditions, ImageTable.ID.LT(UUID(uuid.MustParse(image.Id))))
		} else if *ordering.ID == PaginationDirectionDesc {
			conditions = append(conditions, ImageTable.ID.GT(UUID(uuid.MustParse(image.Id))))
		}
	}
	if ordering.Name != nil {
		if *ordering.Name == PaginationDirectionAsc {
			conditions = append(conditions, ImageTable.Name.LT(String(image.Name)))
		} else if *ordering.Name == PaginationDirectionDesc {
			conditions = append(conditions, ImageTable.Name.GT(String(image.Name)))
		}
	}
	if ordering.CreatedAt != nil {
		if *ordering.CreatedAt == PaginationDirectionAsc {
			conditions = append(conditions, ImageTable.CreatedAt.LT(TimestampzT(image.CreatedAt)))
		} else if *ordering.CreatedAt == PaginationDirectionDesc {
			conditions = append(conditions, ImageTable.CreatedAt.GT(TimestampzT(image.CreatedAt)))
		}
	}
	return conditions
}

func hasNextImageConditions(ordering ListImagesOrdering, image *dm.Image) []BoolExpression {
	var conditions []BoolExpression
	if ordering.ID != nil {
		if *ordering.ID == PaginationDirectionAsc {
			conditions = append(conditions, ImageTable.ID.GT(UUID(uuid.MustParse(image.Id))))
		} else if *ordering.ID == PaginationDirectionDesc {
			conditions = append(conditions, ImageTable.ID.LT(UUID(uuid.MustParse(image.Id))))
		}
	}
	if ordering.Name != nil {
		if *ordering.Name == PaginationDirectionAsc {
			conditions = append(conditions, ImageTable.Name.GT(String(image.Name)))
		} else if *ordering.Name == PaginationDirectionDesc {
			conditions = append(conditions, ImageTable.Name.LT(String(image.Name)))
		}
	}
	if ordering.CreatedAt != nil {
		if *ordering.CreatedAt == PaginationDirectionAsc {
			conditions = append(conditions, ImageTable.CreatedAt.GT(TimestampzT(image.CreatedAt)))
		} else if *ordering.CreatedAt == PaginationDirectionDesc {
			conditions = append(conditions, ImageTable.CreatedAt.LT(TimestampzT(image.CreatedAt)))
		}
	}
	return conditions
}

func (repo *DBImageRepo) imageExists(conditions []BoolExpression) bool {
	statement := ImageTable.SELECT(
		ImageTable.ID,
	).FROM(
		ImageTable,
	)
	for _, condition := range conditions {
		statement = statement.WHERE(condition)
	}
	statement = statement.LIMIT(1)
	dest := []model.ImageTable{}
	err := statement.Query(repo.DB, &dest)
	if err != nil {
		logger.Err(err).Msg("Error checking if image exists")
		return false
	}
	return len(dest) > 0
}

func (repo *DBImageRepo) imageHasPrev(ordering ListImagesOrdering, image *dm.Image) bool {
	conditions := hasPrevImageConditions(ordering, image)
	return repo.imageExists(conditions)
}

func (repo *DBImageRepo) imageHasNext(ordering ListImagesOrdering, image *dm.Image) bool {
	conditions := hasNextImageConditions(ordering, image)
	return repo.imageExists(conditions)
}

func (repo *DBImageRepo) ListImages(filters ListImagesFilters, ordering ListImagesOrdering) (dm.ListImageResult, error) {
	stmt := ImageTable.SELECT(
		ImageTable.AllColumns,
	).FROM(
		ImageTable,
	)
	if filters.CreatedAtGte != nil {
		stmt = stmt.WHERE(ImageTable.CreatedAt.GT_EQ(TimestampzT(*filters.CreatedAtGte)))
	}
	if filters.CreatedAtLte != nil {
		stmt = stmt.WHERE(ImageTable.CreatedAt.LT_EQ(TimestampzT(*filters.CreatedAtLte)))
	}
	if filters.NameContains != "" {
		stmt = stmt.WHERE(ImageTable.Name.LIKE(String("%" + filters.NameContains + "%")))
	}
	if filters.CreatedBy != nil {
		stmt = stmt.WHERE(ImageTable.CreatedByID.EQ(UUID(uuid.MustParse(*filters.CreatedBy))))
	}
	orderByClauses := []OrderByClause{}
	if ordering.ID != nil {
		if *ordering.ID == PaginationDirectionAsc {
			orderByClauses = append(orderByClauses, ImageTable.ID.ASC())
		} else if *ordering.ID == PaginationDirectionDesc {
			orderByClauses = append(orderByClauses, ImageTable.ID.DESC())
		}
	}
	if ordering.Name != nil {
		if *ordering.Name == PaginationDirectionAsc {
			orderByClauses = append(orderByClauses, ImageTable.Name.ASC())
		} else if *ordering.Name == PaginationDirectionDesc {
			orderByClauses = append(orderByClauses, ImageTable.Name.DESC())
		}
	}
	if ordering.CreatedAt != nil {
		if *ordering.CreatedAt == PaginationDirectionAsc {
			orderByClauses = append(orderByClauses, ImageTable.CreatedAt.ASC())
		} else if *ordering.CreatedAt == PaginationDirectionDesc {
			orderByClauses = append(orderByClauses, ImageTable.CreatedAt.DESC())
		}
	}
	if len(orderByClauses) > 0 {
		stmt = stmt.ORDER_BY(orderByClauses...)
	}
	if filters.Limit > 0 {
		stmt = stmt.LIMIT(int64(filters.Limit))
	}
	dest := []model.ImageTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return dm.ListImageResult{}, err
	}
	images := make([]*dm.Image, len(dest))
	for i, image := range dest {
		images[i] = &dm.Image{
			Id:              image.ID.String(),
			Identifier:      image.Identifier,
			CreatedAt:       image.CreatedAt,
			Name:            image.Name,
			ParentId:        utils.SafeDeref(image.ParentID).String(),
			RootId:          utils.SafeDeref(image.RootID).String(),
			UploaderIP:      utils.SafeDeref(image.UploaderIP),
			MIMEType:        image.MimeType,
			NominalWidth:    image.NominalWidth,
			NominalHeight:   image.NominalHeight,
			NominalByteSize: image.NominalByteSize,
		}
	}
	firstImage := images[0]
	lastImage := images[len(dest)-1]

	return dm.ListImageResult{
		Images:  images,
		HasNext: repo.imageHasNext(ordering, lastImage),
		HasPrev: repo.imageHasPrev(ordering, firstImage),
	}, nil
}
