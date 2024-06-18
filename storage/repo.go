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
		Id:         dest.ID.String(),
		Identifier: dest.Identifier,
		Type:       dest.Type,
		Config:     dest.Config,
		IsEnabled:  dest.IsEnabled,
		Priority:   dest.Priority,
	}, nil
}
