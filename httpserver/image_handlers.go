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
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/utils"

	"github.com/google/uuid"
)

type UploadReturn struct {
	Filename   string `json:"filename"`
	URL        string `json:"url"`
	Identifier string `json:"identifier"`
}

func makeUploadHandler(
	conf *HttpServerConfigDef,
	identityManager *IdentityManager,
	storageDefRepo storage.StorageDefRepo,
	imageRepo image.ImageRepo,
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
		storageDefs, err := storageDefRepo.ListStorageDefinitions()
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
			URL:        image.GetURL(conf.ImageDomain, IsSecure(r)),
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

type storedImageWithStorageDef struct {
	*domainmodels.StoredImage
	*domainmodels.StorageDefinition
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
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("ETag", filename)
		w.Header().Set("X-imgdd-si", storageDef.Identifier)
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodHead {
			return
		}
		io.Copy(w, reader)
	}
}
