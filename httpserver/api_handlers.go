package httpserver

// The HTTP API is intention based, not resource based: every endpoint is named
// after the one thing a caller wants to do (createImages), and is invoked with
// POST at /api/v1/<intention>. There are no resource URLs whose meaning
// depends on the HTTP method. New capabilities get new intention endpoints.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/ericls/imgdd/storage"
)

const apiCreateImagesPath = "/api/v1/createImages"
const apiCreateAlbumPath = "/api/v1/createAlbum"

// apiMaxCreateAlbumBodyBytes caps the createAlbum request body, which only
// carries a title.
const apiMaxCreateAlbumBodyBytes = 64 << 10

// apiMaxImagesPerRequest caps how many images one createImages call accepts.
const apiMaxImagesPerRequest = 5

// Bounds for the optional expires_in field, in seconds.
const apiMinExpiresInSeconds = 60
const apiMaxExpiresInSeconds = 365 * 24 * 60 * 60

// apiMultipartOverheadBytes is room for multipart boundaries and headers on
// top of the image bytes themselves.
const apiMultipartOverheadBytes = 1 << 20

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiErrorResponse struct {
	Error apiError `json:"error"`
}

type apiImage struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	MIMEType string `json:"mimeType"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	ByteSize int32  `json:"byteSize"`
	// ExpiresAt is null for images that never expire.
	ExpiresAt *time.Time `json:"expiresAt"`
	// PageURL is the image's page in the web UI.
	PageURL string `json:"pageUrl"`
	// Caption is null for images without one.
	Caption *string `json:"caption"`
	// AlbumURL is null for images uploaded without an album.
	AlbumURL *string `json:"albumUrl"`
}

type apiAlbum struct {
	// ID is the album's public identifier, used as the "album" field of
	// createImages. The internal id is never exposed.
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

func apiAlbumFor(r *http.Request, conf *HttpServerConfigDef, album *domainmodels.Album) apiAlbum {
	return apiAlbum{
		ID:    album.Identifier,
		URL:   albumPageURL(r, conf.WebUIHost, album.Identifier),
		Title: album.Title,
	}
}

type apiCreateAlbumResponse struct {
	Album apiAlbum `json:"album"`
}

type apiCreateImagesResult struct {
	Filename string    `json:"filename"`
	Image    *apiImage `json:"image,omitempty"`
	Error    *apiError `json:"error,omitempty"`
}

type apiCreateImagesResponse struct {
	Results []apiCreateImagesResult `json:"results"`
}

func writeAPIJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeAPIError(w http.ResponseWriter, status int, code string, message string) {
	writeAPIJSON(w, status, apiErrorResponse{Error: apiError{Code: code, Message: message}})
}

// absoluteImageURL makes a relative image URL (used when no image domain is
// configured) absolute, since API clients have no page to resolve it against.
func absoluteImageURL(r *http.Request, imageURL string) string {
	parsed, err := url.Parse(imageURL)
	if err != nil || parsed.IsAbs() {
		return imageURL
	}
	return baseURLFromRequest(r).ResolveReference(parsed).String()
}

// webUIBaseURL is the base URL of the web UI: the configured web UI host
// when there is one (the API may be served from another host), otherwise the
// host of the request.
func webUIBaseURL(r *http.Request, webUIHost string) *url.URL {
	base := baseURLFromRequest(r)
	if webUIHost != "" {
		base.Host = webUIHost
	}
	return base
}

// imagePageURL is the absolute URL of an image's page in the web UI.
func imagePageURL(r *http.Request, webUIHost string, imageId string) string {
	return absoluteImageURL(r, getImageDetailURL(webUIHost, imageId, IsSecure(r)))
}

// albumPageURL is the absolute URL of an album's page in the web UI.
func albumPageURL(r *http.Request, webUIHost string, albumIdentifier string) string {
	return webUIBaseURL(r, webUIHost).JoinPath("albums", albumIdentifier).String()
}

// resolveUploadAlbum reads the optional album form field. It returns the
// album the images should be added to, nil when the field is absent, or an
// error. Only the album's owner may add to it; to everyone else, including
// anonymous callers, the album does not exist.
func resolveUploadAlbum(values []string, orgUser *domainmodels.OrganizationUser, imageRepo image.ImageRepo) (*domainmodels.Album, *uploadError) {
	if len(values) == 0 {
		return nil, nil
	}
	if len(values) > 1 {
		return nil, newUploadError(http.StatusBadRequest, "invalid_album", "Give at most one album")
	}
	notFound := newUploadError(http.StatusNotFound, "album_not_found", "Album not found")
	identifier := strings.TrimSpace(values[0])
	if identifier == "" || orgUser == nil {
		return nil, notFound
	}
	album, err := imageRepo.GetAlbumByIdentifier(identifier)
	if err != nil {
		httpLogger.Error().Err(err).Msg("Unable to get album")
		return nil, internalUploadError()
	}
	if album == nil || album.CreatedById != orgUser.Id {
		return nil, notFound
	}
	return album, nil
}

// parseCaptions reads the optional caption form field for fileCount files: no
// value, one value for every file, or exactly one value per file in request
// order. It returns one caption per file ("" for none).
func parseCaptions(values []string, fileCount int) ([]string, bool) {
	captions := make([]string, fileCount)
	switch len(values) {
	case 0:
		return captions, true
	case 1:
		for i := range captions {
			captions[i] = values[0]
		}
	case fileCount:
		copy(captions, values)
	default:
		return nil, false
	}
	for i, caption := range captions {
		caption = strings.TrimSpace(caption)
		if utf8.RuneCountInString(caption) > domainmodels.ImageMaxCaptionLength {
			return nil, false
		}
		captions[i] = caption
	}
	return captions, true
}

// parseExpiresIn reads the optional expires_in form field (seconds from now)
// and returns when the images should expire, or nil if the field is absent.
func parseExpiresIn(values []string, now time.Time) (*time.Time, bool) {
	if len(values) == 0 {
		return nil, true
	}
	if len(values) > 1 {
		return nil, false
	}
	seconds, err := strconv.ParseInt(strings.TrimSpace(values[0]), 10, 64)
	if err != nil || seconds < apiMinExpiresInSeconds || seconds > apiMaxExpiresInSeconds {
		return nil, false
	}
	expiresAt := now.Add(time.Duration(seconds) * time.Second).UTC()
	return &expiresAt, true
}

var errInvalidAuthorizationHeader = errors.New("invalid authorization header")

// authenticateAPIRequest resolves the caller from an "Authorization: Bearer
// <api key>" header. It returns (nil, nil) for anonymous requests, which carry
// no Authorization header at all. Browser sessions are deliberately not
// accepted by the API.
func authenticateAPIRequest(r *http.Request, identityRepo identity.IdentityRepo) (*domainmodels.OrganizationUser, error) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return nil, nil
	}
	scheme, secret, ok := strings.Cut(header, " ")
	secret = strings.TrimSpace(secret)
	if !ok || !strings.EqualFold(scheme, "Bearer") || secret == "" {
		return nil, errInvalidAuthorizationHeader
	}
	orgUser := identityRepo.GetOrganizationUserByAPIKeySecret(secret)
	if orgUser == nil {
		return nil, errInvalidAuthorizationHeader
	}
	return orgUser, nil
}

// makeAPICreateImagesHandler serves POST /api/v1/createImages.
//
// Request: multipart/form-data with one or more files in the "image" field,
// and optionally:
//   - "expires_in": the number of seconds after which all of the created
//     images expire and are deleted.
//   - "album": the id of an album, owned by the caller, to add the images to.
//   - "caption": one caption for every image, or one per image in order.
//
// Authentication is optional; without an API key the upload is anonymous and
// subject to guest limits and rate limiting.
//
// Response: {"results": [...]} with one entry per file, in request order, each
// holding either "image" or "error". The status is 200 when at least one image
// was created; otherwise it is the status of the first failure.
func makeAPICreateImagesHandler(
	conf *HttpServerConfigDef,
	identityRepo identity.IdentityRepo,
	storageDefRepo storage.StorageDefRepo,
	imageRepo image.ImageRepo,
	limiter *ratelimit.RateLimiter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST")
			return
		}
		orgUser, err := authenticateAPIRequest(r, identityRepo)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Bearer realm="imgdd"`)
			writeAPIError(w, http.StatusUnauthorized, "invalid_api_key", "The API key is missing, malformed, or revoked")
			return
		}
		if !conf.AllowUpload {
			writeAPIError(w, http.StatusForbidden, "upload_disabled", "Image upload is disabled")
			return
		}

		maxBodyBytes := apiMaxImagesPerRequest*conf.ImageMaxUploadBytes + apiMultipartOverheadBytes
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		err = r.ParseMultipartForm(conf.ImageMaxUploadBytes)
		// The form can be populated even when an error is returned, for example
		// when only the query string is malformed.
		if r.MultipartForm != nil {
			defer r.MultipartForm.RemoveAll()
		}
		if err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				writeAPIError(w, http.StatusRequestEntityTooLarge, "request_too_large", "Request body is too large")
				return
			}
			writeAPIError(w, http.StatusBadRequest, "invalid_request", "Expected a multipart/form-data body")
			return
		}

		fileHeaders := r.MultipartForm.File["image"]
		if len(fileHeaders) == 0 {
			writeAPIError(w, http.StatusBadRequest, "no_images", `Attach at least one file in the "image" field`)
			return
		}
		if len(fileHeaders) > apiMaxImagesPerRequest {
			writeAPIError(w, http.StatusBadRequest, "too_many_images", "Too many images in one request")
			return
		}

		expiresAt, ok := parseExpiresIn(r.MultipartForm.Value["expires_in"], time.Now())
		if !ok {
			writeAPIError(w, http.StatusBadRequest, "invalid_expires_in", fmt.Sprintf(
				"expires_in must be a whole number of seconds between %d and %d",
				apiMinExpiresInSeconds, apiMaxExpiresInSeconds,
			))
			return
		}

		album, albumErr := resolveUploadAlbum(r.MultipartForm.Value["album"], orgUser, imageRepo)
		if albumErr != nil {
			writeAPIError(w, albumErr.Status, albumErr.Code, albumErr.Message)
			return
		}
		captions, ok := parseCaptions(r.MultipartForm.Value["caption"], len(fileHeaders))
		if !ok {
			writeAPIError(w, http.StatusBadRequest, "invalid_caption", fmt.Sprintf(
				"Give one caption for all images or one per image, each at most %d characters",
				domainmodels.ImageMaxCaptionLength,
			))
			return
		}
		if album == nil {
			for _, caption := range captions {
				if caption != "" {
					writeAPIError(w, http.StatusBadRequest, "invalid_caption",
						"A caption describes an image in an album; give album too")
					return
				}
			}
		}
		opts := uploadOptions{ExpiresAt: expiresAt}
		var albumURL *string
		if album != nil {
			opts.AlbumId = album.Id
			u := albumPageURL(r, conf.WebUIHost, album.Identifier)
			albumURL = &u
		}

		results := make([]apiCreateImagesResult, len(fileHeaders))
		status := 0
		for i, fh := range fileHeaders {
			results[i].Filename = fh.Filename
			opts.Caption = captions[i]
			uploaded, uploadErr := createUploadedImage(r, conf, orgUser, fh, storageDefRepo, imageRepo, limiter, opts)
			if uploadErr != nil {
				results[i].Error = &apiError{Code: uploadErr.Code, Message: uploadErr.Message}
				if status == 0 {
					status = uploadErr.Status
				}
				continue
			}
			status = http.StatusOK
			results[i].Image = &apiImage{
				ID:        uploaded.ID,
				Filename:  uploaded.Filename,
				URL:       absoluteImageURL(r, uploaded.URL),
				MIMEType:  uploaded.MIMEType,
				Width:     uploaded.Width,
				Height:    uploaded.Height,
				ByteSize:  uploaded.ByteSize,
				ExpiresAt: uploaded.ExpiresAt,
				PageURL:   imagePageURL(r, conf.WebUIHost, uploaded.ID),
				AlbumURL:  albumURL,
			}
			if uploaded.Caption != "" {
				results[i].Image.Caption = &uploaded.Caption
			}
		}
		writeAPIJSON(w, status, apiCreateImagesResponse{Results: results})
	}
}

// readCreateAlbumTitle reads the title from a JSON ({"title": "..."}) or form
// encoded createAlbum body. A missing title means an untitled album.
func readCreateAlbumTitle(r *http.Request) (string, bool) {
	switch requestMediaType(r) {
	case "application/json":
		var body struct {
			Title string `json:"title"`
		}
		// Strict, so {"name": ...} is an error rather than an untitled album.
		// An empty body is an untitled album.
		if err := decodeStrictJSON(r.Body, &body); err != nil && !errors.Is(err, io.EOF) {
			return "", false
		}
		return body.Title, true
	case "multipart/form-data":
		if err := r.ParseMultipartForm(apiMaxCreateAlbumBodyBytes); err != nil {
			return "", false
		}
		defer r.MultipartForm.RemoveAll()
	case "application/x-www-form-urlencoded", "":
		if err := r.ParseForm(); err != nil {
			return "", false
		}
	default:
		return "", false
	}
	values := r.PostForm["title"]
	if len(values) > 1 {
		return "", false
	}
	if len(values) == 0 {
		return "", true
	}
	return values[0], true
}

// makeAPICreateAlbumHandler serves POST /api/v1/createAlbum.
//
// Request: an optional "title", as JSON or form encoded. An API key is
// required, since only the album's owner can add images to it.
//
// Response: {"album": {"id", "url", "title"}}. Pass the id as the "album"
// field of createImages; share the url with whoever should see the album.
func makeAPICreateAlbumHandler(
	conf *HttpServerConfigDef,
	identityRepo identity.IdentityRepo,
	imageRepo image.ImageRepo,
	limiter *ratelimit.RateLimiter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgUser := apiKeyRequest(w, r, identityRepo, limiter, "to create albums")
		if orgUser == nil {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, apiMaxCreateAlbumBodyBytes)
		title, ok := readCreateAlbumTitle(r)
		if !ok {
			writeAPIError(w, http.StatusBadRequest, "invalid_request", `Expected an optional "title" as JSON or form data`)
			return
		}
		title = strings.TrimSpace(title)
		if utf8.RuneCountInString(title) > domainmodels.AlbumMaxTitleLength {
			writeAPIError(w, http.StatusBadRequest, "invalid_title", fmt.Sprintf(
				"title must be at most %d characters", domainmodels.AlbumMaxTitleLength,
			))
			return
		}
		album, err := imageRepo.CreateAlbum(title, orgUser.Id)
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to create album")
			writeAPIError(w, http.StatusInternalServerError, "internal_error", "Internal error")
			return
		}
		writeAPIJSON(w, http.StatusOK, apiCreateAlbumResponse{Album: apiAlbumFor(r, conf, album)})
	}
}
