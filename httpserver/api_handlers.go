package httpserver

// The HTTP API is intention based, not resource based: every endpoint is named
// after the one thing a caller wants to do (createImages), and is invoked with
// POST at /api/v1/<intention>. There are no resource URLs whose meaning
// depends on the HTTP method. New capabilities get new intention endpoints.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/ericls/imgdd/storage"
)

const apiCreateImagesPath = "/api/v1/createImages"

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
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	Filename   string `json:"filename"`
	URL        string `json:"url"`
	MIMEType   string `json:"mimeType"`
	Width      int32  `json:"width"`
	Height     int32  `json:"height"`
	ByteSize   int32  `json:"byteSize"`
	// ExpiresAt is null for images that never expire.
	ExpiresAt *time.Time `json:"expiresAt"`
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
// and optionally "expires_in", the number of seconds after which all of the
// created images expire and are deleted.
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

		results := make([]apiCreateImagesResult, len(fileHeaders))
		status := 0
		for i, fh := range fileHeaders {
			results[i].Filename = fh.Filename
			uploaded, uploadErr := createUploadedImage(r, conf, orgUser, fh, storageDefRepo, imageRepo, limiter, expiresAt)
			if uploadErr != nil {
				results[i].Error = &apiError{Code: uploadErr.Code, Message: uploadErr.Message}
				if status == 0 {
					status = uploadErr.Status
				}
				continue
			}
			status = http.StatusOK
			results[i].Image = &apiImage{
				ID:         uploaded.ID,
				Identifier: uploaded.Identifier,
				Filename:   uploaded.Filename,
				URL:        absoluteImageURL(r, uploaded.URL),
				MIMEType:   uploaded.MIMEType,
				Width:      uploaded.Width,
				Height:     uploaded.Height,
				ByteSize:   uploaded.ByteSize,
				ExpiresAt:  uploaded.ExpiresAt,
			}
		}
		writeAPIJSON(w, status, apiCreateImagesResponse{Results: results})
	}
}
