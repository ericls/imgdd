package model

import "github.com/ericls/imgdd/domainmodels"

type Role struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func FromIdentityRole(identityRole *domainmodels.Role) *Role {
	if identityRole == nil {
		return nil
	}
	role := Role{}
	role.Key = identityRole.Key
	role.Name = identityRole.Name
	return &role
}

func FromIdentityRoles(identityRoles []*domainmodels.Role) []*Role {
	if identityRoles == nil {
		return nil
	}
	roles := make([]*Role, len(identityRoles))
	for i, role := range identityRoles {
		roles[i] = FromIdentityRole(role)
	}
	return roles
}
