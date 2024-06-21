package storage

import (
	"errors"
	"imgdd/domainmodels"
)

func GetStorage(storageDef *domainmodels.StorageDefinition) (Storage, error) {
	if storageDef.StorageType == "s3" {
		backend := GetBackend("s3")
		return backend.FromJSON([]byte(storageDef.Config))
	}
	return nil, errors.New("invalid storage type")
}
