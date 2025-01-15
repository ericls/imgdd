package persister

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/ericls/imgdd/logging"

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
	useCookie  bool
	// headerName is the name of the header to use for session token
	// only used if useCookie is false
	headerName string
}

func NewSessionPersister(redisURI string, contextKey *contextKeyT, cookieName *string, headerName *string) *SessionPersister {
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
	var useCookie bool = true
	var realHeaderName string = ""
	if headerName != nil {
		realHeaderName = *headerName
	}
	if realHeaderName != "" {
		useCookie = false
	}
	return &SessionPersister{
		redisURI,
		realContextKey,
		realCookieName,
		logging.GetLogger("httpserver"),
		useCookie,
		realHeaderName,
	}
}

func (sp *SessionPersister) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cs := newContextSession(nil)
		newContext := context.WithValue(r.Context(), sp.contextKey, cs)
		r_with_cs := r.WithContext(newContext)
		sp.fixReaderHeader(r_with_cs)
		explicitly_wrote_header := false
		wrapped_w := httpsnoop.Wrap(w, httpsnoop.Hooks{
			Write: func(inner_w httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				if !explicitly_wrote_header {
					sp.saveSession(w, r_with_cs)
					sp.fixWriterHeader(w)
				}
				return inner_w
			},
			WriteHeader: func(inner_w httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				explicitly_wrote_header = true
				sp.saveSession(w, r_with_cs)
				sp.fixWriterHeader(w)
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

func (sp *SessionPersister) fixReaderHeader(r *http.Request) {
	if sp.useCookie {
		return
	}
	removeCookie(r, sp.cookieName)
	sessionId := r.Header.Get(sp.headerName)
	cookie := &http.Cookie{
		Name:  sp.cookieName,
		Value: sessionId,
	}
	r.AddCookie(cookie)
}

func removeCookie(r *http.Request, name string) {
	oldCookies := r.Cookies()
	var newCookies []*http.Cookie
	for _, c := range oldCookies {
		if c.Name != name {
			newCookies = append(newCookies, c)
		}
	}
	r.Header.Del("Cookie")

	for i, c := range newCookies {
		pair := c.Name + "=" + c.Value
		if i == 0 {
			r.Header.Set("Cookie", pair)
		} else {
			existing := r.Header.Get("Cookie")
			r.Header.Set("Cookie", existing+"; "+pair)
		}
	}
}

func (sp *SessionPersister) fixWriterHeader(w http.ResponseWriter) {
	headers := w.Header()
	if sp.useCookie {
		return
	}
	if headers.Get("Set-Cookie") == "" {
		return
	}
	sessionId := ""
	cookieLines := headers.Values("Set-Cookie")
	cookies := make([]*http.Cookie, 0)
	for _, cookieLine := range cookieLines {
		cookie, err := http.ParseSetCookie(cookieLine)
		if err != nil {
			continue
		}
		cookies = append(cookies, cookie)
	}
	filteredCookies := make([]*http.Cookie, 0)
	for _, c := range cookies {
		if c.Name != sp.cookieName {
			filteredCookies = append(filteredCookies, c)
		} else {
			sessionId = c.Value
		}
	}
	cookieStrs := make([]string, 0)
	for _, c := range filteredCookies {
		cookieStrs = append(cookieStrs, c.String())
	}
	headers.Set("Set-Cookie", strings.Join(cookieStrs, "; "))

	if sessionId != "" {
		headers.Set(sp.headerName, sessionId)
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
	return nil
}

func (sp *SessionPersister) Clear(w http.ResponseWriter, r *http.Request) error {
	v, exist := sp.getContextSession(r)
	if !exist {
		return nil
	}
	v.session.Options.MaxAge = -1
	v.changed = true
	return nil
}

func (sp *SessionPersister) Rotate(w http.ResponseWriter, r *http.Request) error {
	v, exist := sp.getContextSession(r)
	if !exist {
		return nil
	}
	v.session.Options.MaxAge = -1
	v.changed = true
	store, err := sessionStoreCacheInstance.getOrCreate(sp.redisURI)
	if err != nil {
		sp.logger.Error().Err(err).Msg("Could not get or create session store")
	}
	removeCookie(r, sp.cookieName)
	session, err := store.New(r, sp.cookieName)
	if err != nil {
		sp.logger.Error().Err(err).Msg("Could not create new session")
		return err
	}
	session.ID, err = generateSessionKey(128)
	if err != nil {
		sp.logger.Error().Err(err).Msg("Could not generate session key")
	}
	v.session = session
	v.changed = true
	return nil
}

func generateSessionKey(keyLen int) (string, error) {
	bytes := make([]byte, keyLen)

	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
