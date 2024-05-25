package httpserver

import (
	"context"
	"imgdd/logging"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var store sessions.Store
var logger zerolog.Logger

const session_context_key = ContextKey("context-session")
const session_cookie_name = "session"

func init() {
	logger = logging.GetLogger("httpserver")
	client := redis.NewClient(&redis.Options{
		Addr: strings.TrimPrefix(Config.RedisURIForSession, "redis://"),
	})

	s, err := redisstore.NewRedisStore(context.Background(), client)
	s.KeyPrefix("session_")
	s.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})
	if err != nil {
		logger.Err(err).Msg("failed to create redis session store")
	}
	store = s
}

type ContextSession struct {
	session *sessions.Session
	changed bool
}

func newContextSession(s *sessions.Session) *ContextSession {
	cs := ContextSession{s, false}
	return &cs
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cs := newContextSession(nil)
		newContext := context.WithValue(r.Context(), session_context_key, cs)
		r = r.WithContext(newContext)
		explicitly_wrote_header := false
		wrapped_w := httpsnoop.Wrap(w, httpsnoop.Hooks{
			Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				if !explicitly_wrote_header {
					saveSession(w, r)
				}
				return w.Write
			},
			WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				explicitly_wrote_header = true
				saveSession(w, r)
				return w.WriteHeader
			},
		})
		next.ServeHTTP(wrapped_w, r)
	})
}

func saveSession(w http.ResponseWriter, r *http.Request) {
	v, exist := GetContextSession(r)
	if !exist {
		return
	}
	if v.changed {
		v.session.Save(r, w)
	}
}

func GetSessionValue(r *http.Request, key string) string {
	v, exist := GetContextSession(r)
	if !exist {
		return ""
	}
	vv, e := v.session.Values[key]
	if e {
		return vv.(string)
	}
	return ""
}

func GetContextSession(r *http.Request) (*ContextSession, bool) {
	// Session middleware is not used
	if r.Context().Value(session_context_key) == nil {
		return nil, false
	}
	context_session, ok := r.Context().Value(session_context_key).(*ContextSession)
	// Session middleware is not used and somehow key conflicts
	if !ok {
		return nil, false
	}
	// Session middleware is used but session is not initialized, initialize it
	if context_session.session == nil {
		session, _ := store.Get(r, session_cookie_name)
		context_session.session = session
	}
	return context_session, true
}

func SetSessionValue(w http.ResponseWriter, r *http.Request, key string, value interface{}) {
	v, ok := GetContextSession(r)
	if ok {
		v.session.Values[key] = value
		v.changed = true
	}
	saveSession(w, r)
}

func DeleteSessionValue(w http.ResponseWriter, r *http.Request, key string) {
	v, e := GetContextSession(r)
	if e {
		delete(v.session.Values, key)
		v.changed = true
	}
	saveSession(w, r)
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	v, e := GetContextSession(r)
	if e {
		v.session.Values = make(map[interface{}]interface{})
		v.changed = true
	}
	v.session.Options.MaxAge = -1
	saveSession(w, r)
}
