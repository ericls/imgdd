package model

import "imgdd/domainmodels"

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func FromIdentityOrganization(identityOrg *domainmodels.Organization) *Organization {
	if identityOrg == nil {
		return nil
	}
	org := Organization{}
	org.ID = identityOrg.Id
	org.Name = identityOrg.DisplayName
	org.Slug = identityOrg.Slug
	return &org
}

type User struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	Organization *Organization `json:"organization"`
}

func FromIdentityUser(identityUser *domainmodels.User) *User {
	if identityUser == nil {
		return nil
	}
	u := User{}
	u.ID = identityUser.Id
	u.Name = identityUser.Email
	u.Email = identityUser.Email
	return &u
}
