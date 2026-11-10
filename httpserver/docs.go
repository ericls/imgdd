package httpserver

import (
	_ "embed"
	"net/http"
	"strconv"
	"strings"

	"github.com/ericls/imgdd/domainmodels"
)

// apiDocsMarkdown is the HTTP API reference as markdown, for programs such as
// AI agents that cannot read the web client's /docs/api page. It contains
// {{PLACEHOLDERS}} that are filled in per request; see renderAPIDocsMarkdown.
//
//go:embed docs/api.md
var apiDocsMarkdown string

const apiDocsMarkdownPath = "/docs/api.md"
const llmsTxtPath = "/llms.txt"

// renderAPIDocsMarkdown fills in the API reference's placeholders, so examples
// and limits match this instance and can be copied as they are.
func renderAPIDocsMarkdown(baseURL string) string {
	return strings.NewReplacer(
		"{{BASE_URL}}", baseURL,
		"{{MAX_IMAGES}}", strconv.Itoa(apiMaxImagesPerRequest),
		"{{MIN_EXPIRES_IN}}", strconv.Itoa(apiMinExpiresInSeconds),
		"{{MAX_EXPIRES_IN}}", strconv.Itoa(apiMaxExpiresInSeconds),
		"{{MAX_CAPTION}}", strconv.Itoa(domainmodels.ImageMaxCaptionLength),
		"{{MAX_TITLE}}", strconv.Itoa(domainmodels.AlbumMaxTitleLength),
	).Replace(apiDocsMarkdown)
}

// renderLlmsTxt returns an llms.txt (https://llmstxt.org) pointing language
// models at the plain-text API reference.
func renderLlmsTxt(baseURL string) string {
	return "# imgdd\n\n" +
		"> Self-hosted image hosting with an HTTP API. AI agents can upload screenshots into a live album page that a person watches while the agent works.\n\n" +
		"## Docs\n\n" +
		"- [HTTP API](" + baseURL + apiDocsMarkdownPath + "): authentication, createAlbum, createImages, error codes, and a quick start for agents sharing screenshots\n"
}

func serveText(contentType string, render func(baseURL string) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Use GET", http.StatusMethodNotAllowed)
			return
		}
		body := render(strings.TrimSuffix(baseURLFromRequest(r).String(), "/"))
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		if r.Method == http.MethodHead {
			return
		}
		w.Write([]byte(body))
	}
}

func makeAPIDocsMarkdownHandler() http.HandlerFunc {
	return serveText("text/markdown; charset=utf-8", renderAPIDocsMarkdown)
}

func makeLlmsTxtHandler() http.HandlerFunc {
	return serveText("text/plain; charset=utf-8", renderLlmsTxt)
}
