package model

import (
	"encoding/json"
	"fmt"

	"github.com/ericls/imgdd/domainmodels"
)

// TODO: Maybe move the enum to the storage package
// so that it can be shared with the storage backend
type StorageTypeEnum string

const (
	StorageType_S3    StorageTypeEnum = "s3"
	StorageType_FS    StorageTypeEnum = "fs"
	StorageType_Other StorageTypeEnum = "other"
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

type FSStorageConfig struct {
	MediaRoot string `json:"mediaRoot"`
}

func (FSStorageConfig) IsStorageConfig() {}

type StorageDefinition struct {
	Id         string        `json:"id"`
	Identifier string        `json:"identifier"`
	Config     StorageConfig `json:"config"`
	IsEnabled  bool          `json:"isEnabled"`
	Priority   int           `json:"priority"`
}

type CreateStorageDefinitionInput struct {
	Identifier  string          `json:"identifier"`
	StorageType StorageTypeEnum `json:"storageType"`
	ConfigJSON  string          `json:"configJSON"`
	IsEnabled   bool            `json:"isEnabled"`
	Priority    int             `json:"priority"`
}

type UpdateStorageDefinitionInput struct {
	Identifier string  `json:"identifier"`
	ConfigJSON *string `json:"configJSON,omitempty"`
	IsEnabled  *bool   `json:"isEnabled,omitempty"`
	Priority   *int    `json:"priority,omitempty"`
}

type StorageDefinitionConnectivityResult struct {
	Ok    bool    `json:"ok"`
	Error *string `json:"error,omitempty"`
}

type CheckStorageDefinitionConnectivityInput struct {
	ID string `json:"id"`
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
		storageConfig = conf
	} else if storageType == StorageType_FS {
		fmt.Printf("sd.Config: %s\n", sd.Config)
		var conf FSStorageConfig
		err := json.Unmarshal([]byte(sd.Config), &conf)
		if err != nil {
			return nil, err
		}
		storageConfig = conf
	} else {
		storageConfig = OtherStorageConfig{}
	}
	return &StorageDefinition{
		Id:         sd.Id,
		Identifier: sd.Identifier,
		IsEnabled:  sd.IsEnabled,
		Priority:   int(sd.Priority),
		Config:     storageConfig,
	}, nil
}
