package persister

import (
	"context"
	"imgdd/logging"
	"net/http"
	"strings"
	"sync"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// SessionStoreCache is a cache for the session store keyed by the redis URI.
type sessionStoreCache struct {
	lock            sync.RWMutex
	redisURIToStore map[string]sessions.Store
}

func newSessionStoreCache() *sessionStoreCache {
	return &sessionStoreCache{
		redisURIToStore: make(map[string]sessions.Store),
	}
}

func (c *sessionStoreCache) get(redisURI string) sessions.Store {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.redisURIToStore[redisURI]
}

func (c *sessionStoreCache) set(redisURI string, store sessions.Store) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.redisURIToStore[redisURI] = store
}

func (c *sessionStoreCache) getOrCreate(redisURI string) (sessions.Store, error) {
	if store := c.get(redisURI); store != nil {
		return store, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr: strings.TrimPrefix(redisURI, "redis://"),
	})

	s, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		return nil, err
	}
	s.KeyPrefix("session_")
	s.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})

	c.set(redisURI, s)
	return s, nil
}

type ContextSession struct {
	session *sessions.Session
	changed bool
}

func newContextSession(s *sessions.Session) *ContextSession {
	cs := ContextSession{s, false}
	return &cs
}

type contextKeyT string

var sessionStoreCacheInstance *sessionStoreCache

func init() {
	sessionStoreCacheInstance = newSessionStoreCache()
}

type SessionPersister struct {
	redisURI   string
	contextKey contextKeyT
	cookieName string
	logger     zerolog.Logger
}

func NewSessionPersister(redisURI string, contextKey *contextKeyT, cookieName *string) *SessionPersister {
	var realContextKey contextKeyT
	if contextKey == nil {
		realContextKey = contextKeyT("context-session")
	} else {
		realContextKey = *contextKey
	}
	var realCookieName string
	if cookieName == nil {
		realCookieName = "session"
	} else {
		realCookieName = *cookieName
	}
	return &SessionPersister{
		redisURI,
		realContextKey,
		realCookieName,
		logging.GetLogger("httpserver"),
	}
}

func (sp *SessionPersister) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cs := newContextSession(nil)
		newContext := context.WithValue(r.Context(), sp.contextKey, cs)
		r_with_cs := r.WithContext(newContext)
		explicitly_wrote_header := false
		wrapped_w := httpsnoop.Wrap(w, httpsnoop.Hooks{
			Write: func(inner_w httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				if !explicitly_wrote_header {
					sp.saveSession(w, r_with_cs)
				}
				return inner_w
			},
			WriteHeader: func(inner_w httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				explicitly_wrote_header = true
				sp.saveSession(w, r_with_cs)
				return inner_w
			},
		})
		next.ServeHTTP(wrapped_w, r_with_cs)
	})
}

func (sp *SessionPersister) getContextSession(r *http.Request) (*ContextSession, bool) {
	// Session middleware is not used
	if r.Context().Value(sp.contextKey) == nil {
		sp.logger.Error().Msg("session context key not found. Is the session middleware used?")
		return nil, false
	}
	context_session, ok := r.Context().Value(sp.contextKey).(*ContextSession)
	// Session middleware is not used and somehow key conflicts
	if !ok {
		return nil, false
	}
	// Session middleware is used but session is not initialized, initialize it
	if context_session.session == nil {
		store, err := sessionStoreCacheInstance.getOrCreate(sp.redisURI)
		if err != nil {
			sp.logger.Error().Err(err).Msg("Could not get or create session store")
			return nil, false
		}
		session, _ := store.Get(r, sp.cookieName)
		context_session.session = session
	}
	return context_session, true
}

func (sp *SessionPersister) saveSession(w http.ResponseWriter, r *http.Request) {
	v, exist := sp.getContextSession(r)
	if !exist {
		return
	}
	if v.changed {
		v.session.Save(r, w)
	}
}

// implement Persister interface
func (sp *SessionPersister) Set(w http.ResponseWriter, r *http.Request, key string, value string) error {
	v, exist := sp.getContextSession(r)
	if !exist {
		return nil
	}
	v.session.Values[key] = value
	v.changed = true
	sp.saveSession(w, r)
	return nil
}

func (sp *SessionPersister) Get(r *http.Request, key string) (string, error) {
	v, exist := sp.getContextSession(r)
	if !exist {
		return "", nil
	}
	vv, e := v.session.Values[key]
	if e {
		return vv.(string), nil
	}
	return "", nil
}

func (sp *SessionPersister) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	v, e := sp.getContextSession(r)
	if e {
		delete(v.session.Values, key)
		v.changed = true
	}
	sp.saveSession(w, r)
	return nil
}

func (sp *SessionPersister) Clear(w http.ResponseWriter, r *http.Request) error {
	v, exist := sp.getContextSession(r)
	if !exist {
		return nil
	}
	v.session.Options.MaxAge = -1
	v.changed = true
	sp.saveSession(w, r)
	return nil
}
