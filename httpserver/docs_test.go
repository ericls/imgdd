package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/httpserver"

	"github.com/stretchr/testify/require"
)

var docsTestConf = &httpserver.HttpServerConfigDef{DefaultURLFormat: domainmodels.ImageURLFormat_CANONICAL}

func TestAPIDocsMarkdown(t *testing.T) {
	serve := func(method string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		httpserver.MakeAPIDocsMarkdownHandler(docsTestConf).ServeHTTP(res, httptest.NewRequest(method, "/docs/api.md", nil))
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
	require.Contains(t, body, `"url": "http://example.com/image/8f1c2d7e-0b4a-4c9e-a1f2-3d4e5f6a7b8c.png"`)
	require.Contains(t, body, `"pageUrl": "http://example.com/images/`)
	// Built from the request's Host, so shared caches must not store it.
	require.Equal(t, "private, max-age=300", res.Header().Get("Cache-Control"))

	head := serve(http.MethodHead)
	require.Equal(t, http.StatusOK, head.Code)
	require.Empty(t, head.Body.String())
	require.Equal(t, strconv.Itoa(len(body)), head.Header().Get("Content-Length"))

	post := serve(http.MethodPost)
	require.Equal(t, http.StatusMethodNotAllowed, post.Code)
}

// Links in the docs follow the instance's image domain and web UI host, like
// the URLs the API returns.
func TestAPIDocsMarkdownUsesConfiguredHosts(t *testing.T) {
	conf := &httpserver.HttpServerConfigDef{
		ImageDomain:      "img.example.net",
		WebUIHost:        "app.example.net",
		DefaultURLFormat: domainmodels.ImageURLFormat_CANONICAL,
	}
	req := httptest.NewRequest(http.MethodGet, "/docs/api.md", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	res := httptest.NewRecorder()
	httpserver.MakeAPIDocsMarkdownHandler(conf).ServeHTTP(res, req)
	body := res.Body.String()
	require.Contains(t, body, "POST https://example.com/api/v1/createImages")
	require.Contains(t, body, `"url": "https://img.example.net/8f1c2d7e-0b4a-4c9e-a1f2-3d4e5f6a7b8c.png"`)
	require.Contains(t, body, `"pageUrl": "https://app.example.net/images/`)
	require.Contains(t, body, `"albumUrl": "https://app.example.net/albums/`)
	require.Contains(t, body, "https://app.example.net/profile/api-keys")
}

func TestLlmsTxt(t *testing.T) {
	res := httptest.NewRecorder()
	httpserver.MakeLlmsTxtHandler(docsTestConf).ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/llms.txt", nil))
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
	// Every non-test source file, so a new handler file can't be missed.
	files, err := filepath.Glob("*.go")
	require.NoError(t, err)
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
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
	require.Contains(t, sorted, "image_not_found")
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
