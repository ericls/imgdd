//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package identity

import (
	"database/sql"
	"errors"
	"imgdd/db"
	"imgdd/db/.gen/imgdd/public/model"
	. "imgdd/db/.gen/imgdd/public/table"
	"imgdd/utils"

	dm "imgdd/domainmodels"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type DBIdentityRepo struct {
	db.RepoConn
}

func (repo *DBIdentityRepo) WithTransaction(tx *sql.Tx) db.DBRepo {
	return &DBIdentityRepo{
		RepoConn: repo.RepoConn.WithTransaction(tx),
	}
}

func NewDBIdentityRepo(conn *sql.DB) *DBIdentityRepo {
	return &DBIdentityRepo{
		RepoConn: db.NewRepoConn(conn),
	}
}

var userSelect = SELECT(
	UserTable.AllColumns,
).FROM(UserTable)

type userSelectResult struct {
	model.UserTable
}

func convertUser(jetU *userSelectResult) *dm.User {
	return &dm.User{
		Id:    jetU.ID.String(),
		Email: jetU.Email,
	}
}

func convertRoles(jetRoles []model.RoleTable) []*dm.Role {
	ret := make([]*dm.Role, len(jetRoles))
	for i, r := range jetRoles {
		ret[i] = &dm.Role{
			Id:   r.ID.String(),
			Key:  r.Key,
			Name: r.DisplayName,
		}
	}
	return ret
}

func convertOrganization(jetO *model.OrganizationTable) *dm.Organization {
	return &dm.Organization{
		Id:          jetO.ID.String(),
		DisplayName: jetO.DisplayName,
		Slug:        jetO.Slug,
	}
}

var organizationUserSelectFrom = OrganizationUserTable.LEFT_JOIN(
	UserTable, UserTable.ID.EQ(OrganizationUserTable.UserID),
).
	LEFT_JOIN(
		OrganizationTable, OrganizationTable.ID.EQ(OrganizationUserTable.OrganizationID),
	).
	LEFT_JOIN(
		OrganizationUserRoleTable, OrganizationUserRoleTable.OrganizationUserID.EQ(OrganizationUserTable.ID),
	).
	LEFT_JOIN(
		RoleTable, RoleTable.ID.EQ(OrganizationUserRoleTable.RoleID),
	)

var organizationUserSelect = SELECT(
	OrganizationUserTable.AllColumns,
	UserTable.AllColumns,
	OrganizationTable.AllColumns,
	OrganizationUserRoleTable.AllColumns,
	RoleTable.AllColumns,
).FROM(organizationUserSelectFrom)

type organizationUserSelectResult struct {
	model.OrganizationUserTable
	User         model.UserTable
	Organization model.OrganizationTable
	Roles        []model.RoleTable
}

func convertOrganizationUser(jetOU *organizationUserSelectResult) *dm.OrganizationUser {
	return &dm.OrganizationUser{
		Id:           jetOU.ID.String(),
		Organization: convertOrganization(&jetOU.Organization),
		User:         convertUser(&userSelectResult{UserTable: jetOU.User}),
		Roles:        convertRoles(jetOU.Roles),
	}
}

func (repo *DBIdentityRepo) GetUserById(id string) *dm.User {
	dest := userSelectResult{}
	stmt := userSelect.WHERE(UserTable.ID.EQ(UUID(uuid.MustParse(id))))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	return convertUser(&dest)
}

func (repo *DBIdentityRepo) GetUserByEmail(email string) *dm.User {
	dest := userSelectResult{}
	stmt := userSelect.WHERE(UserTable.Email.EQ(String(email)))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	return convertUser(&dest)
}

func (repo *DBIdentityRepo) GetUserPassword(id string) string {
	dest := model.UserTable{}
	stmt := SELECT(
		UserTable.Password,
	).LIMIT(1).FROM(UserTable).WHERE(UserTable.ID.EQ(UUID(uuid.MustParse(id))))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return ""
	}
	return dest.Password
}

func (repo *DBIdentityRepo) GetUsersByIds(ids []string) []*dm.User {
	dest := []userSelectResult{}
	uuids := make([]Expression, len(ids))
	for i, id := range ids {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := userSelect.WHERE(UserTable.ID.IN(uuids...))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	ret := make([]*dm.User, len(dest))
	for i, d := range dest {
		ret[i] = convertUser(&d)
	}
	return ret
}

func (repo *DBIdentityRepo) GetOrganizationById(id string) *dm.Organization {
	dest := model.OrganizationTable{}
	stmt := SELECT(
		OrganizationTable.AllColumns,
	).LIMIT(1).FROM(OrganizationTable).WHERE(OrganizationTable.ID.EQ(UUID(uuid.MustParse(id))))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	return convertOrganization(&dest)
}

func (repo *DBIdentityRepo) GetOrganizationsByIds(ids []string) []*dm.Organization {
	dest := []model.OrganizationTable{}
	uuids := make([]Expression, len(ids))
	for i, id := range ids {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := SELECT(
		OrganizationTable.AllColumns,
	).FROM(OrganizationTable).WHERE(OrganizationTable.ID.IN(uuids...))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	ret := make([]*dm.Organization, len(dest))
	for i, d := range dest {
		ret[i] = convertOrganization(&d)
	}
	return ret
}

func (repo *DBIdentityRepo) GetOrganizationUserById(id string) *dm.OrganizationUser {
	dest := organizationUserSelectResult{}
	stmt := organizationUserSelect.WHERE(OrganizationUserTable.ID.EQ(UUID(uuid.MustParse(id))))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	return convertOrganizationUser(&dest)
}

func (repo *DBIdentityRepo) GetOrganizationUsersByIds(ids []string) []*dm.OrganizationUser {
	dest := []organizationUserSelectResult{}
	uuids := make([]Expression, len(ids))
	for i, id := range ids {
		uuids[i] = UUID(uuid.MustParse(id))
	}
	stmt := organizationUserSelect.WHERE(OrganizationUserTable.ID.IN(uuids...))
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	ret := make([]*dm.OrganizationUser, len(dest))
	for i, d := range dest {
		ret[i] = convertOrganizationUser(&d)
	}
	return ret
}

func (repo *DBIdentityRepo) createOrganizationUser(organizationId string, userId string) (*dm.OrganizationUser, error) {
	stmt := OrganizationUserTable.INSERT(
		OrganizationUserTable.OrganizationID,
		OrganizationUserTable.UserID,
	).VALUES(organizationId, userId).RETURNING(OrganizationUserTable.ID)
	dest := model.OrganizationUserTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	queryStmt := organizationUserSelect.WHERE(OrganizationUserTable.ID.EQ(UUID(dest.ID)))
	queryDest := organizationUserSelectResult{}
	err = queryStmt.Query(repo.DB, &queryDest)
	if err != nil {
		return nil, err
	}
	return convertOrganizationUser(&queryDest), nil
}

func (repo *DBIdentityRepo) CreateUser(email string, password string) (*dm.User, error) {
	hashsedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	return db.RunInTransaction(repo, func(txRepo *DBIdentityRepo) (*dm.User, error) {
		if err != nil {
			return nil, err
		}
		insertDest := model.UserTable{}
		stmt := UserTable.INSERT(
			UserTable.Email, UserTable.Password,
		).VALUES(email, hashsedPassword).RETURNING(UserTable.AllColumns)
		err = stmt.Query(txRepo.GetDB(), &insertDest)
		if err != nil {
			return nil, err
		}
		dest := userSelectResult{}
		queryStmt := userSelect.WHERE(UserTable.ID.EQ(UUID(insertDest.ID)))
		err = queryStmt.Query(txRepo.GetDB(), &dest)
		if err != nil {
			return nil, err
		}
		return convertUser(&dest), nil
	})
}

func (repo *DBIdentityRepo) CreateOrganization(name, slug string) (*dm.Organization, error) {
	stmt := OrganizationTable.INSERT(
		OrganizationTable.DisplayName,
		OrganizationTable.Slug,
	).VALUES(name, slug).RETURNING(OrganizationTable.AllColumns)
	dest := model.OrganizationTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, err
	}
	return convertOrganization(&dest), nil
}

func (repo *DBIdentityRepo) getRoleByKey(organizationId string, roleKey string) *dm.Role {
	stmt := RoleTable.SELECT(
		RoleTable.AllColumns,
	).WHERE(
		RoleTable.Key.EQ(String(roleKey)).AND(
			RoleTable.OrganizationID.EQ(UUID(uuid.MustParse(organizationId))).OR(RoleTable.OrganizationID.IS_NULL()),
		),
	)
	dest := model.RoleTable{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil
	}
	return &dm.Role{
		Id:   dest.ID.String(),
		Key:  dest.Key,
		Name: dest.DisplayName,
	}
}

func (repo *DBIdentityRepo) AddRoleToOrganizationUser(organizationUserId, roleKey string) error {
	orgUser := repo.GetOrganizationUserById(organizationUserId)
	role := repo.getRoleByKey(orgUser.Organization.Id, roleKey)
	if role == nil {
		return errors.New("role not found")
	}
	stmt := OrganizationUserRoleTable.INSERT(
		OrganizationUserRoleTable.OrganizationUserID,
		OrganizationUserRoleTable.RoleID,
	).VALUES(organizationUserId, role.Id)
	_, err := stmt.Exec(repo.DB)
	return err
}

func (repo *DBIdentityRepo) CreateUserWithOrganization(
	email string, organizationName string, password string,
) (*dm.OrganizationUser, error) {

	return db.RunInTransaction(repo, func(txRepo *DBIdentityRepo) (*dm.OrganizationUser, error) {
		organization, err := txRepo.CreateOrganization(organizationName, utils.Slugify(organizationName))
		if err != nil {
			return nil, err
		}
		user, err := txRepo.CreateUser(email, password)
		if err != nil {
			return nil, err
		}
		// Add user to organization
		orgUser, err := txRepo.createOrganizationUser(organization.Id, user.Id)
		if err != nil {
			return nil, err
		}
		err = txRepo.AddRoleToOrganizationUser(orgUser.Id, "owner")
		if err != nil {
			return nil, err
		}
		orgUser = txRepo.GetOrganizationUserById(orgUser.Id)
		return orgUser, nil
	})

}

func (repo *DBIdentityRepo) GetOrganizationForUser(userId string, maybeOrganizationId string) (*dm.Organization, *dm.OrganizationUser) {
	// if maybeOrganizationId is empty, get the only organization the user is a member of
	if maybeOrganizationId == "" {
		stmt := OrganizationTable.SELECT(
			OrganizationTable.ID.AS("org_id"),
			OrganizationUserTable.ID.AS("organization_user_id"),
		).FROM(
			OrganizationTable.LEFT_JOIN(
				OrganizationUserTable,
				OrganizationUserTable.OrganizationID.EQ(OrganizationTable.ID),
			),
		).WHERE(
			OrganizationUserTable.UserID.EQ(UUID(uuid.MustParse(userId))),
		)
		dest := []struct {
			OrgID              string
			OrganizationUserId string
		}{}
		err := stmt.Query(repo.DB, &dest)
		if err != nil {
			println("Error", err)
			return nil, nil
		}
		if len(dest) != 1 {
			return nil, nil
		}
		return repo.GetOrganizationById(dest[0].OrgID), repo.GetOrganizationUserById(dest[0].OrganizationUserId)
	}
	// if maybeOrganizationId is not empty, return the organization with that ID, if the user is a member of it
	stmt := OrganizationTable.SELECT(
		OrganizationTable.AllColumns,
		OrganizationUserTable.AllColumns,
	).FROM(
		OrganizationTable.LEFT_JOIN(
			OrganizationUserTable,
			OrganizationUserTable.OrganizationID.EQ(OrganizationTable.ID),
		),
	).WHERE(
		OrganizationUserTable.UserID.EQ(UUID(uuid.MustParse(userId))).AND(
			OrganizationTable.ID.EQ(UUID(uuid.MustParse(maybeOrganizationId))),
		),
	)

	dest := struct {
		model.OrganizationTable
		OrganizationUser model.OrganizationUserTable
	}{}
	err := stmt.Query(repo.DB, &dest)
	if err != nil {
		return nil, nil
	}
	return convertOrganization(&dest.OrganizationTable), repo.GetOrganizationUserById(dest.OrganizationUser.ID.String())
}
