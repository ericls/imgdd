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

var RoleTable = newRoleTableTable("public", "role_table", "")

type roleTableTable struct {
	postgres.Table

	// Columns
	ID             postgres.ColumnString
	Key            postgres.ColumnString
	OrganizationID postgres.ColumnString
	DisplayName    postgres.ColumnString
	ExtraAttrs     postgres.ColumnString
	CreatedAt      postgres.ColumnTimestampz
	UpdatedAt      postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type RoleTableTable struct {
	roleTableTable

	EXCLUDED roleTableTable
}

// AS creates new RoleTableTable with assigned alias
func (a RoleTableTable) AS(alias string) *RoleTableTable {
	return newRoleTableTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new RoleTableTable with assigned schema name
func (a RoleTableTable) FromSchema(schemaName string) *RoleTableTable {
	return newRoleTableTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new RoleTableTable with assigned table prefix
func (a RoleTableTable) WithPrefix(prefix string) *RoleTableTable {
	return newRoleTableTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new RoleTableTable with assigned table suffix
func (a RoleTableTable) WithSuffix(suffix string) *RoleTableTable {
	return newRoleTableTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newRoleTableTable(schemaName, tableName, alias string) *RoleTableTable {
	return &RoleTableTable{
		roleTableTable: newRoleTableTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newRoleTableTableImpl("", "excluded", ""),
	}
}

func newRoleTableTableImpl(schemaName, tableName, alias string) roleTableTable {
	var (
		IDColumn             = postgres.StringColumn("id")
		KeyColumn            = postgres.StringColumn("key")
		OrganizationIDColumn = postgres.StringColumn("organization_id")
		DisplayNameColumn    = postgres.StringColumn("display_name")
		ExtraAttrsColumn     = postgres.StringColumn("extra_attrs")
		CreatedAtColumn      = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn      = postgres.TimestampzColumn("updated_at")
		allColumns           = postgres.ColumnList{IDColumn, KeyColumn, OrganizationIDColumn, DisplayNameColumn, ExtraAttrsColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns       = postgres.ColumnList{KeyColumn, OrganizationIDColumn, DisplayNameColumn, ExtraAttrsColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return roleTableTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		Key:            KeyColumn,
		OrganizationID: OrganizationIDColumn,
		DisplayName:    DisplayNameColumn,
		ExtraAttrs:     ExtraAttrsColumn,
		CreatedAt:      CreatedAtColumn,
		UpdatedAt:      UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
