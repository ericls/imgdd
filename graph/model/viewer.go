package model

type Viewer struct {
	Id               string            `json:"id,omitempty"`
	OrganizationUser *OrganizationUser `json:"organizationUser,omitempty"`
}
