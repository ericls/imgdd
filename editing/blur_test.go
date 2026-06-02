package editing

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// makeBlurTestPNG creates a solid-color PNG of the given size.
func makeBlurTestPNG(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// --- NewBlurChange validation ---

func TestNewBlurChange_ValidParams(t *testing.T) {
	_, err := NewBlurChange(BlurParams{
		Region: BlurRegion{X1: 0.1, Y1: 0.1, X2: 0.9, Y2: 0.9},
		Radius: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewBlurChange_CoordOutOfRange(t *testing.T) {
	cases := []BlurRegion{
		{X1: -0.1, Y1: 0, X2: 0.5, Y2: 0.5},
		{X1: 0, Y1: 0, X2: 1.1, Y2: 0.5},
		{X1: 0, Y1: -0.1, X2: 0.5, Y2: 0.5},
		{X1: 0, Y1: 0, X2: 0.5, Y2: 1.1},
	}
	for _, r := range cases {
		if _, err := NewBlurChange(BlurParams{Region: r, Radius: 5}); err == nil {
			t.Errorf("expected error for region %+v", r)
		}
	}
}

func TestNewBlurChange_ZeroArea(t *testing.T) {
	cases := []BlurRegion{
		{X1: 0.5, Y1: 0.1, X2: 0.5, Y2: 0.9}, // x1 == x2
		{X1: 0.1, Y1: 0.5, X2: 0.9, Y2: 0.5}, // y1 == y2
		{X1: 0.8, Y1: 0.1, X2: 0.2, Y2: 0.9}, // x1 > x2
	}
	for _, r := range cases {
		if _, err := NewBlurChange(BlurParams{Region: r, Radius: 5}); err == nil {
			t.Errorf("expected error for zero-area region %+v", r)
		}
	}
}

func TestNewBlurChange_RadiusOutOfRange(t *testing.T) {
	validRegion := BlurRegion{X1: 0.1, Y1: 0.1, X2: 0.9, Y2: 0.9}
	for _, r := range []int{0, -1, 101} {
		if _, err := NewBlurChange(BlurParams{Region: validRegion, Radius: r}); err == nil {
			t.Errorf("expected error for radius %d", r)
		}
	}
}

// --- BlurEditor.Apply ---

func TestBlurApply_PreservesDimensions(t *testing.T) {
	base := makeBlurTestPNG(200, 150, color.RGBA{255, 0, 0, 255})
	c, _ := NewBlurChange(BlurParams{
		Region: BlurRegion{X1: 0.2, Y1: 0.2, X2: 0.8, Y2: 0.8},
		Radius: 5,
	})
	editor := &BlurEditor{}
	result, mime, err := editor.Apply(base, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if mime != "image/png" {
		t.Fatalf("expected image/png, got %s", mime)
	}
	img, err := png.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatal(err)
	}
	b := img.Bounds()
	if b.Dx() != 200 || b.Dy() != 150 {
		t.Fatalf("expected 200x150, got %dx%d", b.Dx(), b.Dy())
	}
}

func TestBlurApply_OnlyBlursRegion(t *testing.T) {
	// Image: left half red, right half blue.
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			if x < 50 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255})
			}
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)

	// Blur only the right half.
	c, _ := NewBlurChange(BlurParams{
		Region: BlurRegion{X1: 0.5, Y1: 0, X2: 1.0, Y2: 1.0},
		Radius: 5,
	})
	editor := &BlurEditor{}
	result, _, err := editor.Apply(buf.Bytes(), c, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, _ := png.Decode(bytes.NewReader(result))
	// Left half pixel should still be pure red.
	r, g, b, _ := out.At(10, 50).RGBA()
	if r != 0xffff || g != 0 || b != 0 {
		t.Fatalf("expected unblurred red pixel at (10,50), got r=%d g=%d b=%d", r>>8, g>>8, b>>8)
	}
}

func TestBlurApply_InvalidBase(t *testing.T) {
	c, _ := NewBlurChange(BlurParams{
		Region: BlurRegion{X1: 0.1, Y1: 0.1, X2: 0.9, Y2: 0.9},
		Radius: 5,
	})
	editor := &BlurEditor{}
	if _, _, err := editor.Apply([]byte("not an image"), c, nil); err == nil {
		t.Fatal("expected error for invalid base image")
	}
}

func TestBlurApply_InvalidJSON(t *testing.T) {
	base := makeBlurTestPNG(100, 100, color.White)
	editor := &BlurEditor{}
	if _, _, err := editor.Apply(base, Change{Type: "blur", Params: []byte("not json")}, nil); err == nil {
		t.Fatal("expected error for invalid params JSON")
	}
}

func TestBlurApply_RadiusIndependence(t *testing.T) {
	// Both radius=1 and radius=50 should complete without error on a real image.
	base := makeBlurTestPNG(300, 300, color.RGBA{128, 64, 200, 255})
	editor := &BlurEditor{}
	for _, radius := range []int{1, 50, 100} {
		c, _ := NewBlurChange(BlurParams{
			Region: BlurRegion{X1: 0.1, Y1: 0.1, X2: 0.9, Y2: 0.9},
			Radius: radius,
		})
		result, _, err := editor.Apply(base, c, nil)
		if err != nil {
			t.Fatalf("radius %d: unexpected error: %v", radius, err)
		}
		if len(result) == 0 {
			t.Fatalf("radius %d: empty result", radius)
		}
	}
}
