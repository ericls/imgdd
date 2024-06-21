package identity

import "imgdd/db"

func AddUserToGroup(groupKey, userEmail string, dbConf *db.DBConfigDef) error {
	conn := db.GetConnection(dbConf)
	defer conn.Close()
	repo := NewDBIdentityRepo(conn)
	user := repo.GetUserByEmail(userEmail)
	_, orgUser := repo.GetOrganizationForUser(user.Id, "")
	return repo.AddRoleToOrganizationUser(orgUser.Id, groupKey)
}
