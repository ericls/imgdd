package editing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"

	_ "golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/ericls/imgdd/utils"
)

func init() {
	Register("watermark", &WatermarkEditor{})
}

type Anchor string

const (
	AnchorTopLeft     Anchor = "top_left"
	AnchorTopRight    Anchor = "top_right"
	AnchorBottomLeft  Anchor = "bottom_left"
	AnchorBottomRight Anchor = "bottom_right"
	AnchorCenter      Anchor = "center"
)

type WatermarkPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type WatermarkParams struct {
	OverlayImageID string            `json:"overlay_image_id"`
	Position       WatermarkPosition `json:"position"`
	Anchor         Anchor            `json:"anchor"`
	Opacity        float64           `json:"opacity"`
	Scale          float64           `json:"scale"`
}

// NewWatermarkChange validates params and builds a Change.
func NewWatermarkChange(params WatermarkParams) (Change, error) {
	if params.Opacity < 0 || params.Opacity > 1 {
		return Change{}, fmt.Errorf("opacity must be between 0 and 1")
	}
	if params.Scale <= 0 || params.Scale > 1 {
		return Change{}, fmt.Errorf("scale must be between 0 (exclusive) and 1")
	}
	if params.Position.X < 0 || params.Position.X > 1 || params.Position.Y < 0 || params.Position.Y > 1 {
		return Change{}, fmt.Errorf("position values must be between 0 and 1")
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return Change{}, fmt.Errorf("failed to serialize params: %w", err)
	}
	return Change{
		Type:   "watermark",
		Params: paramsJSON,
	}, nil
}

type WatermarkEditor struct{}

func (e *WatermarkEditor) Apply(baseBytes []byte, cs Change, fetchImage FetchImageFunc) ([]byte, string, error) {
	var params WatermarkParams
	if err := json.Unmarshal(cs.Params, &params); err != nil {
		return nil, "", fmt.Errorf("invalid watermark params: %w", err)
	}

	if params.Opacity < 0 || params.Opacity > 1 {
		return nil, "", fmt.Errorf("opacity must be between 0 and 1")
	}
	if params.Scale <= 0 || params.Scale > 1 {
		return nil, "", fmt.Errorf("scale must be between 0 (exclusive) and 1")
	}

	// Detect base image MIME type
	baseMime := utils.DetectMIMEType(&baseBytes)

	// Decode base image
	baseImg, _, err := image.Decode(bytes.NewReader(baseBytes))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base image: %w", err)
	}

	// Fetch and decode overlay image
	overlayBytes, err := fetchImage(params.OverlayImageID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch overlay image: %w", err)
	}
	overlayImg, _, err := image.Decode(bytes.NewReader(overlayBytes))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode overlay image: %w", err)
	}

	// Scale overlay relative to base image's shorter dimension
	baseBounds := baseImg.Bounds()
	shortSide := baseBounds.Dx()
	if baseBounds.Dy() < shortSide {
		shortSide = baseBounds.Dy()
	}
	targetSize := int(math.Round(float64(shortSide) * params.Scale))
	if targetSize < 1 {
		targetSize = 1
	}

	overlayBounds := overlayImg.Bounds()
	overlayAspect := float64(overlayBounds.Dx()) / float64(overlayBounds.Dy())
	var scaledW, scaledH int
	if overlayAspect >= 1 {
		scaledW = targetSize
		scaledH = int(math.Round(float64(targetSize) / overlayAspect))
	} else {
		scaledH = targetSize
		scaledW = int(math.Round(float64(targetSize) * overlayAspect))
	}
	if scaledW < 1 {
		scaledW = 1
	}
	if scaledH < 1 {
		scaledH = 1
	}

	// Scale overlay
	scaledOverlay := image.NewRGBA(image.Rect(0, 0, scaledW, scaledH))
	xdraw.CatmullRom.Scale(scaledOverlay, scaledOverlay.Bounds(), overlayImg, overlayBounds, xdraw.Over, nil)

	// Apply opacity
	if params.Opacity < 1.0 {
		applyOpacity(scaledOverlay, params.Opacity)
	}

	// Calculate position
	px := int(math.Round(params.Position.X * float64(baseBounds.Dx())))
	py := int(math.Round(params.Position.Y * float64(baseBounds.Dy())))
	offsetX, offsetY := anchorOffset(params.Anchor, scaledW, scaledH)
	drawX := px + offsetX
	drawY := py + offsetY

	// Composite
	result := image.NewRGBA(baseBounds)
	draw.Draw(result, baseBounds, baseImg, baseBounds.Min, draw.Src)
	draw.Draw(result, image.Rect(drawX, drawY, drawX+scaledW, drawY+scaledH), scaledOverlay, image.Point{}, draw.Over)

	// Encode in the same format as the base image
	var buf bytes.Buffer
	switch baseMime {
	case "image/png":
		err = png.Encode(&buf, result)
	case "image/jpeg":
		err = jpeg.Encode(&buf, result, &jpeg.Options{Quality: 95})
	case "image/gif":
		err = gif.Encode(&buf, result, nil)
	default:
		// Fall back to PNG for formats we can't encode (webp, bmp)
		err = png.Encode(&buf, result)
		baseMime = "image/png"
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode result: %w", err)
	}
	return buf.Bytes(), baseMime, nil
}

func applyOpacity(img *image.RGBA, opacity float64) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			newA := float64(a>>8) * opacity
			// Premultiply RGB by the opacity factor for correct alpha compositing
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(r>>8) * opacity),
				G: uint8(float64(g>>8) * opacity),
				B: uint8(float64(b>>8) * opacity),
				A: uint8(newA),
			})
		}
	}
}

func anchorOffset(anchor Anchor, w, h int) (int, int) {
	switch anchor {
	case AnchorTopLeft:
		return 0, 0
	case AnchorTopRight:
		return -w, 0
	case AnchorBottomLeft:
		return 0, -h
	case AnchorBottomRight:
		return -w, -h
	case AnchorCenter:
		return -w / 2, -h / 2
	default:
		return 0, 0
	}
}
