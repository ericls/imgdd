package httpserver

import (
	"database/sql"
	"encoding/json"
	"imgdd/domainmodels"
	"imgdd/identity"
	"imgdd/image"
	"imgdd/storage"
	"imgdd/utils"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type UploadReturn struct {
	Identifier string `json:"identifier"`
}

func MakeUploadHandler(conn *sql.DB, identityManager *IdentityManager, storageRepo storage.StorageRepo, imageRepo image.ImageRepo) http.HandlerFunc {
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
