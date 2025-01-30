package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TurnStileClient struct {
	secretKey string
}

func NewTurnStileClient(secretKey string) *TurnStileClient {
	return &TurnStileClient{secretKey: secretKey}
}

func (t *TurnStileClient) VerifyCaptcha(ctx context.Context, token string, action string) (bool, error) {
	formData := url.Values{}
	formData.Set("secret", t.secretKey)
	formData.Set("response", token)
	if action != "" {
		formData.Set("action", action)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://challenges.cloudflare.com/turnstile/v0/siteverify", strings.NewReader(formData.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var v struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
		Action     string   `json:"action"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return false, err
	}
	if !v.Success {
		return false, fmt.Errorf("verification failed: %v", v.ErrorCodes)
	}
	if action != "" && v.Action != action {
		return false, fmt.Errorf("unexpected action returned: %s", v.Action)
	}
	return true, nil
}
