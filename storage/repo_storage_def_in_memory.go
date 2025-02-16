package storage

import (
	"sort"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/google/uuid"
)

type InMemoryStorageDefRepo struct {
	storageDefs map[string]*dm.StorageDefinition
}

func NewInMemoryStorageDefRepo() *InMemoryStorageDefRepo {
	return &InMemoryStorageDefRepo{
		storageDefs: make(map[string]*dm.StorageDefinition),
	}
}

func (repo *InMemoryStorageDefRepo) Clear() {
	repo.storageDefs = make(map[string]*dm.StorageDefinition)
}

func (repo *InMemoryStorageDefRepo) AddStorageDefinition(storageDef *dm.StorageDefinition) {
	repo.storageDefs[storageDef.Id] = storageDef
}

func (repo *InMemoryStorageDefRepo) GetStorageDefinitionById(id string) (*dm.StorageDefinition, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	storageDef, ok := repo.storageDefs[id]
	if !ok {
		return nil, nil
	}
	return storageDef, nil
}

func (repo *InMemoryStorageDefRepo) GetStorageDefinitionsByIds(ids []string) ([]*dm.StorageDefinition, error) {
	storageDefs := make([]*dm.StorageDefinition, 0)
	for _, id := range ids {
		_, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		storageDef, ok := repo.storageDefs[id]
		if ok {
			storageDefs = append(storageDefs, storageDef)
		}
	}
	return storageDefs, nil
}

func (repo *InMemoryStorageDefRepo) GetStorageDefinitionByIdentifier(id string) (*dm.StorageDefinition, error) {
	// TODO: build a map keyed by identifier
	for _, storageDef := range repo.storageDefs {
		if storageDef.Identifier == id {
			return storageDef, nil
		}
	}
	return nil, nil
}

func (repo *InMemoryStorageDefRepo) ListStorageDefinitions() ([]*dm.StorageDefinition, error) {
	storageDefs := make([]*dm.StorageDefinition, 0)
	for _, storageDef := range repo.storageDefs {
		storageDefs = append(storageDefs, storageDef)
	}
	sort.Slice(storageDefs, func(i, j int) bool {
		return storageDefs[i].Priority < storageDefs[j].Priority
	})
	return storageDefs, nil
}

func (repo *InMemoryStorageDefRepo) CreateStorageDefinition(storage_type string, config string, identifier string, isEnabled bool, priority int64) (*dm.StorageDefinition, error) {
	storageDef := &dm.StorageDefinition{
		Id:          uuid.New().String(),
		Identifier:  identifier,
		StorageType: dm.StorageTypeName(storage_type),
		Config:      config,
		IsEnabled:   isEnabled,
		Priority:    int32(priority),
	}
	repo.storageDefs[storageDef.Id] = storageDef
	return storageDef, nil
}

func (repo *InMemoryStorageDefRepo) UpdateStorageDefinition(identifier string, storage_type *string, config *string, isEnabled *bool, priority *int64) (*dm.StorageDefinition, error) {
	storageDef, err := repo.GetStorageDefinitionByIdentifier(identifier)
	if err != nil {
		return nil, err
	}
	if storage_type != nil {
		storageDef.StorageType = dm.StorageTypeName(*storage_type)
	}
	if config != nil {
		storageDef.Config = *config
	}
	if isEnabled != nil {
		storageDef.IsEnabled = *isEnabled
	}
	if priority != nil {
		storageDef.Priority = int32(*priority)
	}
	return storageDef, nil
}
