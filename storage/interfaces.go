package storage

import (
	"io"
)

type FileMeta struct {
	ByteSize    int64
	ContentType string
	ETag        string
}

type SeekerReader interface {
	io.Seeker
	io.Reader
}

type Storage interface {
	GetReader(filename string) io.Reader
	Save(file SeekerReader, filename string, mimeType string) error
	GetMeta(filename string) FileMeta
	Delete(filename string) error
	CheckConnection() error
}

type StorageBackend interface {
	FromJSON(jsonConfig []byte) Storage
}
