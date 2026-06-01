package editing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// makeTestPNG creates a solid-color PNG image of the given dimensions.
func makeTestPNG(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func makeParams(overlayID string, x, y float64, anchor Anchor, opacity, scale float64) ChangeSet {
	params := WatermarkParams{
		OverlayImageID: overlayID,
		Position:       WatermarkPosition{X: x, Y: y},
		Anchor:         anchor,
		Opacity:        opacity,
		Scale:          scale,
	}
	paramsJSON, _ := json.Marshal(params)
	return ChangeSet{Type: "watermark", Params: paramsJSON}
}

func TestWatermarkBasic(t *testing.T) {
	base := makeTestPNG(200, 200, color.RGBA{255, 0, 0, 255})
	overlay := makeTestPNG(50, 50, color.RGBA{0, 0, 255, 255})

	cs := makeParams("overlay-1", 0.5, 0.5, AnchorCenter, 1.0, 0.25)

	editor := &WatermarkEditor{}
	result, mime, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		if id == "overlay-1" {
			return overlay, nil
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if mime != "image/png" {
		t.Fatalf("expected image/png, got %s", mime)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}

	// Decode result and verify dimensions match base
	resultImg, err := png.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatal(err)
	}
	bounds := resultImg.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 200 {
		t.Fatalf("expected 200x200, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// The center pixel should not be pure red anymore (overlay is blue)
	r, g, b, _ := resultImg.At(100, 100).RGBA()
	if r == 0xffff && g == 0 && b == 0 {
		t.Fatal("expected center pixel to be modified by overlay")
	}
}

func TestWatermarkZeroOpacity(t *testing.T) {
	base := makeTestPNG(100, 100, color.RGBA{255, 0, 0, 255})
	overlay := makeTestPNG(50, 50, color.RGBA{0, 0, 255, 255})

	cs := makeParams("overlay-1", 0.5, 0.5, AnchorCenter, 0.0, 0.5)

	editor := &WatermarkEditor{}
	result, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// With zero opacity, the image should be unchanged — center should still be red
	resultImg, _ := png.Decode(bytes.NewReader(result))
	r, g, b, _ := resultImg.At(50, 50).RGBA()
	if r != 0xffff || g != 0 || b != 0 {
		t.Fatalf("expected red pixel at center with zero opacity, got r=%d g=%d b=%d", r>>8, g>>8, b>>8)
	}
}

func TestWatermarkAnchors(t *testing.T) {
	tests := []struct {
		anchor   Anchor
		expectDx int
		expectDy int
	}{
		{AnchorTopLeft, 0, 0},
		{AnchorTopRight, -10, 0},
		{AnchorBottomLeft, 0, -10},
		{AnchorBottomRight, -10, -10},
		{AnchorCenter, -5, -5},
	}
	for _, tc := range tests {
		dx, dy := anchorOffset(tc.anchor, 10, 10)
		if dx != tc.expectDx || dy != tc.expectDy {
			t.Errorf("anchor %s: expected (%d,%d), got (%d,%d)", tc.anchor, tc.expectDx, tc.expectDy, dx, dy)
		}
	}
}

func TestWatermarkInvalidOpacity(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)
	overlay := makeTestPNG(10, 10, color.Black)

	cs := makeParams("o", 0.5, 0.5, AnchorCenter, 1.5, 0.1)
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err == nil {
		t.Fatal("expected error for opacity > 1")
	}
}

func TestWatermarkInvalidScale(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)
	overlay := makeTestPNG(10, 10, color.Black)

	cs := makeParams("o", 0.5, 0.5, AnchorCenter, 0.5, 0.0)
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err == nil {
		t.Fatal("expected error for scale == 0")
	}
}

func TestWatermarkOverlayLargerThanBase(t *testing.T) {
	base := makeTestPNG(50, 50, color.RGBA{255, 0, 0, 255})
	overlay := makeTestPNG(500, 500, color.RGBA{0, 255, 0, 255})

	cs := makeParams("big", 0.5, 0.5, AnchorCenter, 0.8, 0.5)
	editor := &WatermarkEditor{}
	result, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	resultImg, _ := png.Decode(bytes.NewReader(result))
	bounds := resultImg.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Fatalf("result should match base dimensions: expected 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestWatermarkCornerPositions(t *testing.T) {
	base := makeTestPNG(200, 200, color.RGBA{255, 0, 0, 255})
	overlay := makeTestPNG(20, 20, color.RGBA{0, 0, 255, 255})

	// Place at bottom-right corner with bottom-right anchor
	cs := makeParams("o", 1.0, 1.0, AnchorBottomRight, 1.0, 0.1)
	editor := &WatermarkEditor{}
	result, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	resultImg, _ := png.Decode(bytes.NewReader(result))
	// The pixel just inside the bottom-right corner should be blue-ish
	r, _, b, _ := resultImg.At(199, 199).RGBA()
	if b == 0 && r == 0xffff {
		t.Fatal("expected bottom-right corner to have overlay color")
	}
}

func TestWatermarkNegativeOpacity(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)
	overlay := makeTestPNG(10, 10, color.Black)

	cs := makeParams("o", 0.5, 0.5, AnchorCenter, -0.5, 0.1)
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err == nil {
		t.Fatal("expected error for negative opacity")
	}
}

func TestWatermarkInvalidBaseImage(t *testing.T) {
	base := []byte("not a valid image")
	overlay := makeTestPNG(10, 10, color.Black)

	cs := makeParams("o", 0.5, 0.5, AnchorCenter, 0.5, 0.1)
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err == nil {
		t.Fatal("expected error for invalid base image")
	}
}

func TestWatermarkInvalidOverlayImage(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)

	cs := makeParams("bad", 0.5, 0.5, AnchorCenter, 0.5, 0.1)
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return []byte("not a valid image"), nil
	})
	if err == nil {
		t.Fatal("expected error for invalid overlay image")
	}
}

func TestWatermarkFetchFailure(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)

	cs := makeParams("missing", 0.5, 0.5, AnchorCenter, 0.5, 0.1)
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return nil, fmt.Errorf("image not found")
	})
	if err == nil {
		t.Fatal("expected error when fetch fails")
	}
}

func TestWatermarkInvalidJSON(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)

	cs := ChangeSet{Type: "watermark", Params: []byte("not json")}
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return base, nil
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON params")
	}
}

func TestWatermarkScaleGreaterThanOne(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)
	overlay := makeTestPNG(10, 10, color.Black)

	cs := makeParams("o", 0.5, 0.5, AnchorCenter, 0.5, 1.5)
	editor := &WatermarkEditor{}
	_, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err == nil {
		t.Fatal("expected error for scale > 1")
	}
}

func TestWatermarkExtremeScale(t *testing.T) {
	base := makeTestPNG(100, 100, color.White)
	overlay := makeTestPNG(10, 10, color.RGBA{0, 0, 255, 255})

	// Scale at maximum (1.0) should still work
	cs := makeParams("o", 0.5, 0.5, AnchorCenter, 0.5, 1.0)
	editor := &WatermarkEditor{}
	result, _, err := editor.Apply(base, cs, func(id string) ([]byte, error) {
		return overlay, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}
}
