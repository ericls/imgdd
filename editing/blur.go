package editing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
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
	// Truncation of small normalized coords can collapse the region to 0×0;
	// expand to at least 1px so the blur is always applied.
	if rx2 <= rx1 {
		if rx1 > 0 {
			rx1--
		} else {
			rx2++
		}
	}
	if ry2 <= ry1 {
		if ry1 > 0 {
			ry1--
		} else {
			ry2++
		}
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
// separable box blur (3 passes for a Gaussian-like effect) to the specified
// pixel rectangle. Each pass is O(region_pixels) regardless of radius.
func applyBoxBlurRegion(src image.Image, rx1, ry1, rx2, ry2, radius int) *image.RGBA {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	rgba := image.NewRGBA(bounds)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			rgba.Set(bounds.Min.X+x, bounds.Min.Y+y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}

	rw := rx2 - rx1
	rh := ry2 - ry1
	kw := 2*radius + 1
	// tmp holds RGBA bytes for the blur region after the horizontal pass.
	tmp := make([]byte, rw*rh*4)

	for range 3 {
		hBoxBlur(rgba, tmp, rx1, ry1, rx2, ry2, rw, radius, kw)
		vBoxBlur(rgba, tmp, rx1, ry1, rx2, ry2, rw, radius, kw)
	}

	return rgba
}

// hBoxBlur performs a horizontal sliding-window box blur over the region
// [rx1,rx2) × [ry1,ry2), reading from rgba.Pix and writing to tmp.
// Samples outside the region are clamped to the nearest edge pixel.
// rx1,ry1,rx2,ry2 are relative (0-based) pixel coordinates within the image.
func hBoxBlur(rgba *image.RGBA, tmp []byte, rx1, ry1, rx2, ry2, rw, radius, kw int) {
	for y := ry1; y < ry2; y++ {
		var r0, g0, b0, a0 int

		// Build initial sum for x = rx1 by summing the full kernel window.
		for k := -radius; k <= radius; k++ {
			sx := rx1 + k
			if sx < rx1 {
				sx = rx1
			} else if sx >= rx2 {
				sx = rx2 - 1
			}
			off := rgba.PixOffset(sx+rgba.Rect.Min.X, y+rgba.Rect.Min.Y)
			r0 += int(rgba.Pix[off])
			g0 += int(rgba.Pix[off+1])
			b0 += int(rgba.Pix[off+2])
			a0 += int(rgba.Pix[off+3])
		}
		tidx := (y - ry1) * rw * 4
		tmp[tidx] = byte(r0 / kw)
		tmp[tidx+1] = byte(g0 / kw)
		tmp[tidx+2] = byte(b0 / kw)
		tmp[tidx+3] = byte(a0 / kw)

		for x := rx1 + 1; x < rx2; x++ {
			// Slide: remove the pixel that left the window (clamped to region).
			lx := x - radius - 1
			if lx < rx1 {
				lx = rx1
			} else if lx >= rx2 {
				lx = rx2 - 1
			}
			loff := rgba.PixOffset(lx+rgba.Rect.Min.X, y+rgba.Rect.Min.Y)
			r0 -= int(rgba.Pix[loff])
			g0 -= int(rgba.Pix[loff+1])
			b0 -= int(rgba.Pix[loff+2])
			a0 -= int(rgba.Pix[loff+3])

			// Add the pixel that entered the window (clamped to region).
			nx := x + radius
			if nx < rx1 {
				nx = rx1
			} else if nx >= rx2 {
				nx = rx2 - 1
			}
			noff := rgba.PixOffset(nx+rgba.Rect.Min.X, y+rgba.Rect.Min.Y)
			r0 += int(rgba.Pix[noff])
			g0 += int(rgba.Pix[noff+1])
			b0 += int(rgba.Pix[noff+2])
			a0 += int(rgba.Pix[noff+3])

			tidx = ((y-ry1)*rw + (x - rx1)) * 4
			tmp[tidx] = byte(r0 / kw)
			tmp[tidx+1] = byte(g0 / kw)
			tmp[tidx+2] = byte(b0 / kw)
			tmp[tidx+3] = byte(a0 / kw)
		}
	}
}

// vBoxBlur performs a vertical sliding-window box blur over the region,
// reading from tmp (horizontal-pass output) and writing back to rgba.Pix.
func vBoxBlur(rgba *image.RGBA, tmp []byte, rx1, ry1, rx2, ry2, rw, radius, kw int) {
	for x := rx1; x < rx2; x++ {
		xi := x - rx1
		var r0, g0, b0, a0 int

		// Build initial sum for y = ry1.
		for k := -radius; k <= radius; k++ {
			sy := ry1 + k
			if sy < ry1 {
				sy = ry1
			} else if sy >= ry2 {
				sy = ry2 - 1
			}
			tidx := ((sy-ry1)*rw + xi) * 4
			r0 += int(tmp[tidx])
			g0 += int(tmp[tidx+1])
			b0 += int(tmp[tidx+2])
			a0 += int(tmp[tidx+3])
		}
		off := rgba.PixOffset(x+rgba.Rect.Min.X, ry1+rgba.Rect.Min.Y)
		rgba.Pix[off] = byte(r0 / kw)
		rgba.Pix[off+1] = byte(g0 / kw)
		rgba.Pix[off+2] = byte(b0 / kw)
		rgba.Pix[off+3] = byte(a0 / kw)

		for y := ry1 + 1; y < ry2; y++ {
			// Remove pixel that left the window (clamped to region).
			ly := y - radius - 1
			if ly < ry1 {
				ly = ry1
			} else if ly >= ry2 {
				ly = ry2 - 1
			}
			ltidx := ((ly-ry1)*rw + xi) * 4
			r0 -= int(tmp[ltidx])
			g0 -= int(tmp[ltidx+1])
			b0 -= int(tmp[ltidx+2])
			a0 -= int(tmp[ltidx+3])

			// Add pixel that entered the window (clamped to region).
			ny := y + radius
			if ny < ry1 {
				ny = ry1
			} else if ny >= ry2 {
				ny = ry2 - 1
			}
			ntidx := ((ny-ry1)*rw + xi) * 4
			r0 += int(tmp[ntidx])
			g0 += int(tmp[ntidx+1])
			b0 += int(tmp[ntidx+2])
			a0 += int(tmp[ntidx+3])

			off = rgba.PixOffset(x+rgba.Rect.Min.X, y+rgba.Rect.Min.Y)
			rgba.Pix[off] = byte(r0 / kw)
			rgba.Pix[off+1] = byte(g0 / kw)
			rgba.Pix[off+2] = byte(b0 / kw)
			rgba.Pix[off+3] = byte(a0 / kw)
		}
	}
}
