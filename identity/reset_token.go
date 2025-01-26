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

const resetPasswordSalt = "imgdd-reset-password"
const resetPasswordTokenValidity = 15 * time.Minute

func makeResetTokenDigest(userId string, currentHashedPassword string, timeInfo time.Time, secret string) string {
	key := fmt.Sprintf("%s.%s", resetPasswordSalt, secret)
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(userId))
	hasher.Write([]byte(currentHashedPassword))
	hasher.Write([]byte(timeInfo.Format(time.RFC3339)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func makeResetPasswordToken(userId string, currentHashedPassword string, currentTime time.Time, secret string) string {
	ts := currentTime.UnixMilli()
	tsB36 := strconv.FormatInt(ts, 36)
	hashedUserWithTs := makeResetTokenDigest(userId, currentHashedPassword, currentTime, secret)
	return fmt.Sprintf("%s-%s", tsB36, hashedUserWithTs)
}

func checkResetPasswordToken(userId string, currentHashedPassword string, currentTime time.Time, secret string, token string) bool {
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
	d1, err1 := hex.DecodeString(digest1)
	d2, err2 := hex.DecodeString(digest2)
	if err1 != nil || err2 != nil {
		return false
	}
	if !hmac.Equal([]byte(d1), []byte(d2)) {
		return false
	}
	return currentTime.Sub(timeThen) <= resetPasswordTokenValidity
}
