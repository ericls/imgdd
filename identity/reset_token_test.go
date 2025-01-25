package identity

import (
	"testing"
	"time"
)

func TestMakeAndCheckToken(t *testing.T) {
	secret := "123"
	userId := "user"
	currentHashedPassword := "password"
	then := time.Now()
	token := makeToken(userId, currentHashedPassword, then, secret)
	now := then.Add(resetPasswordTokenValidity / 2)
	if !checkToken(userId, currentHashedPassword, now, secret, token) {
		t.Errorf("token should be valid")
	}
	now = then.Add(resetPasswordTokenValidity + 1)
	if checkToken(userId, currentHashedPassword, now, secret, token) {
		t.Errorf("token should be invalid, expired")
	}
}
