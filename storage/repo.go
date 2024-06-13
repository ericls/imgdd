package storage

import (
	"database/sql"
	"imgdd/db"
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
