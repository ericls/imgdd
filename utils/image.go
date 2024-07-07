package utils

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func GetImageDimensions(imageBytes []byte) (int32, int32, error) {
	reader := bytes.NewReader(imageBytes)
	// Decode the image
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}

	if img.Width > 2147483647 || img.Height > 2147483647 {
		return 0, 0, errors.New("image dimensions are too large")
	}

	// Return the width and height
	return int32(img.Width), int32(img.Height), nil
}
