package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image/jpeg"
	"image/png"

	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
	pngstructure "github.com/dsoprea/go-png-image-structure"
)

func MaybeRemoveExif(data []byte) ([]byte, error) {

	const (
		JpegMediaType = "jpeg"
		PngMediaType  = "png"
		StartBytes    = 0
		EndBytes      = 0
	)

	jmp := jpegstructure.NewJpegMediaParser()
	pmp := pngstructure.NewPngMediaParser()
	filtered := []byte{}

	if jmp.LooksLikeFormat(data) {

		sl, err := jmp.ParseBytes(data)
		if err != nil {
			return nil, err
		}

		_, rawExif, err := sl.Exif()
		if err != nil {
			return data, nil
		}

		startExifBytes := StartBytes
		endExifBytes := EndBytes

		if bytes.Contains(data, rawExif) {
			for i := 0; i < len(data)-len(rawExif); i++ {
				if bytes.Equal(data[i:i+len(rawExif)], rawExif) {
					startExifBytes = i
					endExifBytes = i + len(rawExif)
					break
				}
			}
			fill := make([]byte, len(data[startExifBytes:endExifBytes]))
			copy(data[startExifBytes:endExifBytes], fill)
		}

		filtered = data

		_, err = jpeg.Decode(bytes.NewReader(filtered))
		if err != nil {
			return nil, errors.New("EXIF removal corrupted " + err.Error())
		}

	} else if pmp.LooksLikeFormat(data) {

		cs, err := pmp.ParseBytes(data)
		if err != nil {
			return nil, err
		}

		_, rawExif, err := cs.Exif()
		if err != nil || len(rawExif) == 0 {
			return data, nil
		}

		startExifBytes := StartBytes
		endExifBytes := EndBytes

		if bytes.Contains(data, rawExif) {
			for i := 0; i < len(data)-len(rawExif); i++ {
				if bytes.Equal(data[i:i+len(rawExif)], rawExif) {
					startExifBytes = i
					endExifBytes = i + len(rawExif)
					break
				}
			}
			fill := make([]byte, len(data[startExifBytes:endExifBytes]))
			copy(data[startExifBytes:endExifBytes], fill)
		}

		filtered = data

		chunks := readPNGChunks(bytes.NewReader(filtered))

		for _, chunk := range chunks {
			if !chunk.CRCIsValid() {
				offset := int(chunk.Offset) + 8 + int(chunk.Length)
				crc := chunk.CalculateCRC()

				buf := new(bytes.Buffer)
				binary.Write(buf, binary.BigEndian, crc)
				crcBytes := buf.Bytes()

				copy(filtered[offset:], crcBytes)
			}
		}

		chunks = readPNGChunks(bytes.NewReader(filtered))
		for _, chunk := range chunks {
			if !chunk.CRCIsValid() {
				return nil, errors.New("EXIF removal failed CRC")
			}
		}

		_, err = png.Decode(bytes.NewReader(filtered))
		if err != nil {
			return nil, errors.New("EXIF removal corrupted " + err.Error())
		}
	}

	return filtered, errors.New("not png or jpg")
}
