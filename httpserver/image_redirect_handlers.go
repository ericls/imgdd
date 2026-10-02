package httpserver

import (
	"mime"
	"net/http"
	"strings"

	"github.com/ericls/imgdd/storage"
	"github.com/gorilla/mux"
)

// The redirect handlers receive the same URLs as the image handlers, but instead of
// serving the images they redirect visitors to the public image detail page.

const imageRedirectPathPrefix = "/i"

// imageDetailPathPrefix is the web UI's public image detail route, `routes.image` in
// web_client/src/routes.ts. TestImageDetailPathMatchesWebClientRoute keeps them in sync.
const imageDetailPathPrefix = "/images/"

func mountImageRedirectRoutes(
	router *mux.Router,
	webUIHost string,
	storageDefRepo storage.StorageDefRepo,
	storedImageRepo storage.StoredImageRepo,
) {
	redirectRouter := router.PathPrefix(imageRedirectPathPrefix).Subrouter()
	redirectRouter.PathPrefix("/image/").Handler(http.StripPrefix(imageRedirectPathPrefix, makeImageRedirectHandler(webUIHost, storedImageRepo)))
	redirectRouter.PathPrefix("/direct").Handler(http.StripPrefix(imageRedirectPathPrefix, makeDirectImageRedirectHandler(webUIHost, storageDefRepo, storedImageRepo)))
}

func getImageDetailURL(maybeWebUIHost string, imageId string, isSecure bool) string {
	path := imageDetailPathPrefix + imageId
	if maybeWebUIHost != "" {
		if isSecure {
			return "https://" + maybeWebUIHost + path
		}
		return "http://" + maybeWebUIHost + path
	}
	return path
}

func redirectToImageDetail(w http.ResponseWriter, r *http.Request, webUIHost string, imageId string) {
	http.Redirect(w, r, getImageDetailURL(webUIHost, imageId, IsSecure(r)), http.StatusFound)
}

func makeImageRedirectHandler(
	webUIHost string,
	storedImageRepo storage.StoredImageRepo,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/image/")
		if filename == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Msg("No filename")
			return
		}
		identifier, ext := splitIdentifierExt(filename)
		mimeType := mime.TypeByExtension("." + ext)
		if mimeType == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Msg("No MIME type")
			return
		}
		storedImages, err := storedImageRepo.GetStoredImageByIdentifierAndMimeType(identifier, mimeType)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str("identifier", identifier).Err(err).Msg("Unable to get stored image")
			return
		}
		for _, si := range storedImages {
			if si != nil && si.Image != nil && si.Image.Id != "" {
				redirectToImageDetail(w, r, webUIHost, si.Image.Id)
				return
			}
		}
		http.Error(w, "Not found", http.StatusNotFound)
		httpLogger.Info().Str("identifier", identifier).Msg("No stored image found")
	}
}

func makeDirectImageRedirectHandler(
	webUIHost string,
	storageDefRepo storage.StorageDefRepo,
	storedImageRepo storage.StoredImageRepo,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL format: <storage_definition_identifier>.<file_identifier>
		urlSegments := strings.Split(r.URL.Path, "/")
		lastSeg := urlSegments[len(urlSegments)-1]
		segments := strings.SplitN(lastSeg, ".", 2)
		if len(segments) != 2 {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Msg("URL is not right")
			return
		}
		storageDefIdentifier := segments[0]
		fileIdentifier := segments[1]
		storageDef, err := storageDefRepo.GetStorageDefinitionByIdentifier(storageDefIdentifier)
		if err != nil || storageDef == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str("storage_definition_identifier", storageDefIdentifier).Err(err).Msg("Unable to get storage definition")
			return
		}
		storedImage, err := storedImageRepo.GetStoredImageByStorageDefinitionIdAndFileIdentifier(storageDef.Id, fileIdentifier)
		if err != nil || storedImage == nil || storedImage.Image == nil || storedImage.Image.Id == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Str("file_identifier", fileIdentifier).Err(err).Msg("Unable to get stored image")
			return
		}
		redirectToImageDetail(w, r, webUIHost, storedImage.Image.Id)
	}
}
