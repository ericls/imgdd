//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package storage

import (
	"database/sql"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/db/.gen/imgdd/public/model"
	. "github.com/ericls/imgdd/db/.gen/imgdd/public/table"
	dm "github.com/ericls/imgdd/domainmodels"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type DBStorageDefRepo struct {
	db.RepoConn
}

func (repo *DBStorageDefRepo) WithTransaction(tx *sql.Tx) DBStorageDefRepo {
	return DBStorageDefRepo{
		RepoConn: repo.RepoConn.WithTransaction(tx),
	}
}

func NewDBStorageDefRepo(conn *sql.DB) *DBStorageDefRepo {
	return &DBStorageDefRepo{
		RepoConn: db.NewRepoConn(conn),
	}
}

func (repo *DBStorageDefRepo) GetStorageDefinitionById(id string) (*dm.StorageDefinition, error) {
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
		StorageType: dm.StorageTypeName(dest.StorageType),
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}

func (repo *DBStorageDefRepo) GetStorageDefinitionsByIds(ids []string) ([]*dm.StorageDefinition, error) {
	uuids := make([]Expression, len(ids))
	for i, id := range ids {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := SELECT(
		StorageDefinitionTable.AllColumns,
	).FROM(StorageDefinitionTable).WHERE(
		StorageDefinitionTable.ID.IN(uuids...),
	).ORDER_BY(
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
			StorageType: dm.StorageTypeName(d.StorageType),
			Config:      d.Config,
			IsEnabled:   d.IsEnabled,
			Priority:    d.Priority,
		}
	}
	return result, nil
}

func (repo *DBStorageDefRepo) GetStorageDefinitionByIdentifier(identifier string) (*dm.StorageDefinition, error) {
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
		StorageType: dm.StorageTypeName(dest.StorageType),
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}

func (repo *DBStorageDefRepo) ListStorageDefinitions() ([]*dm.StorageDefinition, error) {
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
			StorageType: dm.StorageTypeName(d.StorageType),
			Config:      d.Config,
			IsEnabled:   d.IsEnabled,
			Priority:    d.Priority,
		}
	}
	return result, nil
}

func (repo *DBStorageDefRepo) CreateStorageDefinition(storageType string, config string, identifier string, isEnabled bool, priority int64) (*dm.StorageDefinition, error) {
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
		StorageType: dm.StorageTypeName(dest.StorageType),
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}

func (repo *DBStorageDefRepo) UpdateStorageDefinition(identifier string, storage_type *string, config *string, isEnabled *bool, priority *int64) (*dm.StorageDefinition, error) {
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
		StorageType: dm.StorageTypeName(dest.StorageType),
		Config:      dest.Config,
		IsEnabled:   dest.IsEnabled,
		Priority:    dest.Priority,
	}, nil
}
