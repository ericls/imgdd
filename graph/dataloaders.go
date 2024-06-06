package graph

import (
	"context"
	dm "imgdd/domainmodels"
	"imgdd/graph/model"
	"imgdd/identity"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"
)

type ldctxKey string

const loadersKey ldctxKey = "loaders_key"

type Loaders struct {
	UserLoader *dataloadgen.Loader[string, *model.User]
}

func makeUserLoader(identityRepo identity.IdentityRepo) func(c context.Context, keys []string) ([]*model.User, []error) {
	return func(c context.Context, keys []string) ([]*model.User, []error) {
		users := identityRepo.GetUsersByIds(context.Background(), keys)
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

func NewLoaders(identityRepo identity.IdentityRepo) *Loaders {
	return &Loaders{
		UserLoader: dataloadgen.NewLoader(makeUserLoader(identityRepo), dataloadgen.WithWait(time.Millisecond)),
	}
}

func NewLoadersMiddleware(identityRepo identity.IdentityRepo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		l := NewLoaders(identityRepo)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoadersFor(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
