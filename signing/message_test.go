package signing_test

import (
	"testing"

	"github.com/ericls/imgdd/signing"
)

func TestDumpsAndLoadsMessage(t *testing.T) {
	message := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  30,
	}
	key := "123"
	token, err := signing.Dumps(message, key)
	if err != nil {
		t.Fatal(err)
	}
	var loadedMessage struct {
		Name string
		Age  int
	}
	err = signing.Loads(token, &loadedMessage, key)
	if err != nil {
		t.Fatal(err)
	}
	if loadedMessage != message {
		t.Fatalf("expected %v, got %v", message, loadedMessage)
	}
	err = signing.Loads(token+"1", &loadedMessage, key)
	if err == nil {
		t.Fatal("expected token signature error")
	}
}
