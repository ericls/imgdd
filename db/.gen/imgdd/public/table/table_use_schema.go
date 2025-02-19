//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	ImageTable = ImageTable.FromSchema(schema)
	OrganizationTable = OrganizationTable.FromSchema(schema)
	OrganizationUserRoleTable = OrganizationUserRoleTable.FromSchema(schema)
	OrganizationUserTable = OrganizationUserTable.FromSchema(schema)
	RoleTable = RoleTable.FromSchema(schema)
	StorageDefinitionTable = StorageDefinitionTable.FromSchema(schema)
	StoredImageTable = StoredImageTable.FromSchema(schema)
	UserTable = UserTable.FromSchema(schema)
}
