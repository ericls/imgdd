package editing

import (
	"encoding/json"
	"fmt"
)

type ChangeSet struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

// FetchImageFunc retrieves image bytes by image ID.
type FetchImageFunc func(id string) ([]byte, error)

// Editor applies a ChangeSet to base image bytes, producing new image bytes.
type Editor interface {
	Apply(base []byte, cs ChangeSet, fetchImage FetchImageFunc) ([]byte, string, error)
}

var registry = map[string]Editor{}

func Register(changeType string, editor Editor) {
	registry[changeType] = editor
}

func GetEditor(changeType string) (Editor, error) {
	e, ok := registry[changeType]
	if !ok {
		return nil, fmt.Errorf("unknown change type: %s", changeType)
	}
	return e, nil
}
