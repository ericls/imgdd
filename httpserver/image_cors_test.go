package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageCORSMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	cases := []struct {
		name           string
		allowedOrigins []string
		origin         string
		expectedACAO   string
		expectedVary   []string
	}{
		{"allowed origin", []string{"https://imgdd.example"}, "https://imgdd.example", "https://imgdd.example", []string{"Origin"}},
		{"disallowed origin", []string{"https://imgdd.example"}, "https://other.example", "", []string{"Origin"}},
		{"no origin", []string{"https://imgdd.example"}, "", "", []string{"Origin"}},
		{"no allowed origins configured", nil, "https://imgdd.example", "", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := makeImageCORSMiddleware(tc.allowedOrigins)(next)
			req := httptest.NewRequest(http.MethodGet, "/image/abc.png", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if acao := res.Header().Get("Access-Control-Allow-Origin"); acao != tc.expectedACAO {
				t.Fatalf("expected Access-Control-Allow-Origin %q, got %q", tc.expectedACAO, acao)
			}
			vary := res.Header().Values("Vary")
			if len(vary) != len(tc.expectedVary) || (len(vary) > 0 && vary[0] != tc.expectedVary[0]) {
				t.Fatalf("expected Vary %v, got %v", tc.expectedVary, vary)
			}
		})
	}
}
