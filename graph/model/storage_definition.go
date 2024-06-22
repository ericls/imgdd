package model

type StorageTypeEnum string

const (
	StorageType_S3 StorageTypeEnum = "S3"
)

type StorageConfig interface {
	IsStorageConfig()
}

type OtherStorageConfig struct {
	Empty *string `json:"_empty,omitempty"`
}

func (OtherStorageConfig) IsStorageConfig() {}

type S3StorageConfig struct {
	Bucket   string `json:"bucket"`
	Endpoint string `json:"endpoint"`
	Access   string `json:"access"`
	Secret   string `json:"secret"`
}

func (S3StorageConfig) IsStorageConfig() {}

type StorageDefinition struct {
	Identifier string        `json:"identifier"`
	Config     StorageConfig `json:"config"`
	IsEnabled  bool          `json:"isEnabled"`
	Priority   int           `json:"priority"`
}
