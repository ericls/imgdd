package identity

import (
	"imgdd/buildflag"

	"golang.org/x/crypto/bcrypt"
)

const cost = 14

func HashPassword(password string) (string, error) {
	realCost := cost
	if buildflag.Debug == "true" {
		realCost = 1
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), realCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
