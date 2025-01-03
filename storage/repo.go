//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package storage

import (
	"database/sql"
	"imgdd/db"
	"imgdd/db/.gen/imgdd/public/model"
	. "imgdd/db/.gen/imgdd/public/table"
	dm "imgdd/domainmodels"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type DBStorageRepo struct {
	db.RepoConn
}

func (repo *DBStorageRepo) WithTransaction(tx *sql.Tx) DBStorageRepo {
	return DBStorageRepo{
		RepoConn: repo.RepoConn.WithTransaction(tx),
	}
}

func NewDBStorageRepo(conn *sql.DB) *DBStorageRepo {
	return &DBStorageRepo{
		RepoConn: db.NewRepoConn(conn),
	}
}

func (repo *DBStorageRepo) GetStorageDefinitionById(id string) (*dm.StorageDefinition, error) {
	stmt := SELECT(
		StorageDefinitionTable.AllColumns,
	).FROM(StorageDefinitionTable).WHERE(
		StorageDefinitionTable.ID.EQ(UUID(uuid.MustParse(id))),
	)
	dest := model.StorageDefinitionTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return &dm.StorageDefinition{
		Id:          dest.ID.String(),
		Identifier:  dest.Identifier,
		StorageType: dest.StorageType,
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}

func (repo *DBStorageRepo) GetStorageDefinitionByIdentifier(identifier string) (*dm.StorageDefinition, error) {
	stmt := SELECT(
		StorageDefinitionTable.AllColumns,
	).FROM(StorageDefinitionTable).WHERE(
		StorageDefinitionTable.Identifier.EQ(String(identifier)),
	)
	dest := model.StorageDefinitionTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return &dm.StorageDefinition{
		Id:          dest.ID.String(),
		Identifier:  dest.Identifier,
		StorageType: dest.StorageType,
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}

func (repo *DBStorageRepo) ListStorageDefinitions() ([]*dm.StorageDefinition, error) {
	stmt := SELECT(
		StorageDefinitionTable.AllColumns,
	).FROM(StorageDefinitionTable).ORDER_BY(
		StorageDefinitionTable.Priority.ASC(),
	)
	dest := []model.StorageDefinitionTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	result := make([]*dm.StorageDefinition, len(dest))
	for i, d := range dest {
		result[i] = &dm.StorageDefinition{
			Id:          d.ID.String(),
			Identifier:  d.Identifier,
			StorageType: d.StorageType,
			Config:      d.Config,
			IsEnabled:   d.IsEnabled,
			Priority:    d.Priority,
		}
	}
	return result, nil
}

func (repo *DBStorageRepo) CreateStorageDefinition(storageType string, config string, identifier string, isEnabled bool, priority int64) (*dm.StorageDefinition, error) {
	stmt := StorageDefinitionTable.INSERT(
		StorageDefinitionTable.StorageType,
		StorageDefinitionTable.Config,
		StorageDefinitionTable.Identifier,
		StorageDefinitionTable.IsEnabled,
		StorageDefinitionTable.Priority,
	).
		VALUES(
			storageType,
			config,
			identifier,
			isEnabled,
			priority,
		).RETURNING(StorageDefinitionTable.AllColumns)
	dest := model.StorageDefinitionTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return &dm.StorageDefinition{
		Id:          dest.ID.String(),
		Identifier:  dest.Identifier,
		StorageType: dest.StorageType,
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}

func (repo *DBStorageRepo) UpdateStorageDefinition(identifier string, storage_type *string, config *string, isEnabled *bool, priority *int64) (*dm.StorageDefinition, error) {
	// TODO: maybe build a wrapper for this
	updatingInput := model.StorageDefinitionTable{}
	updatingColumns := ColumnList{
		ImageTable.UpdatedAt,
	}
	if storage_type != nil {
		updatingInput.StorageType = *storage_type
		updatingColumns = append(updatingColumns, StorageDefinitionTable.StorageType)
	}
	if config != nil {
		updatingInput.Config = *config
		updatingColumns = append(updatingColumns, StorageDefinitionTable.Config)
	}
	if isEnabled != nil {
		updatingInput.IsEnabled = *isEnabled
		updatingColumns = append(updatingColumns, StorageDefinitionTable.IsEnabled)
	}
	if priority != nil {
		// XXX: Safety
		updatingInput.Priority = int32(*priority)
		updatingColumns = append(updatingColumns, StorageDefinitionTable.Priority)
	}
	stmt := StorageDefinitionTable.UPDATE(
		updatingColumns,
	).MODEL(updatingInput).WHERE(
		StorageDefinitionTable.Identifier.EQ(String(identifier)),
	).RETURNING(StorageDefinitionTable.AllColumns)
	dest := model.StorageDefinitionTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return &dm.StorageDefinition{
		Id:          dest.ID.String(),
		Identifier:  dest.Identifier,
		StorageType: dest.StorageType,
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}

func (repo *DBStorageRepo) GetStoredImageByIdentifierAndMimeType(identifier, mime string) (*dm.StoredImage, error) {
	stmt := SELECT(
		StoredImageTable.AllColumns,
		StorageDefinitionTable.AllColumns,
	).FROM(StoredImageTable.INNER_JOIN(
		ImageTable, ImageTable.ID.EQ(StoredImageTable.ImageID),
	).INNER_JOIN(
		StorageDefinitionTable, StorageDefinitionTable.ID.EQ(StoredImageTable.StorageDefinitionID),
	)).WHERE(
		ImageTable.Identifier.EQ(String(identifier)).
			AND(
				StoredImageTable.IsFileDeleted.EQ(Bool(false)),
			).
			AND(
				ImageTable.MimeType.EQ(String(mime)),
			).
			AND(
				StorageDefinitionTable.IsEnabled.EQ(Bool(true)),
			),
	).ORDER_BY(
		StorageDefinitionTable.Priority.ASC(),
	)
	dest := struct {
		StoredImageTable       model.StoredImageTable
		StorageDefinitionTable model.StorageDefinitionTable
	}{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return &dm.StoredImage{
		Id:             dest.StoredImageTable.ID.String(),
		FileIdentifier: dest.StoredImageTable.FileIdentifier,
		StorageDefinition: &dm.StorageDefinition{
			Id:          dest.StorageDefinitionTable.ID.String(),
			Identifier:  dest.StorageDefinitionTable.Identifier,
			StorageType: dest.StorageDefinitionTable.StorageType,
			Config:      dest.StorageDefinitionTable.Config,
			IsEnabled:   dest.StorageDefinitionTable.IsEnabled,
			Priority:    dest.StorageDefinitionTable.Priority,
		},
	}, nil
}

func (repo *DBStorageRepo) GetStoredImagesByIds(ids []string) ([]*dm.StoredImage, error) {
	uuids := make([]Expression, len(ids))
	for i, id := range ids {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := SELECT(
		StoredImageTable.AllColumns,
		StorageDefinitionTable.AllColumns,
	).FROM(StoredImageTable.INNER_JOIN(
		StorageDefinitionTable, StorageDefinitionTable.ID.EQ(StoredImageTable.StorageDefinitionID),
	)).WHERE(
		StoredImageTable.ID.IN(uuids...),
	).ORDER_BY(
		StorageDefinitionTable.Priority.ASC(),
	)
	dest := []struct {
		StoredImageTable       model.StoredImageTable
		StorageDefinitionTable model.StorageDefinitionTable
	}{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	result := make([]*dm.StoredImage, len(dest))
	for i, d := range dest {
		result[i] = &dm.StoredImage{
			Id:             d.StoredImageTable.ID.String(),
			FileIdentifier: d.StoredImageTable.FileIdentifier,
			StorageDefinition: &dm.StorageDefinition{
				Id:          d.StorageDefinitionTable.ID.String(),
				Identifier:  d.StorageDefinitionTable.Identifier,
				StorageType: d.StorageDefinitionTable.StorageType,
				Config:      d.StorageDefinitionTable.Config,
				IsEnabled:   d.StorageDefinitionTable.IsEnabled,
				Priority:    d.StorageDefinitionTable.Priority,
			},
		}
	}
	return result, nil
}

func (repo *DBStorageRepo) GetStoredImageIdsByImageIds(imageIds []string) (map[string][]string, error) {
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
