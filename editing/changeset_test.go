package editing

import (
	"encoding/json"
	"testing"
)

func TestChangeSetSerialize(t *testing.T) {
	params := WatermarkParams{
		OverlayImageID: "abc-123",
		Position:       WatermarkPosition{X: 0.9, Y: 0.9},
		Anchor:         AnchorBottomRight,
		Opacity:        0.5,
		Scale:          0.15,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	cs := ChangeSet{
		Type:   "watermark",
		Params: paramsJSON,
	}

	data, err := json.Marshal(cs)
	if err != nil {
		t.Fatal(err)
	}

	var decoded ChangeSet
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Type != "watermark" {
		t.Fatalf("expected type 'watermark', got '%s'", decoded.Type)
	}

	var decodedParams WatermarkParams
	if err := json.Unmarshal(decoded.Params, &decodedParams); err != nil {
		t.Fatal(err)
	}
	if decodedParams.OverlayImageID != "abc-123" {
		t.Fatalf("expected overlay_image_id 'abc-123', got '%s'", decodedParams.OverlayImageID)
	}
	if decodedParams.Opacity != 0.5 {
		t.Fatalf("expected opacity 0.5, got %f", decodedParams.Opacity)
	}
}

func TestRegistryGetEditor(t *testing.T) {
	editor, err := GetEditor("watermark")
	if err != nil {
		t.Fatal(err)
	}
	if editor == nil {
		t.Fatal("expected non-nil editor for 'watermark'")
	}
}

func TestRegistryGetUnknownEditor(t *testing.T) {
	_, err := GetEditor("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown editor type")
	}
}
