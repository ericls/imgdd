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

// ApplyResult holds the output of applying a ChangeSet.
type ApplyResult struct {
	Bytes       []byte
	MIMEType    string
	ChangesJSON []byte
}

// ApplyChangeSet orchestrates applying a change set: looks up the editor,
// fetches the base image bytes, applies the edit, and serializes the changes.
func ApplyChangeSet(cs ChangeSet, baseImageId string, fetchImage FetchImageFunc) (*ApplyResult, error) {
	editor, err := GetEditor(cs.Type)
	if err != nil {
		return nil, err
	}

	baseBytes, err := fetchImage(baseImageId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base image: %w", err)
	}

	resultBytes, resultMime, err := editor.Apply(baseBytes, cs, fetchImage)
	if err != nil {
		return nil, fmt.Errorf("failed to apply %s: %w", cs.Type, err)
	}

	changesJSON, err := json.Marshal(cs)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize changes: %w", err)
	}

	return &ApplyResult{
		Bytes:       resultBytes,
		MIMEType:    resultMime,
		ChangesJSON: changesJSON,
	}, nil
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
