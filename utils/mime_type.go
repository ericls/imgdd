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
