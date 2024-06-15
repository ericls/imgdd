//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package db

import (
	. "imgdd/db/.gen/imgdd/public/table"
)

func PopulateBuiltInRoles(dbConfig DBConfigDef) {
	conn := GetConnection(&dbConfig)
	stmt := RoleTable.INSERT(
		RoleTable.Key,
		RoleTable.DisplayName,
	).
		VALUES("site_owner", "Site Owner").
		VALUES("owner", "Owner").
		VALUES("admin", "Admin").
		VALUES("member", "Member").
		VALUES("guest", "Guest").
		ON_CONFLICT(RoleTable.Key, RoleTable.OrganizationID).
		DO_NOTHING()
	res, err := stmt.Exec(conn)
	if err != nil {
		panic(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	println("Inserted", affected, "built-in roles")
}
