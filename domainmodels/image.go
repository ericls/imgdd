package domainmodels

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/ericls/imgdd/utils"
)

type ImageURLFormat string

const (
	ImageURLFormat_CANONICAL      ImageURLFormat = "canonical"
	ImageURLFormat_DIRECT         ImageURLFormat = "direct"
	ImageURLFormat_BACKEND_DIRECT ImageURLFormat = "backend_direct"
)

type Image struct {
	Id              string
	CreatedById     string
	CreatedAt       time.Time
	Name            string
	Identifier      string
	RootId          string
	ParentId        string
	UploaderIP      string
	MIMEType        string
	NominalWidth    int32
	NominalHeight   int32
	NominalByteSize int32
}

func (i *Image) HashStr() string {
	h := fnv.New64a()
	h.Write([]byte(i.Id))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i.NominalWidth))
	h.Write(buf)
	binary.BigEndian.PutUint32(buf, uint32(i.NominalHeight))
	h.Write(buf)
	binary.BigEndian.PutUint32(buf, uint32(i.NominalByteSize))
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (i *Image) GetURL(imageDomain string, isSecure bool, storedImages []*ExternalImageIdentifier, format ImageURLFormat) string {
	if format == ImageURLFormat_CANONICAL {
		return getCanonicalImageURL(imageDomain, i.Identifier, i.MIMEType, isSecure)
	} else if format == ImageURLFormat_DIRECT {
		if len(storedImages) > 0 {
			externalImage := storedImages[0]
			return getDirectImageURL(imageDomain, isSecure, externalImage.StorageDefinitionIdentifier, externalImage.FileIdentifier)
		} else {
			return ""
		}
	} else if format == ImageURLFormat_BACKEND_DIRECT {
		panic("not implemented")
	}
	return ""
}

func getCanonicalImageURL(maybeImageDomain string, identifier string, mimeType string, isSecure bool) string {
	suffix := utils.GetExtFromMIMEType(mimeType)
	if suffix == "" {
		return ""
	}
	filename := identifier + suffix
	if maybeImageDomain != "" {
		if isSecure {
			return "https://" + maybeImageDomain + "/" + filename
		}
		return "http://" + maybeImageDomain + "/" + filename
	}
	return "/image/" + filename
}

func getDirectImageURL(maybeImageDomain string, isSecure bool, storageDefinitionIdentifier string, fileIdentifier string) string {
	filename := fileIdentifier
	if maybeImageDomain != "" {
		if isSecure {
			return "https://" + maybeImageDomain + "/" + storageDefinitionIdentifier + "." + filename
		}
		return "http://" + maybeImageDomain + "/" + storageDefinitionIdentifier + "." + filename
	}
	return "/direct/" + storageDefinitionIdentifier + "." + filename
}
