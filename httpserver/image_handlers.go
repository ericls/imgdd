package httpserver

import (
	"encoding/json"
	"imgdd/domainmodels"
	"imgdd/identity"
	"imgdd/image"
	"imgdd/storage"
	"imgdd/utils"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type UploadReturn struct {
	Filename   string `json:"filename"`
	URL        string `json:"url"`
	Identifier string `json:"identifier"`
}

func getFileExtFromMIMEType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		exts, err := mime.ExtensionsByType(mimeType)
		if err == nil && len(exts) > 0 {
			return exts[0]
		}
	}
	return ""
}

func getImageURL(identifier string, mimeType string, isSecure bool) string {
	maybeImageDomain := Config.ImageDomain
	suffix := getFileExtFromMIMEType(mimeType)
	if suffix == "" {
		return ""
	}
	filename := identifier + suffix
	if maybeImageDomain != "" {
		if isSecure {
			return "https://" + maybeImageDomain + "/image/" + filename
		}
		return "http://" + maybeImageDomain + "/image/" + filename
	}
	return "/image/" + filename
}

func isSecure(
	r *http.Request,
) bool {
	return strings.HasPrefix(r.Proto, "HTTPS/")
}

func makeUploadHandler(
	identityManager *IdentityManager, storageRepo storage.StorageRepo, imageRepo image.ImageRepo,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB
		_, fileHeader, err := r.FormFile("image")
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to read file")
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		src, err := fileHeader.Open()
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to open file")
			http.Error(w, "Unable to open file", http.StatusBadRequest)
			return
		}
		defer src.Close()

		imgBytes, _ := io.ReadAll(src)
		exifRemovedBytes, err := utils.MaybeRemoveExif(imgBytes)
		if err != nil || len(exifRemovedBytes) == 0 {
			exifRemovedBytes = imgBytes
		}

		bytesLength := len(exifRemovedBytes)
		if bytesLength == 0 {
			http.Error(w, "Empty file", http.StatusBadRequest)
			return
		}
		if bytesLength > 10*1024*1024 { // sanity check after removing EXIF
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		declaredMimeType := utils.GetMimeTypeFromFilename(fileHeader.Filename)
		detectedMimeType := utils.DetectMIMEType(&exifRemovedBytes)
		if declaredMimeType != detectedMimeType {
			httpLogger.Warn().Msgf("Declared MIME type: %s, Detected MIME type: %s", declaredMimeType, detectedMimeType)
			http.Error(w, "Declared MIME type does not match detected MIME type", http.StatusBadRequest)
			return
		}

		uploaderIp := ExtractIP(r)

		// Find the best storage definition for storing the image
		storageDefs, err := storageRepo.ListStorageDefinitions()
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to list storage definitions")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		var storageDef *domainmodels.StorageDefinition
		for _, def := range storageDefs {
			if def.IsEnabled {
				storageDef = def
				break
			}
		}
		if storageDef == nil {
			httpLogger.Error().Err(err).Msg("No enabled storage definitions found")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		storageInstance, err := storage.GetStorage(storageDef)
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to get storage")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		orgUser := identity.GetCurrentOrganizationUser(identityManager.ContextUserManager, r.Context())
		var createdById string
		if orgUser != nil {
			createdById = orgUser.Id
		}

		imageIdentifier := uuid.New().String()
		width, height, err := utils.GetImageDimensions(exifRemovedBytes)
		if err != nil {
			httpLogger.Error().Err(err).Msg("Unable to get image dimensions")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		image := domainmodels.Image{
			UploaderIP:      uploaderIp,
			CreatedById:     createdById,
			MIMEType:        detectedMimeType,
			Name:            fileHeader.Filename,
			Identifier:      imageIdentifier,
			NominalByteSize: int32(len(exifRemovedBytes)),
			NominalWidth:    width,
			NominalHeight:   height,
		}
		storedImage, err := imageRepo.CreateAndSaveUploadedImage(&image, exifRemovedBytes, storageDef.Id, storageInstance.Save)
		if err != nil {
			httpLogger.Error().Stack().Err(err).Msg("Unable to save image")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		ret := UploadReturn{
			Filename:   storedImage.Image.Name,
			URL:        getImageURL(storedImage.Image.Identifier, declaredMimeType, isSecure(r)),
			Identifier: storedImage.Image.Identifier,
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

func makeImageHandler(
	storeRepo storage.StorageRepo,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path[len("/image/"):]
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
		storedImage, err := storeRepo.GetStoredImageByIdentifierAndMimeType(
			identifier,
			mimeType,
		)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Error().Str(
				"identifier",
				identifier,
			).Err(err).Msg("Unable to get stored image")
			return
		}
		storageInstance, err := storage.GetStorage(storedImage.StorageDefinition)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			httpLogger.Info().Err(err).Msg("Unable to get storage instance")
			return
		}
		meta := storageInstance.GetMeta(storedImage.FileIdentifier)
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
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Length", strconv.FormatInt(meta.ByteSize, 10))
		if meta.ETag != "" {
			w.Header().Set("ETag", "W/"+meta.ETag)
		}
		w.Header().Set("X-IMGDD-SI", storedImage.StorageDefinition.Identifier)
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodHead {
			return
		}
		io.Copy(w, reader)
	}
}
