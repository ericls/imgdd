package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/google/uuid"
)

const apiAddImagesToAlbumPath = "/api/v1/addImagesToAlbum"
const apiRemoveImagesFromAlbumPath = "/api/v1/removeImagesFromAlbum"

// apiMaxAlbumImagesPerRequest caps how many images one add or remove call
// accepts.
const apiMaxAlbumImagesPerRequest = 100

// apiMaxAlbumImagesBodyBytes caps the add/remove request body: 100 captions
// at their maximum length, plus room for ids and JSON.
const apiMaxAlbumImagesBodyBytes = 1 << 20

type apiAddImagesToAlbumRequest struct {
	Album  string `json:"album"`
	Images []struct {
		ID      string `json:"id"`
		Caption string `json:"caption"`
	} `json:"images"`
}

type apiRemoveImagesFromAlbumRequest struct {
	Album  string   `json:"album"`
	Images []string `json:"images"`
}

type apiAddImagesToAlbumResponse struct {
	Album apiAlbum `json:"album"`
	// Added counts images newly added; images already in the album are left
	// as they were, caption included.
	Added int `json:"added"`
}

type apiRemoveImagesFromAlbumResponse struct {
	Album   apiAlbum `json:"album"`
	Removed int      `json:"removed"`
}

// decodeStrictJSON decodes exactly one JSON value into dest, rejecting
// unknown fields so a misspelled field is an error rather than ignored.
func decodeStrictJSON(r *http.Request, dest any) bool {
	mediaType := strings.TrimSpace(strings.ToLower(strings.Split(r.Header.Get("Content-Type"), ";")[0]))
	if mediaType != "application/json" {
		return false
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		return false
	}
	_, err := decoder.Token()
	return errors.Is(err, io.EOF)
}

// albumEditRequest runs the checks shared by the calls that change an
// album's images: POST, an API key, the per-user rate limit, and a JSON body.
// It returns the caller, or nil after writing an error response.
func albumEditRequest(
	w http.ResponseWriter,
	r *http.Request,
	identityRepo identity.IdentityRepo,
	limiter *ratelimit.RateLimiter,
	dest any,
	shape string,
) *domainmodels.OrganizationUser {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST")
		return nil
	}
	orgUser, err := authenticateAPIRequest(r, identityRepo)
	if err != nil || orgUser == nil {
		w.Header().Set("WWW-Authenticate", `Bearer realm="imgdd"`)
		writeAPIError(w, http.StatusUnauthorized, "invalid_api_key", "An API key is required to change albums")
		return nil
	}
	if limiter.IsRateLimited(orgUser.Id) {
		writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "Rate limited")
		return nil
	}
	r.Body = http.MaxBytesReader(w, r.Body, apiMaxAlbumImagesBodyBytes)
	if !decodeStrictJSON(r, dest) {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Expected a JSON body like "+shape)
		return nil
	}
	return orgUser
}

// resolveOwnedAlbum returns the album with the given identifier if orgUser
// owns it, or nil after writing album_not_found: to anyone but the owner, the
// album does not exist.
func resolveOwnedAlbum(w http.ResponseWriter, identifier string, orgUser *domainmodels.OrganizationUser, imageRepo image.ImageRepo) *domainmodels.Album {
	album, uploadErr := resolveUploadAlbum([]string{identifier}, orgUser, imageRepo)
	if uploadErr != nil {
		writeAPIError(w, uploadErr.Status, uploadErr.Code, uploadErr.Message)
		return nil
	}
	return album
}

// parseAlbumImageIds validates the image ids of an add or remove call.
func parseAlbumImageIds(ids []string) ([]string, bool) {
	if len(ids) == 0 || len(ids) > apiMaxAlbumImagesPerRequest {
		return nil, false
	}
	seen := make(map[string]bool, len(ids))
	parsed := make([]string, len(ids))
	for i, id := range ids {
		u, err := uuid.Parse(strings.TrimSpace(id))
		if err != nil || seen[u.String()] {
			return nil, false
		}
		seen[u.String()] = true
		parsed[i] = u.String()
	}
	return parsed, true
}

func writeInvalidAlbumImages(w http.ResponseWriter) {
	writeAPIError(w, http.StatusBadRequest, "invalid_images", fmt.Sprintf(
		"Give 1 to %d distinct image ids, as returned by createImages", apiMaxAlbumImagesPerRequest,
	))
}

// makeAPIAddImagesToAlbumHandler serves POST /api/v1/addImagesToAlbum.
//
// Request (JSON): {"album": id, "images": [{"id": imageId, "caption": text}]}.
// The caller must own the album and every image, so nobody can put someone
// else's image on their album page. Images already in the album are skipped.
//
// Response: {"album": {"id", "url", "title"}, "added": n}.
func makeAPIAddImagesToAlbumHandler(
	conf *HttpServerConfigDef,
	identityRepo identity.IdentityRepo,
	imageRepo image.ImageRepo,
	limiter *ratelimit.RateLimiter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req apiAddImagesToAlbumRequest
		orgUser := albumEditRequest(w, r, identityRepo, limiter, &req, `{"album": "...", "images": [{"id": "...", "caption": "..."}]}`)
		if orgUser == nil {
			return
		}
		album := resolveOwnedAlbum(w, req.Album, orgUser, imageRepo)
		if album == nil {
			return
		}
		ids := make([]string, len(req.Images))
		for i, img := range req.Images {
			ids[i] = img.ID
		}
		ids, ok := parseAlbumImageIds(ids)
		if !ok {
			writeInvalidAlbumImages(w)
			return
		}
		entries := make([]domainmodels.AlbumImageEntry, len(ids))
		for i, img := range req.Images {
			caption := strings.TrimSpace(img.Caption)
			if utf8.RuneCountInString(caption) > domainmodels.ImageMaxCaptionLength {
				writeAPIError(w, http.StatusBadRequest, "invalid_caption", fmt.Sprintf(
					"Each caption must be at most %d characters", domainmodels.ImageMaxCaptionLength,
				))
				return
			}
			entries[i] = domainmodels.AlbumImageEntry{ImageId: ids[i], Caption: caption}
		}

		images, err := imageRepo.GetImagesByIds(ids)
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to get images")
			writeAPIError(w, http.StatusInternalServerError, "internal_error", "Internal error")
			return
		}
		owned := make(map[string]bool, len(images))
		for _, img := range images {
			if img != nil && img.CreatedById == orgUser.Id {
				owned[img.Id] = true
			}
		}
		for _, id := range ids {
			// Someone else's image is reported the same as a missing one.
			if !owned[id] {
				writeAPIError(w, http.StatusNotFound, "image_not_found", "Image "+id+" not found")
				return
			}
		}

		added, err := imageRepo.AddImagesToAlbum(album.Id, entries, orgUser.Id)
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to add images to album")
			writeAPIError(w, http.StatusInternalServerError, "internal_error", "Internal error")
			return
		}
		writeAPIJSON(w, http.StatusOK, apiAddImagesToAlbumResponse{
			Album: apiAlbumFor(r, conf, album),
			Added: added,
		})
	}
}

// makeAPIRemoveImagesFromAlbumHandler serves POST /api/v1/removeImagesFromAlbum.
//
// Request (JSON): {"album": id, "images": [imageId, ...]}. The caller must own
// the album. The images themselves are not deleted, and ids that are not in
// the album are ignored.
//
// Response: {"album": {"id", "url", "title"}, "removed": n}.
func makeAPIRemoveImagesFromAlbumHandler(
	conf *HttpServerConfigDef,
	identityRepo identity.IdentityRepo,
	imageRepo image.ImageRepo,
	limiter *ratelimit.RateLimiter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req apiRemoveImagesFromAlbumRequest
		orgUser := albumEditRequest(w, r, identityRepo, limiter, &req, `{"album": "...", "images": ["..."]}`)
		if orgUser == nil {
			return
		}
		album := resolveOwnedAlbum(w, req.Album, orgUser, imageRepo)
		if album == nil {
			return
		}
		ids, ok := parseAlbumImageIds(req.Images)
		if !ok {
			writeInvalidAlbumImages(w)
			return
		}
		removed, err := imageRepo.RemoveImagesFromAlbum(album.Id, ids)
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to remove images from album")
			writeAPIError(w, http.StatusInternalServerError, "internal_error", "Internal error")
			return
		}
		writeAPIJSON(w, http.StatusOK, apiRemoveImagesFromAlbumResponse{
			Album:   apiAlbumFor(r, conf, album),
			Removed: removed,
		})
	}
}
