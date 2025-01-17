package graph

import (
	"context"
	"net/http"
	"time"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/graph/model"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/storage"

	"github.com/vikstrous/dataloadgen"
)

type ldctxKey string

const loadersKey ldctxKey = "loaders_key"

type Loaders struct {
	UserLoader                   *dataloadgen.Loader[string, *model.User]
	OrganizationUserLoader       *dataloadgen.Loader[string, *model.OrganizationUser]
	StoredImagesLoader           *dataloadgen.Loader[string, *model.StoredImage]
	StoredImagesByImageIdsLoader *dataloadgen.Loader[string, []*model.StoredImage]
	StorageDefinitionsLoader     *dataloadgen.Loader[string, *model.StorageDefinition]
}

func makeUserLoader(identityRepo identity.IdentityRepo) func(c context.Context, keys []string) ([]*model.User, []error) {
	return func(c context.Context, keys []string) ([]*model.User, []error) {
		users := identityRepo.GetUsersByIds(keys)
		idToUser := make(map[string]*dm.User)
		for _, u := range users {
			if u == nil {
				continue
			}
			idToUser[u.Id] = u
		}
		ret := make([]*model.User, len(keys))
		for i, id := range keys {
			if u, ok := idToUser[id]; ok {
				ret[i] = model.FromIdentityUser(u)
			} else {
				ret[i] = nil
			}
		}
		return ret, nil
	}
}

func makeOrganizationUserLoader(identityRepo identity.IdentityRepo) func(c context.Context, keys []string) ([]*model.OrganizationUser, []error) {
	return func(c context.Context, keys []string) ([]*model.OrganizationUser, []error) {
		orgUsers := identityRepo.GetOrganizationUsersByIds(keys)
		idToUser := make(map[string]*dm.OrganizationUser)
		for _, ou := range orgUsers {
			if ou == nil {
				continue
			}
			idToUser[ou.Id] = ou
		}
		ret := make([]*model.OrganizationUser, len(keys))
		for i, id := range keys {
			if ou, ok := idToUser[id]; ok {
				ret[i] = model.FromIdentityOrganizationUser(ou)
			} else {
				ret[i] = nil
			}
		}
		return ret, nil
	}
}

func makeStoredImagesLoader(storedImageRepo storage.StoredImageRepo, storageDefinitionLoader *dataloadgen.Loader[string, *model.StorageDefinition]) func(c context.Context, keys []string) ([]*model.StoredImage, []error) {
	return func(c context.Context, keys []string) ([]*model.StoredImage, []error) {
		storedImages, err := storedImageRepo.GetStoredImagesByIds(keys)
		if err != nil {
			return nil, []error{err}
		}
		idToStoredImage := make(map[string]*dm.StoredImage)
		storageDefIds := make([]string, 0)
		for _, si := range storedImages {
			if si == nil {
				continue
			}
			storageDefIds = append(storageDefIds, si.StorageDefinitionId)
			idToStoredImage[si.Id] = si
		}
		storageDefs, err := storageDefinitionLoader.LoadAll(c, storageDefIds)
		if err != nil {
			return nil, []error{err}
		}
		idToStorageDef := make(map[string]*model.StorageDefinition)
		for _, sd := range storageDefs {
			if sd == nil {
				continue
			}
			idToStorageDef[sd.Id] = sd
		}
		ret := make([]*model.StoredImage, len(keys))
		for i, id := range keys {
			if si, ok := idToStoredImage[id]; ok {
				ret[i] = model.FromStorageStoredImage(si, idToStorageDef[si.StorageDefinitionId])
			} else {
				ret[i] = nil
			}
		}
		return ret, nil
	}
}

func makeStoredImagesByImageIdsLoader(storedImageRepo storage.StoredImageRepo, storedImageLoader *dataloadgen.Loader[string, *model.StoredImage]) func(c context.Context, imageIds []string) ([][]*model.StoredImage, []error) {
	return func(c context.Context, imageIds []string) ([][]*model.StoredImage, []error) {
		storedImageIdsByImageId, err := storedImageRepo.GetStoredImageIdsByImageIds(imageIds)
		if err != nil {
			return nil, []error{err}
		}
		ret := make([][]*model.StoredImage, len(imageIds))
		thunks := make([]func() ([]*model.StoredImage, error), len(imageIds))
		for i, id := range imageIds {
			storedImageIds := storedImageIdsByImageId[id]
			thunk := storedImageLoader.LoadAllThunk(c, storedImageIds)
			thunks[i] = thunk
		}
		for i, thunk := range thunks {
			ret[i], err = thunk()
			if err != nil {
				return nil, []error{err}
			}
		}
		return ret, nil
	}
}

func makeStorageDefinitionsLoader(storageDefRepo storage.StorageDefRepo) func(c context.Context, keys []string) ([]*model.StorageDefinition, []error) {
	return func(c context.Context, keys []string) ([]*model.StorageDefinition, []error) {
		storageDefs, err := storageDefRepo.GetStorageDefinitionsByIds(keys)
		if err != nil {
			return nil, []error{err}
		}
		idToStorageDef := make(map[string]*dm.StorageDefinition)
		for _, sd := range storageDefs {
			if sd == nil {
				continue
			}
			idToStorageDef[sd.Id] = sd
		}
		ret := make([]*model.StorageDefinition, len(keys))
		for i, id := range keys {
			if sd, ok := idToStorageDef[id]; ok {
				graphSd, err := model.FromStorageDefinition(sd)
				if err != nil {
					return nil, []error{err}
				}
				ret[i] = graphSd
			} else {
				ret[i] = nil
			}
		}
		return ret, nil
	}
}

func NewLoaders(identityRepo identity.IdentityRepo, storageDefRepo storage.StorageDefRepo, storedImageRepo storage.StoredImageRepo) *Loaders {
	userLoader := dataloadgen.NewLoader(makeUserLoader(identityRepo), dataloadgen.WithWait(time.Millisecond))
	organizationUserLoader := dataloadgen.NewLoader(makeOrganizationUserLoader(identityRepo), dataloadgen.WithWait(time.Millisecond))
	storageDefinitionsLoader := dataloadgen.NewLoader(makeStorageDefinitionsLoader(storageDefRepo), dataloadgen.WithWait(time.Millisecond))
	storedImagesLoader := dataloadgen.NewLoader(makeStoredImagesLoader(storedImageRepo, storageDefinitionsLoader), dataloadgen.WithWait(time.Millisecond))
	storedImagesByImageIdsLoader := dataloadgen.NewLoader(makeStoredImagesByImageIdsLoader(storedImageRepo, storedImagesLoader), dataloadgen.WithWait(time.Millisecond))
	return &Loaders{
		UserLoader:                   userLoader,
		OrganizationUserLoader:       organizationUserLoader,
		StoredImagesLoader:           storedImagesLoader,
		StoredImagesByImageIdsLoader: storedImagesByImageIdsLoader,
	}
}

func NewLoadersMiddleware(
	identityRepo identity.IdentityRepo,
	storageDefRepo storage.StorageDefRepo,
	storedImageRepo storage.StoredImageRepo,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		l := NewLoaders(identityRepo, storageDefRepo, storedImageRepo)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoadersFor(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
