package editing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/ericls/imgdd/utils"
)

func init() {
	Register("blur", &BlurEditor{})
}

type BlurRegion struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

type BlurParams struct {
	Region BlurRegion `json:"region"`
	Radius int        `json:"radius"`
}

func NewBlurChange(params BlurParams) (Change, error) {
	if params.Region.X1 < 0 || params.Region.X1 > 1 ||
		params.Region.Y1 < 0 || params.Region.Y1 > 1 ||
		params.Region.X2 < 0 || params.Region.X2 > 1 ||
		params.Region.Y2 < 0 || params.Region.Y2 > 1 {
		return Change{}, fmt.Errorf("region coordinates must be between 0 and 1")
	}
	if params.Region.X1 >= params.Region.X2 || params.Region.Y1 >= params.Region.Y2 {
		return Change{}, fmt.Errorf("region must have positive area (x1<x2, y1<y2)")
	}
	if params.Radius < 1 || params.Radius > 100 {
		return Change{}, fmt.Errorf("radius must be between 1 and 100")
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return Change{}, fmt.Errorf("failed to serialize params: %w", err)
	}
	return Change{
		Type:   "blur",
		Params: paramsJSON,
	}, nil
}

type BlurEditor struct{}

func (e *BlurEditor) Apply(baseBytes []byte, cs Change, _ FetchImageFunc) ([]byte, string, error) {
	var params BlurParams
	if err := json.Unmarshal(cs.Params, &params); err != nil {
		return nil, "", fmt.Errorf("invalid blur params: %w", err)
	}
	if params.Region.X1 >= params.Region.X2 || params.Region.Y1 >= params.Region.Y2 {
		return nil, "", fmt.Errorf("region must have positive area")
	}
	if params.Radius < 1 || params.Radius > 100 {
		return nil, "", fmt.Errorf("radius must be between 1 and 100")
	}

	baseMime := utils.DetectMIMEType(&baseBytes)

	baseImg, _, err := image.Decode(bytes.NewReader(baseBytes))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base image: %w", err)
	}

	bounds := baseImg.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Convert normalized region to pixel coords
	rx1 := int(params.Region.X1 * float64(w))
	ry1 := int(params.Region.Y1 * float64(h))
	rx2 := int(params.Region.X2 * float64(w))
	ry2 := int(params.Region.Y2 * float64(h))

	// Clamp to image bounds
	if rx1 < 0 {
		rx1 = 0
	}
	if ry1 < 0 {
		ry1 = 0
	}
	if rx2 > w {
		rx2 = w
	}
	if ry2 > h {
		ry2 = h
	}

	result := applyBoxBlurRegion(baseImg, rx1, ry1, rx2, ry2, params.Radius)

	var buf bytes.Buffer
	switch baseMime {
	case "image/png":
		err = png.Encode(&buf, result)
	case "image/jpeg":
		err = jpeg.Encode(&buf, result, &jpeg.Options{Quality: 95})
	case "image/gif":
		err = gif.Encode(&buf, result, nil)
	default:
		err = png.Encode(&buf, result)
		baseMime = "image/png"
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode result: %w", err)
	}
	return buf.Bytes(), baseMime, nil
}

// applyBoxBlurRegion copies the full image into an RGBA buffer and applies a
// box blur (iterated 3× for a Gaussian-like effect) only to the specified
// pixel rectangle.
func applyBoxBlurRegion(src image.Image, rx1, ry1, rx2, ry2, radius int) *image.RGBA {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Copy source into an RGBA buffer
	rgba := image.NewRGBA(bounds)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			rgba.Set(bounds.Min.X+x, bounds.Min.Y+y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}

	// Extract the region pixels, blur, write back – three passes for smoother result
	for pass := 0; pass < 3; pass++ {
		blurred := boxBlurPass(rgba, bounds, rx1, ry1, rx2, ry2, radius)
		for y := ry1; y < ry2; y++ {
			for x := rx1; x < rx2; x++ {
				p := blurred[(y-ry1)*(rx2-rx1)+(x-rx1)]
				rgba.SetRGBA(bounds.Min.X+x, bounds.Min.Y+y, color.RGBA{R: p.R, G: p.G, B: p.B, A: p.A})
			}
		}
	}

	return rgba
}

type rgbaPixel struct {
	R, G, B, A uint8
}

func boxBlurPass(src *image.RGBA, bounds image.Rectangle, rx1, ry1, rx2, ry2, radius int) []rgbaPixel {
	rw := rx2 - rx1
	rh := ry2 - ry1
	out := make([]rgbaPixel, rw*rh)

	for y := ry1; y < ry2; y++ {
		for x := rx1; x < rx2; x++ {
			var rSum, gSum, bSum, aSum int
			count := 0
			for dy := -radius; dy <= radius; dy++ {
				sy := bounds.Min.Y + y + dy
				if sy < bounds.Min.Y+ry1 {
					sy = bounds.Min.Y + ry1
				}
				if sy >= bounds.Min.Y+ry2 {
					sy = bounds.Min.Y + ry2 - 1
				}
				for dx := -radius; dx <= radius; dx++ {
					sx := bounds.Min.X + x + dx
					if sx < bounds.Min.X+rx1 {
						sx = bounds.Min.X + rx1
					}
					if sx >= bounds.Min.X+rx2 {
						sx = bounds.Min.X + rx2 - 1
					}
					c := src.RGBAAt(sx, sy)
					rSum += int(c.R)
					gSum += int(c.G)
					bSum += int(c.B)
					aSum += int(c.A)
					count++
				}
			}
			out[(y-ry1)*rw+(x-rx1)] = rgbaPixel{
				R: uint8(rSum / count),
				G: uint8(gSum / count),
				B: uint8(bSum / count),
				A: uint8(aSum / count),
			}
		}
	}
	return out
}
