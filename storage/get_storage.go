package storage

import (
	"errors"

	"github.com/ericls/imgdd/domainmodels"
)

func GetStorage(storageDef *domainmodels.StorageDefinition) (Storage, error) {
	backend := GetBackend(storageDef.StorageType)
	if backend == nil {
		return nil, errors.New("invalid storage type")
	}
	s, err := backend.FromJSONConfig([]byte(storageDef.Config))
	if err != nil {
		return nil, err
	}
	return s, nil
}
