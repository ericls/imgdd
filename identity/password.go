package identity

import (
	"github.com/ericls/imgdd/buildflag"

	"golang.org/x/crypto/bcrypt"
)

const cost = 14

func HashPassword(password string) (string, error) {
	realCost := cost
	if buildflag.IsDebug {
		realCost = bcrypt.MinCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), realCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckPasswordByUserId(userId, password string, identityRepo IdentityRepo) bool {
	hashedPassword := identityRepo.GetUserPassword(userId)
	return CheckPasswordHash(password, hashedPassword)
}
