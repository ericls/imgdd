package httpserver

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/httpserver/imagechecks"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/utils"

	"github.com/google/uuid"
)

// uploadError describes why an uploaded file was rejected. Message is safe to
// show to clients; Code is a stable, machine-readable identifier.
type uploadError struct {
	Status  int
	Code    string
	Message string
}

func newUploadError(status int, code string, message string) *uploadError {
	return &uploadError{Status: status, Code: code, Message: message}
}

func internalUploadError() *uploadError {
	return newUploadError(http.StatusInternalServerError, "internal_error", "Internal error")
}

type uploadedImage struct {
	ID         string
	Identifier string
	Filename   string
	URL        string
	MIMEType   string
	Width      int32
	Height     int32
	ByteSize   int32
	ExpiresAt  *time.Time
	Caption    string
}

// uploadOptions are the optional properties given to every image created by
// one upload request. The zero value means none.
type uploadOptions struct {
	// ExpiresAt is nil for images that never expire.
	ExpiresAt *time.Time
	// AlbumId is the internal id of an album the uploader owns, or "".
	AlbumId string
	// Caption is a description of the image, or "".
	Caption string
}

// effectiveUploadLimit returns the maximum image size for orgUser (nil for
// guests): the guest limit (or global max if unset), overridden by a per-user
// limit that is capped at the global max.
func effectiveUploadLimit(conf *HttpServerConfigDef, orgUser *domainmodels.OrganizationUser) int64 {
	limit := conf.ImageMaxUploadBytes
	if conf.GuestImageMaxUploadBytes > 0 {
		limit = conf.GuestImageMaxUploadBytes
	}
	if orgUser != nil && orgUser.UploadLimitBytes != nil {
		limit = *orgUser.UploadLimitBytes
		if limit > conf.ImageMaxUploadBytes {
			limit = conf.ImageMaxUploadBytes
		}
	}
	return limit
}

// createUploadedImage validates and stores one uploaded file on behalf of
// orgUser (nil for guests). It is shared by every upload entry point so they
// apply the same checks. Callers must have checked that orgUser owns
// opts.AlbumId.
func createUploadedImage(
	r *http.Request,
	conf *HttpServerConfigDef,
	orgUser *domainmodels.OrganizationUser,
	fileHeader *multipart.FileHeader,
	storageDefRepo storage.StorageDefRepo,
	imageRepo image.ImageRepo,
	limiter *ratelimit.RateLimiter,
	opts uploadOptions,
) (*uploadedImage, *uploadError) {
	src, err := fileHeader.Open()
	if err != nil {
		httpLogger.Error().Err(err).Msg("Unable to open file")
		return nil, newUploadError(http.StatusBadRequest, "unreadable_file", "Unable to open file")
	}
	defer src.Close()

	imgBytes, err := io.ReadAll(src)
	if err != nil {
		httpLogger.Error().Err(err).Msg("Unable to read file")
		return nil, newUploadError(http.StatusBadRequest, "unreadable_file", "Unable to read file")
	}
	exifRemovedBytes, err := utils.MaybeRemoveExif(imgBytes)
	if err != nil || len(exifRemovedBytes) == 0 {
		exifRemovedBytes = imgBytes
	}

	bytesLength := len(exifRemovedBytes)
	if bytesLength == 0 {
		return nil, newUploadError(http.StatusBadRequest, "empty_file", "Empty file")
	}

	if int64(bytesLength) > effectiveUploadLimit(conf, orgUser) {
		return nil, newUploadError(http.StatusBadRequest, "file_too_large", "File too large")
	}

	declaredMimeType := utils.GetMimeTypeFromFilename(fileHeader.Filename)
	detectedMimeType := utils.DetectMIMEType(&exifRemovedBytes)
	if declaredMimeType != detectedMimeType {
		httpLogger.Warn().Msgf("Declared MIME type: %s, Detected MIME type: %s", declaredMimeType, detectedMimeType)
		return nil, newUploadError(http.StatusBadRequest, "mime_type_mismatch", "Declared MIME type does not match detected MIME type")
	}

	checkers := make([]imagechecks.Checker, 0)
	if conf.EnableSafeImageCheck && conf.SafeImageCheckEndpoint != "" {
		checkers = append(checkers, imagechecks.MakeImageChecker(conf.SafeImageCheckEndpoint))
	}
	if len(checkers) > 0 {
		imageReader := bytes.NewReader(exifRemovedBytes)
		ok := imagechecks.CheckAll(checkers, imageReader)
		if !ok {
			return nil, newUploadError(http.StatusBadRequest, "bad_image", "Bad image")
		}
	}

	uploaderIp := ExtractIP(r)

	// Find the best storage definition for storing the image
	storageDefs, err := storageDefRepo.ListStorageDefinitions()
	if err != nil {
		httpLogger.Error().Err(err).Msg("Unable to list storage definitions")
		return nil, internalUploadError()
	}
	var storageDef *domainmodels.StorageDefinition
	for _, def := range storageDefs {
		if def.IsEnabled {
			storageDef = def
			break
		}
	}
	if storageDef == nil {
		httpLogger.Error().Msg("No enabled storage definitions found")
		return nil, internalUploadError()
	}

	storageInstance, err := storage.GetStorage(storageDef)
	if err != nil {
		httpLogger.Error().Err(err).Msg("Unable to get storage")
		return nil, internalUploadError()
	}

	var createdById string
	if orgUser != nil {
		createdById = orgUser.Id
	}

	if orgUser == nil && limiter.IsRateLimited(uploaderIp) {
		return nil, newUploadError(http.StatusTooManyRequests, "rate_limited", "Rate limited")
	}

	imageIdentifier := uuid.New().String()
	width, height, err := utils.GetImageDimensions(exifRemovedBytes)
	if err != nil {
		httpLogger.Warn().Err(err).Msg("Unable to get image dimensions")
		return nil, newUploadError(http.StatusBadRequest, "invalid_image", "The file is not a valid image")
	}
	img := domainmodels.Image{
		UploaderIP:      uploaderIp,
		CreatedById:     createdById,
		MIMEType:        detectedMimeType,
		Name:            fileHeader.Filename,
		Identifier:      imageIdentifier,
		NominalByteSize: int32(len(exifRemovedBytes)),
		NominalWidth:    width,
		NominalHeight:   height,
		ExpiresAt:       opts.ExpiresAt,
		AlbumId:         opts.AlbumId,
		Caption:         opts.Caption,
	}
	storedImage, err := imageRepo.CreateAndSaveUploadedImage(&img, detectedMimeType, exifRemovedBytes, storageDef.Id, storageInstance.Save)
	if err != nil {
		httpLogger.Error().Stack().Err(err).Msg("Unable to save image")
		return nil, internalUploadError()
	}
	externalImageIdentifiers := []*domainmodels.ExternalImageIdentifier{
		{
			StorageDefinitionIdentifier: storageDef.Identifier,
			FileIdentifier:              storedImage.FileIdentifier,
		},
	}
	return &uploadedImage{
		ID:         storedImage.Image.Id,
		Identifier: storedImage.Image.Identifier,
		Filename:   storedImage.Image.Name,
		URL:        img.GetURL(conf.ImageDomain, IsSecure(r), externalImageIdentifiers, conf.DefaultURLFormat),
		MIMEType:   detectedMimeType,
		Width:      width,
		Height:     height,
		ByteSize:   int32(len(exifRemovedBytes)),
		ExpiresAt:  storedImage.Image.ExpiresAt,
		Caption:    storedImage.Image.Caption,
	}, nil
}
