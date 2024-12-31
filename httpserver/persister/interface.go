package persister

import "net/http"

// persister provices an interface to persist things for
// http.
// Examples of implementations are:
// - SessionPersister
// - CookiePersister
// - JWTTokenPersister
// - etc.

type Persister interface {
	Set(w http.ResponseWriter, r *http.Request, key string, value string) error
	Get(r *http.Request, key string) (string, error)
	Delete(w http.ResponseWriter, r *http.Request, key string) error
	Clear(w http.ResponseWriter, r *http.Request) error
	Rotate(w http.ResponseWriter, r *http.Request) error
}
