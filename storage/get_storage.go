package storage

import (
	"errors"

	"github.com/ericls/imgdd/domainmodels"
)

func GetStorage(storageDef *domainmodels.StorageDefinition) (Storage, error) {
	if storageDef.StorageType == domainmodels.S3StorageType {
		backend := GetBackend("s3")
		return backend.FromJSONConfig([]byte(storageDef.Config))
	}
	if storageDef.StorageType == domainmodels.FSStorageType {
		backend := GetBackend("fs")
		return backend.FromJSONConfig([]byte(storageDef.Config))
	}
	return nil, errors.New("invalid storage type")
}
