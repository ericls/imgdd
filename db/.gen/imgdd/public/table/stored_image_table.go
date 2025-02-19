//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var StoredImageTable = newStoredImageTableTable("public", "stored_image_table", "")

type storedImageTableTable struct {
	postgres.Table

	// Columns
	ID                  postgres.ColumnString
	ImageID             postgres.ColumnString
	StorageDefinitionID postgres.ColumnString
	FileIdentifier      postgres.ColumnString
	CopiedFromID        postgres.ColumnString
	CreatedAt           postgres.ColumnTimestampz
	UpdatedAt           postgres.ColumnTimestampz
	IsFileDeleted       postgres.ColumnBool

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type StoredImageTableTable struct {
	storedImageTableTable

	EXCLUDED storedImageTableTable
}

// AS creates new StoredImageTableTable with assigned alias
func (a StoredImageTableTable) AS(alias string) *StoredImageTableTable {
	return newStoredImageTableTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new StoredImageTableTable with assigned schema name
func (a StoredImageTableTable) FromSchema(schemaName string) *StoredImageTableTable {
	return newStoredImageTableTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new StoredImageTableTable with assigned table prefix
func (a StoredImageTableTable) WithPrefix(prefix string) *StoredImageTableTable {
	return newStoredImageTableTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new StoredImageTableTable with assigned table suffix
func (a StoredImageTableTable) WithSuffix(suffix string) *StoredImageTableTable {
	return newStoredImageTableTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newStoredImageTableTable(schemaName, tableName, alias string) *StoredImageTableTable {
	return &StoredImageTableTable{
		storedImageTableTable: newStoredImageTableTableImpl(schemaName, tableName, alias),
		EXCLUDED:              newStoredImageTableTableImpl("", "excluded", ""),
	}
}

func newStoredImageTableTableImpl(schemaName, tableName, alias string) storedImageTableTable {
	var (
		IDColumn                  = postgres.StringColumn("id")
		ImageIDColumn             = postgres.StringColumn("image_id")
		StorageDefinitionIDColumn = postgres.StringColumn("storage_definition_id")
		FileIdentifierColumn      = postgres.StringColumn("file_identifier")
		CopiedFromIDColumn        = postgres.StringColumn("copied_from_id")
		CreatedAtColumn           = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn           = postgres.TimestampzColumn("updated_at")
		IsFileDeletedColumn       = postgres.BoolColumn("is_file_deleted")
		allColumns                = postgres.ColumnList{IDColumn, ImageIDColumn, StorageDefinitionIDColumn, FileIdentifierColumn, CopiedFromIDColumn, CreatedAtColumn, UpdatedAtColumn, IsFileDeletedColumn}
		mutableColumns            = postgres.ColumnList{ImageIDColumn, StorageDefinitionIDColumn, FileIdentifierColumn, CopiedFromIDColumn, CreatedAtColumn, UpdatedAtColumn, IsFileDeletedColumn}
	)

	return storedImageTableTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                  IDColumn,
		ImageID:             ImageIDColumn,
		StorageDefinitionID: StorageDefinitionIDColumn,
		FileIdentifier:      FileIdentifierColumn,
		CopiedFromID:        CopiedFromIDColumn,
		CreatedAt:           CreatedAtColumn,
		UpdatedAt:           UpdatedAtColumn,
		IsFileDeleted:       IsFileDeletedColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
