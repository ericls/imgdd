package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ericls/imgdd/httpserver"

	"github.com/stretchr/testify/require"
)

func TestAppHandlerKeepsAlbumLinksPrivate(t *testing.T) {
	handler := httpserver.MakeAppHandler(
		httpserver.WithTemplateFS(os.DirFS("templates")),
		httpserver.WithGoogleAnalyticsID("G-TEST"),
	)
	serve := func(path string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, httptest.NewRequest(http.MethodGet, path, nil))
		require.Equal(t, http.StatusOK, res.Code)
		return res
	}

	album := serve("/albums/Q2xhdWRlQ29kZUFsYnVt")
	require.Equal(t, "noindex, nofollow", album.Header().Get("X-Robots-Tag"))
	require.Contains(t, album.Body.String(), `<meta name="robots" content="noindex, nofollow" />`)
	// The album URL grants access, so it must not reach analytics.
	require.NotContains(t, album.Body.String(), "G-TEST")

	for _, path := range []string{"/", "/images/abc", "/albums"} {
		page := serve(path)
		require.Empty(t, page.Header().Get("X-Robots-Tag"), path)
		require.NotContains(t, page.Body.String(), `name="robots"`, path)
		require.True(t, strings.Contains(page.Body.String(), "G-TEST"), path)
	}
}
