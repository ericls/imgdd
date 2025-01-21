package model

import "github.com/ericls/imgdd/domainmodels"

type StoredImage struct {
	ID                string             `json:"id"`
	FileIdentifier    string             `json:"fileIdentifier"`
	StorageDefinition *StorageDefinition `json:"storageDefinition"`
}

func FromStorageStoredImage(si *domainmodels.StoredImage, sd *StorageDefinition) *StoredImage {
	if si == nil {
		return nil
	}
	return &StoredImage{
		ID:                si.Id,
		FileIdentifier:    si.FileIdentifier,
		StorageDefinition: sd,
	}
}
