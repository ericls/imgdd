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

type RecaptchaClient struct {
	serverKey string
}

func NewRecaptchaClient(secretKey string) *RecaptchaClient {
	return &RecaptchaClient{serverKey: secretKey}
}

func (r *RecaptchaClient) VerifyCaptcha(ctx context.Context, token string, action string) (bool, error) {
	requestBody := url.Values{}
	requestBody.Add("secret", r.serverKey)
	requestBody.Add("response", token)

	googleRecaptchaURL := "https://www.google.com/recaptcha/api/siteverify"
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "POST", googleRecaptchaURL, strings.NewReader(requestBody.Encode()))
	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var v struct {
		Success bool    `json:"success"`
		Score   float64 `json:"score"`
		Action  string  `json:"action"`
	}
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return false, err
	}
	if !v.Success {
		return false, fmt.Errorf("verification failed")
	}
	if action != "" && v.Action != action {
		return false, fmt.Errorf("unexpected action returned: %s", v.Action)
	}
	if v.Score < 0.5 {
		return false, fmt.Errorf("score is too low: %f", v.Score)
	}
	return true, nil
}
