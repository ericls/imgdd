package image

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	dm "imgdd/domainmodels"
	"imgdd/utils"
	"strings"
	"time"
)

type SaveFunc func(file utils.SeekerReader, filename string, mimeType string) error
type PaginationDirection string

const (
	PaginationDirectionAsc  PaginationDirection = "asc"
	PaginationDirectionDesc PaginationDirection = "desc"
)

type ListImagesFilters struct {
	NameContains string
	CreatedAtLte *time.Time
	CreatedAtGte *time.Time
	CreatedBy    *string
	Limit        int
}

// TODO: make a generic pagination utility
type ListImagesOrdering struct {
	ID        *PaginationDirection `json:"id,omitempty"`
	Name      *PaginationDirection `json:"name,omitempty"`
	CreatedAt *PaginationDirection `json:"createdAt,omitempty"`
	Checksum  string
}

type ImageCursor struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

func (c *ImageCursor) B64Encode() string {
	if c == nil {
		return ""
	}
	jsonBytes, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

func ImageCursorB64Decode(encoded string) *ImageCursor {
	if encoded == "" {
		return nil
	}
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	cursor := &ImageCursor{}
	_ = json.Unmarshal(decoded, cursor)
	return cursor
}

func (o *ListImagesOrdering) GetChecksum() string {
	if o == nil {
		return ""
	}
	if o.Checksum != "" {
		return o.Checksum
	}
	jsonBytes, _ := json.Marshal(o)
	h := fnv.New32a()
	h.Write(jsonBytes)
	// (0x%x)
	sum := h.Sum32()
	hexFormated := fmt.Sprintf("%x", sum)
	o.Checksum = hexFormated
	return hexFormated
}

func (o *ListImagesOrdering) GetCursor(i *dm.Image) string {
	if o == nil {
		return ""
	}
	if o.GetChecksum() == "" {
		return ""
	}
	usedFieldNames := make([]string, 0)
	if o.ID != nil {
		usedFieldNames = append(usedFieldNames, "ID")
	}
	if o.Name != nil {
		usedFieldNames = append(usedFieldNames, "Name")
	}
	if o.CreatedAt != nil {
		usedFieldNames = append(usedFieldNames, "CreatedAt")
	}
	cursor := &ImageCursor{}
	for _, fieldName := range usedFieldNames {
		switch fieldName {
		case "ID":
			cursor.ID = i.Id
		case "Name":
			cursor.Name = i.Name
		case "CreatedAt":
			cursor.CreatedAt = i.CreatedAt.Format(time.RFC3339)
		}
	}
	return o.Checksum + "|" + cursor.B64Encode()
}

func (o *ListImagesOrdering) GetImageCursor(encoded string) *ImageCursor {
	if encoded == "" {
		return nil
	}
	parts := strings.SplitN(encoded, "|", 2)
	if len(parts) != 2 {
		return nil
	}
	if parts[0] != o.GetChecksum() {
		return nil
	}
	return ImageCursorB64Decode(parts[1])
}

type ImageRepo interface {
	CreateAndSaveUploadedImage(image *dm.Image, fileBytes []byte, storageDefinitionId string, saveFn SaveFunc) (*dm.StoredImage, error)
	ListImages(filters ListImagesFilters, ordering ListImagesOrdering) (dm.ListImageResult, error)
}
