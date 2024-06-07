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
