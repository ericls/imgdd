package httpserver

import (
	"imgdd/test_support"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestLoginLogout(t *testing.T) {
	identityRepo := &test_support.TestIdentityRepo{}
	identityRepo.Reset()
	contextUserManager := NewContextUserManager("foo", identityRepo)
	testIdentityManager := IdentityManager{
		IdentityRepo:       identityRepo,
		ContextUserManager: contextUserManager,
	}
	orgUser1, err := testIdentityManager.IdentityRepo.CreateUserWithOrganization("test@home.arpa", "test", "test")
	if err != nil {
		t.Fatal(err)
	}
	orgUser2, err := testIdentityManager.IdentityRepo.CreateUserWithOrganization(
		"test2@home.arpa", "test2", "test2",
	)
	if err != nil {
		t.Fatal(err)
	}
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(SessionMiddleware)
	r.Use(RWContextMiddleware)
	r.Use(testIdentityManager.Middleware)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Authenticate(w, r, orgUser1.User.Id, orgUser1.Id)
		w.Write([]byte(""))
	})
	r.HandleFunc("/loginContext", func(w http.ResponseWriter, r *http.Request) {
		testIdentityManager.AuthenticateContext(r.Context(), orgUser2.User.Id, orgUser2.Id)
		w.Write([]byte(""))
	})
	r.HandleFunc("/name", func(w http.ResponseWriter, r *http.Request) {
		authenticationInfo := testIdentityManager.ContextUserManager.GetAuthenticationInfo(r.Context())
		user := authenticationInfo.AuthenticatedUser.User
		if user == nil {
			w.Write([]byte("NO USER"))
		} else {
			w.Write([]byte("test-" + user.Id))
		}
	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		Logout(w, r)
		w.Write([]byte(""))
	})
	t.Run("can login and logout", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()
		handler := r
		handler.ServeHTTP(w, req)
		res := w.Result()
		req2 := httptest.NewRequest(http.MethodGet, "/name", nil)
		cookie := res.Header.Get("Set-Cookie")
		req2.Header.Set("Cookie", cookie)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req2)
		res = w.Result()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "test-"+orgUser1.User.Id {
			t.Fatalf("expected test-%s got: %s", orgUser1.Id, string(data))
		}
		req3 := httptest.NewRequest(http.MethodGet, "/logout", nil)
		req3.Header.Set("Cookie", cookie)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req3)
		req4 := httptest.NewRequest(http.MethodGet, "/name", nil)
		req4.Header.Set("Cookie", cookie)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req4)
		res = w.Result()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "NO USER" {
			t.Fatal("expected NO USER string got", string(data), "<-")
		}
		req5 := httptest.NewRequest(http.MethodGet, "/loginContext", nil)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req5)
		res = w.Result()
		req6 := httptest.NewRequest(http.MethodGet, "/name", nil)
		cookie = res.Header.Get("Set-Cookie")
		req6.Header.Set("Cookie", cookie)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req6)
		res = w.Result()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "test-"+orgUser2.User.Id {
			t.Fatalf("expected test-%s got: %s ", orgUser2.Id, string(data))
		}
	})
}
