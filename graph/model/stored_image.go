package model

import "imgdd/domainmodels"

type StoredImage struct {
	ID                string             `json:"id"`
	StorageDefinition *StorageDefinition `json:"storageDefinition"`
}

func FromStorageStoredImage(si *domainmodels.StoredImage) *StoredImage {
	if si == nil {
		return nil
	}
	storageDef, err := FromStorageDefinition(si.StorageDefinition)
	if err != nil {
		return nil
	}
	return &StoredImage{
		ID:                si.Id,
		StorageDefinition: storageDef,
	}
}
