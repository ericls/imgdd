package storage

import (
	"database/sql"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/utils"
)

type StorageDefSource string

const (
	StorageDefSourceDB   StorageDefSource = "db"
	StorageDefSourceConf StorageDefSource = "conf"
)

type StorageConfigDef struct {
	StorageDefSource StorageDefSource
	StorageDefs      []dm.StorageDefinition
	Conn             *utils.Lazy[*sql.DB]
}

func (c *StorageConfigDef) MakeStorageDefRepo() StorageDefRepo {
	if c.StorageDefSource == StorageDefSourceDB {
		return NewDBStorageDefRepo(c.Conn.Value())
	}
	repo := NewInMemoryStorageDefRepo()
	for _, def := range c.StorageDefs {
		repo.AddStorageDefinition(&def)
	}
	return repo
}

func NewDBStorageConfig(conn *sql.DB) *StorageConfigDef {
	return &StorageConfigDef{
		StorageDefSource: StorageDefSourceDB,
		Conn:             utils.NewLazy(func() *sql.DB { return conn }),
	}
}

func NewConfStorageConfig(storageDefs []dm.StorageDefinition) *StorageConfigDef {
	return &StorageConfigDef{
		StorageDefSource: StorageDefSourceConf,
		StorageDefs:      storageDefs,
	}
}
