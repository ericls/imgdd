package model

import "github.com/ericls/imgdd/domainmodels"

type OrganizationUser struct {
	ID           string        `json:"id"`
	Organization *Organization `json:"organization"`
	User         *User         `json:"user"`
	Roles        []*Role       `json:"roles"`
}

func FromIdentityOrganizationUser(identityOrgUser *domainmodels.OrganizationUser) *OrganizationUser {
	if identityOrgUser == nil {
		return nil
	}
	orgUser := OrganizationUser{}
	orgUser.ID = identityOrgUser.Id
	orgUser.Organization = FromIdentityOrganization(identityOrgUser.Organization)
	orgUser.User = FromIdentityUser(identityOrgUser.User)
	orgUser.Roles = FromIdentityRoles(identityOrgUser.Roles)
	return &orgUser
}
