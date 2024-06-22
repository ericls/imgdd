package model

import (
	"encoding/json"
	"imgdd/domainmodels"
)

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

func FromStorageDefinition(sd *domainmodels.StorageDefinition) (*StorageDefinition, error) {
	storageType := StorageTypeEnum(sd.StorageType)
	var storageConfig StorageConfig
	if storageType == StorageType_S3 {
		var conf S3StorageConfig
		err := json.Unmarshal([]byte(sd.Config), &conf)
		if err != nil {
			return nil, err
		}
		storageConfig = &conf
	} else {
		storageConfig = &OtherStorageConfig{}
	}
	return &StorageDefinition{
		Identifier: sd.Identifier,
		IsEnabled:  sd.IsEnabled,
		Priority:   int(sd.Priority),
		Config:     storageConfig,
	}, nil
}
