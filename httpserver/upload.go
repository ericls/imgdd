package httpserver

import (
	"imgdd/utils"
	"io"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB
	_, fileHeader, err := r.FormFile("myFile")
	if err != nil {
		httpLogger.Err(err).Msg("Unable to read file")
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	src, err := fileHeader.Open()
	if err != nil {
		httpLogger.Err(err).Msg("Unable to open file")
		http.Error(w, "Unable to open file", http.StatusBadRequest)
		return
	}
	defer src.Close()

	imgBytes, _ := io.ReadAll(src)
	exifRemovedBytes, err := utils.MaybeRemoveExif(imgBytes)
	if err != nil || len(exifRemovedBytes) == 0 {
		exifRemovedBytes = imgBytes
	}

	declaredMimeType := utils.GetMimeTypeFromFilename(fileHeader.Filename)
	detectedMimeType := utils.DetectMIMEType(&exifRemovedBytes)
	if declaredMimeType != detectedMimeType {
		httpLogger.Warn().Msgf("Declared MIME type: %s, Detected MIME type: %s", declaredMimeType, detectedMimeType)
		http.Error(w, "Declared MIME type does not match detected MIME type", http.StatusBadRequest)
		return
	}

	uploaderIP := ExtractIP(r)

	// Find the best storage definition for storing the image
	// Create an image record
	// Allocate a stored image record
	// Save the image to the storage backend <- returns an identifier
	// Save the stored image record with
}
