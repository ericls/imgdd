package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSession(t *testing.T) {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(SessionMiddleware)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := GetSessionValue(r, "test")
		SetSessionValue(w, r, "test", "foo")
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
