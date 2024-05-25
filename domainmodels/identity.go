package domainmodels

type OrganizationUser struct {
	Id           string
	Organization *Organization
	User         *User
	Roles        []*Role
}

type User struct {
	Id             string
	OrganizationId string
	Email          string
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
	User
	OrganizationUsers []*OrganizationUser
}
