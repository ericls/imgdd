package imagechecks

import (
	"io"

	"github.com/ericls/imgdd/utils"
)

type Checker func(file io.Reader) bool

func CheckAll(checkers []Checker, file utils.SeekerReader) bool {
	for _, check := range checkers {
		file.Seek(0, io.SeekStart)
		if !check(file) {
			return false
		}
	}
	return true
}
