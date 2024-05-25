package db

import (
	"database/sql"
	"imgdd/internal"
)

var connectStringToDBMap = map[string]*sql.DB{}

func GetConnection(c *DBConfigDef) *sql.DB {
	connectString := c.ConnectString()
	if db, ok := connectStringToDBMap[connectString]; ok {
		return db
	}
	new_db, err := sql.Open("postgres", connectString)
	internal.PanicOnError(err)
	connectStringToDBMap[connectString] = new_db
	return new_db
}
