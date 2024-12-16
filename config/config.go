package config

import (
	"imgdd/db"
	"imgdd/httpserver"
)

type ConfigDef struct {
	db         *db.DBConfigDef
	httpServer *httpserver.HttpServerConfigDef
}
