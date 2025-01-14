package domainmodels

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"time"
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
