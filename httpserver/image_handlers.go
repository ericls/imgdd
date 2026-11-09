package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/ericls/imgdd/storage"
)

type UploadReturn struct {
	ID         string `json:"id"`
	Filename   string `json:"filename"`
	URL        string `json:"url"`
	Identifier string `json:"identifier"`
}

// makeUploadHandler serves /upload, which the web client uses. It
// authenticates with the browser session. Third-party clients should use the
// documented API at /api/v1/createImages instead.
func makeUploadHandler(
	conf *HttpServerConfigDef,
	identityManager *IdentityManager,
	storageDefRepo storage.StorageDefRepo,
	imageRepo image.ImageRepo,
	limiter *ratelimit.RateLimiter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !conf.AllowUpload {
			http.Error(w, "Image upload is disabled", http.StatusForbidden)
			return
		}
		r.ParseMultipartForm(conf.ImageMaxUploadBytes)
		_, fileHeader, err := r.FormFile("image")
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to read file")
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		orgUser := identity.GetCurrentOrganizationUser(identityManager.ContextUserManager, r.Context())
		uploaded, uploadErr := createUploadedImage(r, conf, orgUser, fileHeader, storageDefRepo, imageRepo, limiter, nil)
		if uploadErr != nil {
			http.Error(w, uploadErr.Message, uploadErr.Status)
			return
		}
		ret := UploadReturn{
			ID:         uploaded.ID,
			Filename:   uploaded.Filename,
			URL:        uploaded.URL,
			Identifier: uploaded.Identifier,
		}
		serialized, err := json.Marshal(ret)
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to serialize response")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(serialized)
	}
}

func splitIdentifierExt(filename string) (string, string) {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return parts[0], ""
	}
	return strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1]

}

type storedImageWithStorageDef struct {
	*domainmodels.StoredImage
	*domainmodels.StorageDefinition
}

func copyImageResponse(
	w http.ResponseWriter,
	reader io.Reader,
	imageResponseCache *imageResponseCache,
	cacheKey string,
	contentType string,
	etag string,
	xImgddSI string,
) {
	if imageResponseCache == nil {
		io.Copy(w, reader)
		return
	}
	body, err := io.ReadAll(io.LimitReader(reader, imageResponseCache.maxFileBytes+1))
	if err != nil {
		httpLogger.Info().Err(err).Msg("Unable to read image")
		return
	}
	if int64(len(body)) > imageResponseCache.maxFileBytes {
		w.Write(body)
		io.Copy(w, reader)
		return
	}
	imageResponseCache.put(cacheKey, contentType, etag, xImgddSI, body)
	w.Write(body)
}

func sortStoredImages(storedImages []*domainmodels.StoredImage, storageDefRepo storage.StorageDefRepo) ([]storedImageWithStorageDef, error) {
	var storedImageWithStorageDefs []storedImageWithStorageDef
	storageDefIds := make([]string, 0)
	for _, si := range storedImages {
		if si == nil {
			continue
		}
		storageDefIds = append(storageDefIds, si.StorageDefinitionId)
	}
	storageDefs, err := storageDefRepo.GetStorageDefinitionsByIds(storageDefIds)
	if err != nil {
		return nil, err
	}
	idToStorageDef := make(map[string]*domainmodels.StorageDefinition)
	for _, sd := range storageDefs {
		if sd == nil {
			continue
		}
		idToStorageDef[sd.Id] = sd
	}
	for _, si := range storedImages {
		if si == nil {
			continue
		}
		storageDef := idToStorageDef[si.StorageDefinitionId]
		if storageDef == nil {
			continue
		}
		if !storageDef.IsEnabled {
			continue
		}
		storedImageWithStorageDefs = append(storedImageWithStorageDefs, storedImageWithStorageDef{
			StoredImage:       si,
			StorageDefinition: storageDef,
		})
	}
	if len(storedImageWithStorageDefs) == 0 {
		return nil, fmt.Errorf("no enabled storage definitions found")
	}
	sort.SliceStable(storedImageWithStorageDefs, func(i, j int) bool {
		return storedImageWithStorageDefs[i].StorageDefinition.Priority < storedImageWithStorageDefs[j].StorageDefinition.Priority
	})
	return storedImageWithStorageDefs, nil
}

func makeImageHandler(
	storageDefRepo storage.StorageDefRepo,
	storedImageRepo storage.StoredImageRepo,
	imageResponseCache *imageResponseCache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path[len("/image/"):]
		if filename == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Msg("No filename")
			return
		}
		cacheKey := "image:" + filename
		if entry, ok := imageResponseCache.get(cacheKey); ok {
			writeCachedImageResponse(w, r, entry)
			return
		}
		identifier, ext := splitIdentifierExt(filename)
		mimeType := mime.TypeByExtension("." + ext)
		if mimeType == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Msg("No MIME type")
			return
		}
		providedEtag := r.Header.Get("If-None-Match")
		if providedEtag == filename {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		storedImages, err := storedImageRepo.GetStoredImageByIdentifierAndMimeType(
			identifier,
			mimeType,
		)
		if err != nil || len(storedImages) == 0 {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str(
				"identifier",
				identifier,
			).Err(err).Msg("Unable to get stored image")
			return
		}
		storedImageWithStorageDefs, err := sortStoredImages(storedImages, storageDefRepo)
		if len(storedImageWithStorageDefs) == 0 {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Str(
				"identifier", identifier,
			).Msg("No enabled storage definitions found")
			return
		}
		storedImage := storedImageWithStorageDefs[0].StoredImage
		storageDef := storedImageWithStorageDefs[0].StorageDefinition
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str(
				"storage_definition_id",
				storedImage.StorageDefinitionId,
			).Err(err).Msg("Unable to get storage definition")
		}
		storageInstance, err := storage.GetStorage(storageDef)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Err(err).Msg("Unable to get storage instance")
			return
		}
		meta := storage.GetMetaCached(storageInstance, storageDef.Id, storedImage.FileIdentifier)
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Length", strconv.FormatInt(meta.ByteSize, 10))
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("ETag", filename)
		w.Header().Set("X-imgdd-si", storageDef.Identifier)
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodHead {
			return
		}
		reader := storageInstance.GetReader(storedImage.FileIdentifier)
		if reader == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Str(
				"file_identifier",
				storedImage.FileIdentifier,
			).Msg("Unable to get reader")
			return
		}
		defer reader.Close()
		copyImageResponse(w, reader, imageResponseCache, cacheKey, mimeType, filename, storageDef.Identifier)
	}
}

func makeDirectImageHandler(
	storageDefRepo storage.StorageDefRepo,
	imageResponseCache *imageResponseCache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL format: <storage_definition_identifier>.<file_identifier>
		urlSegments := strings.Split(r.URL.Path, "/")
		if len(urlSegments) < 1 {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Msg("URL is not right")
			return
		}
		lastSeg := urlSegments[len(urlSegments)-1]
		segments := strings.SplitN(lastSeg, ".", 2)
		if len(segments) != 2 {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Msg("URL is not right")
			return
		}
		storageDefIdentifier := segments[0]
		fileIdentifier := segments[1]
		cacheKey := "direct:" + storageDefIdentifier + ":" + fileIdentifier
		if entry, ok := imageResponseCache.get(cacheKey); ok {
			writeCachedImageResponse(w, r, entry)
			return
		}
		if r.Header.Get("If-None-Match") == fileIdentifier {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		storageDef, err := storageDefRepo.GetStorageDefinitionByIdentifier(storageDefIdentifier)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str("storage_definition_identifier", storageDefIdentifier).Err(err).Msg("Unable to get storage definition")
			return
		}
		if !storageDef.IsEnabled {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str("storage_definition_identifier", storageDefIdentifier).Msg("Storage definition is not enabled")
			return
		}
		storageInstance, err := storage.GetStorage(storageDef)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str("storage_definition_identifier", storageDefIdentifier).Err(err).Msg("Unable to get storage instance")
			return
		}
		meta := storage.GetMetaCached(storageInstance, storageDef.Id, fileIdentifier)
		w.Header().Set("Content-Type", meta.ContentType)
		w.Header().Set("Content-Length", strconv.FormatInt(meta.ByteSize, 10))
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("ETag", fileIdentifier)
		w.Header().Set("X-imgdd-si", storageDef.Identifier)
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodHead {
			return
		}
		reader := storageInstance.GetReader(fileIdentifier)
		if reader == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Str("file_identifier", fileIdentifier).Msg("Unable to get reader")
			return
		}
		defer reader.Close()
		copyImageResponse(w, reader, imageResponseCache, cacheKey, meta.ContentType, fileIdentifier, storageDef.Identifier)
	}
}
