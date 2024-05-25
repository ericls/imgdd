package httpserver

import (
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	r := t.Run()
	os.Exit(r)
}
