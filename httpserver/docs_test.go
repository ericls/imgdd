package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/httpserver"

	"github.com/stretchr/testify/require"
)

func TestAPIDocsMarkdown(t *testing.T) {
	serve := func(method string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		httpserver.MakeAPIDocsMarkdownHandler().ServeHTTP(res, httptest.NewRequest(method, "/docs/api.md", nil))
		return res
	}

	res := serve(http.MethodGet)
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "text/markdown; charset=utf-8", res.Header().Get("Content-Type"))
	body := res.Body.String()
	require.True(t, strings.HasPrefix(body, "# imgdd HTTP API"))
	// Placeholders are filled in for this instance.
	require.NotContains(t, body, "{{")
	require.Contains(t, body, "Base URL: `http://example.com`")
	require.Contains(t, body, "POST http://example.com/api/v1/createAlbum")
	require.Contains(t, body, "up to 5 per request")
	require.Contains(t, body, "at most "+strconv.Itoa(domainmodels.ImageMaxCaptionLength)+" characters")
	require.Contains(t, body, "at most "+strconv.Itoa(domainmodels.AlbumMaxTitleLength)+"\ncharacters")
	require.Contains(t, body, "from 60 to 31536000")

	head := serve(http.MethodHead)
	require.Equal(t, http.StatusOK, head.Code)
	require.Empty(t, head.Body.String())
	require.Equal(t, strconv.Itoa(len(body)), head.Header().Get("Content-Length"))

	post := serve(http.MethodPost)
	require.Equal(t, http.StatusMethodNotAllowed, post.Code)
}

func TestLlmsTxt(t *testing.T) {
	res := httptest.NewRecorder()
	httpserver.MakeLlmsTxtHandler().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/llms.txt", nil))
	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "text/plain; charset=utf-8", res.Header().Get("Content-Type"))
	require.True(t, strings.HasPrefix(res.Body.String(), "# imgdd\n\n> "))
	require.Contains(t, res.Body.String(), "(http://example.com/docs/api.md)")
}

// Every error code the API can return must be documented, both in the
// markdown reference and on the web client's docs page.
func TestAPIDocsListEveryErrorCode(t *testing.T) {
	codePattern := regexp.MustCompile(`(?:writeAPIError\(w, [^,]+, |newUploadError\([^,]+, |Code: )"([a-z_]+)"`)
	codes := map[string]bool{}
	for _, file := range []string{"api_handlers.go", "upload.go"} {
		src, err := os.ReadFile(file)
		require.NoError(t, err)
		for _, m := range codePattern.FindAllSubmatch(src, -1) {
			codes[string(m[1])] = true
		}
	}
	sorted := make([]string, 0, len(codes))
	for code := range codes {
		sorted = append(sorted, code)
	}
	sort.Strings(sorted)
	// Guard against the pattern silently matching nothing.
	require.Contains(t, sorted, "album_not_found")
	require.Contains(t, sorted, "rate_limited")
	require.GreaterOrEqual(t, len(sorted), 15)

	markdown, err := os.ReadFile("docs/api.md")
	require.NoError(t, err)
	webDocs, err := os.ReadFile("../web_client/src/app/pages/apiDocs.tsx")
	require.NoError(t, err)
	for _, code := range sorted {
		require.True(t, strings.Contains(string(markdown), "| `"+code+"` |"), "docs/api.md is missing error code %s", code)
		require.True(t, strings.Contains(string(webDocs), `"`+code+`"`), "apiDocs.tsx is missing error code %s", code)
	}
}
