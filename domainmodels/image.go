package domainmodels

type Image struct {
	Id        string
	StorageId string
	Filename  string
	MimeType  string
	ByteSize  int64
	ETag      string
}

type StoredImage struct {
	Id                  string
	StorageDefinitionId string
	Filename            string
	MimeType            string
	ETag                string
	ByteSize            int64
}
