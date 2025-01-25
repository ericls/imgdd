package identity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const resetPasswordSaltPrefix = "imgdd-reset-password"
const resetPasswordTokenValidity = 15 * time.Minute

func makeResetTokenDigest(userId string, currentHashedPassword string, timeInfo time.Time, secret string) string {
	hasher := sha256.New()
	salt := fmt.Sprintf("%s-%s-%s", resetPasswordSaltPrefix, secret, timeInfo.Format(time.RFC3339))
	hasher.Write([]byte(userId))
	hasher.Write([]byte(currentHashedPassword))
	hasher.Write([]byte(salt))
	return hex.EncodeToString(hasher.Sum(nil))
}

func makeToken(userId string, currentHashedPassord string, currentTime time.Time, secret string) string {
	ts := currentTime.UnixMilli()
	tsB36 := strconv.FormatInt(ts, 36)
	hashedUserWithTs := makeResetTokenDigest(userId, currentHashedPassord, currentTime, secret)
	return fmt.Sprintf("%s-%s", tsB36, hashedUserWithTs)
}

func checkToken(userId string, currentHashedPassword string, currentTime time.Time, secret string, token string) bool {
	parts := strings.SplitN(token, "-", 2)
	if len(parts) != 2 {
		return false
	}
	ts := parts[0]
	tsInt, err := strconv.ParseInt(ts, 36, 64)
	if err != nil {
		return false
	}
	timeThen := time.Unix(0, tsInt*int64(time.Millisecond))
	digest1 := makeResetTokenDigest(userId, currentHashedPassword, timeThen, secret)
	digest2 := parts[1]
	if !hmac.Equal([]byte(digest1), []byte(digest2)) {
		return false
	}
	return currentTime.Sub(timeThen) <= resetPasswordTokenValidity
}
