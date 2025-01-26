package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func decodeBase64(data string) ([]byte, error) {
	missingPadding := len(data) % 4
	if missingPadding > 0 {
		data += strings.Repeat("=", 4-missingPadding)
	}
	return base64.URLEncoding.DecodeString(data)
}

func signMessage(message string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(signature)
}

func Dumps(data interface{}, key string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	signature := signMessage(string(jsonData), key)
	encodedData := base64.URLEncoding.EncodeToString(jsonData)
	return fmt.Sprintf("%s.%s", strings.TrimRight(encodedData, "="), strings.TrimRight(signature, "=")), nil
}

func Loads(token string, target interface{}, key string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return errors.New("invalid token format")
	}

	jsonData, err := decodeBase64(parts[0])
	if err != nil {
		return fmt.Errorf("failed to decode JSON data: %w", err)
	}

	providedSignature := parts[1]
	expectedSignature := signMessage(string(jsonData), key)
	decodedProvidedSignature, err1 := decodeBase64(providedSignature)
	decodedExpectedSignature, err2 := decodeBase64(expectedSignature)
	if err1 != nil || err2 != nil {
		return errors.New("invalid token signature")
	}
	if !hmac.Equal(decodedProvidedSignature, decodedExpectedSignature) {
		return errors.New("invalid token signature")
	}

	err = json.Unmarshal(jsonData, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
