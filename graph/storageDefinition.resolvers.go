package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"
	"errors"

	"github.com/ericls/imgdd/graph/model"
	"github.com/ericls/imgdd/storage"
)

// CreateStorageDefinition is the resolver for the createStorageDefinition field.
func (r *mutationResolver) CreateStorageDefinition(ctx context.Context, input model.CreateStorageDefinitionInput) (*model.StorageDefinition, error) {
	repo := r.StorageRepo
	storageType := input.StorageType
	config := input.ConfigJSON
	backend := storage.GetBackend(string(storageType))
	if backend == nil {
		return nil, errors.New("invalid storage type")
	}
	err := backend.ValidateJSONConfig([]byte(config))
	if err != nil {
		return nil, err
	}
	created, err := repo.CreateStorageDefinition(
		string(storageType),
		config,
		input.Identifier,
		input.IsEnabled,
		int64(input.Priority),
	)
	if err != nil {
		return nil, err
	}
	s, err := model.FromStorageDefinition(created)
	return s, err
}

// UpdateStorageDefinition is the resolver for the updateStorageDefinition field.
func (r *mutationResolver) UpdateStorageDefinition(ctx context.Context, input model.UpdateStorageDefinitionInput) (*model.StorageDefinition, error) {
	repo := r.StorageRepo
	config := input.ConfigJSON
	var priority int64
	var priorityPtr *int64
	if input.Priority != nil {
		priority = int64(*input.Priority)
		priorityPtr = &priority
	}
	_, err := repo.GetStorageDefinitionByIdentifier(input.Identifier)
	if err != nil {
		return nil, err
	}
	updated, err := repo.UpdateStorageDefinition(
		input.Identifier,
		nil,
		config,
		input.IsEnabled,
		priorityPtr,
	)
	if err != nil {
		return nil, err
	}
	s, err := model.FromStorageDefinition(updated)
	return s, err
}

// CheckStorageDefinitionConnectivity is the resolver for the checkStorageDefinitionConnectivity field.
func (r *mutationResolver) CheckStorageDefinitionConnectivity(ctx context.Context, input model.CheckStorageDefinitionConnectivityInput) (*model.StorageDefinitionConnectivityResult, error) {
	repo := r.StorageRepo
	storageDefinition, err := repo.GetStorageDefinitionById(input.ID)
	if err != nil {
		return nil, err
	}
	backend, err := storage.GetStorage(storageDefinition)
	if err != nil {
		return nil, err
	}
	err = backend.CheckConnection()
	if err != nil {
		errMessage := err.Error()
		return &model.StorageDefinitionConnectivityResult{
			Ok:    false,
			Error: &errMessage,
		}, nil
	}
	return &model.StorageDefinitionConnectivityResult{
		Ok: true,
	}, nil
}

// Connectivity is the resolver for the connectivity field.
func (r *storageDefinitionResolver) Connectivity(ctx context.Context, obj *model.StorageDefinition) (bool, error) {
	repo := r.StorageRepo
	storageDefinition, err := repo.GetStorageDefinitionById(obj.Id)
	if err != nil {
		return false, err
	}
	backend, err := storage.GetStorage(storageDefinition)
	if err != nil {
		return false, err
	}
	err = backend.CheckConnection()
	if err != nil {
		return false, err
	}
	return true, nil
}

// StorageDefinitions is the resolver for the storageDefinitions field.
func (r *viewerResolver) StorageDefinitions(ctx context.Context, obj *model.Viewer) ([]*model.StorageDefinition, error) {
	repo := r.StorageRepo
	storageDefinitions, err := repo.ListStorageDefinitions()
	if err != nil {
		return nil, err
	}
	ret := make([]*model.StorageDefinition, len(storageDefinitions))
	for i, storageDefinition := range storageDefinitions {
		s, _ := model.FromStorageDefinition(storageDefinition)
		ret[i] = s
	}
	return ret, nil
}

// GetStorageDefinition is the resolver for the getStorageDefinition field.
func (r *viewerResolver) GetStorageDefinition(ctx context.Context, obj *model.Viewer, id string) (*model.StorageDefinition, error) {
	repo := r.StorageRepo
	storageDefinition, err := repo.GetStorageDefinitionById(id)
	if err != nil {
		return nil, err
	}
	s, err := model.FromStorageDefinition(storageDefinition)
	return s, err
}

// StorageDefinition returns StorageDefinitionResolver implementation.
func (r *Resolver) StorageDefinition() StorageDefinitionResolver {
	return &storageDefinitionResolver{r}
}

type storageDefinitionResolver struct{ *Resolver }
