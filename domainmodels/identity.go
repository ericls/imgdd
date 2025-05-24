package domainmodels

type OrganizationUser struct {
	Id           string
	Organization *Organization
	User         *User
	Roles        []*Role
}

func (ou *OrganizationUser) IsSiteOwner() bool {
	for _, role := range ou.Roles {
		if role.Key == "site_owner" {
			return true
		}
	}
	return false
}

func (ou *OrganizationUser) HasAdminAccess() bool {
	for _, role := range ou.Roles {
		if role.Key == "admin" || role.Key == "owner" || role.Key == "site_owner" {
			return true
		}
	}
	return false
}

func (ou *OrganizationUser) CanManage(other *OrganizationUser) bool {
	if ou.IsSiteOwner() {
		return true
	}
	if other == nil {
		return false
	}
	if ou.Id == other.Id {
		return true
	}
	if ou.HasAdminAccess() && other.Organization.Id == ou.Organization.Id {
		return true
	}
	return false
}

type User struct {
	Id    string
	Email string
}

type Role struct {
	Id   string
	Key  string
	Name string
}

type Organization struct {
	Id          string
	DisplayName string
	Slug        string
}

type UserWithOrganizationUsers struct {
	User              *User
	OrganizationUsers []*OrganizationUser
}
