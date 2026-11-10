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

// docsURLs are the instance-specific URLs filled into the docs.
type docsURLs struct {
	// API is the base URL the docs were requested from, used for API calls.
	API string
	// WebUI is the base URL of the web UI, used for pages and album links.
	WebUI string
	// ExampleImage is an image URL in the form this instance returns.
	ExampleImage string
}

const docsExampleImageIdentifier = "8f1c2d7e-0b4a-4c9e-a1f2-3d4e5f6a7b8c"

func docsURLsFor(r *http.Request, conf *HttpServerConfigDef) docsURLs {
	format := conf.DefaultURLFormat
	if format != domainmodels.ImageURLFormat_DIRECT {
		format = domainmodels.ImageURLFormat_CANONICAL
	}
	example := domainmodels.Image{Identifier: docsExampleImageIdentifier, MIMEType: "image/png"}
	exampleURL := example.GetURL(conf.ImageDomain, IsSecure(r), []*domainmodels.ExternalImageIdentifier{{
		StorageDefinitionIdentifier: "s1",
		FileIdentifier:              docsExampleImageIdentifier + ".png",
	}}, format)
	return docsURLs{
		API:          strings.TrimSuffix(baseURLFromRequest(r).String(), "/"),
		WebUI:        strings.TrimSuffix(webUIBaseURL(r, conf.WebUIHost).String(), "/"),
		ExampleImage: absoluteImageURL(r, exampleURL),
	}
}

// renderAPIDocsMarkdown fills in the API reference's placeholders, so examples
// and limits match this instance and can be copied as they are.
func renderAPIDocsMarkdown(urls docsURLs) string {
	return strings.NewReplacer(
		"{{BASE_URL}}", urls.API,
		"{{WEB_UI_URL}}", urls.WebUI,
		"{{EXAMPLE_IMAGE_URL}}", urls.ExampleImage,
		"{{MAX_IMAGES}}", strconv.Itoa(apiMaxImagesPerRequest),
		"{{MIN_EXPIRES_IN}}", strconv.Itoa(apiMinExpiresInSeconds),
		"{{MAX_EXPIRES_IN}}", strconv.Itoa(apiMaxExpiresInSeconds),
		"{{MAX_CAPTION}}", strconv.Itoa(domainmodels.ImageMaxCaptionLength),
		"{{MAX_TITLE}}", strconv.Itoa(domainmodels.AlbumMaxTitleLength),
	).Replace(apiDocsMarkdown)
}

// renderLlmsTxt returns an llms.txt (https://llmstxt.org) pointing language
// models at the plain-text API reference.
func renderLlmsTxt(urls docsURLs) string {
	return "# imgdd\n\n" +
		"> Self-hosted image hosting with an HTTP API. AI agents can upload screenshots into a live album page that a person watches while the agent works.\n\n" +
		"## Docs\n\n" +
		"- [HTTP API](" + urls.API + apiDocsMarkdownPath + "): authentication, createAlbum, createImages, error codes, and a quick start for agents sharing screenshots\n"
}

func serveText(conf *HttpServerConfigDef, contentType string, render func(docsURLs) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Use GET", http.StatusMethodNotAllowed)
			return
		}
		body := render(docsURLsFor(r, conf))
		w.Header().Set("Content-Type", contentType)
		// The body is built from the request's Host and scheme, so a shared
		// cache must not hand one client's copy to another: a forged Host
		// would otherwise point every agent at someone else's server.
		w.Header().Set("Cache-Control", "private, max-age=300")
		w.Header().Set("Vary", "Host, X-Forwarded-Proto")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		if r.Method == http.MethodHead {
			return
		}
		w.Write([]byte(body))
	}
}

func makeAPIDocsMarkdownHandler(conf *HttpServerConfigDef) http.HandlerFunc {
	return serveText(conf, "text/markdown; charset=utf-8", renderAPIDocsMarkdown)
}

func makeLlmsTxtHandler(conf *HttpServerConfigDef) http.HandlerFunc {
	return serveText(conf, "text/plain; charset=utf-8", renderLlmsTxt)
}
