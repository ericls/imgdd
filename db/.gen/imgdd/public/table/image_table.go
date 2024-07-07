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

var ImageTable = newImageTableTable("public", "image_table", "")

type imageTableTable struct {
	postgres.Table

	// Columns
	ID              postgres.ColumnString
	CreatedByID     postgres.ColumnString
	Name            postgres.ColumnString
	Identifier      postgres.ColumnString
	RootID          postgres.ColumnString
	ParentID        postgres.ColumnString
	Changes         postgres.ColumnString
	UploaderIP      postgres.ColumnString
	CreatedAt       postgres.ColumnTimestampz
	UpdatedAt       postgres.ColumnTimestampz
	DeletedAt       postgres.ColumnTimestampz
	NominalWidth    postgres.ColumnInteger
	NominalHeight   postgres.ColumnInteger
	NominalByteSize postgres.ColumnInteger
	MimeType        postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ImageTableTable struct {
	imageTableTable

	EXCLUDED imageTableTable
}

// AS creates new ImageTableTable with assigned alias
func (a ImageTableTable) AS(alias string) *ImageTableTable {
	return newImageTableTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ImageTableTable with assigned schema name
func (a ImageTableTable) FromSchema(schemaName string) *ImageTableTable {
	return newImageTableTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ImageTableTable with assigned table prefix
func (a ImageTableTable) WithPrefix(prefix string) *ImageTableTable {
	return newImageTableTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ImageTableTable with assigned table suffix
func (a ImageTableTable) WithSuffix(suffix string) *ImageTableTable {
	return newImageTableTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newImageTableTable(schemaName, tableName, alias string) *ImageTableTable {
	return &ImageTableTable{
		imageTableTable: newImageTableTableImpl(schemaName, tableName, alias),
		EXCLUDED:        newImageTableTableImpl("", "excluded", ""),
	}
}

func newImageTableTableImpl(schemaName, tableName, alias string) imageTableTable {
	var (
		IDColumn              = postgres.StringColumn("id")
		CreatedByIDColumn     = postgres.StringColumn("created_by_id")
		NameColumn            = postgres.StringColumn("name")
		IdentifierColumn      = postgres.StringColumn("identifier")
		RootIDColumn          = postgres.StringColumn("root_id")
		ParentIDColumn        = postgres.StringColumn("parent_id")
		ChangesColumn         = postgres.StringColumn("changes")
		UploaderIPColumn      = postgres.StringColumn("uploader_ip")
		CreatedAtColumn       = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn       = postgres.TimestampzColumn("updated_at")
		DeletedAtColumn       = postgres.TimestampzColumn("deleted_at")
		NominalWidthColumn    = postgres.IntegerColumn("nominal_width")
		NominalHeightColumn   = postgres.IntegerColumn("nominal_height")
		NominalByteSizeColumn = postgres.IntegerColumn("nominal_byte_size")
		MimeTypeColumn        = postgres.StringColumn("mime_type")
		allColumns            = postgres.ColumnList{IDColumn, CreatedByIDColumn, NameColumn, IdentifierColumn, RootIDColumn, ParentIDColumn, ChangesColumn, UploaderIPColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, NominalWidthColumn, NominalHeightColumn, NominalByteSizeColumn, MimeTypeColumn}
		mutableColumns        = postgres.ColumnList{CreatedByIDColumn, NameColumn, IdentifierColumn, RootIDColumn, ParentIDColumn, ChangesColumn, UploaderIPColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn, NominalWidthColumn, NominalHeightColumn, NominalByteSizeColumn, MimeTypeColumn}
	)

	return imageTableTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:              IDColumn,
		CreatedByID:     CreatedByIDColumn,
		Name:            NameColumn,
		Identifier:      IdentifierColumn,
		RootID:          RootIDColumn,
		ParentID:        ParentIDColumn,
		Changes:         ChangesColumn,
		UploaderIP:      UploaderIPColumn,
		CreatedAt:       CreatedAtColumn,
		UpdatedAt:       UpdatedAtColumn,
		DeletedAt:       DeletedAtColumn,
		NominalWidth:    NominalWidthColumn,
		NominalHeight:   NominalHeightColumn,
		NominalByteSize: NominalByteSizeColumn,
		MimeType:        MimeTypeColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
