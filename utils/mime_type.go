package utils

import (
	"mime"
	"net/http"
	"strings"
)

func DetectMIMEType(data *([]byte)) string {
	buff := make([]byte, 512)
	copy(buff, *data)
	return http.DetectContentType(buff)
}

func GetExtFromFilename(filename string) string {
	ext := filename[strings.LastIndex(filename, ".")+1:]
	return ext
}

func GetMimeTypeFromFilename(filename string) string {
	ext := GetExtFromFilename(filename)
	mime_type := mime.TypeByExtension("." + ext)
	return mime_type
}

func GetExtFromMIMEType(mimeType string) string {
	if mimeType == "image/jpeg" {
		return ".jpg"
	} else if mimeType == "image/png" {
		return ".png"
	} else if mimeType == "image/gif" {
		return ".gif"
	} else if mimeType == "image/bmp" {
		return ".bmp"
	} else if mimeType == "image/webp" {
		return ".webp"
	} else if mimeType == "image/tiff" {
		return ".tiff"
	}

	if strings.HasPrefix(mimeType, "image/") {
		exts, err := mime.ExtensionsByType(mimeType)
		if err == nil && len(exts) > 0 {
			return exts[0]
		}
	}
	return ""
}
