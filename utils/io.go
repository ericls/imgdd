package utils

import "io"

type SeekerReader interface {
	io.Seeker
	io.Reader
}
