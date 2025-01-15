package storage

import (
	"errors"

	"github.com/ericls/imgdd/domainmodels"
)

func GetStorage(storageDef *domainmodels.StorageDefinition) (Storage, error) {
	if storageDef.StorageType == "s3" {
		backend := GetBackend("s3")
		return backend.FromJSONConfig([]byte(storageDef.Config))
	}
	return nil, errors.New("invalid storage type")
}
