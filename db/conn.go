package db

import (
	"context"
	"database/sql"

	"github.com/ericls/imgdd/internal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var connectStringToDBMap = map[string]*sql.DB{}

func GetConnection(c *DBConfigDef) *sql.DB {
	connectString := c.ConnectString()
	if db, ok := connectStringToDBMap[connectString]; ok {
		return db
	}
	poolCfg, err := pgxpool.ParseConfig(connectString)
	poolCfg.MinConns = 2
	internal.PanicOnError(err)
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	internal.PanicOnError(err)
	conn := stdlib.OpenDBFromPool(pool)
	connectStringToDBMap[connectString] = conn
	return conn
}
