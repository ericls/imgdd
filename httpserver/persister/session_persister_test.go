package persister_test

import (
	"imgdd/httpserver/persister"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gorilla/mux"
)

func TestSession(t *testing.T) {
	s := miniredis.RunT(t)
	redis_uri := "redis://" + s.Addr()
	persister := persister.NewSessionPersister(redis_uri, nil, nil)
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(persister.Middleware)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v, err := persister.Get(r, "test")
		if err != nil {
			t.Fatal(err)
		}
		err = persister.Set(w, r, "test", "foo")
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(v))
	})
	t.Run("can read and set session values", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler := r
		handler.ServeHTTP(w, req)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "" {
			t.Fatal("expected empty string")
		}
		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		req2.Header.Set("Cookie", res.Header.Get("Set-Cookie"))
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req2)
		res = w.Result()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "foo" {
			t.Fatal("expected foo got", string(data), "<-")
		}
	})
}
