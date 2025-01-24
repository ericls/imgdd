package imagechecks

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type SafeImageResp struct {
	Safe bool
}

func isImageSafe(file io.Reader, endpoint string) (bool, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "img.jpg")
	if err != nil {
		return false, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return false, err
	}
	err = writer.Close()
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var safeImageResp SafeImageResp
	err = json.Unmarshal(content, &safeImageResp)
	if err != nil {
		return false, err
	}
	return safeImageResp.Safe, nil
}

func MakeImageChecker(endpoint string) Checker {
	return func(file io.Reader) bool {
		safe, err := isImageSafe(file, endpoint)
		if err != nil {
			return false
		}
		return safe
	}
}
