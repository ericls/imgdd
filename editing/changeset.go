package editing

import (
	"encoding/json"
	"fmt"
)

// Change represents a single named edit operation with its parameters.
type Change struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

// ChangeSet is an ordered list of changes applied atomically to produce one image.
type ChangeSet []Change

// FetchImageFunc retrieves image bytes by image ID.
type FetchImageFunc func(id string) ([]byte, error)

// Editor applies a Change to base image bytes, producing new image bytes.
type Editor interface {
	Apply(base []byte, c Change, fetchImage FetchImageFunc) ([]byte, string, error)
}

// ApplyResult holds the output of applying a ChangeSet.
type ApplyResult struct {
	Bytes       []byte
	MIMEType    string
	ChangesJSON []byte
}

// ApplyChangeSet applies each Change in the ChangeSet in order, piping
// output bytes into the next step, then serializes the full ChangeSet.
func ApplyChangeSet(cs ChangeSet, baseImageId string, fetchImage FetchImageFunc) (*ApplyResult, error) {
	currentBytes, err := fetchImage(baseImageId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base image: %w", err)
	}

	var currentMIME string
	for _, c := range cs {
		editor, err := GetEditor(c.Type)
		if err != nil {
			return nil, err
		}
		resultBytes, resultMIME, err := editor.Apply(currentBytes, c, fetchImage)
		if err != nil {
			return nil, fmt.Errorf("failed to apply %s: %w", c.Type, err)
		}
		currentBytes = resultBytes
		currentMIME = resultMIME
	}

	changesJSON, err := json.Marshal(cs)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize changes: %w", err)
	}

	return &ApplyResult{
		Bytes:       currentBytes,
		MIMEType:    currentMIME,
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
