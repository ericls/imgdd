package identity

import (
	"os"
	"testing"

	"github.com/ericls/imgdd/test_support"
)

var TestServiceMan = test_support.NewTestExternalServiceManager()

func TestMain(m *testing.M) {

	TestServiceMan.StartPostgres()
	code := m.Run()
	TestServiceMan.Purge()
	os.Exit(code)
}
