//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type ImageTable struct {
	ID              uuid.UUID `sql:"primary_key"`
	CreatedByID     *uuid.UUID
	Name            string
	Identifier      string
	RootID          *uuid.UUID
	ParentID        *uuid.UUID
	Changes         string
	UploaderIP      *string
	MimeType        string
	NominalWidth    int32
	NominalHeight   int32
	NominalByteSize int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
