package db

import (
	"database/sql"

	"github.com/go-jet/jet/v2/qrm"
)

type DBRepo interface {
	GetDB() qrm.DB
	GetConn() *sql.DB
	GetIsInTransaction() bool
	WithTransaction(tx *sql.Tx) DBRepo
}

func RunInTransaction[Ret any, Repo DBRepo](repo Repo, fn func(txRepo Repo) (Ret, error)) (Ret, error) {
	// TODO: Implement savepoints
	if repo.GetIsInTransaction() {
		return fn(repo)
	}
	tx, err := repo.GetConn().Begin()
	var empty Ret
	if err != nil {
		return empty, err
	}
	txRepo := repo.WithTransaction(tx)
	ret, err := fn(txRepo.(Repo))
	if err != nil {
		tx.Rollback()
		return empty, err
	}
	err = tx.Commit()
	if err != nil {
		return ret, err
	}
	return ret, nil
}

type RepoConn struct {
	DB              qrm.DB
	Conn            *sql.DB
	isInTransaction bool
}

func NewRepoConn(conn *sql.DB) RepoConn {
	return RepoConn{
		DB:              conn,
		Conn:            conn,
		isInTransaction: false,
	}
}

// Implmenting the DBRepo interface
func (repo *RepoConn) GetDB() qrm.DB {
	return repo.DB
}

func (repo *RepoConn) GetConn() *sql.DB {
	return repo.Conn
}

func (repo *RepoConn) GetIsInTransaction() bool {
	return repo.isInTransaction
}

func (repo *RepoConn) WithTransaction(tx *sql.Tx) RepoConn {
	return RepoConn{
		DB:              tx,
		Conn:            repo.Conn,
		isInTransaction: true,
	}
}
